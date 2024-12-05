package http1

import (
	"context"
	"forum/internal/entity"
	smpljwt "forum/pkg/smplJwt"
	"net/http"
	"strings"
)

func (h *Handler) corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http//localhost:8081/")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) identify(role uint, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if role > entity.Roles.Guest {
			header, ok := r.Header["Authorization"]
			if !ok {
				h.errorHandler(w, r, http.StatusUnauthorized, "empty auth header")
				return
			}

			headerParts := strings.Split(header[0], " ")
			if len(headerParts) != 2 {
				h.errorHandler(w, r, http.StatusUnauthorized, "invalid auth header")
				return
			}

			exist, err := h.service.IsTokenExist(r.Context(), headerParts[1])
			if err != nil {
				h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			if !exist {
				h.errorHandler(w, r, http.StatusUnauthorized, "invalid token")
				return
			}
			id, err := smpljwt.ParseToken(headerParts[1], h.secret)
			if err != nil {
				if err == smpljwt.ErrExpiredToken {
					if dberr := h.service.DeleteSessionByToken(r.Context(), headerParts[1]); dberr != nil {
						h.errorHandler(w, r, http.StatusInternalServerError, dberr.Error())
						return
					}
				}
				h.errorHandler(w, r, http.StatusUnauthorized, "invalid token")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "id", id))
			r = r.WithContext(context.WithValue(r.Context(), "token", headerParts[1]))
			next(w, r)
			return
		}

		next(w, r)
	}
}

func (h *Handler) isAlreadyIdentified(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header["Authorization"]
		if ok {
			h.errorHandler(w, r, http.StatusForbidden, "already authorized")
			return
		}
		next(w, r)
	}
}
