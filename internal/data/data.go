package data

import "time"

type Model struct {
	url      string
	interval time.Duration
}

func NewModel(url string, interval time.Duration) Model {
	return Model{
		url: url,
		interval: interval,
	}
}

func (m Model) Interval() time.Duration {
	return m.interval
}

func (m Model) Url() string {
	return m.url
}
