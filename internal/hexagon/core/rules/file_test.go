package rules

import (
	"strings"
	"testing"

	"github.com/pt9912/d-check/internal/hexagon/core/coretest"
	"github.com/pt9912/d-check/internal/hexagon/core/model"
)

// nLines baut Inhalt mit genau n Zeilen (n Zeilenumbrüche) — die Zählung von
// countLines: n Zeilenumbrüche, endet auf einen, also n Zeilen.
func nLines(n int) string {
	return strings.Repeat("x\n", n)
}

// Der Grenzwert selbst — N ist grün, N+1 ist rot. Kein Vorzustand nötig: die
// Bedingung existiert erst mit diesem Slice, es gibt keine Regression zu
// belegen.
func TestCheckFile_MaxLinesGrenzwert(t *testing.T) {
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(5)}
	fsGruen := coretest.NewMemFS(map[string]string{"a.txt": nLines(5)})
	if f := CheckFile(fsGruen, []model.FileRule{r}); f != nil {
		t.Fatalf("5 Zeilen bei max-lines 5 ⇒ befundfrei, got %+v", f)
	}
	fsRot := coretest.NewMemFS(map[string]string{"a.txt": nLines(6)})
	f := CheckFile(fsRot, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileLinesExceeded {
		t.Fatalf("6 Zeilen bei max-lines 5 muss file-lines-exceeded melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "6 Zeilen") || !strings.Contains(f[0].Message, "5") {
		t.Errorf("Meldung nennt Ist- und Soll-Zahl: %q", f[0].Message)
	}
}

// Eine unvollstaendige Schlusszeile zaehlt eine weitere Zeile mit — dieselbe
// Zaehlung, die codepaths/citations bereits teilen (countLines).
func TestCheckFile_MaxLinesUnvollstaendigeSchlusszeile(t *testing.T) {
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(5)}
	fs := coretest.NewMemFS(map[string]string{"a.txt": nLines(5) + "rest ohne umbruch"})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileLinesExceeded {
		t.Fatalf("5 Umbrueche + unvollstaendige Schlusszeile = 6 Zeilen, muss melden: %+v", f)
	}
}

// Bytes werden ROH gezaehlt -- Fenced-Code und Inline-Code zaehlen mit. Das
// ist die Eigenschaft, die den Umweg ueber structure.forbid-pattern (zaehlt
// den BEREINIGTEN Text) unbrauchbar macht (ADR-0088).
func TestCheckFile_MaxBytesZaehltFencedCodeMit(t *testing.T) {
	body := "# T\n\n```\nviel code hier drin\n```\n"
	r := model.FileRule{Files: "a.md", MaxBytes: ptr(len(body) - 1)}
	fs := coretest.NewMemFS(map[string]string{"a.md": body})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileBytesExceeded {
		t.Fatalf("max-bytes knapp unter der rohen Groesse muss trotz Fence melden: %+v", f)
	}
}

// Nullmengen-Haerte: eine Regel, die keine Datei trifft (auch nach Abzug von
// exempt-paths), meldet statt leer zu laufen.
func TestCheckFile_NullmengenHaerte(t *testing.T) {
	r := model.FileRule{Files: "nirgends/*.md", MaxLines: ptr(10)}
	fs := coretest.NewMemFS(map[string]string{"a.md": "x\n"})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileNoMatch {
		t.Fatalf("Glob ohne Treffer muss file-no-match melden, got %+v", f)
	}
}

// exempt-paths kann die Kandidatenmenge bis auf null leeren -- dieselbe
// Nullmengen-Haerte greift dann, kein stiller Leerlauf.
func TestCheckFile_ExemptPathsLeertBisAufNull(t *testing.T) {
	r := model.FileRule{Files: "docs/*.md", MaxLines: ptr(1), ExemptPaths: []string{"docs/a.md"}}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": nLines(5)})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileNoMatch {
		t.Fatalf("die einzige Kandidatin ist ausgenommen ⇒ file-no-match, got %+v", f)
	}
}

// Eine unlesbare Einzeldatei ist FAIL-CLOSED wie bei structure -- kein
// stiller Uebersprung.
func TestCheckFile_UnlesbareDateiFailClosed(t *testing.T) {
	fs := readErrFS{coretest.NewMemFS(map[string]string{"a.txt": nLines(3)})}
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(1)}
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileNoMatch {
		t.Fatalf("unlesbare Datei muss fail-closed melden, got %+v", f)
	}
	if !strings.Contains(f[0].Message, "unlesbar") {
		t.Errorf("die Meldung muss die Ursache nennen: %q", f[0].Message)
	}
}

// hint gewinnt gegen die modul-eigene Meldung eines SCHWELLEN-Befunds --
// dieselbe Form wie structure[].hint.
func TestCheckFile_HintGewinntBeiSchwellenBefund(t *testing.T) {
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(1), Hint: "bitte kuerzen"}
	fs := coretest.NewMemFS(map[string]string{"a.txt": nLines(2)})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Message != "bitte kuerzen" {
		t.Fatalf("hint muss die Meldung ersetzen, got %+v", f)
	}
}

// hint gewinnt NICHT bei file-no-match -- dort hat die Regel nicht gemessen
// (dieselbe Grenze wie structure.MessageFor).
func TestCheckFile_HintGiltNichtBeiNoMatch(t *testing.T) {
	r := model.FileRule{Files: "nirgends/*.md", MaxLines: ptr(1), Hint: "bitte kuerzen"}
	fs := coretest.NewMemFS(map[string]string{"a.md": "x\n"})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Message == "bitte kuerzen" {
		t.Fatalf("file-no-match darf den hint nicht tragen, got %+v", f)
	}
}

// Zwei unabhaengige Schwellen in EINER Regel: beide koennen gleichzeitig
// verletzt sein und melden dann beide.
func TestCheckFile_BeideSchwellenGleichzeitig(t *testing.T) {
	body := nLines(10)
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(5), MaxBytes: ptr(len(body) - 1)}
	fs := coretest.NewMemFS(map[string]string{"a.txt": body})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 2 {
		t.Fatalf("beide Schwellen verletzt ⇒ zwei Befunde, got %+v", f)
	}
}

// Modul-aus: ohne Regeln ist der Befundsatz byte-identisch (nil) --
// DC-QA-02.
func TestCheckFile_OhneRegelnByteIdentisch(t *testing.T) {
	fs := coretest.NewMemFS(map[string]string{"a.txt": nLines(999)})
	if f := CheckFile(fs, nil); f != nil {
		t.Fatalf("ohne Regeln muss der Befundsatz leer bleiben, got %+v", f)
	}
}

// Die Kandidaten-Menge ist NICHT auf Markdown beschraenkt -- anders als
// structure. Eine .sh-Datei ist ein gueltiger Kandidat.
func TestCheckFile_NichtAufMarkdownBeschraenkt(t *testing.T) {
	r := model.FileRule{Files: "tools/*.sh", MaxLines: ptr(2)}
	fs := coretest.NewMemFS(map[string]string{"tools/a.sh": nLines(3)})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileLinesExceeded {
		t.Fatalf("eine Shell-Datei muss wie jede andere geprueft werden, got %+v", f)
	}
}

// Eine leere Datei hat 0 Zeilen (countLines("") == 0) und 0 Bytes -- unter
// jeder positiven Schwelle befundfrei.
func TestCheckFile_LeereDatei(t *testing.T) {
	r := model.FileRule{Files: "a.txt", MaxLines: ptr(0), MaxBytes: ptr(0)}
	fs := coretest.NewMemFS(map[string]string{"a.txt": ""})
	if f := CheckFile(fs, []model.FileRule{r}); f != nil {
		t.Fatalf("leere Datei bei max-lines/max-bytes 0 ⇒ befundfrei, got %+v", f)
	}
}

// Bytes werden ROH gezaehlt, nicht als Runen/Zeichen -- ein Umlaut ist im
// UTF-8-Quelltext zwei Bytes. Eine Schwelle, die die Byte-Zahl korrekt
// zwischen Runen- und Byte-Zaehlung unterscheidet, muss bei drei Umlauten
// (3 Runen, 6 Bytes) rot werden, waehrend eine Runen-Zaehlung grün bliebe.
func TestCheckFile_BytesZaehltRohNichtRunen(t *testing.T) {
	body := "ääält " // drei Umlaute: 3 Runen, aber 6 Bytes; plus vier ASCII-Bytes
	if len([]rune(body)) >= len(body) {
		t.Fatalf("Testvoraussetzung verletzt: Fixture muss mehr Bytes als Runen haben")
	}
	r := model.FileRule{Files: "a.txt", MaxBytes: ptr(len([]rune(body)))}
	fs := coretest.NewMemFS(map[string]string{"a.txt": body})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Reason != ReasonFileBytesExceeded {
		t.Fatalf("max-bytes auf Runenzahl gesetzt muss an der hoeheren Byte-Zahl rot werden, got %+v", f)
	}
}

// Ist schon der Dateibaum nicht lesbar, meldet JEDE Regel -- ein Befund je
// Regel, kein Sammel-Befund, damit die Identitaet der unmessbaren Behauptung
// erhalten bleibt (wie bei structure, TestStructureUnlesbarerBaumFailClosed).
func TestCheckFile_DateibaumUnlesbarJeRegelEinBefund(t *testing.T) {
	fs := listErrFS{coretest.NewMemFS(map[string]string{"a.txt": nLines(1)})}
	rules := []model.FileRule{
		{Files: "a.txt", MaxLines: ptr(1)},
		{Files: "b.txt", MaxLines: ptr(1)},
	}
	f := CheckFile(fs, rules)
	if len(f) != len(rules) {
		t.Fatalf("erwartet je Regel einen Befund, got %+v", f)
	}
	for _, fi := range f {
		if fi.Reason != ReasonFileNoMatch || !strings.Contains(fi.Message, "fail-closed") {
			t.Errorf("erwartet fail-closed-Befund, got %+v", fi)
		}
	}
	if f[0].Target == f[1].Target {
		t.Errorf("die Regel-Identitaet muss erhalten bleiben: %q", f[0].Target)
	}
}

// Grenzwert-Symmetrie fuer max-bytes, eigenstaendig wie bei max-lines: N ist
// gruen, N+1 ist rot.
func TestCheckFile_MaxBytesGrenzwert(t *testing.T) {
	body := nLines(3) // 6 Bytes
	r := model.FileRule{Files: "a.txt", MaxBytes: ptr(len(body))}
	if f := CheckFile(coretest.NewMemFS(map[string]string{"a.txt": body}), []model.FileRule{r}); f != nil {
		t.Fatalf("Byte-Zahl == max-bytes ⇒ befundfrei, got %+v", f)
	}
	rEins := model.FileRule{Files: "a.txt", MaxBytes: ptr(len(body) - 1)}
	f := CheckFile(coretest.NewMemFS(map[string]string{"a.txt": body}), []model.FileRule{rEins})
	if len(f) != 1 || f[0].Reason != ReasonFileBytesExceeded {
		t.Fatalf("Byte-Zahl > max-bytes muss file-bytes-exceeded melden, got %+v", f)
	}
}

// Regel-Identitaet ist der files-Glob; zwei Regeln mit demselben Glob sind
// eine Konfigurations-Dopplung. Geprueft ist hier nur die Kern-Funktion
// (CheckFile prueft keine Dopplung -- das ist Sache des Config-Adapters);
// dieser Test belegt nur, dass die Identitaet im Befund landet.
func TestCheckFile_IdentityImTarget(t *testing.T) {
	r := model.FileRule{Files: "docs/*.md", MaxLines: ptr(1)}
	fs := coretest.NewMemFS(map[string]string{"docs/a.md": nLines(2)})
	f := CheckFile(fs, []model.FileRule{r})
	if len(f) != 1 || f[0].Target != "docs/*.md" {
		t.Fatalf("target muss die Regel-Identitaet (der Glob) sein, got %+v", f)
	}
}
