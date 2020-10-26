package gophy

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// ProteinModel protein model struct
type ProteinModel struct {
	M DiscreteModel
}

// NewProteinModel get new PROTModel pointer
func NewProteinModel() *ProteinModel {
	CharMap := make(map[string][]int)
	CharMap["A"] = []int{0}
	CharMap["R"] = []int{1}
	CharMap["N"] = []int{2}
	CharMap["D"] = []int{3}
	CharMap["C"] = []int{4}
	CharMap["Q"] = []int{5}
	CharMap["E"] = []int{6}
	CharMap["G"] = []int{7}
	CharMap["H"] = []int{8}
	CharMap["I"] = []int{9}
	CharMap["L"] = []int{10}
	CharMap["K"] = []int{11}
	CharMap["M"] = []int{12}
	CharMap["F"] = []int{13}
	CharMap["P"] = []int{14}
	CharMap["S"] = []int{15}
	CharMap["T"] = []int{16}
	CharMap["W"] = []int{17}
	CharMap["Y"] = []int{18}
	CharMap["V"] = []int{19}
	CharMap["-"] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	CharMap["X"] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	CharMap["B"] = []int{2, 3}
	CharMap["Z"] = []int{5, 6}
	dnam := ProteinModel{}
	dnam.M.Alph = AminoAcid
	dnam.M.NumStates = 20
	dnam.M.CharMap = CharMap
	return &dnam
}

// SetRateMatrixJTT set up JTT exchangeabilities
func (d *ProteinModel) SetRateMatrixJTT() {
	d.M.Ex = `58
	54  45
	81  16 528
	56 113  34  10
	57 310  86  49   9
	105  29  58 767   5 323
	179 137  81 130  59  26 119
	27 328 391 112  69 597  26  23
	36  22  47  11  17   9  12   6  16
	30  38  12   7  23  72   9   6  56 229
	35 646 263  26   7 292 181  27  45  21  14
	54  44  30  15  31  43  18  14  33 479 388  65
	15   5  10   4  78   4   5   5  40  89 248   4  43
	194  74  15  15  14 164  18  24 115  10 102  21  16  17
	378 101 503  59 223  53  30 201  73  40  59  47  29  92 285
	475  64 232  38  42  51  32  33  46 245  25 103 226  12 118 477
	9 126   8   4 115  18  10  55   8   9  52  10  24  53   6  35  12
	11  20  70  46 209  24   7   8 573  32  24   8  18 536  10  63  21  71
	298  17  16  31  62  20  45  47  11 961 180  14 323  62  23  38 112  25  16

	0.076748 0.051691 0.042645 0.051544 0.019803 0.040752 0.061830 0.073152 0.022944 0.053761 0.091904 0.058676 0.023826 0.040126 0.050901 0.068765 0.058565 0.014261 0.032102 0.066005`
	// from paml jones.dat

	scanner := bufio.NewScanner(strings.NewReader(d.M.Ex))
	res := [][]float64{}
	for scanner.Scan() {
		tmp := []float64{}
		line := strings.Fields(scanner.Text())
		for _, j := range line {
			f, err := strconv.ParseFloat(j, 64)
			if err != nil {
				log.Fatal(err)
			}
			tmp = append(tmp, f)
		}
		res = append(res, tmp)
	}
	d.M.R = mat.NewDense(20, 20, nil) // exchangeabilities - let diagonals be 0
	for i := 0; i < d.M.NumStates; i++ {
		for j := 0; j < d.M.NumStates; j++ {
			if i == j {
				d.M.R.Set(i, j, 0.0)
			} else if j > i {
				d.M.R.Set(i, j, res[j-1][i])
			} else {
				d.M.R.Set(i, j, res[i-1][j])
			}
		}
	}
	d.M.MBF = res[20]
}

// SetRateMatrixWAG set up WAG exchangeabilities
func (d *ProteinModel) SetRateMatrixWAG() {
	d.M.Ex = `0.551571
	0.509848  0.635346
	0.738998  0.147304  5.429420
	1.027040  0.528191  0.265256  0.0302949
	0.908598  3.035500  1.543640  0.616783  0.0988179
	1.582850  0.439157  0.947198  6.174160  0.021352  5.469470
	1.416720  0.584665  1.125560  0.865584  0.306674  0.330052  0.567717
	0.316954  2.137150  3.956290  0.930676  0.248972  4.294110  0.570025  0.249410
	0.193335  0.186979  0.554236  0.039437  0.170135  0.113917  0.127395  0.0304501 0.138190
	0.397915  0.497671  0.131528  0.0848047 0.384287  0.869489  0.154263  0.0613037 0.499462  3.170970
	0.906265  5.351420  3.012010  0.479855  0.0740339 3.894900  2.584430  0.373558  0.890432  0.323832  0.257555
	0.893496  0.683162  0.198221  0.103754  0.390482  1.545260  0.315124  0.174100  0.404141  4.257460  4.854020  0.934276
	0.210494  0.102711  0.0961621 0.0467304 0.398020  0.0999208 0.0811339 0.049931  0.679371  1.059470  2.115170  0.088836  1.190630
	1.438550  0.679489  0.195081  0.423984  0.109404  0.933372  0.682355  0.243570  0.696198  0.0999288 0.415844  0.556896  0.171329  0.161444
	3.370790  1.224190  3.974230  1.071760  1.407660  1.028870  0.704939  1.341820  0.740169  0.319440  0.344739  0.967130  0.493905  0.545931  1.613280
	2.121110  0.554413  2.030060  0.374866  0.512984  0.857928  0.822765  0.225833  0.473307  1.458160  0.326622  1.386980  1.516120  0.171903  0.795384  4.378020
	0.113133  1.163920  0.0719167 0.129767  0.717070  0.215737  0.156557  0.336983  0.262569  0.212483  0.665309  0.137505  0.515706  1.529640  0.139405  0.523742  0.110864
	0.240735  0.381533  1.086000  0.325711  0.543833  0.227710  0.196303  0.103604  3.873440  0.420170  0.398618  0.133264  0.428437  6.454280  0.216046  0.786993  0.291148  2.485390
	2.006010  0.251849  0.196246  0.152335  1.002140  0.301281  0.588731  0.187247  0.118358  7.821300  1.800340  0.305434  2.058450  0.649892  0.314887  0.232739  1.388230  0.365369  0.314730

	0.0866279 0.043972  0.0390894 0.0570451 0.0193078 0.0367281 0.0580589 0.0832518 0.0244313 0.048466  0.086209  0.0620286 0.0195027 0.0384319 0.0457631 0.0695179 0.0610127 0.0143859 0.0352742 0.0708956`
	// from PAML wag.dat

	scanner := bufio.NewScanner(strings.NewReader(d.M.Ex))
	res := [][]float64{}
	for scanner.Scan() {
		tmp := []float64{}
		line := strings.Fields(scanner.Text())
		for _, j := range line {
			f, err := strconv.ParseFloat(j, 64)
			if err != nil {
				log.Fatal(err)
			}
			tmp = append(tmp, f)
		}
		res = append(res, tmp)
	}
	d.M.R = mat.NewDense(20, 20, nil) // exchangeabilities - let diagonals be 0
	for i := 0; i < d.M.NumStates; i++ {
		for j := 0; j < d.M.NumStates; j++ {
			if i == j {
				d.M.R.Set(i, j, 0.0)
			} else if j > i {
				d.M.R.Set(i, j, res[j-1][i])
			} else {
				d.M.R.Set(i, j, res[i-1][j])
			}
		}
	}
	d.M.MBF = res[20]
}

// SetRateMatrixLG set up LG exchangeabilities
func (d *ProteinModel) SetRateMatrixLG() {
	d.M.Ex = `0.425093
	0.276818 0.751878
	0.395144 0.123954 5.076149
	2.489084 0.534551 0.528768 0.062556
	0.969894 2.807908 1.695752 0.523386 0.084808
	1.038545 0.363970 0.541712 5.243870 0.003499 4.128591
	2.066040 0.390192 1.437645 0.844926 0.569265 0.267959 0.348847
	0.358858 2.426601 4.509238 0.927114 0.640543 4.813505 0.423881 0.311484
	0.149830 0.126991 0.191503 0.010690 0.320627 0.072854 0.044265 0.008705 0.108882
	0.395337 0.301848 0.068427 0.015076 0.594007 0.582457 0.069673 0.044261 0.366317 4.145067
	0.536518 6.326067 2.145078 0.282959 0.013266 3.234294 1.807177 0.296636 0.697264 0.159069 0.137500
	1.124035 0.484133 0.371004 0.025548 0.893680 1.672569 0.173735 0.139538 0.442472 4.273607 6.312358 0.656604
	0.253701 0.052722 0.089525 0.017416 1.105251 0.035855 0.018811 0.089586 0.682139 1.112727 2.592692 0.023918 1.798853
	1.177651 0.332533 0.161787 0.394456 0.075382 0.624294 0.419409 0.196961 0.508851 0.078281 0.249060 0.390322 0.099849 0.094464
	4.727182 0.858151 4.008358 1.240275 2.784478 1.223828 0.611973 1.739990 0.990012 0.064105 0.182287 0.748683 0.346960 0.361819 1.338132
	2.139501 0.578987 2.000679 0.425860 1.143480 1.080136 0.604545 0.129836 0.584262 1.033739 0.302936 1.136863 2.020366 0.165001 0.571468 6.472279
	0.180717 0.593607 0.045376 0.029890 0.670128 0.236199 0.077852 0.268491 0.597054 0.111660 0.619632 0.049906 0.696175 2.457121 0.095131 0.248862 0.140825
	0.218959 0.314440 0.612025 0.135107 1.165532 0.257336 0.120037 0.054679 5.306834 0.232523 0.299648 0.131932 0.481306 7.803902 0.089613 0.400547 0.245841 3.151815
	2.547870 0.170887 0.083688 0.037967 1.959291 0.210332 0.245034 0.076701 0.119013 10.649107 1.702745 0.185202 1.898718 0.654683 0.296501 0.098369 2.188158 0.189510 0.249313

	0.079066 0.055941 0.041977 0.053052 0.012937 0.040767 0.071586 0.057337 0.022355 0.062157 0.099081 0.064600 0.022951 0.042302 0.044040 0.061197 0.053287 0.012066 0.034155 0.069147`
	// lg.dat from paml

	scanner := bufio.NewScanner(strings.NewReader(d.M.Ex))
	res := [][]float64{}
	for scanner.Scan() {
		tmp := []float64{}
		line := strings.Fields(scanner.Text())
		for _, j := range line {
			f, err := strconv.ParseFloat(j, 64)
			if err != nil {
				log.Fatal(err)
			}
			tmp = append(tmp, f)
		}
		res = append(res, tmp)
	}
	d.M.R = mat.NewDense(20, 20, nil) // exchangeabilities - let diagonals be 0
	for i := 0; i < d.M.NumStates; i++ {
		for j := 0; j < d.M.NumStates; j++ {
			if i == j {
				d.M.R.Set(i, j, 0.0)
			} else if j > i {
				d.M.R.Set(i, j, res[j-1][i])
			} else {
				d.M.R.Set(i, j, res[i-1][j])
			}
		}
	}
	d.M.MBF = res[20]
}

// SetupQProt set up Q matrix
// This is scaled (so the branch lengths are going to be proportional to these changes)
// Use SetRateMatrix* and then do this
/*func (d *ProteinModel) SetupQProt() {
	bigpi := mat.NewDiagDense(d.M.NumStates, d.M.BF)
	dQ := mat.NewDense(d.M.NumStates, d.M.NumStates, nil)
	dQ.Mul(d.M.R, bigpi)
	// le and gascuel 08 account of scaling. Gives correct likelihood but for overlong branch lengths
	// var diagSum float64
	// for i := 0; i < d.NumStates; i++ {
	// 	dQ.Set(i, i, 0-sumRow(dQ, i))
	// }
	// for i := 0; i < d.NumStates; i++ {
	// 	diagSum += dQ.At(i, i)
	// }
	// diagSum = -diagSum
	// d.Q = mat.NewDense(d.NumStates, d.NumStates, nil)
	// for i := 0; i < d.NumStates; i++ {
	// 	for j := 0; j < d.NumStates; j++ {
	// 		d.Q.Set(i, j, dQ.At(i, j)/diagSum)
	// 	}
	// }

	// foster-style scaling
	var offdSum float64
	for i := 0; i < d.M.NumStates; i++ {
		for j := 0; j < d.M.NumStates; j++ {
			if i != j {
				offdSum += dQ.At(i, j)
			}
		}
	}
	var dSum float64
	d.M.Q = mat.NewDense(d.M.NumStates, d.M.NumStates, nil)
	for i := 0; i < d.M.NumStates; i++ {
		d.M.Q.Set(i, i, 0-sumRow(dQ, i))
		dSum += d.M.Q.At(i, i) * d.M.BF[i]
	}
	for i := 0; i < d.M.NumStates; i++ {
		for j := 0; j < d.M.NumStates; j++ {
			if i == j {
				d.M.Q.Set(i, i, d.M.Q.At(i, i)/-dSum)
			} else {
				d.M.Q.Set(i, j, dQ.At(i, j)/-dSum)
			}
		}
	}
}
*/
