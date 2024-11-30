package model

import "sync"

type Value struct {
	Str string
	Dec float64
}

type Bucket struct {
	Name string
	Data []Value
}

type Snapshot map[string]Bucket

type Storage struct {
	History map[int]Snapshot
	// each element is key for History map
	Index []int
	// last snapshot number
	Counter int
	// max size
	Limit int
	Lock  sync.RWMutex
}
