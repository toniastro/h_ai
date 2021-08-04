build:
	docker-compose build
up:
	docker-compose up
stop:
	docker-compose stop
test:
	cd ./api; \
	go test -v; \
	cd ../

init: build up