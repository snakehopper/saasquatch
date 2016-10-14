package saasquatch

import (
	"github.com/google/go-querystring/query"
	"net/url"
)

type ReferralService struct {
	client *Client
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
	Id                       string          `json:"id"`
	ReferredUser             *ReferredUser   `json:"referredUser"`
	ReferrerUser             *ReferredUser   `json:"referrerUser"`
	ReferredReward           *ReferredReward `json:"referredReward"`
	ReferrerReward           *ReferredReward `json:"referrerReward"`
	ModerationStatus         string          `json:"moderationStatus"`
	ReferredModerationStatus string          `json:"referredModerationStatus"`
	ReferrerModerationStatus string          `json:"referrerModerationStatus"`
	FraudSignals             *FraudSignals   `json:"fraudSignals,omitempty"`
	DateReferralStarted      int             `json:"dateReferralStarted"`
	DateReferralPaid         int             `json:"dateReferralPaid"`
	DateReferralEnded        int             `json:"dateReferralEnded"`
	DateModerated            int             `json:"dateModerated"`
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

func (s *ReferralService) ListReferrals(opt ListReferralsOptions) (*ListReferralsResult, error) {
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
