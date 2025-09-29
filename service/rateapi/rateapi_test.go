package rateapi

import (
	"net/http"
	"testing"
)

func TestP2pExchange(t *testing.T) {

	srv := Service{
		Client: &http.Client{},
	}

	res, err := srv.GetUSDTtoVESp2pRate()
	if err != nil {
		t.Fatal(err)
	}

	rate := res.ToRateDB()

	t.Log(rate)
}
