package service

import (
	"github.com/o-sokol-o/hub/internal/domain"
)

// Сервис для работы со списками

type ChecklistService struct {
	repo IStoreChecklist
}

// Также в нашем сервисе и понадобится репозиторий
// Добавим его в качестве поля нашей структуры и будем передавать в конструкторе.
func NewChecklistService(repo IStoreChecklist) *ChecklistService {
	return &ChecklistService{repo: repo}
}

// При создании списка мы будем передавать данные на следующий уровень - в репозиторий,
// поэтому в сервисе мы лишь будем возвращать аналогичный метод репозитория.
// Дополнительной логики мы реализовывать не будем.
func (s *ChecklistService) Create(userId int, list domain.CreateChecklist) (int, error) {
	return s.repo.Create(userId, list)
}

// Метод GetAll, который будет возвращать слайс списка вместе с ошибкой.
// func (s *ChecklistService) GetAllChecklist() ([]domain.Checklist, error) {
// 	// В сервисе мы будем вызывать аналогичный метод репозитория, поскольку дополнительной бизнес логики тут нет.
// 	return s.repo.GetAll_Checklist()
// }

// Метод GetAll, который будет принимать id пользователя
// и возвращать слайс списка вместе с ошибкой.
func (s *ChecklistService) GetAllChecklistOfUser(userId int) ([]domain.Checklist, error) {
	// В сервисе мы будем вызывать аналогичный метод репозитория, поскольку дополнительной бизнес логики тут нет.
	return s.repo.GetAll_ChecklistOfUser(userId)
}

func (s *ChecklistService) GetById(userId, listId int) (*domain.Checklist, error) {
	return s.repo.GetById(userId, listId)
}

func (s *ChecklistService) Update(userId int, input domain.UpdateChecklist) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, input)
}

func (s *ChecklistService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}
