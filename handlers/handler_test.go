package handlers_test

import (
	"github.com/rikuya98/go-poke-data-api/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func BenchmarkGetPokeDataHandler(b *testing.B) {
	reqURL := "/pokemon/1?lv=50&ef=252&in=31"
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)
	w := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handlers.GetPokeDataHandler(w, req)
	}
}
