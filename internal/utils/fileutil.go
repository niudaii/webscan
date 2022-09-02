package utils

import (
	"bufio"
	"io/ioutil"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil || info == nil {
		return false
	}
	return !info.IsDir()
}

func ReadLines(filename string) (lines []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return
}

func ReadFile(filename string) (bytes []byte, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return data, nil
}
