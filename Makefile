OBU:
	@go build -o bin/OBU OBU/main.go
	@./bin/OBU

reciever:
	@go build -o bin/reciever ./data_reciever
	@./bin/reciever

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator


aggregator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

.PHONY: OBU aggregator