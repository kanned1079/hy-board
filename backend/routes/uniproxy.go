package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

var (
	trafficRateCache = make(map[uint]float32)
	cacheMu          sync.RWMutex
)

func GetNodeTrafficRate(nodeID uint) float32 {
	cacheMu.RLock()
	rate, exists := trafficRateCache[nodeID]
	cacheMu.RUnlock()
	if exists {
		return rate
	}

	// Fallback to database query
	var node models.Node
	if err := database.DB.First(&node, nodeID).Error; err == nil {
		cacheMu.Lock()
		trafficRateCache[nodeID] = node.TrafficRate
		cacheMu.Unlock()
		return node.TrafficRate
	}

	return 1.0
}

func UpdateNodeTrafficRateCache(nodeID uint, rate float32) {
	cacheMu.Lock()
	trafficRateCache[nodeID] = rate
	cacheMu.Unlock()
}

func DeleteNodeTrafficRateCache(nodeID uint) {
	cacheMu.Lock()
	delete(trafficRateCache, nodeID)
	cacheMu.Unlock()
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

	if !node.Show {
		c.JSON(http.StatusForbidden, gin.H{
			"error":    "node is disabled",
			"disabled": true,
		})
		return
	}

	// Update node's updated_at timestamp to register a heartbeat
	database.DB.Model(&node).Update("updated_at", time.Now())

	// Update traffic rate cache
	UpdateNodeTrafficRateCache(node.ID, node.TrafficRate)

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
	nodeIDStr := c.Query("node_id")

	var users []models.User
	query := database.DB.Where("status = ? AND (expired_at IS NULL OR expired_at > ?)", 1, time.Now())

	if nodeIDStr != "" {
		if nodeID, err := strconv.ParseUint(nodeIDStr, 10, 32); err == nil {
			var node models.Node
			if err := database.DB.First(&node, uint(nodeID)).Error; err == nil {
				if !node.Show {
					c.JSON(http.StatusForbidden, gin.H{
						"error":    "node is disabled",
						"disabled": true,
					})
					return
				}
				// Parse allowed group IDs
				parts := strings.Split(node.GroupIDs, ",")
				var allowedGroupIDs []uint
				for _, p := range parts {
					if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
						allowedGroupIDs = append(allowedGroupIDs, uint(id))
					}
				}
				if len(allowedGroupIDs) > 0 {
					query = query.Where("group_id IN ?", allowedGroupIDs)
				} else {
					query = query.Where("group_id = ?", node.GroupID)
				}
			}
		}
	}

	err := query.Find(&users).Error
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

	trafficRate := GetNodeTrafficRate(payload.NodeID)

	tx := database.DB.Begin()
	for _, t := range payload.Traffic {
		upMultiplied := uint64(float64(t.Up) * float64(trafficRate))
		downMultiplied := uint64(float64(t.Down) * float64(trafficRate))

		// Increment user's used traffic
		tx.Model(&models.User{}).Where("id = ?", t.UserID).UpdateColumn("used_traffic", gorm.Expr("used_traffic + ?", upMultiplied+downMultiplied))

		// Log traffic
		log := models.TrafficLog{
			UserID:    t.UserID,
			NodeID:    payload.NodeID,
			Up:        upMultiplied,
			Down:      downMultiplied,
			CreatedAt: time.Now(),
		}
		tx.Create(&log)
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
