package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/AleF83/airbag/config"

	jwkfetch "github.com/Soluto/fetch-jwk"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	log "github.com/sirupsen/logrus"
)

// NewAuthMiddleware creates new auth middleware
func NewAuthMiddleware(providers []config.JWTProvider, UnauthenticatedRoutes []*regexp.Regexp) func(http.Handler) http.Handler {
	initJWKs(providers)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if matches(UnauthenticatedRoutes, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, jwkfetch.FromIssuerClaim())
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					log.WithError(err).Error("Error while validating request JWT")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if err == request.ErrNoTokenInRequest {
					log.WithError(err).Error("Error while validating request JWT")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.WithError(err).Error("Error while validating request JWT")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !token.Valid {
				log.Error("Error while validating request JWT: token is invalid.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = validate(token, providers)
			if err != nil {
				log.WithError(err).Error("Error while validating request JWT")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})

	}
}

func validate(t *jwt.Token, providers []config.JWTProvider) error {
	claims := t.Claims.(jwt.MapClaims)
	iss := claims["iss"].(string)
	auds := extractAudience(claims)

	for _, p := range providers {
		if iss == p.Issuer && contains(auds, p.Audience) {
			return nil
		}
	}
	log.Printf("Unmatched iss (%s) or aud(%v)\n", iss, auds)
	return fmt.Errorf("iss or aud claims in JWT is invalid")
}

func extractAudience(c jwt.MapClaims) []string {
	switch c["aud"].(type) {
	case []interface{}:
		auds := make([]string, len(c["aud"].([]interface{})))
		for i, value := range c["aud"].([]interface{}) {
			auds[i] = value.(string)
		}
		return auds
	case []string:
		return c["aud"].([]string)
	default:
		return []string{c["aud"].(string)}
	}
}

func initJWKs(providers []config.JWTProvider) {
	jwkProviders := make([]jwkfetch.JWKProvider, len(providers))
	for i, jwtP := range providers {
		jwkProviders[i] = jwkfetch.JWKProvider{
			Issuer: jwtP.Issuer,
			JWKURL: jwtP.JWKURL,
		}
	}
	jwkfetch.Init(jwkProviders)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func matches(regs []*regexp.Regexp, e string) bool {
	for _, r := range regs {
		if r.MatchString(e) {
			return true
		}
	}
	return false
}
