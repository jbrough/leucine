.PHONY: example

example:
	mkdir -p examples/generated
	rm -rf examples/generated/*
	go run cmd/split.go -n=5 -in=examples/ -out=examples/generated/
	go run cmd/select.go -search=Frog -in=examples/generated/ -out=examples/generated/frog.fa
	go run cmd/align.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j
	go run cmd/align.go -query=data/sars2.fa -candidates=examples/generated -n 5 -j | go run cmd/score.go -min=40 | go run cmd/pretty.go
