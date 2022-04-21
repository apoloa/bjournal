package ui

import (
	"github.com/apoloa/bjournal/model"
	"testing"
)

func TestIncreaseIndex(t *testing.T) {
	list := NewList()
	log1 := model.NewLog("1", model.Irrelevant)
	list.AddItem(&log1, nil)

}
