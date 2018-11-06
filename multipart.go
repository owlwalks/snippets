package snippets

import (
	"bytes"
	"io"
	"io/ioutil"
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

func receiveMultipart(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("upload")
	if err != nil {
		return
	}
	defer file.Close()

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		return
	}

	if err = tmp.Close(); err != nil {
		return
	}
}
