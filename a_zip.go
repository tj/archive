package archive

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

type zipArchive struct {
	writer *zip.Writer
}

func newZipArchive(dst io.Writer) *zipArchive {
	return &zipArchive{
		writer: zip.NewWriter(dst),
	}
}

func (a *zipArchive) addBytes(path string, contents []byte) error {
	h := &zip.FileHeader{
		Name: path,
	}
	f, err := a.writer.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = f.Write(contents)
	return err
}

func (a *zipArchive) addFile(path string, info os.FileInfo, f *os.File) error {
	if !info.Mode().IsRegular() {
		return errors.New("Only regular files supported: " + path)
	}
	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	h.Name = path
	zf, err := a.writer.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = io.Copy(zf, f)
	return err
}

func (a *zipArchive) close() error {
	return a.writer.Close()
}
