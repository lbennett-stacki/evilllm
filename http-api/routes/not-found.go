package routes

import (
	"fmt"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not found")
}
