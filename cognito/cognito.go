package cognito

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/saucesteals/extend"
)

var (
	clientId     = "79k2g0t0ujq2tfchb23d5j6htk"
	userPoolName = "pN4CuZHEc"
)

type Cognito struct {
	csrp *srpAuthentication

	accessToken  string
	refreshToken string
	expiry       time.Time

	http *http.Client

	mu sync.Mutex
}

var (
	_ extend.Authenticator = (*Cognito)(nil)
)

func NewCognito(auth AuthParams) *Cognito {
	csrp := newSRP(auth)
	return &Cognito{
		csrp: csrp,
		http: http.DefaultClient,
	}
}

func (c *Cognito) Expiry() time.Time {
	return c.expiry
}

func (c *Cognito) SetHTTPClient(http *http.Client) {
	c.http = http
}

func (c *Cognito) request(ctx context.Context, target string, body any, response any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://cognito-idp.us-east-1.amazonaws.com/", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept-language":    {"en-US,en;q=0.7"},
		"cache-control":      {"no-cache"},
		"dnt":                {"1"},
		"origin":             {"https://app.paywithextend.com"},
		"pragma":             {"no-cache"},
		"priority":           {"u=1, i"},
		"referer":            {"https://app.paywithextend.com/"},
		"sec-ch-ua":          {`"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"sec-gpc":            {"1"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"},
		"X-Amz-User-Agent":   {"aws-amplify/5.0.4 auth framework/1"},
		"X-Amz-Target":       {target},
		"Content-Type":       {"application/x-amz-json-1.1"},
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("status code %s: %s", res.Status, string(body))
	}

	return json.NewDecoder(res.Body).Decode(response)
}

func (c *Cognito) Login(ctx context.Context) (string, error) {
	userChallenge, err := c.userSrpAuth(ctx)
	if err != nil {
		return "", fmt.Errorf("user auth: %w", err)
	}

	if userChallenge.ChallengeName != "PASSWORD_VERIFIER" {
		return "", fmt.Errorf("unexpected user challenge name: %s", userChallenge.ChallengeName)
	}

	deviceAuth, err := c.userPasswordVerifier(ctx, userChallenge.ChallengeParameters)
	if err != nil {
		return "", fmt.Errorf("user password verifier: %w", err)
	}

	if deviceAuth.ChallengeName != "DEVICE_SRP_AUTH" {
		return "", fmt.Errorf("unexpected device challenge name: %s", deviceAuth.ChallengeName)
	}

	deviceChallenge, err := c.deviceSrpAuth(ctx, deviceAuth.Session)
	if err != nil {
		return "", fmt.Errorf("device srp auth: %w", err)
	}

	tokens, err := c.devicePasswordVerifier(ctx, userChallenge.ChallengeParameters["USER_ID_FOR_SRP"], deviceChallenge.ChallengeParameters)
	if err != nil {
		return "", fmt.Errorf("device password verifier: %w", err)
	}

	c.accessToken = tokens.AuthenticationResult.AccessToken
	c.refreshToken = tokens.AuthenticationResult.RefreshToken
	c.expiry = time.Now().Add(time.Duration(tokens.AuthenticationResult.ExpiresIn) * time.Second)

	return c.accessToken, nil
}

type refreshAuthParameters struct {
	RefreshToken string `json:"REFRESH_TOKEN"`
	DeviceKey    string `json:"DEVICE_KEY"`
}

type refreshPayload struct {
	ClientId       string                `json:"ClientId"`
	AuthFlow       string                `json:"AuthFlow"`
	AuthParameters refreshAuthParameters `json:"AuthParameters"`
	ClientMetadata struct{}              `json:"ClientMetadata"`
}

type refreshAuthenticationResult struct {
	IdToken     string `json:"IdToken"`
	AccessToken string `json:"AccessToken"`
	ExpiresIn   int    `json:"ExpiresIn"`
}

type refreshResponse struct {
	AuthenticationResult refreshAuthenticationResult `json:"AuthenticationResult"`
	ChallengeParameters  struct{}                    `json:"ChallengeParameters"`
}

func (c *Cognito) Refresh(ctx context.Context) (string, error) {
	var res refreshResponse
	err := c.request(ctx, "AWSCognitoIdentityProviderService.InitiateAuth", refreshPayload{
		ClientId: clientId,
		AuthFlow: "REFRESH_TOKEN_AUTH",
		AuthParameters: refreshAuthParameters{
			RefreshToken: c.refreshToken,
			DeviceKey:    c.csrp.auth.DeviceKey,
		},
		ClientMetadata: struct{}{},
	}, &res)
	if err != nil {
		return "", err
	}

	c.accessToken = res.AuthenticationResult.AccessToken
	c.expiry = time.Now().Add(time.Duration(res.AuthenticationResult.ExpiresIn) * time.Second)

	return res.AuthenticationResult.AccessToken, nil
}

func (c *Cognito) GetAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.refreshToken == "" {
		return c.Login(ctx)
	}

	if expiresSoon(c.expiry) {
		return c.Refresh(ctx)
	}

	return c.accessToken, nil
}

func expiresSoon(expiry time.Time) bool {
	return time.Now().Add(time.Minute * 5).After(expiry)
}
