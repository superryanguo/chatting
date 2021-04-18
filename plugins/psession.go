package psession

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	CesKey = "csekey"
)

var (
	sessName = "csess"
	store   *sessions.CookieStore
)

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	//TODO: should we hardcode one key if not export this var?
}

func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session, error {
	return store.Get(r, sessName)
}
