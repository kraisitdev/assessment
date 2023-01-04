package handler

import (
	"net/http"

	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseById(c echo.Context) error {
	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, model.Err{Message: "paramId is not found value"})
	}

	resTb := model.ResultTbExpenses{}
	row := h.DB.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", paramId)

	err := row.Scan(&resTb.Id, &resTb.Title, &resTb.Amount, &resTb.Note, pq.Array(&resTb.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Err{Message: err.Error()})
	}

	resBody := model.ResponseExpenses{
		Id:     resTb.Id,
		Title:  resTb.Title,
		Amount: resTb.Amount,
		Note:   resTb.Note,
		Tags:   resTb.Tags,
	}

	return c.JSON(http.StatusOK, resBody)
}
