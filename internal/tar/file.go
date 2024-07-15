package tar

import (
	"archive/tar"
	"bytes"
	"github.com/radium-rtf/coderunner_lib/file"
	"io"
)

func NewTarArchive(files []file.File, uid int) (io.Reader, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 777,
			Size: int64(len(file.Content)),

			Uid: uid,
		}

		if err := tarWriter.WriteHeader(hdr); err != nil {
			return nil, err
		}

		if _, err := tarWriter.Write([]byte(file.Content)); err != nil {
			return nil, err
		}
	}

	return &buf, tarWriter.Close()
}
