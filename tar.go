package snippets

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func createTar(name string, files []string) error {
	tarFile, err := os.Create(name)
	if err != nil {
		return err
	}

	tWriter := tar.NewWriter(tarFile)

	for _, fName := range files {
		content, err := ioutil.ReadFile(fName)
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: fName,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tWriter.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := tWriter.Write(content); err != nil {
			return err
		}
	}

	if err = tWriter.Close(); err != nil {
		return err
	}

	return tarFile.Close()
}

func extractTar(name string, extractTo string) error {
	tarFile, err := os.Open(name)
	if err != nil {
		return err
	}
	tReader := tar.NewReader(tarFile)

	if err = os.MkdirAll(extractTo, os.ModePerm); err != nil {
		return err
	}

	for {
		hdr, err := tReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		destpath := filepath.Join(extractTo, hdr.Name)
		if !strings.HasPrefix(destpath, filepath.Clean(extractTo)+string(os.PathSeparator)) {
			return errors.New("zip slip")
		}

		content, err := ioutil.ReadAll(tReader)
		if err != nil {
			return err
		}

		if hdr.FileInfo().IsDir() {
			if err = os.MkdirAll(destpath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(destpath), os.ModePerm); err != nil {
				return err
			}

			extractedF, err := os.Create(destpath)
			if err != nil {
				return err
			}

			if _, err = extractedF.Write(content); err != nil {
				return err
			}

			if err = extractedF.Close(); err != nil {
				return err
			}
		}
	}

	return tarFile.Close()
}
