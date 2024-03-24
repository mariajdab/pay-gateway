package AcqBank

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckTxn(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bank-sim.com/cards/245-3646-7477", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body.String())
}
