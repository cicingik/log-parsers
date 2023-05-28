
.PHONY: run
run:  ## Run applicaton (via go run main.go)
	@go run main.go -d ./metrics/csv/ -t csv --startTime 2022-01-01T00:00:00.00Z --endTime 2022-02-28T00:00:00.52Z
