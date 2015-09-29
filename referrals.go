package saasquatch

import (
	"fmt"
)

type ReferralsService struct {
	client *Client
}

type ReferralCodes struct {
	Code         string `json:"code"`
	DateCreated  int    `json:"dateCreated"`
	ReferrerName string `json:"referrerName"`
	Reward       Reward `json:"reward"`
}

type Reward struct {
	Type                  string `json:"type"`
	Unit                  string `json:"unit"`
	Credit                int    `json:"credit"`
	DiscountPercent       int    `json:"discountPercent"`
	MonthsDiscountIsValid int    `json:"monthsDiscountIsValid"`

	//Only works with types FEATURE
	FeatureType string `json:"featureType,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
}

func (s *ReferralsService) LookupReferralCode(cd string) (*ReferralCodes, error) {
	u := fmt.Sprintf("code/%v", cd)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rc := new(ReferralCodes)
	if err := s.client.Do(req, rc); err != nil {
		return nil, err
	}

	return rc, nil
}
