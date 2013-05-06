package clog

import (
	"log"
	"os"
	"sync"
)

// A File.
// Implements io.Writer.
type File struct {
	currentFile *os.File
	filename    string
	mtx         sync.Mutex
}

// Opens or creates a new File using the specified file name.
//
// 		var Log *clog.Clog = clog.NewClog()
// 		file := clog.NewFile("some_file.log")
//		Log.AddOutput(file, clog.LevelWarning)
//		defer file.Close()
func NewFile(filename string) *File {
	return &File{nil, filename, sync.Mutex{}}
}

func (this *File) Write(p []byte) (n int, err error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.currentFile == nil {
		if err := this.createFile(); err != nil {
			return 0, err
		}
	}
	return this.currentFile.Write(p)
}

func (this *File) Close() {
	if this.currentFile != nil {
		this.currentFile.Close()
	}
}

func (this *File) createFile() error {
	if this.currentFile != nil {
		this.currentFile.Close()
	}
	var err error
	this.currentFile, err = os.OpenFile(
		this.filename,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666)
	if err != nil {
		log.Printf("Clog: Unable to open %v for writing\n", this.filename)
		return err
	}
	return nil
}
