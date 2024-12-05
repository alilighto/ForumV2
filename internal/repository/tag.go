package repository

import (
	"context"
	"database/sql"
	"net/http"
)

type TagRepository struct {
	db *sql.DB
}

func newTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) CreateTags(ctx context.Context, tagsName []string) (int, error) {
	query := "INSERT OR IGNORE INTO tags(name) VALUES ($1);"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	for _, tag := range tagsName {
		if _, err := prep.ExecContext(ctx, tag); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func (r *TagRepository) GetTagsIDByName(ctx context.Context, tagsName []string) ([]uint, int, error) {

	ids := []uint{}
	query := "SELECT id FROM tags WHERE name = $1;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	for _, tag := range tagsName {
		var id uint
		if err = prep.QueryRowContext(ctx, tag).Scan(&id); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		ids = append(ids, id)
	}

	return ids, http.StatusOK, nil
}

func (r *TagRepository) CreateTagsAndPostCon(ctx context.Context, tagsID []uint, postID uint) (int, error) {
	query := "INSERT INTO tag_and_post(tag_id, post_id) VALUES($1, $2);"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for _, tagID := range tagsID {
		if _, err := prep.ExecContext(ctx, tagID, postID); err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}
