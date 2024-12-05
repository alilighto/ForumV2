package http1

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/service"
	"forum/pkg/config"
	"net/http"
	"text/template"
)

type Handler struct {
	service *service.Service
	secret  string
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Role    uint
}

func NewHandler(service *service.Service, secret string) *Handler {
	return &Handler{
		service: service,
		secret:  secret,
	}
}

func (h *Handler) InitRoutes(conf *config.Conf) *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web/src"))
	mux.Handle("/src/", http.StripPrefix("/src/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./web/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err = tmpl.Execute(w, fmt.Sprintf("%v:%v", conf.API.Host, conf.API.Port)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
	mux.HandleFunc("/api/is-valid", h.isValidToken)
	routes := h.createRoutes()
	for _, route := range routes {
		if route.Role == entity.Roles.Authorized {
			mux.Handle(route.Path, h.corsMiddleWare(h.isAlreadyIdentified(route.Handler)))
		} else {
			mux.Handle(route.Path, h.corsMiddleWare(h.identify(route.Role, route.Handler)))
		}
	}
	return mux
}

func (h *Handler) createRoutes() []Route {
	return []Route{
		{
			Path:    "/api/signup",
			Handler: h.signUp,
			Role:    entity.Roles.Authorized,
		},
		{
			Path:    "/api/signin",
			Handler: h.signIn,
			Role:    entity.Roles.Authorized,
		},
		{
			Path:    "/api/signout",
			Handler: h.signOut,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/",
			Handler: h.profile,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/posts/",
			Handler: h.getAllPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/liked-posts/",
			Handler: h.getAllLikedPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/profile/disliked-posts/",
			Handler: h.getAllDisLikedPostsByUserID,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/post/create",
			Handler: h.createPost,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/posts/",
			Handler: h.getALLPosts,
			Role:    entity.Roles.Guest,
		},
		{
			Path:    "/api/post/",
			Handler: h.getPostbyID,
			Role:    entity.Roles.Guest,
		},
		{
			Path:    "/api/post/vote",
			Handler: h.votePost,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/post/delete/",
			Handler: h.deletePost,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/comment/create",
			Handler: h.createComment,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/comment/vote",
			Handler: h.voteComment,
			Role:    entity.Roles.User,
		},
		{
			Path:    "/api/comment/delete/",
			Handler: h.deleteComment,
			Role:    entity.Roles.User,
		},
	}
}
