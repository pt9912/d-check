# Welle welle-96: `vcs` folgt einer Umbenennung — ein Link-Nachzug auf die umbenannte Datei ist keine Drift

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-11-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Eine `Accepted`-ADR, die einen Link auf eine Datei nur nachzieht, weil diese Datei in
derselben Commit-Range **umbenannt** wurde, ist für das Modul `vcs` keine Drift mehr —
opt-in, ohne den Schalter byte-identisch. Heute geht mit `vcs.ignore-link-targets` nur
ein Nachzug auf eine Datei **gleichen Namens** durch; ein Nachzug auf einen neuen Namen
meldet `core-drift-vcs`, obwohl die Baseline ihn ausdrücklich erlaubt (Modul 4,
§Nachzug ist keine Überschreibung). Die Identität „dieselbe Datei unter neuem Namen“
liest das Modul aus der git-Historie der Range; ein Umbiegen auf eine andere,
bestehende Datei, ein Löschen-und-neu-Anlegen ohne erkannte Umbenennung und jede
Änderung an Anker, Linktext oder übriger Zeile bleiben Drift. Der Lauf bleibt
netzlos, liest `.git` nur lesend und schreibt nichts.

Anlass ist [Issue #18](https://github.com/pt9912/d-check/issues/18) aus dem Kurs-Repo:
eine Layout-Migration mit englischen Dateinamen (im Verzeichnis `spec/` wird `spezifikation.md` zu
`specification.md`) trifft jede `Accepted`-ADR, die die Datei verlinkt — im
eigenen Bestand 59 von 107 (Zählung 2026-11-23, Kommando in
[`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) §3).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Der Auftraggeber hat [Issue #18](https://github.com/pt9912/d-check/issues/18) als
  Welle angefordert (2026-11-23, in der Sitzung).
- `origin/main` und der lokale Hauptzweig stehen gleich (`git fetch`, Abstand 0/0),
  `next/` ist leer, `in-progress/` trägt nur die Roadmap, `open/` nur die zwei Slices
  dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle sind `done/`.
- `make fullbuild` grün.
- **Das Mehr gegenüber den Slice-DoDs** (repo-weit, in keiner einzelnen DoD):
  - **Ohne den Schalter nichts Neues:** `make blackbox-probe` gegen das Vorher-Image
    der letzten Release-Version (`REF=v0.92.0`) zeigt für den Eigenlauf und ein
    Fixture ohne den neuen Schalter **dieselbe** Ausgabe (stdout, stderr und Exit
    getrennt). Eine opt-in-Erweiterung darf keinen bisher grünen oder roten Lauf
    verschieben.
  - **Der Anlass trägt am gebauten Image:** in einem Wegwerf-Klon dieses Repos
    benennt ein Commit `spec/spezifikation.md` per `git mv` um und zieht in jeder
    `Accepted`-ADR nur das Link-Ziel nach. `make adr-check RANGE=HEAD~1..HEAD` meldet
    ohne den Schalter je ADR einen `core-drift-vcs`, mit ihm keinen; die Gegenprobe —
    derselbe Nachzug, die Datei aber gelöscht und mit anderem Inhalt neu angelegt —
    meldet wieder. Die Tabelle Probe → Ergebnis steht in `welle-96-results.md`.
- Das Trigger-Audit (seit slice-241) ist über die Welle-Closure vollzogen und in
  `welle-96-results.md` belegt.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) | Anforderung, Spezifikation und ADR: `vcs` folgt einer Umbenennung; Entscheide festschreiben | [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) |
| [`slice-292`](open/slice-292-vcs-umbenennung-modulkern.md) | Kern und git-Adapter: Umbenennungen der Range lesen, Link-Ziele über sie gleichsetzen, Schalter | [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in) in der Fassung aus [`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle. Eine spätere Baseline-Hebung, die Dateinamen der Spec
  umbenennt, läuft ohne diese Welle mit rotem `make adr-check` in der CI.
- Wird blockiert von: keiner.
- **Reihenfolge der Slices ist fachlich bindend:** erst
  [`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) (Anforderung und
  Entscheide stehen, bevor Code entsteht), dann
  [`slice-292`](open/slice-292-vcs-umbenennung-modulkern.md). Das WIP-Limit 1 hält sie
  ohnehin nacheinander.
- **Entscheide, die [`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md)
  festschreibt** (mit Empfehlung; ohne Entscheid des Auftraggebers beginnt kein Code):
  woran eine Umbenennung erkannt wird (Inhalts-Gleichheit oder Ähnlichkeits-Schwelle),
  ob nur die Endpunkte der Range oder jeder Commit gelesen wird (und damit, ob
  A → B → C zählt), eigener Schalter oder Erweiterung von `ignore-link-targets`, was mit
  mehrdeutigen Umbenennungen und Verzeichnis-Zielen geschieht, ob das eigene Profil den
  Schalter einschaltet.
- **Entschieden mit dem Auftraggeber (2026-11-23):** die acht Entscheide in
  [`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) §3 mit den
  Empfehlungen bestätigt; nach dem Architect-Verdikt am selben Tag auch Entscheid 5 in
  der seiten-gleichen Fassung, die Folge-Entscheide 9–14 aus
  [`ADR-0115`](../adr/0115-vcs-folgt-umbenennung.md) und die Schärfung des
  Out-of-Scope-Satzes von
  [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in).
  Nach Review R2 am selben Tag bestätigt: Entscheidung 9 und 10 in der neuen Fassung
  (beide in der Sache geändert); Entscheidung 9 am 2026-11-24 erneut bestätigt, mit den
  drei Fällen, in denen sie nachsichtiger ist als die zuerst bestätigte Fassung
  ([`slice-291`](in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md) §3).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Die Umbenennung der immutablen Datei selbst** — eine umbenannte `Accepted`-ADR
  bleibt `core-drift-vcs` („der Pfad einer immutablen Datei ist stabil“); das Issue
  fragt nicht danach.
- **Ein Nachzug auf eine gleichnamige andere Datei** — die offene Grenze von
  `ignore-link-targets`
  ([`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/observation.md))
  wird nicht geschärft; eine Verschärfung wäre nicht opt-in-neutral und ein eigener
  Vorgang. **Eine Ausnahme folgt aus Entscheid 10:** sind beide Schalter an, gleicht ein
  Ziel, das an einer Umbenennung teilnimmt, keinem nicht teilnehmenden — ein Nachzug
  zwischen beiden ist Drift, auch bei gleichem Dateinamen; zwischen zwei nicht
  teilnehmenden Zielen bleibt die Grenze, wie sie ist.
- **Die Umbenennung der eigenen Spec-Dateien** — gehört zur Baseline-Hebung, die sie
  verlangt, nicht zu dieser Welle.
- **Andere Module, die git lesen** (`metrics`, `tracked`, `commits`) — `metrics` zählt
  eine Umbenennung weiter als Löschung plus Einfügung.
- **Release, Tag, Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins) und
  die Antwort auf das Issue — eigener Vorgang nach der Welle
  ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-96-results.md` (Geschwister im Ruheort `done/`, entsteht bei der Closure)
Zähler: [Beobachtungs-Register](observations/README.md)
