package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_UserGenerateLink(t *testing.T) {
	// Setup
	e := echo.New()
	util.InitLogger(e)

	query := map[string]string{
		"bundleId":    "com.nihao",
		"playerToken": "d123",
		"country":     "cn",
		"os_version":  "10.01",
	}

	query_empty_bundleid := map[string]string{
		"bundleId":    "",
		"playerToken": "d123",
		"country":     "cn",
		"os_version":  "10.01",
	}

	tests := []struct {
		query  map[string]string
		code   int
		hasErr bool
	}{
		{query, 200, false},
		{query_empty_bundleid, 404, true},
	}

	for _, test := range tests {
		req, err := util.NewHTTPReqeust(echo.POST, "/promotion/generatelink", test.query, nil, nil)
		if assert.NoError(t, err) {
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUserService := &service.UserMockSV{}
			mockUserService.On("GetPlayTokenByDeviceID", "d123").Return("u123", nil)
			service.Injector.MapTo(mockUserService, new(service.UserService))

			mockItuneSV := &service.ItunesMockSV{}
			mockItuneSV.On("GenerateAdLink", "com.nihao", "cn", "u123").Return("http://sdf/id1233?at=123", int64(1233), nil)
			service.Injector.MapTo(mockItuneSV, new(service.ItunesService))

			controller := NewUserController()
			err := controller.GenerateLink(c)

			// Assertions
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.code, rec.Code)
		}
	}
}
