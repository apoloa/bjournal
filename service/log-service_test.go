package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDayParser(t *testing.T) {
	date := time.Date(2002, 1, 1, 23, 59, 59, 0, time.UTC)
	dateString := parseDay(date)
	assert.Equal(t, "01.01.2002", dateString)
}
