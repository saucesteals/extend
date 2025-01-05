package extend

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CreateVirtualCardOptions struct {
	// CreditCardID is the ID of the credit card to use for the virtual card
	CreditCardID string   `json:"creditCardId"`
	DisplayName  string   `json:"displayName"`
	BalanceCents int      `json:"balanceCents"`
	Currency     Currency `json:"currency"`
	Notes        string   `json:"notes"`
	// ValidTo is the date the card expires (date only)
	ValidTo time.Time `json:"-"`
	// Recipient is the email of the recipient
	Recipient string `json:"recipient"`
}

type createVirtualCardOptions struct {
	CreateVirtualCardOptions
	ValidTo string `json:"validTo"`
}

func (a *Client) CreateVirtualCard(ctx context.Context, options CreateVirtualCardOptions) (*VirtualCard, error) {
	payload := createVirtualCardOptions{
		CreateVirtualCardOptions: options,
		ValidTo:                  options.ValidTo.Format("2006-01-02"),
	}
	var response VirtualCardResponse
	err := a.jsonRequest(ctx, http.MethodPost, "/virtualcards", payload, &response)
	if err != nil {
		return nil, err
	}

	return &response.VirtualCard, nil
}

type UpdateVirtualCardOptions struct {
	CreditCardID string `json:"creditCardId"`
	DisplayName  string `json:"displayName"`
	BalanceCents int    `json:"balanceCents"`
	Recurs       bool   `json:"recurs"`

	// ValidTo is the date the card expires (date only)
	ValidTo time.Time `json:"-"`

	Currency           Currency `json:"currency"`
	ReceiptRulesExempt bool     `json:"receiptRulesExempt"`
}

type updateVirtualCardOptions struct {
	UpdateVirtualCardOptions
	ValidTo string `json:"validTo"`
}

func (a *Client) UpdateVirtualCard(ctx context.Context, id string, options UpdateVirtualCardOptions) (*VirtualCard, error) {
	payload := updateVirtualCardOptions{
		UpdateVirtualCardOptions: options,
		ValidTo:                  options.ValidTo.Format("2006-01-02"),
	}
	var response VirtualCardResponse
	err := a.jsonRequest(ctx, http.MethodPut, fmt.Sprintf("/virtualcards/%s", id), payload, &response)
	if err != nil {
		return nil, err
	}

	return &response.VirtualCard, nil
}

func (a *Client) GetVirtualCard(ctx context.Context, id string) (*VirtualCard, error) {
	var response VirtualCardResponse
	err := a.jsonRequest(ctx, http.MethodGet, fmt.Sprintf("/virtualcards/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.VirtualCard, nil
}

func (a *Client) CancelVirtualCard(ctx context.Context, id string) (*VirtualCard, error) {
	var response VirtualCardResponse
	err := a.jsonRequest(ctx, http.MethodPut, fmt.Sprintf("/virtualcards/%s/cancel", id), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.VirtualCard, nil
}

func (a *Client) CloseVirtualCard(ctx context.Context, id string) (*VirtualCard, error) {
	var response VirtualCardResponse
	err := a.jsonRequest(ctx, http.MethodPut, fmt.Sprintf("/virtualcards/%s/close", id), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.VirtualCard, nil
}

type VirtualCardResponse struct {
	VirtualCard VirtualCard `json:"virtualCard"`
}

type ListVirtualCardsOptions struct {
	PaginationOptions
	CardholderOrViewer string
	Issued             bool
	Statuses           []VirtualCardStatus
}

type ListVirtualCardsResponse struct {
	PaginationResponse
	VirtualCards []VirtualCard `json:"virtualCards"`
}

func (r ListVirtualCardsResponse) Items() []VirtualCard {
	return r.VirtualCards
}

func (c *Client) ListVirtualCards(options *ListVirtualCardsOptions) *Paginator[VirtualCard, ListVirtualCardsResponse] {
	return newPaginator[VirtualCard, ListVirtualCardsResponse](c, options.PaginationOptions, "/virtualcards", url.Values{
		"cardholderOrViewer": {options.CardholderOrViewer},
		"issued":             {strconv.FormatBool(options.Issued)},
		"statuses":           {join(options.Statuses, ",")},
	})
}

type VirtualCardType string

const (
	VirtualCardTypeStandard VirtualCardType = "STANDARD"
)

type VirtualCardStatus string

const (
	VirtualCardStatusActive    VirtualCardStatus = "ACTIVE"
	VirtualCardStatusCancelled VirtualCardStatus = "CANCELLED"
	VirtualCardStatusClosed    VirtualCardStatus = "CLOSED"
)

type VirtualCardImage struct {
	ID                  string `json:"id"`
	ContentType         string `json:"contentType"`
	Urls                Asset  `json:"urls"`
	TextColorRGBA       string `json:"textColorRGBA"`
	HasTextShadow       bool   `json:"hasTextShadow"`
	ShadowTextColorRGBA string `json:"shadowTextColorRGBA"`
}

type VirtualCardFeatures struct {
	Recurrence       bool `json:"recurrence"`
	MccControl       bool `json:"mccControl"`
	QboReportEnabled bool `json:"qboReportEnabled"`
}

type VirtualCardIssuer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type VirtualCard struct {
	ID     string            `json:"id"`
	Status VirtualCardStatus `json:"status"`

	RecipientID string `json:"recipientId"`
	Recipient   User   `json:"recipient"`

	CardholderID string `json:"cardholderId"`
	Cardholder   User   `json:"cardholder"`

	Vcn          *string `json:"vcn,omitempty"`
	SecurityCode *string `json:"securityCode,omitempty"`

	LastUpdatedBy *User `json:"lastUpdatedBy,omitempty"`

	CardImage          VirtualCardImage `json:"cardImage"`
	CardType           string           `json:"cardType"`
	DisplayName        string           `json:"displayName"`
	Currency           string           `json:"currency"`
	LimitCents         int              `json:"limitCents"`
	BalanceCents       int              `json:"balanceCents"`
	SpentCents         int              `json:"spentCents"`
	LifetimeSpentCents int              `json:"lifetimeSpentCents"`
	Last4              string           `json:"last4"`
	NumberFormat       string           `json:"numberFormat"`

	InactiveSince *Time `json:"inactiveSince"`
	Expires       *Time `json:"expires"`
	ValidFrom     *Time `json:"validFrom"`
	ValidTo       *Time `json:"validTo"`
	CreatedAt     *Time `json:"createdAt"`
	UpdatedAt     *Time `json:"updatedAt"`
	ActiveUntil   *Time `json:"activeUntil"`

	Timezone     string `json:"timezone"`
	CreditCardID string `json:"creditCardId"`
	Recurs       bool   `json:"recurs"`

	Address        Address             `json:"address"`
	Features       VirtualCardFeatures `json:"features"`
	HasPlasticCard bool                `json:"hasPlasticCard"`

	Network               string            `json:"network"`
	CompanyName           string            `json:"companyName"`
	CreditCardDisplayName VirtualCardIssuer `json:"issuer"`
	ReceiptRulesExempt    bool              `json:"receiptRulesExempt"`
	IsBillPay             bool              `json:"isBillPay"`
}
