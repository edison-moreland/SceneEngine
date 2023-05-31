package main

import (
	"time"
)

// rollingAverage is used to smooth out large spike in the elapsed render time
type rollingAverage struct {
	samples    []time.Duration
	sampleSize int
}

func (r *rollingAverage) HasSamples() bool {
	return len(r.samples) > 0
}

func (r *rollingAverage) Sample(t time.Duration) {
	if len(r.samples) == r.sampleSize {
		r.samples = r.samples[1:]
	}
	r.samples = append(r.samples, t)
}

func (r *rollingAverage) Average() time.Duration {
	average := time.Duration(0)
	count := 0
	for _, s := range r.samples {
		average += s
		count += 1
	}

	return average / time.Duration(count)
}

func (r *rollingAverage) Reset() {
	r.samples = r.samples[:0]
}
