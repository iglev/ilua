all: mod test

mod:
	#go mod vendor -v
	go mod tidy -v

test:
	go test

