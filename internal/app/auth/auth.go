package auth

import (
	"SSO-GC/api/request"
	"SSO-GC/config"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"github.com/patrickmn/go-cache"
	"net/url"
	"time"
)

var (
	c = cache.New(5*time.Minute, 10*time.Minute)
)

// Assuming CustomClaims is a struct that represents your custom claims
type CustomClaims struct {
	// Define the fields according to your requirements
}

func GetOpenIDConfiguration() (map[string]interface{}, error) {
	cacheKey := "openid-configuration"
	if cachedConfig, found := c.Get(cacheKey); found {
		return cachedConfig.(map[string]interface{}), nil
	}

	cfg := config.LoadConfig()

	client := resty.New()
	client.RemoveProxy()
	resp, err := client.R().Get(cfg.SsoIssuer + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &configMap); err != nil {
		return nil, err
	}

	c.Set(cacheKey, configMap, cache.DefaultExpiration)
	return configMap, nil
}

// Converts custom claims to JSON
func (cc *CustomClaims) ToJson() (string, error) {
	data, err := json.Marshal(cc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Placeholder for creating custom claims based on CreateTokensParams
func createCustomClaims(params *request.CreateTokensParams) (*CustomClaims, error) {
	// Placeholder implementation
	return &CustomClaims{}, nil
}

// Sends a POST request to the token endpoint with the form data and custom claims
func getTokensFromTokenEndpoint(authHeader string, formData url.Values) (map[string]interface{}, error) {
	openidConfig, err := GetOpenIDConfiguration()
	if err != nil {
		return nil, err
	}

	tokenEndpoint, ok := openidConfig["token_endpoint"].(string)
	if !ok {
		return nil, errors.New("token endpoint not found in OpenID configuration")
	}

	log.Info(tokenEndpoint)

	client := resty.New()
	client.RemoveProxy()

	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", authHeader).
		SetFormDataFromValues(formData).
		Post(tokenEndpoint)

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		log.Error(response.Error())
		return nil, errors.New("failed to get tokens from token endpoint")
	}

	var tokens map[string]interface{}
	if err := json.Unmarshal(response.Body(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// Orchestrates the process of getting tokens
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
	log.Info(formData)

	// Assuming clientAuthHeader is available or retrieved from somewhere
	clientAuthHeader := "Bearer your_auth_token_here"
	tokens, err := getTokensFromTokenEndpoint(clientAuthHeader, formData)
	if err != nil {
		return nil, err
	}

	// Placeholder for updateUser logic
	// Assuming updateUser updates user information based on tokens and returns nil on success
	// err = updateUser(tokens)
	// if err != nil {
	//     return nil, err
	// }

	return tokens, nil
}
