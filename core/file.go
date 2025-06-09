package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

type (
	File struct{}
)

func (file *File) ExtractData(path string) (map[string]any, error) {
	nameFile := filepath.Base(path)
	ext := filepath.Ext(nameFile)

	pathAbsolute, errAbs := filepath.Abs(path)

	if errAbs != nil {
		return nil, errAbs
	}
	var object map[string]any
	output, err := os.ReadFile(pathAbsolute)

	if err != nil {
		return nil, err
	}

	switch ext {
	case ".json":
		errParser := json.Unmarshal(output, &object)
		if errParser != nil {
			return nil, errParser
		}

	case ".yaml":
		errParser := yaml.Unmarshal(output, &object)
		if errParser != nil {
			return nil, errParser
		}
	}

	fmt.Println(object)
	return object, nil

}
