# Review-Report: slice-210 — 2026-09-07

**Review-Art:** Code — geprüft wird der Diff gegen Plan, Hard Rules und Kanon.

**Gegenstand:** slice-210, Commit-Range `1a32c617^..HEAD` (vier Commits:
`1a32c617` Vorprüfungen · `971aaa3d` Beanspruchung · `9306e18b` DoD 1 ·
`8f08d48b` DoD 2/3).

**Skill:** `.harness/skills/reviewer.md` @ 1.15.0 · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** claude-opus-5[1m] · **Datum:** 2026-09-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-NNN` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als **Tag + Pfad
> in Inline-Code** statt als Link (`v<X.Y.Z>` · `regelwerk/<datei>.md`
> §<Abschnitt>). Der vendored Baum trägt genau einen Tag; der Sprung löscht den
> alten, und ein Link darauf färbt beim nächsten Bump ein Artefakt rot, das
> niemand mehr anfassen darf.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-210 (Baseline-Form, acht Abschnitte)
- `AGENTS.md` §3.1–§3.9, §5, §6
- Registereintrag `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`
  (`observation.md`, `state.md`, drei Evidence-Dateien)
- Geschwister-Einträge `BEO-ALL/commit-message-overclaims-work`,
  `BEO-ALL/citation-stretched-beyond-scope`,
  `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`,
  `BEO-ALL/rule-drawn-from-occasion-not-inventory`,
  `BEO-ALL/begruendung-traegt-entscheidung-nicht`
- `MR-013`, `MR-053`, `MR-054`, `MR-070`
- `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice,
  §Zwei Schritte vor der Modus-Begründung
- `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register
- `v6.5.0` · `regelwerk/grundlagen-traceability.md` §Herkunfts-Anker
- `v6.5.0` · `templates/docs/plan/planning/slice.template.md`
- Vorherige Findings am gleichen Gegenstand: Review-Runde 1 zu slice-209
  (Kollision `MR-070`/`MR-013`), Review-Runde 2 zu slice-205 (`R2-M-1`,
  die Fehlzählung, die dieser Slice verkörpert), Review-Runde 1 zu slice-208

**Eigene Läufe** (echte Ausgabe):

```
$ make doc-check
d-check: 706 Datei(en) geprüft, 0 Befund(e)

$ make gates
[gates] baseline-verify + workflow-pins + doc-check + lint + test +
        arch-check + coverage-gate + semgrep + gate-consistency +
        planning-check green
        (doc-check/targets/planning/workflows je 706 Datei(en), 0 Befund(e);
         semgrep: 55 Regeln auf 63 Dateien, 0 findings; test: gesamt 0 Befund(e))

$ make review-coverage
d-check: 706 Datei(en) geprüft, 0 Befund(e)
```

---

## Findings

### F-1 — Die Abgrenzungs-Zeile der neuen Regel zeigt auf den falschen Nachbarn

- `kategorie`: MEDIUM
- `quelle`: Maintainability · `AGENTS.md` §5 (DoD (1) dieses Slice)
- `pfad`: `AGENTS.md:560-562`
- `befund`: Der neue Absatz sagt *„**Andere Frage als der Absatz darüber:** Der
  prüft, ob der **Schluss** weiter reicht als die gemessene Menge"*. Der Absatz
  darüber (`AGENTS.md:538-553`) regelt die drei Vorprüfungen und die
  `d-check:cite`-Direktiven und enthält weder *Schluss* noch *gemessene Menge*;
  der Absatz, der das prüft, steht **darunter** (`AGENTS.md:568-575`,
  `commit-message-overclaims-work`). Damit trifft genau der Satz ins Leere, den
  DoD (1) erzeugt hat und den `8f08d48b` als Schutz gegen eine spätere
  Dopplungs-Lesart benennt (*„Ohne diesen Satz wäre sie beim nächsten Aufräumen
  als Dopplung zu lesen"*); ein Leser, der der Ortsangabe folgt, findet keine
  Abgrenzung und muss sie neu herleiten oder die Regel für redundant halten.
- `verifizierbar`: nein — kein Gate deckt Richtungsangaben in Prosa; per Lesen
  der beiden Nachbarabsätze in derselben Datei reproduzierbar.
- `klasse`: Abgrenzungs-Zeiger benennt den falschen Nachbarn

### F-2 — Der Geschwister-Eintrag wird als *verkörpert* geführt, seine `state.md` sagt *gemischt*

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 (*„Vor jedem Verweis das **Feld** lesen, nicht den
  Titel"*) · `BEO-ALL/citation-stretched-beyond-scope`
- `pfad`: slice-210 §1 (Zeile 42) und §8 (Zeile 237)
- `befund`: Beide Stellen führen
  `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` als *„(10×, verkörpert)"*.
  Dessen `state.md` beginnt mit *„**Stand:** gemischt, differenziert nach
  Instanz (ein einzelnes Wort trüge zu wenig)"* und weist zwei seiner vier
  Instanzen ausdrücklich als *kein formgültiger Ausgang* und eine als *geplant
  (welle-86)* aus. Der Plan fasst also genau in das eine Wort zusammen, gegen
  das die Quelle sich verwahrt, und friert das mit der Closure in `done/` ein;
  ein späterer Sichtungs-Schritt liest den Eintrag als abgeschlossen und
  übersieht die zwei ausgangslosen Instanzen. Die Entscheidung des
  Ausschlusses selbst hängt nicht daran — sie ruht auf der verschiedenen
  **Frage** (Menge vs. Kategorie), nicht auf dem Stand.
- `verifizierbar`: nein — kein Gate liest Stand-Felder gegeneinander; per
  Vergleich der beiden Zeilen mit der zitierten `state.md` reproduzierbar.
- `klasse`: Stand-Feld auf ein Wort verkürzt, das die Quelle ablehnt

### F-3 — Die `state.md` speichert einen abgeleiteten Zähler neben der Belegliste

- `kategorie`: MEDIUM
- `quelle`: `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Das
  Beobachtungs-Register (*„Der Zähler wird abgeleitet, nicht geführt … Es gibt
  kein Feld, in das man ihn schreibt, und deshalb keines, das falsch stehen
  kann"*) · `v6.5.0` · `templates/docs/plan/planning/slice.template.md` §2
  (*„**kein Zaehler wird gesetzt**, er folgt aus den Dateien"*)
- `pfad`: `docs/plan/planning/observations/BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand/state.md:1`
- `befund`: Der dritte Satz lautet *„Drei Vorgänge, **vier** Instanzen —
  dreimal fand sie ein unabhängiger Review, einmal der Zählende selbst."*
  *„Drei Vorgänge"* ist der Zähler des Eintrags (drei Evidence-Dateien),
  ausgeschrieben in ein Feld; er muss bei jedem weiteren Beleg von Hand
  nachgezogen werden, und der Eintrag erklärt im selben Satzblock, dass er wach
  bleibt. Dass die Kopie driftet, ist hier nicht Prognose, sondern gemessen: bei
  **unveränderter** Evidence-Basis (`slice-205.md` seit `a0a3e61e`,
  `slice-207.md` seit `33c24e11`, `slice-208.md` seit `0b6f276e`) stand in
  derselben Datei nacheinander *„zwei Vorgänge, vier Instanzen"* (`33c24e11`),
  *„Drei Vorgänge, fünf Instanzen"* (`b2b4acf5`) und jetzt *„Drei Vorgänge,
  vier Instanzen"* (`9306e18b`). Die zweite Größe (*Instanz*) ist zusätzlich
  eine, die der Kanon ausdrücklich **nicht** zählt (*„Zwei Funde im selben
  Vorgang sind … eine Gelegenheit … nicht die Zahl der Funde"*). Nachgezählt
  ist der aktuelle Stand korrekt (1 + 2 + 1 = 4; drei Review-Funde über
  `R2-M-1` zu slice-205, den zweiten Fund in `evidence/slice-207.md` und `F`
  in Review-Runde 1 zu slice-208, ein Selbst-Fund) — der Befund gilt der
  gespeicherten Form, nicht der Zahl.
- `verifizierbar`: nein — der Kanon verzichtet bewusst auf ein Zähler-Feld und
  damit auf einen Sensor darauf; per `ls evidence/` gegen den Satz reproduzierbar.
- `klasse`: gespeicherter Zähler neben der Belegliste

### F-4 — Die neue Regel nennt keinen Ort, an dem die Form stehen muss

- `kategorie`: MEDIUM
- `quelle`: Maintainability · `AGENTS.md` §5
- `pfad`: `AGENTS.md:554-567`
- `befund`: Die Regel verlangt, die Form des Gegenstands werde *„ausgeschrieben,
  bevor gezählt wird"*, benennt aber weder das Artefakt noch die Stelle — anders
  als ihre beiden Nachbarn, die *„Commit-Botschaft oder Closure-Notiz"* bzw.
  *„vor jedem Verweis"* adressieren. Damit ist ihre Erfüllung nicht feststellbar:
  Der Skill-Anker 17 macht *„steht die **Form** des Gegenstands nirgends?"* zur
  Finding-Bedingung, ohne dass *irgendwo* einen Geltungsbereich hätte — eine im
  Kopf des Zählenden gebildete Form erfüllt den Wortlaut ebenso wie eine
  geschriebene. Der Lauf demonstriert das an sich selbst: Seine eigene
  Kennzahl *„vier Instanzen"* ruht auf einer Form von *Instanz*, die weder im
  Plan noch im Register noch in einer Commit-Botschaft ausgeschrieben ist, und
  nach dem Wortlaut der Regel lässt sich nicht sagen, ob das ein Verstoß ist.
  Die deklarierte Grenze (*„gilt der Methode, nicht der Sorgfalt … macht ihn
  auffindbar"*) deckt das nicht: *auffindbar* setzt voraus, dass die Form dort
  steht, wo der Findende sucht.
- `verifizierbar`: nein — die Regel ist als Urteil deklariert; per Anwendung des
  Anker-17-Wortlauts auf die eigene Kennzahl dieses Slice reproduzierbar.
- `klasse`: Regel ohne benannten Ort ihrer Erfüllung

### F-5 — Zwei verschiedene Zähler-Stände für denselben Eintrag im selben Plan

- `kategorie`: LOW
- `quelle`: `BEO-ALL/eigene-menge-gemessen-fremde-behauptet` (verkörperte
  Klasse: *wer ändert die Menge, die ich zähle?*)
- `pfad`: slice-210 §6 (Zeile 175) gegen §8 (Zeile 242)
- `befund`: §6 führt `rule-drawn-from-occasion-not-inventory` mit *„7×"*, §8
  mit *„8×, zuletzt slice-209"*. Der Eintrag hat acht Evidence-Dateien; die
  achte (`slice-209.md`) kam mit `3cdc3c69` dazu, also nach dem Schreiben von
  §6 in `open/` und vor dem Schreiben von §8 in `1a32c617`. §6 wurde nicht
  nachgezogen. Ein Leser des in `done/` eingefrorenen Plans entnimmt je nach
  Abschnitt einen anderen Stand.
- `verifizierbar`: nein — kein Gate hält Zahlen in Plan-Prosa gegen das
  Register; per `ls …/rule-drawn-from-occasion-not-inventory/evidence/`
  reproduzierbar.
- `klasse`: Zähler-Stand im eigenen Dokument nicht nachgezogen

### F-6 — Tabellenzeile 17 verengt die Klasse gegenüber ihrem Anker

- `kategorie`: LOW
- `quelle`: `.harness/skills/reviewer.md` §Die sechzehn Prüffragen (*„die
  Tabelle trägt keine Ausnahme und keine Begründung, sondern ausschließlich die
  Frage"*)
- `pfad`: `.harness/skills/reviewer.md:47` gegen `:103-118`
- `befund`: Die Zeile fragt *„Zählt eine Messung ein Muster, das dem Gegenstand
  nur ähnelt — **und** steht die Form des Gegenstands nirgends?"* und macht das
  Fehlen der Form damit zur **Bedingung** des Findings. Der Anker darunter
  benennt als Klasse allein *„Messung zählt einen Proxy statt des
  Gegenstands"*; die Form ist dort **Arbeitsanweisung** zur Erkennung, nicht
  Teil des Tatbestands. Ein Reviewer, der nur aus der Tabelle arbeitet — die
  Skill-Datei rechnet mit ihm (*„Wer nur diese Tabelle liest"*) —, meldet eine
  Proxy-Messung nicht, deren Form irgendwo geschrieben ist, obwohl sie den
  falschen Gegenstand zählt.
- `verifizierbar`: nein — Konsistenz zwischen den zwei Skill-Ebenen ist
  ungewächtert; per Vergleich der beiden Stellen reproduzierbar.
- `klasse`: Tabellen-Frage verengt den Anker, den sie indiziert

### F-7 — Herkunfts-Anker steht im Abschnitts-Rumpf, nicht in seiner Überschrift

- `kategorie`: INFO
- `quelle`: `v6.5.0` · `regelwerk/grundlagen-traceability.md` §Herkunfts-Anker
  (*„beim Abschnitt in dessen Überschrift"*)
- `pfad`: `AGENTS.md:567`
- `befund`: Die `state.md` nennt als Zielort `AGENTS.md` §5; der Anker
  `seit slice-210` steht am Ende des Regel-Bullets, die Überschrift
  *„## 5. Dokumentations-Regeln"* trägt ihn nicht. Für einen Abschnitt mit
  mehreren Steering-Loop-Regeln ist die Kanon-Form nicht durchführbar, und der
  Bestand weicht durchgehend gleich ab (`seit welle-82` in `AGENTS.md:575`,
  `seit slice-147` in `:589`, `seit welle-73` in `:360`). Dieser Slice führt
  die Abweichung also nicht ein, sondern setzt sie fort; deklariert ist sie in
  keinem `MR-<NNN>`. Gemeldet nach `AGENTS.md` §1 (*„Melde den Widerspruch"*),
  nicht als Mangel dieses Laufs.
- `verifizierbar`: nein — die Anker-Paarung ist in diesem Repo nicht
  mechanisiert.
- `klasse`: Kanon-Abweichung im Bestand, undeklariert

---

## Negativbefunde

- **geprüft, ohne Befund: Die Lücken-Frage von DoD (1) trägt.** `AGENTS.md` §5
  Absatz `commit-message-overclaims-work` (`:568-575`) sagt *„ihr Schluss reicht
  nicht weiter als die gemessene Menge"* — eine Aussage über den **Schluss**;
  Skill-Anker 8 (`:98-100`) nennt den Grund für seine Kategorie wörtlich
  *„weil die Messung stimmt und nur ihre Reichweite überdehnt ist"* und schreibt
  die korrekte Messung damit ausdrücklich als **Voraussetzung** fest. Die
  Selbstauskunft des Trägers deckt die neue Klasse also nicht; der neue Absatz
  ist eine zweite Frage, keine zweite Quelle. Die Vorsichtsmaßnahme gegen
  `citation-stretched-beyond-scope`, die der Plan sich auferlegt, ist
  eingehalten — §5 wird nicht enger gelesen, als der Träger sich selbst liest.
- **geprüft, ohne Befund: Der Zielort-Entscheid trägt.** Das erste Argument
  (zwei Handlungen, zwei Zeitpunkte — Messender *vor*, Prüfender *nach* dem
  Zählen) ist am geschriebenen Text nachvollziehbar: `AGENTS.md:556`
  (*„ausgeschrieben, bevor gezählt wird"*) gegen `reviewer.md:110-111`
  (*„verlange die Form vor der Zahl … sieh dir drei Treffer an"*). Die
  Gegenprobe aus `begruendung-traegt-entscheidung-nicht` hält: fiele das zweite
  Argument (die Zahl) weg, bliebe die Skill-Hälfte durch die zweite Handlung
  begründet.
- **geprüft, ohne Befund: Die Geschwister stehen tatsächlich an beiden Orten.**
  `commit-message-overclaims-work` in `AGENTS.md:568-575` und `reviewer.md:92-102`;
  `citation-stretched-beyond-scope` in `AGENTS.md:576-589` und
  `reviewer.md:119-136`. Die Prämisse des Plans („*beide*" ist die verdächtige
  Antwort) ist korrekt belegt.
- **geprüft, ohne Befund: Nachgezählt — die korrigierte Zahl stimmt und ist
  überall angekommen.** Vier Instanzen über drei Vorgänge (slice-205: ein Fund;
  slice-207: *„Zweimal im selben Slice"*; slice-208: ein Fund); drei davon fand
  ein unabhängiger Review — Review-Runde 2 zu slice-205 (`R2-M-1`, *„Die
  Zählmethode zählt keine Anforderungen"*), der zweite Fund in
  `evidence/slice-207.md` (*„der zweite erst im unabhängigen Review"*) und
  Review-Runde 1 zu slice-208 (*„von den acht Ausschluss-Abschnitten nennen
  zwei einen Folge-Slice"*) —, einen der Zählende selbst
  (*„der erste fiel beim Klassifizieren auf"*). Kein Rückstand von *„viermal
  von fünf"*: die Korrektur steht in slice-210 §3 (Zeilen 90, 124, 135), §6
  (Zeile 180), `reviewer.md:116-118`, `state.md:1` und in den Botschaften von
  `9306e18b` und `8f08d48b`. Ein repo-weiter Griff nach *„fünf Instanzen"* /
  *„viermal von fünf"* außerhalb von `docs/reviews/` liefert null Treffer.
- **geprüft, ohne Befund: Abgrenzung zu `eigene-menge-gemessen-fremde-behauptet`
  verschwimmt für einen §5-Leser nicht.** Der Geschwister-Eintrag ist **nicht**
  in `AGENTS.md` §5 verkörpert — seine `state.md` weist Instanz 1 als
  verkörpert *seit slice-161* über eine Sensor-Änderung aus, nicht über eine
  Regel-Zeile. In §5 steht daher nur eine der beiden Messregeln; eine
  Verwechslung ist dort nicht möglich. Die ausdrückliche Abgrenzung führt
  `observation.md` des neuen Eintrags (unverändert, korrekt).
- **geprüft, ohne Befund: `AGENTS.md` §3.7 / Zustandsfeld-Form der `state.md`.**
  Der Eintrag nennt Zustand (*verkörpert*) und Beleg als auflösbare Anker (zwei
  Links plus `seit slice-210`) und hat die Chronik-Sätze des Vorgängers
  (*„Die Schwelle ist mit dem dritten Beleg (slice-208) erreicht; der Ausgang
  wird dort zu verkörpert"*) fallen lassen. Beide Zielorte lösen auf
  (`make doc-check` grün) und tragen den Anker: `AGENTS.md:567` und
  `reviewer.md:118`. Der Skill-Anker ist damit strenger geführt als die
  Geschwister-Anker 8 und 9, die keinen `seit`-Anker tragen, obwohl ihre
  `state.md` sie als Zielort nennt. Zur Platzierung siehe F-7; zur dritten
  Satzhälfte siehe F-3.
- **geprüft, ohne Befund: §8, die drei Vorprüfungen.** Beide
  `d-check:cite`-Direktiven ankern auf die **vorschreibenden** Zeilen
  (`v6.5.0` · `regelwerk/modul-05-planning-harness.md`:268-269 und :274) und
  zitieren wortgleich — identisch zu der von slice-209 etablierten Form, und
  `citations` läuft im inneren Loop grün. Der Nachtlauf-Block trägt
  planmäßig keine Direktive (`MR-054`). Die Sub-Area-Wahl (`*`) ist mit dem
  Ausschluss eines Sensors begründet; das Register ist mit 37 Verzeichnissen
  vollständig durchgegangen (nachgezählt: 36 `BEO-ALL/` + 1 `BEO-HARN/`), und
  die sechs als einschlägig geführten Einträge sind es. Die acht Abschnitte
  entsprechen `v6.5.0` · `templates/docs/plan/planning/slice.template.md`.
- **geprüft, ohne Befund: §3.3 / `MR-013`, der Beanspruchungs-Commit.**
  `971aaa3d` zeigt `R100` auf der Slice-Datei (reiner Move, Rename-Detection
  sicher) und trägt genau die gekoppelten Verweise mit: den Roadmap-Flip
  (byte-exakte Umkehr von `78a6c130`, derselbe Absatz und dieselbe Leerzeile)
  sowie zwei Pfad-Nachzüge (`done/slice-209`, `state.md`). Vier Dateien, keine
  inhaltliche Änderung an der bewegten Datei.
- **geprüft, ohne Befund: `MR-070` kollidiert nicht.** Sein Geltungsbereich ist
  *„jede mechanische Ersetzung über mehr als eine Datei"* und nimmt den
  Lifecycle-Pfad-Nachzug ausdrücklich aus; die neue Regel gilt einer **Zählung**.
  Gemeinsam ist beiden die Aufforderung, das Ergebnis statt der Aggregat-Zahl
  anzusehen — sie treffen aber verschiedene Operationen und widersprechen
  einander an keiner Stelle. Auch `AGENTS.md` §3.1 (Mess-Rangfolge *Produkt vor
  `grep`/`awk`*) kollidiert nicht: dort geht es um das Werkzeug, hier um die
  Methode.
- **geprüft, ohne Befund: §1-Abgrenzung nicht ausgeweitet.** Der Diff berührt
  `AGENTS.md`, `.harness/skills/reviewer.md`, den Plan, die Roadmap, eine
  `state.md` und eine Pfad-Zeile in `done/slice-209`. Kein Sensor entstanden,
  der Geschwister-Eintrag unangetastet, die drei Evidence-Dateien der Vorgänge
  unverändert (kein Retrofit).
- **geprüft, ohne Befund: Commit-Botschaften gegen `AGENTS.md` §5.** Die
  genannten Proben sind gelaufen und die Zahlen stimmen: `make gates` grün mit
  zehn Gliedern und 706 Dateien (eigener Lauf), `make doc-check` 706/0,
  37 Register-Verzeichnisse, vier Dateien mit einem Rename in `971aaa3d`,
  Skill-Bump 1.14.0→1.15.0. Der einzige Satz, der mehr behauptet, als der Text
  trägt — *„Beide Texte sagen ausdrücklich, worin sie sich von ihrem Nachbarn
  unterscheiden"* in `8f08d48b` —, ist als F-1 geführt und dort nicht doppelt
  gemeldet.
- **geprüft, ohne Befund: `AGENTS.md` §3.1/§3.2/§3.6/§3.9.** Kein Code, kein
  Gate-Skript, keine Schwelle, keine Suppression, keine Workflow-Referenz
  berührt; `make gates` und `make review-coverage` grün.
- **geprüft, ohne Befund: `.d-check.yml` und Makefile unberührt** — der Slice
  ändert keine Modul-Konfiguration und kein Target, die Sensors-Tabelle in
  `harness/README.md` bleibt deckungsgleich (`make gate-consistency` grün).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Abgrenzungs-Zeiger benennt den falschen
Nachbarn · Stand-Feld auf ein Wort verkürzt, das die Quelle ablehnt ·
gespeicherter Zähler neben der Belegliste · Regel ohne benannten Ort ihrer
Erfüllung · Zähler-Stand im eigenen Dokument nicht nachgezogen ·
Tabellen-Frage verengt den Anker, den sie indiziert · Kanon-Abweichung im
Bestand, undeklariert

## Verdikt

**Merge-blockierend:** ja — vier MEDIUM. Die Substanz des Slice trägt: Die
Lücke gegenüber `AGENTS.md` §5 ist echt und aus der Selbstauskunft des Trägers
belegt statt herbeigelesen, der Zwei-Orte-Entscheid ruht auf einem Argument,
das die Gegenprobe übersteht, die Zahlen des Registers sind nachgezählt korrekt
und die Korrektur ist lückenlos propagiert, und der Beanspruchungs-Commit hält
`MR-013` sauber. Blockierend ist, dass die eine Zeile, die diese Substanz vor
einer späteren Dopplungs-Lesart schützen soll, auf den falschen Absatz zeigt
(F-1) — und dass der Slice zwei Formen seiner eigenen Klasse mitträgt: eine
Quelle, die auf ein Wort verkürzt wird, das sie ablehnt (F-2), und einen
gespeicherten Zähler, dessen Drift in dieser Datei bereits dreimal
stattgefunden hat (F-3). F-4 fragt, woran die neue Regel überhaupt als erfüllt
zu erkennen ist; solange das offen ist, ist die Skill-Hälfte ihr einziger
wirksamer Teil.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report
ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat.
