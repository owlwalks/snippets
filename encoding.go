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
