docker-mysql: ##build MySQL database docker image named chat
	@docker build -t mysql:chat mysql/.
docker-redis: ##build Redis docker image named chat
	@docker build -t redis:chat redis/.
docker-server: ##build Server docker image named chat
	@docker build -t chat33:1.0.0 service/.

docker-image: docker-mysql docker-redis docker-server

docker-compose: ##
	-@docker network create chatnet
	@docker-compose up -d --build

docker-compose-stop:
	@docker-compose stop

docker-compose-down:
	@docker-compose down
