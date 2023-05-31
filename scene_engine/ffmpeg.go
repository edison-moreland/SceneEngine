package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func ffmpegEncodeVideo(ctx context.Context, config messages.Config, framesFolder string) (*exec.Cmd, error) {
	inputImages := path.Join(framesFolder, "frame-%d.png")
	outputVideo := path.Clean(path.Join(framesFolder, "..", path.Base(framesFolder)+".mp4"))

	fmt.Printf("%s, %s \n", inputImages, outputVideo)

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-framerate", strconv.Itoa(int(config.FrameSpeed)),
		"-i", inputImages,
		"-pix_fmt", "yuv420p",
		"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
		"-vcodec", "h264",
		"-acodec", "aac",
		outputVideo,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd, cmd.Start()

	//	"/usr/bin/env",
	//	[
	//	"ffmpeg"; "-y"
	//	"-framerate"; Format.int[USize](config.frames_per_second)
	//	"-i"; input_files
	//	"-pix_fmt"; "yuv420p"
	//	"-vf"; "pad=ceil(iw/2)*2:ceil(ih/2)*2"
	//	"-vcodec"; "h264"
	//	"-acodec"; "aac"
	//	output_file
	//]
}
