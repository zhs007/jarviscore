package jarviscore

import (
	"net/http"
	"runtime"
	"strconv"
)

// InitPprof - init pprof
func InitPprof(cfg *Config) error {

	if cfg.Pprof.GoRoutineURL != "" {
		mux := http.NewServeMux()

		mux.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		})

		http.ListenAndServe(cfg.Pprof.GoRoutineURL, mux)
	}

	if cfg.Pprof.BaseURL != "" {
		http.ListenAndServe(cfg.Pprof.BaseURL, nil)
	}

	return nil
}
