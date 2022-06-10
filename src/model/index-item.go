package model

import "path"

type IndexItem struct {
	Name    string `json:"name" yaml:"name"`
	Url     string `json:"url" yaml:"url"`
	FullUrl string `json:"-" yaml:"-"`
}

func NewIndexItem(name, url, baseUrl string) IndexItem {
	return IndexItem{
		Name:    name,
		Url:     url,
		FullUrl: path.Join(baseUrl, url),
	}
}
