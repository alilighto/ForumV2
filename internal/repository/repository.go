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
	GetAllByTag(ctx context.Context, tagName string) ([]entity.Post, int, error)
	GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error)
	GetAllByUserID(ctx context.Context, userID uint) ([]entity.Post, int, error)
	GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool) ([]entity.Post, int, error)
}

type Tag interface {
	CreateTags(ctx context.Context, tagsName []string) (int, error)
	GetTagsIDByName(ctx context.Context, tagsName []string) ([]uint, int, error)
	CreateTagsAndPostCon(ctx context.Context, tagsID []uint, postID uint) (int, error)
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
	Tag
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:    newUserRepository(db),
		Session: newSessionRepository(db),
		Post:    newPostRepository(db),
		Tag:     newTagRepository(db),
		Comment: newCommentRepository(db),
	}
}
