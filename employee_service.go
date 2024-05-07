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
