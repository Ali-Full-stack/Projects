package middleware

import (
	"bookstore/auth"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
)

func RolePasswordMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("role")
		password := r.Header.Get("password")

		if role == "" {
			http.Error(w, "Missing role", http.StatusUnauthorized)
			return
		}
		role = strings.ToLower(role)

		if role == "admin" {
			if password != os.Getenv("admin_password") {
				http.Error(w, "Admin Password is incorrect !", http.StatusInternalServerError)
				return
			}
		}

		enforcer, err := casbin.NewEnforcer("./casbin/auth.conf", "./casbin/auth.csv")
		if err != nil {
			log.Println("error: casbin connection:", err)
			http.Error(w, "Something went wrong !", http.StatusInternalServerError)
			return
		}

		ok, _ := enforcer.Enforce(role, r.URL.Path, r.Method)
		if !ok {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckJwtTokenMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		token :=r.Header.Get("token")
		if err :=auth.VerifyToken(token); err != nil {
			http.Error(w, "token is expired !", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
