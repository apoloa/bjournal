package ui

import (
	"bjournal/model"
	"testing"
)

func TestIncreaseIndex(t *testing.T) {
	list := NewList()
	log1 := model.NewLog("1", model.Irrelevant)
	list.AddItem(&log1, nil)

}
