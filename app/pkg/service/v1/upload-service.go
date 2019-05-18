package v1

import (
	"fmt"
	"io"
	"os"
	"time"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/pkg/api/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

// UploadServiceServer is implementation of v1.userServiceServer proto interface
type UploadServiceServer struct{}

// NewUploadServiceServer creates Upload service
func NewUploadServiceServer() *UploadServiceServer {
	return &UploadServiceServer{}
}

func uploadFailed(stream v1.UploadService_UploadServer, uploadError error) error {
	err := stream.SendAndClose(&v1.UploadStatusResponse{
		Api:     apiVersion,
		Message: "Upload failed",
		Code:    v1.UploadStatusCode_Failed,
	})
	if err != nil {
		return err
	}

	return status.Error(codes.Unknown, "failed to upload file-> "+uploadError.Error())
}

// Upload new file
func (s *UploadServiceServer) Upload(stream v1.UploadService_UploadServer) error {
	filedata, err := stream.Recv()
	if err != nil {
		return uploadFailed(stream, err)
	}

	f, err := os.Create(fmt.Sprintf(
		"./app/%s_%v.%s",
		filedata.GetMeta().GetFileName(),
		time.Now(),
		filedata.GetMeta().GetFileType(),
	))
	if err != nil {
		return uploadFailed(stream, err)
	}

	defer f.Close()

	for {
		filedata, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return uploadFailed(stream, err)
		}

		chunks := filedata.GetContent()
		_, err = f.Write(chunks)
		if err != nil {
			return uploadFailed(stream, err)
		}
	}

	err = stream.SendAndClose(&v1.UploadStatusResponse{
		Api:     apiVersion,
		Message: "Upload succeeded",
		Code:    v1.UploadStatusCode_Ok,
	})
	return err
}
