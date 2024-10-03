package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log/slog"
)

type Task interface {
	Result() (interface{}, error)
}

type SnowflakeServiceMux struct {
	functions map[string]func(headers *SnowflakeFunctionHeader, args []interface{}) Task
	tasks map[string]Task
}

func NewSnowflakeServiceMux() *SnowflakeServiceMux {
	return &SnowflakeServiceMux{
		functions: make(map[string]func(headers *SnowflakeFunctionHeader, args []interface{}) Task, 0),
		tasks: make(map[string]Task, 0),
	}
}

func (s *SnowflakeServiceMux) ServiceFunc(functionName string, functionSignature string, f func(headers *SnowflakeFunctionHeader, args []interface{}) Task) {
	s.functions[functionName + functionSignature] = f
}

func (s *SnowflakeServiceMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "path", r.URL.Path)

	headers := readHeaders(r)

	var result interface{}
	var err error

	switch r.Method {
	case "GET":
		// GET Requests are checking if the given QueryId is completed and potentially returning the result
		task, ok := s.tasks[headers.SfExternalFunctionCurrentQueryId]
		if !ok {
			slog.Info("Current Query Id Not Found", "QueryId", headers.SfExternalFunctionCurrentQueryId)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error Current Query Id=%s Not Found", headers.SfExternalFunctionCurrentQueryId)))
			return
		}
		
		result, err = task.Result()
	case "POST":
		// POST Requests are starting the task and potentially returning the result
		var request SnowflakeFunctionRequest

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			slog.Info("Error decoding json from service function", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Service Error"))
			return
		}

		if len(request.Data) != 1 {
			slog.Info("Error expected request to have only one row of data", "row", len(request.Data))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Service Error"))
			return
		}

		functionKey := headers.SfExternalFunctionName + headers.SfExternalFunctionSignature
		f, ok := s.functions[functionKey]
		if !ok {
			slog.Info("Service Function Not Found", "FunctionName", headers.SfExternalFunctionName, "FunctionSignature", headers.SfExternalFunctionSignature)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 Service Function %s Not Found", functionKey)))
			return
		}

		// Create Task
		task := f(headers, request.Data[0])
		result, err = task.Result()
		s.tasks[headers.SfExternalFunctionCurrentQueryId] = task
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	// Task is still processing
	if result == nil && err == nil {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("Task Processing %s", headers.SfExternalFunctionCurrentQueryId)))
		return
	}

	// Task either encountered an error or completed so delete the task
	delete(s.tasks, headers.SfExternalFunctionCurrentQueryId)

	// Task encountered an error
	if err != nil {
		slog.Info("Error encountered running service function", "QueryId", headers.SfExternalFunctionCurrentQueryId, "FunctionName", headers.SfExternalFunctionName, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error running service function %s", err.Error())))
		return
	}

	response := SnowflakeFunctionResponse{
		Data: [][]interface{}{{
			0,
			result,
		}},
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Info("Error processing request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Service Error"))
		return
	}
}

// https://docs.snowflake.com/en/sql-reference/external-functions-data-format#header-format
type SnowflakeFunctionHeader struct {
	SfExternalFunctionFormat           string //always "json"
	SfExternalFunctionFormatVersion    string // always "1.0"
	SfExternalFunctionCurrentQueryId   string
	SfExternalFunctionQueryBatchId     string
	SfExternalFunctionName             string
	SfExternalFunctionNameBase64       string
	SfExternalFunctionSignature        string
	SfExternalFunctionSignatureBase64  string
	SfExternalFunctionReturnType       string
	SfExternalFunctionReturnTypeBase64 string
}

func readHeaders(r *http.Request) *SnowflakeFunctionHeader {
	return &SnowflakeFunctionHeader{
		SfExternalFunctionFormat:           r.Header.Get("sf-external-function-format"),
		SfExternalFunctionFormatVersion:    r.Header.Get("sf-external-function-format-version"),
		SfExternalFunctionCurrentQueryId:   r.Header.Get("sf-external-function-current-query-id"),
		SfExternalFunctionQueryBatchId:     r.Header.Get("sf-external-function-query-batch-id"),
		SfExternalFunctionName:             r.Header.Get("sf-external-function-name"),
		SfExternalFunctionNameBase64:       r.Header.Get("sf-external-function-name-base64"),
		SfExternalFunctionSignature:        r.Header.Get("sf-external-function-signature"),
		SfExternalFunctionSignatureBase64:  r.Header.Get("sf-external-function-signature-base64"),
		SfExternalFunctionReturnType:       r.Header.Get("sf-external-function-return-type"),
		SfExternalFunctionReturnTypeBase64: r.Header.Get("sf-external-function-return-type-base64"),
	}
}

type SnowflakeFunctionRequest struct {
	Data [][]interface{} `json:"data"`
}

type SnowflakeFunctionResponse SnowflakeFunctionRequest
