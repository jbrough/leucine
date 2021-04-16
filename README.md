Leucine Search
==============

A text search engine based on fasta files for deterministic searching of whole
sequence spaces.

Current performance exceeds NCBI Blast out the box for searches against ~500k
proteins and scales linearly. There is no pre-indexing, the only setup cost is
downloading the fasta files you'd like to searcch.

Everything is a stream from querying the fastas to outputting the results, and
all outputs can be piped from one process to another, eg to score and format
matches.


Quick Start
------------

Run `make example`

This will reformat and search the UniProtKB fasta and GenBank Sequence files
into sequential fasta files, select some query proteins, and run a search.


About
-----

Leucine implements its own matching and scoring and is designed around streaming
from flat file. It does not need a database or require any pre-indexing.

An example of some output:
 
```
{
  "align": [
    "Query   155  AIVLQLPQGTTLPKGFYAEGSRGGSQASSRSSSRSRNSSRNSTPGSSRGTSPAR   209",
    "                      T          S      SSRSSSR R SSR S P  SR            ",
    "Sbjct     8  KMASVRFMVTPTKIDDIPGLSDTSPDLSSRSSSRVRFSSRESVPETSRSEPMSE    62"
  ],
  "query": "sp|P0DTC9|NCAP_SARS2 Nucleoprotein OS=Severe acute respiratory syndrome coronavirus 2 OX=2697049 GN=N PE=1 SV=1",
  "sbjct": "sp|Q9UHW9|S12A6_HUMAN Solute carrier family 12 member 6 OS=Homo sapiens OX=9606 GN=SLC12A6 PE=1 SV=2",
  "score": 49
}
```


Performance
-----------

For example, quering the 17 SARS-CoV-2 proteins (included in `data/`) against
all 390k human proteins in SwissProtKB takes 14 secs when looking for matches
that are at least 8AA long (for smaller matches of 5AA it takes about 22
seconds).

This is a linear search of all sub-sequences in all the query proteins against
all sub-sequences in the candidate proteins. For the search above, that's 1
billion alignment tests at 9 million a second, as you can see from the stats
below after a test run:

```
{
  "query": "data/sars2.fa",
  "candidates": "human/",
  "stats": {
    "sequeneces_searched": 389097,
    "alignments_found": 282,
    "alignments_tested": 114004830,
    "alginment_tests_per_sec": 9989867,
    "runtime_secs": 14.617510989,
    "stats": [
      {
        "sequeneces_searched": 9097,
        "alignments_found": 6,
        "alignments_tested": 1995209,
        "alginment_tests_per_sec": 418485,
        "runtime_secs": 4.767694313,
        "fasta_file": "human/uniprot.human.fa.20.fa"
      },
      {
        "sequeneces_searched": 20000,
        "alignments_found": 10,
        "alignments_tested": 4334219,
        "alginment_tests_per_sec": 432148,
        "runtime_secs": 10.02946162,
        "fasta_file": "human/uniprot.human.fa.19.fa"
      },
```

This search is parallelised by splitting the files in a preparatory step.

I'm currently testing this on a laptop-grade Linode, their smallest server with
dedicated CPU. It's possible I'm being allowed some burst :) But I expect
comparable performance on my chromebook - I need to test that. The only
constraints are storage and IO.


TODO
----

Adding more tools that these rough results can be piped into - solvent
accessible surface area is top of the list and in progress.

Another facet to the scoring would be to take into account Amino Acid 
Exchangeability (https://www.ncbi.nlm.nih.gov/pmc/articles/PMC1449787/).

Probabilistic subsitution matrices like Blossom would be unlikely to add any
information to these kinds of queries that compare small regions of sequences,
but EX has potential.
