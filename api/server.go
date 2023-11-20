package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sensors/internal/app"
)

type Server struct {
	microserviceServer app.MicroserviceServer
}

func NewServer(ctx context.Context, microserviceServer app.MicroserviceServer) *Server {
	return &Server{
		microserviceServer: microserviceServer,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/group/{groupName}/transparency/average", s.getGroupTransparencyAverage)
	router.HandleFunc("/group/{groupName}/temperature/average", s.getGroupTemperatureAverage)
	router.HandleFunc("/group/{groupName}/species", s.getGroupSpecies)
	router.HandleFunc("/group/{groupName}/species/top/{n}", s.getTopNGroupSpecies)
	router.HandleFunc("/region/temperature/min", s.getRegionMinTemperature)
	router.HandleFunc("/region/temperature/max", s.getRegionMaxTemperature)
	router.HandleFunc("/sensor/{codeName}/temperature/average", s.getCodenameTemperatureAverage)
	http.Handle("/", router)
	return http.ListenAndServe(":8080", router)
}

func (s *Server) getGroupTransparencyAverage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	averageTransparency, err := s.microserviceServer.GetGroupTransparencyAverage(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error calculating transparency average", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"group": groupName, "averageTransparency": averageTransparency})
}

func (s *Server) getGroupTemperatureAverage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	averageTemperature, err := s.microserviceServer.GetGroupTemperatureAverage(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error calculating temperature average", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"group": groupName, "averageTemperature": averageTemperature})
}

func (s *Server) getGroupSpecies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	speciesList, err := s.microserviceServer.GetGroupSpecies(r.Context(), groupName)
	if err != nil {
		http.Error(w, "Error fetching species list", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"group": groupName, "speciesList": speciesList})
}

func (s *Server) getTopNGroupSpecies(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupName := vars["groupName"]
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		http.Error(w, "Error fetching top n species list, invalid n value", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	fromStr := query.Get("from")
	tillStr := query.Get("till")

	var fromTime, tillTime *time.Time

	from, err := strconv.Atoi(fromStr)
	if err != nil && fromStr != "" {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	till, err := strconv.Atoi(tillStr)
	if err != nil && tillStr != "" {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	from_time := time.Unix(int64(from), 0)

	till_time := time.Unix(int64(till), 0)

	if fromStr != "" {
		fromTime = &from_time
	}

	if tillStr != "" {
		tillTime = &till_time
	}

	speciesList, err := s.microserviceServer.GetTopNGroupSpecies(r.Context(), groupName, n, fromTime, tillTime)
	if err != nil {
		http.Error(w, "Error fetching top n species list", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"group": groupName, "speciesList": speciesList})
}

func (s *Server) getRegionMinTemperature(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	xMin, err := strconv.ParseFloat(params.Get("xMin"), 64)
	if err != nil {
		http.Error(w, "Invalid xMin parameter", http.StatusBadRequest)
		return
	}

	xMax, err := strconv.ParseFloat(params.Get("xMax"), 64)
	if err != nil {
		http.Error(w, "Invalid xMax parameter", http.StatusBadRequest)
		return
	}

	yMin, err := strconv.ParseFloat(params.Get("yMin"), 64)
	if err != nil {
		http.Error(w, "Invalid yMin parameter", http.StatusBadRequest)
		return
	}

	yMax, err := strconv.ParseFloat(params.Get("yMax"), 64)
	if err != nil {
		http.Error(w, "Invalid yMax parameter", http.StatusBadRequest)
		return
	}

	zMin, err := strconv.ParseFloat(params.Get("zMin"), 64)
	if err != nil {
		http.Error(w, "Invalid zMin parameter", http.StatusBadRequest)
		return
	}

	zMax, err := strconv.ParseFloat(params.Get("zMax"), 64)
	if err != nil {
		http.Error(w, "Invalid zMax parameter", http.StatusBadRequest)
		return
	}

	minTemperature, err := s.microserviceServer.GetRegionMinTemperature(r.Context(), xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		http.Error(w, "Error fetching region min temperature", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"Min Temperature": minTemperature})
}

func (s *Server) getRegionMaxTemperature(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// Parse query parameters
	xMin, err := strconv.ParseFloat(params.Get("xMin"), 64)
	if err != nil {
		http.Error(w, "Invalid xMin parameter", http.StatusBadRequest)
		return
	}

	xMax, err := strconv.ParseFloat(params.Get("xMax"), 64)
	if err != nil {
		http.Error(w, "Invalid xMax parameter", http.StatusBadRequest)
		return
	}

	yMin, err := strconv.ParseFloat(params.Get("yMin"), 64)
	if err != nil {
		http.Error(w, "Invalid yMin parameter", http.StatusBadRequest)
		return
	}

	yMax, err := strconv.ParseFloat(params.Get("yMax"), 64)
	if err != nil {
		http.Error(w, "Invalid yMax parameter", http.StatusBadRequest)
		return
	}

	zMin, err := strconv.ParseFloat(params.Get("zMin"), 64)
	if err != nil {
		http.Error(w, "Invalid zMin parameter", http.StatusBadRequest)
		return
	}

	zMax, err := strconv.ParseFloat(params.Get("zMax"), 64)
	if err != nil {
		http.Error(w, "Invalid zMax parameter", http.StatusBadRequest)
		return
	}

	maxTemperature, err := s.microserviceServer.GetRegionMaxTemperature(r.Context(), xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		http.Error(w, "Error fetching region max temperature", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"Max Temperature": maxTemperature})
}

func (s *Server) getCodenameTemperatureAverage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codeName := vars["codeName"]

	query := r.URL.Query()
	fromStr := query.Get("from")
	tillStr := query.Get("till")

	from, err := strconv.Atoi(fromStr)
	if err != nil && fromStr != "" {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	till, err := strconv.Atoi(tillStr)
	if err != nil && tillStr != "" {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	fromTime := time.Unix(int64(from), 0)

	tillTime := time.Unix(int64(till), 0)

	averageTemperature, err := s.microserviceServer.GetCodeNameTemperatureAverage(r.Context(), codeName, fromTime, tillTime)
	if err != nil {
		http.Error(w, "Error calculating temperature average", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"codeName": codeName, "averageTemperature": averageTemperature})
}
