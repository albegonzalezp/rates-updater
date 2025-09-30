package consts

const (
	EurVesBankTransfer = iota + 1
	EurUsdDolarCash
	EurVesMobilePay

	EurUsdZelle
	EurUsdBankTransfer
	EurUsdDigitalWallet

	EurCopBankTransfer
	EurCopNequi

	EurPenBankTransfer

	VesEurBankTransfer
	VesUsdZelle
)

/*
 1. EUR/VES Transferencia Bancaria (Esta OK)
 2. EUR/VES** Dolares en Efectivo (Usamos la API para el PAR EUR/USD Zelle)
 3. EUR/VES  Pago Movil (Igual que EUR/VES Transferencia ref: 337)

 4. EUR/USD Zelle
 5. EUR/USD Transferencia
 6. EUR/USD Digital Wallet

 7. EUR/COP Transferencia
 8. EUR/COP Nequi

 9. EUR/PEN Transferencia

 10. VES/EUR Transferencia (seria: 1 DIVIDIDO ENTRE la tasa base de EUR/VES(TRANFERENCIA) que tenemos. hoy seria: 1/ 337.6512: 0.002961.)
 11. VES/USD Zelle Para calcular tasa Base en  este par seria: 1 DIVIDIDO ENTRE la tasa del API de binance del par USDT/VES (sell).
que hoy marca  290 VES y para hoy seria: 1/ 290: 0.003448


*/
