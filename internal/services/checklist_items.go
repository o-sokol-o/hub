package service

import (
	"github.com/AquaEngineering/AquaHub/internal/domain"
)

type ChecklistItemService struct {
	repo     IStoreChecklistItem
	listRepo IStoreChecklist
}

func NewChecklistItemService(repo IStoreChecklistItem, listRepo IStoreChecklist) *ChecklistItemService {
	return &ChecklistItemService{repo: repo, listRepo: listRepo}
}

func (s *ChecklistItemService) Create(userId, listId int, item domain.ChecklistItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// list does not exists or does not belongs to user
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *ChecklistItemService) GetAll(userId, listId int) ([]domain.ChecklistItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *ChecklistItemService) GetById(userId, listId, itemId int) (domain.ChecklistItem, error) {
	return s.repo.GetById(userId, listId, itemId)
}

func (s *ChecklistItemService) Update(userId, listId, itemId int, input domain.UpdateChecklistItem) error {
	return s.repo.Update(userId, listId, itemId, input)
}

func (s *ChecklistItemService) Delete(userId, listId, itemId int) error {
	return s.repo.Delete(userId, listId, itemId)
}
