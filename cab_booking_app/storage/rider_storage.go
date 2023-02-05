package storage

import (
	"errors"
	"log"

	"github.com/hhow09/cab_booking_app/model"
)

type RiderStorage struct {
	count int
	dict  map[int]*model.Rider
}

func NewRiderStorage() RiderStorage {
	return RiderStorage{count: 0, dict: map[int]*model.Rider{}}
}

func (gs *RiderStorage) Add(el model.Rider) (int, error) {
	id := gs.count
	el.SetId(id)
	if _, exist := gs.dict[id]; exist {
		log.Fatal("value already exist")
		return -1, errors.New("value already exist")
	}
	gs.count += 1
	gs.dict[id] = &el
	return id, nil
}

func (gs *RiderStorage) Get(id int) *model.Rider {
	if _, exist := gs.dict[id]; !exist {
		log.Fatal("not exist")
		return nil
	} else {
		return gs.dict[id]
	}
}
