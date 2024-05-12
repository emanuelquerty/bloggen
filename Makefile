build:
	go build -o bin/bloggen cmd/*.go

run: build
	bin/bloggen create --from cmd/markdown

preview: build
	bin/bloggen preview --from cmd/markdown