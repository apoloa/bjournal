package model

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)

type DailyLog struct {
	key      string `yaml:"-"`
	basePath string `yaml:"-"`
	Logs     []Log  `yaml:"items"`
}

func NewDailyLog(date, basePath string) DailyLog {
	return DailyLog{
		key:      date,
		basePath: basePath,
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

func DailyFrom(from []byte, date string, dir string) (DailyLog, error) {
	dailyLog := DailyLog{}
	err := yaml.Unmarshal(from, &dailyLog)
	if err != nil {
		return dailyLog, err
	}
	dailyLog.key = date
	dailyLog.basePath = dir
	dailyLog.fullRead()
	return dailyLog, nil
}

func (d DailyLog) ToBytes() ([]byte, error) {
	return yaml.Marshal(d)
}
