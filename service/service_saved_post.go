package service

import (
	"database/sql"
	"log"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/google/uuid"
)

type SavedPostService struct {
	repo repository.SavedPost
}

func NewSavedPostService(repo repository.SavedPost) *SavedPostService {
	return &SavedPostService{repo: repo}
}

func (s *SavedPostService) Add(savedPost models.SavedPostAction) error {

	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("service_saved_post: create - uuid generate error: ", err)
		return err
	}
	return s.repo.Add(models.SavedPost{ID: id.String(), UserID: savedPost.UserID, PostID: savedPost.PostID})
}

func (s *SavedPostService) Remove(post_id string) error {
	return s.repo.Remove(post_id)
}

func (s *SavedPostService) GetByID(user_id, post_id string) (*models.SavedPost, error) {
	return s.repo.GetByID(user_id, post_id)
}

func (s *SavedPostService) GetAll(user_id string) ([]*models.Post, error) {
	list, err := s.repo.GetAll(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Post{}, nil
		}
		log.Println("service_saved_post: GetAll - repo error: ", err)
		return nil, err
	}

	if list == nil {
		return []*models.Post{}, nil
	}

	return list, nil
}
