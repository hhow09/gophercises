package storage

import (
	"github.com/hhow09/cab_booking_app/model"
)

type TripStorage struct {
	dict map[int][]model.Trip // [userId] -> trip
}

func (ts *TripStorage) Add(t model.Trip) {
	rider_id := t.GetRider().Id
	ts.dict[rider_id] = append(ts.dict[rider_id], t)
}

func NewTripStorage() TripStorage {
	return TripStorage{map[int][]model.Trip{}}
}
