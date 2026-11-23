# Slice slice-292: Kern und git-Adapter — `vcs` liest die Umbenennungen der Range und setzt Link-Ziele über sie gleich

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-96-vcs-folgt-umbenennung`](../welle-96-vcs-folgt-umbenennung.md).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
in der Fassung aus [`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md)
und [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[`ADR-0024`](../../adr/0024-vcs-immutable-gate.md),
[`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).

**Berührte Spec-Stellen:** — (Verfahren und Schema schreibt
[`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md); kein neuer
Grund-Code, §4 der Spezifikation bleibt). Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-11-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit `vcs.follow-renames` liest das Modul die Umbenennungen der Range über den
git-Port (fail-closed), setzt ein Link-Ziel über sie gleich, wie
[`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) es festschreibt, und
ein reiner Nachzug auf die umbenannte Datei ist keine Drift mehr — ohne den Schlüssel
byte-identisch; jede Randbedingung der Anforderung hat einen Test, der ohne die Regel aus
dem richtigen Grund rot läuft, und die ADR ist `Accepted`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anforderung, Verfahren, ADR-Wortlaut ändern** — gehören
  [`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md). Zeigt der Code,
  dass ein Entscheid nicht trägt, geht das als Plan-Änderung an den Architect, nicht als
  stille Abweichung im Code; die ADR ändert nur der Architect.
- **`metrics`, `tracked`, `commits`** — lesen git, folgen aber keiner Umbenennung;
  `LineChanges` bleibt ohne Rename-Erkennung (sein Vertrag sagt es).
- **Die Grenze „gleichnamige andere Datei“** von `ignore-link-targets` — Welle-Out-of-Scope.
- **CHANGELOG, Handbuch-Kopf, Release** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17). Was ein Spiegel-Test schon im
  Feature-Commit verlangt, kommt mit und steht in §3.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] git-Port und -Adapter liefern die Umbenennungen der Range nach der Anforderung
      (Commit-Menge, Inhalts-Gleichheit, Kette, Mehrdeutigkeit, `--staged`), ein
      unlesbarer Commit oder Baum ist Exit 2; Tests gegen echte git-Fixtures.
- [ ] Regel und Konfigurations-Rand: `vcs.follow-renames` wird gelesen (falscher Typ ⇒
      Exit 2), das Gleichsetzen trägt jedes Akzeptanzkriterium der Anforderung, ohne
      Schlüssel ist die Ausgabe byte-identisch; die Negativfälle aus dem Plan von
      [`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) stehen als
      Testfälle, und für jeden Test ist ein Rot ohne die Regel gezeigt
      (Baseline-Regelwerk `modul-11-verification.md` §Bewusstes Brechen).
- [ ] Spiegel und Profil: `--print-config`, Handbuch-Abschnitt des Moduls, READMEs und
      [`harness/sensors/adr-check.md`](../../../../harness/sensors/adr-check.md) nennen den
      Schlüssel (vor dem Editieren aufgelistet,
      [`MR-025`](../../../../harness/conventions.md#mr-025)); `.d-check.yml` schaltet ihn
      ein **im selben Commit**, der
      [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) auf `Accepted` hebt und den
      Index nachzieht (Entscheid 14); [`AGENTS.md`](../../../../AGENTS.md) §3.5 und
      [`harness/README.md`](../../../../harness/README.md) §Traceability rules nennen
      [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) neben
      [`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/vcs.go` | update | neue Methode am Port, gerufen nur bei aktivem Schalter: die Abbildung Quelle → eindeutiges Ende (Vertrag nach [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13: Commit-Menge, Inhalts-Gleichheit nach Blob-Kennung und Art ohne Ausführbar-Bit, Kette entlang der Abstammung, Mehrdeutigkeit, Basis kein Vorfahr ⇒ leer, fail-closed auf jedes nötige Objekt; Kette und Mehrdeutigkeit entscheidet der Adapter, das Gleichsetzen der Kern) |
| `internal/adapter/driven/git/git.go` | update | eigener Blob-Vergleich über die Bäume je Commit und Eltern-Commit bzw. HEAD gegen Index; Teilbäume gleicher Kennung bleiben ungelesen; die Rename-Erkennung von go-git wird **nicht** benutzt; kein Diff-Walker, der einen unlesbaren Unterbaum still überspringt. **Dazu (Plan-Änderung nach dem Architect-Befund S-1, 2026-11-23):** `AllPaths` nimmt im gestagten Modus heute jeden Index-Eintrag, auch einen Gitlink (`git.go` Z. 95–98), die Spezifikation (§[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2) sagt „jeder Datei-Pfad (regulär/ausführbar/symlink) … bzw. am staged Index“ — der Code folgt der Spezifikation (Rang 2): dieselben drei Arten wie im Baum, ein Filter für Enumeration und Umbenennungs-Abbildung |
| `internal/hexagon/core/model/config.go` | update | `VCSConfig` um den Schalter |
| `internal/adapter/driven/configyaml/configyaml.go` | update | Schlüssel lesen und am Rand prüfen (Exit 2) |
| `internal/hexagon/core/rules/vcs.go`, `vcs_link_targets.go` | update | Umbenennungen nur bei aktivem Schalter lesen; das Gleichsetzen im Auflöser der Link-Ziele |
| `internal/adapter/driving/cli/config_template.go` | update | `--print-config`: Zeile des Schlüssels im `vcs`-Block |
| `.d-check.yml` | update | `follow-renames: true` im `vcs`-Block (Entscheid 7), im selben Commit wie der Übergang der ADR nach `Accepted` (Entscheid 14); Kommentar nennt die ADR |
| `docs/user/benutzerhandbuch.md`, `harness/sensors/adr-check.md` | update | nennen `ignore-link-targets` und die Grenze „anderer Name bleibt Drift“, die der Schalter ändert (Workflow-Skelett `implement-slice` Schritt 17); der Handbuch-Kopf bleibt Release-Prep |
| `README.md`, `README.de.md` | — | nennen die Grenze „anderer Name“ nicht; ein Satz zum neuen Schalter ist README-Prosa und gehört in die Release-Prep ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), außer ein Spiegel-Test verlangt ihn schon hier |
| [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md), `docs/plan/adr/README.md` | update | Status `Accepted` und Index — durch den Architect |
| `AGENTS.md` §3.5, `harness/README.md` §Traceability rules | update | nennen [`ADR-0115`](../../adr/0115-vcs-folgt-umbenennung.md) neben [`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) (Entscheid 14) |
| `docs/user/abdeckung-tests.md` | update | von `make abdeckung` abgeleitet; `make test` hält sie gegen die Ableitung |
| Testdateien neben Adapter, Regel und Konfigurations-Rand; der Test-Doppelgänger des Ports in `internal/hexagon/core/rules/vcs_test.go` | neu / update | Happy/Boundary/Negative je Kriterium von [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in); jeder Implementierer des Ports folgt der neuen Methode |

**Spiegel vor dem Editieren auflisten** ([`MR-025`](../../../../harness/conventions.md#mr-025)):
die Tabelle oben ist der Stand der Planung. Vor dem ersten Edit läuft die Suche nach
`ignore-link-targets` und `IgnoreLinkTargets` über das ganze Repo (ohne `done/`,
`docs/reviews/` und `CHANGELOG.md`, die Historie tragen); jeder Treffer, der in der
Tabelle fehlt, kommt als Plan-Nachtrag vor den Code. Stand der Planung (2026-11-23):

```bash
grep -rln "ignore-link-targets\|IgnoreLinkTargets" --exclude-dir=.git --exclude-dir=done . \
  | grep -vE "^./(docs/plan/planning|docs/reviews)/|CHANGELOG"
```

**Testfälle aus der Anforderung** (Negativliste nach Schritt 19 des Workflow-Skeletts; die
Liste folgt den bestätigten Entscheiden 1–14 und dem Architect-Verdikt):

- `git mv A B` und Nachzug in derselben `Accepted`-ADR, ein Commit ⇒ mit Schalter kein
  Befund, ohne Schalter `core-drift-vcs` (heutiges Verhalten, byte-identisch).
- Derselbe Nachzug, A gelöscht und B mit anderem Inhalt neu angelegt ⇒ Befund.
- Umbiegen von A auf eine bestehende, nicht umbenannte Datei C ⇒ Befund.
- Nachzug mit geändertem Anker oder Linktext ⇒ Befund.
- Umbenennen und im selben Commit ändern ⇒ Befund (keine Inhalts-Gleichheit).
- Kette über drei Commits (A → B, Inhaltsänderung an B, B → C) ⇒ folgt (Entscheid 2).
- Zwei verschwundene oder zwei neue Pfade mit demselben Blob; Kette mit zwei Enden ⇒
  Befund (Mehrdeutigkeit).
- Merge als HEAD mit einer umbenennenden Seite ⇒ die Umbenennung trägt (Entscheid 11);
  parallele Linien, die denselben Pfad verschieden umbenennen ⇒ Befund.
- Basis kein Vorfahr von HEAD ⇒ keine Umbenennung, Urteil wie ohne Schalter
  (Entscheid 12), kein Exit 2.
- Unveränderter Link auf eine umbenannte Datei in einer ADR, die nur einen erlaubten
  Anhang bekommt ⇒ kein Befund (seiten-gleich, Entscheid 9).
- Tausch A ↔ B in einem Commit ⇒ kein Pfad nimmt teil (Quelle und Ende zugleich).
- Beide Schalter, umbenannte Quelle, Nachzug auf eine gleichnamige andere Datei ⇒ Befund
  (Entscheid 10).
- Ende der Kette in HEAD gelöscht ⇒ keine Umbenennung.
- Nur das Ausführbar-Bit geändert ⇒ gleich; Symlink gegen Datei mit gleichem Blob ⇒
  nicht gleich (Entscheid 13).
- Link auf ein umbenanntes Verzeichnis ⇒ wie heute; Link auf eine Datei darin ⇒ folgt.
- Die immutable Datei selbst umbenannt ⇒ Befund wie heute.
- `--staged` mit `git mv` und Nachzug im Index ⇒ kein Befund; `git mv` in einem früheren
  Commit, nur der Nachzug gestagt ⇒ Befund (benannte Grenze).
- Unlesbares Objekt eines Commits der Range ⇒ Exit 2.
- Gitlink im Index unter `vcs.paths` mit `--staged` ⇒ kein Kandidat (wie im Baum); ohne
  die Angleichung wäre er einer (Rot-Beleg am Vorzustand).
- Zehn Läufe derselben Eingabe ⇒ byte-gleiche Ausgabe (Messmethode von
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md)
liegt in `done/` (Anforderung, Spezifikation und ADR stehen, die Entscheide sind bestätigt),
das Architect-Verdikt zu diesem Plan liegt unter `docs/reviews/`, und das WIP-Limit ist
frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Port und Adapter brauchen
  allein mehr als eine Review-Sitzung — dann wird der Adapter-Teil ein eigener Slice
  derselben Welle vor dem Regel-Teil.
- `in-progress` → `open` (blockiert): go-git liefert für den Index oder für ältere
  Pack-Formate die Blob-Kennungen nicht so, wie die Inhalts-Gleichheit sie braucht — dann
  geht die Frage an die Anforderung zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Port, Adapter und Regel mit Tests je Kriterium und gezeigtem Rot, Spiegel
grün, ADR `Accepted`, `make gates` grün), Review- und Verifikations-Report liegen unter
`docs/reviews/`, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Auflöser der Link-Ziele trug schon für `ignore-link-targets` acht Review-Runden
  (Reports `docs/reviews/2026-10-09-slice-267-vcs-pfad-nachzug-r1.md` bis `-r8.md`); eine
  zweite Gleichsetzung an derselben Stelle kann die Filter dort verschieben
  ([`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md))
  — **Ausgang:** *entfallen*, wenn die bestehenden Tests von
  `vcs_pfad_nachzug_test.go` unverändert grün bleiben und ein Test beide Schalter
  zusammen fährt.
- git-Fixtures, die Inhalt gleicher Länge in einem Zug neu schreiben, fallen sporadisch
  aus ([`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md))
  — Umbenennungs-Fixtures mit gleichen Blobs liegen nahe an dieser Klasse —
  **Ausgang:** *entfallen*, wenn jedes neue Fixture über die Helfer von
  `internal/hexagon/core/coretest/` läuft und keines dieselbe Datei gleich lang neu
  schreibt.
- Der Commit-Walk liest mehr Objekte als heute; ein großer Push-Bereich in der CI wird
  langsamer — **Ausgang:** *entfallen*, wenn `make adr-check` über eine Range von
  100 Commits dieses Repos mit Schalter höchstens doppelt so lange läuft wie ohne
  (Messung in der Closure-Notiz); sonst *eingetreten* mit Folge-Slice.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist der Default `*` (`ALL`) — Port,
git-Adapter, Regel, Konfigurations-Rand, Doku; keine Sub-Area mit eigener
Modus-Deklaration (`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-23), Sub-Area `*`. Treffer:
[`BEO-ALL/shared-lexicon-drifts-at-edges`](../observations/BEO-ALL/shared-lexicon-drifts-at-edges/observation.md)
(verkörpert in welle-74, Restklasse offen — Risiko 1 in §6),
[`BEO-ALL/racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)
(offen — Risiko 2 in §6),
[`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
(4× — die Spiegel-Suche in §3 nimmt Handbuch, READMEs und Sensor-Datei auf),
[`BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet`](../observations/BEO-ALL/erkennungs-regex-nur-gegen-positivfaelle-getestet/observation.md)
(7× — die Negativliste in §3),
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/observation.md)
(offen, 1× — Welle-Out-of-Scope) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — der Port liest Bäume
außerhalb der Scan-Menge; die Zusage dort schreibt
[`slice-291`](../in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) fest). Der Zähler
folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-23: `image-scan` grün (2026-11-23T12:51:15Z),
`upstream-drift` rot (2026-11-23T08:29:02Z, vor den Pin-Hebungen desselben Tages); beim
Start dieses Slice neu lesen.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
