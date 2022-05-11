up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
	cd ./src/saiEthManager && go build -o ../../microservices/saiEthManager/build/sai-eth-manager
	cd ./src/saiGNMonitor && go build -o ../../microservices/saiGNMonitor/build/sai-gn-monitor
	cp ./src/saiEthManager/config/config.json ./microservices/saiEthManager/build/config.json
	cp ./src/saiGNMonitor/config/config.json ./microservices/saiGNMonitor/build/config.json

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

logs:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logn:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-gn-monitor