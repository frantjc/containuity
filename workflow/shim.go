package workflow

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"time"
)

var (
	//go:embed sqncshim
	Shim []byte
)

func init() {
	if _, err := shimTarArchive(); err != nil {
		panic("github.com/frantjc/sequence/workflow.Shim is not able to be tarballed")
	}
}

func shimTarArchive() (io.Reader, error) {
	tarball := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(tarball)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    "sqncshim",
		Size:    int64(len(Shim)),
		Mode:    0777,
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	}

	if _, err := io.Copy(tarWriter, bytes.NewBuffer(Shim)); err != nil {
		return nil, err
	}

	return tarball, nil
}
