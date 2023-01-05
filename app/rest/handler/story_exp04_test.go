//go:build unit
// +build unit

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

func TestGetExpenseAll(t *testing.T) {
	// Arrange
	resTb := []model.ResultTbExpenses{{
		Id:     "1",
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}, {
		Id:     "2",
		Title:  "iPhone 14 Pro Max 1TB",
		Amount: 66900,
		Note:   "birthday gift from my love",
		Tags:   []string{"gadget"},
	}}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Msgf("Sqlmock Error: %s", err.Error())
	}

	mockSqlStmt := "SELECT id, title, amount, note, tags FROM expenses"
	mockReturnRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(resTb[0].Id, resTb[0].Title, resTb[0].Amount, resTb[0].Note, pq.Array(resTb[0].Tags)).
		AddRow(resTb[1].Id, resTb[1].Title, resTb[1].Amount, resTb[1].Note, pq.Array(resTb[1].Tags))

	mock.ExpectQuery(regexp.QuoteMeta(mockSqlStmt)).
		WillReturnRows((mockReturnRow))

	h := handler{db}
	c := e.NewContext(req, rec)

	resBody := []model.ResponseExpenses{{
		Id:     "1",
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}, {
		Id:     "2",
		Title:  "iPhone 14 Pro Max 1TB",
		Amount: 66900,
		Note:   "birthday gift from my love",
		Tags:   []string{"gadget"},
	}}

	resBodyJson, _ := json.Marshal(resBody)
	expectedResBodyStr := string(resBodyJson)
	expectedResStateCode := http.StatusOK

	// Act
	err = h.GetExpenseAll(c)
	actualResBodyStr := strings.TrimSpace(rec.Body.String())
	actualResStateCode := rec.Code

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, expectedResStateCode, actualResStateCode)
		assert.Equal(t, expectedResBodyStr, actualResBodyStr)
	}
}
