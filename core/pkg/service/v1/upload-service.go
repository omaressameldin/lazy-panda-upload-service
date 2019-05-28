package v1

import (
	"context"

	writer "github.com/omaressameldin/lazy-panda-upload-service/core/internal/writer"
	v1 "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/api/v1"
	uploader "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/uploader"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

// UploadServiceServer is implementation of v1.userServiceServer proto interface
type UploadServiceServer struct {
	fileUploader uploader.Uploader
}

// NewUploadServiceServer creates Upload service
func NewUploadServiceServer(fileUploader uploader.Uploader) *UploadServiceServer {
	return &UploadServiceServer{fileUploader}
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
	f, err := writer.CreateTmpFile(stream)
	if err != nil {
		return uploadFailed(stream, err)
	}
	defer writer.DeleteTmp(f)

	if err = writer.WriteToFile(stream, f); err != nil {
		return uploadFailed(stream, err)
	}

	if err = s.fileUploader.UploadFile(f.Name()); err != nil {
		return err
	}

	return stream.SendAndClose(&v1.UploadStatusResponse{
		Api:     apiVersion,
		Message: "Upload succeeded",
		Code:    v1.UploadStatusCode_Ok,
	})
}

func (s *UploadServiceServer) Delete(
	ctx context.Context,
	req *v1.DeleteRequest,
) (*v1.DeleteResponse, error) {
	if err := s.fileUploader.DeleteFile(req.Url); err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete File-> "+err.Error())
	}

	return &v1.DeleteResponse{
		Api: apiVersion,
	}, nil
}
