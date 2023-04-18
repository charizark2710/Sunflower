package ValidatorConst

import (
	urlconst "RDIPs-BE/constant/URLConst"
	"RDIPs-BE/validators"
)

var ValidatorsMap = map[string]interface{}{
	"POST" + urlconst.PostDevice:     validators.PostDeviceValidator,
	"PUT" + urlconst.PutDetailDevice: validators.UpdateDeviceValidator,
	// "DELETE" + urlconst.DeleteDevice: services.DeleteDevice,

	//Performances
	"POST" + urlconst.PostPerformance:     validators.PostPerformanceValidator,
	"PUT" + urlconst.PutDetailPerformance: validators.UpdatePerformanceValidator,

	//History
	"POST" + urlconst.PostHistory: validators.PostHistoryValidator,
}
