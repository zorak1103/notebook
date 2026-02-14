package models

import "time"

// Config represents a key-value configuration entry
type Config struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
