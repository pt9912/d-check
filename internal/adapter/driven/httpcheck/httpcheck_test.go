// Adapter-Tests gegen httptest.Server — nie gegen echte URLs
// (Slice-Risiko slice-008; Netz-Nichtdeterminismus).
package httpcheck_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pt9912/d-check/internal/adapter/driven/httpcheck"
)

func adapter() *httpcheck.Adapter { return httpcheck.New(2 * time.Second) }

// DC-FA-EXT-001 Happy: Status < 400 → kein Fehlerzustand.
func TestCheck_Erreichbar(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	res := adapter().Check(srv.URL)
	if res.Status != 200 || res.Timeout || res.TooManyRedirects || res.TransportError != "" {
		t.Fatalf("res = %+v", res)
	}
}

// DC-FA-EXT-001 Negative: HTTP 404.
func TestCheck_Status404(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	if res := adapter().Check(srv.URL); res.Status != 404 {
		t.Fatalf("res = %+v", res)
	}
}

// §DC-FA-EXT-001.a: HEAD mit GET-Fallback bei 405/501.
func TestCheck_HeadFallbackAufGet(t *testing.T) {
	var methods []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	res := adapter().Check(srv.URL)
	if res.Status != 200 {
		t.Fatalf("res = %+v", res)
	}
	if len(methods) != 2 || methods[0] != http.MethodHead || methods[1] != http.MethodGet {
		t.Fatalf("methods = %v (HEAD, dann GET erwartet)", methods)
	}
}

// §DC-FA-EXT-001.a: Redirects bis 5 Stationen gefolgt; Station 6 →
// TooManyRedirects.
func TestCheck_RedirectKette(t *testing.T) {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	hop := func(n, last int) {
		mux.HandleFunc(fmt.Sprintf("/r%d", n), func(w http.ResponseWriter, r *http.Request) {
			if n >= last {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/r%d", n+1), http.StatusFound)
		})
	}
	for i := 0; i <= 7; i++ {
		hop(i, 5) // Kette endet bei /r5 mit 200 — genau 5 Redirects
	}
	if res := adapter().Check(srv.URL + "/r0"); res.Status != 200 || res.TooManyRedirects {
		t.Fatalf("5 Redirects müssen erlaubt sein: %+v", res)
	}

	mux2 := http.NewServeMux()
	srv2 := httptest.NewServer(mux2)
	defer srv2.Close()
	mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound) // Endlos-Kette
	})
	if res := adapter().Check(srv2.URL); !res.TooManyRedirects {
		t.Fatalf("Endlos-Kette: %+v (TooManyRedirects erwartet)", res)
	}
}

// DC-FA-EXT-001 Negative: Timeout.
func TestCheck_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	a := httpcheck.New(50 * time.Millisecond)
	if res := a.Check(srv.URL); !res.Timeout {
		t.Fatalf("res = %+v (Timeout erwartet)", res)
	}
}

// Transportfehler (Verbindung verweigert) → TransportError.
func TestCheck_Transportfehler(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close() // Port wieder frei → Verbindungsfehler
	res := adapter().Check(url)
	if res.TransportError == "" || res.Timeout || res.Status != 0 {
		t.Fatalf("res = %+v (TransportError erwartet)", res)
	}
}

// Nicht parsbare URL → TransportError ohne Netzzugriff.
func TestCheck_UngueltigeURL(t *testing.T) {
	res := adapter().Check("http://ungültig mit leerzeichen/")
	if res.TransportError == "" || res.Status != 0 {
		t.Fatalf("res = %+v (Parse-Fehler erwartet)", res)
	}
}
