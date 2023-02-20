package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLastDate(t *testing.T) {
	a := assert.New(t)
	d, _ := getLastDay(2023, 2)
	a.Equal(d, 28)
}

func RandomWeather() Weather {
	rand.Seed(time.Now().UnixNano())
	return Weather{
		High: 20 + rand.Float64()*(40-20),
		Low:  10 + rand.Float64()*(40-10),
		Rain: rand.Float64(),
	}
}

// all success
func TestWeathersSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ret := RandomWeather()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
	}))
	defer ts.Close() // shutdown mock server after test finished

	res, err := quoteApi(ts.URL, 2022, time.April)
	ass := assert.New(t)
	ass.Nil(err)
	ass.Equal(len(res), 30)
	for i := range res {
		ass.IsType(Weather{}, res[i])
	}
}

func TestWeathersPartialSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ret := RandomWeather()
		w.Header().Set("Content-Type", "application/json")
		// intenionally return error
		if strings.Contains(r.URL.RawQuery, "10") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(struct{}{})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ret)
	}))
	defer ts.Close() // shutdown mock server after test finished

	res, err := quoteApi(ts.URL, 2022, time.April)
	ass := assert.New(t)
	ass.Equal(len(res), 0)
	ass.EqualError(err, "incomplete data")

}
