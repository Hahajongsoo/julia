package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
	"time"
)

type AdminMemoService interface {
	GetAdminMemo() (*models.AdminMemo, error)
	SaveAdminMemo(memo *models.AdminMemo) (*models.AdminMemo, error)
	DeleteAdminMemo() error
}

type adminMemoService struct {
	memoRepo repositories.AdminMemoRepository
}

func NewAdminMemoService(memoRepo repositories.AdminMemoRepository) AdminMemoService {
	return &adminMemoService{
		memoRepo: memoRepo,
	}
}

func (s *adminMemoService) GetAdminMemo() (*models.AdminMemo, error) {
	return s.memoRepo.GetAdminMemo()
}

func (s *adminMemoService) SaveAdminMemo(memo *models.AdminMemo) (*models.AdminMemo, error) {
	memo.UpdatedAt = time.Now()
	return s.memoRepo.SaveAdminMemo(memo)
}

func (s *adminMemoService) DeleteAdminMemo() error {
	return s.memoRepo.DeleteAdminMemo()
}
