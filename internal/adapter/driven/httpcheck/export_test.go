// White-Box-Test (package httpcheck): der Grenzwert-Test übersteuert das
// unexportierte Body-Limit auf einen kleinen Wert, statt echte 64 MiB zu
// übertragen (R2-C-3). Der Dateiname export_test.go ist von testpackage
// ausgenommen (Default-skip-regexp).
package httpcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// DC-FA-SRC-001.a Schritt 3: ein Body über dem Limit ⇒ TooLarge (kein
// Hash-Schritt). maxBody hier auf 8 Byte übersteuert (Adapter-Feld, kein Global).
func TestFetch_BodyLimitUeberschritten(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("deutlich mehr als acht bytes"))
	}))
	defer srv.Close()

	a := New(2 * time.Second)
	a.maxBody = 8
	res := a.Fetch(srv.URL)
	if !res.TooLarge || res.Body != nil {
		t.Fatalf("res = %+v (TooLarge, kein Body erwartet)", res)
	}
}
