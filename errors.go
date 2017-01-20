package main

import (
	"errors"
	"io/ioutil"
	"log"
)

// fileErrReader reads successive files until it reads without error, then stops reading.
type fileErrReader struct {
	file []byte
	err  error
}

func newFileErrReader() fileErrReader { return fileErrReader{err: errors.New("")} }

func (fr *fileErrReader) read(filename string) {
	if fr.err == nil {
		return
	}
	fr.file, fr.err = ioutil.ReadFile(filename)
}

func (fr fileErrReader) String() string { return string(fr.file) }

// Log the error and return the value.
func logErr(s string, err error) string {
	if err != nil {
		log.Println(err)
	}
	return s
}
