package handler

import (
	"github.com/emiliocc5/online-store-api/pkg/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	aValidClientId = 1
)

func Test_GetCart_Successful(t *testing.T) {
	cartServiceMock := &CartServiceMock{}

	cartServiceMock.On("GetCart", aValidClientId).Return(getMockedValidCartResponse(), nil)

	x := []string{"1"}
	var header = make(map[string][]string)
	header["clientId"] = x
	request := http.Request{Header: header}
	context := gin.Context{
		Request: &request,
	}

	handl := handler.CartHandlerImpl{CartService: cartServiceMock}

	handl.HandleGetCart(&context)

	assert.Equal(t, http.StatusOK, context.Status)
}
