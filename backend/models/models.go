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
	Password       string         `gorm:"not null;size:255" json:"-"`
	V2rayUUID      string         `gorm:"uniqueIndex;size:36" json:"v2ray_uuid"`
	TrojanPassword string         `gorm:"size:64" json:"trojan_password"`
	SpeedLimit     uint32         `gorm:"default:0" json:"speed_limit"`   // Mbps (0 = unlimited)
	DeviceLimit    uint32         `gorm:"default:0" json:"device_limit"`  // Max devices (0 = unlimited)
	TotalTraffic   uint64         `gorm:"default:0" json:"total_traffic"` // Bytes
	UsedTraffic    uint64         `gorm:"default:0" json:"used_traffic"`  // Bytes
	ExpiredAt      time.Time      `json:"expired_at"`
	Status         int8           `gorm:"default:1" json:"status"`   // 1 = Active, 0 = Disabled
	IsAdmin        bool           `gorm:"default:false" json:"is_admin"`
	Balance        float64        `gorm:"default:0.0" json:"balance"`
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
}

type TrafficLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	NodeID    uint      `gorm:"index" json:"node_id"`
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
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Title     string         `gorm:"not null;size:255" json:"title"`
	Status    string         `gorm:"default:'open';size:20" json:"status"` // open, closed
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
