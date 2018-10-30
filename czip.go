package snippets

import (
	"archive/zip"
	"io/ioutil"
	"os"
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
