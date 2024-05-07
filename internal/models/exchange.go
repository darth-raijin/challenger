package models

type Exchange string // Exchange represents a stock exchange.

const (
	KRAKEN  Exchange = "kraken"
	BINANCE Exchange = "binance"
	BITTREX Exchange = "bittrex"
	UNISWAP Exchange = "uniswap"
)
