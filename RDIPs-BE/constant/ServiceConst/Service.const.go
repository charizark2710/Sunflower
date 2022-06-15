package ServiceConst

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/services"
)

var ServicesMap = map[string]interface{}{
	"GET/devices":  services.GetAllDevices,
	"POST/devices": services.PostDevice,
}
