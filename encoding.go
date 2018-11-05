package snippets

import (
	"encoding/csv"
	"os"
)

func parseCSV(name string) (<-chan []string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	out := make(chan []string)
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				file.Close()
				close(out)
				break
			}
			out <- record
		}
	}()

	return out, nil
}

func writeCsv(name string, rows <-chan []string) error {
	fWriter, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fWriter.Close()

	writer := csv.NewWriter(fWriter)
	for row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
		writer.Flush()
	}

	return nil
}
