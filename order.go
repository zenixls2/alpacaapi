package alpacaapi

import (
	"time"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/requestgen"
)

//go:generate -command GetRequest requestgen -method GET -responseType .APIResponse
//go:generate -command PostRequest requestgen -method POST -responseType .APIResponse

type OrderService struct {
	client *RestClient
}

type OrderType string

const (
	OrderTypeMarket       OrderType = "market"
	OrderTypeLimit        OrderType = "limit"
	OrderTypeStop         OrderType = "stop"
	OrderTypeStopLimit    OrderType = "stop_limit"
	OrderTypeTrailingStop OrderType = "trailing_stop"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

type OrderTimeInForce string

const (
	// for equity, options
	OrderTimeInForceDay OrderTimeInForce = "day"
	// for equity
	OrderTimeInForceGTC OrderTimeInForce = "gtc"
	// for equity
	OrderTimeInForceOPG OrderTimeInForce = "opg"
	// for equity
	OrderTimeInForceIOC OrderTimeInForce = "ioc"
	// for equity
	OrderTimeInForceFOK OrderTimeInForce = "fok"
	// for crypto
	OrderTimeInForceCLS OrderTimeInForce = "cls"
)

type OrderClass string

const (
	OrderClassSimple  OrderClass = "simple"
	OrderClassBracket OrderClass = "bracket"
	OrderClassOCO     OrderClass = "oco"
	OrderClassOTO     OrderClass = "oto"
	OrderClassMLEG    OrderClass = "mleg"
	OrderClassDefault OrderClass = ""
)

type PositionIntent string

const (
	PositionIntentBuy2Open   PositionIntent = "buy_to_open"
	PositionIntentBuy2Close  PositionIntent = "buy_to_close"
	PositionIntentSell2Open  PositionIntent = "sell_to_open"
	PositionIntentSell2Close PositionIntent = "sell_to_close"
)

type OrderStatus string

const (
	OrderStatusNew                OrderStatus = "new"
	OrderStatusPartiallyFilled    OrderStatus = "partially_filled"
	OrderStatusFilled             OrderStatus = "filled"
	OrderStatusDoneForDay         OrderStatus = "done_for_day"
	OrderStatusCanceled           OrderStatus = "canceled"
	OrderStatusExpired            OrderStatus = "expired"
	OrderStatusReplaced           OrderStatus = "replaced"
	OrderSpacePendingCancel       OrderStatus = "pending_cancel"
	OrderStatusPendingReplace     OrderStatus = "pending_replace"
	OrderStatusAccepted           OrderStatus = "accepted"
	OrderStatusPendingNew         OrderStatus = "pending_new"
	OrderStatusAcceptedForBidding OrderStatus = "accepted_for_bidding"
	OrderStatusStopped            OrderStatus = "stopped"
	OrderStatusRejected           OrderStatus = "rejected"
	OrderStatusSuspended          OrderStatus = "suspended"
	OrderStatusCalculated         OrderStatus = "calculated"
)

type Leg struct {
	Side           OrderSide        `json:"side" param:"side"`
	PositionIntent PositionIntent   `json:"position_intent,omitempty" param:"position_intent,omitempty"`
	Symbol         string           `json:"symbol,required" param:"symbol,required"`
	RatioQty       fixedpoint.Value `json:"ratio_qty,string,required" param:"ratio_qty,required"`
}

type AdvancedInstruction struct {
	Algorithm     OrderAlgorithm   `json:"algorithm,omitempty"`
	Destination   Exchange         `json:"destination,omitempty"`
	DisplayQty    fixedpoint.Value `json:"display_qty,string,omitempty"`
	StartTime     time.Time        `json:"start_time,omitempty"`
	EndTime       time.Time        `json:"end_time,omitempty"`
	maxPercentage fixedpoint.Value `json:"max_percentage,string,omitempty"`
}

type OrderAlgorithm string

const (
	OrderAlgorithmVWAP OrderAlgorithm = "VWAP"
	OrderAlgorithmTWAP OrderAlgorithm = "TWAP"
	OrderAlgorithmTMA  OrderAlgorithm = "TMA"
)

type Exchange string

const (
	ExchangeNYSE   Exchange = "NYSE"
	ExchangeNASDAQ Exchange = "NASDAQ"
	ExchangeARCA   Exchange = "ARCA"
)

type TakeProfit struct {
	LimitPrice           fixedpoint.Value    `json:"limit_price,string"`
	StopLoss             fixedpoint.Value    `json:"stop_loss,string"`
	PositionIntent       string              `json:"position_intent,omitempty"`
	AdvancedInstructions AdvancedInstruction `json:"advanced_instructions,omitempty"`
}

type StopLoss struct {
	StopPrice  fixedpoint.Value `json:"stop_price,string" param:"stop_price,string"`
	LimitPrice fixedpoint.Value `json:"limit_price,string,omitempty" param:"limit_price,string,omitempty"`
}

// reference: https://docs.alpaca.markets/reference/postorder
type Order struct {
	ID             string           `json:"id"`
	ClientOrderID  string           `json:"client_order_id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      *time.Time       `json:"updated_at"`
	SubmittedAt    *time.Time       `json:"submitted_at"`
	FilledAt       *time.Time       `json:"filled_at"`
	ExpiredAt      *time.Time       `json:"expired_at"`
	CanceledAt     *time.Time       `json:"canceled_at"`
	FailedAt       *time.Time       `json:"failed_at"`
	ReplacedAt     *time.Time       `json:"replaced_at"`
	ReplacedBy     string           `json:"replaced_by"`
	Replaces       string           `json:"replaces"`
	AssetID        string           `json:"asset_id"`
	Symbol         string           `json:"symbol"`
	AssetClass     string           `json:"asset_class"`
	Notional       fixedpoint.Value `json:"notional,string"`
	Qty            fixedpoint.Value `json:"qty,string"`
	FilledQty      fixedpoint.Value `json:"filled_qty,string"`
	FilledAvgPrice fixedpoint.Value `json:"filled_avg_price,string"`
	OrderClass     string           `json:"order_class"`
	OrderType      string           `json:"order_type"`
	Type           OrderType        `json:"type"`
	Side           OrderSide        `json:"side"`
	TimeInForce    string           `json:"time_in_force"`
	LimitPrice     fixedpoint.Value `json:"limit_price,string"`
	StopPrice      fixedpoint.Value `json:"stop_price,string"`
	Status         OrderStatus      `json:"status"`
	ExtendedHours  bool             `json:"extended_hours"`
	Legs           []Order          `json:"legs"`
	TrailPercent   fixedpoint.Value `json:"trail_percent,string,omitempty"`
	TrailPrice     fixedpoint.Value `json:"trail_price,string,omitempty"`
	Hwm            fixedpoint.Value `json:"hwm,string,omitempty"`
	PositionIntent string           `json:"position_intent"`
}

//go:generate requestgen -method POST -url "/v2/orders" -type CreateOrderRequest -responseType Order
type CreateOrderRequest struct {
	client               requestgen.AuthenticatedAPIClient
	Symbol               string               `json:"symbol" param:"symbol"`
	Qty                  *fixedpoint.Value    `json:"qty,string,omitempty" param:"qty,omitempty"`
	Notional             *fixedpoint.Value    `json:"notional,string,omitempty" param:"notional,omitempty"`
	Side                 OrderSide            `json:"side,required" param:"side,required" validValues:"buy,sell"`
	Type                 OrderType            `json:"type,required" param:"type,required" validValues:"market,limit,stop,stop_limit,trailing_stop"`
	TimeInForce          OrderTimeInForce     `json:"time_in_force,required" param:"time_in_force,required" validValues:"day,gtc,opg,ioc,fok,cls"`
	LimitPrice           *fixedpoint.Value    `json:"limit_price,string,omitempty" param:"limit_price,omitempty"`
	StopPrice            *fixedpoint.Value    `json:"stop_price,string,omitempty" param:"stop_price,omitempty"`
	TrailPrice           *fixedpoint.Value    `json:"trail_price,string,omitempty" param:"trail_price,omitempty"`
	TrailPercent         *fixedpoint.Value    `json:"trail_percent,string,omitempty" param:"trail_percent,omitempty"`
	ExtendedHours        bool                 `json:"extended_hours,omitempty" param:"extended_hours,omitempty"`
	ClientOrderID        *string              `json:"client_order_id,omitempty" param:"client_order_id,omitempty"`
	OrderClass           OrderClass           `json:"order_class,omitempty" param:"order_class,omitempty" validValues:"simple,bracket,oco,oto,mleg"`
	Legs                 []Leg                `json:"legs,omitempty" param:"legs,omitempty"`
	TakeProfit           *TakeProfit          `json:"take_profit,omitempty" param:"take_profit,omitempty"`
	StopLoss             *StopLoss            `json:"stop_loss,omitempty" param:"stop_loss,omitempty"`
	PositionIntent       *PositionIntent      `json:"position_intent,omitempty" param:"position_intent,omitempty" validValues:"buy_to_open,buy_to_close,sell_to_open,sell_to_close"`
	AdvancedInstructions *AdvancedInstruction `json:"advanced_instructions,omitempty" param:"advanced_instructions,omitempty"`
}

func (s *OrderService) NewCreateOrderRequest() *CreateOrderRequest {
	return &CreateOrderRequest{client: s.client}
}
