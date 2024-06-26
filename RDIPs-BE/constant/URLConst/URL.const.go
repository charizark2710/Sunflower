package urlconst

const (
	GetAllDevices   = "/devices"
	PostDevice      = "/devices"
	GetDetailDevice = "/devices/:id"
	PutDetailDevice = "/devices/:id"
	DeleteDevice    = "/devices/:id"
	GetLogOfDevice  = "/device/:deviceID/logs/:dateMilisec"
	PostLogOfDevice = "/device/:deviceID/logs"

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

	// Authenticate
	GetLoginScreen = "/login"
	Callback       = "/callback"

	//Keycloak - Users
	PostKeycloakUser    = "/users"
	GetKeycloakUsers    = "/users"
	GetKeycloakUserById = "/users/:id"
	PutKeycloakUsers    = "/users/:id"
	DeleteKeycloakUser  = "/users/:id"

	//Keycloak - Groups
	GetKeycloakGroups    = "/groups"
	GetKeycloakGroupById = "/groups/:id"
	DeleteKeycloakGroup  = "/groups/:id"
	PostKeycloakGroup    = "/groups"
	PutKeycloakGroup     = "/groups/:id"
)
