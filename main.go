package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/albegonzalezp/ratesupdater/config"
	"github.com/albegonzalezp/ratesupdater/database"
	"github.com/albegonzalezp/ratesupdater/models"
	"github.com/albegonzalezp/ratesupdater/service/rateapi"
	"gorm.io/gorm/clause"
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

	// Get USD ->
	rts3, err := exSrv.GetRates("USD")
	if err != nil {
		log.Fatal(err)
	}

	// Append new results to ratesDB
	ratesDb = append(ratesDb, rts2.ToRateDB()...)
	ratesDb = append(ratesDb, rts3.ToRateDB()...)

	// Update in db or create if not exists
	if err := db.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "base_code"}, {Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"rate", "updated_at"}),
	}).Create(&ratesDb).Error; err != nil {
		log.Fatal(err)
	}

	// Get rates from p2p
	rtsp2p, err := exSrv.GetUSDTtoVESp2pRate()
	if err != nil {
		log.Fatal(err)
	}

	// base code: USD; code: VES; rate: 0.0004..
	rateP2p := rtsp2p.ToRateDB()
	if rateP2p.Rate != 0 {
		ratesDb = append(ratesDb, rateP2p)
	}

	for _, rate := range ratesDb {

		// If EUR -> VES we only update the payment method Dolares en Efectivo and TRansferencia bancaria
		if rate.BaseCode == "EUR" && rate.Code == "VES" {

			log.Println("Iterating over EUR/VES")

			// get the rate EUR -> USD
			var rateEurUsd models.Rate
			if err := db.Db.Where("base_code = ? AND code = ?", "EUR", "USD").First(&rateEurUsd).Error; err != nil {
				log.Fatal(err.Error())
			}

			log.Println("Current EUR/USD rate from db: ", rateEurUsd.Rate)

			log.Println("Rate USDT/VES from p2p: ", rateP2p.Rate)

			// get the rate USD -> VES from p2p
			rt := rateEurUsd.Rate * rateP2p.Rate

			log.Println("Rate final for EUR/VES: ", rt)

			// Update the rates table with the calculated EUR/VES rate
			if err := db.Db.Model(&models.Rate{}).
				Where("base_code = ? AND code = ?", "EUR", "VES").
				Update("rate", rt).Error; err != nil {
				log.Fatal(err.Error())
			}

			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? AND (payment_method_type = ? or payment_method_type = ?)", 1, 3, false, 2, 1).
				Update("rate", rt).Error; err != nil {
				log.Fatal(err.Error())
			}

		}

		// If EUR -> USD update the rate for all payment methods
		if rate.BaseCode == "EUR" && rate.Code == "USD" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 2, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		if rate.BaseCode == "EUR" && rate.Code == "COP" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 4, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		if rate.BaseCode == "EUR" && rate.Code == "PEN" {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 5, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		if rate.BaseCode == "VES" && rate.Code == "USD" {

			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 3, 2, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		if rate.BaseCode == "VES" && rate.Code == "EUR" {

			rate := 1 / rateP2p.Rate

			log.Println("Rate for VES/EUR: ", rate)

			// get the rate USD -> EUR
			var rateUsdEur models.Rate
			if err := db.Db.Where("base_code = ? AND code = ?", "USD", "EUR").First(&rateUsdEur).Error; err != nil {
				log.Fatal(err.Error())
			}

			rate = rate / rateUsdEur.Rate

			log.Println("Final rate for VES/EUR: ", rate)

			// Update the rates table with the calculated VES/EUR rate
			if err := db.Db.Model(&models.Rate{}).
				Where("base_code = ? AND code = ?", "VES", "EUR").
				Update("rate", rate).Error; err != nil {
				log.Fatal(err.Error())
			}

			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 3, 1, false).
				Update("rate", rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

	}

	log.Println("Updated rates. Took:", time.Since(start))

}
