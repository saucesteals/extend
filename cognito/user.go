package cognito

import (
	"context"
)

type userSrpAuth struct {
	AuthFlow       string            `json:"AuthFlow"`
	ClientID       string            `json:"ClientId"`
	AuthParameters map[string]string `json:"AuthParameters"`
	ClientMetadata struct{}          `json:"ClientMetadata"`
}

type userSrpResponse struct {
	ChallengeName       string            `json:"ChallengeName"`
	ChallengeParameters map[string]string `json:"ChallengeParameters"`
}

func (c *Cognito) userSrpAuth(ctx context.Context) (*userSrpResponse, error) {
	var res userSrpResponse
	err := c.request(ctx, "AWSCognitoIdentityProviderService.InitiateAuth", userSrpAuth{
		AuthFlow:       "USER_SRP_AUTH",
		ClientID:       clientId,
		AuthParameters: c.csrp.GetAuthParams(),
		ClientMetadata: struct{}{},
	}, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type userPasswordVerifier struct {
	ChallengeName      string            `json:"ChallengeName"`
	ClientID           string            `json:"ClientId"`
	ChallengeResponses map[string]string `json:"ChallengeResponses"`
	ClientMetadata     struct{}          `json:"ClientMetadata"`
}

type userPasswordVerifierResponse struct {
	ChallengeName       string   `json:"ChallengeName"`
	ChallengeParameters struct{} `json:"ChallengeParameters"`
	Session             string   `json:"Session"`
}

func (c *Cognito) userPasswordVerifier(ctx context.Context, challengeParameters map[string]string) (*userPasswordVerifierResponse, error) {
	challengeResponses, err := c.csrp.PasswordVerifierChallenge(challengeParameters)
	if err != nil {
		return nil, err
	}

	var res userPasswordVerifierResponse
	err = c.request(ctx, "AWSCognitoIdentityProviderService.RespondToAuthChallenge", userPasswordVerifier{
		ChallengeName:      "PASSWORD_VERIFIER",
		ClientID:           clientId,
		ChallengeResponses: challengeResponses,
		ClientMetadata:     struct{}{},
	}, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
