# Slice slice-290: Modulkern `linelength`: Konfiguration lesen, Zeilen messen, Ausnahmen anwenden, Grund-Code

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-95-linelength-modul`](../welle-95-linelength-modul.md).

**Bezug:** [`DC-FA-LLEN-001`](../../../../spec/lastenheft.md#dc-fa-llen-001--zeilenlänge-in-markdown-dateien-modul-linelength-opt-in) und
[ADR-0114](../../adr/0114-linelength-modul.md) aus
[`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl).

**Berührte Spec-Stellen:** `spezifikation.md` §4 (die Zeile des Grund-Codes),
`lastenheft.md` (Modul-Aufzählung in
[`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl) und
Glossar). Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-11-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das opt-in-Modul `linelength` ist registriert, liest seinen Block am
Konfigurations-Rand (Fehlkonfiguration ⇒ Exit 2), misst jede Zeile der Markdown-Dateien
seiner Scan-Menge nach den in [`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md)
bestätigten Entscheiden und meldet jede überlange Zeile mit ihrem Grund-Code — jede
Randbedingung der Anforderung hat einen Test, der ohne die Regel aus dem richtigen
Grund rot läuft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anforderung, Verfahren, ADR ändern** — gehören
  [`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md). Zeigt der Code, dass ein
  Entscheid nicht trägt, geht das als Plan-Änderung zurück an die Anforderung, nicht als
  stille Abweichung im Code.
- **Das Modul im eigenen Profil (`.d-check.yml`) einschalten** — Welle-Out-of-Scope:
  der Bestand führt bewusst lange Zeilen; das Einschalten wäre eine Umformatierung,
  keine Lieferung dieses Slice.
- **Eine eigene Fence- oder Tabellen-Erkennung** — die Ausnahmen kommen aus der
  geteilten Lexik; ein zweites Muster daneben ist genau die Klasse
  [`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md).
- **`--fix-patch`, `--doctor`-Reparaturvorschlag, `--suggest-config`-Erkennung** — eigene
  Vorgänge; der Diagnose-Klartext des Grund-Codes kommt mit (Lockstep-Test), mehr nicht.
- **CHANGELOG, Handbuch-Kopf, Release** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). Was ein Spiegel-Test schon im
  Feature-Commit verlangt (Modul-Bullets der READMEs, Modul-Abschnitt und
  Modul-Aufzählung im Handbuch), kommt mit und steht in §3.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Konfigurations-Rand: der Block `linelength` wird gelesen, jede Fehlkonfiguration der
      Anforderung (fehlende oder < 1 Schwelle, unbekannter Schlüssel, ungültiger
      `scope`/`exempt-paths`) ist Exit 2 mit Test; `--print-config` und `--enable` kennen
      das Modul.
- [ ] Regel: Zählung in Codepunkten, Grenze genau auf der Schwelle befundfrei, jede
      Ausnahme einzeln schaltbar und getestet (an und aus), ein Befund je Zeile mit Ort,
      Länge und Grenze, Reihenfolge nach Pfad und Zeile, Bezugsmenge (Zeilen wie bei `file`
      gezählt, die über (c) ausgenommenen als eigene Zahl), leere Menge kein Fehler; die
      Negativfälle aus dem Plan von
      [`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md) stehen als Testfälle, und
      für jeden Test ist ein Rot ohne die Regel gezeigt (Baseline-Regelwerk
      `modul-11-verification.md` §Bewusstes Brechen). Der Kopplungs-Test fährt dieselbe Eingabe
      durch `linelength` und die bisherigen Konsumenten von `proseLineSet`/`tableRowLine` —
      Pflicht nach [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md)
      Entscheidung 4, keine Option (Architect-Verdikt F-3).
- [ ] Spiegel: Grund-Code in `AllReasons()`/Diagnose-Klartext und §4 der Spezifikation,
      Modul in der Aufzählung von
      [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl),
      im Glossar, in den READMEs und im Handbuch — die Spiegel-Deckungstests sind grün und
      fangen das neue Modul (vor dem Editieren aufgelistet,
      [`MR-025`](../../../../harness/conventions.md#mr-025)).
- [ ] `make gates` grün.
- [ ] [ADR-0114](../../adr/0114-linelength-modul.md) `Accepted`, Index nachgezogen — der
      Status-Übergang durch den Architect, nachdem Review und Verifikation dieses Slice keine
      Änderung an der Entscheidung mehr verlangen (Architect-Verdikt F-8).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/core/model/config.go` | update | Konfigurations-Typ des Moduls; Modulliste (`ValidModules()`) |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Block lesen und am Rand prüfen (Exit 2) |
| Regel-Datei `linelength` im Verzeichnis `internal/hexagon/core/rules/` | neu | Zeilen messen, Ausnahmen über `proseLineSet`/`tableRowLine`, Befunde, Bezugsmenge |
| `internal/hexagon/core/rules/run.go` | update | Modul im Lauf verdrahten (Scan-Menge mit `scope`/`exempt-paths`) |
| `internal/hexagon/core/app/diagnose.go` | update | `AllReasons()` und Diagnose-Klartext des Grund-Codes (Lockstep-Test) |
| `internal/adapter/driving/cli/config_template.go` | update | `--print-config`: Zeile „Verfügbar“ und Block `linelength` |
| `spec/spezifikation.md` §4 | update | Zeile des Grund-Codes — kommt **mit** dem Code, der Lockstep-Test hält die Tabelle gegen `AllReasons()` |
| `spec/lastenheft.md` (Modul-Aufzählung in [`DC-FA-CLI-002`](../../../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl), Glossar) | update | nennt `linelength` erst, wenn das Modul registriert ist; Version und Historie nach [`MR-032`](../../../../harness/conventions.md#mr-032) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/user/operations.md` | update | Modul-Bullets und -Aufzählungen, die ein Spiegel-Test schon im Feature-Commit verlangt; der Handbuch-Kopf bleibt Release-Prep |
| `docs/user/abdeckung-tests.md` | update | von `make abdeckung` abgeleitet; `make test` hält sie gegen die Ableitung |
| Testdateien neben Regel und Konfigurations-Rand | neu | Happy/Boundary/Negative je Randbedingung der Anforderung; Spiegel-Deckung |

**Spiegel vor dem Editieren auflisten** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Tabelle oben ist der Stand der Planung, nicht der Suche. Vor dem ersten Edit läuft die
Aufzählungs-Suche nach einem bestehenden Modulnamen (etwa `pinsync`) über das ganze Repo;
jeder Treffer, der in der Tabelle fehlt, kommt als Plan-Nachtrag vor den Code.

**Testfälle aus der Anforderung** (Negativliste nach Schritt 19 des Workflow-Skeletts; die
Liste folgt den bestätigten Entscheiden 1–9 und dem
[Architect-Verdikt](../../../reviews/2026-11-22-welle-95-architect-verdikt.md); maßgeblich ist
die Anforderung, sobald sie steht):

- Zeile mit `max-chars` Zeichen ⇒ befundfrei; mit `max-chars + 1` ⇒ ein Befund — mit
  `exempt-unbreakable: false` oder ohne Wort über der Grenze; eine Zeile aus einem einzigen Wort
  mit `max-chars + 1` Zeichen ist unter Variante H im Default ausgenommen (Implementer-Hinweis).
- Zeile aus Mehrbyte-Zeichen (Umlaut, Emoji) an der Grenze ⇒ gezählt in Codepunkten,
  nicht in Bytes.
- `\r\n`-Zeilenende ⇒ das `\r` zählt nicht; Tab zählt als ein Zeichen; ungültiges UTF-8
  nach der Regel der Anforderung.
- Überlange Zeile in Fenced-Code (Backtick- und Tilde-Zaun, eingerückter Zaun) ⇒ mit
  Ausnahme an befundfrei, mit Ausnahme aus ein Befund; die Zaunzeile selbst ebenso.
- `~~~`-Zeile in einem Backtick-Block ⇒ der Block schließt früh, die Zeilen danach werden
  gemessen; eingerückter Code-Block (vier Leerzeichen, ohne Zaun) ⇒ gemessen (Verdikt F-3).
- Nicht geschlossener Fence bis Dateiende ⇒ der Rest der Datei ist ausgenommen, wie die
  Anforderung es als Grenze benennt; gemeldet wird der offene Zaun von `spans`
  (`fence-unclosed`), nicht von `linelength` (Verdikt F-3).
- Überlange `|`-Zeile ohne Trennzeile ⇒ ausgenommen (die Lexik verlangt keine Trennzeile);
  eingerückte `|`-Zeile, Zeile in einem Zitat (`> |`) und GFM-Zeile ohne führende Pipe ⇒
  gemessen; `exempt-code-blocks: false` mit `exempt-tables: true` ⇒ eine `|`-Zeile im Zaun wird
  gemessen (Verdikt F-4).
- Ausnahme (c) nach Variante H (Entscheid 5 (c), neu bestätigt 2026-11-23): Zeile mit einem Wort
  (längste Folge ohne Space und Tab), das allein länger als `max-chars` ist (lange URL) ⇒
  ausgenommen, auch mit umbrechbarem Text davor; Zeile über der Grenze, deren längstes Wort genau
  `max-chars` lang ist ⇒ Befund; Zeile nur
  aus kurzen Wörtern, deren letztes über die Grenze ragt ⇒ Befund; ein Tab trennt Wörter wie ein
  Space; ein geschütztes Leerzeichen (U+00A0) trennt nicht; letzte Zeile `…\r` ohne `\n` mit
  genau `max-chars` Zeichen ⇒ befundfrei (Architect-Befund A5).
- Zeile über der Grenze mit `<!-- d-check:ignore … -->` ⇒ Befund (Entscheid 9).
- Bezugsmenge: Datei mit und ohne abschließendes `\n` ⇒ dieselbe Zeilenzahl wie bei `file`;
  die über (c) ausgenommenen Zeilen stehen als eigene Zahl (Verdikt F-7).
- Zwei Dateien mit Befunden ⇒ Reihenfolge Pfad, dann Zeile; zweimal derselbe Lauf ⇒
  byte-gleiche Ausgabe.
- Keine Markdown-Datei in der Menge ⇒ kein Befund, Exit 0, Bezugsmenge `0`; explizit leere
  `linelength.scope.roots` ⇒ prüft nichts, kein Fehler (Entscheid 8,
  [`DC-FA-CONF-002`](../../../../spec/lastenheft.md#dc-fa-conf-002--modul-lokaler-scan-scope));
  Nicht-Markdown-Datei in `scan.roots` mit langer Zeile ⇒ kein Befund.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md)
liegt in `done/` (Anforderung, Spezifikation und ADR stehen, die Entscheide sind bestätigt),
und das WIP-Limit ist frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Spiegel-Suche ergibt so viele
  Stellen, dass Regel und Spiegel nicht in einer Review-Sitzung prüfbar sind — dann wird
  der Spiegel-Nachzug ein eigener Slice derselben Welle.
- `in-progress` → `open` (blockiert): die geteilte Lexik liefert für die Ausnahmen (a) oder
  (b) eine Antwort, die der Anforderung widerspricht, und eine Änderung der Lexik verschöbe
  andere Module — dann geht die Frage an die Anforderung zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Modul registriert, Tests je Randbedingung mit gezeigtem Rot, Spiegel grün,
`make gates` grün), Review- und Verifikations-Report liegen unter `docs/reviews/`,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein nicht geschlossener Fence nimmt über die geteilte Lexik den Rest der Datei aus der
  Prüfung — ein stilles Grün über einem Teil der Datei
  ([`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md),
  Fence-Hälfte, 2×) — **Ausgang:** *entfallen*, wenn ein Test das in der Anforderung
  festgelegte Verhalten beim offenen Fence belegt und die Anforderung für die Meldung auf
  `spans` (`fence-unclosed`) verweist — ein zweiter Wächter in `linelength` wäre eine zweite
  Antwort (Architect-Verdikt F-3); *weiter offen* — ein Beleg, der den Zähler auf 3× hebt —,
  wenn der Ausfall in keinem Modul sichtbar wird.
- `proseLineSet`/`tableRowLine` sind paket-interne Helfer anderer Module; sie für ein neues
  Modul zu nutzen bindet dessen Verhalten an deren Änderungen
  ([`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md),
  verkörpert über Kopplungs-Tests; Pflicht nach [ADR-0054](../../adr/0054-geteilte-lexik-bindet-ihre-konsumenten.md) Entscheidung 4, DoD-Punkt 2) —
  **Ausgang:** *entfallen*, wenn der Kopplungs-Test
  dieselbe Eingabe durch `linelength` und die bisherigen Konsumenten fährt; *eingetreten*
  (Rückführung nach `next`, §4), wenn die Helfer für das Modul geändert werden müssten.
- Eine Randbedingung der Anforderung bleibt ohne Test
  ([`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md),
  1×) — **Ausgang:** *entfallen*, wenn jede Randbedingung der Anforderung einen Test mit
  gezeigtem Rot hat (Zuordnung im Plan-Nachtrag); *weiter offen* — ein Beleg in diesem
  Eintrag —, wenn eine fehlt.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:373-374 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) —
Regelmodul, Konfigurations-Rand, CLI-Vorlage, Spiegel-Tests; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-22), Sub-Area `*`. Treffer:
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(verkörpert in den Spiegel-Deckungstests — die DoD verlangt, dass sie das neue Modul
fangen),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(verkörpert in Schritt 19 des Workflow-Skeletts — die Testfall-Liste in §3 ist seine
Anwendung),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — ein eigenes Fence-Muster wäre der dritte Beleg),
[`BEO-ALL/prosa-verschwindet-in-inline-code-spanne`](../observations/BEO-ALL/prosa-verschwindet-in-inline-code-spanne/observation.md)
(offen, 2× — Risiko 1 in §6),
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert über Kopplungs-Tests — Risiko 2 in §6),
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen, 1× — Risiko 3 in §6) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — das Modul liest nur seine
Scan-Menge). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-22: **beide rot** — `upstream-drift` mit fünf veralteten
Pins (Go, Semgrep, a-check, `docker/login-action`, Digest `golang:1.27.3`), `image-scan` mit
behebbaren CRITICAL/HIGH im publizierten Image; Einzelheiten im Plan von
[`slice-289`](../in-progress/slice-289-linelength-anforderung-und-adr.md) §8. Berührt das Modul nicht;
eigene Pin-Hebung vor dem nächsten Release. Beim Start dieses Slice neu lesen.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
