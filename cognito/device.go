package cognito

import (
	"context"
)

type deviceSrpAuth struct {
	ChallengeName      string            `json:"ChallengeName"`
	ClientID           string            `json:"ClientId"`
	ChallengeResponses map[string]string `json:"ChallengeResponses"`
	Session            string            `json:"Session"`
}

type deviceSrpResponse struct {
	ChallengeName       string            `json:"ChallengeName"`
	ChallengeParameters map[string]string `json:"ChallengeParameters"`
}

func (c *Cognito) deviceSrpAuth(ctx context.Context, session string) (*deviceSrpResponse, error) {
	var res deviceSrpResponse
	err := c.request(ctx, "AWSCognitoIdentityProviderService.RespondToAuthChallenge", deviceSrpAuth{
		ChallengeName:      "DEVICE_SRP_AUTH",
		ClientID:           clientId,
		ChallengeResponses: c.csrp.GetDeviceAuthParams(),
		Session:            session,
	}, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type devicePasswordVerifier struct {
	ChallengeName      string            `json:"ChallengeName"`
	ClientID           string            `json:"ClientId"`
	ChallengeResponses map[string]string `json:"ChallengeResponses"`
}

type initialAuthenticationResult struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
	ExpiresIn    int    `json:"ExpiresIn"`
}

type devicePasswordVerifierResponse struct {
	AuthenticationResult initialAuthenticationResult `json:"AuthenticationResult"`
	ChallengeParameters  struct{}                    `json:"ChallengeParameters"`
}

func (c *Cognito) devicePasswordVerifier(ctx context.Context, userId string, challengeParameters map[string]string) (*devicePasswordVerifierResponse, error) {
	challengeResponses, err := c.csrp.DevicePasswordVerifierChallenge(userId, challengeParameters)
	if err != nil {
		return nil, err
	}

	var res devicePasswordVerifierResponse
	err = c.request(ctx, "AWSCognitoIdentityProviderService.RespondToAuthChallenge", devicePasswordVerifier{
		ChallengeName:      "DEVICE_PASSWORD_VERIFIER",
		ClientID:           clientId,
		ChallengeResponses: challengeResponses,
	}, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
