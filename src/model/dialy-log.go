package model

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type DailyLog struct {
	key      string    `yaml:"-"`
	basePath string    `yaml:"-"`
	Date     time.Time `json:"-" yaml:"-"`
	Logs     []Log     `json:"logs" yaml:"items"`
}

func NewDailyLog(date, basePath string) DailyLog {
	return DailyLog{
		key:      date,
		basePath: basePath,
		Date:     time.Now(),
		Logs:     []Log{},
	}
}

func (d *DailyLog) fullRead() {
	for index, item := range d.Logs {
		if item.Url != nil {
			filePath := path.Join(d.basePath, *item.Url)
			file, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}
			stringFile := string(file)
			d.Logs[index].Text = &stringFile
		}
	}
}

func (d *DailyLog) setParent() {
	for id, parent := range d.Logs {
		if parent.SubLogs == nil {
			continue
		}
		if len(*parent.SubLogs) == 0 {
			continue
		}
		for index, _ := range *parent.SubLogs {
			(*parent.SubLogs)[index].Parent = &d.Logs[id]
		}
	}
}

func (d *DailyLog) addUUID() {
	for index, _ := range d.Logs {
		d.Logs[index].Id = uuid.NewString()
	}
}

func DailyFrom(from []byte, dateTime time.Time, date string, dir string) (DailyLog, error) {
	dailyLog := DailyLog{}
	err := yaml.Unmarshal(from, &dailyLog)
	if err != nil {
		return dailyLog, err
	}
	dailyLog.key = date
	dailyLog.Date = dateTime
	dailyLog.basePath = dir
	dailyLog.fullRead()
	dailyLog.setParent()
	dailyLog.addUUID()
	return dailyLog, nil
}

func (d DailyLog) ToBytes() ([]byte, error) {
	return yaml.Marshal(d)
}
