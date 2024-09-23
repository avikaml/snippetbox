package main

import "github.com/avikaml/snippetbox/internal/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

