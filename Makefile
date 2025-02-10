dev:
	wgo run ./cmd/main.go
gen-docs:
	swag init -g swagger.go -d ./internal -o ./docs/api
gen-mock:
	mockery
test:
	export TESTCONTAINERS_RYUK_DISABLED=true && go test ./... -coverprofile=cover.out && go tool cover -html=cover.out
show-coverage:
	go tool cover -html=cover.out