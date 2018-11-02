package snippets

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func createZip(name string, files []string) error {
	zipFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zWriter := zip.NewWriter(zipFile)

	for _, fName := range files {
		reader, err := os.Open(fName)
		if err != nil {
			return err
		}

		writer, err := zWriter.Create(fName)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}

		if err = reader.Close(); err != nil {
			return err
		}
	}

	err = zWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func extractZip(name string, extractTo string) error {
	zReader, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zReader.Close()

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

			if _, err = io.Copy(extractedF, rc); err != nil {
				return err
			}

			if err = extractedF.Close(); err != nil {
				return err
			}
		}

		if err = rc.Close(); err != nil {
			return err
		}
	}

	return nil
}
