package core

import (
	"bufio"
	"os"
)

func GetAPIVersionFromFile(file string) (string, error) {
	var version string
	if _, err := os.Stat(file); err != nil {
		return version, err
	}

	f, err := os.Open(file)
	if err != nil {
		return version, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		version = scanner.Text()
	}

	err = scanner.Err()
	return version, err
}
