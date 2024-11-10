package todo

import (
	"context"
	"time"

	"github.com/abrahammegantoro/to-do-list-be/domain"
)

type TodoRepository interface {
	Fetch(ctx context.Context, limit int64, offset int64) ([]domain.Todo, error)
	GetByID(ctx context.Context, id int64) (domain.Todo, error)
	GetByUserID(ctx context.Context, userID int64, limit int64, offset int64, category *string, priorityLevel *string, keyword *string) ([]domain.Todo, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	Store(ctx context.Context, td *domain.Todo) error
	Update(ctx context.Context, td *domain.Todo) error
	Delete(ctx context.Context, id int64) error
}

type TodoService struct {
	todoRepository TodoRepository
}

func NewTodoService(td TodoRepository) *TodoService {
	return &TodoService{
		todoRepository: td,
	}
}

func (t *TodoService) Fetch(ctx context.Context, page int64, limit int64) (res []domain.Todo, err error) {
	offset := (page - 1) * limit

	res, err = t.todoRepository.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return
}

func (t *TodoService) GetByID(ctx context.Context, id int64) (res domain.Todo, err error) {
	res, err = t.todoRepository.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (t *TodoService) GetByUserID(ctx context.Context, userID int64, page int64, limit int64, category *string, priorityLevel *string, keyword *string) (res []domain.Todo, err error) {
    offset := (page - 1) * limit

    res, err = t.todoRepository.GetByUserID(ctx, userID, limit, offset, category, priorityLevel, keyword)
    if err != nil {
        return nil, err
    }

    return
}

func (t *TodoService) GetAllCategories(ctx context.Context) (res []string, err error) {
	return t.todoRepository.GetAllCategories(ctx)
}

func (t *TodoService) Store(ctx context.Context, td *domain.Todo) (err error) {
	td.CreatedAt = time.Now()
	td.UpdatedAt = time.Now()
	return t.todoRepository.Store(ctx, td)
}

func (t *TodoService) Update(ctx context.Context, td *domain.Todo) (err error) {
	existedTodo, err := t.todoRepository.GetByID(ctx, td.ID)
	if err != nil {
		return
	}
	if existedTodo == (domain.Todo{}) {
		return domain.ErrNotFound
	}

	td.UpdatedAt = time.Now()
	return t.todoRepository.Update(ctx, td)
}

func (t *TodoService) Delete(ctx context.Context, id int64) (err error) {
	existedTodo, err := t.todoRepository.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedTodo == (domain.Todo{}) {
		return domain.ErrNotFound
	}

	return t.todoRepository.Delete(ctx, id)
}
