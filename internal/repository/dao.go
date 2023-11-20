package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// SensorGroup represents the structure of a sensor group
type SensorGroup struct {
	ID   int
	Name string
}

// Sensor represents the structure of a sensor
type Sensor struct {
	ID       int
	GroupID  int
	Codename string
	Index    int
	X        float64
	Y        float64
	Z        float64
	DataRate int
}

// SensorData represents the structure of sensor data
type SensorData struct {
	ID               int
	SensorID         int
	Temperature      float64
	Transparency     int
	FishSpeciesName  string
	FishSpeciesCount int
	CreatedAt        time.Time
}

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root@123"
	dbname   = "sensors_db"
)

type DAO interface {
	NewSensorQuery() SensorQuery
}

type dao struct {
	DB *sql.DB
}

func NewDAO(db *sql.DB) DAO {
	return &dao{
		DB: db,
	}
}

func (d *dao) NewSensorQuery() SensorQuery {
	return &sensorQuery{
		db: d.DB,
	}
}

func NewDB() (*sql.DB, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
