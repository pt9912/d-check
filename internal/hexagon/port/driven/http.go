package driven

// HTTPResult ist das Ergebnis einer Erreichbarkeits-Prüfung
// (DC-FA-EXT-001): höchstens eine der Fehler-Bedingungen ist gesetzt;
// sonst trägt Status den letzten HTTP-Status.
type HTTPResult struct {
	// Status ist der letzte HTTP-Status (0, wenn keiner vorliegt).
	Status int
	// Timeout: Zeitlimit überschritten.
	Timeout bool
	// TooManyRedirects: Redirect-Kette länger als REDIRECT_MAX.
	TooManyRedirects bool
	// TransportError beschreibt DNS-/Verbindungsfehler ("" = keiner).
	TransportError string
}

// FetchResult ist das Ergebnis eines Inhalts-Abrufs des Moduls sources
// (DC-FA-SRC-001): bei Erfolg trägt Body die rohen Antwort-Bytes der finalen
// Antwort (Redirects gefolgt) und Status den HTTP-Status; sonst ist genau eine
// Fehler-Bedingung gesetzt. TooLarge markiert ein überschrittenes Body-Limit
// (> 64 MiB). Body ist auf das Limit begrenzt gelesen.
type FetchResult struct {
	// Status ist der HTTP-Status der finalen Antwort (0 bei Transportfehler).
	Status int
	// Body sind die gelesenen Antwort-Bytes (nil bei Fehler/TooLarge).
	Body []byte
	// Timeout: Zeitlimit überschritten.
	Timeout bool
	// TooManyRedirects: Redirect-Kette länger als REDIRECT_MAX.
	TooManyRedirects bool
	// TooLarge: Body-Limit überschritten.
	TooLarge bool
	// TransportError beschreibt DNS-/Verbindungsfehler ("" = keiner).
	TransportError string
}

// HTTPChecker ist der driven Port der Netz-Module: Check prüft die
// Erreichbarkeit einer URL (HEAD, Fallback GET — Modul external,
// spec/spezifikation.md §DC-FA-EXT-001.a); Fetch holt den Inhalt (GET,
// Redirects gefolgt, größenbegrenzt — Modul sources, §DC-FA-SRC-001.a). Der
// HTTP-Adapter ist die einzige Netzwerk-Tür von d-check (DC-QA-03; beide
// Module strikt opt-in).
type HTTPChecker interface {
	Check(url string) HTTPResult
	Fetch(url string) FetchResult
}
