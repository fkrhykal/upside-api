dev:
	wgo run ./cmd/main.go
gen-docs:
	swag init -g swagger.go -d ./internal -o ./docs/api
gen-mock:
	mockery