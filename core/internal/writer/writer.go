package writer

import (
	"fmt"
	"io"
	"os"
	"time"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/api/v1"
)

const (
	tmpPath = "."
)

func CreateTmpFile(stream v1.UploadService_UploadServer) (*os.File, error) {
	filedata, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	f, err := os.Create(fmt.Sprintf(
		"%s/%s_%v.%s",
		tmpPath,
		filedata.GetMeta().GetFileName(),
		time.Now(),
		filedata.GetMeta().GetFileType(),
	))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// writetoFile writes stream content to already opened file
func WriteToFile(stream v1.UploadService_UploadServer, f *os.File) error {
	defer f.Close()

	for {
		filedata, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		chunks := filedata.GetContent()
		if _, err = f.Write(chunks); err != nil {
			return err
		}
	}

	return nil
}

func DeleteTmp(f *os.File) {
	if f == nil {
		return
	}
	os.Remove(f.Name())
}
