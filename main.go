package main

import (
	"flag"
	"log/slog"
	"net/http"
)

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

type SnowflakeFunctionResponseHeader struct {
	ContentMD5 string `json:"Content-MD5"`
}

func handleSnowflakeFunctionGET(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "path", r.URL.Path)

	headers := readHeaders(r)

	switch headers.SfExternalFunctionName {
	case "MYFUNCTION":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"data\": [[0, \"myfunction\"]]}"))
		return
	default:
		slog.Info("Error processing request unknown function name", "functionName", headers.SfExternalFunctionName)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Service Error"))
		return
	}
}

type FunctionMYFUNCTIONOptions struct {
	Arg1 string
}

func handleSnowflakeFunctionPOST(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "path", r.URL.Path)

	headers := readHeaders(r)

	switch headers.SfExternalFunctionName {
	case "MYFUNCTION":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"data\": [[0, \"myfunction\"]]}"))
		return
	default:
		slog.Info("Error processing request unknown function name", "functionName", headers.SfExternalFunctionName)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Service Error"))
		return
	}
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "address for server to run")
	flag.Parse()

	handler := http.NewServeMux()
	handler.HandleFunc("GET /rpc", handleSnowflakeFunctionGET)
	handler.HandleFunc("POST /rpc", handleSnowflakeFunctionPOST)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
	})

	slog.Info("Server listening", "addr", *addr)
	http.ListenAndServe(*addr, handler)
}
