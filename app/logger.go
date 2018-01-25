package app

import (
    "github.com/eriklindqvist/recepies_auth/log"
    "net/http"
    "fmt"
)

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Info(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, name))
        inner.ServeHTTP(w, r)
    })
}
