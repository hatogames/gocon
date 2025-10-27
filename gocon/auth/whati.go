package auth_all

import (
	"fmt"
	"gocon/db/mini"
	"net/http"
)

func Whati(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Fprint(w, "none")
		return
	}

	session, err := mini.GetSession(cookie.Value)
	if err != nil {
		fmt.Fprint(w, "none")
		return
	}

	fmt.Fprint(w, session.Typ)
	fmt.Fprint(w, session.Expires)

}
