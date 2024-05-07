package main

import (
	"sync"
)

type Employee struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

// EmployeeStore holds employees and necessary synchronization primitives
type EmployeeStore struct {
	sync.RWMutex
	Employees map[int]Employee
	NextID    int
}
