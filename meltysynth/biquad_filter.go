package meltysynth

import (
	"math"
)

var resonancePeakOffset = float32(1 - 1/math.Sqrt(2))

type biQuadFilter struct {
	synthesizer *Synthesizer
	active      bool
	a0          float32
	a1          float32
	a2          float32
	a3          float32
	a4          float32
	x1          float32
	x2          float32
	y1          float32
	y2          float32
}

func newBiQuadFilter(synthesizer *Synthesizer) *biQuadFilter {
	result := new(biQuadFilter)
	result.synthesizer = synthesizer
	return result
}

func (filter *biQuadFilter) clearBuffer() {
	filter.x1 = 0
	filter.x2 = 0
	filter.y1 = 0
	filter.y2 = 0
}

func (filter *biQuadFilter) setLowPassFilter(cutoffFrequency float32, resonance float32) {

	if cutoffFrequency < 0.499*float32(filter.synthesizer.SampleRate) {

		filter.active = true

		// This equation gives the Q value which makes the desired resonance peak.
		// The error of the resultant peak height is less than 3%.
		q := resonance - resonancePeakOffset/(1+6*(resonance-1))

		w := 2 * math.Pi * float64(cutoffFrequency) / float64(filter.synthesizer.SampleRate)
		cosw := math.Cos(w)
		alpha := math.Sin(w) / float64(2*q)

		b0 := (1 - cosw) / 2
		b1 := 1 - cosw
		b2 := (1 - cosw) / 2
		a0 := 1 + alpha
		a1 := -2 * cosw
		a2 := 1 - alpha

		filter.setCoefficients(float32(a0), float32(a1), float32(a2), float32(b0), float32(b1), float32(b2))

	} else {

		filter.active = false
	}
}

func (filter *biQuadFilter) process(block []float32) {

	blockLength := len(block)

	if filter.active {

		for t := 0; t < blockLength; t++ {

			input := block[t]
			output := filter.a0*input + filter.a1*filter.x1 + filter.a2*filter.x2 - filter.a3*filter.y1 - filter.a4*filter.y2

			filter.x2 = filter.x1
			filter.x1 = input
			filter.y2 = filter.y1
			filter.y1 = output

			block[t] = output
		}

	} else {

		filter.x2 = block[blockLength-2]
		filter.x1 = block[blockLength-1]
		filter.y2 = filter.x2
		filter.y1 = filter.x1
	}
}

func (filter *biQuadFilter) setCoefficients(a0 float32, a1 float32, a2 float32, b0 float32, b1 float32, b2 float32) {
	filter.a0 = b0 / a0
	filter.a1 = b1 / a0
	filter.a2 = b2 / a0
	filter.a3 = a1 / a0
	filter.a4 = a2 / a0
}