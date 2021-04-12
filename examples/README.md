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
rather than because of filesystem or limits.

`$ make example`

```
rm -rf examples/generated/*
go run cmd/split.go -n=5 -in=examples/ -out=examples/generated/ | jq
{
  "source": "examples/uniprot_tr.fasta",
  "splits": [
    "examples/generated/uniprot_tr.1.fa",
    "examples/generated/uniprot_tr.2.fa"
  ],
  "runtime_secs": 0.00017247
}
{
  "source": "examples/uniprot_sprot.fasta",
  "splits": [
    "examples/generated/uniprot_sprot.1.fa",
    "examples/generated/uniprot_sprot.2.fa"
  ],
  "runtime_secs": 0.000334381
}
{
  "source": "examples/",
  "destination": "examples/generated/",
  "stats": [
    {
      "source": "examples/uniprot_tr.fasta",
      "splits": [
        "examples/generated/uniprot_tr.1.fa",
        "examples/generated/uniprot_tr.2.fa"
      ],
      "runtime_secs": 0.00017247
    },
    {
      "source": "examples/uniprot_sprot.fasta",
      "splits": [
        "examples/generated/uniprot_sprot.1.fa",
        "examples/generated/uniprot_sprot.2.fa"
      ],
      "runtime_secs": 0.000334381
    }
  ],
  "runtime_secs": 0.000456201
}
go run cmd/select.go -search=Frog -in=examples/generated/ -out=examples/generated/frog.fa | jq
{
  "source": "examples/generated/",
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "searched": 3,
      "selected": 0,
      "runtime_secs": 0.00011169
    }
  ],
  "runtime_secs": 0
}
{
  "source": "examples/generated/",
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "searched": 3,
      "selected": 0,
      "runtime_secs": 0.00011169
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "searched": 5,
      "selected": 0,
      "runtime_secs": 0.00011386
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "searched": 2,
      "selected": 1,
      "runtime_secs": 0.00015176
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.1.fa",
      "searched": 5,
      "selected": 3,
      "runtime_secs": 0.0001562
    }
  ],
  "runtime_secs": 0
}
{
  "source": "examples/generated/",
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "searched": 3,
      "selected": 0,
      "runtime_secs": 0.00011169
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "searched": 5,
      "selected": 0,
      "runtime_secs": 0.00011386
    }
  ],
  "runtime_secs": 0
}
{
  "source": "examples/generated/",
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "searched": 3,
      "selected": 0,
      "runtime_secs": 0.00011169
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "searched": 5,
      "selected": 0,
      "runtime_secs": 0.00011386
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "searched": 2,
      "selected": 1,
      "runtime_secs": 0.00015176
    }
  ],
  "runtime_secs": 0
}
{
  "source": "examples/generated/",
  "destination": "examples/generated/frog.fa",
  "query": "Frog",
  "stats": [
    {
      "fasta_file": "examples/generated/uniprot_tr.2.fa",
      "searched": 3,
      "selected": 0,
      "runtime_secs": 0.00011169
    },
    {
      "fasta_file": "examples/generated/uniprot_tr.1.fa",
      "searched": 5,
      "selected": 0,
      "runtime_secs": 0.00011386
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.2.fa",
      "searched": 2,
      "selected": 1,
      "runtime_secs": 0.00015176
    },
    {
      "fasta_file": "examples/generated/uniprot_sprot.1.fa",
      "searched": 5,
      "selected": 3,
      "runtime_secs": 0.0001562
    }
  ],
  "runtime_secs": 0.00063067
}
go run cmd/align.go -query=examples/generated/frog.fa -candidates=examples/generated -n 5 -j | jq
{
  "sequeneces_searched": 3,
  "alignments_found": 0,
  "alignments_tested": 459,
  "alginment_tests_per_sec": 2954808,
  "runtime_secs": 0.00015534,
  "fasta_file": "examples/generated/uniprot_tr.2.fa"
}
{
  "sequeneces_searched": 1,
  "alignments_found": 0,
  "alignments_tested": 0,
  "alginment_tests_per_sec": 0,
  "runtime_secs": 9.219e-06,
  "fasta_file": "examples/generated/uniprot_sprot.1.fa"
}
{
  "sequeneces_searched": 1,
  "alignments_found": 0,
  "alignments_tested": 0,
  "alginment_tests_per_sec": 0,
  "runtime_secs": 4.25e-06,
  "fasta_file": "examples/generated/uniprot_sprot.2.fa"
}
{
  "sequeneces_searched": 1,
  "alignments_found": 0,
  "alignments_tested": 0,
  "alginment_tests_per_sec": 0,
  "runtime_secs": 4.9e-06,
  "fasta_file": "examples/generated/frog.fa"
}
{
  "qid": "Q6GZX4",
  "qi": 171,
  "cid": "O08452",
  "ci": 425,
  "w": "DKRVD"
}
{
  "qid": "Q6GZX2",
  "qi": 220,
  "cid": "C0LL04",
  "ci": 36,
  "w": "GKVPA"
}
{
  "sequeneces_searched": 5,
  "alignments_found": 2,
  "alignments_tested": 1447,
  "alginment_tests_per_sec": 2846637,
  "runtime_secs": 0.000508319,
  "fasta_file": "examples/generated/uniprot_tr.1.fa"
}
{
  "query": "examples/generated/frog.fa",
  "candidates": "examples/generated",
  "stats": {
    "sequeneces_searched": 11,
    "alignments_found": 2,
    "alignments_tested": 1906,
    "alginment_tests_per_sec": 5801445,
    "runtime_secs": 0.000682028,
    "stats": [
      {
        "sequeneces_searched": 3,
        "alignments_found": 0,
        "alignments_tested": 459,
        "alginment_tests_per_sec": 2954808,
        "runtime_secs": 0.00015534,
        "fasta_file": "examples/generated/uniprot_tr.2.fa"
      },
      {
        "sequeneces_searched": 1,
        "alignments_found": 0,
        "alignments_tested": 0,
        "alginment_tests_per_sec": 0,
        "runtime_secs": 9.219e-06,
        "fasta_file": "examples/generated/uniprot_sprot.1.fa"
      },
      {
        "sequeneces_searched": 1,
        "alignments_found": 0,
        "alignments_tested": 0,
        "alginment_tests_per_sec": 0,
        "runtime_secs": 4.25e-06,
        "fasta_file": "examples/generated/uniprot_sprot.2.fa"
      },
      {
        "sequeneces_searched": 1,
        "alignments_found": 0,
        "alignments_tested": 0,
        "alginment_tests_per_sec": 0,
        "runtime_secs": 4.9e-06,
        "fasta_file": "examples/generated/frog.fa"
      },
      {
        "sequeneces_searched": 5,
        "alignments_found": 2,
        "alignments_tested": 1447,
        "alginment_tests_per_sec": 2846637,
        "runtime_secs": 0.000508319,
        "fasta_file": "examples/generated/uniprot_tr.1.fa"
      }
    ]
  }
}
```
