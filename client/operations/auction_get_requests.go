package operations

import (
	"encoding/json"

	"github.com/regner/albionmarket-client/client/albionstate"
	"github.com/regner/albionmarket-client/client/uploader"
	"github.com/regner/albionmarket-client/log"
)

type AuctionGetRequestsResponse struct {
	MarketOrders []string `mapstructure:"0"`
}

func (op AuctionGetRequestsResponse) Process(state *albionstate.AlbionState) {
	log.Debug("Got response to AuctionGetOffers operation...")

	if state.LocationId == 0 {
		log.Warn("The players location has not yet been set. Pleas transition zones so the location can be identified.")
		return
	}

	orders := []*marketOrder{}

	for _, v := range op.MarketOrders {
		order := &marketOrder{}

		err := json.Unmarshal([]byte(v), order)
		if err != nil {
			log.Errorf("Problem converting market order to internal struct: %v", err)
		}

		orders = append(orders, order)
	}

	if len(orders) > 0 {
		ingestRequest := marketUpload{
			Orders:     orders,
			LocationID: state.LocationId,
		}

		data, err := json.Marshal(ingestRequest)
		if err != nil {
			log.Errorf("Error while marshalling payload for market orders: %v", err)
			return
		}

		uploader.SendToIngest([]byte(string(data)), "marketorders")
	}
}
