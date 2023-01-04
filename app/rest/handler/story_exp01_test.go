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

func TestInsertExpense(t *testing.T) {
	// Arrange
	reqBody := model.RequestExpenses{
		Title:  "strawberry smoothie",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	reqBodyJson, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewReader(reqBodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Msgf("Sqlmock Error: %s", err.Error())
	}

	mockSqlStmt := "INSERT INTO expenses (id, title, amount, note, tags) values (default, $1, $2, $3, $4) RETURNING id"
	mockReturnRow := sqlmock.NewRows([]string{"id"}).AddRow("1")

	mock.ExpectQuery(regexp.QuoteMeta(mockSqlStmt)).
		WithArgs(reqBody.Title, reqBody.Amount, reqBody.Note, pq.Array(reqBody.Tags)).
		WillReturnRows((mockReturnRow))

	h := handler{db}
	c := e.NewContext(req, rec)

	resBody := model.ResponseExpenses{
		Id:     "1",
		Title:  "strawberry smoothie",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	resBodyJson, _ := json.Marshal(resBody)
	expectedResBodyStr := string(resBodyJson)
	expectedResStateCode := http.StatusCreated

	// Act
	err = h.InsertExpense(c)
	actualResBodyStr := strings.TrimSpace(rec.Body.String())
	actualResStateCode := rec.Code

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, expectedResStateCode, actualResStateCode)
		assert.Equal(t, expectedResBodyStr, actualResBodyStr)
	}
}
