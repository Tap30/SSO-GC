package request

import (
	"net/url"
)

type CreateTokensParams struct {
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	GrantType    string `form:"grant_type"`
	RefreshToken string `form:"refresh_token"`
	CodeVerifier string `form:"code_verifier"`
}

func (p *CreateTokensParams) ToValues() url.Values {
	values := url.Values{}
	values.Set("code", p.Code)
	values.Set("redirect_uri", p.RedirectURI)
	values.Set("grant_type", p.GrantType)
	values.Set("code_verifier", p.CodeVerifier)
	values.Set("refresh_token", p.RefreshToken)
	return values
}
