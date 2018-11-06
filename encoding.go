package snippets

import (
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
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

func readXML(name string) (interface{}, error) {
	// <?xml version="1.0" encoding="UTF-8"?>
	// <Person>
	// 	<FullName>Grace R. Emlin</FullName>
	// 	<Email where="home">
	// 		<Addr>gre@example.com</Addr>
	// 	</Email>
	// 	<Email where='work'>
	// 		<Addr>gre@work.com</Addr>
	// 	</Email>
	// 	<Group>
	// 		<Value>Friends</Value>
	// 		<Value>Squash</Value>
	// 	</Group>
	// </Person>
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	data := struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"FullName"`
		Email   []struct {
			Where string `xml:"where,attr"`
			Addr  string
		} `xml:"Email"`
		Groups []string `xml:"Group>Value"`
	}{}
	err = xml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeXML(name string, data interface{}) error {
	fWriter, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fWriter.Close()

	enc := xml.NewEncoder(fWriter)
	enc.Indent("  ", "    ")
	// type Person struct {
	// 	XMLName   xml.Name `xml:"Person"`
	// 	ID        int      `xml:"id,attr"`
	// 	FirstName string   `xml:"Name>First"`
	// 	LastName  string   `xml:"Name>Last"`
	// 	Age       int      `xml:"Age"`
	// 	Height    float32  `xml:"Height,omitempty"`
	// 	Comment   string   `xml:",comment"`
	// }
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}
