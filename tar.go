package snippets

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func createTar(name string, files []string) error {
	fWriter, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fWriter.Close()

	writer := tar.NewWriter(fWriter)

	for _, fName := range files {
		reader, err := os.Open(fName)
		if err != nil {
			return err
		}

		stat, err := reader.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: fName,
			Mode: 0600,
			Size: stat.Size(),
		}
		if err := writer.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(writer, reader); err != nil {
			return err
		}

		if err = reader.Close(); err != nil {
			return err
		}
	}

	if err = writer.Close(); err != nil {
		return err
	}

	return nil
}

func extractTar(name string, extractTo string) error {
	fReader, err := os.Open(name)
	if err != nil {
		return err
	}
	defer fReader.Close()

	reader := tar.NewReader(fReader)

	if err = os.MkdirAll(extractTo, os.ModePerm); err != nil {
		return err
	}

	for {
		hdr, err := reader.Next()
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

		if hdr.FileInfo().IsDir() {
			if err = os.MkdirAll(destpath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(destpath), os.ModePerm); err != nil {
				return err
			}

			writer, err := os.Create(destpath)
			if err != nil {
				return err
			}

			if _, err = io.Copy(writer, reader); err != nil {
				return err
			}

			if err = writer.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
