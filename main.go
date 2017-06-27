package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"github.com/gorilla/mux"
)

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

func init() {
	// configure logger
	log.SetLevel(log.DebugLevel)

	ln := net.Listen()

	// Amazon API Gateway binary media types are supported out of the box.
	// If you don't send or receive binary data, you can safely set it to nil.
	Handle = apigatewayproxy.New(ln, nil).Handle

	// Gorilla mux
	r := mux.NewRouter()
	// setup routes
	r.HandleFunc("/health", handleHome).Methods("GET")

	go http.Serve(ln, r)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	status := &Status{
		State:       "OK",
		CurrentTime: time.Now(),
	}

	log.WithField("State", status.State).WithField("CurrentTime", status.CurrentTime).Debug("Response Type")

	j, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write(j)
}

// Main is required for completeness is not used by the shim. Can be omitted but linters nolikey
func main() {}

// Status simple return type
type Status struct {
	State       string    `json:"state"`
	CurrentTime time.Time `json:"currentTime"`
}
