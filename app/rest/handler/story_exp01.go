package handler

import (
	"net/http"

	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) InsertExpense(c echo.Context) error {
	reqBody := model.RequestExpenses{}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Err{Message: err.Error()})
	}

	resTb := model.ResultTbExpenses{}
	row := h.DB.QueryRow("INSERT INTO expenses (id, title, amount, note, tags) values (default, $1, $2, $3, $4) RETURNING id",
		reqBody.Title, reqBody.Amount, reqBody.Note, pq.Array(reqBody.Tags))

	err = row.Scan(&resTb.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Err{Message: err.Error()})
	}

	resBody := model.ResponseExpenses{
		Id:     resTb.Id,
		Title:  reqBody.Title,
		Amount: reqBody.Amount,
		Note:   reqBody.Note,
		Tags:   reqBody.Tags,
	}

	return c.JSON(http.StatusCreated, resBody)
}
