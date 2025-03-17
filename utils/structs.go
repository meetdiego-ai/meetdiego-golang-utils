package utils

import "time"

type TaskItem struct {
	ID     string `json:"id"`
	Value  string `json:"value"`
	TaskId string `json:"taskId"`
	Status string `json:"status"`
}

type Task struct {
	UserId    string     `json:"userId"`
	Type      string     `json:"type"`
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Label     string     `json:"label"`
	TaskItems []TaskItem `json:"taskItems"`
}
