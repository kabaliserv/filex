dev-start:
	docker-compose -f docker/docker-compose.dev.yml up

testbuild:
	go build -o /dev/null github.com/kabaliserv/filex/cmd/kbs-filex

test:
	go test ./...