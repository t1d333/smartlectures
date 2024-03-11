.PHONY: build down


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
