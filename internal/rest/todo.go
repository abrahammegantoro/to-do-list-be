package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TodoService interface {
	Fetch(ctx context.Context, page int64, limit int64) ([]domain.Todo, error)
	GetByID(ctx context.Context, id int64) (domain.Todo, error)
	GetByUserID(ctx context.Context, userID int64, page int64, limit int64, category *string, priorityLevel *string, keyword *string) ([]domain.Todo, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	Store(ctx context.Context, td *domain.Todo) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, td *domain.Todo) error
}

type TodoHandler struct {
	Service TodoService
}

const defaultLimit = 10

func NewTodoHandler(e *echo.Group, svc TodoService) {
	handler := &TodoHandler{
		Service: svc,
	}

	e.GET("", handler.GetByUserID)
	e.GET("/:id", handler.GetByID)
	e.GET("/categories", handler.GetAllCategories)
	e.POST("", handler.Store)
	e.DELETE("/:id", handler.Delete)
	e.PUT("/:id", handler.Update)
}

func (t *TodoHandler) FetchTodo(c echo.Context) error {
	limitString := c.QueryParam("limit")

	limit, err := strconv.Atoi(limitString)
	if err != nil || limit == 0 {
		limit = defaultLimit
	}

	pageString := c.QueryParam("page")
	page, err := strconv.Atoi(pageString)
	if err != nil || page == 0 {
		page = 1
	}

	ctx := c.Request().Context()

	listTd, err := t.Service.Fetch(ctx, int64(page), int64(limit))
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "success",
		"data":    listTd,
	})
}

func (t *TodoHandler) GetByUserID(c echo.Context) error {
    userId := c.Get("userId").(int64)

    limitString := c.QueryParam("limit")
    limit, err := strconv.Atoi(limitString)
    if err != nil || limit == 0 {
        limit = defaultLimit
    }

    pageString := c.QueryParam("page")
    page, err := strconv.Atoi(pageString)
    if err != nil || page == 0 {
        page = 1
    }

    category := c.QueryParam("category")
    priorityLevel := c.QueryParam("priority_level")
    keyword := c.QueryParam("keyword")

    ctx := c.Request().Context()

    listTd, err := t.Service.GetByUserID(ctx, userId, int64(page), int64(limit), &category, &priorityLevel, &keyword)
    if err != nil {
        logrus.Error(err)
        return c.JSON(getStatusCode(err), map[string]interface{}{
            "status":  getStatusCode(err),
            "message": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status":  http.StatusOK,
        "message": "success",
        "data":    listTd,
    })
}


func (t *TodoHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	id := int64(idP)
	ctx := c.Request().Context()

	td, err := t.Service.GetByID(ctx, id)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "success",
		"data":    td,
	})
}

func (t *TodoHandler) GetAllCategories(c echo.Context) error {
	ctx := c.Request().Context()
	listCategory, err := t.Service.GetAllCategories(ctx)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "success",
		"data":    listCategory,
	})
}

func (t *TodoHandler) Store(c echo.Context) (err error) {
	userId := c.Get("userId").(int64)

	var todo domain.Todo
	todo.UserID = userId
	
	err = c.Bind(&todo)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	var ok bool
	if ok, err = isRequestValid(&todo); !ok {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	ctx := c.Request().Context()
	err = t.Service.Store(ctx, &todo)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "success",
		"data":    todo,
	})
}

func (t *TodoHandler) Delete(c echo.Context) (err error) {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = t.Service.Delete(ctx, id)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "item successfully deleted",
	})
}

func (t *TodoHandler) Update(c echo.Context) (err error) {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	id := int64(idP)
	userId := c.Get("userId").(int64)

	var todo domain.Todo
	todo.UserID = userId
	err = c.Bind(&todo)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&todo); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	todo.ID = id
	err = t.Service.Update(ctx, &todo)
	if err != nil {
		logrus.Error(err)
		return c.JSON(getStatusCode(err), map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "success",
		"data":    todo,
	})
}
