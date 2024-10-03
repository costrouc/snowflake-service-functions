package main

import (
	"fmt"
	"flag"
	"log/slog"
	"net/http"
	"os/exec"
)

type EchoServiceFunction struct {
	headers *SnowflakeFunctionHeader
}

func (s *EchoServiceFunction) Result() (interface{}, error) {
	return s.headers, nil
}

type ShellServiceFunction struct {
	command string
	args []string
}

func (s *ShellServiceFunction) Result() (interface{}, error) {
	out, err := exec.Command(s.command, s.args...).Output()
	if err != nil {
		return nil, fmt.Errorf("running shell command=%s args=%s, %w", s.command, s.args, err)
	}

	return string(out), nil
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "address for server to run")
	flag.Parse()

	serviceMux := NewSnowflakeServiceMux()
	serviceMux.ServiceFunc("DEBUG", "(TEXT VARCHAR)", func(headers *SnowflakeFunctionHeader, args []interface{}) Task {
		return &EchoServiceFunction{
			headers: headers,
		}
	})

	serviceMux.ServiceFunc("SHELL", "(CMD VARCHAR)", func(headers *SnowflakeFunctionHeader, args []interface{}) Task {
		return &ShellServiceFunction{
			command: "date",
			args: []string{},
		}
	})

	handler := http.NewServeMux()
	handler.Handle("/rpc", serviceMux)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
	})

	slog.Info("Server listening", "addr", *addr)
	http.ListenAndServe(*addr, handler)
}
