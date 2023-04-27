package model

import (
	"time"
)

type Performance struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DocumentName string    `json:"document_name" validate:"required"`
	Payload      string    `json:"payload"`
}

type SysPerformance struct {
	Id           string    `gorm:"default:gen_random_uuid();primaryKey;column:id;type:uuid"`
	CreatedAt    time.Time `gorm:"column:created_at;"`
	UpdatedAt    time.Time `gorm:"column:updated_at;"`
	DocumentName string    `gorm:"column:document_name"`
}

func (SysPerformance) TableName() string {
	return "sunflower.sys_performance"
}

func (in SysPerformance) ConvertToJson(out *Performance) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.DocumentName = in.DocumentName
}

func (in Performance) ConvertToDB(out *SysPerformance) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.DocumentName = in.DocumentName
}
