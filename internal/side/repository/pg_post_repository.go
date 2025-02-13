package repository

import (
	"slices"
	"strconv"
	"strings"

	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/oklog/ulid/v2"
)

// id
// body
// created_at
// updated_at
// author_id
// side_id
type PgPostRepository struct {
	logger log.Logger
}

func NewPgPostRepository(logger log.Logger) PostRepository[db.SqlExecutor] {
	return &PgPostRepository{
		logger: logger,
	}
}

func (pr *PgPostRepository) Save(ctx db.DBContext[db.SqlExecutor], post *entity.Post) error {
	query := `INSERT INTO posts (id, body, created_at, updated_at, author_id, side_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := ctx.
		Executor().
		ExecContext(ctx, query, post.ID.String(), post.Body, post.CreatedAt, post.UpdatedAt, post.Author.ID, post.Side.ID)
	if err != nil {
		pr.logger.Warnf("%+v", err)
		return err
	}
	return nil
}

func (pr *PgPostRepository) FindManyWithLimit(ctx db.DBContext[db.SqlExecutor], limit int) (entity.Posts, error) {
	query := `SELECT 
	p.id, p.body, p.created_at, p.updated_at,
	u.id, u.username,
	s.id, s.nick, s.description, s.created_at
	FROM posts AS p
	JOIN sides AS s ON p.side_id = s.id
	JOIN users AS u ON p.author_id = u.id
	ORDER BY p.id DESC
	LIMIT $1`

	var posts entity.Posts

	rows, err := ctx.Executor().QueryContext(ctx, query, limit)
	if err != nil {
		pr.logger.Warnf("%+v", err)
		return entity.EmptyPosts, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		post := &entity.Post{
			Author: &entity.Author{},
			Side:   &entity.Side{},
		}
		if err := rows.Scan(
			&id, &post.Body, &post.CreatedAt, &post.UpdatedAt,
			&post.Author.ID, &post.Author.Username,
			&post.Side.ID, &post.Side.Nick, &post.Side.Description, &post.Side.CreatedAt,
		); err != nil {
			pr.logger.Warnf("%+v", err)
			return entity.EmptyPosts, err
		}
		postID, err := ulid.Parse(id)
		if err != nil {
			pr.logger.Warnf("%+v", err)
			return entity.EmptyPosts, err
		}
		post.ID = postID
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		pr.logger.Warnf("%+v", err)
		return entity.EmptyPosts, err
	}

	return posts, nil
}

func (pr *PgPostRepository) FindManyWithULIDCursor(ctx db.DBContext[db.SqlExecutor], cursor pagination.ULIDCursor) (entity.Posts, error) {
	var q strings.Builder
	q.WriteString(`SELECT 
		p.id, p.body, p.created_at, p.updated_at,
		u.id, u.username,
		s.id, s.nick, s.description, s.created_at
	FROM posts AS p
	JOIN sides AS s ON p.side_id = s.id
	JOIN users AS u ON p.author_id = u.id`)

	var args []any
	var orderBy string

	next, isNextCursor := cursor.(*pagination.NextULIDCursor)
	if isNextCursor {
		if next.ID() != nil {
			q.WriteString(` WHERE p.id <= $1`)
			args = append(args, next.ID().String())
		}
		orderBy = ` ORDER BY p.id DESC`
	}

	prev, isPrevCursor := cursor.(*pagination.PrevULIDCursor)
	if isPrevCursor {
		if prev.ID() != nil {
			q.WriteString(` WHERE p.id >= $1`)
			args = append(args, prev.ID().String())
		}
	}

	q.WriteString(orderBy)

	if cursor.Limit() > 0 {
		q.WriteString(` LIMIT $` + strconv.Itoa(len(args)+1))
		args = append(args, cursor.Limit()+2)
	}

	rows, err := ctx.Executor().QueryContext(ctx, q.String(), args...)
	if err != nil {
		pr.logger.Warnf("%+v", err)
		return entity.EmptyPosts, err
	}
	defer rows.Close()

	var posts entity.Posts

	for rows.Next() {
		var post entity.Post
		var id string
		post.Author = &entity.Author{}
		post.Side = &entity.Side{}

		if err := rows.Scan(
			&id, &post.Body, &post.CreatedAt, &post.UpdatedAt,
			&post.Author.ID, &post.Author.Username,
			&post.Side.ID, &post.Side.Nick, &post.Side.Description, &post.Side.CreatedAt,
		); err != nil {
			pr.logger.Warnf("%+v", err)
			return entity.EmptyPosts, err
		}

		postID, err := ulid.Parse(id)
		if err != nil {
			return entity.EmptyPosts, err
		}
		post.ID = postID
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		pr.logger.Warnf("%+v", err)
		return entity.EmptyPosts, err
	}

	if isPrevCursor {
		slices.SortStableFunc(posts, func(a, b *entity.Post) int {
			return b.ID.Compare(a.ID)
		})
	}

	return posts, nil
}
