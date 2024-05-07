# employee api in Go with pagination and concurrency.
For running the server run command:  go run main.go</br>
You hit the Post api and insert the data as follows:</br>
curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","position":"Developer","salary":50000}' http://localhost:8080/employees</br>
Get the data with pagination by running http://localhost:8080/employees?page=1&limit=10 on browser.</br>


