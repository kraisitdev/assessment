//go:build integration
// +build integration

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kraisitdev/assessment/app/rest/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestITGetExpenseAll(t *testing.T) {
	// Setup server
	port := os.Getenv("PORT")
	host := fmt.Sprintf("localhost:%s", port)
	endpoint := fmt.Sprintf("http://localhost:%s/expenses", port)

	eh := echo.New()
	go func(e *echo.Echo) {
		h := NewApp(false)

		e.GET("/expenses", h.GetExpenseAll)
		e.Start(":" + port)
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", host, 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	resBody := []model.ResponseExpenses{
		{
			Id:     "100",
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		},
		{
			Id:     "1",
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		},
		{
			Id:     "200",
			Title:  "apple smoothie",
			Amount: 89,
			Note:   "no discount",
			Tags:   []string{"beverage"},
		}}

	resBodyJson, _ := json.Marshal(resBody)
	expectedResBodyStr := string(resBodyJson)
	expectedResStateCode := http.StatusOK

	actualResBodyStr := strings.TrimSpace(string(byteBody))
	actualResStateCode := resp.StatusCode

	if assert.NoError(t, err) {
		assert.Equal(t, expectedResStateCode, actualResStateCode)
		assert.Equal(t, expectedResBodyStr, actualResBodyStr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
