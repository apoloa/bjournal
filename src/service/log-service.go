package service

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/apoloa/bjournal/src/model"
	"github.com/apoloa/bjournal/src/utils"
	"github.com/apoloa/bjournal/src/utils/timeconv"
	zerolog "github.com/rs/zerolog/log"
)

const indexFile = "index.yaml"

type LogService struct {
	baseDir string
	cache   map[string]model.DailyLog
	Index   model.Index
}

func NewLogService(baseDir string) *LogService {
	return &LogService{
		baseDir: baseDir,
		cache:   make(map[string]model.DailyLog),
		Index:   readIndex(baseDir),
	}
}

func readIndex(baseDir string) model.Index {
	indexPath := path.Join(baseDir, indexFile)
	data, err := os.ReadFile(indexPath)
	if err != nil {
		zerolog.Print("Error reading the index log", err)
		return model.Index{}
	}
	index, err := model.IndexFromFile(baseDir, data)
	if err != nil {
		zerolog.Print("Error parsing the index log", err)
		return model.Index{}
	}
	return index
}

func (m *LogService) ReadDay(date time.Time) (model.DailyLog, error) {
	dateString := timeconv.TimeToDayString(date)
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
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Print("Error reading the file")
		log.Print(err.Error())
		return model.NewDailyLog(date, m.baseDir), nil
	}
	return model.DailyFrom(file, dateTime, date, m.baseDir)
}

func (m *LogService) AddNewLog(date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := timeconv.TimeToDayString(date)
	dailyLog, _ := m.cache[dateString]
	dailyLog.Logs = append(dailyLog.Logs, model.NewLog(name, category))
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) AppendNewLog(uuid string, date time.Time, name string, category model.Category) (model.DailyLog, error) {
	dateString := timeconv.TimeToDayString(date)
	dailyLog, _ := m.cache[dateString]
	for index, appendLog := range dailyLog.Logs {
		if appendLog.Id == uuid {
			dailyLog.Logs[index].AppendNewSubLog(name, category)
		}
	}
	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) MoveExistingLog(date time.Time, previousLog model.Log) (model.DailyLog, error) {
	dateString := timeconv.TimeToDayString(date)
	dailyLog, _ := m.cache[dateString]
	if previousLog.IsComplete() || previousLog.IsMigrated() || previousLog.IsIrrelevant() {
		if previousLog.SubLogs != nil {
			for _, item := range *previousLog.SubLogs {
				if item.IsATask() {
					dailyLog.Logs = append(dailyLog.Logs, item)
				}
			}
		}
	} else {
		// Only keep the task logs.
		if previousLog.SubLogs != nil {
			var notDeletedSubLogs []model.Log
			for _, item := range *previousLog.SubLogs {
				if item.IsATask() {
					notDeletedSubLogs = append(notDeletedSubLogs, item)
				}
			}
			previousLog.SubLogs = &notDeletedSubLogs
		}
		dailyLog.Logs = append(dailyLog.Logs, previousLog)
	}

	m.cache[dateString] = dailyLog
	return m.SaveLog(date)
}

func (m *LogService) getPreviousFileName(from time.Time) (time.Time, string, error) {
	files, err := os.ReadDir(m.baseDir)
	if err != nil {
		log.Print(err)
		return time.Now(), "", err
	}
	startTime := time.Time{}
	previousFileName := ""
	actualDate := timeconv.TimeToDayString(from)
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			extension := filepath.Ext(filename)
			dateName := filename[0 : len(filename)-len(extension)]
			if dateName == actualDate {
				continue
			}
			toTime, err := timeconv.StringToDayTime(dateName)
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

func (m *LogService) GetPreviousDate(from time.Time) (model.DailyLog, error) {
	date, dateString, err := m.getPreviousFileName(from)
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
	dateString := timeconv.TimeToDayString(date)
	dailyLog, _ := m.cache[dateString]
	filePath := path.Join(m.baseDir, fmt.Sprintf("%v.yaml", dateString))
	bytes, err := dailyLog.ToBytes()
	if err != nil {
		return dailyLog, err
	}
	err = os.WriteFile(filePath, bytes, 0666)
	if err != nil {
		return dailyLog, err
	}
	return dailyLog, nil
}

func (m *LogService) SaveIndex() {
	indexPath := path.Join(m.baseDir, indexFile)
	bytes, err := m.Index.ToBytes()
	if err != nil {
		log.Print("Error converting the index log", err, indexPath)
		return
	}
	err = os.WriteFile(indexPath, bytes, 0666)
	if err != nil {
		log.Print("Error saving the index log", err, indexPath)
		return
	}

}

func (m *LogService) OpenIndexItem(index model.IndexItem) {
	// TODO: Move the editor to a config file
	err := utils.RunEditor("nvim", index.FullUrl)
	if err != nil {
		log.Print("Error opening the editor", err, index.FullUrl)
	}
}

func (m *LogService) CreateIndexItem(name string) {
	output := m.escapeName(timeconv.TimeToDayString(time.Now()), name)
	output += ".md"

	indexItem := model.NewIndexItem(name, output, m.baseDir)

	_, err := os.Create(indexItem.FullUrl)
	if err != nil {
		log.Print("Error creating the file", err, indexItem.FullUrl)
		return
	}
	m.Index.Items = append(m.Index.Items, indexItem)
	m.SaveIndex()
	m.OpenIndexItem(indexItem)
}

func (m *LogService) escapeName(time, name string) string {
	output := strings.ToUpper(fmt.Sprintf("%v_%v", time, name))
	output = strings.ReplaceAll(output, " ", "_")
	return output
}
