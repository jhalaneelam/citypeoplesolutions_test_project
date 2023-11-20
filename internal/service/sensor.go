package service

import (
	"context"
	"strconv"
	"time"

	"github.com/sensors/internal/repository"
)

type SensorService interface {
	GetGroupTransparencyAverage(ctx context.Context, groupName string) (float64, error)
	GetGroupTemperatureAverage(ctx context.Context, groupName string) (float64, error)
	GetGroupSpecies(ctx context.Context, groupName string) (map[string]int, error)
	GetTopNGroupSpecies(ctx context.Context, groupName string, n int, from, till *time.Time) (map[string]int, error)
	GetRegionMinTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (float64, error)
	GetRegionMaxTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (float64, error)
	GetCodeNameTemperatureAverage(ctx context.Context, codeName string, from, till time.Time) (float64, error)
}

type sensorService struct {
	dao repository.DAO
}

func NewSensorService(dao repository.DAO) SensorService {
	return &sensorService{dao: dao}
}

func (s *sensorService) GetGroupTransparencyAverage(ctx context.Context, groupName string) (averageTransparency float64, err error) {

	// Check Redis cache first
	cacheKey := "transparency:" + groupName
	redisClient := repository.NewClient()
	val, err := redisClient.Get(redisClient.Context(), cacheKey).Result()
	if err == nil {
		averageTransparency, _ = strconv.ParseFloat(val, 64)
		return
	}

	averageTransparency, err = s.dao.NewSensorQuery().FetchAverageTransparency(groupName)
	if err != nil {
		return
	}

	redisClient.Set(redisClient.Context(), cacheKey, averageTransparency, 10*time.Second)

	return
}

func (s *sensorService) GetGroupTemperatureAverage(ctx context.Context, groupName string) (averageTemperature float64, err error) {

	// Check Redis cache first
	cacheKey := "temperature:" + groupName
	redisClient := repository.NewClient()
	val, err := redisClient.Get(redisClient.Context(), cacheKey).Result()
	if err == nil {
		averageTemperature, _ = strconv.ParseFloat(val, 64)
		return
	}

	averageTemperature, err = s.dao.NewSensorQuery().FetchAverageTemperature(groupName)
	if err != nil {
		return
	}

	redisClient.Set(redisClient.Context(), cacheKey, averageTemperature, 10*time.Second)

	return
}

func (s *sensorService) GetGroupSpecies(ctx context.Context, groupName string) (speciesList map[string]int, err error) {

	speciesList, err = s.dao.NewSensorQuery().FetchSpeciesList(groupName)
	if err != nil {
		return
	}

	return
}

func (s *sensorService) GetTopNGroupSpecies(ctx context.Context, groupName string, n int, from, till *time.Time) (speciesList map[string]int, err error) {

	speciesList, err = s.dao.NewSensorQuery().FetchTopNSpeciesList(groupName, n, from, till)
	if err != nil {
		return
	}

	return
}

func (s *sensorService) GetRegionMinTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (minTemperature float64, err error) {

	minTemperature, err = s.dao.NewSensorQuery().GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		return
	}

	return
}

func (s *sensorService) GetRegionMaxTemperature(ctx context.Context, xMin, xMax, yMin, yMax, zMin, zMax float64) (maxTemperature float64, err error) {

	maxTemperature, err = s.dao.NewSensorQuery().GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		return
	}

	return
}

func (s *sensorService) GetCodeNameTemperatureAverage(ctx context.Context, codeName string, from, till time.Time) (averageTemperature float64, err error) {

	// Check Redis cache first
	cacheKey := "temperature:" + codeName
	redisClient := repository.NewClient()
	val, err := redisClient.Get(redisClient.Context(), cacheKey).Result()
	if err == nil {
		averageTemperature, _ = strconv.ParseFloat(val, 64)
		return
	}

	averageTemperature, err = s.dao.NewSensorQuery().FetchCodeNameAverageTemperature(codeName, from, till)
	if err != nil {
		return
	}

	redisClient.Set(redisClient.Context(), cacheKey, averageTemperature, 10*time.Second)

	return
}
