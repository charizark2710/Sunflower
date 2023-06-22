package urlconst

const (
	GetAllDevices   = "/devices"
	PostDevice      = "/devices"
	GetDetailDevice = "/devices/:id/"
	PutDetailDevice = "/devices/:id"
	DeleteDevice    = "/devices/:id"

	//Performances
	GetAllPerformances = "/performances"
	// PostPerformance      = "/performances"
	GetDetailPerformance = "/performances/:deviceId"
	PutDetailPerformance = "/performances/:deviceId"

	//History
	// PostHistory      = "/history"
	GetDetailHistory = "/history/:deviceId"
	PutDetailHistory = "/history/:deviceId"

	//Weather
	GetWeatherForecast = "/weather/forecast"
)
