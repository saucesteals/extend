package extend

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

func (c *Client) GetBulkVirtualCardUpload(ctx context.Context, uploadId string) (*BulkVirtualCardUpload, error) {
	var response bulkVirtualCardUploadResponse
	err := c.jsonRequest(ctx, http.MethodGet, fmt.Sprintf("/bulkvirtualcarduploads/%s", uploadId), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response.BulkVirtualCardUpload, nil
}

type BulkVirtualCardUploadStatus string

const (
	BulkVirtualCardUploadStatusInitiated BulkVirtualCardUploadStatus = "Initiated"
	BulkVirtualCardUploadStatusCompleted BulkVirtualCardUploadStatus = "Completed"
)

type BulkVirtualCardUploadTask struct {
	TaskID        string                      `json:"taskId"`
	Status        BulkVirtualCardUploadStatus `json:"status"`
	VirtualCardID string                      `json:"virtualCardId"`
}

type BulkVirtualCardUpload struct {
	ID           string                      `json:"id"`
	UserID       string                      `json:"userId"`
	CreditCardID string                      `json:"creditCardId"`
	CreatedAt    Time                        `json:"createdAt"`
	UpdatedAt    Time                        `json:"updatedAt"`
	Tasks        []BulkVirtualCardUploadTask `json:"tasks"`
}

type bulkVirtualCardUploadResponse struct {
	BulkVirtualCardUpload BulkVirtualCardUpload `json:"bulkVirtualCardUpload"`
}

func (c *Client) BulkCreateVirtualCards(ctx context.Context, cardId string, options []BulkCreateVirtualCard) (*BulkVirtualCardPushResponse, error) {
	var csv strings.Builder
	csv.WriteString(`"Card Type","en-US","Virtual Card User Email","Card Name","Credit Limit","Active Until Date (MM/DD/YYYY)","Notes"`)
	for _, option := range options {
		csv.WriteString("\n")
		csv.WriteString(
			fmt.Sprintf(
				`"%s","en-US","%s","%s",%.2f,"%s","%s"`,
				option.CardType, option.Recipient, option.DisplayName, float64(option.BalanceCents)/100, option.ValidTo.Format("01/02/2006"), option.Notes,
			),
		)
	}

	body := new(bytes.Buffer)
	form := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="virtual_cards.csv"`)
	h.Set("Content-Type", "text/csv")
	file, err := form.CreatePart(h)
	if err != nil {
		return nil, err
	}
	file.Write([]byte(csv.String()))
	form.Close()

	var response BulkVirtualCardPushResponse
	err = c.request(ctx, http.MethodPost, fmt.Sprintf("/creditcards/%s/bulkvirtualcardpush", cardId), form.FormDataContentType(), body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type BulkCreateVirtualCard struct {
	CardType VirtualCardType

	// Recipient is the email of the recipient
	Recipient    string
	DisplayName  string
	BalanceCents int

	// ValidTo is the date the card expires (date only)
	ValidTo time.Time
	Notes   string
}

type BulkVirtualCardRecord struct {
	CreditCardID   string `json:"creditCardId"`
	Recipient      string `json:"recipient"`
	Cardholder     string `json:"cardholder"`
	DisplayName    string `json:"displayName"`
	Direct         bool   `json:"direct"`
	BalanceCents   int    `json:"balanceCents"`
	Currency       string `json:"currency"`
	ValidToDate    []int  `json:"validToDate"`
	Recurs         bool   `json:"recurs"`
	HasPlasticCard bool   `json:"hasPlasticCard"`
	SingleExactPay bool   `json:"singleExactPay"`
	IsPush         bool   `json:"isPush"`
	IsRequest      bool   `json:"isRequest"`
	UntilDate      []int  `json:"untilDate"`
}

type BulkVirtualCardTask struct {
	TaskID string                      `json:"taskId"`
	Status BulkVirtualCardUploadStatus `json:"status"`
	Record BulkVirtualCardRecord       `json:"record"`
}

type BulkVirtualCardPush struct {
	BulkVirtualCardUploadID string                `json:"bulkVirtualCardUploadId"`
	Tasks                   []BulkVirtualCardTask `json:"tasks"`
}

type BulkVirtualCardPushResponse struct {
	BulkVirtualCardPush BulkVirtualCardPush `json:"bulkVirtualCardPush"`
	InvalidEmails       []string            `json:"invalidEmails"`
	CsvVirtualCardPush  BulkVirtualCardPush `json:"csvVirtualCardPush"`
}
