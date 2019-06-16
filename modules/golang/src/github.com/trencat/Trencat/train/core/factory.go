package core

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Factory provides a way to generate core structs with random (but reasonable) data
type Factory struct {
	rand *rand.Rand
}

// NewFactory declares and initialises a Factory instance with a random seed
func NewFactory() Factory {
	seed := time.Now().UnixNano()
	fmt.Printf("Factory seed set to %d", seed)
	return Factory{
		rand: rand.New(rand.NewSource(seed)),
	}
}

// NewFactorySeed declares and initialises a Factory instance with a given seed
func NewFactorySeed(seed int64) Factory {
	return Factory{
		rand: rand.New(rand.NewSource(seed)),
	}
}

// GetTrain generate a Train instance with random specifications
func (f *Factory) GetTrain() Train {
	return Train{
		ID:            f.rand.Int(),
		Length:        50 + 1e2*f.rand.Float64(),             // Range [50m, 150m]
		Mass:          3e5 + 4e5*f.rand.Float64(),            // Range [3e5, 7e5]
		MassFactor:    1 + f.rand.Float64()/5,                // Range [  1, 1.5]
		MaxForce:      2e5 + 4e5*f.rand.Float64(),            // Range [2e5, 6e5]
		MaxBrake:      2e5 + 4e5*f.rand.Float64(),            // Range [2e5, 6e5]
		ResistanceLin: 1e-9 + (1e-6-1e-9)*f.rand.Float64(),   // Range [1e-9, 1e-6]
		ResistanceQua: 3e-11 + (1e-9-1e-11)*f.rand.Float64(), // Range [3e-11, 1e-9]
	}
}

// GetTrack generates a slice of consecutive Track instances with random specifications
func (f *Factory) GetTrack(number int, minLength float64, maxLength float64,
	slope bool, bend bool, tunnel bool) ([]Track, error) {

	if number <= 0 {
		fail := fmt.Sprintf("Attempt to Factory.GetTrack. Parameter %d is negative", number)
		return nil, errors.New(fail)
	}

	railroad := make([]Track, number)

	for i := 0; i < number; i++ {
		ID := f.rand.Int()
		track := Track{
			ID:          ID,
			Length:      minLength + (maxLength-minLength)*f.rand.Float64(), // Range [minLength, maxLength]
			MaxVelocity: 5 + (30-5)*f.rand.Float64(),                        // Range [5, 30]
			BendRadius:  math.Inf(1),
		}

		if slope && f.rand.Intn(2) == 1 {
			var sign float64 = 1
			r := f.rand.Intn(2)
			if r == 0 {
				sign = -1
			}
			track.Slope = sign * 0.01 * f.rand.Float64() // Range [-0.01, 0.01]
		}

		if bend && f.rand.Intn(2) == 1 {
			track.BendRadius = 60 + (1e3-60)*f.rand.Float64() // Range [60m, 1000m]
		}

		if tunnel && f.rand.Intn(2) == 1 {
			track.Tunnel = true
		}

		railroad[i] = track
	}

	return railroad, nil
}
