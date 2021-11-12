package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

func ReadOneLine(p string) (string, error) {
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	idx := bytes.IndexByte(buf, '\n')
	if idx == -1 {
		return "", fmt.Errorf("file does not contain a complete line: %s", p)
	}
	if idx != len(buf)-1 {
		return "", errors.New("trailing junk after first line")
	}
	return string(buf[:len(buf)-1]), nil
}
