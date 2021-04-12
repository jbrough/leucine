Examples
--------

These example files are the `head -50` output of `uniprot_trembl.fasta` and
`uniprot_sprot.fasta`

These are single "interleaved" format multi-fasta files. blastr uses these
files as a flat-file database after some minimal processing. See the commands
in the `makefile` for an example.

The full files can be downloaded from uniprot.org or (the far faster)
epi.ac.uk FTP server:

https://ftp.ebi.ac.uk/pub/databases/swissprot/release_compressed/

The full trEMBL file contains 200 million fastas but will process in a few
minutes. The default config splits it into 20 files of 20 million sequential
fastas (40 million lines). This is to give some organisation around parallelism
rather than because of filesystem or memory limits - everything is read as a
stream.

`$ make example`

```
rm -rf examples/generated/*
go run cmd/split.go -n=5 -in=examples/ -out=examples/generated/ | go run cmd/pretty.go
{
  "runtime_secs": 0.0002063,
  "source": "examples/uniprot_tr.fasta",
  "splits": [
    "examples/generated/uniprot_tr.1.fa",
    "examples/generated/uniprot_tr.2.fa"
  ]
}
{
  "runtime_secs": 0.00020423,
  "source": "examples/uniprot_sprot.fasta",
  "splits": [
    "examples/generated/uniprot_sprot.1.fa",
    "examples/generated/uniprot_sprot.2.fa"
  ]
}
{
  "destination": "examples/generated/",
  "runtime_secs": 0.000408089,
  "source": "examples/",
  "stats": [
    {
      "runtime_secs": 0.0002063,
      "source": "examples/uniprot_tr.fasta",
      "splits": [
        "examples/generated/uniprot_tr.1.fa",
        "examples/generated/uniprot_tr.2.fa"
      ]
    },
    {
      "runtime_secs": 0.00020423,
      "source": "examples/uniprot_sprot.fasta",
      "splits": [
        "examples/generated/uniprot_sprot.1.fa",
        "examples/generated/uniprot_sprot.2.fa"
      ]
    }
  ]
}
go run cmd/select.go -search=Frog -in=examples/generated/ -out=examples/generated/frog.fa | go run cmd/pretty.go
{
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "runtime_secs": 0,
  "source": "examples/generated/",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "runtime_secs": 0.00003762,
      "searched": 3,
      "selected": 0
    }
  ]
}
{
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "runtime_secs": 0,
  "source": "examples/generated/",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "runtime_secs": 0.00003762,
      "searched": 3,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "runtime_secs": 0.00009335,
      "searched": 5,
      "selected": 0
    }
  ]
}
{
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "runtime_secs": 0,
  "source": "examples/generated/",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "runtime_secs": 0.00003762,
      "searched": 3,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "runtime_secs": 0.00009335,
      "searched": 5,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "runtime_secs": 0.00017282,
      "searched": 2,
      "selected": 1
    }
  ]
}
{
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "runtime_secs": 0,
  "source": "examples/generated/",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "runtime_secs": 0.00003762,
      "searched": 3,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "runtime_secs": 0.00009335,
      "searched": 5,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "runtime_secs": 0.00017282,
      "searched": 2,
      "selected": 1
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.1.fa",
      "runtime_secs": 0.000310521,
      "searched": 5,
      "selected": 3
    }
  ]
}
{
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "runtime_secs": 0.000600361,
  "source": "examples/generated/",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "runtime_secs": 0.00003762,
      "searched": 3,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "runtime_secs": 0.00009335,
      "searched": 5,
      "selected": 0
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "runtime_secs": 0.00017282,
      "searched": 2,
      "selected": 1
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.1.fa",
      "runtime_secs": 0.000310521,
      "searched": 5,
      "selected": 3
    }
  ]
}
go run cmd/align.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j | go run cmd/pretty.go
{
  "qi": 171,
  "qid": "Q6GZX4",
  "qn": "sp|Q6GZX4|001R_FRG3G Putative transcription factor 001R OS=Frog virus 3 (isolate Goorha) OX=654924 GN=FV3-001R PE=4 SV=1",
  "qs": {
    "a": "QHVAAFLKELRHSKQYENVNLIHYILTDKRVDIQHLEKDLVKDFKALVESAHRMRQGHMI",
    "x": 144,
    "y": 204
  },
  "si": 425,
  "sid": "O08452",
  "sn": "tr|O08452|O08452_9EURY Alpha-amylase OS=Pyrococcus furiosus OX=2261 GN=amyA PE=3 SV=1",
  "ss": {
    "a": "WVGRWVYVPKFAGACIHEYTGNLGGWVDKRVDSSGWVYLEAPPHDPANGYYGYSVWSYCG",
    "x": 398,
    "y": 458
  },
  "w": "DKRVD"
}
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
