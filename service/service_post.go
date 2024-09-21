package service

import (
	"database/sql"
	"log"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/google/uuid"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(user_id string, post models.CreatePost) (*models.Post, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("service_post: create - uuid generate error: ", err)
		return nil, err
	}

	newPost, err := s.repo.Create(models.Post{ID: id.String(), UserID: user_id, Title: post.Title, Body: post.Body, Published: post.Published})
	if err != nil {
		log.Println("service_post: create - repo error: ", err)
		return nil, err
	}

	return newPost, nil
}

func (s *PostService) GetByID(id string) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		log.Println("service_post: GetByID - repo error: ", err)
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPersonal(user_id string) ([]*models.Post, error) {
	list, err := s.repo.GetPersonal(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Post{}, nil
		}
		log.Println("service_post: GetPersonal - repo error: ", err)
		return nil, err
	}

	if list == nil {
		return []*models.Post{}, nil
	}

	return list, nil
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Post{}, nil
		}
		log.Println("service_post: GetAll - repo error: ", err)
		return nil, err
	}

	if list == nil {
		return []*models.Post{}, nil
	}

	return list, nil
}

func (s *PostService) Update(id string, post models.UpdatePost) (*models.Post, error) {
	dbPost, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("service_post: update - not found dbPost")
			return nil, err
		}
		log.Println("service_post: update - get post repo error: ", err)
		return nil, err
	}

	if post.Title == "" {
		post.Title = dbPost.Title
	}
	if post.Body == "" {
		post.Body = dbPost.Body
	}
	if post.Published == nil {
		post.Published = &dbPost.Published
	}

	updatePost, err := s.repo.Update(id, post)
	if err != nil {
		log.Println("service_post: update - update repo error: ", err)
		return nil, err
	}

	return updatePost, nil
}

func (s *PostService) Delete(id string) error {
	return s.repo.Delete(id)
}
