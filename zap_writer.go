package logkit

import (
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type ConsoleWriter struct {
	out io.Writer
}

// NewConsoleWriter returns a console writer.
func NewConsoleWriter() io.Writer {
	return &ConsoleWriter{os.Stdout}
}

func (cw *ConsoleWriter) Write(p []byte) (n int, err error) {
	return cw.out.Write(p)
}

type FileWriter struct {
	out io.Writer
}

// NewFileWriter
func NewFileWriter(file string, maxSize, maxAge, maxBackups int, compress bool) io.Writer {
	out := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  true,
	}
	return &FileWriter{out}
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	return fw.out.Write(p)
}

type zapWriter struct {
	w io.Writer
}

// newZapWriter
func newZapWriter(w io.Writer) (zw *zapWriter) {
	return &zapWriter{w: w}
}

// Write implements io.Writer.
func (zw *zapWriter) Write(p []byte) (n int, err error) {
	return zw.w.Write(p)
}
