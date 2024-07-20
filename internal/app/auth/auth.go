package auth

import (
	"SSO-GC/config"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	c = cache.New(5*time.Minute, 10*time.Minute)
)

// GetOpenIDConfiguration fetches the OpenID configuration from the server or cache
func GetOpenIDConfiguration() (map[string]interface{}, error) {
	cacheKey := "openid-configuration"
	if cachedConfig, found := c.Get(cacheKey); found {
		return cachedConfig.(map[string]interface{}), nil
	}

	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	client := resty.New()
	resp, err := client.R().Get(appConfig.SsoIssuer + "/.well-known/openid-configuration")
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
