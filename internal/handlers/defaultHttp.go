package handlers

import "net/http"

func (h *Handler) RunServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, h.setMetrics)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}

	return nil
}
