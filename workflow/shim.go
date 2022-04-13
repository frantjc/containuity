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
	shim     []byte
	//go:embed sqncshim-uses
	shimUses []byte
	shimName = "sqncshim"
)

func init() {
	go func () {
		if _, err := shimTarArchive(); err != nil {
			panic("github.com/frantjc/sequence/workflow.shim is not able to be tarballed")
		}
	
		if _, err := shimUsesTarArchive(); err != nil {
			panic("github.com/frantjc/sequence/workflow.shimUses is not able to be tarballed")
		}
	}()
}

func shimTarArchive() (io.Reader, error) {
	tarball := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(tarball)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    shimName,
		Size:    int64(len(shim)),
		Mode:    0777,
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	}

	if _, err := io.Copy(tarWriter, bytes.NewBuffer(shim)); err != nil {
		return nil, err
	}

	return tarball, nil
}

func shimUsesTarArchive() (io.Reader, error) {
	tarball := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(tarball)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    shimName,
		Size:    int64(len(shim)),
		Mode:    0777,
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	}

	if _, err := io.Copy(tarWriter, bytes.NewBuffer(shimUses)); err != nil {
		return nil, err
	}

	return tarball, nil
}
