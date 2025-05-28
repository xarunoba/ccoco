package ccoco

import "errors"

type (
	File struct {
		Name    string
		Content *FileContent
	}
	FileContent struct {
		Files []string `json:"files"`
	}
)

func (f *File) CheckState() error {
	if f.Name == "" {
		return errors.New("file name is empty")
	}
	if f.Content == nil {
		return errors.New("file content is nil")
	}
	if err := f.Content.CheckState(); err != nil {
		return err
	}
	return nil
}

func (fc *FileContent) CheckState() error {
	if len(fc.Files) == 0 {
		return errors.New("file content is empty")
	}
	return nil
}
