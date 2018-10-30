package snippets

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func createZip(name string, files []string) error {
	zipFile, err := os.Create(name)
	if err != nil {
		return err
	}

	zWriter := zip.NewWriter(zipFile)

	for _, fName := range files {
		content, err := ioutil.ReadFile(fName)
		if err != nil {
			return err
		}
		fWriter, err := zWriter.Create(fName)
		if err != nil {
			return err
		}
		_, err = fWriter.Write(content)
		if err != nil {
			return err
		}
	}

	err = zWriter.Close()
	if err != nil {
		return err
	}

	return zipFile.Close()
}

func extractZip(name string, extractTo string) error {
	zReader, err := zip.OpenReader(name)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(extractTo, os.ModePerm); err != nil {
		return err
	}

	for _, f := range zReader.File {
		destpath := filepath.Join(extractTo, f.Name)

		if !strings.HasPrefix(destpath, filepath.Clean(extractTo)+string(os.PathSeparator)) {
			return errors.New("zip slip")
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() {
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

	return zReader.Close()
}
