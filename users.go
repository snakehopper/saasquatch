package saasquatch

import (
	"fmt"
)

type UsersServices struct {
	client *Client
}

type User struct {
	Id           string `json:"id"`
	AccountId    string `json:"accountId"`
	ReferralCode string `json:"referralCode,omitempty"`
	Email        string `json:"email"`
	ImageUrl     string `json:"imageUrl,omitempty"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

func (s *UsersServices) CreateOrUpdateUser(payload User) (*User, error) {
	req, err := s.client.NewRequest("POST", "user", payload)
	if err != nil {
		return nil, err
	}

	resp := new(User)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UsersServices) LookupUser(accountId, userId string) (*User, error) {
	u := fmt.Sprintf("account/%s/user/%s", accountId, userId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp := new(User)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
