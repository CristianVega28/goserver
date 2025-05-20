package models

import "time"

type (
	BaseModel struct {
		Id         uint
		Created_at time.Time
		Updated_at time.Time
	}
)

func Migration() {

}
