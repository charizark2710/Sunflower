package ServiceConst

import (
	urlconst "RDIPs-BE/constant/URLConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/services"
)

type ServiceFn func(*commonModel.ServiceContext) (commonModel.ResponseTemplate, error)

var ServicesMap = map[string]ServiceFn{
	"GET" + urlconst.GetAllDevices:    services.GetAllDevices,
	"POST" + urlconst.PostDevice:      services.PostDevice,
	"GET" + urlconst.GetDetailDevice:  services.GetDetailDevice,
	"PUT" + urlconst.PutDetailDevice:  services.UpdateDevice,
	"DELETE" + urlconst.DeleteDevice:  services.DeleteDevice,
	"GET" + urlconst.GetLogOfDevice:   services.GetLogOfDevice,
	"POST" + urlconst.PostLogOfDevice: services.PostLogOfDevice,

	//Performances
	"GET" + urlconst.GetAllPerformances: services.GetAllPerformances,
	// "POST" + urlconst.PostPerformance:     services.PostPerformance,
	"GET" + urlconst.GetDetailPerformance: services.GetDetailPerformance,
	"PUT" + urlconst.PutDetailPerformance: services.PutPerformance,

	//History
	// "POST" + urlconst.PostHistory:     services.PostHistory,
	"GET" + urlconst.GetDetailHistory: services.GetDetailHistory,
	"PUT" + urlconst.PutDetailHistory: services.UpdateHistory,

	//Weather
	"GET" + urlconst.GetWeatherForecast: services.GetWeatherForecast,

	//Keycloak
	"POST" + urlconst.PostLogin:            services.Login,
	"GET" + urlconst.GetKeycloakUsers:      services.GetKeycloakUsers,
	"POST" + urlconst.PostKeycloakUser:     services.PostKeycloakUser,
	"PUT" + urlconst.PutKeycloakUsers:      services.PutKeycloakUser,
	"DELETE" + urlconst.DeleteKeycloakUser: services.DeleteKeycloakUser,
}

var ServiceMapMQTT = map[string]string{
	// "GetAllDevices":   "GET" + urlconst.GetAllDevices,
	"PostDevice": "POST" + urlconst.PostDevice,
	// "GetDetailDevice": "GET" + urlconst.GetDetailDevice,
	"PutDetailDevice": "PUT" + urlconst.PutDetailDevice,
	// "DeleteDevice":    "DELETE" + urlconst.DeleteDevice,
	// "GetLogOfDevice":  "GET" + urlconst.GetLogOfDevice,
	"PostLogOfDevice": "POST" + urlconst.PostLogOfDevice,

	//Performances
	// "GetAllPerformances":   "GET" + urlconst.GetAllPerformances,
	// "GetDetailPerformance": "GET" + urlconst.GetDetailPerformance,
	"PutDetailPerformance": "PUT" + urlconst.PutDetailPerformance,

	//History
	// "GetDetailHistory": "GET" + urlconst.GetDetailHistory,
	"PutDetailHistory": "PUT" + urlconst.PutDetailHistory,
	//Weather
	"GetWeatherForecast": "GET" + urlconst.GetWeatherForecast,
}
