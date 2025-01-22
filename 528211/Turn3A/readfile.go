package Turn3A

import (
	"errors"
	"os"
)

// ReadFile reads a file and returns its content as a string, or an error if unsuccessful.
func ReadFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "", err
	}

	if stat.Size() == 0 {
		return "", errors.New("file is empty")
	}

	content := make([]byte, stat.Size())
	_, err = f.Read(content)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
