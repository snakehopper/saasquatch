package saasquatch

import (
	"fmt"
)

type RewardBalancesServices struct {
	client *Client
}

type ListRewardOptions struct {
	AccountId         string `url:"accountId"`
	UserId            string `url:"userId,omitempty"`
	RewardTypeFilter  string `url:"rewardTypeFilter,omitempty"`
	FeatureTypeFilter string `url:"featureTypeFilter,omitempty"`
}

type RewardBalances struct {
	Type                    string `json:"type"`
	Unit                    string `json:"unit,omitempty"`
	Count                   int    `json:"count,omitempty"`
	FeatureType             string `json:"featureType,omitempty"`
	TotalAssignedCredit     int    `json:"totalAssignedCredit,omitempty"`
	TotalRedeemedCredit     int    `json:"totalRedeemedCredit,omitempty"`
	TotalDiscountPercent    int    `json:"totalDiscountPercent,omitempty"`
	ReferredDiscountPercent int    `json:"referredDiscountPercent,omitempty"`
	ReferrerDiscountPercent int    `json:"referrerDiscountPercent,omitempty"`
}

func (s *RewardBalancesServices) ListRewardBalances(opt ListRewardOptions) ([]RewardBalances, error) {
	u := fmt.Sprintf("reward/balance")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp := new([]RewardBalances)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return *resp, nil
}

type DebitRewardOptions struct {
	AccountId string `json:"accountId"`
	Unit      string `json:"unit"`
	Amount    int    `json:"amount"`
}

type BalanceDebitted struct {
	CreditRedeemed  int    `json:"creditRedeemed"`
	CreditAvailable int    `json:"creditAvailable"`
	Unit            string `json:"unit"`
}

func (s *RewardBalancesServices) DebitRewardBalance(opt DebitRewardOptions) (*BalanceDebitted, error) {
	req, err := s.client.NewRequest("POST", "credit/bulkredeem", opt)
	if err != nil {
		return nil, err
	}

	resp := new(BalanceDebitted)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
