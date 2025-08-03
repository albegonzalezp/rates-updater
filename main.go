package main

import (
	"github.com/albegonzalezp/ratesupdater/config"
	"github.com/albegonzalezp/ratesupdater/database"
	"github.com/albegonzalezp/ratesupdater/models"
	"github.com/albegonzalezp/ratesupdater/service/rateapi"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	start := time.Now()

	// Load config according to environment.
	if err := config.LoadEnvironmentVariables(); err != nil {
		log.Fatal(err)
	}

	// Connects to db
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Migration
	if err := db.AutoMigrate(&models.Rate{}); err != nil {
		log.Fatal(err)
	}

	// Init client http
	client := http.Client{
		Timeout: time.Second * 30,
	}

	// Init exchange api service
	exSrv := rateapi.NewService(os.Getenv("EXCHANGE_API_URL"), os.Getenv("EXCHANGE_API_KEY"), &client)

	// Get EUR ->
	rts, err := exSrv.GetRates("EUR")
	if err != nil {
		log.Fatal(err)
	}

	// Convert rates from API to db rates
	ratesDb := rts.ToRateDB()

	// Get VES ->
	rts2, err := exSrv.GetRates("VES")
	if err != nil {
		log.Fatal(err)
	}

	// Append new results to ratesDB
	ratesDb = append(ratesDb, rts2.ToRateDB()...)

	// Update in db or create if not exists
	if err := db.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "base_code"}, {Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"rate", "updated_at"}),
	}).Create(&ratesDb).Error; err != nil {
		log.Fatal(err)
	}

	log.Println("Updated rates. Took:", time.Since(start))

}
