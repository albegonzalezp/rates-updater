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

	// Get rates

	for _, rate := range ratesDb {

		// If EUR -> VES we only update the payment method Dolares en Efectivo
		if rate.BaseCode == "EUR" && rate.Code == "VES" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? AND payment_method = ?", 1, 3, false, 3).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err)
			}
		}

		// If EUR -> USD update the rate for all payment methods
		if rate.BaseCode == "EUR" && rate.Code == "USD" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 2, false).
				Update("rate", rate.Rate); err != nil {
				log.Fatal(err)
			}
		}

		if rate.BaseCode == "EUR" && rate.Code == "COL" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 4, false).
				Update("rate", rate.Rate); err != nil {
				log.Fatal(err)
			}
		}

		if rate.BaseCode == "EUR" && rate.Code == "PEN" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 5, false).
				Update("rate", rate.Rate); err != nil {
				log.Fatal(err)
			}
		}

	}

	log.Println("Updated rates. Took:", time.Since(start))

}
