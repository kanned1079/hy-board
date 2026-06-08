package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
)

// Helper to check if a user group is in the allowed node groups list
func isUserInNodeGroups(userGroupId uint, nodeGroupIdsStr string) bool {
	if nodeGroupIdsStr == "" {
		return false
	}
	parts := strings.Split(nodeGroupIdsStr, ",")
	target := strconv.Itoa(int(userGroupId))
	for _, part := range parts {
		if strings.TrimSpace(part) == target {
			return true
		}
	}
	return false
}

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

	var allNodes []models.Node
	database.DB.Where("show = ?", true).Find(&allNodes)

	var nodes []models.Node
	for _, node := range allNodes {
		if isUserInNodeGroups(user.GroupID, node.GroupIDs) {
			nodes = append(nodes, node)
		}
	}

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
		case "Shadowsocks":
			// ss://method:password@hostname:port#tag
			method := "aes-256-gcm"
			if node.Settings != "" {
				var settingsMap map[string]interface{}
				if err := json.Unmarshal([]byte(node.Settings), &settingsMap); err == nil {
					if m, ok := settingsMap["method"].(string); ok && m != "" {
						method = m
					}
				}
			}
			ssUserInfo := fmt.Sprintf("%s:%s", method, user.TrojanPassword)
			encodedUserInfo := base64.RawURLEncoding.EncodeToString([]byte(ssUserInfo))
			subContent += fmt.Sprintf("ss://%s@%s:%d#%s\n",
				encodedUserInfo, node.Address, node.Port, node.Name)
		}
	}

	// Base64 encode the output configuration for standard clients (Shadowrocket / Surge / V2rayN)
	encodedSub := base64.StdEncoding.EncodeToString([]byte(subContent))
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, encodedSub)
}
