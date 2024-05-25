build:
	go build -o bin/bloggen cmd/bloggen/*.go

run: build
	bin/bloggen create --from cmd/bloggen/markdown

preview: build
	bin/bloggen preview --from cmd/bloggen/markdown