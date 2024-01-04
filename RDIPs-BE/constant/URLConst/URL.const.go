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
	GetDetailPerformance = "/performances/:id"
	PutDetailPerformance = "/performances/:deviceId"

	//History
	// PostHistory      = "/history"
	GetDetailHistory = "/history/:id"
	PutDetailHistory = "/history/:deviceId"

	//Weather
	GetWeatherForecast = "/weather/forecast"

	//Keycloak
	PostLogin         = "/login"
	PostKeycloakUsers = "/users"
	GetKeycloakUsers  = "/users"
)
