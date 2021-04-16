package lib

import (
	"errors"
	"strings"
)

const AACSV = `
Alanine;Ala;A;Aliphatic;Nonpolar;Neutral;1.8;;;89.094;8.76;GCN
Arginine;Arg;R;Basic;Basic polar;Positive;−4.5;;;174.203;5.78;MGR, CGY (coding codons can also be expressed by: CGN, AGR)
Asparagine;Asn;N;Amide;Polar;Neutral;−3.5;;;132.119;3.93;AAY
Aspartic acid;Asp;D;Acid;Acidic polar;Negative;−3.5;;;133.104;5.49;GAY
Cysteine;Cys;C;Sulfuric;Nonpolar;Neutral;2.5;250;0.3;121.154;1.38;UGY
Glutamine;Gln;Q;Amide;Polar;Neutral;−3.5;;;146.146;3.9;CAR
Glutamic acid;Glu;E;Acid;Acidic polar;Negative;−3.5;;;147.131;6.32;GAR
Glycine;Gly;G;Aliphatic;Nonpolar;Neutral;−0.4;;;75.067;7.03;GGN
Histidine;His;H;Basic aromatic;Basic polar;Positive, 10% Neutral, 90%;−3.2;211;5.9;155.156;2.26;CAY
Isoleucine;Ile;I;Aliphatic;Nonpolar;Neutral;4.5;;;131.175;5.49;AUH
Leucine;Leu;L;Aliphatic;Nonpolar;Neutral;3.8;;;131.175;9.68;YUR, CUY (coding codons can also be expressed by: CUN, UUR)
Lysine;Lys;K;Basic;Basic polar;Positive;−3.9;;;146.189;5.19;AAR
Methionine;Met;M;Sulfuric;Nonpolar;Neutral;1.9;;;149.208;2.32;AUG
Phenylalanine;Phe;F;Aromatic;Nonpolar;Neutral;2.8;257, 206, 188;0.2, 9.3, 60.0;165.192;3.87;UUY
Proline;Pro;P;Cyclic;Nonpolar;Neutral;−1.6;;;115.132;5.02;CCN
Serine;Ser;S;Hydroxylic;Polar;Neutral;−0.8;;;105.093;7.14;UCN, AGY
Threonine;Thr;T;Hydroxylic;Polar;Neutral;−0.7;;;119.119;5.53;ACN
Tryptophan;Trp;W;Aromatic;Nonpolar;Neutral;−0.9;280, 219;5.6, 47.0;204.228;1.25;UGG
Tyrosine;Tyr;Y;Aromatic;Polar;Neutral;−1.3;274, 222, 193;1.4, 8.0, 48.0;181.191;2.91;UAY
Valine;Val;V;Aliphatic;Nonpolar;Neutral;4.2;;;117.148;6.73;GUN`

type AminoAcids struct {
	aa map[string]AminoAcid
	a3 map[string]string
}

func (a AminoAcids) GetAAA(code string) (*AminoAcid, error) {
	c, ok := a.a3[code]
	if !ok {
		return nil, errors.New("NOT FOUND: AminoAcids.GetAAA() called with " + code)
	}
	r, _ := a.aa[c]
	return &r, nil
}

func (a AminoAcids) GetA(code string) (*AminoAcid, error) {
	r, ok := a.aa[code]
	if !ok {
		return nil, errors.New("NOT FOUND: AminoAcids.GetA() called with " + code)
	}
	return &r, nil
}

type AminoAcid struct {
	Name string
	A    string
	AAA  string
}

func AA() *AminoAcids {
	a := AminoAcids{
		make(map[string]AminoAcid),
		make(map[string]string),
	}

	for i, row := range strings.Split(AACSV, "\n") {
		if i == 0 {
			continue
		}
		v := strings.Split(row, ";")
		a3 := strings.ToUpper(v[1])
		aa := AminoAcid{
			Name: v[0], A: v[2], AAA: a3,
		}
		a.a3[a3] = v[2]
		a.aa[v[2]] = aa
	}

	return &a
}
