package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Hello(t *testing.T) {
	type wantS struct {
		body string
		code int
	}
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	tests := []struct {
		name     string
		params   generated.HelloParams
		mockFunc func(ctx context.Context, input repository.GetTestByIdInput)
		want     wantS
	}{
		{
			name: "success",
			params: generated.HelloParams{
				Id: "1",
			},
			mockFunc: func(ctx context.Context, input repository.GetTestByIdInput) {
				mockRepo.EXPECT().GetTestById(ctx, input).Return(repository.QueryOutput{Name: "user"}, nil)
			},
			want: wantS{
				body: `{"message":"Hello User user"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "failed",
			params: generated.HelloParams{
				Id: "1",
			},
			mockFunc: func(ctx context.Context, input repository.GetTestByIdInput) {
				mockRepo.EXPECT().GetTestById(ctx, input).Return(repository.QueryOutput{}, errors.New("err"))
			},
			want: wantS{
				body: `{"message":"User not found"}`,
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			test.mockFunc(c.Request().Context(), repository.GetTestByIdInput{Id: test.params.Id})
			s := Server{
				Repository: mockRepo,
			}
			err := s.Hello(c, test.params)
			if err != nil {
				t.Errorf("Hello() error = %+v", err)
				return
			}
			if !assert.Equal(t, test.want.code, rec.Code) {
				t.Errorf("Hello() code = %v, want %v", rec.Code, test.want.code)
			}
			if !assert.Equal(t, strings.TrimSpace(test.want.body), strings.TrimSpace(rec.Body.String())) {
				t.Errorf("Hello() body = %+v, want %+v", rec.Body.String(), test.want.body)
			}
		})
	}
}

func Test_GetMyProfile(t *testing.T) {
	type wantS struct {
		body string
		code int
	}
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	tests := []struct {
		name     string
		mockFunc func(ctx context.Context)
		token    string
		want     wantS
		err      string
	}{
		{
			name:  "success",
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiJzc3Nzc3MiLCJwaG9uZV9udW1iZXIiOiI5MjIyMCJ9.NrROF0UGK88Oq2073VE_A55-HkANbRMHajbrBDPnf8w",
			want: wantS{
				body: `{"name":"ssssss","phone_number":"92220"}`,
				code: http.StatusOK,
			},
		},
		{
			name:  "failed",
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want: wantS{
				body: ``,
				code: http.StatusOK,
			},
			err: "code=403, message=Token is not valid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			req.Header.Set("Authorization", test.token)

			s := Server{
				Repository: mockRepo,
			}
			err := s.GetMyProfile(c)
			if err != nil && err.Error() != test.err {
				t.Errorf("Hello() err = %v, want %v", err.Error(), test.err)
			}
			if !assert.Equal(t, test.want.code, rec.Code) {
				t.Errorf("Hello() code = %v, want %v", rec.Code, test.want.code)
			}
			if !assert.Equal(t, strings.TrimSpace(test.want.body), strings.TrimSpace(rec.Body.String())) {
				t.Errorf("Hello() body = %+v, want %+v", rec.Body.String(), test.want.body)
			}
		})
	}
}

func Test_UpdateMyProfile(t *testing.T) {
	type wantS struct {
		body string
		code int
	}
	pn1 := "+62888732928"
	pn2 := "123"
	fn1 := "namakuu"
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	tests := []struct {
		name     string
		params   generated.UpdateMyProfileParams
		mockFunc func()
		token    string
		err      string
		want     wantS
	}{
		{
			name: "success update phone number",
			params: generated.UpdateMyProfileParams{
				PhoneNumber: &pn1,
			},
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiJzc3Nzc3MiLCJwaG9uZV9udW1iZXIiOiI5MjIyMCJ9.NrROF0UGK88Oq2073VE_A55-HkANbRMHajbrBDPnf8w",
			mockFunc: func() {
				mockRepo.EXPECT().GetUserData(gomock.Any(), gomock.Any()).Return(repository.QueryOutput{}, nil)
				mockRepo.EXPECT().UpdatePhoneNumber(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: wantS{
				body: `{"message":"Successfully updated user data"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "success update name",
			params: generated.UpdateMyProfileParams{
				FullName: &fn1,
			},
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiJzc3Nzc3MiLCJwaG9uZV9udW1iZXIiOiI5MjIyMCJ9.NrROF0UGK88Oq2073VE_A55-HkANbRMHajbrBDPnf8w",
			mockFunc: func() {
				mockRepo.EXPECT().UpdateName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: wantS{
				body: `{"message":"Successfully updated user data"}`,
				code: http.StatusOK,
			},
		},
		{
			name:     "no request params",
			params:   generated.UpdateMyProfileParams{},
			token:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiJzc3Nzc3MiLCJwaG9uZV9udW1iZXIiOiI5MjIyMCJ9.NrROF0UGK88Oq2073VE_A55-HkANbRMHajbrBDPnf8w",
			mockFunc: func() {},
			err:      "code=400, message=invalid request data",
			want: wantS{
				body: ``,
				code: http.StatusOK,
			},
		},
		{
			name: "invalid phone number",
			params: generated.UpdateMyProfileParams{
				PhoneNumber: &pn2,
			},
			token:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiJzc3Nzc3MiLCJwaG9uZV9udW1iZXIiOiI5MjIyMCJ9.NrROF0UGK88Oq2073VE_A55-HkANbRMHajbrBDPnf8w",
			mockFunc: func() {},
			err:      "code=400, message=invalid request data",
			want: wantS{
				body: ``,
				code: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			req.Header.Set("Authorization", test.token)

			test.mockFunc()
			s := Server{
				Repository: mockRepo,
			}
			err := s.UpdateMyProfile(c, test.params)
			if err != nil && err.Error() != test.err {
				t.Errorf("UpdateMyProfile() err = %v, want %v", err.Error(), test.err)
			}
			if !assert.Equal(t, test.want.code, rec.Code) {
				t.Errorf("UpdateMyProfile() code = %v, want %v", rec.Code, test.want.code)
			}
			if !assert.Equal(t, strings.TrimSpace(test.want.body), strings.TrimSpace(rec.Body.String())) {
				t.Errorf("UpdateMyProfile() body = %+v, want %+v", rec.Body.String(), test.want.body)
			}
		})
	}
}

func Test_PostLogin(t *testing.T) {
	type wantS struct {
		body string
		code int
	}

	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	tests := []struct {
		name     string
		params   generated.PostLoginParams
		mockFunc func()
		token    string
		err      string
		want     wantS
	}{
		{
			name: "success",
			params: generated.PostLoginParams{
				PhoneNumber: "+62888732928",
				Password:    "aaaaA1&",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetUserData(gomock.Any(), gomock.Any()).Return(repository.QueryOutput{
					Password: "$2a$10$aWgB75Bp8DtTygKxoKXUsuGA2cE/eycXqT3YvoproS9BAJiu5fpSS",
				}, nil)
				mockRepo.EXPECT().Logged(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: wantS{
				body: `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmdWxsX25hbWUiOiIiLCJwaG9uZV9udW1iZXIiOiIrNjI4ODg3MzI5MjgifQ.mMEb2U8m4xcS___RaSvTJ5vPaM6wss9TbjWI3rFfGGU"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "wrong password",
			params: generated.PostLoginParams{
				PhoneNumber: "+62888732928",
				Password:    "aabaA1&",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetUserData(gomock.Any(), gomock.Any()).Return(repository.QueryOutput{
					Password: "$2a$10$aWgB75Bp8DtTygKxoKXUsuGA2cE/eycXqT3YvoproS9BAJiu5fpSS",
				}, nil)
			},
			err: "code=400, message=Invalid phone number or password",
			want: wantS{
				body: ``,
				code: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()
			s := Server{
				Repository: mockRepo,
			}
			err := s.PostLogin(c, test.params)
			if err != nil && err.Error() != test.err {
				t.Errorf("PostLogin() err = %v, want %v", err.Error(), test.err)
			}
			if !assert.Equal(t, test.want.code, rec.Code) {
				t.Errorf("PostLogin() code = %v, want %v", rec.Code, test.want.code)
			}
			if !assert.Equal(t, strings.TrimSpace(test.want.body), strings.TrimSpace(rec.Body.String())) {
				t.Errorf("PostLogin() body = %+v, want %+v", rec.Body.String(), test.want.body)
			}
		})
	}
}

func Test_PostSignup(t *testing.T) {
	type wantS struct {
		body string
		code int
	}
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	tests := []struct {
		name     string
		params   generated.PostSignupParams
		mockFunc func()
		token    string
		err      string
		want     wantS
	}{
		{
			name: "success",
			params: generated.PostSignupParams{
				FullName: "aaa",	
				PhoneNumber: "+62888732928",
				Password:    "aabaA1&",
			},
			mockFunc: func() {
				mockRepo.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(repository.QueryOutput{ID: 1}, nil)
			},
			want: wantS{
				body: `{"ID":"1","message":"Successfully signed up"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "data exists",
			params: generated.PostSignupParams{
				FullName: "aaa",	
				PhoneNumber: "+62888732928",
				Password:    "aabaA1&",
			},
			mockFunc: func() {
				mockRepo.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(repository.QueryOutput{}, errors.New("data exists"))
			},
			want: wantS{
				body: ``,
				code: http.StatusOK,
			},
			err: "code=400, message=Account already exists",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mockFunc()
			s := Server{
				Repository: mockRepo,
			}
			err := s.PostSignup(c, test.params)
			if err != nil && err.Error() != test.err {
				t.Errorf("PostSignup() err = %v, want %v", err.Error(), test.err)
			}
			if !assert.Equal(t, test.want.code, rec.Code) {
				t.Errorf("PostSignup() code = %v, want %v", rec.Code, test.want.code)
			}
			if !assert.Equal(t, strings.TrimSpace(test.want.body), strings.TrimSpace(rec.Body.String())) {
				t.Errorf("PostSignup() body = %+v, want %+v", rec.Body.String(), test.want.body)
			}
		})
	}
}
