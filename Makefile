run:
	go run .

build:
	go build

upload:
	rsync -auv . uber:src/albumbot-go
