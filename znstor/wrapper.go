package znstor

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

func Wrapper(inner http.Handler, name string, logfile io.Writer, login, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// Check is a authorized request

		userLogin, userPassowrd, _ := r.BasicAuth()
		if userLogin != login || userPassowrd != password {
			sendMessage(w, http.StatusUnauthorized, traceFunctionName(), "")
			return
		}

		// Check is request assign to private pool?
		vars := mux.Vars(r)
		pool := vars["pool"]

		// TODO: declare privat pool list in configuration file
		if pool == "rpool" || pool == "zroot" {
			w.WriteHeader(http.StatusForbidden)
			sendMessage(w, http.StatusForbidden, traceFunctionName(), "PERMISSION DENIED on POOL: "+pool)
			return
		}

		// log requests
		log.SetOutput(logfile)
		log.Println("---")
		log.Printf(
			"%s\t%s\t%s\t",
			name,
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		// Processing
		inner.ServeHTTP(w, r)

		log.Printf("%s\n", time.Since(start))
	})
}
