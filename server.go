package gameoflife

import (
  "context"
  "encoding/json"
  "fmt"

  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/gorilla/mux"
)

type Server struct {
  Changes chan Change
  Syncs   chan Sync
  Inits   chan Init
}

func (s *Server) Run(port string) {
  var wait time.Duration

  r := mux.NewRouter()
  api := r.PathPrefix("/api").Subrouter()

  api.HandleFunc("/init", s.initHandle).Methods(http.MethodPost)

  api.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) { s.syncHandle(w, r) })

  api.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) { s.updateHandle(w, r) }).Methods(http.MethodPost)

  srv := &http.Server{
    Addr:         ":" + port,
    WriteTimeout: time.Second * 15,
    ReadTimeout:  time.Second * 15,
    IdleTimeout:  time.Second * 60,
    Handler:      r,
  }

  go func() {
    if err := srv.ListenAndServe(); err != nil {
      fmt.Println(err)
    }
  }()

  end := make(chan os.Signal, 1)
  signal.Notify(end, os.Interrupt, syscall.SIGTERM)
  <-end

  ctx, cancel := context.WithTimeout(context.Background(), wait)
  defer cancel()

  srv.Shutdown(ctx)
  os.Exit(0)
}

func (s *Server) syncHandle(w http.ResponseWriter, r *http.Request) {
  var resp Sync

  err := json.NewDecoder(r.Body).Decode(&resp)
  if err != nil {
    HttpErrWrite("Error while reading body", err, http.StatusUnprocessableEntity, w)
    return
  }

  s.Syncs <- resp

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
}

func HttpErrWrite(msg string, err error, status int, w http.ResponseWriter) {
  log.Println(msg, err)

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(status)
  if err := json.NewEncoder(w).Encode(&HttpErrorBody{Err: msg + err.Error()}); err != nil {
    fmt.Println(err)
  }
  return
}

func (s *Server) updateHandle(w http.ResponseWriter, r *http.Request) {
  var resp []Change

  err := json.NewDecoder(r.Body).Decode(&resp)
  if err != nil {
    HttpErrWrite("Error while reading body", err, http.StatusUnprocessableEntity, w)
    return
  }

  for _, chg := range resp {
    s.Changes <- chg
  }

  w.Header().Set("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
}

func (s *Server) initHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("content-type", "application/json")
  var resp Init

  err := json.NewDecoder(r.Body).Decode(&resp)
  if err != nil {
    HttpErrWrite("Error while reading body", err, http.StatusUnprocessableEntity, w)
    return
  }

  w.WriteHeader(http.StatusOK)
  s.Inits <- resp
}

type HttpErrorBody struct {
  Err string `json:"error"`
}
