package urlconst

const (
	GetAllDevices   = "/devices"
	PostDevice      = "/devices"
	GetDetailDevice = "/devices/:id"
	PutDetailDevice = "/devices/:id"
	DeleteDevice    = "/devices/:id"

	//Performances
	GetAllPerformances   = "/performances"
	PostPerformance      = "/performances"
	GetDetailPerformance = "/performances/:id"
	PutDetailPerformance = "/performances/:id"
	PostHistory          = "/history"
	GetDetailHistory     = "/history/:id"
	GetHistoriesOfDevice = "/history/devices/:id"
)
