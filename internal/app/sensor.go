package app

import (
	"context"
	"time"
)

func (m *MicroserviceServer) GetGroupTransparencyAverage(ctx context.Context, groupName string) (averageTransparency float64, err error) {

	averageTransparency, err = m.SensorService.GetGroupTransparencyAverage(ctx, groupName)
	return
}

func (m *MicroserviceServer) GetGroupTemperatureAverage(ctx context.Context, groupName string) (averageTemperature float64, err error) {

	averageTemperature, err = m.SensorService.GetGroupTemperatureAverage(ctx, groupName)
	return
}

func (m *MicroserviceServer) GetGroupSpecies(ctx context.Context, groupName string) (speciesList map[string]int, err error) {

	speciesList, err = m.SensorService.GetGroupSpecies(ctx, groupName)
	return
}

func (m *MicroserviceServer) GetTopNGroupSpecies(ctx context.Context, groupName string, n int, from, till *time.Time) (speciesList map[string]int, err error) {

	speciesList, err = m.SensorService.GetTopNGroupSpecies(ctx, groupName, n, from, till)
	return
}

func (m *MicroserviceServer) GetRegionMinTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (minTemperature float64, err error) {

	minTemperature, err = m.SensorService.GetRegionMinTemperature(ctx, xMin, xMax, yMin, yMax, zMin, zMax)
	return
}

func (m *MicroserviceServer) GetRegionMaxTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (maxTemperature float64, err error) {

	maxTemperature, err = m.SensorService.GetRegionMaxTemperature(ctx, xMin, xMax, yMin, yMax, zMin, zMax)
	return
}

func (m *MicroserviceServer) GetCodeNameTemperatureAverage(ctx context.Context, codeName string, from, till time.Time) (averageTemperature float64, err error) {

	averageTemperature, err = m.SensorService.GetCodeNameTemperatureAverage(ctx, codeName, from, till)
	return
}
