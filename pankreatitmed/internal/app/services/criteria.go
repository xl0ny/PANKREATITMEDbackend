// internal/app/service/criteria_service.go
package services

import (
	"errors"
	"fmt"
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// минимальный контракт репозитория, НУЖНЫЙ именно этому сервису

type CriteriaService interface {
	List(q string) ([]ds.Criterion, error)
	Get(id uint) (*ds.Criterion, error)
	Create(c *ds.Criterion) error
	Update(id uint, in *request.UpdateCriterion) error
	Delete(id uint) error
	ToDraft(id, user_id uint) error
	DeleteImage(client *minio.Client, critID uint, c *gin.Context) error
}

type criteriaService struct {
	repo CriteriaRepoPort
}

func NewCriteriaService(repo CriteriaRepoPort) CriteriaService {
	return &criteriaService{repo: repo}
}

func (s *criteriaService) List(q string) ([]ds.Criterion, error) {
	return s.repo.GetCriteria(q)
}

func (s *criteriaService) Get(id uint) (*ds.Criterion, error) {
	return s.repo.GetCriterionByID(id)
}

func (s *criteriaService) Create(c *ds.Criterion) error {
	id, err := s.repo.GetSeq()
	if err != nil {
		return err
	}
	c.ID = id

	err = s.repo.CreateCriterion(c)
	if err != nil {
		s.repo.ResetCriterionSequence()
	}
	return err
}

func (s *criteriaService) Update(id uint, in *request.UpdateCriterion) error {
	return s.repo.UpdateCriterion(id, in)
}

func (s *criteriaService) Delete(id uint) error {
	return s.repo.DeleteCriterion(id)
}

func (s *criteriaService) ToDraft(id, user_id uint) error {
	oi, err := s.repo.GetOrCreateDraftPankreatitOrder(user_id)
	println(oi)
	if err != nil {
		return err
	}
	return s.repo.AddItem(oi.ID, id)
}

// func (s *criteriaService) UpdateImage(id uint, in *request.UpdateCriterion) error {
//
// }
// TODO че тут за поиск ошибки такой ущербный, разобраться
func (s *criteriaService) DeleteImage(client *minio.Client, critID uint, c *gin.Context) error {
	objectName, err := s.repo.GetImageName(critID)
	if objectName == "" || err != nil {
		return nil
	}
	fmt.Println(objectName)
	if err != nil {
		return err
	}
	parts := strings.SplitN(objectName, "/services-images/", 2)
	if len(parts) == 2 {
		objectName := parts[1] // "service_13/IMG_5587.jpeg"
		return client.RemoveObject(c, "services-images", objectName, minio.RemoveObjectOptions{})
	} else if len(parts) == 1 {
		objectName := parts[0] // "service_13/IMG_5587.jpeg"
		return client.RemoveObject(c, "services-images", objectName, minio.RemoveObjectOptions{})
	} else {
		return errors.New("Invalid old image url to delete")
	}
}
