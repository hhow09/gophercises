package app

import (
	"testing"

	"github.com/hhow09/cab_booking_app/model"

	"github.com/go-playground/assert/v2"
)

func TestUpdateLocation(t *testing.T) {
	x, y := 100, 200
	app := NewApp()
	id := app.RegisterCab(x, y)
	assert.NotEqual(t, id, -1)
	err := app.UpdateCabLocation(id, 200, 300)
	assert.Equal(t, err, nil)
	updated := app.GetCab(id)
	assert.Equal(t, updated.GetLocation().X, 200)
	assert.Equal(t, updated.GetLocation().Y, 300)
}

func TestMatch(t *testing.T) {
	cabs := []model.Location{{X: 500, Y: 600}, {X: 700, Y: 800}, {X: 100, Y: 200}}
	app := NewApp()
	for _, location := range cabs {
		app.RegisterCab(location.X, location.Y)
	}
	rider1 := app.RegisterRider(0, 0)
	dest := model.Location{X: 7000, Y: 8000}
	cab1, err := app.BookCab(rider1, dest)
	assert.Equal(t, err, nil)
	assert.Equal(t, cab1, 2)
	rider2 := app.RegisterRider(1000, 800)
	cab2, err := app.BookCab(rider2, dest)
	assert.Equal(t, err, nil)
	assert.Equal(t, cab2, 1)
}
