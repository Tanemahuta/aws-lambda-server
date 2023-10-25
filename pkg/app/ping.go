package app

import (
	"fmt"
	"net/http"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "ok")
}
