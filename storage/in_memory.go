package storage

import "GOFolder/models"

var (
	Users       = make(map[int]models.User)
	Messages    = make(map[int]models.Message)
	UserCounter = 1
	MsgCounter  = 1
)
