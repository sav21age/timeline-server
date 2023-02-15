package domain

import "time"

type Pagination struct {
	Skip  int
	Limit int
}

type Datination struct {
	MinDate time.Time
	MaxDate time.Time
}