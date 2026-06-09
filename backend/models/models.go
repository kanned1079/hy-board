package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Email          string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username       string         `gorm:"size:100" json:"username"`
	Password       string         `gorm:"not null;size:255" json:"-"`
	V2rayUUID      string         `gorm:"uniqueIndex;size:36" json:"v2ray_uuid"`
	TrojanPassword string         `gorm:"size:64" json:"trojan_password"`
	SpeedLimit     uint32         `gorm:"default:0" json:"speed_limit"`   // Mbps (0 = unlimited)
	DeviceLimit    uint32         `gorm:"default:0" json:"device_limit"`  // Max devices (0 = unlimited)
	TotalTraffic   uint64         `gorm:"default:0" json:"total_traffic"` // Bytes
	UsedTraffic    uint64         `gorm:"default:0" json:"used_traffic"`  // Bytes
	ExpiredAt      time.Time      `json:"expired_at"`
	Status         int8           `gorm:"default:1" json:"status"` // 1 = Active, 0 = Disabled
	IsAdmin        bool           `gorm:"default:false" json:"is_admin"`
	Balance        float64        `gorm:"default:0.0" json:"balance"`
	GroupID        uint           `gorm:"default:1" json:"group_id"` // Subscription Level (1 = S1, 2 = S2, 99 = Admin)
	RegisterIP     string         `gorm:"size:45" json:"register_ip"`
}

type Node struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Type        string         `gorm:"not null;size:20" json:"type"` // V2ray, Vless, Trojan, Shadowsocks
	Address     string         `gorm:"not null;size:255" json:"address"`
	Port        uint16         `gorm:"not null" json:"port"`
	TrafficRate float32        `gorm:"default:1.0" json:"traffic_rate"`
	Settings    string         `gorm:"type:text" json:"settings"` // JSON string config
	Show        bool           `gorm:"default:true" json:"show"`
	GroupID     uint           `gorm:"default:1" json:"group_id"`                      // Deprecated in favor of GroupIDs
	GroupIDs    string         `gorm:"type:varchar(255);default:'1'" json:"group_ids"` // Comma-separated group IDs required (e.g. "1,2")
}

type TrafficLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	NodeID    uint      `gorm:"index" json:"node_id"`
	Node      Node      `gorm:"foreignKey:NodeID" json:"node"`
	Up        uint64    `gorm:"default:0" json:"up"`
	Down      uint64    `gorm:"default:0" json:"down"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

type Announcement struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"not null;size:255" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Show      bool           `gorm:"default:true" json:"show"`
}

type Knowledge struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"not null;size:255" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Show      bool           `gorm:"default:true" json:"show"`
	Sort      int            `gorm:"default:0" json:"sort"`
}

type Ticket struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
	UserID    uint            `gorm:"index" json:"user_id"`
	User      User            `gorm:"foreignKey:UserID" json:"user"`
	Title     string          `gorm:"not null;size:255" json:"title"`
	Status    string          `gorm:"default:'open';size:20" json:"status"` // open, closed
	Messages  []TicketMessage `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE" json:"messages"`
}

type TicketMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TicketID  uint      `gorm:"index" json:"ticket_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
}

type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
}

type Plan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Price       float64        `gorm:"not null;default:0.0" json:"price"`
	Traffic     uint64         `gorm:"not null;default:0" json:"traffic"`      // Total traffic in GB
	SpeedLimit  uint32         `gorm:"not null;default:0" json:"speed_limit"`  // Mbps (0 = unlimited)
	DeviceLimit uint32         `gorm:"not null;default:0" json:"device_limit"` // Max devices (0 = unlimited)
	ExpiryDays  uint32         `gorm:"not null;default:30" json:"expiry_days"` // Duration of subscription plan
	GroupID     uint           `gorm:"not null;default:1" json:"group_id"`     // Maps to subscription level Group ID
	Show        bool           `gorm:"default:true" json:"show"`
}

type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null;size:100" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
