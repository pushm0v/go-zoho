package file

type FileUpload interface {
}

type fileUpload struct {
}

func NewFileUpload() FileUpload {
	return &fileUpload{}
}
