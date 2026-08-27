package coretest

import (
	"errors"
	"fmt"
)

// ErrGleichLangeNeuschreibung ist der Befund von GitFixtureRewriteHazard — als
// Sentinel, damit der Negativ-Selbsttest ihn prüfen kann, statt auf einen
// Meldungstext zu vergleichen.
var ErrGleichLangeNeuschreibung = errors.New("fixture mit gleich langem Inhalt neu geschrieben")

// GitFixtureRewriteHazard meldet eine gleich lange Neuschreibung derselben
// git-Fixture-Datei. Beide Fixture-Helfer des Repos rufen sie — sonst gäbe es
// zwei Antworten auf dieselbe Frage.
//
// Zusage: die Metadaten-Abkürzung der Änderungserkennung überspringt das Hashen
// einer Datei nur, wenn Größe, Änderungszeit und Modus dem Index-Eintrag
// gleichen **und** die Datei älter ist als der Index. Die Größe ist davon die
// Bedingung, die ein Fixture selbst herstellt — unterscheiden sich die Längen,
// wird immer gehasht. Trifft die Abkürzung dagegen zu, sieht der folgende
// Commit den Baum als sauber („cannot create empty commit: clean working tree")
// und der Test wird zeitabhängig.
//
// Abgrenzung: gleiche Länge ist **notwendig, nicht hinreichend**. Scharf wird
// eine Stelle erst, wenn zwischen der zweiten Schreibung und dem nächsten
// Status ein index-schreibender Aufruf liegt. Die Prüfung ist deshalb
// konservativ und meldet auch Paare, die nie scharf waren.
//
// Kopplung: `oldLen` ist die Länge des **vorhandenen Datei-Inhalts**, `-1` wenn
// es keinen gibt. Das Lesen bleibt beim Aufrufer — dieser Kern hat kein
// Dateisystem, und die Entscheidung braucht keins.
//
// Grenze: gemessen wird der Inhalt **auf der Platte**, während die Gefahr am
// **Index-Eintrag** hängt. Eine ungestagte Zwischenschreibung verschiebt beides
// gegeneinander, und dieser Fall bleibt ungesehen — ebenso jeder Schreibweg,
// der an den Helfern vorbeiführt.
func GitFixtureRewriteHazard(name string, oldLen, newLen int) error {
	if oldLen < 0 || oldLen != newLen {
		return nil
	}
	return fmt.Errorf("%w: %q, %d Byte — die Metadaten-Abkürzung der "+
		"Änderungserkennung kann das übersehen; Länge unterscheiden",
		ErrGleichLangeNeuschreibung, name, newLen)
}
