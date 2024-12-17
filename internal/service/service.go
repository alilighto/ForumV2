package service

import (
	"context"

	"forum/internal/entity"
	"forum/internal/repository"
)

type User interface {
	Create(ctx context.Context, user entity.User) (int, error)
	SignIn(ctx context.Context, user entity.User) (string, int, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, int, error)
}

type Session interface {
	IsTokenExist(ctx context.Context, token string) (bool, error)
	DeleteSessionByToken(ctx context.Context, token string) error
	DeleteSessionByUserID(ctx context.Context, userID uint) error
}

type Post interface {
	CreatePost(ctx context.Context, input entity.Post) (uint, int, error)
	DeletePostByID(ctx context.Context, postID uint, userID uint) (int, error)
	UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error)
	GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error)
	GetAllByCategory(ctx context.Context, tagName string, limit, offset int) ([]entity.Post, int, error)
	GetAllByUserID(ctx context.Context, userID uint) ([]entity.Post, int, error)
	GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool) ([]entity.Post, int, error)
}

type Category interface {
	GetAllCategorys(ctx context.Context) ([]entity.Category, int, error)
}

type Comment interface {
	CreateComment(ctx context.Context, input entity.Comment) (int, error)
	DeleteComment(ctx context.Context, commentID uint, userID uint) (int, error)
	UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error)
}

type Service struct {
	User
	Session
	Post
	Comment
	Category
}

func NewService(repo *repository.Repository, secret string) *Service {
	return &Service{
		User:     newUserService(repo.User, repo.Session, secret),
		Session:  newSessionService(repo.Session),
		Post:     newPostService(repo.Post, repo.Category),
		Comment:  newCommentService(repo.Comment),
		Category: newCategoryService(repo.Category),
	}
}
