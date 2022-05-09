package ui

import (
	model2 "github.com/apoloa/bjournal/src/model"
	"testing"
)

func TestIncreaseIndex(t *testing.T) {
	list := NewList()
	log1 := model2.NewLog("1", model2.Irrelevant)
	list.AddItem(&log1, nil)

}
