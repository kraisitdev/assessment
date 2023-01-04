u-test:
	go test -v --tags=unit ./...

it-test-up:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

it-test-down:
	docker-compose -f docker-compose.test.yml down
