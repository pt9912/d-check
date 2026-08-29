# Slice slice-176: Die Planungs-Regeln werden zugestellt, wenn sie gelten

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §1 (das Prinzip: Hard Rules und
**Pointer**, keine Duplikation) und §3 (der Abschnitt, der über die Hälfte der
Datei trägt);
[`MR-054`](../../../../harness/conventions.md#mr-054) (die Beleg-Form, die
dieser Slice übernimmt);
[`MR-043`](../../../../harness/conventions.md#mr-043) (der Werkzeug-Einstieg
importiert `AGENTS.md` — der Mechanismus, den dieser Slice ergänzt);
[`MR-042`](../../../../harness/conventions.md#mr-042) (die Präzedenz für einen
**werkzeug-lokalen** Träger und seine Nicht-Zusagen, und die Klasse von Lücken,
die eine Befehls-Zerlegung hat).

**Berührte Spec-Stellen:** — (Werkzeug-Konfiguration und Konventionsform; keine
Produkt-Anforderung berührt).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**`AGENTS.md` ist auf über 500 Zeilen gewachsen, und die Regeln, die eine
Sitzung gekostet haben, standen nicht darin.**

Gemessen: die Datei trägt 511 Zeilen gegen 236 der Baseline-Vorlage; allein
**§3 Harte Regeln** sind 280 davon — mehr als die Hälfte, und die meisten
dieser Regeln gelten für **einen** Pfad. Die Werkzeug-Dokumentation nennt genau
diese Lage und genau dieses Mittel:

> target under 200 lines per CLAUDE.md file. Longer files consume more context
> and reduce adherence. If your instructions are growing large, use path-scoped
> rules so instructions load only when Claude works with matching files.

**Der Import löst es nicht.** [`CLAUDE.md`](../../../../CLAUDE.md) zieht
`AGENTS.md` über `@AGENTS.md` in jeden Lauf — dieselbe Dokumentation sagt dazu,
Imports helfen der *Organisation*, nicht dem Kontext: sie laden beim Start
vollständig.

**Was gefehlt hat, ist belegt.** In der auslösenden Sitzung liefen drei Slices
(168, 169, 170) durch, ohne dass der Zyklus (`Spec → ADR → Plan → Code →
Review → Verifikation → Closure`) im Kontext war. Er steht in `modul-01`, das
niemand geöffnet hatte — und in `AGENTS.md` steht er nicht, korrekterweise:
§6 beschreibt den **Implementer**-Workflow, und der endet bei Schritt 8. Folge:
Review und Verifikation fielen aus, das Beobachtungs-Register blieb
unangetastet, drei Slices gingen mit offenem DoD-Haken nach `done/`.

## Der naheliegende Träger trägt nicht, und das ist jetzt der Kern

**Der erste Schnitt setzte auf `.claude/rules/` mit `paths`-Frontmatter — die
Annahme ist widerlegt, bevor eine Zeile geschrieben war.** Nach Auskunft des
Werkzeugs (Hinweis des Auftraggebers, 2026-08-29) hängt eine pfad-gebundene
Regel im **Auto-Modus** nur ein, wenn die passende Datei über die **dedizierten
Werkzeuge** (`Read`, `Edit`, `Write`) angefasst wird. Jeder Zugriff über die
Shell — `cat`, `sed`, `awk`, ein Skript — geht daran vorbei.

**Das trifft dieses Repo härter als die meisten.** Der Auto-Modus weist
ausdrücklich an, Datei-Änderungen über die Shell zu machen, wo sie es kann; in
der Sitzung, die diesen Slice ausgelöst hat, wurden die Planungs-Dateien
überwiegend mit `awk` und `sed` geschrieben, und die dedizierten Werkzeuge kamen
fast nur für ganz neue Dateien zum Einsatz. **Unter dieser Arbeitsweise hätte
die Regel im Anlassfall nicht geladen** — der Träger hätte genau dort
geschwiegen, wo er gebraucht wurde.

**Ein Kanal, der aussieht wie einer, der zustellt, ist schlimmer als keiner.**
Das ist die Klasse, die dieses Repo als [`BEO-023`](../observations.md) führt,
eine Ebene über dem Wächter: nicht eine Prüfung, die nicht greift, sondern eine
**Zustellung**, die nicht ankommt — und beides ist von außen grün.

## 2. Vorgehen

**Hooks statt Rules.** Ein Hook wird vom Harness ausgeführt, unabhängig davon,
welches Werkzeug der Lauf benutzt. Das ist die einzige Variante, die
deterministisch zustellt.

1. **Ein `PreToolUse`-Hook als Pilot, nicht fünf Regel-Dateien.** Er matcht
   `Bash|Edit|Write`, prüft die Nutzlast auf einen Bezug zu
   `docs/plan/planning/` und speist die Planungs-Regeln über
   `hookSpecificOutput.additionalContext` ein. Damit greift er auf **beiden**
   Wegen — dem dedizierten Werkzeug-Pfad und der Shell.
2. **Was immer gilt, geht in `SessionStart`.** Für `SessionStart`,
   `UserPromptSubmit` und `UserPromptExpansion` wird stdout unmittelbar
   Kontext. Der Zyklus gehört dorthin, wenn er für jeden Lauf gilt; nur das
   **Pfad-Gebundene** braucht den `PreToolUse`-Weg. Welche Hälfte was trägt,
   ist Teil dieses Slice.
3. **Die Zustellung stellt zu, sie dupliziert nicht.** Der eingespeiste Text
   trägt den Zyklus, die Lifecycle-Kanten, die Closure-Pflichten und die
   Rollen-Trennung in verdichteter Form — jeweils mit **zwei** Bindungen an die
   Quelle: einer `d-check:cite`-Direktive für den Wortlaut
   ([`MR-054`](../../../../harness/conventions.md#mr-054)) **und** einem
   Markdown-Link auf den Regelwerk-Abschnitt für die Vollform. Damit ist sie
   keine dritte Wahrheitsquelle, sondern die kontextsensitive Auslieferung der
   zweiten.
4. **Die Verweise sind gate-geprüft, und das ist der Grund für sie.** Der Text
   lebt in einer Datei unter `.claude/`, die im Scan-Bereich liegt
   (`scan.roots: ["."]`) — also halten `links`, `anchors` und `citations` ihn
   gegen seine Quelle. Eine Zustellung, die still von ihrer Quelle abdriftet,
   wäre schlimmer als keine.
5. **Proben statt Zusage.** Der Hook bekommt eine Proben-Datei nach dem Muster
   von `make guard-probe`: je Fall Erwartung und Ergebnis — Shell-Zugriff auf
   eine Planungs-Datei, dedizierter Werkzeug-Zugriff, ein Zugriff außerhalb des
   Pfades (muss schweigen). Ohne wiederholbare Proben wäre die Zusage eine
   Erinnerung.
6. **`AGENTS.md` gibt erst danach ab — und nur, was die Zustellung wirklich
   trägt.** Die Verdichtung ist ein **eigener Schritt nach der Messung**, kein
   Vorgriff: erst zeigen, dass die Zustellung greift, dann kürzen.
7. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.


### Eine Referenz-Implementierung liegt vor — und drei Deltas dazu

**Ein Schwester-Projekt fährt den Mechanismus bereits** (Hinweis des
Auftraggebers, 2026-08-29): ein Hook `rules-on-touch.sh`, der die
`paths`-Frontmatter der Regel-Dateien selbst liest, aus der Tool-Nutzlast
Kandidaten-Pfade sammelt — **auch aus dem Bash-Kommando** — und bei einem
Treffer den Regel-Rumpf über `additionalContext` einspeist. Er markiert je
Sitzung und Regel einmal, damit nichts doppelt kommt, und überspringt Regeln
**ohne** `paths:`, statt sie ungefragt einzuspielen. Damit ist die Frage „geht
das überhaupt?" beantwortet; offen bleiben drei Entscheide.

**Delta 1 — `jq` scheidet aus, und das ist hier entschieden.** Die Referenz
liest die Hook-Eingabe mit `jq`. Dieses Repo hat genau das schon einmal
verworfen: der Befehls-Wächter wurde in reinem `bash` + `awk` gebaut, und
`nightly-state.sh` trägt den Grund als Kommentar — das Werkzeug steht **nicht**
in der Host-Klasse aus [`AGENTS.md`](../../../../AGENTS.md) §3.1. Der Ersatz
liegt bereit: die JSON-Extraktion des Wächters lebt schon als eigenes
`awk`-Programm unter `tools/harness/` und liest dieselbe Nutzlast-Form.

**Delta 2 — `PostToolUse` oder `PreToolUse`?** Die Referenz nimmt
`PostToolUse`. Das hat zwei Vorzüge: der Hook feuert nur bei **Erfolg**, und er
mischt sich nicht in den Berechtigungs-Fluss. Der Preis: die Regel kommt für den
**nächsten** Schritt, nicht für den laufenden — wer eine Planungs-Datei
schreibt, hat sie beim Schreiben noch nicht. `PreToolUse` dreht beides um. Der
Entscheid gehört in diesen Slice, nicht in eine Annahme.

**Delta 3 — zwei Glob-Dialekte wären eine Drift-Quelle.** Die Referenz baut sich
`glob2re` selbst, mit eigener `**`-Semantik. Dieses Repo beantwortet
Glob-Fragen über `path.Match`. Zwei Dialekte für dieselbe Frage sind die
Gestalt, die hier als [`BEO-003`](../observations.md) geführt wird — eine
geteilte Lexik driftet, weil jeder Konsument sie selbst vorbereitet. Welcher
Dialekt gilt, ist zu entscheiden und zu **schreiben**, nicht zu erben.

**Ein vierter Punkt, der keine Entscheidung braucht, aber eine Zeile:** die
Einmal-je-Sitzung-Markierung ist richtig und hat einen Preis. Nach einer
Kontext-Verdichtung ist die Regel aus dem Kontext verschwunden, die Markierung
aber noch da — sie kommt dann nicht wieder. Das gehört zu den Nicht-Zusagen.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Kürzen von `AGENTS.md`.** Solange nicht gemessen ist, dass die
  Zustellung greift, wäre jede gestrichene Zeile ein Verlust ohne Ersatz. Die
  Datei bleibt in diesem Slice unverändert.
- **Keine `.claude/rules/`-Datei.** Der Träger ist widerlegt (§1); ihn
  „zusätzlich" anzulegen hieße, zwei Kanäle zu pflegen, von denen einer
  nachweislich schweigt.
- **Kein Ersatz für Gates.** Eingespeister Text ist Kontext, **keine erzwungene
  Konfiguration**. welle-86 bleibt davon unberührt.
- **Keine Aussage über andere Werkzeuge.** Der Hook wirkt in **einem**
  Werkzeug; `AGENTS.md` bleibt die werkzeug-neutrale Datei.
- **Keine Erweiterung der Wächter-Zerlegung.** Der Hook liest die Nutzlast, um
  einen Pfad-Bezug zu erkennen — er entscheidet nichts und blockiert nichts.

## 4. Definition of Done

- [ ] Der Hook ist in [`.claude/settings.json`](../../../../.claude/settings.json)
      verdrahtet, matcht `Bash|Edit|Write` und speist über
      `hookSpecificOutput.additionalContext` ein; die Zustell-Datei bleibt
      **unter 60 Zeilen** — eine Zustellung, die selbst zu lang ist, verfehlt
      ihren Zweck.
- [ ] **Beide Zugriffswege gemessen:** ein `sed`/`awk`-Zugriff auf eine Datei
      unter `docs/plan/planning/` löst den Hook aus, ein `Edit`/`Write`-Zugriff
      ebenso, ein Zugriff außerhalb des Pfades **nicht**. Je Fall Erwartung und
      Ergebnis, Ausgabe im Slice.
- [ ] Jeder normative Satz der Zustellung trägt **beide** Bindungen: eine
      `d-check:cite`-Direktive auf den Wortlaut **und** einen Link auf den
      Regelwerk-Abschnitt; `make doc-check` prüft beides grün.
- [ ] Die Proben laufen über ein `make`-Target und melden je Fall Erwartung und
      Ergebnis mit Fehlschlag-Zähler — Muster: `make guard-probe`.
- [ ] Die **Nicht-Zusagen** stehen geschrieben, nicht implizit: werkzeug-lokal,
      kein Gate, kein Ersatz für `AGENTS.md`, und der Unterschied zwischen
      *„der Hook hat gefeuert"* und *„der Text hat gewirkt"*; dazu, dass die
      Einmal-je-Sitzung-Markierung eine Kontext-Verdichtung **überlebt**, die
      eingespeiste Regel aber nicht.
- [ ] **Kein `jq`.** Die Hook-Eingabe wird mit `bash` + `awk` gelesen — dieselbe
      Klasse, die [`AGENTS.md`](../../../../AGENTS.md) §3.1 zusagt, und
      dieselbe, in der der Befehls-Wächter schon liest. Die vorhandene
      `awk`-Extraktion wird **wiederverwendet**, nicht nachgebaut.
- [ ] **Der Event-Entscheid steht geschrieben:** `PreToolUse` (Regel wirkt auf
      den laufenden Schritt, mischt sich aber in den Berechtigungs-Fluss) oder
      `PostToolUse` (feuert nur bei Erfolg, Regel wirkt ab dem nächsten
      Schritt) — mit Begründung, nicht als Übernahme.
- [ ] **Ein Glob-Dialekt, und er ist benannt.** Entweder `path.Match` wie im
      Produkt oder ein eigener mit `**` — geschrieben, nicht geerbt.
- [ ] Der Nachfolge-Entscheid ist benannt — ob weitere Pfade folgen und was
      `AGENTS.md` dann abgeben kann.
- [ ] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Messbar ist, dass der Hook feuert — nicht, dass der Text gewirkt hat.**
  `additionalContext` wird als Systemerinnerung eingespeist; ob sie das
  Verhalten ändert, ist dieselbe unbeantwortbare Frage wie beim Anlassfall, der
  darin bestand, dass niemand etwas vermisste. — **Ausgang:** *(bei Closure)*
- **Die Pfad-Erkennung in einer Shell-Nutzlast ist eine Heuristik.**
  [`MR-042`](../../../../harness/conventions.md#mr-042) buchstabiert für den
  Befehls-Wächter aus, was eine Zerlegung übersieht — Wrapper, Sub-Shells,
  wort-interne Splices. Dieselbe Klasse gilt hier mit umgekehrtem Vorzeichen:
  übersieht sie den Pfad, schweigt die Zustellung **still**. —
  **Ausgang:** *(bei Closure)*
- **Ein werkzeug-lokaler Träger ist ungebunden, sobald das Werkzeug wechselt**
  ([`MR-042`](../../../../harness/conventions.md#mr-042)). Anders als beim
  Wächter ist das hier **inhärent**: Kontext lässt sich nur dorthin einspeisen,
  wo einer ist. — **Ausgang:** *(bei Closure)*
- **Zwei Orte für dieselbe Aussage sind die Drift, die dieses Repo als
  [`BEO-010`](../observations.md) führt.** Die Zustellung und `AGENTS.md` sagen
  beide etwas über Planung. Solange `AGENTS.md` nichts abgibt (§3), ist das
  **Duplikation** — und der Slice endet in genau diesem Zustand. Die Auflösung
  liegt im Nachfolger, und bis dahin ist der Zustand benannt statt übersehen. —
  **Ausgang:** *(bei Closure)*
- **Die Bindung an das Regelwerk kostet beim Bump doppelt:** der Link trägt den
  `<tag>` ([`MR-021`](../../../../harness/conventions.md#mr-021)), die
  `cite`-Spanne die Zeilennummer
  ([`MR-051`](../../../../harness/conventions.md#mr-051)). Jede
  Zustell-Datei erhöht den Bump-Aufwand — und dieser Slice legt die erste an. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `open`, falls die Proben zeigen, dass der
Hook den **Shell**-Weg nicht zuverlässig erkennt — dann trägt auch dieser Kanal
die Pfad-Bindung nicht, und der Befund ist ein anderer: es bleibt, was für
**jeden** Lauf gilt (`SessionStart`), und die Pfad-Bindung entfällt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.claude/` (Werkzeug-Konfiguration) und
  `docs/plan/planning/` (der adressierte Pfad). Beide fallen unter den Default
  `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-024`): [`BEO-024`](../observations.md) — ein Zustell-Kanal hängt an der
  **Arbeitsweise** statt am Inhalt: dieser Slice ist die erste Instanz und
  zugleich ihre Antwort; [`BEO-023`](../observations.md) — ein Wächter, der nie
  fangen konnte: eine Zustellung, die nie ankommt, ist dieselbe Gestalt eine
  Ebene höher, weshalb die Proben als DoD-Punkt stehen;
  [`BEO-010`](../observations.md) — eine Liste mit Spiegeln außerhalb ihrer
  Datei: Zustellung und `AGENTS.md` §3 wären genau das, solange keine abgibt,
  und das steht als Risiko in §5; [`BEO-012`](../observations.md) — eine Quelle
  über ihren Geltungsbereich hinaus zitiert: eine verdichtete Zustellung ist
  genau diese Gefahr, weshalb jeder Satz **zwei** Bindungen an die Quelle
  trägt. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen:** entfällt in `open/` — der Block entsteht
  **spätestens bei der Beanspruchung** (`open→in-progress`), weil ein zum
  Planungszeitpunkt gelesener Stand bis dahin veraltet wäre
  ([`MR-053`](../../../../harness/conventions.md#mr-053)).

Slice-ID: slice-176. Betroffene IDs:
[`MR-054`](../../../../harness/conventions.md#mr-054),
[`MR-042`](../../../../harness/conventions.md#mr-042),
[`MR-043`](../../../../harness/conventions.md#mr-043). Module: `links`,
`anchors`, `citations`. Gates: `make gates`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter den
Default: Doc führt, Code folgt. Eine Werkzeug-Konfiguration plus ihre Bindung an
vorhandene Quellen; kein Produkt-Code, kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)
