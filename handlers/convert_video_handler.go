package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func VideoToMP4(path string, data []byte) (*bufio.Reader, error) {
	cmd := exec.Command(path,
		"-i", // read from stdin
		"-codec",
		"copy",
		"-", // read from stdout
	)
	// stdin read
	cmd.Stdin = bytes.NewBuffer(data)

	var buffer bytes.Buffer
	// stdout set buffer
	cmd.Stdout = &buffer

	if cmd.Run() != nil {
		return nil, fmt.Errorf("Cannot run or found ffmpeg file by path: %s", path)
	}

	// reading of stdout
	reader := bufio.NewReader(&buffer)

	return reader, nil
}
