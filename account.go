package alpacaapi

import (
	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/requestgen"
)

//go:generate -command GetRequest requestgen -method GET -responseType .APIResponse
//go:generate -command PostRequest requestgen -method POST -responseType .APIResponse

type AccountService struct {
	client *RestClient
}

// reference: https://docs.alpaca.markets/reference/getaccount-1
type Account struct {
	ID                       string           `json:"id"`
	AccountNumber            string           `json:"account_number"`
	Status                   AccountStatus    `json:"status"`
	Currency                 string           `json:"currency"`
	Cash                     fixedpoint.Value `json:"cash"`
	PortfolioValue           fixedpoint.Value `json:"portfolio_value"`
	NonMarginableBuyingPower fixedpoint.Value `json:"non_marginable_buying_power"`
	AccruedFees              fixedpoint.Value `json:"accrued_fees"`
	PendingTransferIn        fixedpoint.Value `json:"pending_transfer_in"`
	PendingTransferOut       fixedpoint.Value `json:"pending_transfer_out"`
	PatternDayTrader         bool             `json:"pattern_day_trader"`
	TradeSuspendedByUser     bool             `json:"trade_suspended_by_user"`
	TradingBlocked           bool             `json:"trading_blocked"`
	TransfersBlocked         bool             `json:"transfers_blocked"`
	AccountBlocked           bool             `json:"account_blocked"`
	CreatedAt                string           `json:"created_at"`
	ShortingEnabled          bool             `json:"shorting_enabled"`
	LongMarketValue          fixedpoint.Value `json:"long_market_value"`
	ShortMarketValue         fixedpoint.Value `json:"short_market_value"`
	Equity                   fixedpoint.Value `json:"equity"`
	LastEquity               fixedpoint.Value `json:"last_equity"`
	Multiplier               fixedpoint.Value `json:"multiplier"`
	BuyingPower              fixedpoint.Value `json:"buying_power"`
	InitialMargin            fixedpoint.Value `json:"initial_margin"`
	MaintenanceMargin        fixedpoint.Value `json:"maintenance_margin"`
	SMA                      fixedpoint.Value `json:"sma"`
	DaytradeCount            int              `json:"daytrade_count"`
	BalanceAsof              string           `json:"balance_asof"`
	LastMaintenanceMargin    fixedpoint.Value `json:"last_maintenance_margin"`
	DaytradingBuyingPower    fixedpoint.Value `json:"daytrading_buying_power"`
	RegtBuyingPower          fixedpoint.Value `json:"regt_buying_power"`
	OptionsBuyingPower       fixedpoint.Value `json:"options_buying_power"`
	OptionsApprovedLevel     OptionsLevel     `json:"options_approved_level"`
	OptionsTradingLevel      OptionsLevel     `json:"options_trading_level"`
	IntradayAdjustment       fixedpoint.Value `json:"intraday_adjustment"`
	PendingRegTafFees        fixedpoint.Value `json:"pending_reg_taf_fees"`
}

//go:generate GetRequest -url "/v2/account" -type GetAccountRequest -responseType Account
type GetAccountRequest struct {
	client requestgen.AuthenticatedAPIClient
}

func (s *AccountService) NewGetAccountRequest() *GetAccountRequest {
	return &GetAccountRequest{
		client: s.client,
	}
}
