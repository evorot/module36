package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

// Создание объекта DB
func New(url string) (*postgres, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &postgres{db: db}, nil
}

func (p *postgres) PostsMany(posts []storage.Post) error {
	for _, post := range posts {
		err := p.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postgres) Posts(quantity int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM posts ORDER BY pubtime DESC LIMIT $1;", quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (p *postgres) AddPost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO posts (title, content, pubtime, link) VALUES ($1,$2, $3, $4);", post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		return err
	}
	return nil
}

func (p postgres) UpdatePost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"UPDATE posts "+
			"SET title = $1, "+
			"content = $2, "+
			"pubtime = $3,"+
			"link = $4 "+
			"WHERE id = $5", post.Title, post.Content, post.PubTime, post.Link, post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgres) DeletePost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"DELETE FROM posts WHERE id=$1;", post.ID)
	if err != nil {
		return err
	}
	return nil
}
