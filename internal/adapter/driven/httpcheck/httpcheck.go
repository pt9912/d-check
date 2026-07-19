// Package httpcheck ist der driven HTTP-Adapter (Modul external):
// HEAD mit GET-Fallback bei 405/501, Timeout, Redirect-Limit
// (spec/spezifikation.md §DC-FA-EXT-001.a; spec/architecture.md §2).
// Einziger net/http-Ort des Repos (arch-check R2) — und damit die
// einzige Netzwerk-Tür (DC-QA-03).
package httpcheck

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pt9912/d-check/internal/hexagon/port/driven"
)

// redirectMax ist REDIRECT_MAX (spec/spezifikation.md §3, fest).
const redirectMax = 5

// maxFetchBytes ist das Body-Limit des Moduls sources (≤ 64 MiB,
// spec/spezifikation.md §DC-FA-SRC-001.a Schritt 3).
const maxFetchBytes = 64 << 20

// errTooManyRedirects markiert eine Redirect-Kette über REDIRECT_MAX.
var errTooManyRedirects = errors.New("zu viele Redirects")

// Adapter implementiert driven.HTTPChecker über net/http.
type Adapter struct {
	client *http.Client
}

// New erzeugt einen Adapter mit dem gegebenen Gesamt-Timeout pro
// Prüfung (DC-FA-EXT-001: konfigurierbar, Default 10 s).
func New(timeout time.Duration) *Adapter {
	return &Adapter{client: &http.Client{
		Timeout: timeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) > redirectMax {
				return errTooManyRedirects
			}
			return nil
		},
	}}
}

// Check prüft die Erreichbarkeit: HEAD, bei 405/501 Fallback auf GET
// (Body verworfen).
func (a *Adapter) Check(url string) driven.HTTPResult {
	res := a.request(http.MethodHead, url)
	if res.Status == http.StatusMethodNotAllowed || res.Status == http.StatusNotImplemented {
		return a.request(http.MethodGet, url)
	}
	return res
}

func (a *Adapter) request(method, url string) driven.HTTPResult {
	req, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	if err != nil {
		return driven.HTTPResult{TransportError: err.Error()}
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return classifyError(err)
	}
	defer func() { _ = resp.Body.Close() }()
	// Body verworfen (GET-Fallback) — Drain begrenzt: kleine Bodies
	// erhalten Keep-Alive, große werden nicht voll heruntergeladen.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	return driven.HTTPResult{Status: resp.StatusCode}
}

// Fetch holt den Inhalt einer URL (GET; Redirects bis REDIRECT_MAX über den
// geteilten Client wie Check; Body ≤ maxFetchBytes) — der Netz-Port des Moduls
// sources (spec/spezifikation.md §DC-FA-SRC-001.a Schritt 3). Ein Body über dem
// Limit ⇒ TooLarge (kein Hash-Schritt). Transportfehler werden über classifyError
// klassifiziert (geteilte Fehler-Zuordnung mit Check).
func (a *Adapter) Fetch(url string) driven.FetchResult {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return driven.FetchResult{TransportError: err.Error()}
	}
	resp, err := a.client.Do(req)
	if err != nil {
		c := classifyError(err)
		return driven.FetchResult{
			Timeout:          c.Timeout,
			TooManyRedirects: c.TooManyRedirects,
			TransportError:   c.TransportError,
		}
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFetchBytes+1))
	if err != nil {
		return driven.FetchResult{Status: resp.StatusCode, TransportError: err.Error()}
	}
	if int64(len(body)) > maxFetchBytes {
		return driven.FetchResult{Status: resp.StatusCode, TooLarge: true}
	}
	return driven.FetchResult{Status: resp.StatusCode, Body: body}
}

// classifyError ordnet Transport-Fehler den Port-Bedingungen zu.
func classifyError(err error) driven.HTTPResult {
	if errors.Is(err, errTooManyRedirects) {
		return driven.HTTPResult{TooManyRedirects: true}
	}
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) ||
		(errors.As(err, &netErr) && netErr.Timeout()) {
		return driven.HTTPResult{Timeout: true}
	}
	return driven.HTTPResult{TransportError: err.Error()}
}
