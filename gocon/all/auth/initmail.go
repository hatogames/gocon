package auth

import (
	"fmt"
	"net/http"
)

func Initmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "initmail")
}
