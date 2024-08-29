package openstack

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type ImageBufReader struct {
	io.Reader
	Name             string
	TotalSize        int
	readSize         int
	logProgressChunk int
	nextLogSize      int
}

func (buf *ImageBufReader) increment(n int) {
	buf.readSize += n
}

func (buf *ImageBufReader) Read(p []byte) (int, error) {
	n, err := buf.Reader.Read(p)
	defer buf.increment(n)
	if buf.readSize >= buf.nextLogSize {
		logging.Debug("read %s %d %%", buf.Name, buf.readSize*100/buf.TotalSize)
		buf.nextLogSize = min(buf.nextLogSize+buf.logProgressChunk, buf.TotalSize)
	}
	return n, err
}

func NewImageBufReader(name string, reader io.ReadCloser, size int) *ImageBufReader {
	return &ImageBufReader{
		Name:             name,
		TotalSize:        size,
		Reader:           bufio.NewReaderSize(reader, 1024*1024*8),
		logProgressChunk: size / 10,
	}
}

func NewImageBufReaderFromFile(imageFile string) (*ImageBufReader, error) {
	fileInfo, err := os.Stat(imageFile)
	if err != nil {
		return nil, fmt.Errorf("get image stat failed: %s", err)
	}
	reader, err := os.Open(imageFile)
	if err != nil {
		return nil, fmt.Errorf("open image faile: %s", err)
	}
	return NewImageBufReader(imageFile, reader, int(fileInfo.Size())), nil
}
