package handler

import (
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

func TestGetExpenseById(t *testing.T) {
	// Arrange
	paramId := "1"

	resTb := model.ResultTbExpenses{
		Id:     paramId,
		Title:  "strawberry smoothie",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Msgf("Sqlmock Error: %s", err.Error())
	}

	mockSqlStmt := "SELECT id, title, amount, note, tags FROM expenses WHERE id=$1"
	mockReturnRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(resTb.Id, resTb.Title, resTb.Amount, resTb.Note, pq.Array(resTb.Tags))

	mock.ExpectQuery(regexp.QuoteMeta(mockSqlStmt)).
		WithArgs(paramId).
		WillReturnRows((mockReturnRow))

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(paramId)

	resBody := model.ResponseExpenses{
		Id:     "1",
		Title:  "strawberry smoothie",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	resBodyJson, _ := json.Marshal(resBody)
	expectedResBodyStr := string(resBodyJson)
	expectedResStateCode := http.StatusOK

	// Act
	err = h.GetExpenseById(c)
	actualResBodyStr := strings.TrimSpace(rec.Body.String())
	actualResStateCode := rec.Code

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, expectedResStateCode, actualResStateCode)
		assert.Equal(t, expectedResBodyStr, actualResBodyStr)
	}
}
