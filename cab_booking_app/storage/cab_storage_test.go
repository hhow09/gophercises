package storage

import (
	"testing"

	"github.com/hhow09/cab_booking_app/model"
	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	x, y := 100, 200
	cab := model.NewCab(x, y)
	cs := NewCabStorage()
	id, err := cs.Add(cab)
	assert.Nil(t, err)
	c := cs.Get(id)
	assert.Equal(t, (*c).Id, id)
	assert.Equal(t, c.GetLocation().X, x)
	assert.Equal(t, c.GetLocation().Y, y)
}
