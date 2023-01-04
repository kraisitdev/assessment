package handler

import (
	"net/http"

	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenseById(c echo.Context) error {
	paramId := c.Param("id")
	if paramId == "" {
		return c.JSON(http.StatusBadRequest, model.Err{Message: "paramId is not found value"})
	}

	reqBody := model.RequestExpenses{}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Err{Message: err.Error()})
	}

	resTb := model.ResultTbExpenses{}
	row := h.DB.QueryRow("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING *",
		paramId, reqBody.Title, reqBody.Amount, reqBody.Note, pq.Array(reqBody.Tags))

	err = row.Scan(&resTb.Id, &resTb.Title, &resTb.Amount, &resTb.Note, pq.Array(&resTb.Tags))
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
