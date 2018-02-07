package app

import (
	"os"
	"fmt"
	"net/http"
	"regexp"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"github.com/eriklindqvist/recepies_auth/log"
	l "github.com/eriklindqvist/recepies_api/app/lib"
	c "github.com/eriklindqvist/recepies_api/app/controllers"
	jwt "github.com/dgrijalva/jwt-go"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route
type Endpoint func(w http.ResponseWriter, r *http.Request) ([]byte, error)

type Scope struct {
  Entity string `json:"ent"`
  Actions []string `json:"act"`
}

type User struct {
    Username string `json:"usr"`
    Scopes []Scope `json:"scp"`
    jwt.StandardClaims
}

func getSession() *mgo.Session {
		host := "mongodb://" + l.Getenv("MONGODB_HOST", "localhost")
    s, err := mgo.Dial(host)
		log.Info(fmt.Sprintf("host: %s", host))
    // Check if connection error, is mongo running?
    if err != nil {
        log.Panic(err.Error())
    }
    return s
}

func getPublicKey() *rsa.PublicKey {
	var (
		pub *rsa.PublicKey
		pem []byte
    err error
	)

	if pem, err = ioutil.ReadFile(l.Getenv("KEYFILE", "public.rsa")); err != nil {
    log.Panic(err.Error())
  }

	if pub, err = jwt.ParseRSAPublicKeyFromPEM(pem); err != nil {
		log.Panic(err.Error())
	}

	return pub
}

var rc = c.NewRecipeController(getSession())
var routes = NewRoutes(*rc)
var publicKey = getPublicKey()

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler = Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func validates(scopes []Scope, entity string, action string) bool {
	for _, scope:= range scopes {
		if (scope.Entity == entity) {
			for _, act := range scope.Actions {
				if act == action {
					return true
				}
			}
			return false
		}
	}
	return false
}

func contains(arr []string, str string) bool {
  for _, a := range arr {
      if a == str {
          return true
      }
  }
  return false
}

func Handle(entity string, action string, endpoint Endpoint, w http.ResponseWriter, r *http.Request) {
	protected := strings.Split(os.Getenv("PROTECTED_ENDPOINTS"),",")

	if contains(protected, entity+":"+action) {
		authorization := r.Header.Get("Authorization")
		regex, _ := regexp.Compile("(?:Bearer *)([^ ]+)(?: *)")
		matches := regex.FindStringSubmatch(authorization)

		if len(matches) != 2 {
			log.Err(fmt.Sprintf("Malformed Authorization header: %s", authorization))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		 	return
		}

		jwtToken := matches[1]

		// parse token
		token, err := jwt.ParseWithClaims(jwtToken, &User{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			log.Err(fmt.Sprintf("Could not parse token: %s", err.Error()))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// extract claims
		user, ok := token.Claims.(*User)

		if !ok || !token.Valid {
			log.Err(fmt.Sprintf("Invalid token: %s", token))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if (!validates(user.Scopes, entity, action)) {
			log.Err(fmt.Sprintf("Insufficient privileges: %s", user.Scopes))
			http.Error(w, "Unauthorized: Insufficient privileges", http.StatusUnauthorized)
			return
		}
	}

	setContentType(w)

	body, err := endpoint(w, r)

	if err != nil {
		switch e := err.(type) {
		case l.Error:
			http.Error(w, e.Error(), e.Status())
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	writeBody(w, body)
}

func setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func writeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func writeBody(w http.ResponseWriter, body []byte) {
	fmt.Fprintf(w, "%s", body)
}
