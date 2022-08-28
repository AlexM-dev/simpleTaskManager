package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type storage struct {
	db *pgxpool.Pool
}

type task struct {
	id         int
	opened     int64
	closed     int64
	authorID   int
	assignedID int
	title      string
	content    string
}

func CreateTask(authorID, assignedID int, title, content string) task {
	return task{authorID: authorID, assignedID: assignedID, title: title, content: content}
}

func New(c string) (*storage, error) {
	db, err := pgxpool.Connect(context.Background(), c)
	if err != nil {
		return nil, err
	}
	s := storage{
		db: db,
	}
	return &s, nil
}

func (s *storage) Tasks(taskID, authorID int) ([]task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []task
	for rows.Next() {
		var t task
		err = rows.Scan(
			&t.id,
			&t.opened,
			&t.closed,
			&t.authorID,
			&t.assignedID,
			&t.title,
			&t.content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *storage) NewTask(t task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.title,
		t.content,
	).Scan(&id)
	return id, err
}

func (s *storage) UpdateTask(id, assignedID int, title, content string) {
	s.db.QueryRow(context.Background(), `
		UPDATE tasks
		SET title = $1, content = $2
		WHERE id = $3;
		`,
		title,
		content,
		id,
	)
}

func (s *storage) DeleteTask(id int) {
	s.db.QueryRow(context.Background(), `
		DELETE FROM tasks WHERE id = $1;
		`,
		id,
	)
}

func (s *storage) AddUSer(name string) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO users (name) VALUES ($1) RETURNING id;
		`,
		name,
	).Scan(&id)
	return id, err
}
