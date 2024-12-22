package repository

import (
	"context"
	"database/sql"

	"forum/internal/entity"
)

type User interface {
	Create(ctx context.Context, user entity.User) (int, error)
	GetUserIDByEmail(ctx context.Context, email string) (entity.User, int, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, int, error)
}

type Session interface {
	IsTokenExist(ctx context.Context, token string) (bool, error)
	DeleteSessionByToken(ctx context.Context, token string) error
	PostSession(ctx context.Context, session entity.Session) (int, error)
	DeleteSessionByUserID(ctx context.Context, userID uint) error
}

type Post interface {
	CreatePost(ctx context.Context, input entity.Post) (uint, int, error)
	DeletePostByID(ctx context.Context, PostID uint, userID uint) (int, error)
	UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error)
	GetAllByCategory(ctx context.Context, categoryName string, limit, offset int) ([]entity.Post, int, error)
	GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error)
	GetAllByUserID(ctx context.Context, userID uint, limit, offset int) ([]entity.Post, int, error)
	GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool, limit, offset int) ([]entity.Post, int, error)
}

type Category interface {
	CreateCategorys(ctx context.Context, categorysName []string) (int, error)
	GetCategorysIDByName(ctx context.Context, categorysName []string) ([]uint, int, error)
	CreateCategorysAndPostCon(ctx context.Context, categorysID []uint, postID uint) (int, error)
	CategoryExist(ctx context.Context, categoryName string) (bool, int, error)
	GetAllCategories(ctx context.Context) ([]entity.Category, int, error)
}

type Comment interface {
	CreateComment(ctx context.Context, input entity.Comment) (int, error)
	DeleteComment(ctx context.Context, commentID uint, userID uint) (int, error)
	UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error)
}

type Repository struct {
	Post
	User
	Session
	Category
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:     newUserRepository(db),
		Session:  newSessionRepository(db),
		Post:     newPostRepository(db),
		Category: newCategoryRepository(db),
		Comment:  newCommentRepository(db),
	}
}
