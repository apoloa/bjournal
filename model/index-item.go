package model

import "path"

type IndexItem struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	FullUrl string `yaml:"-"`
}

func NewIndexItem(name, url, baseUrl string) IndexItem {
	return IndexItem{
		Name:    name,
		Url:     url,
		FullUrl: path.Join(baseUrl, url),
	}
}
