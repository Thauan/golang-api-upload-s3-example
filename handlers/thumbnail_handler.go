package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

func TakeFrame(path string, data []byte) (io.Reader, error) {
	reader := ReadFrameAsJpeg("./sample_data/in1.mp4", 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		panic(err)
	}
	err = imaging.Save(img, "./sample_data/out1.jpeg")
	if err != nil {
		panic(err)
	}

	return reader, nil
}
