package models

import (
	"time"
)


type Channel struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    CreatedAt    time.Time `json:"created_at"`
    TonalityTag  string    `json:"tonality_tag"`
    Description  string    `json:"description"`
}

type ChannelUser struct {
    ChannelID int64 `json:"channel_id"`
    UserID    int64 `json:"user_id"`
}