package repositories

import (
	"encoding/json"
	"julia/internal/models"
	"os"
	"path/filepath"
	"time"
)

type AdminMemoRepository interface {
	GetAdminMemo() (*models.AdminMemo, error)
	SaveAdminMemo(memo *models.AdminMemo) (*models.AdminMemo, error)
	DeleteAdminMemo() error
}

type adminMemoRepository struct {
	filePath string
}

func NewAdminMemoRepository(dataDir string) AdminMemoRepository {
	return &adminMemoRepository{
		filePath: filepath.Join(dataDir, "admin_memo.json"),
	}
}

func (r *adminMemoRepository) GetAdminMemo() (*models.AdminMemo, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return &models.AdminMemo{
			Content:   "",
			UpdatedAt: time.Now(),
		}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var memo models.AdminMemo
	if err := json.Unmarshal(data, &memo); err != nil {
		return nil, err
	}

	return &memo, nil
}

func (r *adminMemoRepository) SaveAdminMemo(memo *models.AdminMemo) (*models.AdminMemo, error) {
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(memo, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return nil, err
	}

	return memo, nil
}

func (r *adminMemoRepository) DeleteAdminMemo() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(r.filePath)
}
