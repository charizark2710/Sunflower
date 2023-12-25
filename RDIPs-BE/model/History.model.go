package model

import "time"

type History struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LogPath   string    `json:"log_path" validate:"required"`
	Payload   string    `json:"payload"`
}

type SysHistory struct {
	Id        string    `gorm:"default:gen_random_uuid();primaryKey;column:id;type:uuid"`
	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
	LogPath   string    `gorm:"column:log_path;type:string;size:256"`
}

func (SysHistory) TableName() string {
	return "sunflower.sys_history"
}

func (in *SysHistory) ConvertToJson(out *History) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.LogPath = in.LogPath
}

func (in *History) ConvertToDB(out *SysHistory) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.LogPath = in.LogPath
}
