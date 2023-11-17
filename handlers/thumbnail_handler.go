package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)

	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf, nil
}

func TakeFrame(filePath string, frame int) (io.Reader, error) {
	reader, err := ReadFrameAsJpeg(filePath, frame)

	if err != nil {
		fmt.Printf("could not take a frame: %v", err)
	}

	return reader, nil
}
