package routes

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Health handler called")
	fmt.Fprintf(w, "ok")
}
