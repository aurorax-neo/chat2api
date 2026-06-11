package service

import (
	"chat2api/app/conf"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminConfigPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", adminConfigHTML)
}

func GetAdminConfig(c *gin.Context) {
	cfg, err := conf.LoadAdminConfig()
	if err != nil {
		adminConfigError(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func SaveAdminConfig(c *gin.Context) {
	var req conf.AdminConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, err := conf.SaveAdminConfig(req)
	if err != nil {
		adminConfigError(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func adminConfigError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, conf.ErrAdminConfigUnavailable) {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
