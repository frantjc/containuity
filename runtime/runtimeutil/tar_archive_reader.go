package runtimeutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"time"
)

func NewSingleFileTarArchiveReader(name string, b []byte) (io.Reader, error) {
	tarArchive := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(tarArchive)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    name,
		Size:    int64(len(b)),
		Mode:    0777,
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	}

	if _, err := io.Copy(tarWriter, bytes.NewReader(b)); err != nil {
		return nil, err
	}

	return tarArchive, nil
}
