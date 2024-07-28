package auth

import (
	"SSO-GC/api/request"
	"SSO-GC/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"github.com/patrickmn/go-cache"
	"net/url"
	"time"
)

var (
	c   = cache.New(5*time.Minute, 10*time.Minute)
	cfg = config.LoadConfig() // Load the configuration once
)

// CustomClaims represents your custom claims in the authentication token.
type CustomClaims struct {
	// Define the fields according to your requirements
}

// ToJson converts CustomClaims to a JSON string.
func (cc *CustomClaims) ToJson() (string, error) {
	data, err := json.Marshal(cc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// createCustomClaims generates custom claims based on provided parameters.
func createCustomClaims(params *request.CreateTokensParams) (*CustomClaims, error) {
	// Implementation specific to your business logic
	return &CustomClaims{}, nil
}

// GetOpenIDConfiguration retrieves the OpenID configuration, caching it to optimize performance.
func GetOpenIDConfiguration() (map[string]interface{}, error) {
	cacheKey := "openid-configuration"
	if cachedConfig, found := c.Get(cacheKey); found {
		return cachedConfig.(map[string]interface{}), nil
	}

	client := resty.New().RemoveProxy()
	resp, err := client.R().Get(cfg.SsoIssuer + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}

	var configMap map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &configMap); err != nil {
		return nil, err
	}

	c.Set(cacheKey, configMap, cache.DefaultExpiration)
	return configMap, nil
}

// GetTokens orchestrates the token retrieval process.
func GetTokens(params *request.CreateTokensParams) (map[string]interface{}, error) {
	customClaims, err := createCustomClaims(params)
	if err != nil {
		return nil, err
	}

	customClaimsJson, err := customClaims.ToJson()
	if err != nil {
		return nil, err
	}

	formData := params.ToValues()
	formData.Add("custom_claims", customClaimsJson)

	clientAuthHeader := createClientAuthHeader(cfg.ClientId, cfg.ClientSecret)
	return getTokensFromTokenEndpoint(clientAuthHeader, formData)
}

// getTokensFromTokenEndpoint sends a POST request to fetch tokens.
func getTokensFromTokenEndpoint(authHeader string, formData url.Values) (map[string]interface{}, error) {
	openidConfig, err := GetOpenIDConfiguration()
	if err != nil {
		return nil, err
	}

	tokenEndpoint, ok := openidConfig["token_endpoint"].(string)
	if !ok {
		return nil, errors.New("token endpoint not found in OpenID configuration")
	}

	client := resty.New().RemoveProxy()
	response, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": authHeader,
		}).
		SetFormDataFromValues(formData).
		Post(tokenEndpoint)

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		log.Error(response.Error(), " ", response.StatusCode())
		return nil, errors.New("failed to get tokens from token endpoint")
	}

	var tokens map[string]interface{}
	if err := json.Unmarshal(response.Body(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetUserInfo retrieves user information using the access token.
func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	openidConfig, err := GetOpenIDConfiguration()
	if err != nil {
		return nil, err
	}

	userInfoEndpoint, ok := openidConfig["userinfo_endpoint"].(string)
	if !ok {
		return nil, errors.New("userinfo endpoint not found in OpenID configuration")
	}

	client := resty.New().RemoveProxy()
	response, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(userInfoEndpoint)

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		log.Error(response.Error(), " ", response.StatusCode())
		return nil, errors.New("failed to get user info from userinfo endpoint")
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(response.Body(), &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// createClientAuthHeader generates the HTTP Basic Authentication header.
func createClientAuthHeader(clientId, clientSecret string) string {
	auth := clientId + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
