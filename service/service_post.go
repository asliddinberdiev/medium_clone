package service

import (
	"github.com/asliddinberdiev/medium_clone/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}
