package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

func TestDayToString(t *testing.T) {
	date := time.Date(2002, 1, 1, 23, 59, 59, 0, time.UTC)
	dateString := timeToString(date)
	assert.Equal(t, "01.01.2002", dateString)
}

func TestStringToDate(t *testing.T) {
	stringFromTimeValue, err := stringToTime("19.02.2022")
	assert.Nil(t, err)
	assert.Equal(t, 19, stringFromTimeValue.Day())
	assert.Equal(t, time.February, stringFromTimeValue.Month())
	assert.Equal(t, 2022, stringFromTimeValue.Year())

	stringFromTimeValue, err = stringToTime("index")
	assert.Equal(t, 1, stringFromTimeValue.Day())
	assert.Equal(t, time.January, stringFromTimeValue.Month())
	assert.Equal(t, 1, stringFromTimeValue.Year())
	assert.NotNil(t, err)
}

func TestTimeFromInit(t *testing.T) {
	date := time.Time{}
	assert.Equal(t, 1, date.Day())
	assert.Equal(t, time.January, date.Month())
	assert.Equal(t, 1, date.Year())
}

func TestSplitCorrectlyFile(t *testing.T) {
	filename := "19.02.2022.yaml"
	extension := filepath.Ext(filename)
	dateName := filename[0 : len(filename)-len(extension)]
	assert.Equal(t, "19.02.2022", dateName)
}

func TestLoadPreviousDay(t *testing.T) {
	dir, err := ioutil.TempDir("", "load_previous_day")
	assert.Nil(t, err)

	todayPath := path.Join(dir, fmt.Sprintf("%v.yaml", timeToString(time.Now())))
	_, err = os.Create(todayPath)
	assert.Nil(t, err)

	yesterdayPath := path.Join(dir, fmt.Sprintf("%v.yaml", timeToString(time.Now().Add(-24*time.Hour))))
	_, err = os.Create(yesterdayPath)
	assert.Nil(t, err)

	specificDayDate := time.Date(2002, time.August, 22, 2, 20, 20, 20, time.UTC)
	specificDayPath := path.Join(dir, fmt.Sprintf("%v.yaml", timeToString(specificDayDate)))
	_, err = os.Create(specificDayPath)
	assert.Nil(t, err)

	logService := NewLogService(dir)

	_, name, err := logService.getPreviousFileName()
	assert.Equal(t, timeToString(time.Now().Add(-24*time.Hour)), name)
	assert.Nil(t, err)

	err = os.RemoveAll(dir)
	assert.Nil(t, err)
}
