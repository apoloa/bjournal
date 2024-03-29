package model

import (
	"path"

	"gopkg.in/yaml.v3"
)

type Index struct {
	Items []IndexItem `json:"items" yaml:"items"`
}

func IndexFromFile(basePath string, from []byte) (Index, error) {
	index := Index{}
	err := yaml.Unmarshal(from, &index)
	if err != nil {
		return index, err
	}
	for i, item := range index.Items {
		index.Items[i].FullUrl = path.Join(basePath, item.Url)
	}
	return index, nil
}

func (i Index) ToBytes() ([]byte, error) {
	return yaml.Marshal(i)
}
