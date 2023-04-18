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
	// "POST" + urlconst.PostPerformance:     services.PostPerformance,
	// "PUT" + urlconst.PutDetailPerformance: services.PutPerformance,

	//History
	// "POST" + urlconst.PostHistory:     services.PostHistory,
}
