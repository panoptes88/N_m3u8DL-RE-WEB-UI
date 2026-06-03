package model

import (
	"time"

	"gorm.io/gorm"
)

// DownloadProfile 下载方案
type DownloadProfile struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Name             string         `gorm:"size:255;not null" json:"name"`                 // 方案名称
	Domain           string         `gorm:"size:255;index" json:"domain"`                  // 域名，用于自动匹配
	ThreadCount      int            `gorm:"default:32" json:"thread_count"`
	RetryCount       int            `gorm:"default:15" json:"retry_count"`
	Headers          string         `gorm:"size:2048" json:"headers,omitempty"`
	BaseURL          string         `gorm:"size:1024" json:"base_url,omitempty"`
	DelAfterDone     bool           `gorm:"default:true" json:"del_after_done"`
	BinaryMerge      bool           `gorm:"default:false" json:"binary_merge"`
	AutoSelect       bool           `gorm:"default:false" json:"auto_select"`
	Key              string         `gorm:"size:512" json:"key,omitempty"`
	DecryptionEngine string         `gorm:"size:32;default:'MP4DECRYPT'" json:"decryption_engine"`
	CustomArgs       string         `gorm:"size:2048" json:"custom_args,omitempty"`
	CustomProxy      string         `gorm:"size:512" json:"custom_proxy,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DownloadProfile) TableName() string {
	return "download_profiles"
}
