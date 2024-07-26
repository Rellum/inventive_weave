package types

import (
	"time"
)

type Data struct {
	Creators []struct {
		Id    string
		Email string
	}
	Products []struct {
		Id         string
		CreatorId  string
		CreateTime time.Time
	}
}
