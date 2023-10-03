package ui

import (
	"testing"

	model2 "github.com/apoloa/bjournal/src/model"
)

func TestIncreaseIndex(t *testing.T) {
	list := NewList()
	log1 := model2.NewLog("1", model2.Irrelevant)
	list.AddItem(&log1, nil)

}
