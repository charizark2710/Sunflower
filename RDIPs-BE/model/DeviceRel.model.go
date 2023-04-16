package model

type DeviceRel struct {
	DeviceID      string `json:"device_id"`
	HistoryID     string `json:"history_id"`
	PerformanceID string `json:"performance_id"`
}

type DeviceRelWithDetail struct {
	DeviceID      string      `json:"device_id"`
	HistoryID     string      `json:"history_id"`
	History       History     `json:"history,omitempty"`
	PerformanceID string      `json:"performance_id"`
	Performance   Performance `json:"performance,omitempty"`
}

type SysDeviceRel struct {
	DeviceID      string         `gorm:"column:device_id;type:uuid"`
	HistoryID     string         `gorm:"column:history_id;type:uuid"`
	History       SysHistory     `gorm:"foreignKey:HistoryID"`
	PerformanceID string         `gorm:"column:performance_id;type:uuid"`
	Performance   SysPerformance `gorm:"foreignKey:PerformanceID"`
}

func (SysDeviceRel) TableName() string {
	return "sunflower.sys_device_rel"
}

func (in SysDeviceRel) ConvertToJson(out *DeviceRel) {
	out.HistoryID = in.HistoryID
	out.PerformanceID = in.PerformanceID
	out.DeviceID = in.DeviceID
}

func (in SysDeviceRel) ConvertToJsonWithDetail(out *DeviceRelWithDetail) {
	out.HistoryID = in.HistoryID
	out.PerformanceID = in.PerformanceID
	out.DeviceID = in.DeviceID
	in.History.ConvertToJson(&out.History)
	in.Performance.ConvertToJson(&out.Performance)

}

func (in DeviceRel) ConvertToDB(out *SysDeviceRel) {
	out.HistoryID = in.HistoryID
	out.PerformanceID = in.PerformanceID
	out.DeviceID = in.DeviceID
}
