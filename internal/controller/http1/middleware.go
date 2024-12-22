package http1

import (
	"context"
	"net/http"

	"forum/internal/entity"
	smpljwt "forum/pkg/smplJwt"
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
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					h.errorHandler(w, r, http.StatusUnauthorized, "no token cookie")
					return
				}
				h.errorHandler(w, r, http.StatusInternalServerError, "failed to get cookie")
				return
			}

			exist, err := h.service.IsTokenExist(r.Context(), cookie.Value)
			if err != nil {
				h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			if !exist {
				h.errorHandler(w, r, http.StatusUnauthorized, "invalid token")
				return
			}

			id, err := smpljwt.ParseToken(cookie.Value, h.secret)
			if err != nil {
				if err == smpljwt.ErrExpiredToken {
					if dberr := h.service.DeleteSessionByToken(r.Context(), cookie.Value); dberr != nil {
						h.errorHandler(w, r, http.StatusInternalServerError, dberr.Error())
						return
					}
				}
				h.errorHandler(w, r, http.StatusUnauthorized, "invalid token")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "id", id))
			r = r.WithContext(context.WithValue(r.Context(), "token", cookie.Value))
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			next(w, r)
			return
		}

		id, err := smpljwt.ParseToken(cookie.Value, h.secret)
		if err != nil {
			next(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "id", id))
		r = r.WithContext(context.WithValue(r.Context(), "token", cookie.Value))

		next(w, r)
	}
}

func (h *Handler) isAlreadyIdentified(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("Authorization")
		if err == nil {
			h.errorHandler(w, r, http.StatusForbidden, "already authorized")
			return
		}
		next(w, r)
	}
}
