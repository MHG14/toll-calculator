OBU:
	@go build -o bin/OBU OBU/main.go
	@./bin/OBU

reciever:
	@go build -o bin/reciever ./data_reciever
	@./bin/reciever

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

.PHONY: OBU