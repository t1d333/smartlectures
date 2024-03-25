.PHONY: build down prod recognizer


build:
	@docker compose up --build -d
	
down:
	@docker compose down

prod:
	@mv ./docker-compose.yml ./docker-compose-dev.yml 
	@cp ./prod/docker-compose.yml ./docker-compose.yml
	@docker compose up --build -d
	@rm ./docker-compose.yml
	@mv ./docker-compose-dev.yml ./docker-compose.yml  

swagger:
	@docker compose up swagger -d
	@docker compose up nginx -d

lint:
	@golangci-lint run ./...

recognizer:
	@python3 -m grpc_tools.protoc -Irecognizer=protos \
  	--python_out=. --pyi_out=. --grpc_python_out=. \
  	./protos/recognizer.proto
