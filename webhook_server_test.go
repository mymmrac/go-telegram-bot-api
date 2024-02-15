package telego

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fasthttp/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestFastHTTPWebhookServer_RegisterHandler(t *testing.T) {
	require.Implements(t, (*WebhookServer)(nil), FastHTTPWebhookServer{})

	addr := testAddress(t)

	s := FastHTTPWebhookServer{
		Logger:      testLoggerType{},
		Server:      &fasthttp.Server{},
		Router:      router.New(),
		SecretToken: "secret",
	}

	go func() {
		err := s.Start(addr)
		assert.NoError(t, err)
	}()

	err := s.RegisterHandler("/", func(_ context.Context, data []byte) error {
		if len(data) == 0 {
			return nil
		}

		return errTest
	})
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.Header.Set(WebhookSecretTokenHeader, s.SecretToken)
		s.Server.Handler(ctx)

		assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	})

	t.Run("error_method", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/")
		ctx.Request.Header.SetMethod(fasthttp.MethodGet)
		s.Server.Handler(ctx)

		assert.Equal(t, fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode())
	})

	t.Run("error_handler", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		ctx.Request.Header.Set(WebhookSecretTokenHeader, s.SecretToken)
		ctx.Request.SetBody([]byte("err"))
		s.Server.Handler(ctx)

		assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
	})

	t.Run("secret_token_invalid", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/")
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)
		s.Server.Handler(ctx)

		assert.Equal(t, fasthttp.StatusUnauthorized, ctx.Response.StatusCode())
	})

	err = s.Stop(context.Background())
	require.NoError(t, err)
}

func TestHTTPWebhookServer_RegisterHandler(t *testing.T) {
	require.Implements(t, (*WebhookServer)(nil), HTTPWebhookServer{})

	t.Run("error_start_fail", func(t *testing.T) {
		s := HTTPWebhookServer{
			Logger:   testLoggerType{},
			Server:   &http.Server{}, //nolint:gosec
			ServeMux: http.NewServeMux(),
		}

		testAddr := testAddress(t)
		go func() {
			err := http.ListenAndServe(testAddr, nil) //nolint:gosec
			assert.NoError(t, err)
		}()

		time.Sleep(time.Millisecond * 10)

		err := s.Start(testAddr)
		require.Error(t, err)
	})

	t.Run("end_to_end", func(t *testing.T) {
		s := HTTPWebhookServer{
			Logger:      testLoggerType{},
			Server:      &http.Server{}, //nolint:gosec
			ServeMux:    http.NewServeMux(),
			SecretToken: "secret",
		}

		go func() {
			err := s.Start(testAddress(t))
			assert.NoError(t, err)
		}()

		err := s.RegisterHandler("/", func(_ context.Context, data []byte) error {
			if len(data) == 0 {
				return nil
			}

			return errTest
		})
		require.NoError(t, err)

		t.Run("success", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(WebhookSecretTokenHeader, s.SecretToken)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusOK, rc.Code)
		})

		t.Run("error_method", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rc.Code)
		})

		t.Run("error_handler", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("err"))
			req.Header.Set(WebhookSecretTokenHeader, s.SecretToken)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusInternalServerError, rc.Code)
		})

		t.Run("error_read", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", errReader{})
			req.Header.Set(WebhookSecretTokenHeader, s.SecretToken)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusInternalServerError, rc.Code)
		})

		t.Run("error_close", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", errReaderCloser{reader: strings.NewReader("ok")})
			req.Header.Set(WebhookSecretTokenHeader, s.SecretToken)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusInternalServerError, rc.Code)
		})

		t.Run("secret_token_invalid", func(t *testing.T) {
			rc := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil)

			s.Server.Handler.ServeHTTP(rc, req)

			assert.Equal(t, http.StatusUnauthorized, rc.Code)
		})

		err = s.Stop(context.Background())
		require.NoError(t, err)
	})
}

type errReader struct{}

func (e errReader) Read(_ []byte) (n int, err error) {
	return 0, errTest
}

type errReaderCloser struct {
	reader io.Reader
}

func (e errReaderCloser) Close() error {
	return errTest
}

func (e errReaderCloser) Read(b []byte) (n int, err error) {
	return e.reader.Read(b)
}

func TestMultiBotWebhookServer_RegisterHandler(t *testing.T) {
	require.Implements(t, (*WebhookServer)(nil), &MultiBotWebhookServer{})

	ts := &testServer{}
	s := &MultiBotWebhookServer{
		Server: ts,
	}

	assert.Equal(t, 0, ts.started)
	assert.Equal(t, 0, ts.stopped)
	assert.Equal(t, 0, ts.registered)

	err := s.Start("")
	require.NoError(t, err)
	assert.Equal(t, 1, ts.started)

	err = s.Start("")
	require.NoError(t, err)
	assert.Equal(t, 1, ts.started)

	err = s.RegisterHandler("", nil)
	require.NoError(t, err)
	assert.Equal(t, 1, ts.registered)

	err = s.RegisterHandler("", nil)
	require.NoError(t, err)
	assert.Equal(t, 2, ts.registered)

	err = s.Stop(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, ts.stopped)

	err = s.Stop(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, ts.stopped)
}

type testServer struct {
	started    int
	stopped    int
	registered int
}

func (t *testServer) Start(_ string) error {
	t.started++
	return nil
}

func (t *testServer) Stop(_ context.Context) error {
	t.stopped++
	return nil
}

func (t *testServer) RegisterHandler(_ string, _ WebhookHandler) error {
	t.registered++
	return nil
}

func TestNoOpWebhookServer(t *testing.T) {
	require.Implements(t, (*WebhookServer)(nil), NoOpWebhookServer{})

	registered := false
	s := NoOpWebhookServer{
		RegisterHandlerFunc: func(path string, handler WebhookHandler) error {
			registered = true
			return nil
		},
	}

	err := s.Start("")
	require.NoError(t, err)
	err = s.Stop(nil) //nolint:staticcheck
	require.NoError(t, err)
	err = s.RegisterHandler("", nil)
	require.NoError(t, err)
	assert.True(t, registered)
}

func TestFuncWebhookServer(t *testing.T) {
	require.Implements(t, (*WebhookServer)(nil), FuncWebhookServer{})

	ts := &testServer{}
	s1 := FuncWebhookServer{
		Server: ts,
	}

	assert.Equal(t, 0, ts.started)
	assert.Equal(t, 0, ts.stopped)
	assert.Equal(t, 0, ts.registered)

	err := s1.Start("")
	require.NoError(t, err)

	err = s1.RegisterHandler("", nil)
	require.NoError(t, err)

	err = s1.Stop(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 1, ts.started)
	assert.Equal(t, 1, ts.stopped)
	assert.Equal(t, 1, ts.registered)

	started := 0
	stopped := 0
	registered := 0
	s2 := FuncWebhookServer{
		Server: ts,
		StartFunc: func(_ string) error {
			started++
			return nil
		},
		StopFunc: func(_ context.Context) error {
			stopped++
			return nil
		},
		RegisterHandlerFunc: func(_ string, _ WebhookHandler) error {
			registered++
			return nil
		},
	}

	err = s2.Start("")
	require.NoError(t, err)

	err = s2.RegisterHandler("", nil)
	require.NoError(t, err)

	err = s2.Stop(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 1, ts.started)
	assert.Equal(t, 1, ts.stopped)
	assert.Equal(t, 1, ts.registered)

	assert.Equal(t, 1, started)
	assert.Equal(t, 1, stopped)
	assert.Equal(t, 1, registered)
}
