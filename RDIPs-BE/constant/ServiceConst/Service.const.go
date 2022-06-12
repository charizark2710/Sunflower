package ServiceConst

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/services"
)

var ServicesMap = map[string]interface{}{
	"/devices": services.GetAllDevices,
}
