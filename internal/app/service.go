package app

import "github.com/sensors/internal/service"

type MicroserviceServer struct {
	SensorService service.SensorService
}

func NewMicroservice(sensorService service.SensorService) *MicroserviceServer {

	return &MicroserviceServer{
		SensorService: sensorService,
	}
}
