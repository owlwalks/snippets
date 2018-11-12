package snippets

import (
	"bufio"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func sendMultipart(url string, field string, name string) (*http.Response, error) {
	r, w := io.Pipe()
	mWriter := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer mWriter.Close()

		fWriter, err := mWriter.CreateFormFile(field, name)
		if err != nil {
			return
		}

		file, err := os.Open(name)
		if err != nil {
			return
		}
		defer file.Close()

		if _, err = io.Copy(fWriter, file); err != nil {
			return
		}
	}()

	return http.Post(url, mWriter.FormDataContentType(), r)
}

func receiveMultipart(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("upload")
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	b, _ := reader.Peek(512)
	contentType := http.DetectContentType(b)
	if !strings.HasPrefix(contentType, "image") {
		return
	}

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
