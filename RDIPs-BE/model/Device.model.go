package model

import "time"

type statusEnum string
type deviceType string

const (
	Active   statusEnum = "Active"
	Sleep    statusEnum = "Sleep"
	Warning  statusEnum = "Warning"
	Error    statusEnum = "Error"
	Critical statusEnum = "Critical"
	Disable  statusEnum = "Disable"
)

const (
	Common deviceType = "Common"
	Family deviceType = "Family"
)

type Devices struct {
	Id             string       `json:"id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Type           deviceType   `json:"type"`
	Status         statusEnum   `json:"status"`
	LifeTime       time.Time    `json:"life_time"`
	FirwareVersion int          `json:"firmware_ver"`
	AppVersion     int          `json:"app_ver"`
	ParentID       string       `json:"parentID"`
	Parent         *Devices     `json:"parent,omitempty"`
	Name           string       `json:"name" validate:"required"`
	HistoryID      string       `json:"historyID,omitempty"`
	History        *History     `json:"history,omitempty"`
	PerformanceID  string       `json:"performanceID,omitempty"`
	Performance    *Performance `json:"performance,omitempty"`
}

type SysDevices struct {
	// User      SysUser    `gorm:"column:user;type:uuid"`
	Id             string       `gorm:"default:gen_random_uuid();primaryKey;column:id;type:uuid"`
	CreatedAt      time.Time    `gorm:"column:created_at;"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;"`
	Type           deviceType   `gorm:"column:type;"`
	Status         statusEnum   `gorm:"column:status;default:active"`
	LifeTime       time.Time    `gorm:"column:life_time"`
	FirwareVersion int          `gorm:"column:firmware_ver;type:integer"`
	AppVersion     int          `gorm:"column:app_ver;type:integer"`
	ParentID       string       `gorm:"column:parent;uniqueIndex;default:NULL"`
	Parent         *SysDevices  `gorm:"foreignKey:ParentID"`
	Name           string       `gorm:"column:name;unique"`
	DeviceRel      SysDeviceRel `gorm:"foreignKey:DeviceID"`
}

func (SysDevices) TableName() string {
	return "sunflower.sys_devices"
}

func (in SysDevices) ConvertToJson(out *Devices) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Type = in.Type
	out.Status = in.Status
	out.LifeTime = in.LifeTime
	out.FirwareVersion = in.FirwareVersion
	out.AppVersion = in.AppVersion
	out.ParentID = in.ParentID
	if in.Parent != nil {
		out.Parent = &Devices{}
		in.Parent.ConvertToJson(out.Parent)
	}
	out.Name = in.Name
	if in.DeviceRel.HistoryID != "" {
		if in.DeviceRel.History.Id == "" {
			out.HistoryID = in.DeviceRel.HistoryID
		} else {
			out.History = &History{}
			in.DeviceRel.History.ConvertToJson(out.History)
		}
	}

	if in.DeviceRel.PerformanceID != "" {
		if in.DeviceRel.Performance.Id == "" {
			out.PerformanceID = in.DeviceRel.PerformanceID
		} else {
			out.Performance = &Performance{}
			in.DeviceRel.Performance.ConvertToJson(out.Performance)
		}
	}

}

func (in Devices) ConvertToDB(out *SysDevices) {
	out.Id = in.Id
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Type = in.Type
	out.Status = in.Status
	out.LifeTime = in.LifeTime
	out.FirwareVersion = in.FirwareVersion
	out.AppVersion = in.AppVersion
	out.ParentID = in.ParentID
	if in.Parent != nil {
		out.Parent = &SysDevices{}
		in.Parent.ConvertToDB(out.Parent)
	}
	out.Name = in.Name
}
