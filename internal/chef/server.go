package chef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hidalgopl/secureapi-boatswain/internal/http/headers"
	"github.com/hidalgopl/secureapi-boatswain/internal/http/mediatypes"
)

type Server struct {
	starter starter
}

func NewServer() (*Server, error) {
	starter := newStarter()

	server := &Server{
		starter: starter,
	}
	return server, nil
}

func (s *Server) AttachRoutes(router *mux.Router) {
	router.HandleFunc("/tests/schedule", s.start()).
		Methods("POST")
}

func (s *Server) start() http.HandlerFunc {
	// auth between Boatswain and Web happens via Auth header.
	// Secret key is shared between those two services via K8s secret.
	decodeReq := func(r *http.Request) (scheduleTestsReq, error) {
		var req scheduleTestsReq
		err := json.NewDecoder(r.Body).Decode(&req)
		return req, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := checkAuth(r)
		if err != nil {
			encodeErrRsp(w, http.StatusUnauthorized, "auth missing")
			return
		}
		req, err := decodeReq(r)
		if err != nil {
			encodeErrRsp(w, http.StatusUnprocessableEntity, "cannot parse start chef req body")
			return
		}
		rsp, err := s.starter.start(r.Context(), req)
		if err != nil {
			encodeErrRsp(w, http.StatusInternalServerError, "cannot start chef")
			return
		}
		w.Header().Set(headers.ContentType, mediatypes.ApplicationJSONUtf8)
		w.WriteHeader(http.StatusOK)
		body := <- rsp
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(body)
		callbackUrl := fmt.Sprintf("%s/%s", "secureapi-web.default.local.cluster./sec-test", req.TestSuiteId)
		go http.Post(
			callbackUrl,
			mediatypes.ApplicationJSONUtf8,
			b,
		)

	}
}
func checkAuth(request *http.Request) error {
	authKey := request.Header.Get("X-SecureAPI-Secret-Key")
	if authKey != os.Getenv("X_SECUREAPI_SECRET_KEY") {
		return errors.New("missing auth")
	}
	return nil
}

func encodeErrRsp(w http.ResponseWriter, statusCode int, msg string) {
	errBody := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}
	encodeRsp(w, statusCode, errBody)
}

func encodeRsp(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set(headers.ContentType, mediatypes.ApplicationJSONUtf8)
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("cannot encode body")
	}
}
