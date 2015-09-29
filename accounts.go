package saasquatch

import (
	"fmt"
)

type AccountsServices struct {
	client *Client
}

type Account struct {
	Id           string        `json:"id"`
	Currency     string        `json:"currency"`
	Subscription *Subscription `json:"subscription,omitempty"`
	Referral     *Referral     `json:"referral,omitempty"`
}

type Subscription struct {
	Status               string  `json:"status"`
	BillingIntervalType  string  `json:"billingIntervalType,omitempty"`
	BillingIntervalValue int     `json:"billingIntervalValue,omitempty"`
	Value                float32 `json:"value,omitempty"`
}

type Referral struct {
	Code string `json:"code"`
}

func (s *AccountsServices) CreateOrUpdateAccount(payload Account) (*Account, error) {
	req, err := s.client.NewRequest("POST", "accountsync", payload)
	if err != nil {
		return nil, err
	}

	resp := new(Account)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AccountsServices) LookupAccount(accountId string) (*Account, error) {
	u := fmt.Sprintf("account/%s", accountId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp := new(Account)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
