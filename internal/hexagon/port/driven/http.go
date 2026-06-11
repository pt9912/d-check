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

// HTTPChecker ist der driven Port des Moduls external: prüft die
// Erreichbarkeit einer URL (HEAD, Fallback GET — spec/spezifikation.md
// §DC-FA-EXT-001.a). Der HTTP-Adapter ist die einzige Netzwerk-Tür
// von d-check (DC-QA-03).
type HTTPChecker interface {
	Check(url string) HTTPResult
}
