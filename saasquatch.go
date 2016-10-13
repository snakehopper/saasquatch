package saasquatch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"

	"github.com/google/go-querystring/query"
)

const (
	HOST = "http://app.referralsaasquatch.com"
)

type Tenant struct {
	Alias      string
	ApiKey     string
	SecureMode bool
}

func NewTenant(alias, apiKey string) *Tenant {
	tn := &Tenant{alias, apiKey, true}
	return tn
}

func (tn *Tenant) DisableSecureMode() *Tenant {
	tn.SecureMode = false
	return tn
}

func (tn Tenant) NewMobileWidget(uid, email, firstName string) *MobileWidget {
	return NewMobileWidget(&tn, uid, email, firstName)
}

type MobileWidget struct {
	tenant *Tenant

	userId            string
	accountId         string
	paymentProviderId string

	//optional parameters
	firstName string
	lastName  string
	email     string
	locale    string
}

func NewMobileWidget(tn *Tenant, uid, email, firstName string) *MobileWidget {
	return &MobileWidget{tenant: tn,
		userId:            uid,
		accountId:         uid,
		email:             email,
		firstName:         firstName,
		paymentProviderId: "NULL"}
}

func (mw *MobileWidget) WithAccount(aid string) *MobileWidget {
	mw.accountId = aid
	return mw
}

func (mw *MobileWidget) WithPayment(pid string) *MobileWidget {
	mw.paymentProviderId = pid
	return mw
}

func (mw *MobileWidget) WithFirstName(name string) *MobileWidget {
	mw.firstName = name
	return mw
}

func (mw *MobileWidget) WithLastName(name string) *MobileWidget {
	mw.lastName = name
	return mw
}

func (mw *MobileWidget) WithEmail(email string) *MobileWidget {
	mw.email = email
	return mw
}

func (mw *MobileWidget) WithLocale(locale string) *MobileWidget {
	mw.locale = locale
	return mw
}

func (mw MobileWidget) checksum(v url.Values) string {
	keys := make([]string, 0)
	for k := range v {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var plain string
	for _, k := range keys {
		plain += v.Get(k)
	}

	h := hmac.New(sha256.New, []byte(mw.tenant.ApiKey))
	h.Write([]byte(plain))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (mw MobileWidget) BuildUrl() (string, error) {
	u, err := url.Parse(HOST)
	if err != nil {
		return "", err
	}

	u.Path = fmt.Sprintf("/a/%s/widgets/mobilewidget", mw.tenant.Alias)

	q := u.Query()
	q.Set("userId", mw.userId)
	q.Set("tenantAlias", mw.tenant.Alias)
	q.Set("accountId", mw.accountId)

	//Payment Provider Id SHOULD omitted when checksum if it is NULL
	if mw.paymentProviderId != "NULL" {
		q.Set("paymentProviderId", mw.paymentProviderId)
	}

	//Optional parameters
	q.Set("firstName", mw.firstName)
	q.Set("lastName", mw.lastName)
	q.Set("email", mw.email)

	if mw.tenant.SecureMode {
		if sum := mw.checksum(q); sum == "" {
			return "", fmt.Errorf("checksum error")
		} else {
			q.Set("checksum", sum)
		}
	}

	q.Set("paymentProviderId", mw.paymentProviderId)
	if mw.locale != "" {
		q.Set("locale", mw.locale)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// A Client manages communication with the SaaSquatch API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	apiKey string

	BaseURL *url.URL

	Referral       *ReferralService
	ReferralCode   *ReferralCodeService
	Accounts       *AccountsServices
	Users          *UsersServices
	RewardBalances *RewardBalancesServices
}

func NewClient(httpClient *http.Client, tenantAlias, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u := fmt.Sprintf("https://app.referralsaasquatch.com/api/v1/%s/", tenantAlias)
	baseURL, _ := url.Parse(u)

	c := &Client{client: httpClient, apiKey: apiKey, BaseURL: baseURL}
	c.Referral = &ReferralService{client: c}
	c.ReferralCode = &ReferralCodeService{client: c}
	c.Accounts = &AccountsServices{client: c}
	c.Users = &UsersServices{client: c}
	c.RewardBalances = &RewardBalancesServices{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("ApiKey", c.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return err
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
