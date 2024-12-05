package http1

import (
	"encoding/json"
	"fmt"
	"forum/internal/entity"
	smpljwt "forum/pkg/smplJwt"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	var input entity.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, fmt.Sprintf("invalid json input: %v", err.Error()))
		return
	}
	status, err := h.service.User.Create(r.Context(), input)
	if err != nil {
		h.errorHandler(w, r, status, fmt.Sprintf("bad request: %v", err.Error()))
		return
	}
	w.WriteHeader(status)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	var input entity.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, fmt.Sprintf("invalid json input: %v", err.Error()))
		return
	}
	token, status, err := h.service.SignIn(r.Context(), input)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	}); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strUserID := r.URL.Path[len("/api/profile/"):]
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}
	if userID < 0 {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}
	user, status, err := h.service.User.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) isValidToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	header, ok := r.Header["Authorization"]
	if !ok {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"checker": false,
		}); err != nil {
			h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	headerParts := strings.Split(header[0], " ")
	if len(headerParts) != 2 {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"checker": false,
		}); err != nil {
			h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

	exist, err := h.service.IsTokenExist(r.Context(), headerParts[1])
	if err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if !exist {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"checker": false,
		}); err != nil {
			h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}
	_, err = smpljwt.ParseToken(headerParts[1], h.secret)
	if err != nil {
		if err == smpljwt.ErrExpiredToken {
			if dberr := h.service.DeleteSessionByToken(r.Context(), headerParts[1]); dberr != nil {
				h.errorHandler(w, r, http.StatusInternalServerError, dberr.Error())
				return
			}
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"checker": false,
		}); err != nil {
			h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"checker": true,
	}); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) signOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	token := r.Context().Value("token").(string)
	if err := h.service.DeleteSessionByToken(r.Context(), token); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
