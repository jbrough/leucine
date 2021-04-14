Leucine Search
==============

A text search engine based on fasta files for deterministic searching of whole
sequence spaces.


Quick Start
------------

Run `make example`

This will reformat and search the UniProtKB fasta and GenBank Sequence files
into sequential fasta files, select some query proteins, and run a search.


About
-----

Leucine implements its own matching and scoring and is designed around streaming
from flat file, and it does not need a database or require any pre-indexing.

On a normal laptop you can expect to be able to run several million alignments a
second - enough to compare all human proteins against a selection of other
proteins (eg a smaller organism like a virus), in a few seconds.
 
_Do any of the subsequences in any of the 30 proteins in an organism appear
anywhere, at any position, in any human protein?_

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


TODO
----

Adding more tools that these rough results can be piped into - solvent
accessible surface area is top of the list.

Another facet to the scoring would be to take into account Amino Acid 
Exchangeability (https://www.ncbi.nlm.nih.gov/pmc/articles/PMC1449787/).

Probabilistic subsitution matrices like Blossom would be unlikely to add any
information to these kinds of queries that compare small regions of sequences,
but EX has potential.
