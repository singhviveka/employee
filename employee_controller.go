import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

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
