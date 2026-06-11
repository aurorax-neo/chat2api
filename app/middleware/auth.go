package middleware

import (
	"chat2api/app/common"
	"chat2api/app/conf"
	"strings"

	"github.com/gin-gonic/gin"
)

func V1Auth(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	localToken := strings.TrimSpace(strings.TrimPrefix(authToken, "Bearer "))
	appConf := conf.GetApp()
	hasAccessTokenPrefix := appConf.HasAccessTokenPrefix()
	if authToken == "" {
		common.ErrorResponse(c, 401, missingAPIKeyMessage(hasAccessTokenPrefix), nil)
		return
	}
	if _, ok := appConf.DirectAccessToken(localToken); ok {
		c.Next()
		return
	}
	if appConf.DirectAccessTokenPrefixMatched(localToken) {
		common.ErrorResponse(c, 401, "Invalid access token for access_token_prefix mode. Use Authorization: Bearer <configured-prefix><real_access_token>.", nil)
		return
	}
	if strings.HasPrefix(authToken, "Bearer eyJhbGciOiJSUzI1NiI") {
		c.Next()
		return
	}
	if len(appConf.Auth.AccessTokens) == 0 {
		common.ErrorResponse(c, 401, noLocalAPIKeysMessage(hasAccessTokenPrefix), nil)
		return
	}
	if !common.IsStrInArray(localToken, appConf.Auth.AccessTokens) {
		common.ErrorResponse(c, 401, incorrectAPIKeyMessage(hasAccessTokenPrefix), nil)
		return
	}
	c.Next()
}

func AdminAuth(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	localToken := strings.TrimSpace(strings.TrimPrefix(authToken, "Bearer "))
	if authToken == "" || localToken == "" {
		common.ErrorResponse(c, 401, "Admin API requires a configured local API key.", nil)
		return
	}
	appConf := conf.GetApp()
	if len(appConf.Auth.AccessTokens) == 0 || !common.IsStrInArray(localToken, appConf.Auth.AccessTokens) {
		common.ErrorResponse(c, 401, "Incorrect admin API key.", nil)
		return
	}
	c.Next()
}

func missingAPIKeyMessage(hasAccessTokenPrefix bool) string {
	if hasAccessTokenPrefix {
		return "You didn't provide an API key. Use Authorization: Bearer <local-api-key>, or Authorization: Bearer <configured-prefix><real_access_token> for access_token_prefix mode."
	}
	return "You didn't provide an API key. You need to provide your API key in an Authorization header using Bearer auth (i.e. Authorization: Bearer YOUR_KEY)."
}

func noLocalAPIKeysMessage(hasAccessTokenPrefix bool) string {
	if hasAccessTokenPrefix {
		return "No local API keys are configured. Use Authorization: Bearer <configured-prefix><real_access_token> for access_token_prefix mode."
	}
	return "No local API keys are configured"
}

func incorrectAPIKeyMessage(hasAccessTokenPrefix bool) string {
	if hasAccessTokenPrefix {
		return "Incorrect API key provided. Use a configured local API key, or Authorization: Bearer <configured-prefix><real_access_token> for access_token_prefix mode."
	}
	return "Incorrect API key provided."
}
