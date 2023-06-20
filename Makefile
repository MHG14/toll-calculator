OBU:
	@go build -o bin/OBU OBU/main.go
	@./bin/OBU

reciever:
	@go build -o bin/reciever data_reciever/main.go
	@./bin/reciever

.PHONY: OBU