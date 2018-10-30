package snippets

import (
	"archive/tar"
	"io/ioutil"
	"os"
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
