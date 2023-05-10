package ServiceConst

import (
	urlconst "RDIPs-BE/constant/URLConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/services"
)

type ServiceFn func(*commonModel.ServiceContext) (commonModel.ResponseTemplate, error)

var ServicesMap = map[string]ServiceFn{
	"GET" + urlconst.GetAllDevices:   services.GetAllDevices,
	"POST" + urlconst.PostDevice:     services.PostDevice,
	"GET" + urlconst.GetDetailDevice: services.GetDetailDevice,
	"PUT" + urlconst.PutDetailDevice: services.UpdateDevice,
	"DELETE" + urlconst.DeleteDevice: services.DeleteDevice,

	//Performances
	"GET" + urlconst.GetAllPerformances: services.GetAllPerformances,
	// "POST" + urlconst.PostPerformance:     services.PostPerformance,
	"GET" + urlconst.GetDetailPerformance: services.GetDetailPerformance,
	"PUT" + urlconst.PutDetailPerformance: services.PutPerformance,

	//History
	// "POST" + urlconst.PostHistory:     services.PostHistory,
	"GET" + urlconst.GetDetailHistory: services.GetDetailHistory,
	"PUT" + urlconst.PutDetailHistory: services.UpdateHistory,
}
