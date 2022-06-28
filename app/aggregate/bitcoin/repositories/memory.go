package repositories

import (
	"github.com/F0rzend/simple-go-webserver/app/aggregate/bitcoin/entity"
)

var _ entity.BTCRepository = &MemoryBTCRepository{}

type MemoryBTCRepository struct {
	bitcoin entity.BTCPrice
}

func NewMemoryBTCRepository(initialPrice entity.USD) (*MemoryBTCRepository, error) {
	btcPrice := entity.NewBTCPrice(initialPrice)

	return &MemoryBTCRepository{
		bitcoin: btcPrice,
	}, nil
}

func (r *MemoryBTCRepository) Get() entity.BTCPrice {
	return r.bitcoin
}

func (r *MemoryBTCRepository) SetPrice(price entity.USD) error {
	btcPrice := entity.NewBTCPrice(price)

	r.bitcoin = btcPrice
	return nil
}
