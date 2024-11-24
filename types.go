package extend

import (
	"fmt"
	"time"
)

type Organization struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joinedAt"`
	Explicit bool   `json:"explicit"`
}

type Asset struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

type Address struct {
	Address1 string `json:"address1"`
	City     string `json:"city"`
	Province string `json:"province"`
	Postal   string `json:"postal"`
	Country  string `json:"country"`
}

type User struct {
	ID               string       `json:"id"`
	FirstName        string       `json:"firstName"`
	LastName         string       `json:"lastName"`
	Email            string       `json:"email"`
	PhoneIsoCountry  string       `json:"phoneIsoCountry"`
	AvatarType       string       `json:"avatarType"`
	CreatedAt        string       `json:"createdAt"`
	UpdatedAt        string       `json:"updatedAt"`
	Currency         string       `json:"currency"`
	Locale           string       `json:"locale"`
	Timezone         string       `json:"timezone"`
	Verified         bool         `json:"verified"`
	Organization     Organization `json:"organization"`
	OrganizationID   string       `json:"organizationId"`
	OrganizationRole string       `json:"organizationRole"`
}

type Currency string

const (
	CurrencyUSD Currency = "USD"
)

type Time struct {
	time.Time
}

var (
	timeLayout = "2006-01-02T15:04:05.000"
)

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(timeLayout) + `+0000"`), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) < 30 {
		return fmt.Errorf("invalid time format: %s", string(data))
	}

	str := string(data[1 : len(data)-6])
	v, err := time.Parse(timeLayout, str)
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

func join[T ~string](values []T, sep string) string {
	v := ""
	for i, value := range values {
		if i > 0 {
			v += sep
		}
		v += string(value)
	}
	return v
}

type apiErrorDetail struct {
	Field        string `json:"field"`
	Error        string `json:"error"`
	InvalidValue string `json:"invalidValue"`
}

type apiErrorResponse struct {
	ErrorMessage string           `json:"error"`
	Details      []apiErrorDetail `json:"details"`
}

func (e apiErrorResponse) Error() string {
	message := "extend: " + e.ErrorMessage
	for _, detail := range e.Details {
		message = fmt.Sprintf("%s (%s: %s)", message, detail.Field, detail.Error)
	}
	return message
}
