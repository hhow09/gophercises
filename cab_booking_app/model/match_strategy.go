package model

import "math"

type MatchStrategyInterface interface {
	Match(start Location, candidateMap map[int]*Cab) (*Cab, error)
}

type MatchStrategy struct {
}

type DefaultMatchStrategy struct {
}

func calcDist(la, lb Location) float64 {
	return math.Sqrt(math.Pow(float64(la.X-lb.X), 2) + math.Pow(float64(la.Y-lb.Y), 2))
}

func (s *DefaultMatchStrategy) Match(start Location, candidateMap map[int]*Cab) (*Cab, error) {
	min_dist := math.MaxFloat64
	var matchedCab *Cab
	for _, candidate := range candidateMap {
		dist := calcDist(start, candidate.GetLocation())
		if dist < float64(min_dist) {
			min_dist = dist
			matchedCab = candidate
		}
	}
	return matchedCab, nil
}
