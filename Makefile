run: stop up

up:
	docker-compose -f docker-compose.yaml up -d --build

stop:
	docker-compose -f docker-compose.yaml stop

down:
	docker-compose -f docker-compose.yaml down

test:
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yaml down --volumes
