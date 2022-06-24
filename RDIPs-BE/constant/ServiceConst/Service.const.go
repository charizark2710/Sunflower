package ServiceConst

import (
	urlconst "github.com/charizark2710/Sunflower/RDIPs-BE/constant/URLConst"
	"github.com/charizark2710/Sunflower/RDIPs-BE/services"
)

var ServicesMap = map[string]interface{}{
	"GET" + urlconst.GetAllDevices:   services.GetAllDevices,
	"POST" + urlconst.PostDevice:     services.PostDevice,
	"GET" + urlconst.GetDetailDevice: services.GetDetailDevice,
	"PUT" + urlconst.PutDetailDevice: services.UpdateDevice,
	"DELETE" + urlconst.DeleteDevice: services.DeleteDevice,
}
