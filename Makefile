all: mod test

mod:
	go mod vendor -v

test:
	go test

