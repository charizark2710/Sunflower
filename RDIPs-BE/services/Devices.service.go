package services

import (
	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
	"github.com/charizark2710/Sunflower/RDIPs-BE/utils"
)

func GetAllDevices() (model.Device, error) {
	utils.Log(constant.Info, "Test")
	return model.Device{Id: "123", Name: "TEST"}, nil
}
