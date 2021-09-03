package xml

import (
	"encoding/xml"
	"io"
	"strings"
)

func GetStartTag(bytes []byte) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(bytes)))

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return "", nil
			}
			return "", err
		}
		if se, ok := t.(xml.StartElement); ok {
			return se.Name.Local, nil
		}
	}
}
