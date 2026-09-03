# Slice slice-176: Der Zyklus ist im Kontext, bevor jemand ihn braucht

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht). Der **Anlass** liegt in
[welle-86](welle-86-closure-uebergang-durchsetzen.md); dort ist die
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

### Zwei Träger wurden geprüft, einer davon verworfen

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
`paths:`. **Belegt ist damit eine Sitzungsart, nicht „immer":** in einer
frischen interaktiven Hauptsitzung laden sie (beobachtet, §4); in einem
**Subagenten**-Kontext desselben Werkzeugs sind sie **nicht** präsent (§5). Der
Kanal hängt also nicht nur am Werkzeug, sondern am **Sitzungstyp** darin —
beobachtet in einem Schwester-Projekt,
das den ganzen Regelwerk-Baum so einhängt (Auskunft des Auftraggebers,
2026-08-30).

1. **Die Auswahl folgt dem Anlassfall, nicht dem Vollständigkeits-Wunsch:**
   `modul-01` (der Zyklus, dessen Fehlen Review und Verifikation kostete),
   `modul-05` (Planning-Form und Vorprüfungen), `modul-06` (Roadmap und
   Beobachtungs-Register — die zweite ausgefallene Pflicht), `modul-08`
   (Rollen-Trennung). **Gemessen: 805 Zeilen** von 4787. Die übrigen 22 Module
   bleiben draußen, bis ihr Fehlen belegt ist. **Gezählt sind Dateien:** der
   Baum trägt 26 Dateien — 17 `modul-*`, 8 `grundlagen-*` und eine `README.md`;
   eingehängt sind vier Module, offen also 13 Module bzw. 22 Dateien.
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


### Was diese Zustellung **nicht** zusagt

Sieben Nicht-Zusagen, ausgeschrieben statt implizit:

1. **Sie ist werkzeug-lokal.** Nur dieses Werkzeug liest `.claude/rules/`; ein
   fremder Agent liest [`AGENTS.md`](../../../../AGENTS.md) direkt und sieht von
   den vier Modulen nichts ([`MR-042`](../../../../harness/conventions.md#mr-042)).
   Anders als beim Befehls-Wächter ist das hier **inhärent** — Kontext lässt
   sich nur dorthin einspeisen, wo einer ist.
2. **Sie ist kein Gate.** Eingespeister Text ist Kontext, keine erzwungene
   Konfiguration. Die Durchsetzung des Closure-Übergangs bleibt
   [welle-86](welle-86-closure-uebergang-durchsetzen.md); daran ändert diese
   Zustellung nichts.
3. **Ihr Inhalt ist nicht gate-geprüft.** Der Scanner folgt Symlinks nicht in
   die Prüfmenge (gemessen) — `links`, `anchors` und `citations` sehen die
   Ziele über den Alias nicht. Gehalten wird stattdessen zweierlei: die
   **Integrität** der Ziele per SHA (`make baseline-verify`) und die
   **Auflösung** der Aliase seit
   [`MR-055`](../../../../harness/conventions.md#mr-055). Ungeprüft bleibt ein
   Alias, der auf eine Datei **außerhalb** des gepinnten Baums zeigt.
4. **Sie ersetzt `AGENTS.md` nicht.** Die Datei bleibt die werkzeug-neutrale
   Quelle und gibt in diesem Slice **nichts** ab; die Doppelung zwischen ihren
   §3/§5 und den Modulen besteht fort und ist als Risiko geführt.
5. **Anwesenheit ist nicht Wirkung.** Belegt ist, dass die vier Module im
   Kontext stehen. Ob sie das Verhalten ändern, ist nicht belegbar — der
   Anlassfall bestand darin, dass niemand etwas vermisste, und er lässt sich
   nicht wiederholen.
6. **Sie erreicht keinen Subagenten.** Gemessen in den Review- und
   Verifikations-Läufen zu diesem Slice: ihr Projekt-Kontext führt `CLAUDE.md`,
   `AGENTS.md` und die Nutzer-Memory — **keines der vier Module**. Das trifft
   ausgerechnet die Rollen, deren Ausfall den Anlassfall bildete: ein Review-
   oder Verifikations-Lauf bekommt `modul-01` und `modul-06` nicht mit. Der
   Kanal hängt am **Sitzungstyp**, nicht nur am Werkzeug — [`BEO-024`](../observations.md)
   eine Ebene tiefer.
7. **Sie überlebt nicht jedes Dateisystem.** Auf einem Checkout ohne
   Symlink-Unterstützung (`core.symlinks=false`) werden aus den Aliasen reguläre
   Dateien mit dem Pfad als Inhalt; die Zustellung ist dann eine Textzeile, und
   der Wächter sieht keinen Symlink mehr
   ([`MR-055`](../../../../harness/conventions.md#mr-055) §Grenzen).

### Der Nachfolge-Entscheid

**Der Hook ist nicht verworfen, sondern nachrangig.** Er löst dieselbe Aufgabe
pfad-gebunden und verdichtet, mit vier beweglichen Teilen (Pfad-Heuristik,
Event-Wahl, Glob-Dialekt, JSON-Extraktion ohne `jq`). Er lohnt sich, **wenn der
Preis drückt** — und der ist jetzt beziffert: **29,4k Token**. Das Kriterium
steht damit fest, statt Geschmacksfrage zu sein: solange das Kontextfenster den
Betrag trägt, ist die Zustellung ohne bewegliche Teile die bessere; wird er
knapp, ist der Hook der Ausweg.

**Weitere Module kommen erst, wenn ihr Fehlen belegt ist.** Die vier sind
gewählt, weil ihr Fehlen einen Ausfall verursacht hat — nicht, weil sie die
wichtigsten scheinen. Für jedes weitere gilt derselbe Maßstab; wer alle 26
einhängt, zahlt 4787 Zeilen für eine Lücke, die vier erklären.

**`AGENTS.md` kann jetzt zum ersten Mal etwas abgeben — und das ist ein eigener
Schnitt.** Solange die Module nicht im Kontext waren, wäre jede gestrichene
Zeile ein Verlust ohne Ersatz gewesen. Jetzt ist die Frage beantwortbar: welche
Aussagen in `AGENTS.md` §3 und §5 stehen bereits in `modul-05`/`modul-06`, und
welche sind **Adaptionen**, die der Kanon nicht trägt? Nur die erste Gruppe darf
weichen — die zweite ist der Grund, warum die Datei existiert. Der Schnitt
braucht eine **Messung** der Überlappung, keinen Rotstift.

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

- [x] `.claude/rules/` trägt **vier** Symlinks auf `modul-01`, `modul-05`,
      `modul-06`, `modul-08` des gepinnten Regelwerk-Baums; jeder löst auf
      (`readlink -e`, alle vier OK), und die Summe der Zielzeilen ist **805**
      von 4787 — gemessen, nicht geschätzt.
- [x] **Die Zustellung ist belegt, und der Beleg nennt die Dateien.** `/memory`
      in einer **frischen** Sitzung führt die vier Module einzeln auf — mit
      ihren **aufgelösten Zielpfaden** (`.harness/baseline/v5.12.0/regelwerk/…`),
      nicht mit den Alias-Namen: die Mechanik folgt dem Symlink also und meldet,
      was sie wirklich gelesen hat. `AGENTS.md` steht daneben als `@-imported`,
      die vier Module als eigene Einträge.
      **Herkunft des Belegs, benannt:** `/memory` und `/context` sind
      Werkzeug-Anzeigen einer Sitzung des Auftraggebers; sie liegen in keiner
      Datei dieses Repos und sind von hier aus nicht reproduzierbar. Was das
      Repo selbst hält, ist die **Auflösung** der Aliase (`make baseline-probe`).
- [x] **Der Preis ist gemessen, nicht geschätzt.** Dieselbe frische Sitzung
      zeigt in `/context` die Kategorie *Memory files* bei **58.3k** Token gegen
      **28.9k** in der bauenden — **+29.4k**. Die Größenordnung passt: die vier
      Ziele wiegen **59 868 Bytes**, also rund **2,0 Bytes je Token**, für
      deutsche Prosa mit Tabellen und Inline-Code der erwartete Wert.
- [x] **Der Bump-Träger ist gewächtert, nicht nur benannt:** `make
      baseline-verify` meldet einen toten Symlink unter `.claude/rules/` —
      **gemessen** an einer Probe, mit Erwartung und Ergebnis; geführt als
      [`MR-055`](../../../../harness/conventions.md#mr-055).
- [x] Die **Nicht-Zusagen** stehen geschrieben (§2, **sieben** Punkte):
      werkzeug-lokal · kein Gate · Inhalt nicht gate-geprüft (Symlink) · kein
      Ersatz für `AGENTS.md` · Anwesenheit ist nicht Wirkung · erreicht keinen
      Subagenten · überlebt nicht jedes Dateisystem.
- [x] Der Nachfolge-Entscheid ist benannt (§2) — ob der Hook folgt, ob weitere
      Module dazukommen, und was `AGENTS.md` dann abgeben kann.
- [x] `make gates` grün (Exit explizit) — 602 Dateien, 0 Befunde, Exit 0;
      **unabhängiger Review** (1 HIGH · 8 MEDIUM · 4 LOW · 2 INFO, blockierend)
      und **Verifikation** gegen DoD (A-1 bis A-14) in eigenen Kontexten
      gelaufen, beide Berichte unter [`docs/reviews/`](../../../reviews/),
      alle Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **29,4k Token kommen in jeden Lauf, auch in den, der nichts mit Planung zu tun**
  (gemessen, §4). Zusammen mit `AGENTS.md` sind das über 1300 Zeilen Kontext,
  bevor die erste Frage gestellt ist — bei diesem Modell 2,9 % des Fensters, bei
  einem 200k-Modell rund 15 %. Ob das billiger ist als eine Heuristik,
  ist die Wette dieses Slice. — **Ausgang:** weiter offen. Der **Preis** ist
  gemessen, der **Nutzen** nicht — und er ist es aus demselben Grund wie im
  Anlassfall: niemand vermisst, was er nicht kennt. Was sich seither belegen
  lässt, ist die Anwesenheit.
- **Zwei Orte für dieselbe Aussage** — die Module und `AGENTS.md` §3/§5 sagen
  beide etwas über Planung und Closure. Solange `AGENTS.md` nichts abgibt (§3),
  ist das die Drift, die dieses Repo als [`BEO-010`](../observations.md)
  führt. — **Ausgang:** weiter offen, und um eine Stelle schärfer. `AGENTS.md`
  §1 nennt die vier Module jetzt ausdrücklich, die Doppelung ist damit
  **deklariert** statt still — gemessen ist sie nicht. Der Schnitt, der
  `AGENTS.md` etwas abnehmen könnte, braucht diese Messung und steht in §2 als
  eigener.
- **Die Symlinks sind nicht gate-geprüft** (gemessen, §2 Punkt 3) — ihr
  **Brechen** ist es seit diesem Slice: `make baseline-verify` meldet einen
  toten Alias. Ungeprüft bleibt ihr **Inhalt**, und das wiegt wenig: das Ziel
  liegt im SHA-gepinnten Baum. Ungeprüft bleibt auch ein Alias, der auf eine
  Datei **außerhalb** des Pins zeigt — er löst auf und passiert
  ([`MR-055`](../../../../harness/conventions.md#mr-055) §Grenze). —
  **Ausgang:** weiter offen — und der Wächter hielt bei seiner Ablieferung
  weniger, als dieser Absatz behauptete: ein Alias in einem Unterverzeichnis
  und einer mit Punkt-Namen passierten still. Beides ist behoben und mit neun
  Proben belegt; die Grenzen sind von einer auf sechs ausgeschrieben.
- **Ein werkzeug-lokaler Träger ist ungebunden, sobald das Werkzeug wechselt**
  ([`MR-042`](../../../../harness/conventions.md#mr-042)). Anders als beim
  Wächter ist das hier **inhärent**: Kontext lässt sich nur dorthin einspeisen,
  wo einer ist. — **Ausgang:** weiter offen, mit einer gemessenen zweiten Kante:
  die Zustellung erreicht auch **innerhalb** dieses Werkzeugs keinen
  Subagenten — also genau die Review- und Verifikations-Rollen, deren Fehlen
  den Anlassfall auslöste.
- **Dass etwas geladen ist, heißt nicht, dass es wirkt.** Der Anlassfall bestand
  darin, dass niemand etwas vermisste; er lässt sich nicht wiederholen. Belegbar
  ist die Anwesenheit, nicht die Wirkung. — **Ausgang:** weiter offen, und
  dieser Punkt bleibt es dauerhaft. Er ist kein Rest, sondern die Form der
  Zusage.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei, beansprucht am 2026-08-30 —
slice-172 und slice-178 warten beide auf Vorbedingungen, dieser Slice auf
nichts.

**Rückführungen:** `in-progress` → `open`, falls der Beleg zeigt, dass die
Module **nicht** geladen werden — dann trägt auch dieser Kanal nicht, und der
Befund ist ein anderer: dann bleibt nur der Hook, und die Reihenfolge kehrt
sich um.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `.claude/` (Werkzeug-Konfiguration), `harness/`
  (die Bump-Prozedur) und `tools/harness/` (das Gate-Skript, das die dritte
  Frage bekommt). Alle drei fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration) — `tools/harness/` führt zwar eine eigene Deklaration, sie
  lautet aber ebenfalls Greenfield.
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
  die Bump-Prozedur als DoD-Punkt; [`BEO-022`](../observations.md) — eine Regel
  tritt in Kraft, bevor ihre Zustellung existiert: **der Eintrag, der diesen
  Slice namentlich führt**, und seine Prozedur verlangt, in derselben Änderung
  den **Lesepfad** zu benennen — deshalb nennt `AGENTS.md` §1 die vorgeladenen
  Module (§2). Die Regel, die diesen Schritt vorschreibt:

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

**GF (Greenfield, Repo-Default)** — alle drei berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Vier Symlinks, eine Prozedur-Zeile und eine
dritte Frage im Gate-Skript; kein Produkt-Code, kein Fremdsystem, keine
Reconciliation.

## 9. Closure-Notiz (nach `done/`)

**Geliefert, und der Kanal war ein anderer als der geplante.** Der Slice war
zweimal geschnitten: erst auf `.claude/rules/` mit `paths`-Frontmatter, dann —
nach der Werkzeug-Auskunft, dass eine solche Regel im Auto-Modus nur beim
Zugriff über die dedizierten Werkzeuge lädt — auf einen Hook. Geliefert hat
schließlich der erste Kanal in der dritten Form: **vier Symlinks ohne
`paths`**. Ein Eintrag ohne diese Bedingung hat den Werkzeugweg-Vorbehalt
nicht und lädt beim Sitzungsstart. Der Hook ist damit nicht verworfen, sondern
nachrangig — und er hat seit diesem Slice ein Kriterium statt einer Meinung:
er lohnt, wenn der Preis drückt, und der ist mit 29,4k Token beziffert.

**Der Beleg ist eine fremde Sitzung, und das steht so da.** `/memory` einer
frischen Hauptsitzung nennt die vier Module mit ihren **aufgelösten**
Zielpfaden, `/context` die Kategorie *Memory files* bei 58.3k gegen 28.9k in
der bauenden. Beides sind Werkzeug-Anzeigen, die in keiner Datei dieses Repos
liegen und von hier aus nicht reproduzierbar sind; reproduzierbar ist die
**Auflösung** der Aliase (`make baseline-probe`). Die Grenze war zuerst nur in
der Commit-Botschaft benannt — die Verifikation hat sie als A-9 zurück in den
Plan geholt.

**Der Wächter, den dieser Slice mitbrachte, hielt weniger als sein Wort.**
[`MR-055`](../../../../harness/conventions.md#mr-055) sagte an drei Stellen
„jeder Symlink"; der erste Glob war einstufig und punktblind. Ein Alias in einem
Unterverzeichnis von `.claude/rules/` und einer mit Punkt-Namen passierten still
— gemessen vom Review, nicht vermutet. Das ist genau die Bauform, gegen die
dieses Repo [`BEO-023`](../observations.md) führt: ein Wächter, der aussieht wie
einer, der fängt. Behoben mit `find` statt Glob, mit einer eigenen Meldung für
den Schleifenfall, mit `readlink` in der Werkzeug-Vorprüfung — und mit neun
Proben unter `make baseline-probe`, weil eine Zusage ohne wiederholbare Probe
eine Erinnerung ist. Die Grenzen stehen jetzt zu sechst in
[`MR-055`](../../../../harness/conventions.md#mr-055), statt zu einer.

**Vier Deklarationen sagten „zwei Hälften", während der Code drei Fragen
stellte.** Die schwerste saß in der Sensors-Tabelle von
[`harness/README.md`](../../../../harness/README.md) — der Tabelle, auf die
`AGENTS.md` §4 selbst für „Details und Bindungen" verweist. Das ist keine
Lücke, sondern eine Falschaussage an der Stelle, die als genauere gilt. Alle
vier Flächen sind nachgezogen.

**Zwei Planungs-Stellen behaupteten das Gegenteil des Gelieferten.**
[welle-86](welle-86-closure-uebergang-durchsetzen.md) §4 schrieb „der
Rules-Kanal fällt weg" und nannte diesen Slice als Beleg; die `Stand`-Zelle von
[`BEO-024`](../observations.md) führte „Hooks statt Regeln" als Antwort. Beide
waren zum Zeitpunkt ihres Schreibens richtig und sind es durch die Lieferung
nicht mehr geblieben — ein Zustandsfeld, das seinen Zustand verpasst hat.
Beide sind fortgeschrieben.

**Die Klasse, die den Anlassfall trug, ist verschoben, nicht getilgt.** Die
Zustellung erreicht **keinen Subagenten**-Kontext — also genau die Review- und
Verifikations-Rollen, deren Fehlen `BEO-022` zweimal gefunden hat, und die
auch diesen Slice geprüft haben. Sie steht als sechste von sieben
Nicht-Zusagen im Plan und in der `Stand`-Zelle des Eintrags.

**Fortgeschrieben:** [`BEO-022`](../observations.md) (die Beobachtung zur
Zustellung ist beantwortet; die Klasse verschiebt sich in den
Subagenten-Kontext), [`BEO-024`](../observations.md) (die Antwort war nicht der
Hook), [`MR-055`](../../../../harness/conventions.md#mr-055) (Anker auf
`modul-02` statt auf die Verzeichniskonvention, `## Begründung`, sechs
Grenzen), `AGENTS.md` §1 und §4, die Sensors-Tabelle und die Bump-Prozedur.

**Was offen bleibt und wohin es gehört:** die Messung der Überlappung zwischen
`AGENTS.md` §3/§5 und `modul-05`/`modul-06` — sie ist die Vorbedingung dafür,
dass die Datei zum ersten Mal etwas abgeben kann, und ein eigener Schnitt.
Weitere Module kommen erst, wenn ihr Fehlen belegt ist; derselbe Maßstab, mit
dem diese vier gewählt wurden.
