// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sensors/api"
	"github.com/sensors/internal/app"
	"github.com/sensors/internal/repository"
	"github.com/sensors/internal/service"
)

func main() {

	// Initialize PostgreSQL DB
	db, err := repository.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var wg sync.WaitGroup

	redisClient := repository.NewClient()

	// Phase 1: One-time "Kickoff" Phase
	GenerateSensorGroupsAndSensors(context.Background(), db, redisClient)

	// Phase 2: Regularly Repeated Phase for Data Generation

	wg.Add(1)
	go generateSensorDataRegularly(&wg, context.Background(), db, redisClient)

	// Phase 3: Regularly Repeated Phase for Aggregate Statistics
	wg.Add(1)
	go aggregateStatisticsRegularly(&wg, context.Background(), db)

	dao := repository.NewDAO(db)

	sensor := service.NewSensorService(dao)

	server := api.NewServer(context.Background(), *app.NewMicroservice(sensor))
	wg.Add(1)
	go func() {
		log.Fatal(server.Start())
		wg.Done()
	}()
	wg.Wait()
}

func GenerateSensorGroupsAndSensors(ctx context.Context, db *sql.DB, redisClient *redis.Client) {

	createTables(db)

	groups := []string{"alpha", "beta", "gamma"}
	for index, group := range groups {
		sensorGroup := repository.SensorGroup{Name: group}

		insertData := insertSensorGroup(db, sensorGroup)

		if insertData {
			for sensorIndex := 0; sensorIndex < 3; sensorIndex++ {
				sensor := repository.Sensor{
					GroupID:  index + 1,
					Codename: fmt.Sprintf("%s%d", group, sensorIndex+1),
					Index:    sensorIndex + 1,
					X:        rand.Float64() * 10,
					Y:        rand.Float64() * 10,
					Z:        rand.Float64() * 10,
					DataRate: 60,
				}

				insertSensor(db, sensor)

				fishSpecies := []string{"Atlantic Cod", "Sailfish", "Tuna", "Salmon", "Trout", "Barracuda"}

				randomFish := getRandomFishSpecies(fishSpecies)

				data := repository.SensorData{
					SensorID:         sensor.ID,
					Temperature:      generateTemperature(sensor.Z),
					Transparency:     generateTransparency(redisClient, sensor.Z),
					FishSpeciesName:  randomFish,
					FishSpeciesCount: rand.Intn(20),
					CreatedAt:        time.Now(),
				}

				insertSensorData(db, data)
			}
		}
	}
}

// createTables creates the necessary tables if they do not exist
func createTables(db *sql.DB) {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS sensors (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL,
			codename VARCHAR(255) NOT NULL,
			index INT NOT NULL,
			x FLOAT NOT NULL,
			y FLOAT NOT NULL,
			z FLOAT NOT NULL,
			data_rate INT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS sensor_data (
			id SERIAL PRIMARY KEY,
			sensor_id INT NOT NULL,
			temperature FLOAT NOT NULL,
			transparency INT NOT NULL,
			fish_species_name VARCHAR(255) NOT NULL,
			fish_species_count INT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS aggregated_statistics (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL,
			average_temperature DOUBLE PRECISION,
			average_transparency INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES sensor_groups(id)
		);
	`)
	if err != nil {
		panic(err)
	}
}

func insertSensorGroup(db *sql.DB, group repository.SensorGroup) (isEmpty bool) {
	var count int
	var lastInsertID int64
	isEmpty = false
	err := db.QueryRow("SELECT COUNT(*) FROM sensor_groups WHERE name = $1", group.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		isEmpty = true
		err = db.QueryRow(`
			INSERT INTO sensor_groups (name)
			VALUES ($1)
			RETURNING id;
		`, group.Name).Scan(&lastInsertID)
		if err != nil {
			panic(err)
		}
		group.ID = int(lastInsertID)
	}
	return
}

func insertSensor(db *sql.DB, sensor repository.Sensor) {
	var lastInsertID int64
	err := db.QueryRow(`
		INSERT INTO sensors (group_id, codename, index, x, y, z, data_rate)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`, sensor.GroupID, sensor.Codename, sensor.Index, sensor.X, sensor.Y, sensor.Z, sensor.DataRate).Scan(&lastInsertID)
	if err != nil {
		panic(err)
	}

	sensor.ID = int(lastInsertID)
}

// insertSensorData inserts sensor data into the sensor_data table
func insertSensorData(db *sql.DB, data repository.SensorData) {
	_, err := db.Exec(`
		INSERT INTO sensor_data (sensor_id, temperature, transparency, fish_species_name, fish_species_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6);
	`, data.SensorID, data.Temperature, data.Transparency, data.FishSpeciesName, data.FishSpeciesCount, data.CreatedAt)
	if err != nil {
		panic(err)
	}
}

// generateAndInsertSensorData generates random sensor data and inserts it into the database
func generateSensorDataRegularly(wg *sync.WaitGroup, ctx context.Context, db *sql.DB, redisClient *redis.Client) {
	defer wg.Done()

	for {

		rows, err := db.QueryContext(ctx, "SELECT id, group_id, codename, index, x, y, z, data_rate FROM sensors")
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var sensor repository.Sensor
			if err := rows.Scan(&sensor.ID, &sensor.GroupID, &sensor.Codename, &sensor.Index, &sensor.X, &sensor.Y, &sensor.Z, &sensor.DataRate); err != nil {
				panic(err)
			}

			fishSpecies := []string{"Atlantic Cod", "Sailfish", "Tuna", "Salmon", "Trout", "Barracuda"}

			randomFish := getRandomFishSpecies(fishSpecies)
			data := repository.SensorData{
				SensorID:         sensor.ID,
				Temperature:      generateTemperature(sensor.Z),
				Transparency:     generateTransparency(redisClient, sensor.Z),
				FishSpeciesName:  randomFish,
				FishSpeciesCount: rand.Intn(20),
				CreatedAt:        time.Now(),
			}

			insertSensorData(db, data)
		}

		rows.Close()
	}
}

// generateTemperature generates temperature based on the depth (Z-axis)
func generateTemperature(depth float64) float64 {
	// Adjust the formula based on your requirements
	return 10 + depth*2 + rand.Float64()*5
}

// generateTransparency generates transparency with the constraint of not differing too much for nearby sensors
func generateTransparency(redisClient *redis.Client, depth float64) int {
	// Adjust the formula based on your requirements
	baseTransparency := rand.Intn(100)
	redisKey := fmt.Sprintf("transparency:depth:%.2f", depth)

	// Check Redis cache for previous transparency value
	if val, err := redisClient.Get(redisClient.Context(), redisKey).Result(); err == nil {
		cachedTransparency, _ := strconv.Atoi(val)
		// Ensure that the difference with the previous value is within a specified range
		if abs(baseTransparency-cachedTransparency) <= 10 {
			return baseTransparency
		}
	}

	// Store the new transparency value in Redis
	redisClient.Set(redisClient.Context(), redisKey, baseTransparency, 0)

	return baseTransparency
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getRandomFishSpecies(fishSpecies []string) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomIndex := random.Intn(len(fishSpecies))

	return fishSpecies[randomIndex]
}
func aggregateStatisticsRegularly(wg *sync.WaitGroup, ctx context.Context, db *sql.DB) {
	defer wg.Done()
	for {
		aggregateStatistics(db)
	}
}

func aggregateStatistics(db *sql.DB) {
	_, err := db.Exec(`
		INSERT INTO aggregated_statistics (group_id, average_temperature, average_transparency, created_at)
		SELECT 
			s.group_id,
			AVG(sd.temperature) AS average_temperature,
			AVG(sd.transparency) AS average_transparency,
			NOW() AS created_at
		FROM sensors s
		JOIN sensor_data sd ON s.id = sd.sensor_id
		GROUP BY s.group_id;
	`)
	if err != nil {
		panic(err)
	}
}
