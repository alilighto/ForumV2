package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

type PostRepository struct {
	db *sql.DB
}

func newPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAllByTag(ctx context.Context, tagName string) ([]entity.Post, int, error) {
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username
	FROM
		post p
		INNER JOIN tag_and_post tp ON p.id = tp.post_id
		INNER JOIN tags t ON tp.tag_id = t.id
		INNER JOIN users u ON u.id = p.user_id
	WHERE t.name = $1;
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()
	posts := []entity.Post{}
	rows, err := prep.QueryContext(ctx, tagName)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tags, status, err := r.getTagsByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetAllByUserID(ctx context.Context, userID uint) ([]entity.Post, int, error) {
	posts := []entity.Post{}
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username
	FROM
		post p
		INNER JOIN users u ON u.id = p.user_id
	WHERE p.user_id = $1;
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()
	rows, err := prep.QueryContext(ctx, userID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tags, status, err := r.getTagsByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool) ([]entity.Post, int, error) {
	posts := []entity.Post{}
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username
	FROM
		post p
	INNER JOIN users u on p.user_id = u.id
	INNER JOIN post_vote pv ON p.id = pv.post_id
	WHERE
		pv.user_id = $1 AND pv.vote = %d
	`
	if islike {
		query = fmt.Sprintf(query, 1)
	} else {
		query = fmt.Sprintf(query, 0)
	}
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rows, err := prep.QueryContext(ctx, userID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tags, status, err := r.getTagsByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error) {
	var post entity.Post
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username,
		COALESCE(COUNT(CASE WHEN pv.vote = 1 THEN 1 END), 0) AS voting,
		COALESCE(COUNT(CASE WHEN pv.vote = 0 THEN 1 END), 0) AS voting1
	FROM
		post p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_vote pv ON p.id = pv.post_id
	WHERE
		p.id = $1
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return post, http.StatusInternalServerError, err
	}
	if err := prep.QueryRowContext(ctx, postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName, &post.Likes, &post.Dislikes); err != nil {
		return post, http.StatusNotFound, err
	}
	tags, status, err := r.getTagsByPostID(ctx, postID)
	if err != nil {
		return post, status, err
	}
	post.Tags = tags

	comments, status, err := r.getCommentsByPostID(ctx, postID)
	if err != nil {
		return post, status, err
	}
	post.Comments = comments
	return post, http.StatusOK, nil
}

func (r *PostRepository) getCommentsByPostID(ctx context.Context, postID uint) ([]entity.Comment, int, error) {
	query := `
	SELECT 
		c.id,
		c.user_id,
		c.data,
		u.username,
		COALESCE(COUNT(CASE WHEN cv.vote = 1 THEN 1 END), 0) AS voting,
		COALESCE(COUNT(CASE WHEN cv.vote = 0 THEN 1 END), 0) AS voting1
	FROM 
		comment c
		INNER JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_vote cv ON c.id = cv.comment_id
	WHERE 
		c.post_id = $1
	GROUP BY
		c.id, c.user_id, c.data, u.username;
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, postID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	comments := []entity.Comment{}
	for rows.Next() {
		comment := entity.Comment{}
		if err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.Data, &comment.UserName, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		comment.PostID = postID
		comments = append(comments, comment)
	}
	return comments, http.StatusOK, nil
}

func (r *PostRepository) getTagsByPostID(ctx context.Context, postID uint) ([]string, int, error) {
	query := `
	SELECT
		t.name
	FROM 
		tags t
		INNER JOIN tag_and_post tp ON tp.post_id = $1 and t.id = tp.tag_id
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rows, err := prep.QueryContext(ctx, postID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, http.StatusBadRequest, err
		}
		tags = append(tags, tag)
	}
	return tags, http.StatusOK, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, input entity.Post) (uint, int, error) {
	query := `INSERT INTO post(user_id, title, data) VALUES($1, $2, $3) RETURNING id;`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	defer prep.Close()
	var id uint
	if err = prep.QueryRowContext(ctx, input.UserID, input.Title, input.Data).Scan(&id); err != nil {
		return 0, http.StatusBadRequest, err
	}
	return id, http.StatusOK, nil
}

func (r *PostRepository) DeletePostByID(ctx context.Context, PostID uint, userID uint) (int, error) {
	query := `DELETE FROM post WHERE id = $1 AND user_id = $2`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, PostID, userID); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (r *PostRepository) UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error) {
	query := "SELECT vote FROM post_vote WHERE user_id = $1 and post_id = $2;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	var vote int
	if err := prep.QueryRowContext(ctx, input.UserID, input.PostID).Scan(&vote); err != nil {
		if err == sql.ErrNoRows {
			query = "INSERT INTO post_vote(user_id, post_id, vote) VALUES($1, $2, $3);"
			if _, err = r.db.ExecContext(ctx, query, input.UserID, input.PostID, input.Vote); err != nil {
				return http.StatusBadRequest, err
			}
		} else {
			return http.StatusInternalServerError, err
		}
	} else {
		if vote == input.Vote {
			query = "DELETE FROM post_vote WHERE user_id = $1 and post_id = $2;"
			if _, err := r.db.ExecContext(ctx, query, input.UserID, input.PostID); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			query = "UPDATE post_vote SET vote = $1 WHERE user_id = $2 and post_id = $3;"
			if _, err := r.db.ExecContext(ctx, query, input.Vote, input.UserID, input.PostID); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}
	return http.StatusOK, nil
}
