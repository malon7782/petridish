package main

type Event struct {
	Day int
	Msg string
}

type Logger struct {
	Events []Event
	Num    int
}

func (l *Logger) Add(day int, msg string) {
	l.Events = append(l.Events, Event{ Day: day, Msg: msg })
}
