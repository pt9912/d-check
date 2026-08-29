# Slice slice-177: Eine `structure`-Regel sagt, was zu tun ist

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**. Sein Closure-Grund geht über die eigene DoD nicht
hinaus (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht); der **Anlass** liegt in
[welle-86](../welle-86-closure-uebergang-durchsetzen.md), die **Fähigkeit**
gehört jeder `structure`-Regel.

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(das erweiterte Modul);
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
(die Befund-Zeile ist ein Vertrag);
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(die Diagnose und ihre `fixCandidate`-Abgrenzung);
[ADR-0069](../../adr/0069-zellenlaenge-als-strukturbedingung.md) (dieselbe
Bauform: eine neue `structure`-Bedingung samt Lastenheft-Bump).

**Berührte Spec-Stellen:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
sowie ihre `.a`-Verfeinerungen in der Spezifikation.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Ein Befund sagt heute, was falsch ist, aber nicht, was der Regel-Autor
wollte.**

`section-forbidden` heißt „hier steht eine verbotene Wendung". Welche Zusage
diese Regel hütet und was der Leser jetzt tun soll, steht nirgends im Lauf — es
steht bestenfalls als Kommentar in der Konfiguration, also an einem Ort, den
niemand öffnet, wenn ein Gate rot wird.

**Der Grund-Code kann das nicht tragen, und zwar aus einem Struktur-Grund.**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
führt zu jeder Bedingung eine **Reparatur**-Spalte — für `forbid-pattern` lautet
sie *„die Wendung ersetzen"*. Das ist richtig und bleibt richtig: der Grund-Code
sagt die **Art** des Defekts. Aber **eine** Art bedient viele Regeln. Jede
`forbid`-Regel dieses Repos bekäme denselben Satz, obwohl sie verschiedene
Zusagen hüten — die eine einen erfundenen Risiko-Ausgang, die andere einen
offenen DoD-Haken.

**Die fehlende Größe ist regel-lokal, nicht code-lokal.** Sie gehört an die
Regel, weil nur die Regel weiß, welche Zusage sie hütet.

**Der Anlass ist gemessen und liegt in diesem Repo:**
[slice-172](../open/slice-172-closure-uebergang-waechtern.md) hat den Sensor für
den offenen DoD-Haken fertig entworfen und beidseitig gemessen (37 Befunde ohne
Bestands-Ausnahme, null mit ihr, Positiv-Probe grün) — und ist an genau diesem
Punkt zurückgeführt worden: seine eigene §2 verlangt eine Meldung, die sagt,
*was zu tun ist*, und das Schema kennt sie nicht.

## 2. Vorgehen

1. **Ein neues Feld `structure[].hint`**: freier Text an **einer** Regel. Leer
   gesetzt ⇒ Exit 2 — ein leerer Hinweis sagt nichts zu, dieselbe Härte, die
   `planning.closure.boilerplate` für den leeren Eintrag führt.
2. **Die Befund-Zeile bekommt eine vierte Spalte, und nur wenn sie gefüllt
   ist.** [`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate)
   sagt **ein Befund pro Zeile** zu; eine Fortsetzungszeile bräche das. Eine
   vierte tab-getrennte Spalte bricht es nicht — wer auf Tab trennt und die
   Felder 1–3 liest, liest weiter dasselbe. **Nur wenn gefüllt**, damit die
   Ausgabe jeder Regel ohne `hint` **byte-identisch** bleibt; das ist zu
   messen, nicht zu behaupten.
3. **`--json`/`--yaml`** tragen `hint` je `findings`-Eintrag (additiv);
   **`--doctor`** rendert ihn als eigene Zeile unter `Stelle:`.
4. **Abgrenzung gegen `fixCandidate`, ausdrücklich.** Der Fix-Kandidat ist
   **abgeleitet** und nur dort, wo er *eindeutig ableitbar* ist
   ([`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus));
   er speist `--repair` und wird zu einem anwendbaren Patch. Ein `hint` ist
   **verfasst**. Beides zu vermischen hieße, Autoren-Prosa in die
   Patch-Pipeline zu geben — das ist der Grund für ein eigenes Feld und nicht
   für `fixCandidate.note`.
5. **ADR** für die Entscheidung (vierte Spalte statt Fortsetzungszeile;
   regel-lokal statt code-lokal; `hint` ≠ `fixCandidate`), Lastenheft-Bump,
   Spezifikations-Verfeinerung, Handbuch.
6. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein `hint` für andere Module.** `links`, `ids`, `matrix` und die übrigen
  konfigurieren keine Regel-Liste mit Autoren-Text; ihre Befunde entstehen aus
  dem Dokument, nicht aus einer benannten Zusage. Wer das später will, misst
  erst den Bedarf.
- **Keine Änderung an `fixCandidate` oder `--repair`.** Siehe §2 Punkt 4.
- **Keine Prüfung des Hinweis-Textes.** Ob er stimmt, ist Urteil — dieselbe
  Klasse wie ein Kommentar (`AGENTS.md` §3.7). Der Sensor prüft, dass er nicht
  leer ist.
- **Keine Anwendung auf den DoD-Haken-Sensor.** Das ist
  [slice-172](../open/slice-172-closure-uebergang-waechtern.md); dieser Slice
  liefert die Fähigkeit, nicht ihren ersten Konsumenten.

## 4. Definition of Done

- [ ] `structure[].hint` ist im Schema, in
      [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Bump + Historie)
      und in [`spec/spezifikation.md`](../../../../spec/spezifikation.md)
      geführt; leerer Wert ⇒ Exit 2, mit Test.
- [ ] **Byte-Identität gemessen:** ein voller `make gates`-Lauf vor und nach der
      Änderung liefert für alle Regeln **ohne** `hint` dieselbe Ausgabe; die
      Messung steht in der Commit-Botschaft, nicht als Behauptung.
- [ ] Die vierte Spalte erscheint in der Befund-Zeile, `hint` in `--json` und
      `--yaml`, eine eigene Zeile in `--doctor` — je mit Test.
- [ ] Eine ADR begründet die drei Entscheide aus §2 Punkt 5 und ist im
      [ADR-Index](../../adr/README.md) eingetragen.
- [ ] Das [Benutzerhandbuch](../../../user/benutzerhandbuch.md) führt das Feld
      dort, wo es die übrigen `structure`-Schlüssel führt.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein Hinweis ist Autoren-Text und altert wie ein Kommentar.** Er kann
  unwahr werden, ohne dass ein Gate es merkt — die Klasse, die dieses Repo als
  [`BEO-013`](../observations.md) führt, nur ohne Wächter. —
  **Ausgang:** *(bei Closure)*
- **Die Ausgabe-Zusage wird geweitet.** Ein Konsument, der auf genau drei
  Tab-Felder besteht, sieht ab jetzt bei manchen Regeln vier. Die
  Nicht-Änderung für Regeln ohne `hint` ist messbar, die Weitung selbst bleibt
  eine Vertrags-Änderung. — **Ausgang:** *(bei Closure)*
- **Das Feld lädt zur Ausrede ein.** Ein schlecht benannter Grund-Code lässt
  sich mit einem Hinweis zudecken, statt den Code zu schärfen. —
  **Ausgang:** *(bei Closure)*
- **Die Fähigkeit entsteht für einen einzigen Konsumenten**
  ([`BEO-011`](../observations.md)): slice-172. Ob weitere Regeln sie nutzen,
  ist nicht gemessen. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): unmittelbar — der Slice entsteht als
Vorbedingung von [slice-172](../open/slice-172-closure-uebergang-waechtern.md),
der dafür nach `open/` zurückgeführt wurde.

**Rückführungen:** `in-progress` → `open`, falls sich beim Bau zeigt, dass die
vierte Spalte doch Konsumenten bricht, die das Repo führt — dann ist der
Entscheid ein anderer (`--doctor`-only), und die ADR trägt ihn.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `internal/hexagon/core/` (Kern: Modell und Regel) und
  `internal/adapter/driven/` (Report-Rendering). Beide fallen unter den Default
  `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-022`): [`BEO-011`](../observations.md) — Regel aus dem Anlass: die
  Fähigkeit entsteht für **einen** Konsumenten, und das steht als Risiko in §5;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt: ein
  alternder Hinweis ist die textliche Spielart derselben Klasse;
  [`BEO-022`](../observations.md) — eine Regel tritt in Kraft, bevor ihre
  Zustellung existiert: hier gering, weil ein neues **optionales** Feld
  niemanden blockiert; die Handbuch-Zeile ist trotzdem DoD-Punkt.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).

Slice-ID: slice-177. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-CLI-004`](../../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate),
[`DC-FA-CLI-007`](../../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus).
Module: `structure`. Gates: `make gates`, `make test`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Ein neues optionales Konfigurationsfeld plus
sein Rendering; kein Fremdsystem, keine Reconciliation, kein Bestand, der
umgestellt werden müsste.

## 9. Closure-Notiz (nach `done/`)
