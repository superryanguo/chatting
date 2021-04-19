package psession

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	log "github.com/micro/go-micro/v2/logger"
)

const (
	CtsKey = "ctskey"
)

var (
	sessName = "csess"
	store    *sessions.CookieStore
)

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	//TODO: should we hardcode one key if not export this var?
}

func GetSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	ses, err := store.Get(r, sessName)
	if ses != nil && ses.Values[CtsKey] == nil {
		sesId := uuid.New().String()
		ses.Values[CtsKey] = sesId
		log.Debug("GetCtSession->New Session generate id=", sesId)
		ses.Save(r, w)
	}

	return ses, err
}
