# Slice slice-025: Diagnose-Modus (`--doctor`)

**Status:** open (geplant).

**Welle:** welle-15-doctor-repair (Trigger: Change Request 0.15.0
akzeptiert; erste Slice der Welle, keine Vorbedingung).

**Bezug:**
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(Hauptanforderung),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)
(`ids` muss aktiv sein, damit `id-unlinked`-Kandidaten entstehen),
[`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
(Exit-Codes der Modi),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(Modus ersetzt Default-stdout; nicht mit `--json` kombinierbar),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)
(byte-identische Diagnose),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(read-only, stdout-only).

**Autor:** pt9912. **Datum:** 2026-06-18.

---

## 1. Ziel

`d-check --doctor` macht einen Lese-Durchgang wie eine normale Prüfung,
gibt aber statt der knappen Befund-Zeilen eine **erklärende, nach Datei
und Regel gruppierte Diagnose** auf stdout aus: je Befund den Grund-Code
in Klartext und — wo aus dem Befund **eindeutig ableitbar** — einen
**Fix-Kandidaten** (vorgeschlagene Änderung, **nicht angewendet**).
Read-only, stdout-only.

Diese Slice führt das **Fix-Kandidaten-Modell** als eigene Core-Funktion
ein; slice-026 (`--repair`) rendert dieselben Kandidaten zum Patch — eine
Quelle, zwei Ausgaben.

## 2. Definition of Done

- [ ] **Spezifikation** zu [`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus) (Spec-Abschnitt `…007.a` entsteht in dieser Slice): Diagnose-Format (Gruppierung
  nach Datei/Regel, Grund-Klartext-Mapping je Grund-Code), Fix-Kandidaten-
  Modell (welche Grund-Codes einen *eindeutigen* Kandidaten liefern — v1:
  `id-unlinked` → Definitions-Link), Determinismus-Festlegung,
  `--json`-Inkompatibilität (→ exit 2).
- [ ] **Grund-Klartext-Mapping** für alle aktiven Grund-Codes
  (Stand 14: `target-missing`, `repo-escape`, `symlink`, `anchor-missing`,
  `id-unlinked`, `codepath-missing`, `matrix-inactive`, `matrix-forbidden`,
  `external-redirects`, `external-status`, `external-timeout`,
  `span-unclosed`, `span-nested-link`, `hostpath-forbidden`) — mit
  Vollständigkeits-Sicherung (Test/Gate gegen neue, unkartierte Codes).
- [ ] **Implementierung** im CLI-/Core-Layer: Flag `--doctor`,
  Diagnose-Renderer, Fix-Kandidaten-Ableitung als wiederverwendbare
  Core-Funktion (Eingang für slice-026).
- [ ] **Default-Ausgabe ersetzt** (nicht ergänzt): `--doctor` gibt die
  Diagnose *statt* der Default-Befund-Zeilen
  ([`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate))
  auf stdout aus.
- [ ] **Exit-Codes** ([`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)):
  0/1/2; Negativtest `--doctor --json` → exit 2.
- [ ] **Read-only-Beleg** ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  und **Determinismus-Beleg** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
  10× byte-identisch).
- [ ] **Doku** unter `docs/user/`; `make gates` grün; Closure-Notiz.

## 3. Plan (vor Code)

| Datei | Art | Begründung |
|---|---|---|
| [`spec/spezifikation.md`](../../../../spec/spezifikation.md) | update | Spec-Abschnitt zu [`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus): Diagnose-Format, Fix-Kandidaten-Modell, Determinismus |
| `internal/…` (CLI + Diagnose-Renderer + Kandidaten-Funktion) | update/neu | Flag, Gruppierung, Grund-Klartext, Kandidaten-Ableitung |
| `docs/user/operations.md` | update | Option `--doctor` dokumentieren |

Das Lastenheft ([`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)) ist mit Change Request 0.15.0 bereits
gesetzt — kein weiterer Vertrags-Change in dieser Slice (außer einer
Schärfung, falls die Spezifikation eine aufdeckt).

## 4. Trigger

Change Request 0.15.0 akzeptiert (committet). Keine Slice-Vorbedingung
(erste Slice der Welle).

## 5. Closure-Trigger

DoD vollständig inkl. Grund-Klartext-Vollständigkeit, Determinismus-/
read-only-Beleg und grüner Gates.

## 6. Risiken und offene Punkte

- **Grund-Klartext-Vollständigkeit:** jeder künftige Grund-Code braucht
  einen Mapping-Eintrag — sonst stille Lücke. Mit einem Test/Gate gegen
  die Grund-Code-Konstanten absichern, nicht per Hand pflegen.
- **Fix-Kandidaten-Eindeutigkeit:** nur `id-unlinked` ist v1 sicher
  ableitbar (Definition über die `ids`-Mechanik bekannt). Nicht
  überdehnen — `target-missing`/`span-*` sind Best-Guess und gehören in
  die breite `--repair`-Stufe (slice-026), nicht in die Kandidaten der
  Diagnose als „eindeutig".
- **Determinismus** ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)):
  Gruppierung darf keine Map-Iterationsreihenfolge nach außen lecken —
  stabil sortieren wie der bestehende Befund-Pfad.

## 7. Closure-Notiz (nach `done/`)

**Umsetzung:** Vertrag
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(Lastenheft 0.15.0) + Spezifikation
[`DC-FA-CLI-007.a`](../../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus);
`core.FixCandidateFor`/`ReasonText`/`AllReasons` (Fix-Kandidaten-Modell +
Grund-Klartext, alle als Funktionen statt Paket-Global), `report.Doctor`
(gruppierter Renderer), CLI-Flag `--doctor` mit `--json`-Inkompatibilität.
`make gates` grün (Coverage 94,5 %).

**Belege:**

- Happy/Boundary/Negative + Determinismus (10×) als CLI-Akzeptanztests;
  Core-Tests für Klartext-Vollständigkeit und Kandidaten-Ableitung.
- Default-Format **ersetzt** (Test: keine knappe Befund-Zeile unter
  `--doctor`).
- read-only ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit))
  und Determinismus ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).

**Fix-Kandidaten-Modell:** v1 nur `id-unlinked` → Link auf das
Definitions-`target` (Datei-Ebene). Bewusst konservativ — Best-Guess-Fälle
liefern keinen Kandidaten und gehören in die breite `--repair`-Stufe. Die
Ableitung (`FixCandidateFor`) ist die wiederverwendbare Eingabe für
slice-026: eine Quelle, zwei Ausgaben.

**Lerneintrag:** `FixCandidateFor` als Core-Funktion trennt das *Was*
(eindeutiger Fix) vom *Wie* der Ausgabe (Diagnose hier, Patch in
slice-026). Die Vollständigkeits-Prüfung gegen die Reason-Konstanten
verwandelt eine stille Lücke (neuer Grund-Code ohne Klartext) in einen
Test-Bruch. Der Verzicht auf das Paket-Global (Map/Slice → Funktion) hielt
`gochecknoglobals` ohne Lint-Ausnahme grün.

**Review R1** (Self-Review,
[Report](../../../reviews/2026-06-18-slice-025-doctor.md)): HIGH 0 /
MEDIUM 0 / LOW 1 / INFO 4 — freigegeben. LOW-1 (Map-pro-Aufruf) als
bewusster Tausch gegen ein Global akzeptiert; INFO-1 (Verzeichnis-Target-
Kandidat nur Datei-Ebene) als Forward-Note an slice-026 übergeben.

**Folge-Slice:** slice-026 (`--repair`) — rendert die Fix-Kandidaten zum
Patch; Trigger „slice-025 done" ist eingetreten.

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (Spec-/Code-/Doku-Arbeit; Greenfield-Default).
