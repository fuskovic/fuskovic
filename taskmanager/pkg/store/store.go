package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type (
	// Store describes the behavior of the store.
	Store interface {
		Add(t Task) error
		Remove(ID int) error
		Do(ID int) error
		List() (Tasks, error)
		Completed() (Tasks, error)
		Close()
	}
	store struct{ db *bolt.DB }
)

// New returns a Store interface
func New() (Store, error) {
	db, err := bolt.Open(dbPath(), 0600, nil)
	if err != nil {
		return nil, err
	}

	s := &store{db: db}
	if err := s.InitTasks(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) InitTasks() error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	if _, err := tx.CreateBucketIfNotExists([]byte("tasks")); err != nil {
		return err
	}

	return tx.Commit()
}

// Add adds a new task to incompleted tasks.
func (s *store) Add(t Task) error {
	if err := s.db.View(func(tx *bolt.Tx) error {
		getResult := tx.Bucket([]byte("tasks")).Get([]byte(t.Description))
		if getResult != nil || len(getResult) > 0 {
			fmt.Println("reached bool check block")
			return fmt.Errorf("%s already exists", t.Description)
		}
		return nil
	}); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		tasksBucket := tx.Bucket([]byte("tasks"))

		ID, _ := tasksBucket.NextSequence()
		t.ID = int(ID)

		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return tasksBucket.Put([]byte(strconv.Itoa(t.ID)), data)
	})
}

// Remove removes an existing task from incompleted tasks.
func (s *store) Remove(ID int) error {
	if err := s.db.View(func(tx *bolt.Tx) error {
		getResult := tx.Bucket([]byte("tasks")).Get([]byte(strconv.Itoa(ID)))
		if getResult == nil || len(getResult) == 0 {
			return fmt.Errorf("%d does not exist", ID)
		}
		return nil
	}); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("tasks")).Delete([]byte(strconv.Itoa(ID)))
	})
}

// List lists all incompleted tasks.
func (s *store) List() (Tasks, error) {
	var incompleteTasks Tasks

	err := s.db.View(func(tx *bolt.Tx) error {
		tasksBucket := tx.Bucket([]byte("tasks"))
		cursor := tasksBucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var task Task

			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}

			if !task.IsCompleted {
				incompleteTasks = append(incompleteTasks, task)
			}
		}
		return nil
	})

	return incompleteTasks, err
}

// Do marks a task as completed.
func (s *store) Do(ID int) error {
	var task Task

	if err := s.db.View(func(tx *bolt.Tx) error {
		getResult := tx.Bucket([]byte("tasks")).Get([]byte(strconv.Itoa(ID)))

		if getResult == nil || len(getResult) == 0 {
			return fmt.Errorf("%d does not exist", ID)
		}

		if err := json.Unmarshal(getResult, &task); err != nil {
			return err
		}

		if task.IsCompleted {
			return fmt.Errorf("%d is already completed", ID)
		}
		return nil
	}); err != nil {
		return err
	}

	task.IsCompleted = true

	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("tasks")).Put([]byte(strconv.Itoa(ID)), data)
	})
}

// Completed lists all completed tasks.
func (s *store) Completed() (Tasks, error) {
	var completedTasks Tasks

	err := s.db.View(func(tx *bolt.Tx) error {
		tasksBucket := tx.Bucket([]byte("tasks"))
		cursor := tasksBucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var task Task

			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}

			if task.IsCompleted {
				completedTasks = append(completedTasks, task)
			}
		}
		return nil
	})

	return completedTasks, err
}

// Close terminates the connection to boltDB.
func (s *store) Close() {
	s.db.Close()
}

func dbPath() string {
	wd, _ := os.Getwd()
	return path.Join(wd, "boltDB/task.db")
}
