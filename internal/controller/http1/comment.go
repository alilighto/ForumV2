package http1

import (
	"encoding/json"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value("id").(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	input := entity.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = uint(userID)
	if status, err := h.service.Comment.CreateComment(r.Context(), input); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) voteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value("id").(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	var input entity.CommentVote
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = uint(userID)
	if status, err := h.service.Comment.UpsertCommentVote(r.Context(), input); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value("id").(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	strCommentID := strings.TrimPrefix(r.URL.Path, "/api/comment/delete/")
	CommentID, err := strconv.ParseUint(strCommentID, 10, 64)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, fmt.Sprintf("Invalid id: %v", CommentID))
		return
	}

	if status, err := h.service.Comment.DeleteComment(r.Context(), uint(CommentID), uint(userID)); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
