package model

import (
	"github.com/apoloa/bjournal/src/utils"
)

type Log struct {
	Parent    *Log     `json:"-" yaml:"-"`
	Id        string   `json:"-" yaml:"-"`
	Name      string   `json:"name" yaml:"name"`
	Mark      Category `json:"mark" yaml:"mark"`
	Important bool     `json:"important" yaml:"important"`
	Url       *string  `json:"url,omitempty" yaml:"url,omitempty"`
	Text      *string  `json:"-" yaml:"-"`
	SubLogs   *[]Log   `json:"sub_logs,omitempty" yaml:"subLogs,omitempty"`
}

func NewLog(name string, category Category) Log {
	return Log{
		Name:      name,
		Mark:      category,
		Important: false,
		Url:       nil,
		Text:      nil,
	}
}

func (l *Log) GetName() string {
	if l.Mark == Irrelevant {
		return utils.Strikethrough(l.Name)
	}
	return l.Name
}

func (l *Log) AppendNewSubLog(name string, category Category) {
	if l.SubLogs == nil {
		l.SubLogs = &[]Log{}
	}
	*l.SubLogs = append(*l.SubLogs, NewLog(name, category))
}

func (l *Log) MarkAsComplete() {
	if l.Mark == Task {
		l.Mark = Complete
	}
}

func (l *Log) MarkAsIrrelevant() {
	if l.Mark == Task {
		l.Mark = Irrelevant
	}
}

func (l *Log) MarkAsMigrated() {
	if l.Mark == Task {
		l.Mark = Migrated
		if l.SubLogs != nil {
			for i, _ := range *l.SubLogs {
				(*l.SubLogs)[i].MarkAsMigrated()
			}
		}
	}
}

func (l *Log) IsATask() bool {
	return l.Mark == Task
}

func (l *Log) IsComplete() bool {
	return l.Mark == Complete
}

func (l *Log) IsMigrated() bool {
	return l.Mark == Migrated
}
