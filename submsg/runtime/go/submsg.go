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
type MsgSender func(id MsgId, length uint32, body io.Reader)
type MsgReceiver func(id MsgId, body io.Reader) error

type msg struct {
	id     MsgId
	length uint32
	body   io.Reader
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

	return func(id MsgId, length uint32, body io.Reader) {
		sendQueue <- msg{
			id:     id,
			length: length,
			body:   body,
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

			err = binary.Write(childIn, byteOrder, msg.length)
			if err != nil {
				err = fmt.Errorf("%w: writing msg length to child", err)
				panic(err)
			}

			if msg.length <= 0 {
				break
			}

			bWritten, err := io.CopyN(childIn, msg.body, int64(msg.length))
			if err != nil {
				err = fmt.Errorf("%w: writing msg body to child", err)
				panic(err)
			}

			if uint32(bWritten) != msg.length {
				err = fmt.Errorf("only %d/%d bytes written", bWritten, msg.length)
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

			limitedReader := io.LimitReader(childOut, int64(msgLength))

			err = receiver(msgId, limitedReader)
			if err != nil {
				err = fmt.Errorf("%w: handling msg %d", err, msgId)
				panic(err)
			}

			// Throw away anything not read by the handler
			_, _ = io.ReadAll(limitedReader)
		}
	}
}
