protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/api/v1/user.proto

go build .

<!-- grpc -->
<!-- server -->
go run main.go server

<!-- client -->
go run main.go client '{"Name":"adi","Id":1}'


<!-- rest -->
http://localhost:8080/createTodo

