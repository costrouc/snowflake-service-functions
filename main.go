package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
)

// https://docs.snowflake.com/en/sql-reference/external-functions-data-format#header-format
type SnowflakeFuctionHeader struct {
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

type SnowflakeFunctionResponseHeader struct {
	ContentMD5 string `json:"Content-MD5"`
}

type SnowflakeFunctionResponse struct {
	StatusCode ResponseStatusCode              `json:"statusCode"`
	Body       string                          `json:"body"`
	Headers    SnowflakeFunctionResponseHeader `json:"headers"`
}

type ResponseStatusCode int64

// Status Codes
// https://docs.snowflake.com/en/sql-reference/external-functions-data-format#status-code
// 200 - Batch processed successfully.
// 202 - Batch received and still being processed.
const (
	BatchProcessed  ResponseStatusCode = 200
	BatchProcessing ResponseStatusCode = 202
)

func handleSnowflakeFunctionGET(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "path", r.URL.Path)

	responseData := SnowflakeFunctionResponse{
		StatusCode: BatchProcessed,
	}

	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		slog.Info("Error processing request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleSnowflakeFunctionPOST(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "path", r.URL.Path)
	w.Write([]byte("post asdf"))
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "address for server to run")
	flag.Parse()

	handler := http.NewServeMux()
	handler.HandleFunc("GET /rpc", handleSnowflakeFunctionGET)
	handler.HandleFunc("POST /rpc", handleSnowflakeFunctionPOST)

	slog.Info("Server listening", "addr", *addr)
	http.ListenAndServe(*addr, handler)
}
