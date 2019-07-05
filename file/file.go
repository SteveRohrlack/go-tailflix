package file

import (
	"errors"
	"io"
	"os"
	"strings"
	"tailflix/math"
	"time"
)

/*
Open opens a file at the given path
it returns an error if opening is not possible or if the given path points to a directory
*/
func Open(filePath string) (file *os.File, err error) {

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, errors.New("tailflix cannot handle directories, please specify a regular file")
	}

	file, err = os.Open(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

/*
Preview returns the last "numberOfLines" lines of the given file content
it may return less because the line length is currently approximated with 80 bytes
*/
func Preview(file *os.File, numberOfLines int) string {

	const lineSeparator = "\n"

	lineLength := int64(80 * numberOfLines)

	readBuffer := make([]byte, lineLength)

	file.Seek(-lineLength, 2)
	file.Read(readBuffer)

	lines := strings.Split(string(readBuffer), lineSeparator)
	lineCount := len(lines) - 1

	previewStartIndex := lineCount - math.Min(lineCount, numberOfLines)
	previewLines := lines[previewStartIndex:lineCount]

	return strings.Join(previewLines, lineSeparator)
}

/*
Watch watches the given file and sends newly added bytes (to the end of the file)
to the data channel (one byte per send)
*/
func Watch(file *os.File, data chan<- byte, errors chan<- error) {

	const sendDataDelay = 600 * time.Microsecond
	const idleDelay = 200 * time.Millisecond

	readBuffer := make([]byte, 1)
	file.Seek(-1, 2)

	var err error

	for {
		_, err = file.Read(readBuffer)

		if err == io.EOF {
			// idle
			time.Sleep(idleDelay)
			continue
		}

		if err != nil {
			errors <- err
			break
		}

		data <- readBuffer[0]
		time.Sleep(sendDataDelay)
	}
}
