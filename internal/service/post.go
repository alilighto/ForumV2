package service

import (
	"context"
	"errors"
	"forum/internal/entity"
	"forum/internal/repository"
	"log"
	"net/http"
	"strings"
)

type PostService struct {
	postRepo repository.Post
	tagRepo  repository.Tag
}

func newPostService(postRepo repository.Post, tagRepo repository.Tag) *PostService {
	return &PostService{
		postRepo: postRepo,
		tagRepo:  tagRepo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, input entity.Post) (uint, int, error) {
	if input.Data == "" || len(input.Data) > 10000 {
		return 0, http.StatusBadRequest, errors.New("data is empty")
	} else if input.Title == "" || len(input.Title) > 58 {
		return 0, http.StatusBadRequest, errors.New("title is empty")
	} else if len(input.Tags) == 0 || len(input.Tags) > 5 {
		return 0, http.StatusBadRequest, errors.New("tags is empty")
	}
	for _, tag := range input.Tags {
		if len(tag) == 0 || len(tag) > 20 {
			return 0, http.StatusBadRequest, errors.New("invalid tag")
		}
	}
	postID, status, err := s.postRepo.CreatePost(ctx, input)
	if err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	input.Tags = append(input.Tags, "ALL")
	if status, err := s.tagRepo.CreateTags(ctx, input.Tags); err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	tagIDS, status, err := s.tagRepo.GetTagsIDByName(ctx, input.Tags)
	if err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	if status, err := s.tagRepo.CreateTagsAndPostCon(ctx, tagIDS, postID); err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	return postID, http.StatusOK, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error) {
	return s.postRepo.GetPostByID(ctx, postID)
}

func (s *PostService) DeletePostByID(ctx context.Context, postID uint, userID uint) (int, error) {
	return s.postRepo.DeletePostByID(ctx, postID, userID)
}

func (s *PostService) UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error) {
	if input.Vote != 0 && input.Vote != 1 {
		return http.StatusBadRequest, errors.New("invalid vote")
	}
	return s.postRepo.UpsertPostVote(ctx, input)
}

func (s *PostService) GetAllByTag(ctx context.Context, tagName string) ([]entity.Post, int, error) {
	if strings.TrimSpace(tagName) == "" {
		return nil, http.StatusBadRequest, errors.New("invalid tag")
	}
	return s.postRepo.GetAllByTag(ctx, tagName)
}

func (s *PostService) GetAllByUserID(ctx context.Context, userID uint) ([]entity.Post, int, error) {
	return s.postRepo.GetAllByUserID(ctx, userID)
}

func (s *PostService) GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool) ([]entity.Post, int, error) {
	return s.postRepo.GetAllLikedPostsByUserID(ctx, userID, islike)
}
