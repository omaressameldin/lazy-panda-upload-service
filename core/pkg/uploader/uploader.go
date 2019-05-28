package uploader

type Uploader interface {
	UploadFile(pathToFile string) error
	DeleteFile(url string) error
}
