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

func TakeFrame(currentFile string, outputFile string) (io.Reader, error) {
	// currentDir := utils.RootDir()
	// outputTempPath := fmt.Sprintf("%s/%s", currentDir+"/"+GetEnvWithKey("OUTPUT_TEMPORARY_PATH"), outputFile)
	// inputTempPath := fmt.Sprintf("%s/%s", currentDir, currentFile)

	reader := ReadFrameAsJpeg("./sample_data/1.mp4", 5)

	img, err := imaging.Decode(reader)

	if err != nil {
		print("Imaging Open error")
	}

	// imageOutput := fmt.Sprintf("%s", outputTempPath)

	err2 := imaging.Save(img, "./sample_data/out2.jpeg")

	if err2 != nil {
		print("Imaging Save error")
	}

	return reader, nil
}

// func ReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, error) {
// 	buf := bytes.NewBuffer(nil)

// 	err := ffmpeg.Input(inFileName).
// 		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
// 		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
// 		WithOutput(buf, os.Stdout).
// 		Run()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return buf, nil
// }

// func TakeFrame(filePath string, frame int) (io.Reader, error) {
// 	reader, err := ReadFrameAsJpeg(filePath, frame)

// 	if err != nil {
// 		fmt.Printf("could not take a frame: %v", err)
// 	}

// 	return reader, nil
// }
