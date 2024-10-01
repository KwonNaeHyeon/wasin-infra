package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	JOBNAME_ROUTERS = "routers" // 모니터링 라우터 관리용 Job Name
)

type APIHandler interface {
	AddTarget(rw http.ResponseWriter, req *http.Request)
	RemoveTarget(rw http.ResponseWriter, req *http.Request)
	ListTargets(rw http.ResponseWriter, req *http.Request)
}

type apiHandler struct {
	config Config
}

func NewAPIHander(config Config) APIHandler {
	return &apiHandler{
		config: config,
	}
}

func (ah *apiHandler) AddTarget(rw http.ResponseWriter, req *http.Request) {
	slog.Info("API Received: Add Target.")

	var body AddTargetRequest
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = ah.config.AddTarget(JOBNAME_ROUTERS, body.IP)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ah.config.Save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AddTargetResponse{Msg: "Successfully Added."}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (ah *apiHandler) RemoveTarget(rw http.ResponseWriter, req *http.Request) {
	slog.Info("API Received: Remove Target.")

	var body RemoveTargetRequest
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = ah.config.RemoveTarget(JOBNAME_ROUTERS, body.IP)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ah.config.Save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AddTargetResponse{Msg: "Successfully Removed."}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (ah *apiHandler) ListTargets(rw http.ResponseWriter, req *http.Request) {
	slog.Info("API Received: List Targets.")

	var body ListTargetRequest
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	targets, err := ah.config.ListTargets(JOBNAME_ROUTERS)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListTargetsResponse{
		Msg:     "Successfully fetched the targets.",
		Targets: targets,
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}
