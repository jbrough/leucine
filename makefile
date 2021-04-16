.PHONY: example

test:
	go test ./... -bench=. -v

example:
	mkdir -p examples/generated
	rm -rf examples/generated/*
	go run cmd/main.go -split=5 -src=examples/ -dst=examples/generated/ import
	go run cmd/main.go -search=Frog -src=examples/generated/ -dst=examples/generated/frog.fa select
	go run cmd/main.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j search
	go run cmd/main.go -query=data/sars2.fa -candidates=examples/generated -n 5 -j search | go run cmd/main.go -min=40 -j score | go run cmd/main.go pretty
