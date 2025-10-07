package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/zenixls2/alpacaapi"
)

func main() {
	alpacaapi.DebugRequestResponse = true
	client := alpacaapi.NewClient()
	id := strings.TrimSpace(os.Getenv("APCA_API_KEY_ID"))
	key := strings.TrimSpace(os.Getenv("APCA_API_SECRET_KEY"))
	fmt.Println("ID:", strings.TrimSpace(id))
	fmt.Println("Key:", strings.TrimSpace(key))
	client.SetAuthByAPIKey(id, key)
	/*req := client.AccountService.NewGetAccountRequest()
	resp, err := req.Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)*/
	req := client.OrderService.NewCreateOrderRequest()
	req.SetSymbol("AAPL")
	req.SetQty(fixedpoint.NewFromFloat(1))
	req.SetSide("buy")
	req.SetType("limit")
	req.SetTimeInForce("day")
	req.SetOrderClass("simple")
	req.SetLimitPrice(fixedpoint.NewFromFloat(170))
	/*req.SetLegs([]alpacaapi.Leg{{
	Side:           "buy",
	Symbol:         "AAPL",
	RatioQty:       1,
	PositionIntent: "buy_to_open"}})*/
	req.SetPositionIntent("buy_to_open")
	// print json
	resp, err := req.Do(context.Background())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println(resp)
}
