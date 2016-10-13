package saasquatch

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
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

type ListReferralsOptions struct {
	ReferringAccountId       string `url:"referringAccountId,omitempty"`
	ReferringUserId          string `url:"referringUserId,omitempty"`
	DateReferralPaid         string `url:"dateReferralPaid,omitempty"`
	DateReferralEnded        string `url:"dateReferralEnded,omitempty"`
	ReferredModerationStatus string `url:"referredModerationStatus,omitempty"`
	ReferrerModerationStatus string `url:"referrerModerationStatus,omitempty"`
	Limit                    int    `url:"limit,omitempty"`
	Offset                   int    `url:"offset,omitempty"`
}

type ListReferralsResult struct {
	Count      int              `json:"count"`
	TotalCount int              `json:"totalCount"`
	Referrals  []ReferralObject `json:"referrals"`
}

type ReferralObject struct {
	Id                       string         `json:"id"`
	ReferredUser             ReferredUser   `json:"referredUser"`
	ReferrerUser             ReferredUser   `json:"referrerUser"`
	ReferredReward           ReferredReward `json:"referredReward"`
	ReferrerReward           ReferredReward `json:"referrerReward"`
	ModerationStatus         string         `json:"moderationStatus"`
	ReferredModerationStatus string         `json:"referredModerationStatus"`
	ReferrerModerationStatus string         `json:"referrerModerationStatus"`
	FraudSignals             FraudSignals   `json:"fraudSignals"`
	DateReferralStarted      int            `json:"dateReferralStarted"`
	DateReferralPaid         int            `json:"dateReferralPaid"`
	DateReferralEnded        int            `json:"dateReferralEnded"`
	DateModerated            int            `json:"dateModerated"`
}

type ReferredUser struct {
	Id             string `json:"id"`
	AccountId      string `json:"accountId"`
	ReferralCode   string `json:"referralCode"`
	Email          string `json:"email"`
	ImageUrl       string `json:"imageUrl"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	FirstSeenIP    string `json:"firstSeenIP"`
	LastSeenIP     string `json:"lastSeenIP"`
	DateCreated    int    `json:"dateCreated"`
	EmailHash      string `json:"emailHash"`
	ReferralSource string `json:"referralSource"`
	Locale         string `json:"locale"`
	//ShareLinks     string `json:"-"`
}

type ReferredReward struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	DateGiven       int    `json:"dateGiven"`
	DateExpires     int    `json:"dateExpires"`
	DateCancelled   int    `json:"dateCancelled"`
	Cancellable     bool   `json:"cancellable"`
	RewardSource    string `json:"rewardSource"`
	Unit            string `json:"unit"`
	DiscountPercent int    `json:"discountPercent"`
	FeatureType     string `json:"featureType"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Quantity        int    `json:"quantity"`
	AssignedCredit  int    `json:"assignedCredit"`
	RedeemedCredit  int    `json:"redeemedCredit"`
	Currency        string `json:"currency"`
}

type FraudSignals struct {
	Name  FraudSignal `json:"name"`
	Ip    FraudSignal `json:"ip"`
	Email FraudSignal `json:"email"`
	Rate  FraudSignal `json:"rate"`
}

type FraudSignal struct {
	Message string `json:"message"`
	Score   int    `json:"score"`
}

func (s *ReferralsService) ListReferrals(opt ListReferralsOptions) (*ListReferralsResult, error) {
	ul, err := url.Parse("referrals")
	if err != nil {
		return nil, err
	}

	v, _ := query.Values(opt)
	ul.RawQuery = v.Encode()
	req, err := s.client.NewRequest("GET", ul.String(), nil)
	if err != nil {
		return nil, err
	}

	resp := new(ListReferralsResult)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
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
