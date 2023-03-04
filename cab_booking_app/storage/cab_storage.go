package storage

import (
	"errors"
	"log"

	"github.com/hhow09/cab_booking_app/model"
)

type CabStorage struct {
	count int
	dict  map[int]*model.Cab
}

func NewCabStorage() CabStorage {
	return CabStorage{count: 0, dict: map[int]*model.Cab{}}
}

func (cs *CabStorage) GetDict() map[int]*model.Cab {
	return cs.dict
}

func (cs *CabStorage) Add(el model.Cab) (int, error) {
	id := cs.count
	el.SetId(id)
	if _, exist := cs.dict[id]; exist {
		log.Fatal("value already exist")
		return -1, errors.New("value already exist")
	}
	cs.count += 1
	cs.dict[id] = &el
	return id, nil
}

func (cs *CabStorage) Get(id int) *model.Cab {
	if _, exist := cs.dict[id]; !exist {
		log.Fatal("not exist")
		return nil
	} else {
		return cs.dict[id]
	}
}
