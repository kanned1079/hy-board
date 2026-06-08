package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"hy-board-backend/database"
	"hy-board-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TrafficLogPayload struct {
	UserID uint   `json:"user_id"`
	Up     uint64 `json:"up"`
	Down   uint64 `json:"down"`
}

type PushPayload struct {
	NodeID  uint                `json:"node_id"`
	Traffic []TrafficLogPayload `json:"traffic"`
}

// GetNodeConfig handles GET /api/v1/server/UniProxy/config
func GetNodeConfig(c *gin.Context) {
	nodeIDStr := c.Query("node_id")
	nodeID, err := strconv.ParseUint(nodeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node_id"})
		return
	}

	var node models.Node
	if err := database.DB.First(&node, uint(nodeID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	// Deserialize settings JSON
	var settingsMap map[string]interface{}
	if node.Settings != "" {
		_ = json.Unmarshal([]byte(node.Settings), &settingsMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           node.ID,
			"port":         node.Port,
			"address":      node.Address,
			"type":         node.Type,
			"traffic_rate": node.TrafficRate,
			"settings":     settingsMap,
		},
	})
}

// GetNodeUsers handles GET /api/v1/server/UniProxy/user
func GetNodeUsers(c *gin.Context) {
	var users []models.User
	// Fetch active users who are not expired
	err := database.DB.Where("status = ? AND (expired_at IS NULL OR expired_at > ?)", 1, time.Now()).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	type NodeUserResponse struct {
		ID          uint   `json:"id"`
		UUID        string `json:"uuid"`
		Password    string `json:"password"`
		SpeedLimit  uint32 `json:"speed_limit"`
		DeviceLimit uint32 `json:"device_limit"`
	}

	res := make([]NodeUserResponse, len(users))
	for i, u := range users {
		res[i] = NodeUserResponse{
			ID:          u.ID,
			UUID:        u.V2rayUUID,
			Password:    u.TrojanPassword,
			SpeedLimit:  u.SpeedLimit,
			DeviceLimit: u.DeviceLimit,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

// PushNodeData handles POST /api/v1/server/UniProxy/push
func PushNodeData(c *gin.Context) {
	var payload PushPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	tx := database.DB.Begin()
	for _, t := range payload.Traffic {
		// Increment user's used traffic
		tx.Model(&models.User{}).Where("id = ?", t.UserID).UpdateColumn("used_traffic", gorm.Expr("used_traffic + ?", t.Up+t.Down))

		// Log traffic
		log := models.TrafficLog{
			UserID:    t.UserID,
			NodeID:    payload.NodeID,
			Up:        t.Up,
			Down:      t.Down,
			CreatedAt: time.Now(),
		}
		tx.Create(&log)
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
