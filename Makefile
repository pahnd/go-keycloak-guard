build:
	go build -o bin/main main.go

run:
	go run main.go


compile:
	# Compiling for Linux x86 64 bit
	go mod tidy && GOOS=linux GOARCH=amd64 go build -o bin/keycloak-guard main.go

docker-compile:
	docker-compose --profile compile up

kong-start:
	docker-compose --profile compile up
	docker-compose --profile kong up


kong-stop:
	docker-compose down

docker-clean-up:
	docker-compose down --volumes --remove-orphans
	docker volume prune -f
	docker network prune -f
	docker container prune -f
	docker image prune -a -f

release:
	mkdir -p bin/linux_amd64
	mkdir -p bin/linux_arm64
	go mod tidy && GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/keycloak-guard main.go
	go mod tidy && GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64/keycloak-guard main.go
	cp schema.lua bin/linux_amd64/schema.lua
	cp schema.lua bin/linux_arm64/schema.lua
	tar -cvzf bin/linux_amd64.tar.gz bin/linux_amd64
	tar -cvzf bin/linux_arm64.tar.gz bin/linux_arm64
	rm -rf bin/linux_amd64
	rm -rf bin/linux_arm64
