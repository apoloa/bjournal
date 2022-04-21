package service

import (
	"fmt"
	"github.com/apoloa/bjournal/model"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"
)

const layout = "02.01.2006"

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

func timeToString(date time.Time) string {
	return fmt.Sprintf("%02d.%02d.%v", date.Day(), int(date.Month()), date.Year())
}

func stringToTime(date string) (time.Time, error) {
	return time.Parse(layout, date)
}

func (m *LogService) ReadDay(date time.Time) (model.DailyLog, error) {
	dateString := timeToString(date)
	if val, ok := m.cache[dateString]; ok {
		return val, nil
	} else {
		dailyLog, err := m.ReadDailyLog(date, dateString)
		if err != nil {
			return model.DailyLog{}, err
		}
		m.cache[dateString] = dailyLog
		return dailyLog, nil
	}
}

func (m *LogService) ReadDailyLog(dateTime time.Time, date string) (model.DailyLog, error) {
	filePath := path.Join(m.baseDir, fmt.Sprintf("%v.yaml", date))
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Print("Error reading the file")
		log.Print(err.Error())
		return model.NewDailyLog(date, m.baseDir), nil
	}
	return model.DailyFrom(file, dateTime, date, m.baseDir)
}

func (m *LogService) AddNewLog(date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := timeToString(date)
	dailyLog, _ := m.cache[dateString]
	dailyLog.Logs = append(dailyLog.Logs, model.NewLog(name, category))
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) AppendNewLog(uuid string, date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := timeToString(date)
	dailyLog, _ := m.cache[dateString]
	for index, appendLog := range dailyLog.Logs {
		if appendLog.Id == uuid {
			dailyLog.Logs[index].AppendNewSubLog(name, category)
		}
	}
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) getPreviousFileName() (time.Time, string, error) {
	files, err := ioutil.ReadDir(m.baseDir)
	if err != nil {
		log.Print(err)
		return time.Now(), "", err
	}
	startTime := time.Time{}
	previousFileName := ""
	actualDate := timeToString(time.Now())
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			extension := filepath.Ext(filename)
			dateName := filename[0 : len(filename)-len(extension)]
			if dateName == actualDate {
				continue
			}
			toTime, err := stringToTime(dateName)
			if err != nil {
				continue
			}
			if startTime.Before(toTime) {
				startTime = toTime
				previousFileName = dateName
			}
		}
	}
	return startTime, previousFileName, nil
}

func (m *LogService) GetPreviousDate() (model.DailyLog, error) {
	date, dateString, err := m.getPreviousFileName()
	if err != nil {
		return model.DailyLog{}, err
	}
	if val, ok := m.cache[dateString]; ok {
		return val, nil
	} else {
		dailyLog, err := m.ReadDailyLog(date, dateString)
		if err != nil {
			return model.DailyLog{}, err
		}
		m.cache[dateString] = dailyLog
		return dailyLog, nil
	}
}

func (m *LogService) SaveLog(date time.Time) (model.DailyLog, error) {
	dateString := timeToString(date)
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
