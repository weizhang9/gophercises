package db

import (
	"encoding/binary"
	"time"
	
	"github.com/boltdb/bolt"
)

var (
	taskBucket = []byte("tasks")

	db *bolt.DB
)

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) error {
	return db.Update(func(tx *bolt.Tx) error  {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		key := itob(id64)
		return b.Put(key, []byte(task))
	})
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key: btoi(k),
				Value: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(taskBucket).Delete(itob(uint64(key)))
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}