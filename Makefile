unit_test:
	go test -v ./...

it_test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

it_test_down:
	docker-compose -f docker-compose.test.yml down

run_server:
	docker-compose up --build

run_server_down:
	docker-compose down