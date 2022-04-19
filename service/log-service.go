package service

import (
	"bjournal/model"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path"
	"time"
)

type LogService struct {
	baseDir string
	cache   map[string]model.DailyLog
	index   model.Index
}

func NewLogService(baseDir string) *LogService {
	return &LogService{
		baseDir: baseDir,
		cache:   make(map[string]model.DailyLog),
	}
}

func parseDay(date time.Time) string {
	return fmt.Sprintf("%02d.%02d.%v", date.Day(), int(date.Month()), date.Year())
}

func (m *LogService) ReadDay(date time.Time) (model.DailyLog, error) {
	dateString := parseDay(date)
	log.Print(dateString)
	if val, ok := m.cache[dateString]; ok {
		return val, nil
	} else {
		dailyLog, err := m.ReadDailyLog(dateString)
		if err != nil {
			return model.DailyLog{}, err
		}
		m.cache[dateString] = dailyLog
		return dailyLog, nil
	}
}

func (m *LogService) ReadDailyLog(date string) (model.DailyLog, error) {
	filePath := path.Join(m.baseDir, fmt.Sprintf("%v.yaml", date))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Print("Error reading the file")
		log.Print(err.Error())
		return model.NewDailyLog(date, m.baseDir), nil
	}
	return model.DailyFrom(file, date, m.baseDir)
}

func (m *LogService) AddNewLog(date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := parseDay(date)
	dailyLog, _ := m.cache[dateString]
	dailyLog.Logs = append(dailyLog.Logs, model.NewLog(name, category))
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) AppendNewLog(uuid string, date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := parseDay(date)
	dailyLog, _ := m.cache[dateString]
	for index, appendLog := range dailyLog.Logs {
		if appendLog.Id == uuid {
			dailyLog.Logs[index].AppendNewSubLog(name, category)
		}
	}
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) GetPreviousDate() (model.DailyLog, error) {

}

func (m *LogService) SaveLog(date time.Time) (model.DailyLog, error) {
	dateString := parseDay(date)
	dailyLog, _ := m.cache[dateString]
	filePath := path.Join(m.baseDir, fmt.Sprintf("%v.yaml", dateString))
	bytes, err := dailyLog.ToBytes()
	if err != nil {
		return dailyLog, err
	}
	err = ioutil.WriteFile(filePath, bytes, 0666)
	if err != nil {
		return dailyLog, err
	}
	return dailyLog, nil
}
