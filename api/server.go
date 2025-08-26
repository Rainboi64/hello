package api

import (
	"fmt"
	"net/http"
)

func ShortnerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}