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
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// demonstration purpose - we have 2 fields named: text_field and file_file
	// they will be parsed in order
	var _ = []string{"text_field", "file_field"}

	// parse text field
	var text = make([]byte, 512)
	p, err := reader.NextPart()
	// one more field to parse, EOF is considered as failure here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "text_field" {
		http.Error(w, "text_field is expected", http.StatusBadRequest)
	}

	_, err = p.Read(text)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse file field
	p, err = reader.NextPart()
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "file_field" {
		http.Error(w, "file_field is expected", http.StatusBadRequest)
	}

	buf := bufio.NewReader(p)
	sniff, _ := buf.Peek(512)
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
	// here we should have text and f to be processed accordingly
}
