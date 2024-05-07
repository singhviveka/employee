package main

import (
	"sync"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (store *EmployeeStore) CreateEmployee(emp Employee) int {
	store.Lock()
	defer store.Unlock()
	emp.ID = store.NextID
	store.Employees[store.NextID] = emp
	store.NextID++
	return emp.ID
}

func (store *EmployeeStore) GetEmployeeByID(id int) (Employee, bool) {
	store.RLock()
	defer store.RUnlock()
	emp, exists := store.Employees[id]
	return emp, exists
}

func (store *EmployeeStore) UpdateEmployee(id int, emp Employee) bool {
	store.Lock()
	defer store.Unlock()
	if _, exists := store.Employees[id]; exists {
		emp.ID = id
		store.Employees[id] = emp
		return true
	}
	return false
}

func (store *EmployeeStore) DeleteEmployee(id int) bool {
	store.Lock()
	defer store.Unlock()
	if _, exists := store.Employees[id]; exists {
		delete(store.Employees, id)
		return true
	}
	return false
}

func main() {
	store := &EmployeeStore{Employees: make(map[int]Employee), NextID: 1}
	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			pageStr := r.URL.Query().Get("page")
			limitStr := r.URL.Query().Get("limit")
			page, _ := strconv.Atoi(pageStr)
			limit, _ := strconv.Atoi(limitStr)
			if page < 1 {
				page = 1
			}
			if limit < 1 {
				limit = 10
			}
			store.RLock()
			employees := make([]Employee, 0, len(store.Employees))
			for _, emp := range store.Employees {
				employees = append(employees, emp)
			}
			store.RUnlock()
			start := (page - 1) * limit
			if start > len(employees) {
				start = len(employees)
			}
			end := start + limit
			if end > len(employees) {
				end = len(employees)
			}
			json.NewEncoder(w).Encode(employees[start:end])
		case "POST":
			var emp Employee
			if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id := store.CreateEmployee(emp)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}


