package entity

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/F0rzend/simple-go-webserver/app/common"

	bitcoinEntity "github.com/F0rzend/simple-go-webserver/app/aggregate/bitcoin/entity"
)

type User struct {
	ID        uint64
	Name      string
	Username  string
	Email     *mail.Address
	Balance   Balance
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrNameEmpty = common.NewApplicationError(
		http.StatusBadRequest,
		"Name cannot be empty",
	)
	ErrUsernameEmpty = common.NewApplicationError(
		http.StatusBadRequest,
		"Username cannot be empty",
	)
	ErrInvalidEmail = common.NewApplicationError(
		http.StatusBadRequest,
		"You must provide a valid email",
	)
)

func NewUser(
	id uint64,
	name string,
	username string,
	email string,
	btcBalance float64,
	usdBalance float64,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	if username == "" {
		return nil, ErrUsernameEmpty
	}

	addr, err := ParseEmail(email)
	if err != nil {
		return nil, err
	}

	usdAmount, err := bitcoinEntity.NewUSD(usdBalance)
	if err != nil {
		return nil, err
	}

	btcAmount, err := bitcoinEntity.NewBTC(btcBalance)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Name:      name,
		Username:  username,
		Email:     addr,
		Balance:   NewBalance(usdAmount, btcAmount),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ParseEmail(email string) (*mail.Address, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrInvalidEmail
	}
	return addr, nil
}

func (u *User) ChangeUSDBalance(action bitcoinEntity.USDAction, amount bitcoinEntity.USD) error {
	switch action {
	case bitcoinEntity.DepositUSDAction:
		return u.deposit(amount)
	case bitcoinEntity.WithdrawUSDAction:
		return u.withdraw(amount)
	default:
		return bitcoinEntity.ErrInvalidUSDAction
	}
}

func (u *User) deposit(amount bitcoinEntity.USD) error {
	u.Balance.USD = u.Balance.USD.Add(amount)

	return nil
}

func (u *User) withdraw(amount bitcoinEntity.USD) error {
	if u.Balance.USD.LessThan(amount) {
		return ErrInsufficientFunds
	}

	updatedUSD, err := u.Balance.USD.Sub(amount)
	if err != nil {
		return err
	}
	u.Balance.USD = updatedUSD

	return nil
}

func (u *User) ChangeBTCBalance(action bitcoinEntity.BTCAction, amount bitcoinEntity.BTC, price bitcoinEntity.BTCPrice) error {
	switch action {
	case bitcoinEntity.BuyBTCAction:
		return u.buyBTC(amount, price)
	case bitcoinEntity.SellBTCAction:
		return u.sellBTC(amount, price)
	default:
		return bitcoinEntity.ErrInvalidBTCAction
	}
}

func (u *User) buyBTC(amount bitcoinEntity.BTC, price bitcoinEntity.BTCPrice) error {
	if u.Balance.USD.LessThan(price.GetPrice()) {
		return ErrInsufficientFunds
	}

	updatedUSD, err := u.Balance.USD.Sub(amount.ToUSD(price))
	if err != nil {
		return err
	}

	u.Balance.USD = updatedUSD
	u.Balance.BTC = u.Balance.BTC.Add(amount)

	return nil
}

func (u *User) sellBTC(amount bitcoinEntity.BTC, price bitcoinEntity.BTCPrice) error {
	if u.Balance.BTC.LessThan(amount) {
		return ErrInsufficientFunds
	}

	updatedBTC, err := u.Balance.BTC.Sub(amount)
	if err != nil {
		return err
	}

	u.Balance.BTC = updatedBTC
	u.Balance.USD = u.Balance.USD.Add(amount.ToUSD(price))

	return nil
}
