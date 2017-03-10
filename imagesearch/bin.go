package imagesearch

type Bins []Bin

/*
// between 0 and 1. 1 = total match, 0 = no match
func (bins Bins) CountsAsDecimals() []float64 {
	counts := make([]float64, len(bins))

	var total uint
	for _, bin := range bins {
		total += bin.count
	}

	for index, bin := range bins {
		counts[index] = float64(bin.count) / float64(total)
	}

	return counts

}
*/
type Bin struct {
	CountAsDecimal float64
}

func NewBins(counts []uint) Bins {
	var total uint
	for _, count := range counts {
		total += count
	}

	bins := make(Bins, len(counts))
	for index, count := range counts {
		bins[index].CountAsDecimal = float64(total) / float64(count)
	}

	return bins
}

/*
// between 0 and 1. 1 = total match, 0 = no match
func (bin Bin) CountsAsDecimal() []float64 {
	counts := make([]float64, len(bins))

	var total uint
	for _, bin := range bins {
		total += bin.count
	}

	for index, bin := range bins {
		counts[index] = float64(bin.count) / float64(total)
	}

	return counts

}
*/
/*
// between 0 and 1. 1 = total match, 0 = no match
func (bin *Bin) MatchQuality(otherBin *Bin) float64 {

}
*/
