package models

import "time"

// Order represents the structure of a trading order from any given exchange
type Order struct {
	Exchange   string      // Identifier for the exchange
	BuySymbol  string      // Trading symbol for the token being bought
	SellSymbol string      // Trading symbol for the token being bought
	Type       OrderType   // Type of the order, e.g., limit, market
	Side       OrderSide   // Buy or sell
	Quantity   float64     // Quantity of the asset to buy or sell
	Price      float64     // Price per unit for limit orders; ignored if market order
	Status     OrderStatus // Status of the order, e.g., pending, completed
	Timestamp  time.Time   // Timestamp when the order was placed or processed
}

// OrderType defines the type of the order
type OrderType int32

const (
	OrderTypeUnknown         OrderType = 0 // Should not be used
	OrderTypeMarket          OrderType = 1 // Market order
	OrderTypeLimit           OrderType = 2 // Limit order
	OrderTypeStopLoss        OrderType = 3 // Stop loss order
	OrderTypeStopLossLimit   OrderType = 4 // Stop loss limit order
	OrderTypeTakeProfit      OrderType = 5 // Take profit order
	OrderTypeTakeProfitLimit OrderType = 6 // Take profit limit order
)

// OrderSide defines which side of the trade the order is on, buy or sell
type OrderSide int32

const (
	OrderSideUndefined OrderSide = 0 // Default Value, should not be used
	OrderSideBuy       OrderSide = 1 // Buy order
	OrderSideSell      OrderSide = 2 // Sell order
)

// OrderStatus defines the current status of the order
type OrderStatus int32

const (
	OrderStatusPending         OrderStatus = 0 // Order has been placed but not processed yet
	OrderStatusFilled          OrderStatus = 1 // Order has been completely filled
	OrderStatusPartiallyFilled OrderStatus = 2 // Order has been partially filled
	OrderStatusCanceled        OrderStatus = 3 // Order has been canceled
	OrderStatusRejected        OrderStatus = 4 // Order has been rejected
	OrderStatusExpired         OrderStatus = 5 // Order has expired
)
