package imagesearch

type Bins []Bin

func NewBinsFromValues(values []float64) Bins {
	var bins Bins
	for _, value := range values {
		bins = append(bins, Bin{value})
	}
	return bins
}

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
		bins[index].CountAsDecimal = float64(count) / float64(total)
	}

	return bins
}
