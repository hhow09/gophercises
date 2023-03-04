package app

import (
	"errors"
	"log"
	"time"

	"github.com/hhow09/cab_booking_app/model"
	"github.com/hhow09/cab_booking_app/storage"
)

type App struct {
	cabStorage    storage.CabStorage
	riderStorage  storage.RiderStorage
	matchStrategy model.MatchStrategyInterface
	tripStorage   storage.TripStorage
}

func (a *App) RegisterCab(x, y int) int {
	cab := model.NewCab(x, y)
	id, err := a.cabStorage.Add(cab)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return id
}

func (a *App) GetCab(id int) model.Cab {
	cab := a.cabStorage.Get(id)
	return *cab
}

func (a *App) RegisterRider(x, y int) int {
	rider := model.NewRider(x, y)
	id, err := a.riderStorage.Add(rider)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return id
}

func (a *App) UpdateCabLocation(id int, x, y int) error {
	cab := a.cabStorage.Get(id)
	if cab == nil {
		return errors.New("cab not exist")
	}
	cab.UpdateLocation(x, y)
	return nil
}

func (a *App) BookCab(riderId int, destination model.Location) (int, error) {
	rider := a.riderStorage.Get(riderId)
	if rider == nil {
		return -1, errors.New("cannot find rider")
	}
	cab, err := a.matchStrategy.Match(rider.Location, a.cabStorage.GetDict())
	if err != nil {
		return -1, err
	}
	// TODO create trip
	t := model.NewTrip(*rider, *cab, destination)
	a.tripStorage.Add(t)
	// save to trip storage

	go func() {
		time.Sleep(time.Second * 1)
		t.SetStatus(model.FINISHED)
	}()

	return cab.Id, nil
}

func WithStrategy(s model.MatchStrategyInterface) func(app *App) {
	return func(app *App) {
		app.matchStrategy = s
	}
}

func NewApp(options ...func(*App)) App {
	app := App{cabStorage: storage.NewCabStorage(), riderStorage: storage.NewRiderStorage(), tripStorage: storage.NewTripStorage(), matchStrategy: &model.DefaultMatchStrategy{}}
	return app
}
