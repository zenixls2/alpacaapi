package main

import (
	"context"
	"fmt"

	"github.com/zenixls2/alpacaapi"
)

func main() {
	client := alpacaapi.NewClient()
	client.SetAuthByAPIKey("AKC709A5B8381NJTN57L", "tHhGa6g0HgLdpdz3oRZrGvbAmQhfeYmZKUVeAENy")
	req := client.AccountService.NewGetAccountRequest()
	resp, err := req.Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

}
