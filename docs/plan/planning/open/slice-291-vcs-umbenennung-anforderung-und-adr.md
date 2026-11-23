# Slice slice-291: Anforderung, Spezifikation und ADR — `vcs` folgt einer Umbenennung; Entscheide festschreiben

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
(das Kriterium „Boundary (Pfad-Nachzug)“ wird geändert, ein neues kommt hinzu),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)
(Scope), [`ADR-0024`](../../adr/0024-vcs-immutable-gate.md),
[`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).

**Berührte Spec-Stellen:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(Absatz „Pfad-Nachzug (opt-in)“, Akzeptanzkriterien, Version und Historie);
`spezifikation.md` §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 und 4 (Lesen der Umbenennungen,
Normierung) und §2 (Schema `vcs.*`); `architecture.md` bleibt unberührt (die Rolle
„git-Port“ besteht). Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

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

**Ziel:** Die geänderte Anforderung
[`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
steht im Lastenheft, Verfahren und Schema in der Spezifikation, die Entscheidung in
einer ADR (`Proposed` bis zur Closure von
[`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md)) — und die Entscheide unten sind
**vor dem ersten Code** vom Auftraggeber bestätigt und festgeschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, git-Adapter, Konfigurations-Rand, Tests** — gehören
  [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md); dieser Slice schreibt, was
  der Code einlösen soll, damit die Reviews dort gegen eine Anforderung prüfen.
- **Das eigene Profil (`.d-check.yml`)** — der Schalter wird erst wirksam, wenn der
  Code steht; das Einschalten gehört, falls Entscheid 7 es bestätigt, zu
  [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md).
- **Die Grenze „gleichnamige andere Datei“ von `ignore-link-targets`** — bleibt, wie
  sie ist (Welle-Out-of-Scope); die neue Anforderung nennt sie nur als bestehende
  Grenze.
- **CHANGELOG, READMEs, Handbuch** — Release-Prep
  ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 17), kein Feature-Commit.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
      ist im Lastenheft geändert (Version gehoben, Historie-Zeile nach
      [`MR-032`](../../../../harness/conventions.md#mr-032)): der Absatz zum Pfad-Nachzug
      nennt das Folgen einer Umbenennung samt Grenzen, das Kriterium „Boundary
      (Pfad-Nachzug)“ widerspricht ihm nicht mehr, und je Entscheid steht ein
      Akzeptanzkriterium — mindestens die vier aus dem Issue, die Kette, die
      Mehrdeutigkeit und „ohne Schalter byte-identisch“.
- [ ] Die Spezifikation führt das Verfahren (welche Commits und Bäume gelesen werden,
      woran eine Umbenennung erkannt wird, wie ein Ziel über sie gleichgesetzt wird,
      fail-closed-Fälle) und den Schlüssel im Schema; jede Aussage über die Grenzen
      (welche Eingaben das Modul liest, [`AGENTS.md`](../../../../AGENTS.md) §3.8) ist
      am Gegenstand gezählt oder gemessen, nicht aus der Beschreibung übernommen.
- [ ] Die ADR (Entscheide samt verglichenen Alternativen, darunter „Ähnlichkeits-Schwelle“
      und „Erweiterung von `ignore-link-targets`“) liegt als `Proposed` vor, ist im
      [ADR-Index](../../adr/README.md) geführt und trägt `## Re-Evaluierungs-Trigger`.
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
| `spec/lastenheft.md` | update | [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in): Absatz, Kriterien, Version und Historie |
| `spec/spezifikation.md` | update | §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 und 4, Schema-Zeile des neuen Schlüssels, Historie |
| ADR im Verzeichnis `docs/plan/adr/` (nächste freie Nummer, Titel „vcs-folgt-umbenennung“) | neu | Entscheidung samt Alternativen; `Proposed` |
| `docs/plan/adr/README.md` | update | ADR-Index ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 4) |
| Plan von [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md) | update | die Testfälle, die aus den bestätigten Entscheiden folgen |

**Anlass gezählt** (Regel 14/18 der Dokumentations-Regeln — Form vor Zahl): eine
`Accepted`-ADR „verlinkt“ die Spezifikation, wenn eine Zeile ein Markdown-Link-Ziel
`(…spec/spezifikation.md…)` trägt. Am 2026-11-23: 59 von 107 `Accepted`-ADRs;
66 nennen den Dateinamen in irgendeiner Form. Kommando:

```bash
acc=$(grep -l '^\*\*Status:\*\* Accepted' docs/plan/adr/[0-9]*.md)
echo "$acc" | wc -l                                                  # 107
echo "$acc" | xargs grep -lE '\]\([^)]*spec/spezifikation\.md' | wc -l  # 59
```

**Festzuschreibende Entscheide** — vom Auftraggeber am 2026-11-23 mit den
Empfehlungen bestätigt (unverändert 1–8):

1. **Woran eine Umbenennung erkannt wird: Inhalts-Gleichheit, keine
   Ähnlichkeits-Schwelle.** Ein Pfad A, der vor einem Schritt existiert und danach
   nicht, und ein Pfad B, der danach existiert und vorher nicht, sind eine
   Umbenennung A → B genau dann, wenn beide dasselbe Blob tragen (git: `R100`).
   Verworfen: eine Ähnlichkeits-Schwelle (git-Default 50 %; die Vorbelegung der
   Rename-Erkennung von go-git ist am Gegenstand nachzulesen, nicht zu übernehmen) —
   sie ließe „umbenennen und umschreiben“ durch, und genau das ist das Umbiegen auf
   ein anderes Dokument, das die Anforderung als Drift hält; außerdem wäre die
   Schwelle ein Parameter, den niemand entschieden hat. Die eigene Hard Rule
   ([`AGENTS.md`](../../../../AGENTS.md) §3.3) macht den Move-Commit rein — die
   Form, die Inhalts-Gleichheit erkennt.
2. **Was ein „Schritt“ ist: jeder Commit der Range, mit Kette.** Gelesen wird jeder
   Nicht-Merge-Commit, der von HEAD aus erreichbar ist und von BASE aus nicht
   (dieselbe Menge wie das Modul `commits`), je gegen seinen Eltern-Commit; im
   `--staged`-Modus ist der einzige Schritt HEAD → Index. Umbenennungen verketten
   sich (A → B, später B → C ergibt A → C). Verworfen: nur die Endpunkte BASE und
   HEAD vergleichen — nach §3.3 folgt die Inhaltsänderung der umbenannten Datei als
   eigener Commit, oft im selben Push; an den Endpunkten trügen A und C dann
   verschiedene Blobs, und der Nachzug bliebe rot, obwohl jeder einzelne Schritt
   sauber ist.
3. **Mehrdeutigkeit bleibt Drift.** Tragen in einem Schritt mehrere verschwundene
   oder mehrere neue Pfade dasselbe Blob, oder führt die Kette von A zu zwei
   verschiedenen Enden, gilt für A keine Umbenennung. Ein Merge-Commit trägt selbst
   keine Umbenennung bei (seine Seiten stehen als eigene Commits in der Menge).
4. **Eigener Schalter `vcs.follow-renames`** (bool, Default `false`), unabhängig von
   `vcs.ignore-link-targets`; aus ⇒ byte-identisch
   ([`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)). Verworfen:
   die Bedeutung von `ignore-link-targets` erweitern (Vorschlag 1 des Issues) — das
   lockerte das Gate jedes Repos, das den Schlüssel schon setzt, ohne dass es darum
   gebeten hat ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
5. **Was gleichgesetzt wird: nur das aufgelöste Ziel.** Ein Link-Ziel, das in BASE
   auf den Pfad A auflöst, gilt als Ziel auf B, wenn die Range A → B führt; mit HEAD
   verglichen wird der aufgelöste Pfad samt Anker, sonst nichts. Anker, Linktext,
   Titel, übrige Zeile und Zahl der Links bleiben roh wie heute; ein Ziel, das an
   keiner Umbenennung teilnimmt, behandelt das Modul wie bisher (roh oder, mit
   `ignore-link-targets`, auf Dateiname und Anker normiert). Was ein Link ist,
   beantwortet weiter die Erkennung des Moduls `links`. Nur **Datei**-Ziele folgen —
   ein Verzeichnis trägt kein Blob; eine Verzeichnis-Umbenennung erreicht die Links
   auf Dateien darin über deren eigene Umbenennung, ein Link auf das Verzeichnis
   selbst bleibt wie heute.
6. **Grenzen und fail-closed** ([`AGENTS.md`](../../../../AGENTS.md) §3.8): das Modul
   liest dafür Bäume jedes Commits der Range, über die Scan-Menge und `vcs.paths`
   hinaus — eine Umbenennung irgendwo im Repo zählt. Die Zusage gilt dort nur für
   Inhalts-Gleichheit, nicht für „derselbe Gegenstand“. Ein unlesbarer Commit oder
   Baum ist ein Fehler (Exit 2), keine stille Lücke — dieselbe Politik wie
   `AllPaths`. Kein neuer Grund-Code: es bleibt `core-drift-vcs`. Die Umbenennung
   der immutablen Datei selbst bleibt Drift.
7. **Eigenes Profil:** `.d-check.yml` schaltet `follow-renames: true` im `vcs`-Block
   ein ([`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md)). Begründung:
   [`ADR-0103`](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) hat für dieses
   Repo schon entschieden, dass ein Pfad-Nachzug keine Drift ist; die neue ADR trägt
   die Ausweitung auf den neuen Namen (Lockerung nur per ADR,
   [`AGENTS.md`](../../../../AGENTS.md) §3.6). Alternative: erst mit der
   Baseline-Hebung einschalten, die die Spec-Dateien umbenennt.
8. **ADR-Status:** die ADR bleibt `Proposed`, bis der Mechanismus in
   [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md) steht und reviewt ist
   ([`BEO-ALL/adr-accepted-bevor-der-mechanismus-steht`](../observations/BEO-ALL/adr-accepted-bevor-der-mechanismus-steht/observation.md),
   2×).

**Gegenprobe vor dem Schreiben (Schritt 19 des Workflow-Skeletts):** die Negativfälle
stehen als Testfälle für [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md) —
Umbiegen auf eine bestehende Datei C; Löschen und neu Anlegen mit anderem Inhalt;
umbenennen und im selben Commit ändern; zwei Kopien desselben Blobs; Kette über drei
Commits; Kette mit Verzweigung; Anker- oder Linktext-Änderung zusammen mit dem
Nachzug; Link auf ein umbenanntes Verzeichnis; Umbenennung in einem Merge; die
immutable Datei selbst umbenannt; `--staged` mit `git mv` und Nachzug im Index.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Welle `welle-96-vcs-folgt-umbenennung` ist offen,
das Architect-Verdikt zu den Entscheiden liegt unter `docs/reviews/`, und das WIP-Limit
ist frei (`in-progress/` trägt nur die Roadmap).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Commit-Kette aus
  Entscheid 2 verlangt eine eigene Anforderung an den git-Port, die über
  [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  hinausgeht.
- `in-progress` → `open` (blockiert): der Auftraggeber entscheidet einen der Punkte
  1, 2 oder 4 nicht; ohne Entscheid beginnt kein Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Lastenheft, Spezifikation und ADR liegen vor und `make gates` ist grün),
Review- und Verifikations-Report liegen unter `docs/reviews/`, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Das Kriterium „Boundary (Pfad-Nachzug)“ sagt heute ausdrücklich, ein Nachzug auf
  eine Datei mit anderem Namen bleibe Drift; die Änderung muss es umschreiben, nicht
  daneben ein zweites, widersprechendes stellen
  ([`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md),
  2×) — **Ausgang:** *entfallen*, wenn das Review Absatz, Kriterium und ADR
  nebeneinander gelesen und keinen Widerspruch gefunden hat; sonst *eingetreten* mit
  dem Fix-Commit als Beleg und einer Evidenz-Datei im Register.
- Der gestagte Check sieht nur den Index gegen HEAD: liegt der `git mv` in einem
  früheren Commit und wird nur der Nachzug gestagt, fehlt die Umbenennung in der
  Range, und der Hook bleibt rot
  ([`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/observation.md),
  1×); die Anforderung muss das als Grenze nennen, die CI-Range sieht beide Commits —
  **Ausgang:** *weiter offen*, wenn die Anforderung die Grenze nennt (sie ist eine
  Eigenschaft des gestagten Modus, kein Fehler dieses Slice); *entfallen*, wenn der
  Architect einen Weg bestätigt, der den früheren Move-Commit mitliest.
- Die Commit-Menge „von HEAD erreichbar, von BASE nicht“ ist bei einem Force-Push
  oder einer Basis, die kein Vorfahr ist, eine andere als erwartet
  ([`ADR-0109`](../../adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md) regelt
  die CI-Range) — **Ausgang:** *entfallen*, wenn die Spezifikation die Commit-Menge
  und den Fall „Basis kein Vorfahr“ (fail-closed, Exit 2) festschreibt und
  [`slice-292`](../open/slice-292-vcs-umbenennung-modulkern.md) einen Test dafür trägt.

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
Lastenheft, Spezifikation, ADR; keine Sub-Area mit eigener Modus-Deklaration
(`tools/harness/` bleibt unberührt).

<!-- d-check:cite .harness/baseline/v6.18.1/regelwerk/modul-05-planning-harness.md:379-379 -->

> **Offene Beobachtungen sichten.** Das

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand
2026-11-23), Sub-Area `*`. Treffer:
[`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/observation.md)
(offen, 1× — dieselbe Achse, bleibt Welle-Out-of-Scope; die Anforderung nennt die
Grenze weiter),
[`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/observation.md)
(offen, 1× — Risiko 2 in §6),
[`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)
(offen, 2× — Risiko 1 in §6; ein dritter Beleg wäre die Schwelle),
[`BEO-ALL/adr-accepted-bevor-der-mechanismus-steht`](../observations/BEO-ALL/adr-accepted-bevor-der-mechanismus-steht/observation.md)
(offen, 2× — Entscheid 8),
[`BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft`](../observations/BEO-ALL/neue-anforderung-gegen-globalen-ausschluss-ungeprueft/observation.md)
(offen — die geänderte Anforderung wird gegen die globalen Ausschlüsse des Lastenhefts
gehalten, das Review liest sie mit),
[`BEO-ALL/externer-link-auf-umgezogene-datei`](../observations/BEO-ALL/externer-link-auf-umgezogene-datei/observation.md)
(offen, 1× — betrifft die spätere Spec-Umbenennung, nicht diese Welle) und
[`BEO-ALL/module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
(verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.8 — Entscheid 6 beantwortet die
Frage). Der Zähler folgt aus den Dateien der Closure.

**Vorgelagert — Nachtlauf-Stand lesen** ([`MR-053`](../../../../harness/conventions.md#mr-053)):
`make nightly-state` am 2026-11-23: `image-scan` **grün** (jüngster Lauf
2026-11-23T12:51:15Z, nach dem Release v0.92.0). `upstream-drift` **rot** (jüngster Lauf
2026-11-23T08:29:02Z) — der Lauf liegt **vor** den Pin-Hebungen desselben Tages (Go
1.27.4, Semgrep 1.186.0, a-check v0.24.0, `docker/login-action` 4.8.0); planmäßige
Fremd-Drift, die der nächste Lauf neu bewertet. Berührt dieses Modul nicht.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
