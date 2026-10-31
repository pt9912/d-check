# Slice slice-288: Eigenes Profil — die zwei Doppel-Pins von `Dockerfile` und `Makefile` im Gate, beidseitig belegt

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
[`AGENTS.md`](../../../../AGENTS.md) §3.6 (die Aufnahme eines unabhängig lauffähigen
Wächters in `make gates` ist kein ADR-Anlass; ein Verweis auf die tragende ADR genügt).

**Berührte Spec-Stellen:** — (die Spec-Stellen gehören
`slice-286`; dieser Slice konfiguriert nur).

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

**Ziel:** d-check prüft mit `pinsync` im eigenen Gate (`make doc-check`), dass
`GO_VERSION` und `GOLANGCI_LINT_VERSION` im `Dockerfile` und im `Makefile` gleich
sind — und der Lauf ist **beidseitig** belegt: unverändert grün, bei einer Abweichung
in genau einer Datei rot mit dem Grund-Code des Moduls.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Weitere Pins** (`@sha256:`-Digests der Basis-Images, `a-check`-Version, Baseline-Tag,
  Pins in `a-check.mk`) — gemessen ist nur das Paar `GO_VERSION`/`GOLANGCI_LINT_VERSION`
  (je zwei Stellen, je eine Zeile); die Digests gehören zur Digest-Achse der Nachtläufe
  ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)) und ein Pin pro Zusage
  braucht einen eigenen Bestandsbeleg. Der Eintrag
  [`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
  bleibt davon unberührt.
- **Änderungen am Modul selbst** — gehören `slice-287`;
  zeigt der Selbstlauf einen Fehler des Moduls, geht der Slice zurück oder ein
  Folge-Slice entsteht, statt das Modul hier zu flicken.
- **Eine Senkung oder Lockerung eines bestehenden Gates** — der Block kommt zu den
  bestehenden hinzu; kein Schwellenwert, kein Ausnahme-Eintrag wird berührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die `.d-check.yml` führt das Modul mit zwei Pins zu je zwei Quellen (vier Ausdrücke);
      jeder Ausdruck hat am Bestand **genau einen** Treffer (Zählbefehl im Plan-Nachtrag
      mit Ergebnis), das `v` von `GOLANGCI_LINT_VERSION` steht außerhalb der Fanggruppe.
- [ ] Die Brechprobe liegt als Beleg vor: mit einem in nur einer Datei gehobenen Wert endet
      `make doc-check` rot mit dem Grund-Code des Moduls und der Fundstelle; die Rücknahme
      ist wieder grün (beide Richtungen, echte Ausgabe im Report).
- [ ] Das Gate-Verzeichnis (`harness/README.md` §Sensors, `harness/sensors/doc-check.md`)
      nennt die neue Zusage von `doc-check` mit ihrer Bindung (die Anforderung aus
      `slice-286`, die ADR).
- [ ] `make gates` grün.
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
| `.d-check.yml` | update | Modul `pinsync` mit zwei Pins; Eintrag in der Modulliste von `doc-check` |
| `harness/README.md` | update | Bindung der Zeile `make doc-check` |
| `harness/sensors/doc-check.md` | update | Zusage und Grenze des neuen Teils |

**Bestandsmessung (vor dem Eintrag, Ergebnis im Plan-Nachtrag):** am 2026-10-31 gezählt mit
`grep -cE '^ARG GO_VERSION=[^ ]+$' Dockerfile` (1),
`grep -cE '^ARG GOLANGCI_LINT_VERSION=v?[^ ]+$' Dockerfile` (1),
`grep -cE '^GO_VERSION[[:space:]]*\?=[[:space:]]*[^ ]+$' Makefile` (1) und
`grep -cE '^GOLANGCI_LINT_VERSION[[:space:]]*\?=[[:space:]]*v?[^ ]+$' Makefile` (1);
beide Dateien tragen kein CRLF (`grep -c $'\r'` ⇒ 0). Der Implementer zählt vor dem Eintrag
erneut — der Stand kann sich bis dahin bewegt haben.

**Plan-Nachtrag vor dem Code (2026-10-31).** Bestandsmessung wiederholt: 1 / 1 / 1 / 1 Treffer, kein
CR (`grep -c $'\r'` ⇒ 0 in beiden Dateien); die Zeilen sind `Dockerfile:31`, `Dockerfile:32`,
`Makefile:13`, `Makefile:14`. Das Makefile schreibt `?=` mit ausgerichteten Leerzeichen, die
Ausdrücke dafür tragen `\s*`. Zwei Dateien stehen zusätzlich in der Änderungsliste, weil die
Modulliste von `.d-check.yml` Spiegel hat (Hinweis in der Datei selbst, gezählt mit
`grep -rln 'disable citations\|"citations"' .` außerhalb von Plan, Reviews und Baseline): die
Variable `FOCUS_DISABLE` im `Makefile` (sonst liefe `pinsync` in jedem fokussierten Gate mit) und
`netlessDocModules()` im Test von `internal/adapter/driven/configyaml/gate_consistency_test.go`
(die Messmethode für [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)). Die übrigen Treffer (`configyaml.go`, `config.go`, `run.go`,
`citations.go`) sind Modul-Code von `citations` und tragen keine Liste der Default-Module.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-287` liegt in
`done/` und das gebaute Image trägt das Modul (der Selbstlauf braucht es im
`doc-check`-Image).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Selbstlauf zeigt, dass mehr
  als das Paar der Pins nötig ist, um das Gate sinnvoll zu machen.
- `in-progress` → `open` (blockiert): das Modul fängt einen der vier Ausdrücke am Bestand
  nicht genau einmal, und die Ursache liegt im Modul, nicht im Ausdruck.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, beide Richtungen der Brechprobe mit echter Ausgabe im
Report, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein Pin-Bump in nur einer Datei bricht künftig `make gates` und damit jede
  Dependabot-Hebung der beiden Pins — **Ausgang:** *entfallen*: das ist die Zusage;
  `.github/dependabot.yml` führt `docker` bewusst nicht
  ([ADR-0067](../../adr/0067-dependabot-als-hebender-kanal.md)), die Hebung der beiden
  Pins ist ein manueller Akt (Go-Bump zieht das Pendant nach), und die Meldung nennt beide
  Fundstellen.
- Das Modul liest `Dockerfile` und `Makefile` außerhalb der Scan-Menge (kein
  Markdown, §3.8) — **Ausgang:** *entfallen*, die Anforderung nennt die Grenze und die
  Brechprobe belegt, dass beide Dateien tatsächlich gelesen werden (ein Wert nur in der
  `Makefile` gehoben ⇒ rot).

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
Konfiguration des eigenen Profils und Gate-Verzeichnis; `tools/harness/` bleibt unberührt.

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen, Sub-Area
`*`. Treffer:
[`BEO-ALL/pin-bump-mirrors-ungated`](../observations/BEO-ALL/pin-bump-mirrors-ungated/observation.md)
(8×, ohne formgültigen Ausgang — dieser Slice deckt davon nur die Wertgleichheit zweier
Pins und löst den Eintrag **nicht** auf; die Closure-Notiz sagt das ausdrücklich, statt
einen Ausgang zu behaupten) und
[`BEO-ALL/registry-vs-authority-table-drift`](../observations/BEO-ALL/registry-vs-authority-table-drift/observation.md)
(die Bindung im Gate-Verzeichnis muss dem Makefile und dem Index gleichen). Der Zähler folgt
aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-10-31: `upstream-drift` grün (jüngster Lauf 2026-10-31T06:31:58Z),
`image-scan` grün (2026-10-31T09:41:43Z). Nichts zu lesen, nichts ausstehend.

### Sub-Area: `*` (Default)

Alle berührten Sub-Areas GF.
