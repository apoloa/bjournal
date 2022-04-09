package model

type Log struct {
	Name      string   `yaml:"name"`
	Mark      Category `yaml:"mark"`
	Important bool     `yaml:"important"`
	Url       *string  `yaml:"url,omitempty"`
	Text      *string  `yaml:"-"`
	SubLogs   *[]Log   `yaml:"subLogs,omitempty"`
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
