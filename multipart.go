package snippets

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func sendMultipart(url string, field string, name string) (*http.Response, error) {
	buf := new(bytes.Buffer)
	mWriter := multipart.NewWriter(buf)
	fWriter, err := mWriter.CreateFormFile(field, name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	io.Copy(fWriter, file)

	if err = mWriter.Close(); err != nil {
		return nil, err
	}

	return http.Post(url, mWriter.FormDataContentType(), buf)
}
