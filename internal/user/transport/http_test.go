package transport

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/qreasio/go-starter-kit/internal/user"
	"github.com/qreasio/go-starter-kit/internal/user/mock"
	"github.com/qreasio/go-starter-kit/pkg/log"
	"github.com/qreasio/go-starter-kit/pkg/model"
)

// NewUserHTTP returns ne UserHTTP struct instance
func NewUserHTTPMock(ctrl *gomock.Controller, validator *validator.Validate) UserHTTP {
	user1 := model.User{
		Firstname: "Isak",
		Lastname:  "Rickyanto",
		Email:     "isak@ricky.com",
		Created:   time.Now(),
		Updated:   time.Now(),
	}
	user2 := model.User{
		Firstname: "Fafa",
		Lastname:  "Tjan",
		Email:     "fafa@tjan.com",
		Created:   time.Now(),
		Updated:   time.Now(),
	}
	mockListUsers := make([]model.User, 0)
	mockListUsers = append(mockListUsers, user1, user2)

	svc := mock.NewMockService(ctrl)

	req := &user.ListUsersRequest{Pagination: *model.NewPagination()}

	svc.EXPECT().ListUsers(gomock.Any(), req).Return(mockListUsers, nil)

	return UserHTTP{svc: svc, log: log.New()}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestUserHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userHTTP := NewUserHTTPMock(ctrl, validator.New())
	handler := RegisterUserRouter(userHTTP)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code: 200, got: %d", resp.StatusCode)
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("Invalid content type, expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	if body == "" {
		t.Fatalf(body)
	}
}
