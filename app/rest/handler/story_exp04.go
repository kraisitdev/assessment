package handler

import (
	"net/http"

	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseAll(c echo.Context) error {

	resTb := model.ResultTbExpenses{}
	rows, err := h.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Err{Message: err.Error()})
	}

	resBody := []model.ResponseExpenses{}
	for rows.Next() {
		err := rows.Scan(&resTb.Id, &resTb.Title, &resTb.Amount, &resTb.Note, pq.Array(&resTb.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Err{Message: err.Error()})
		}

		resBody = append(resBody, model.ResponseExpenses{
			Id:     resTb.Id,
			Title:  resTb.Title,
			Amount: resTb.Amount,
			Note:   resTb.Note,
			Tags:   resTb.Tags,
		})
	}

	return c.JSON(http.StatusOK, resBody)
}
