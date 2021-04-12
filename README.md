Leucine Search
==============

A text search engine based on fasta files for deterministic matching of whole
sequence spaces in ways that are easy to adapt and reason about.
 
The idea behind Luecine is toolkit of small programs that kind be arranged into
pipelines to make it easy to search large datasets without needing anything
other than ubiquitous `fasta` files, the format that all biological sequence
information is shared in.

Flat files are faster to read than a database, and can be searched without any
pre-indexing.

To get working instantly with single files containing hundreds of millions of
protein sequences on an averergae laptop, you need to stream your data, so the
toolkit makes much use of unix pipes to combine processing steps:

```
go run cmd/align.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j | go run cmd/local.go | go run cmd/pretty.go 
{
  "align": [
    "Query  144  QHVAAFLKELRHSKQYENVNLIHYILTDKRVDIQHLEKDLVKDFKALVESAHRMRQGHMI  204",
    "                                       DKRVD                                 ",
    "Sbjct  398  WVGRWVYVPKFAGACIHEYTGNLGGWVDKRVDSSGWVYLEAPPHDPANGYYGYSVWSYCG  458"
  ],
  "query": "sp|Q6GZX4|001R_FRG3G Putative transcription factor 001R OS=Frog virus 3 (isolate Goorha) OX=654924 GN=FV3-001R PE=4 SV=1",
  "sbjct": "tr|O08452|O08452_9EURY Alpha-amylase OS=Pyrococcus furiosus OX=2261 GN=amyA PE=3 SV=1"
}
```

At the moment, it is only capable of finding exact matches, such as

_Do any of the subsequences in any of the 30 proteins in an organism appear
anywhere, at any position, in any human protein?_

but can do this in a few seconds. It's a proof-of-concept with a view to
incorporating common BLAST features like heuristics and subsitution matrices
into processing pipelines.
