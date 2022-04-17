package gcus

// stores a number of prices
type currencyInfoChunk struct {
	currencyName string
	prices       []float64
	pointer      int
	size         int

	averagePrice         float64
	maximumChangeInPrice float64
	changeInPrice        float64
	percentChangeInPrice float64
}

const (
	chunkSize = 5
)

var (
	chunks             map[string]*currencyInfoChunk
	previousPrediction string
)

func setupMarkov(currencies []string) {
	previousPrediction = "Generating first prediction ..."
	chunks = make(map[string]*currencyInfoChunk, len(currencies))
	for _, curr := range currencies {
		c := currencyInfoChunk{}
		c.currencyName = curr
		c.pointer = 0
		c.size = chunkSize
		c.prices = make([]float64, c.size)
		chunks[curr] = &c
	}
}

func addPriceData(name string, val float64) string {
	index := chunks[name].pointer
	chunks[name].prices[index] = val
	chunks[name].pointer++
	if chunks[name].pointer == chunks[name].size {
		chunks[name].pointer = 0
		calculateChunkProperties(chunks[name])
		previousPrediction = stateToString(determineState(chunks[name], 0.0001))
	}

	return previousPrediction
}

// calculateChunkProperties determines the properties of a chunk
func calculateChunkProperties(c *currencyInfoChunk) {
	c.changeInPrice = c.prices[len(c.prices)-1] - c.prices[0]
	max := c.prices[0]
	min := max
	var sum float64
	sum = 0
	for _, i := range c.prices {
		if i > max {
			max = i
		}

		if i < min {
			min = i
		}

		sum += i
	}

	c.maximumChangeInPrice = max - min
	c.averagePrice = sum / float64(c.size)
	c.percentChangeInPrice = c.changeInPrice / c.averagePrice
}

// determineState finds out whether c fits into one of three states-- 0: falling, 1: stagnating, 2: rising
func determineState(c *currencyInfoChunk, percentBin float64) int {
	if c.percentChangeInPrice < 0.0-percentBin/100.0/2.0 {
		return 0
	}

	if c.percentChangeInPrice > 0.0+percentBin/100.0/2.0 {
		return 2
	}

	return 1
}

// converts a state of a currencyInfoChunk into a string
func stateToString(i int) string {
	if i == 0 {
		return "Currency Falling"
	}

	if i == 1 {
		return "Currency Rising"
	}

	return "Currency Stable"
}
