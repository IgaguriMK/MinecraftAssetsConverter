package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	bytes, err := ioutil.ReadFile("./indexes/legacy.json")
	if err != nil {
		panic(err)
	}

	var indexes Indexes
	err = json.Unmarshal(bytes, &indexes)
	if err != nil {
		panic(err)
	}

	for filename, index := range indexes.Objects {
		CopyFile(filename, index)
	}
}

func CopyFile(filename string, index FileIndex) {
	dst := createVirtual(filename)
	defer dst.Close()

	srcPath := objPath(index.Hash)
	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	fmt.Println(filename, "<-", srcPath)

	var w io.Writer = dst

	if HasExtensions(filename, ".mcmeta", ".txt") {
		w = NewCrLfWriter(w)
	}

	io.Copy(w, src)
}

func objPath(hash string) string {
	md := hash[0:2]
	return "objects/" + md + "/" + hash
}

func createVirtual(path string) *os.File {
	path = "virtual/legacy/" + path

	ps := strings.Split(path, "/")
	dirs := ps[0 : len(ps)-1]
	dirPath := strings.Join(dirs, "/")

	err := os.MkdirAll(dirPath, os.ModeDir)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return file
}

type Indexes struct {
	Objects map[string]FileIndex `json:"objects"`
}

type FileIndex struct {
	Hash string `json:"hash"`
	Size int    `json:"size"`
}

func HasExtensions(name string, extensions ...string) bool {
	for _, ex := range extensions {
		if strings.HasSuffix(name, ex) {
			return true
		}
	}

	return false
}

type CrLfWriter struct {
	dst io.Writer
}

func NewCrLfWriter(dst io.Writer) *CrLfWriter {
	return &CrLfWriter{dst}
}

func (cr *CrLfWriter) Write(p []byte) (n int, err error) {
	tmpBs := make([]byte, 0, len(p))

	for _, b := range p {
		if b == byte('\n') {
			tmpBs = append(tmpBs, byte('\r'), byte('\n'))
		} else {
			tmpBs = append(tmpBs, b)
		}
	}

	n, err = cr.dst.Write(tmpBs)
	return
}
