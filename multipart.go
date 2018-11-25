package snippets

import (
	"bufio"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for {
		p, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if p.FormName() != "field_name" {
			continue
		}

		buf := bufio.NewReader(p)
		sniff, err := buf.Peek(512)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := http.DetectContentType(sniff)
		if contentType != "application/zip" {
			http.Error(w, "file type not allowed", http.StatusBadRequest)
			return
		}

		f, err := ioutil.TempFile("", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var maxSize int64 = 32 << 20
		lmt := io.LimitReader(p, maxSize+1)
		written, err := io.Copy(f, lmt)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if written > maxSize {
			os.Remove(f.Name())
			http.Error(w, "file size over limit", http.StatusBadRequest)
			return
		}
		// notice if f gets here then it is still in tmp, should move it to your persistent dir
	}
}
