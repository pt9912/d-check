# Slice slice-287: Modulkern `pinsync` — Konfiguration lesen, Quellen fangen, Werte vergleichen, Grund-Codes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-94-pinsync-modul`](../welle-94-pinsync-modul.md).

**Bezug:** die Anforderung aus `slice-286`
(Kennung steht nach deren Closure fest),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`DC-FA-CONF-001`](../../../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei),
[`DC-FA-TGT-001`](../../../../spec/lastenheft.md#dc-fa-tgt-001--deklarations-konsistenz-zwischen-doku-und-build-targets-modul-targets-opt-in)
(Vorbild: hermetisches Lesen benannter Dateien).

**Berührte Spec-Stellen:** die Anforderung und die Spezifikations-Abschnitte aus
`slice-286`; `architecture.md` bleibt
unberührt. Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das opt-in-Modul `pinsync` läuft end-to-end: es liest den
Konfigurationsblock, fängt je Quelle den Wert und meldet ungleiche Werte, einen
Ausdruck ohne Treffer und einen mit mehreren, alles hermetisch und in fester
Reihenfolge; jeder Eingabefehler ist Exit 2.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Anforderung, die Grund-Codes als Spezifikation und die ADR** — stehen nach
  `slice-286` fest; dieser Slice löst sie
  ein und ändert sie nicht. Zeigt der Code eine Lücke der Anforderung, ist das eine
  Plan-Änderung vor dem Code, kein Nachtrag im Bericht.
- **Das Scharfschalten im eigenen Profil** — gehört
  `slice-288`: ein Modul wird erst gebaut und
  geprüft, dann im eigenen Gate verdrahtet.
- **CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). Den Handbuch-**Abschnitt** des
  Moduls liefert dieser Slice, weil er den öffentlichen Vertrag beschreibt.
- **Ein Normalisierungs-Schalter und jede Aktualitätsprüfung** — vom Auftraggeber
  ausgeschlossen; ein neuer Schlüssel im Block wäre eine Plan-Änderung.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Konfiguration des Moduls wird gelesen und am Rand geprüft: jeder
      Eingabefehler der Anforderung (fehlende oder unlesbare Datei, ungültiger
      Ausdruck, Anzahl der Fanggruppen, Pin mit weniger als zwei Quellen, doppelter
      Pin-Name, Pfad mit `..`/absolut/Glob, unbekannter Schlüssel) endet mit Exit 2 und
      hat einen Test; `--print-config` und die Vorlage tragen den Block.
- [ ] Die Prüfung fängt je Quelle den Wert zeilenweise und meldet die drei Befund-Klassen
      mit den Grund-Codes der Spezifikation in der festgelegten Reihenfolge; jede
      Randbedingung der Anforderung hat einen Test, **mit** den Negativfällen aus
      `slice-286` §3.
- [ ] Das Modul ist registriert (Modulliste, `--enable`, Spiegel-Deckungstests,
      `gate-consistency`-Test) und im Benutzerhandbuch beschrieben, mit dem
      `v?`-Muster als Beispiel.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die Modul-Aufzählungen, soweit ein öffentlicher Vertrag berührt ist.
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
| `internal/hexagon/core/model/config.go` | update | Konfigurations-Typ des Moduls; Modulliste |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Block lesen und am Rand prüfen (Exit 2) |
| Regel-Datei `pinsync` im Verzeichnis `internal/hexagon/core/rules/` | neu | Quellen lesen, Wert fangen, vergleichen, Grund-Codes |
| `internal/hexagon/core/rules/run.go` | update | Modul im Lauf verdrahten |
| `internal/adapter/driving/cli/` (Vorlage, `--enable`) | update | `--print-config`, Modulliste |
| `docs/user/benutzerhandbuch.md`, `docs/user/operations.md` | update | Abschnitt des Moduls, Modul-Aufzählungen |
| Testdateien neben den genannten Dateien | neu / update | Happy/Boundary/Negative je Randbedingung der Anforderung; Spiegel-Deckung; `gate-consistency` |

Der Implementer erweitert die Liste im ersten Lauf, einschließlich einer Zeile je Testdatei
mit der Akzeptanzkriterien-Kennung (Modul 9). **Vor dem Code-Commit** (Schritt 17 des
Workflow-Skeletts): die Modul-Aufzählungen (Modulliste, README-Listen, Handbuch,
`operations.md`, `print-config`-Vorlage, Spiegel-Tests) werden mit `grep` über den Namen
eines Bestandsmoduls gezählt (etwa `mentions`), jeder Treffer steht hier im Nachtrag.

**Negativfälle der Fang-Logik (Schritt 19):** Kommentarzeile mit demselben Text; Wert in einer
`FROM`-Zeile; zwei Zeilen mit Treffer; Treffer mit `v`-Präfix außerhalb und innerhalb der
Gruppe; CRLF am Zeilenende; Ausdruck ohne Fanggruppe; leere Fanggruppe; leere Datei;
Datei ohne abschließenden Zeilenumbruch.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-286`
liegt in `done/` (die Anforderung und die ADR stehen) und das WIP-Limit ist frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Diff passt nicht in eine
  Review-Sitzung, oder der Slice braucht zusätzlich einen Port (Datei-Lesen
  hinter dem Port statt direkt), den die Anforderung nicht vorsieht.
- `in-progress` → `open` (blockiert): `slice-286`
  wird zurückgeführt oder die Anforderung ändert sich im Lauf.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review- und Verifikations-Report liegen unter
`docs/reviews/`, und die Verifikation hat die Testbehauptungen der Anforderung gebrochen
und rot gesehen (Modul 11 §Bewusstes Brechen); Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Das Lesen der Quelldateien läuft im Kern am Dateisystem vorbei am Port (§3.8, Hexagon) —
  **Ausgang:** *entfallen*, wenn der Kern die Datei über den vorhandenen Lese-Port bekommt
  (so liest `targets` seine Makefiles); *eingetreten* (neuer Port ⇒ Rückführung nach
  `next`), wenn dieser Port die Quelle nicht liefert.
- Ein Ausdruck mit Rückwärts-Verweis oder Lookaround ist in RE2 ungültig — **Ausgang:**
  *entfallen*, die Anforderung macht ihn zum Eingabefehler (Exit 2), der Test belegt die
  Meldung.
- Die Zusage „kein Schreiben“ unter `:ro` — **Ausgang:** *entfallen* mit Beleg aus dem
  `image-test`-Lauf unter `:ro` in `slice-288`;
  der Kern selbst öffnet nur lesend.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln
(das Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde).

Entsteht bei der Closure.

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
Regelmodul, Konfigurations-Rand, Gerüst, Spiegel-Tests; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/modulliste-spiegel-ungegated`](../observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md)
(3×, verkörpert in den Spiegel-Deckungstests — die DoD verlangt, dass sie das neue
Modul fangen),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(verkörpert in Schritt 19 des Workflow-Skeletts — die Negativliste in §3 ist seine
Anwendung),
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — die Quelldateien sind Eingaben
außerhalb der Scan-Menge),
[`BEO-ALL/gemeinsame-leser-sind-zwei-familien`](../observations/BEO-ALL/gemeinsame-leser-sind-zwei-familien/observation.md)
(offen, 2× — der Zeilen-Leser darf kein eigenes Muster neben dem geteilten sein; ein
dritter Beleg in diesem Slice wäre die Schwelle) und
[`BEO-ALL/spec-randbedingung-ohne-test`](../observations/BEO-ALL/spec-randbedingung-ohne-test/observation.md)
(offen — jede Randbedingung der Anforderung braucht einen Test). Der Zähler folgt aus den
Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-31: `upstream-drift` grün (jüngster Lauf 2026-10-31T06:31:58Z),
`image-scan` grün (2026-10-31T09:41:43Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
