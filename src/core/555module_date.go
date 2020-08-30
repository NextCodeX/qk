package core

import "time"

type Datetime struct {
	val       time.Time
	Timestamp int64
	Format    string
}
