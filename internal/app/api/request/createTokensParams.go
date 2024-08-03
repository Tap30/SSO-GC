package request

import (
	"net/url"
)

// CreateTokensParams holds the parameters needed for creating tokens in OAuth flows.
type CreateTokensParams struct {
	Code         string `form:"code"`          // Code is the authorization code received from the authorization server.
	RedirectURI  string `form:"redirect_uri"`  // RedirectURI is the URL to which the user will be sent after authorization.
	GrantType    string `form:"grant_type"`    // GrantType is the type of grant being requested (e.g., authorization_code, refresh_token).
	RefreshToken string `form:"refresh_token"` // RefreshToken is the token used to obtain additional access tokens.
	CustomClaim  string `form:"custom_claim"`  // CustomClaim is the custom claim to be added to the token.
}

// ToValues converts the CreateTokensParams to url.Values for HTTP requests.
func (p *CreateTokensParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("code", p.Code)
	values.Set("redirect_uri", p.RedirectURI)
	values.Set("grant_type", p.GrantType)
	values.Set("refresh_token", p.RefreshToken)
	values.Set("custom_claim", p.CustomClaim)
	return values
}
