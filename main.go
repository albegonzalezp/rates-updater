package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/albegonzalezp/ratesupdater/config"
	"github.com/albegonzalezp/ratesupdater/consts"
	"github.com/albegonzalezp/ratesupdater/database"
	"github.com/albegonzalezp/ratesupdater/models"
	"github.com/albegonzalezp/ratesupdater/service/rateapi"
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

	// Get rates from p2p
	rtsp2p, err := exSrv.GetUSDTtoVESp2pRate()
	if err != nil {
		log.Fatal(err)
	}

	// base code: USD; code: VES; rate: 0.0004..
	rateP2p := rtsp2p.ToRateDB()

	// Append new results to ratesDB
	ratesDb = append(ratesDb, rts2.ToRateDB()...)
	ratesDb = append(ratesDb, rts3.ToRateDB()...)
	ratesDb = append(ratesDb, rateP2p)

	for _, rate := range ratesDb {

		// If EUR -> VES we only update the payment method Pago movil and Transferencia bancaria
		// With Binance P2P rate
		if rate.BaseCode == consts.EURISO && rate.Code == consts.VESISO {

			// get the rate EUR -> USD
			var rateEurUsd models.Rate
			if err := db.Db.Where("base_code = ? AND code = ?", consts.EURISO, consts.USDISO).First(&rateEurUsd).Error; err != nil {
				log.Fatal(err.Error())
			}

			// get the rate USD -> VES from p2p
			rt := rateEurUsd.Rate * rateP2p.Rate

			// Update the rates table with the calculated EUR/VES rate
			if err := db.Db.Model(&models.Rate{}).
				Where("base_code = ? AND code = ?", consts.EURISO, consts.VESISO).
				Update("rate", rt).Error; err != nil {
				log.Fatal(err.Error())
			}

			// EUR/VES Transferencia y Pago Movil
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? AND (payment_method_type = ? or payment_method_type = ?)",
					consts.EUR,
					consts.VES,
					false,
					consts.EurVesBankTransfer,
					consts.EurVesMobilePay).
				Update("rate", rt).Error; err != nil {
				log.Fatal(err.Error())
			}

		}

		// If EUR -> USD update the rate for all payment methods and update payment method Dolar cash with rate EUR/USD
		if rate.BaseCode == consts.EURISO && rate.Code == consts.USDISO {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", consts.EUR, consts.USD, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}

			// Update EUR/VES Dolar Cash with EUR/USD rate
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? and payment_method_type = ?", consts.EUR, consts.VES, false, consts.EurUsdDolarCash).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}

		}

		// EUR/COP
		if rate.BaseCode == consts.EURISO && rate.Code == consts.COPISO {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 4, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		// EUR/PEN
		if rate.BaseCode == consts.EURISO && rate.Code == consts.PENISO {
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ?", 1, 5, false).
				Update("rate", rate.Rate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		// VES/USD, Zelle
		if rate.BaseCode == consts.VESISO && rate.Code == consts.USDISO {

			/*
				Para calcular tasa Base en  este par seria: 1 DIVIDIDO ENTRE la tasa del API de binance del par USDT/VES (sell).
				que hoy marca  290 VES y para hoy seria: 1/ 290: 0.003448
			*/

			finalRate := 1 / rateP2p.Rate

			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? and payment_method_type = ?", consts.VES, consts.USD, false, consts.VesUsdZelle).
				Update("rate", finalRate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		// VES/EUR, Bank Transfer
		if rate.BaseCode == consts.VESISO && rate.Code == consts.EURISO {

			/*
				Para calcular tasa Base en  este par seria: 1 DIVIDIDO ENTRE la tasa base de EUR/VES(TRANFERENCIA) que tenemos. hoy seria: 1/ 337.6512: 0.002961.
				que también nos deje ajustar margen y comisión para obtener la tasa final.
			*/

			rate := models.PaymentMethod{}
			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("payment_method_type = ? and is_custom = ?", consts.EurVesBankTransfer, false).
				First(&rate).Error; err != nil {
				log.Fatal(err.Error())
			}

			finalRate := 1 / rate.Rate
			// Update the rates table with the calculated VES/EUR rate
			if err := db.Db.Model(&models.Rate{}).
				Where("base_code = ? AND code = ?", consts.VESISO, consts.EURISO).
				Update("rate", finalRate).Error; err != nil {
				log.Fatal(err.Error())
			}

			if err := db.Db.Model(&models.PaymentMethod{}).
				Where("currency_from = ? AND currency_to = ? AND is_custom = ? and payment_method_type = ?", consts.VES, consts.EUR, false, consts.VesEurBankTransfer).
				Update("rate", finalRate).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

	}

	log.Println("Updated rates. Took:", time.Since(start))

}
