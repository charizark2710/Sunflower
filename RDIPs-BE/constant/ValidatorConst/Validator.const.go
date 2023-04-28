package ValidatorConst

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/validators"
)

var ValidatorsMap = map[string]validators.Validator{
	"POST" + urlconst.PostDevice:     validators.DeviceValidator{},
	"PUT" + urlconst.PutDetailDevice: validators.DeviceValidator{},
	// "DELETE" + urlconst.DeleteDevice: services.DeleteDevice,

	//Performances
	// "POST" + urlconst.PostPerformance:     validators.PostPerformanceValidator,
	"PUT" + urlconst.PutDetailPerformance: validators.PerformanceValidator{},

	//History
	// "POST" + urlconst.PostHistory: validators.PostHistoryValidator,
}
