package model

import "time"

type MonthlyLog struct {
	Days map[string][]Log `yaml:"days"`
	Logs []Log            `yaml:"items"`
}

func NewMonthlyLog(date time.Time) MonthlyLog {
	return MonthlyLog{}
}
