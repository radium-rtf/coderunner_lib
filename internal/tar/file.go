package tar

import (
	"archive/tar"
	"bytes"
	"github.com/radium-rtf/coderunner_lib/file"
	"io"
)

const (
	modeALL = 777
)

func NewTarArchive(files []file.File, uid int) (io.Reader, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: modeALL,
			Size: file.Content.Size(),

			Uid: uid,
		}

		if err := tarWriter.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if _, err := tarWriter.Write(file.Content.GetBytes()); err != nil {
			return nil, err
		}
	}

	return &buf, tarWriter.Close()
}
