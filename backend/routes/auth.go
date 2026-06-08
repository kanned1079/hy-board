package routes

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
)

// Subscribe handles subscription links: GET /api/v1/client/subscribe
func Subscribe(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Token is required")
		return
	}

	var user models.User
	// For subscription, token can be matched via v2ray_uuid (since it is a unique user key)
	if err := database.DB.Where("v2ray_uuid = ? AND status = ?", token, 1).First(&user).Error; err != nil {
		c.String(http.StatusUnauthorized, "Invalid token")
		return
	}

	var nodes []models.Node
	database.DB.Where("show = ?", true).Find(&nodes)

	var subContent string
	for _, node := range nodes {
		// Output simple raw proxy configurations (like vmess://, vless://, trojan://)
		switch node.Type {
		case "V2ray":
			// Simple VMess json encoded in base64
			vmessJSON := fmt.Sprintf(`{"v":"2","ps":"%s","add":"%s","port":%d,"id":"%s","aid":"0","net":"tcp","type":"none","host":"","path":"","tls":""}`,
				node.Name, node.Address, node.Port, user.V2rayUUID)
			encoded := base64.StdEncoding.EncodeToString([]byte(vmessJSON))
			subContent += "vmess://" + encoded + "\n"
		case "Vless":
			subContent += fmt.Sprintf("vless://%s@%s:%d?type=tcp&security=none#%s\n",
				user.V2rayUUID, node.Address, node.Port, node.Name)
		case "Trojan":
			subContent += fmt.Sprintf("trojan://%s@%s:%d#%s\n",
				user.TrojanPassword, node.Address, node.Port, node.Name)
		}
	}

	// Base64 encode the output configuration for standard clients (Shadowrocket / Surge / V2rayN)
	encodedSub := base64.StdEncoding.EncodeToString([]byte(subContent))
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, encodedSub)
}
