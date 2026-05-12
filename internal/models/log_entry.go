package models

type LogEntry struct {
	Level string `json:"level" binding:"required"`
	Timestamp string `json:"ts" binding:"required"`
	Caller string `json:"caller" binding:"required"`
	Message string `json:"msg" binding:"required"`
}
