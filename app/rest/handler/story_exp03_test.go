//go:build unit
// +build unit

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseById(t *testing.T) {
	// Arrange
	paramId := "200"

	reqBody := model.RequestExpenses{
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	reqBodyJson, _ := json.Marshal(reqBody)

	resTb := model.ResultTbExpenses{
		Id:     paramId,
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/:id", bytes.NewReader(reqBodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Msgf("Sqlmock Error: %s", err.Error())
	}

	mockSqlStmt := "UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING *"
	mockReturnRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(resTb.Id, resTb.Title, resTb.Amount, resTb.Note, pq.Array(resTb.Tags))

	mock.ExpectQuery(regexp.QuoteMeta(mockSqlStmt)).
		WithArgs(paramId, reqBody.Title, reqBody.Amount, reqBody.Note, pq.Array(reqBody.Tags)).
		WillReturnRows((mockReturnRow))

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(paramId)

	resBody := model.ResponseExpenses{
		Id:     "200",
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	resBodyJson, _ := json.Marshal(resBody)
	expectedResBodyStr := string(resBodyJson)
	expectedResStateCode := http.StatusOK

	// Act
	err = h.UpdateExpenseById(c)
	actualResBodyStr := strings.TrimSpace(rec.Body.String())
	actualResStateCode := rec.Code

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, expectedResStateCode, actualResStateCode)
		assert.Equal(t, expectedResBodyStr, actualResBodyStr)
	}
}
