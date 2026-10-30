# Slice slice-284: Messung aus git — VCS-Port, Adapter und Prüfung gegen die Obergrenze

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-93

**Bezug:** die Anforderung des Moduls `metrics` aus [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) (die Hälfte „Messung gegen die Obergrenze“),
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) und
[`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in) (Port-Vorbild, Fail-closed-Haltung),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus) (gleiche Historie ⇒ gleiche Summe),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) (lesend, netzlos, `:ro`).

**Berührte Spec-Stellen:** die in [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) angelegten Stellen (Anforderung,
Messsemantik der Summe, Grund-Codes `metrics-over-limit` und Fail-closed-Fall) — gelesen;
die Zusage der `numstat`-Gleichheit wird hier belegt oder korrigiert.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-30.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Modul `metrics` misst pro Plan mit Block die Summe aus Einfügungen
und Löschungen zwischen der Basis des Blocks und `HEAD` über die Pfade des Blocks —
über den VCS-Port, ohne git-Binary und ohne Netz — und meldet eine Messung über der
Obergrenze, wenn der Plan keinen Ausnahme-Eintrag mit Grund trägt; ohne `.git` endet
der Lauf mit Exit 2.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Block-Lexik, Formel, Config-Leser, leerer Grund** — Bestand dieses Moduls aus
  [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle), davor `done/`; dieser Slice nimmt seine Ausgabe (Obergrenze, Pfade, Basis,
  Ausnahme) als gegeben und ändert sie nicht.
- **Eine Rename-Erkennung nach git-Art, Binärdatei-Gewichtung, Kopien-Erkennung** —
  der Entscheid steht in der Anforderung ([`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle), Entscheid 4); eine andere Semantik wäre eine
  Änderung der Anforderung, nicht dieses Slice.
- **Weitere VCS-Port-Methoden und eine Änderung der bestehenden** — der Port wächst um
  genau eine Methode; die Verträge von `vcs`/`commits` bleiben byte-gleich (Schicht-
  Abgrenzung: kein Umbau fremder Module im selben Slice). Auch `resolveVCS` wächst nur
  um den Fall für `metrics`; das Verhalten für `vcs`, `commits` und `tracked` bleibt.
- **Release-Prep (Handbuch, `operations.md`, README-Modulliste)** — ein anderer Vorgang
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Port und Adapter** (Entscheid 8 der ADR; bei Entscheid für das Zählen im Kern
      entfällt die Port-Methode, und die Zeilen-Zählung steht im Regel-Paket): der VCS-Port trägt eine Methode für die Zeilen-Summe (literale Pfade, ausgenommener Teilbaum) über
      Pfade zwischen zwei Ständen; der reine-Go-Adapter implementiert sie ohne
      git-Binary. **Gleichheit mit `git diff --numstat` (bzw. `--no-renames`, je
      Entscheid 4) ist an einem Fixture mit allen Randformen belegt** — Hinzufügen,
      Löschen, Ändern, leere Datei, Binärdatei, Rename, Pfad außerhalb der Pathspec,
      Zeilenenden-Wechsel — nicht behauptet.
- [ ] **Prüfung „Messung ≤ Obergrenze“:** Messung über der Obergrenze ohne Ausnahme
      ist `metrics-over-limit`; Messung auf der Grenze und darunter ist grün;
      Ausnahme **mit** Grund unterdrückt den Befund, Ausnahme **ohne** Grund bleibt der
      Befund aus [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle); Test referenziert, Negativliste (§3) getestet.
- [ ] **Config-Schlüssel `metrics.excluded`:** Pfade (literal, ohne Glob, `..`, führenden `/`, Steuerzeichen) kommen in den Config-Rand; Verstöße ⇒ Exit 2.
- [ ] **Fail-closed und Betrieb:** ohne `.git` Exit 2 (nicht „0 Befunde“); eine Basis,
      die nicht auflöst, ist Exit 2 mit Klartext; der Lauf bleibt unter `:ro` und gegen
      ein Repo mit anderem Eigentümer lauffähig — am gebauten Image belegt
      (`make image-test`-Fall in `tools/image-test.sh`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Benutzerhandbuch und `operations.md` **nicht** (Release-Prep, §1).
- [ ] Grund-Code `metrics-over-limit`: Zeile in der Tabelle von `spec/spezifikation.md` §4, Eintrag in `AllReasons()` und `--doctor`-Klartext (`TestAllReasonsDeckungGegenSpezifikationGrundCodes`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/vcs.go` | update | eine Methode für die Zeilen-Summe über literale Pfade zwischen zwei Ständen, mit ausgenommenem Teilbaum; Basis auch als Kurz-Hash auflösbar |
| `internal/adapter/driven/git/git.go` | update | Implementierung über den reinen-Go-Diff, ohne git-Binary |
| `metrics.go` im Regel-Paket | update | Messung abfragen, Vergleich mit der Obergrenze, Ausnahme-Logik |
| `internal/hexagon/core/rules/run.go` (`runPostPasses`) | update | Modul in die nicht-datei-scannenden Läufe aufnehmen; Fail-closed ohne `.git` wie `vcs`/`commits` |
| `internal/adapter/driving/cli/cli.go` (`resolveVCS`, Aufruf von `RunWithVCS`) | update | dritter Fall neben „Adapter mit Range“ (`vcs`/`commits`) und „Adapter ohne Range“ (`tracked`): Adapter, **Basis aus dem Block**, `HEAD` fest |
| Test-Doubles des `driven.VCS`-Ports (u. a. in `vcs_test.go`, `tracked_test.go` im Regel-Paket) | update | wachsen mit der neuen Port-Methode mit |
| `internal/adapter/driven/git/git_test.go` | update | Fixture-Gleichheit mit `git diff --numstat`, Randformen aus der DoD |
| `metrics_test.go` im Regel-Paket | update | Happy/Boundary/Negative nach der Anforderung; Negativliste unten |
| `tools/image-test.sh` | update | `:ro` und fremder Eigentümer am gebauten Image (`make image-test`) |
| `internal/hexagon/core/model/config.go` · `internal/adapter/driven/configyaml/configyaml.go` · `internal/adapter/driving/cli/config_template.go` | update | `metrics.excluded` am Config-Rand (`model.ValidLiteralPath`), Gerüst-Zeile |
| `internal/hexagon/core/app/diagnose.go` | update | Grund-Code `metrics-over-limit` in `AllReasons()` samt Klartext |
| `internal/adapter/driven/git/git_lines_test.go` · `internal/adapter/driven/configyaml/metrics_test.go` | neu · update | Fixture-Gleichheit mit `git diff --numstat` in eigener Datei neben `git_test.go`; Config-Test für `excluded` |
| `spec/spezifikation.md` · `docs/user/abdeckung-tests.md` | update | Zeile `metrics-over-limit` in der Tabelle von §4, Abdeckung per `make abdeckung` |

**Negativliste der Messung (Schritt 19, vor dem Code-Commit zu testen):** Pfad im Block
existiert weder an der Basis noch an `HEAD` (Messung 0, kein Fehler — Entscheid in der
Anforderung) · Basis ist `HEAD` (Summe 0) · Basis ist kein Vorfahre von `HEAD` (Exit 2) · Basis als Kurz-Hash (7 Zeichen) und als voller Hash · konfigurierter Teilbaum, den ein Pfad des Blocks enthält (ausgenommen) · nicht eingecheckte Änderung (zählt nicht) · dieselbe Datei über zwei Pfade (zählt einmal) ·
Basis auf einem anderen Zweig · shallow-Klon ohne Basis · Binärdatei (`-` in numstat, zählt 0 mit Hinweis) ·
Datei, die nur Rechte ändert · Rename mit und ohne Inhaltsänderung · Pathspec mit
Verzeichnis-Präfix und mit `*`/`[` im Namen (literal, kein Glob) · Pfad mit Leerzeichen und Nicht-ASCII · lokaler
Branchname, der wie eine Revision aussieht ([`BEO-ALL/range-leerfall-mit-lokalem-branchnamen`](../observations/BEO-ALL/range-leerfall-mit-lokalem-branchnamen/observation.md)) ·
frisch erzeugtes Fixture-Repo, dessen Index „racily clean“ ist ([`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-283`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) liegt in `done/` (Obergrenze, Pfade,
Basis und Ausnahme sind als Ausgabe des Kerns vorhanden), und welle-93 ist offen. Der
Übergang landet auf dem Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Fixture-Gleichheit zeigt,
  dass der reine-Go-Diff die numstat-Semantik nicht trägt und eine eigene Zeilen-Zählung
  nötig wird — dann gehört die Zählung in einen eigenen Slice vor diesem.
- `in-progress` → `open` (blockiert — Carveout?): die Gleichheit lässt sich ohne
  git-Binary nicht herstellen; dann ist die Zusage in der Anforderung zu ändern oder ein
  Carveout zu eröffnen (die Entscheidung trifft der Auftraggeber, nicht der Lauf).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, und am gebauten Image
liefert das Fixture-Repo (Plan mit Block, Basis-Commit, Änderungen über und unter der
Grenze) die erwarteten Befunde — unter `:ro` und mit fremdem Eigentümer des Repos —
während der Lauf ohne `metrics`-Konfiguration gegen das Vorher-Image byte-gleich ist.
Die ADR aus [`slice-282`](../welle-93-metrics-modul.md#4-slices-in-dieser-welle) geht von `Proposed` auf `Accepted`; ist die `numstat`-Zusage
zu korrigieren, geschieht das vorher, solange die ADR noch nicht `Accepted` ist.
Dazu die Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der reine-Go-Diff zählt anders als `git diff --numstat` (Zeilenenden, Binärerkennung,
  Diff-Algorithmus bei verschobenen Blöcken); die Summe wäre dann stabil, aber nicht die,
  die das Anlass-Repo bisher gemessen hat — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Die Messung liest `.git` und Pfade außerhalb der Scan-Wurzeln; die Zusage der
  Scan-Menge ([`AGENTS.md`](../../../../AGENTS.md) §3.8) gilt dort nicht, und ein nicht
  auflösbarer Pfad kann still 0 ergeben — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)
- Große Diffs über viele Pfade kosten Laufzeit; [`DC-QA-01`](../../../../spec/lastenheft.md#dc-qa-01--performance) gilt für das Modul nur, soweit
  es aktiviert ist — **Ausgang:** wird bei der Closure gesetzt (eingetreten / entfallen / weiter offen)

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

Wird bei der Closure gefüllt.

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
VCS-Port, git-Adapter, Regelmodul; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/range-leerfall-mit-lokalem-branchnamen`](../observations/BEO-ALL/range-leerfall-mit-lokalem-branchnamen/observation.md)
(1×, ungeprüft — die Basis des Blocks ist ein Name, der wie ein Zweig aussehen kann;
Negativliste §3),
[`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)
(1×, offen — das Fixture-Repo ist frisch erzeugt; Negativliste §3) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(3×, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — `.git` und die gemessenen
Pfade liegen außerhalb der Scan-Menge, §6). Erreicht eine der beiden ersten mit diesem
Slice ihr zweites Auftreten, wird das in §7 benannt; die dritte Evidenz wäre eine Lücke
und brauchte einen eigenen Slice. Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-30: `upstream-drift` grün (jüngster Lauf 2026-10-30T06:41:15Z), `image-scan`
grün (2026-10-30T10:04:41Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
