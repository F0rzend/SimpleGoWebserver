package entity

import (
	"sort"
)

var (
	DepositUSDAction  = USDAction{"deposit"}
	WithdrawUSDAction = USDAction{"withdraw"}

	BuyBTCAction  = BTCAction{"buy"}
	SellBTCAction = BTCAction{"sell"}
)

var (
	usdActions = map[string]USDAction{
		"deposit":  DepositUSDAction,
		"withdraw": WithdrawUSDAction,
	}
	btcActions = map[string]BTCAction{
		"buy":  BuyBTCAction,
		"sell": SellBTCAction,
	}
)

type (
	USDAction action
	BTCAction action
)

type action struct {
	a string
}

func (a action) String() string {
	return a.a
}

func NewUSDAction(action string) (USDAction, error) {
	usdAction, ok := usdActions[action]
	if !ok {
		return USDAction{}, ErrInvalidUSDAction
	}

	return usdAction, nil
}

func GetUSDActions() (actions []string) {
	for action := range usdActions {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	return actions
}

func NewBTCAction(action string) (BTCAction, error) {
	btcAction, ok := btcActions[action]
	if !ok {
		return BTCAction{}, ErrInvalidBTCAction
	}
	return btcAction, nil
}

func GetBTCActions() (actions []string) {
	for action := range btcActions {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	return actions
}
