package main

import (
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	creatorClient pb.CreatorsClient,
) {
	mux.HandleFunc("/activity", activity(creatorClient))
	mux.Handle("/", fsHandler())
}
