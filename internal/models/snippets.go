package models

// model = DAL

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error){
	return 0, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Return the last 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error){
	return nil, nil
}