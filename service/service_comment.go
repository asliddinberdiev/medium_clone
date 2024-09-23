package service

import (
	"database/sql"
	"log"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/google/uuid"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(user_id string, comment models.CreateComment) (*models.Comment, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("service_comment: create - uuid generate error: ", err)
		return nil, err
	}

	newComment, err := s.repo.Create(models.Comment{ID: id.String(), UserID: user_id, PostID: comment.PostID, Body: comment.Body})
	if err != nil {
		log.Println("service_comment: create - repo error: ", err)
		return nil, err
	}

	return newComment, nil
}

func (s *CommentService) GetAll(post_id string) ([]*models.Comment, error) {
	list, err := s.repo.GetAll(post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Comment{}, nil
		}
		log.Println("service_comment: GetAll - repo error: ", err)
		return nil, err
	}

	if list == nil {
		return []*models.Comment{}, nil
	}

	return list, nil
}

func (s *CommentService) GetByID(id string) (*models.Comment, error) {
	return s.repo.GetByID(id)
}

func (s *CommentService) Update(id, body string) (*models.Comment, error) {
	return s.repo.Update(id, body)
}

func (s *CommentService) Delete(comment_id string) error {
	return s.repo.Delete(comment_id)
}
