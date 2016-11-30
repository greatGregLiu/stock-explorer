package main

import (
	"encoding/json"
	"fmt"

	"time"

	"./finance/yahoo"
	emitter "github.com/emitter-io/go"
)

func main() {
	provider := yahoo.NewProvider()
	o := emitter.NewClientOptions()

	// Set the message handler
	o.SetOnMessageHandler(func(c emitter.Emitter, msg emitter.Message) {
		fmt.Printf("Received message: %s %v\n", msg.Payload(), msg.Topic())

		// Parse the request
		var request map[string]string
		if err := json.Unmarshal(msg.Payload(), &request); err != nil {
			fmt.Println("Error: Unable to parse the request")
		}

		quotes, err := provider.GetQuotes(request["symbol"])
		if err != nil {
			fmt.Println("Error: Unable to process the request")
		}

		response, _ := json.Marshal(quotes[0])
		c.Publish("aST2oXP-iDd09T-dumFL8_GIBf-oTvOw", "quote-response/"+request["reply"], response)
	})

	// Create a new emitter client and connect to the broker
	c := emitter.NewClient(o)
	sToken := c.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		panic("Error on Client.Connect(): " + sToken.Error().Error())
	}

	// Subscribe to the request channel
	c.Subscribe("FKLs16Vo7W4RjYCvU86Nk0GvHNi5AK8t", "quote-request")

	fmt.Println("Service running...")
	for {
		time.Sleep(100)
	}

	//symbols := []string{"WSR", "APLE", "BRG", "EPR", "OAKS", "GOOD"}
	/*for _, quote := range quotes {
		fmt.Printf("%v\n", quote.Name)
		fmt.Printf("Last year dividend frequency: %v\n", quote.GetLastYearDividendFrequency())
		fmt.Printf("Growth: %v, %v\n", quote.GetGrowth(), quote.GetRevenueGrowth())
		fmt.Printf("Profitability: %v, %v\n", quote.GetProfitability(), quote.GetFFOGrowth())

		//spew.Dump(quote)

	}*/

}
