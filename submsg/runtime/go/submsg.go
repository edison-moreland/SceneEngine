package submsg

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// TODO: Auto detect byte order???
var byteOrder = binary.LittleEndian

var (
	ErrMsgIdUnknown = errors.New("msg id unknown")
)

type MsgId uint32
type MsgSender func(id MsgId, body []byte)
type MsgReceiver func(id MsgId, body []byte) error

type msg struct {
	id   MsgId
	body []byte
}

// Start submsg, a framework for communicating with child processes through stdin/out
func Start(ctx context.Context, bin string, receiver MsgReceiver) (MsgSender, error) {
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stderr = os.Stderr

	sendQueue := make(chan msg)
	childIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	go msgSender(ctx, childIn, sendQueue)

	childOut, err := cmd.StdoutPipe()
	if err != nil {
		close(sendQueue)
		return nil, err
	}
	go msgReceiver(ctx, childOut, receiver)

	if err := cmd.Start(); err != nil {
		close(sendQueue)
		return nil, err
	}

	return func(id MsgId, body []byte) {
		sendQueue <- msg{
			id:   id,
			body: body,
		}
	}, nil
}

func msgSender(ctx context.Context, childIn io.WriteCloser, queue chan msg) {
	for {
		select {
		case <-ctx.Done():
			childIn.Close()
			close(queue)
			return
		case msg := <-queue:
			err := binary.Write(childIn, byteOrder, msg.id)
			if err != nil {
				err = fmt.Errorf("%w: writing msg id to child", err)
				panic(err)
			}

			err = binary.Write(childIn, byteOrder, uint32(len(msg.body)))
			if err != nil {
				err = fmt.Errorf("%w: writing msg length to child", err)
				panic(err)
			}

			if len(msg.body) <= 0 {
				break
			}

			bWritten, err := childIn.Write(msg.body)

			if bWritten != len(msg.body) {
				err = fmt.Errorf("only %d/%d bytes written", bWritten, len(msg.body))
				panic(err)
			}
		}
	}
}

func msgReceiver(ctx context.Context, childOut io.Reader, receiver MsgReceiver) {
	recvBuffer := make([]byte, 8)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			bRead, err := childOut.Read(recvBuffer)
			if err != nil {
				err = fmt.Errorf("%w: reading msg header", err)
				panic(err)
			}

			if bRead != 8 {
				err = fmt.Errorf("only %d/8 bytes read", bRead)
				panic(err)
			}

			msgId := MsgId(byteOrder.Uint32(recvBuffer[:4]))
			msgLength := byteOrder.Uint32(recvBuffer[4:8])

			body := make([]byte, msgLength)
			_, err = io.ReadAtLeast(childOut, body, int(msgLength))
			if err != nil {
				err = fmt.Errorf("%w: reading body of msg %d", err, msgId)
				panic(err)
			}

			err = receiver(msgId, body)
			if err != nil {
				err = fmt.Errorf("%w: handling msg %d", err, msgId)
				panic(err)
			}
		}
	}
}
