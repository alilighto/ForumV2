package http1

import (
	"encoding/json"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) getALLPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	tag := r.URL.Path[len("/api/posts/"):]
	posts, status, err := h.service.Post.GetAllByTag(r.Context(), tag)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) getPostbyID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strPostID := r.URL.Path[len("/api/post/"):]
	id, err := strconv.Atoi(strPostID)
	if err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, "invalid post id")
		return
	}
	post, status, err := h.service.Post.GetPostByID(r.Context(), uint(id))
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(post); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	var input entity.Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	id := r.Context().Value("id").(int)
	if id < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	input.UserID = uint(id)
	postID, status, err := h.service.Post.CreatePost(r.Context(), input)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"post_id": postID,
	}); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strPostID := strings.TrimPrefix(r.URL.Path, "/api/post/delete/")
	postID, err := strconv.ParseUint(strPostID, 10, 64)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, fmt.Sprintf("Invalid id: %v", postID))
		return
	}
	userID := r.Context().Value("id").(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	h.service.Post.DeletePostByID(r.Context(), uint(postID), uint(userID))
}

func (h *Handler) votePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value("id").(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	var input entity.PostVote
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = uint(userID)
	if status, err := h.service.Post.UpsertPostVote(r.Context(), input); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllPostsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strUserID := r.URL.Path[len("/api/profile/posts/"):]
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}
	if userID < 0 {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}

	posts, status, err := h.service.Post.GetAllByUserID(r.Context(), uint(userID))
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) getAllLikedPostsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strUserID := r.URL.Path[len("/api/profile/liked-posts/"):]
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}
	if userID < 0 {
		h.errorHandler(w, r, http.StatusNotFound, "not found")
		return
	}
	posts, status, err := h.service.Post.GetAllLikedPostsByUserID(r.Context(), uint(userID), true)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) getAllDisLikedPostsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	strUserID := r.URL.Path[len("/api/profile/disliked-posts/"):]
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		h.errorHandler(w, r, http.StatusNotFound, err.Error())
		return
	}
	if userID < 0 {
		h.errorHandler(w, r, http.StatusNotFound, err.Error())
		return
	}
	posts, status, err := h.service.Post.GetAllLikedPostsByUserID(r.Context(), uint(userID), false)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
