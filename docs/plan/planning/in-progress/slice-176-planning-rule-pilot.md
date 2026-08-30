# Slice slice-176: Der Zyklus ist im Kontext, bevor jemand ihn braucht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht). Der **Anlass** liegt in
[welle-86](../welle-86-closure-uebergang-durchsetzen.md); dort ist die
Zustellung ausdrücklich **nicht** der Closure-Trigger.

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §1 (das Prinzip: Hard Rules und
**Pointer**, keine Duplikation);
[`MR-021`](../../../../harness/conventions.md#mr-021) (der Bump zieht
vendored-Pfad-Verweise auf den neuen Tag — hier kommt ein neuer Träger dazu);
[`MR-043`](../../../../harness/conventions.md#mr-043) (der Werkzeug-Einstieg
importiert `AGENTS.md`);
[`MR-042`](../../../../harness/conventions.md#mr-042) (die Präzedenz für einen
**werkzeug-lokalen** Träger und seine Nicht-Zusagen).

**Berührte Spec-Stellen:** — (Werkzeug-Konfiguration; keine Produkt-Anforderung
berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Die Regeln, die eine Sitzung gekostet haben, standen nicht im Kontext.**

In der auslösenden Sitzung liefen drei Slices (168, 169, 170) durch, ohne dass
der Zyklus (`Spec → ADR → Plan → Code → Review → Verifikation → Closure`) im
Kontext war. Er steht in `modul-01`, das niemand geöffnet hatte — und in
`AGENTS.md` steht er nicht, korrekterweise: §6 beschreibt den
**Implementer**-Workflow, und der endet bei Schritt 8. Folge: Review und
Verifikation fielen aus, das Beobachtungs-Register blieb unangetastet, drei
Slices gingen mit offenem DoD-Haken nach `done/`.

**`AGENTS.md` kann die Lücke nicht schließen, indem es wächst.** Die Datei
trägt **527** Zeilen gegen 236 der Baseline-Vorlage — und ist damit selbst schon
der Fall, vor dem die Werkzeug-Dokumentation warnt:

> target under 200 lines per CLAUDE.md file. Longer files consume more context
> and reduce adherence.

**Der Import löst es nicht.** [`CLAUDE.md`](../../../../CLAUDE.md) zieht
`AGENTS.md` über `@AGENTS.md` in jeden Lauf — Imports helfen der
*Organisation*, nicht dem Kontext.

## Zwei Träger wurden geprüft, einer davon verworfen

**Pfad-gebundene Regeln (`.claude/rules/` mit `paths:`) tragen nicht.** Nach
Werkzeug-Auskunft (2026-08-29) hängt eine solche Regel im **Auto-Modus** nur
ein, wenn die Datei über die **dedizierten** Werkzeuge (`Read`, `Edit`, `Write`)
angefasst wird; jeder Zugriff über die Shell — `cat`, `sed`, `awk` — geht daran
vorbei. Dieses Repo arbeitet überwiegend über die Shell: **im Anlassfall hätte
die Regel nicht geladen.** Das ist die Klasse
[`BEO-024`](../observations.md) — ein Zustell-Kanal, der an der **Arbeitsweise**
hängt statt am Inhalt.

**Ein Hook könnte es**, deterministisch und werkzeug-unabhängig — er ist als
Alternative in §3 geführt und bleibt der spätere Optimierungsschritt. Er bringt
aber bewegliche Teile mit: eine Pfad-Heuristik über die Bash-Nutzlast, die
Entscheidung `PreToolUse` gegen `PostToolUse`, einen eigenen Glob-Dialekt und
eine JSON-Extraktion ohne `jq` ([`AGENTS.md`](../../../../AGENTS.md) §3.1).
Jedes davon kann **still** danebengreifen — und still danebengreifen ist genau
die Klasse, gegen die dieser Slice gebaut ist.

## 2. Vorgehen

**Vier Symlinks, kein Code.** `.claude/rules/` bekommt Symlinks auf die vier
Regelwerk-Module, deren Fehlen den Anlassfall verursacht hat. Sie tragen **kein**
`paths:` und laden deshalb **immer** — beobachtet in einem Schwester-Projekt,
das den ganzen Regelwerk-Baum so einhängt (Auskunft des Auftraggebers,
2026-08-30).

1. **Die Auswahl folgt dem Anlassfall, nicht dem Vollständigkeits-Wunsch:**
   `modul-01` (der Zyklus, dessen Fehlen Review und Verifikation kostete),
   `modul-05` (Planning-Form und Vorprüfungen), `modul-06` (Roadmap und
   Beobachtungs-Register — die zweite ausgefallene Pflicht), `modul-08`
   (Rollen-Trennung). **Gemessen: 805 Zeilen** von 4787. Die übrigen 22 Module
   bleiben draußen, bis ihr Fehlen belegt ist.
2. **Symlink statt Kopie.** Eine Kopie wäre eine zweite Wahrheitsquelle und
   driftete beim Bump still; der Symlink zeigt auf den **gepinnten** Baum, und
   dessen Inhalt hält `make baseline-verify` per SHA — Drift ist dort
   ausgeschlossen.
3. **Der Preis der Symlink-Form ist benannt: die Dateien werden nicht
   gate-geprüft.** Gemessen an einem Probe-Repo: eine **echte** Datei unter
   `.claude/rules/` wird gescannt und ihr toter Link gemeldet, ein **Symlink**
   nicht. Für diesen Fall wiegt das wenig — das Ziel liegt im SHA-gepinnten
   Baum, `links`/`anchors` darüber wären Doppelarbeit. **Es wiegt genau an einer
   Stelle:** beim Bump.
4. **Der Bump bekommt den neuen Träger — als Sensor, nicht als Prozedur.**
   [`MR-021`](../../../../harness/conventions.md#mr-021) bindet *Markdown-Links*
   auf den gepinnten Baum und überlässt die Menge dem Zensus der Bump-Prozedur.
   Ein Symlink ist keiner: ein Zensus, der nach Markdown-Links sucht, findet
   ihn nicht, und beim Bump bricht er **still** — `links` sieht ihn nicht
   (nicht gescannt), `sha256sum -c` und die Manifest-Deckung bleiben grün.
   `make baseline-verify` prüft deshalb als **dritte** Frage, dass jeder
   Symlink unter `.claude/rules/` auflöst; es läuft in `make gates`. Geführt
   als [`MR-055`](../../../../harness/conventions.md#mr-055).
5. **`AGENTS.md` gibt nichts ab.** Erst wenn gemessen ist, dass die Zustellung
   greift, ist zu entscheiden, was die Datei abgeben kann. Bis dahin ist die
   Doppelung benannt (§5), nicht aufgelöst.
6. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Hook.** Er ist der spätere Schritt, nicht der erste: deterministisch
  und pfad-gebunden, aber mit Heuristik, Event-Entscheid, Glob-Dialekt und
  eigener JSON-Extraktion. Wer ihn baut, sollte vorher wissen, ob Zustellung
  überhaupt wirkt — und genau das misst dieser Slice.
- **Keine `paths:`-Frontmatter.** Sie würden die Regeln im Auto-Modus stumm
  schalten (§1). Eines in den Symlink zu schreiben hieße außerdem, den
  gepinnten Baum zu ändern.
- **Kein Kürzen von `AGENTS.md`.** Solange nicht gemessen ist, dass die
  Zustellung greift, wäre jede gestrichene Zeile ein Verlust ohne Ersatz.
- **Keine weiteren 22 Module.** Wer alles einhängt, zahlt 4787 Zeilen für eine
  Lücke, die vier erklären.
- **Kein Ersatz für Gates.** Eingespeister Text ist Kontext, **keine erzwungene
  Konfiguration**. welle-86 bleibt davon unberührt.

## 4. Definition of Done

- [ ] `.claude/rules/` trägt **vier** Symlinks auf `modul-01`, `modul-05`,
      `modul-06`, `modul-08` des gepinnten Regelwerk-Baums; jeder löst auf, und
      die Summe der Zielzeilen steht als Zahl im Slice.
- [ ] **Die Zustellung ist belegt, nicht behauptet:** beim nächsten
      Sitzungsstart führt der Kontext die vier Module — belegt über `/context`
      oder den `InstructionsLoaded`-Hook, mit der Ausgabe im Slice. **Das ist
      der eine Punkt, den der bauende Lauf nicht selbst messen kann**; er gehört
      in die Closure, nicht in eine Zusage.
- [ ] **Der Bump-Träger ist gewächtert, nicht nur benannt:** `make
      baseline-verify` meldet einen toten Symlink unter `.claude/rules/` —
      **gemessen** an einer Probe, mit Erwartung und Ergebnis; geführt als
      [`MR-055`](../../../../harness/conventions.md#mr-055).
- [ ] Die **Nicht-Zusagen** stehen geschrieben: werkzeug-lokal, kein Gate,
      **nicht gate-geprüft** (Symlink), kein Ersatz für `AGENTS.md`.
- [ ] Der Nachfolge-Entscheid ist benannt — ob der Hook folgt, ob weitere
      Module dazukommen, und was `AGENTS.md` dann abgeben kann.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **805 Zeilen kommen in jeden Lauf, auch in den, der nichts mit Planung zu tun
  hat.** Zusammen mit `AGENTS.md` (527) sind das über 1300 Zeilen Kontext,
  bevor die erste Frage gestellt ist. Ob das billiger ist als eine Heuristik,
  ist die Wette dieses Slice. — **Ausgang:** *(bei Closure)*
- **Zwei Orte für dieselbe Aussage** — die Module und `AGENTS.md` §3/§5 sagen
  beide etwas über Planung und Closure. Solange `AGENTS.md` nichts abgibt (§3),
  ist das die Drift, die dieses Repo als [`BEO-010`](../observations.md)
  führt. — **Ausgang:** *(bei Closure)*
- **Die Symlinks sind nicht gate-geprüft** (gemessen, §2 Punkt 3) — ihr
  **Brechen** ist es seit diesem Slice: `make baseline-verify` meldet einen
  toten Alias. Ungeprüft bleibt ihr **Inhalt**, und das wiegt wenig: das Ziel
  liegt im SHA-gepinnten Baum. Ungeprüft bleibt auch ein Alias, der auf eine
  Datei **außerhalb** des Pins zeigt — er löst auf und passiert
  ([`MR-055`](../../../../harness/conventions.md#mr-055) §Grenze). —
- **Ein werkzeug-lokaler Träger ist ungebunden, sobald das Werkzeug wechselt**
  ([`MR-042`](../../../../harness/conventions.md#mr-042)). Anders als beim
  Wächter ist das hier **inhärent**: Kontext lässt sich nur dorthin einspeisen,
  wo einer ist. — **Ausgang:** *(bei Closure)*
- **Dass etwas geladen ist, heißt nicht, dass es wirkt.** Der Anlassfall bestand
  darin, dass niemand etwas vermisste; er lässt sich nicht wiederholen. Belegbar
  ist die Anwesenheit, nicht die Wirkung. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, beansprucht am 2026-08-30 —
slice-172 und slice-178 warten beide auf Vorbedingungen, dieser Slice auf
nichts.

**Rückführungen:** `in-progress` → `open`, falls der Beleg zeigt, dass die
Module **nicht** geladen werden — dann trägt auch dieser Kanal nicht, und der
Befund ist ein anderer: dann bleibt nur der Hook, und die Reihenfolge kehrt
sich um.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.claude/` (Werkzeug-Konfiguration) und `harness/`
  (die Bump-Prozedur). Beide fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-024`): [`BEO-024`](../observations.md) — ein Zustell-Kanal hängt an der
  **Arbeitsweise** statt am Inhalt: dieser Slice ist die erste Instanz und
  wählt seinen Träger genau danach; [`BEO-010`](../observations.md) — eine
  Aussage an zwei Orten: die Module und `AGENTS.md` wären genau das, solange
  keiner abgibt, und das steht als Risiko in §5;
  [`BEO-023`](../observations.md) — ein Wächter, der nie fangen konnte: hier
  ohne Wächter, aber mit derselben Frage — der Beleg muss die **Zustellung**
  zeigen, nicht ihre Konfiguration; [`BEO-008`](../observations.md) — die
  Spiegel einer Pin-Hebung sind mehrere Klassen, gehoben wird nur die
  grep-bare: die Symlinks sind ein **neuer** Spiegel, und genau deshalb steht
  die Bump-Prozedur als DoD-Punkt. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z,
  `image-scan.yml` 2026-08-29T10:07:43Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt ([`MR-054`](../../../../harness/conventions.md#mr-054)).
Slice-ID: slice-176. Betroffene IDs:
[`MR-021`](../../../../harness/conventions.md#mr-021),
[`MR-042`](../../../../harness/conventions.md#mr-042),
[`MR-043`](../../../../harness/conventions.md#mr-043). Module: —
(Werkzeug-Konfiguration; kein Produkt-Modul berührt). Gates: `make gates`,
`make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Vier Symlinks und eine Prozedur-Zeile; kein
Produkt-Code, kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
