.PHONY: example

example:
	rm -rf examples/generated/*
	go run cmd/split.go -n=5 -in=examples/ -out=examples/generated/ | jq
	go run cmd/select.go -search=Frog -in=examples/generated/ -out=examples/generated/frog.fa | jq
	go run cmd/align.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j | jq
