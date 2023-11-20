package repository

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

type SensorQuery interface {
	FetchAverageTransparency(groupName string) (averageTransparency float64, err error)
	FetchAverageTemperature(groupName string) (averageTransparency float64, err error)
	FetchSpeciesList(groupName string) (speciesList map[string]int, err error)
	FetchTopNSpeciesList(groupName string, n int, from, till *time.Time) (speciesList map[string]int, err error)
	GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (minTemperature float64, err error)
	GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (maxTemperature float64, err error)
	FetchCodeNameAverageTemperature(codeName string, from, till time.Time) (averageTransparency float64, err error)
}

type sensorQuery struct {
	db *sql.DB
}

func (s *sensorQuery) FetchAverageTransparency(groupName string) (averageTransparency float64, err error) {
	rows, err := s.db.Query("SELECT ROUND(AVG(sd.transparency), 2) AS transparency FROM sensors s JOIN sensor_data sd ON s.id = sd.sensor_id JOIN sensor_groups sg ON sg.id = s.group_id WHERE sg.name = $1 GROUP BY s.group_id;", groupName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalTransparency float64
	var rowCount int

	for rows.Next() {
		var transparency float64
		if err := rows.Scan(&transparency); err != nil {
			return 0, err
		}
		totalTransparency += transparency
		rowCount++
	}

	if rowCount == 0 {
		return 0, fmt.Errorf("no transparency data found for the group")
	}

	return totalTransparency / float64(rowCount), nil
}

func (s *sensorQuery) FetchAverageTemperature(groupName string) (averageTemperature float64, err error) {
	rows, err := s.db.Query("SELECT AVG(sd.temperature) AS temperature FROM sensors s JOIN sensor_data sd ON s.id = sd.sensor_id JOIN sensor_groups sg ON sg.id = s.group_id WHERE sg.name = $1 GROUP BY s.group_id;", groupName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalTemperature float64
	var rowCount int

	for rows.Next() {
		var temperature float64
		if err := rows.Scan(&temperature); err != nil {
			return 0, err
		}
		totalTemperature += temperature
		rowCount++
	}

	if rowCount == 0 {
		return 0, fmt.Errorf("no transparency data found for the group")
	}

	return roundToPrecision(totalTemperature/float64(rowCount), 2), nil
}

func roundToPrecision(value float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(value*shift) / shift
}

func (s *sensorQuery) FetchSpeciesList(groupName string) (speciesList map[string]int, err error) {
	rows, err := s.db.Query("SELECT fish_species_name, COUNT(*) as count FROM sensor_data sd JOIN sensors s ON sd.sensor_id = s.id JOIN sensor_groups sg ON s.group_id = sg.id	WHERE sg.name = $1 GROUP BY fish_species_name;", groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fishSpeciesName string
	var count int

	speciesList = make(map[string]int)

	for rows.Next() {
		err := rows.Scan(&fishSpeciesName, &count)
		if err != nil {
			return nil, err
		}
		speciesList[fishSpeciesName] = count
	}

	return
}

func (s *sensorQuery) FetchTopNSpeciesList(groupName string, n int, from, till *time.Time) (speciesList map[string]int, err error) {

	query := `
		SELECT fish_species_name AS species, COUNT(*) AS count
		FROM sensor_data sd
		JOIN sensors s ON sd.sensor_id = s.id
		WHERE s.group_id = (SELECT id FROM sensor_groups WHERE name = $1)
		`
	// Add time range conditions if parameters are provided
	if from != nil {
		query += " AND sd.created_at >= $2"
	}

	if till != nil {
		query += " AND sd.created_at <= $3"
	}

	query += fmt.Sprintf(" GROUP BY fish_species_name ORDER BY count DESC LIMIT %d", n)

	// query += `
	// 		 GROUP BY fish_species_name
	// 		ORDER BY count DESC
	// 		LIMIT $4
	// 	`

	// Execute the SQL query
	var rows *sql.Rows

	// Check the number of parameters to bind
	if from == nil && till == nil {
		rows, err = s.db.Query(query, groupName)
	} else if from != nil && till == nil {
		rows, err = s.db.Query(query, groupName, *from)
	} else if from == nil && till != nil {
		rows, err = s.db.Query(query, groupName, *till)
	} else {
		rows, err = s.db.Query(query, groupName, *from, *till)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fishSpeciesName string
	var count int

	speciesList = make(map[string]int)

	for rows.Next() {
		err := rows.Scan(&fishSpeciesName, &count)
		if err != nil {
			return nil, err
		}
		speciesList[fishSpeciesName] = count
	}

	return
}

func (s *sensorQuery) GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (minTemperature float64, err error) {
	err = s.db.QueryRow("SELECT COALESCE(MIN(temperature), 0.0) AS min_temperature	FROM sensor_data sd JOIN sensors s ON sd.sensor_id = s.id WHERE (x >= $1 AND x <= $2) AND (y >= $3 AND y <= $4) AND (z >= $5 AND z <= $6)", xMin, xMax, yMin, yMax, zMin, zMax).Scan(&minTemperature)
	if err != nil {
		return 0, err
	}
	return
}

func (s *sensorQuery) GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (maxTemperature float64, err error) {
	err = s.db.QueryRow("SELECT COALESCE(MAX(temperature), 0.0) AS min_temperature	FROM sensor_data sd JOIN sensors s ON sd.sensor_id = s.id WHERE (x >= $1 AND x <= $2) AND (y >= $3 AND y <= $4) AND (z >= $5 AND z <= $6)", xMin, xMax, yMin, yMax, zMin, zMax).Scan(&maxTemperature)
	if err != nil {
		return 0, err
	}
	return
}

func (s *sensorQuery) FetchCodeNameAverageTemperature(codeName string, from, till time.Time) (averageTemperature float64, err error) {
	rows, err := s.db.Query("SELECT COALESCE(AVG(temperature), 0.0) AS avg_temperature FROM sensor_data WHERE sensor_id = (SELECT id FROM sensors WHERE codename = $1) AND created_at BETWEEN $2 AND $3;", codeName, from, till)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalTemperature float64
	var rowCount int

	for rows.Next() {
		var temperature float64
		if err := rows.Scan(&temperature); err != nil {
			return 0, err
		}
		totalTemperature += temperature
		rowCount++
	}

	if rowCount == 0 {
		return 0, fmt.Errorf("no transparency data found for the group")
	}

	return roundToPrecision(totalTemperature/float64(rowCount), 2), nil
}
