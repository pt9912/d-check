# Welle 86 — Der Closure-Übergang trägt seine Vorbedingungen — Closure-Notiz

**Welle:** welle-86-closure-uebergang-durchsetzen
**Abschluss:** 2026-09-03
**Verantwortlich:** pt9912

## Was wurde geliefert?

Fünf Slices lösen die vier Vorbedingungen ein, die welle-86 §1 als **offen**
auswies — DoD-Häkchen gesetzt, Beobachtungs-Register fortgeschrieben, Review
fand statt, die Prüfung hängt am Übergang, nicht am Zustand:

- [slice-172](welle-86/slice-172-closure-uebergang-waechtern.md): eine
  `structure`-Regel (`max-open-tasks: 0`) meldet jeden offenen DoD-Haken
  eines `done/`-Slice, mit Bestands-Ausnahme bis `slice-170`.
- [slice-173](welle-86/slice-173-review-report-deckung.md): das neue Modul `reviews`
  ([`DC-FA-RVW-001`](../../../../spec/lastenheft.md#dc-fa-rvw-001--review-report-deckung-modul-reviews-opt-in),
  [ADR-0081](../../adr/0081-reviews-modul.md)) — jeder `done/`-Slice mit
  Review-Zusage braucht mindestens einen Report unter `docs/reviews/`.
- [slice-174](welle-86/slice-174-register-deckung.md): die maschinelle Hälfte der
  Register-Paarung — eine zitierte `BEO-<NNN>` hat eine Registerzeile.
- [slice-175](welle-86/slice-175-uebergangs-waechter.md): `.githooks/pre-commit` und
  die PR-/Push-CI erkennen einen Slice-Closure-Übergang (Rename/Add nach
  `docs/plan/planning/done/slice-*.md`) und lösen `make verify-closure-notes`
  aus.
- [slice-192](welle-86/slice-192-uebergangswaechter-nachliefern.md): welle-86s
  eigene Trigger-Prüfung fand, dass zwei der vier Vorbedingungen beim
  lokalen `mv`-Commit noch nicht griffen (`.d-check.closure.yml` kannte
  weder `planning.observations` noch `reviews`) — behoben,
  [ADR-0082](../../adr/0082-uebergangswaechter-reviews-observations.md).

**Alle vier Vorbedingungen sind jetzt real geprüft, nicht nur beschrieben:**
vier isolierte Proben, je eine pro Vorbedingung, real als `git commit`-Versuch
gefahren und abgewiesen (Belege in Verifikation unten).

## Was hat funktioniert?

**Die Trigger-Prüfung selbst als Werkzeug, nicht als Formsache.** Eine
Behauptung „vier Vorbedingungen sind jetzt durchgesetzt" hätte die Lücke aus
slice-192 nicht gefunden — erst der reale `git commit`-Versuch je
Vorbedingung zeigte, dass zwei von ihnen nur über `make gates` bzw. gar
nicht liefen. Dieselbe Disziplin, mit der slice-175 sein eigenes Ergebnis
bewies (zwei reale Proben, sofort zurückgenommen), trug auch die
Welle-Schlussprüfung.

**Wiederverwendung statt neuer Module.** Jede der fünf Lücken wurde mit
einer bereits vorhandenen Fähigkeit geschlossen (`structure`, `planning`,
`reviews`) — die einzige neue Prüf-Logik der ganzen Welle ist das Modul
`reviews` selbst (slice-173); alles Übrige ist Bindung an einen neuen
Auslöser.

**Der Wächter fing sich selbst.** slice-192s eigener erster `mv`-Versuch
wurde vom gerade geschärften Wächter abgewiesen — derselbe Fund wie bei
slice-173/175 (eine `unabhängiger Review`-DoD-Zusage ohne persistierten
Report), diesmal sofort statt erst später gemessen.

## Was ging anders als geplant?

- **Der Scope wuchs von vier auf fünf Slices**, durch die eigene
  Trigger-Prüfung selbst ausgelöst — kein Vorsatz, sondern Messung
  (Modul 6 Schritt 2, Trigger-Audit:
  [ADR-0081](../../adr/0081-reviews-modul.md)s eigener
  Re-Evaluierungs-Trigger „Aufnahme in `gates`" traf genau zu, weil
  welle-86 der volle Wellen-Zyklus war, den
  [ADR-0081](../../adr/0081-reviews-modul.md) selbst benannte).
- **Drei Review-Reports fehlten rückwirkend.** slice-173, slice-175 und
  slice-192 selbst hatten alle eine `unabhängiger Review`-DoD-Zusage ohne
  persistierten Report unter `docs/reviews/` — die Reviews fanden real
  statt (unabhängige Sub-Agenten-Läufe, Befunde eingearbeitet), nur nie als
  Artefakt festgehalten. Alle drei wurden aus den tatsächlichen Befunden
  nachgezogen, nicht rückwirkend erfunden.
- **Zwei Werkzeug-Fehler wurden vom unabhängigen Review vor der Closure
  gefunden** (slice-173: Mehrzeilen-DoD-Items übersehen, redundante
  Fail-Closed-Meldung; slice-175: `awk`-Feldtrennung auf Whitespace statt
  Tab) — beide behoben, mit Regressionstests bzw. Proben belegt.

## Steering-Loop-Einträge

Keine Beobachtung erreichte während dieser Welle neu die 3×-Schwelle.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations.md`](../observations.md). Kein
Eintrag erreichte während dieser Welle die 3×-Schwelle neu.

## Folge-Slices

Keiner verbindlich. Ein Stop-Hook zur früheren Erkennung eines
Übergangs-Verstoßes (vor dem Commit-Versuch, nicht erst bei ihm) bleibt ein
benannter, unabhängiger Kandidat — slice-175 §3 grenzt ihn ausdrücklich aus,
weil der git-Hook plus CI die Repo-Invariante bereits vollständig trägt.

## Verifikation

- `make gates` Exit 0 (zehn Glieder) auf jedem der Commits dieser Welle.
- `make fullbuild` Exit 0 — 51 Requirements, 0 Waisen — Image-Hash
  `sha256:2c3967317b040827ec56673e952752156508425da527061bbd1a64242ea23c95`.
- **Vier reale Proben, isoliert gefahren, jede als eigener `git commit`-Versuch:**
  - DoD-Häkchen offen ⇒ `section-tasks-open`, abgewiesen.
  - Zitierte, nicht registrierte `BEO-<NNN>` ⇒ `observation-unregistered`,
    abgewiesen.
  - Review-Zusage ohne Report ⇒ `review-missing`, abgewiesen.
  - Die Bindung an den Übergang selbst — durch die Form aller drei Proben
    bewiesen (echte `git commit`-Versuche, nicht `make verify-closure-notes`
    im Leerlauf) sowie durch slice-175s eigene Positiv-/Negativ-Probe
    (Commit `07afe62`, sofort zurückgenommen).
  Alle vier Proben wurden nach dem Lauf sofort per `git reset` entfernt;
  kein Test-Artefakt reist in die echte Historie.
- Drei unabhängige Reviews (slice-173, slice-175, slice-192) und drei
  unabhängige Verifikationen, alle Befunde eingearbeitet.
