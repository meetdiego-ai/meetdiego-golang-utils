package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/meetdiego-ai/meetdiego-golang-utils/utils"
)

type MySQLStorage struct {
	DB *sql.DB
}

func NewMySQLStorage() *MySQLStorage {
	// Initialize database connection pool
	db, err := sql.Open("mysql", os.Getenv("MYSQL_URI_GO"))
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 3)

	return &MySQLStorage{
		DB: db,
	}
}

func CreateTaskArtefact(uuid string, taskId string, artefactType string, artefactValue string) error {
	db := NewMySQLStorage()
	_, err := db.DB.Exec("INSERT INTO task_artefact (id, taskId, type, value) VALUES (?, ?, ?, ?)", uuid, taskId, artefactType, artefactValue)
	if err != nil {
		fmt.Printf("Error creating task artefact: %v\n", err)
		return err
	}
	return nil
}

func UpdateTaskStatus(uuid string, status string) error {
	db := NewMySQLStorage()
	_, err := db.DB.Exec("UPDATE task SET status = ? WHERE id = ? ", status, uuid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskItemStatus(uuid string, status string) error {
	db := NewMySQLStorage()
	_, err := db.DB.Exec("UPDATE task_item SET status = ? WHERE id = ? ", status, uuid)
	if err != nil {
		return err
	}
	return nil
}

func GetTaskItems(taskId string) ([]utils.TaskItem, error) {
	db := NewMySQLStorage()
	taskItems, err := db.DB.Query("SELECT id, value, taskId, status FROM task_item WHERE taskId = ?", taskId)
	if err != nil {
		return nil, err
	}
	defer taskItems.Close()

	var values []utils.TaskItem
	for taskItems.Next() {
		var value utils.TaskItem
		err = taskItems.Scan(&value.ID, &value.Value, &value.TaskId, &value.Status)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}
