# Slice slice-132: Je Hard Rule eine Antwort — welcher Gate-Lauf trägt sie?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-84-durchsetzung](../welle-84-durchsetzung.md).

**Bezug:** Baseline-Regelwerk
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.11.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(*„Jede Hard Rule liegt in zwei Quadranten … nur in einem ist halb
durchgesetzt"*) und
[`modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin)](../../../../.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md);
[`AGENTS.md`](../../../../AGENTS.md) §3 und §5.

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

`AGENTS.md` §3 trägt neun Hard Rules, dazu die Botschafts-Regel in §5. Dieser
Slice beantwortet für **jede einzelne**: welcher Gate-Lauf trägt ihre
Feedback-Hälfte — mit Target **und** Befund-Code — oder ist sie einseitig?

**Was dieser Slice ausdrücklich nicht behauptet:** dass die Mehrheit einseitig
ist. Belegt sind heute genau **drei** (§3.8, §3.9 und die §5-Zeile, alle aus
welle-83 und dort als ohne Gate ausgewiesen). Für die übrigen ist die Frage
**offen**, nicht beantwortet — und beide Ausgänge sind Ergebnisse.

**Der Beleg ist ein Gate-Lauf, keine Zuordnung.** Ein Target zu nennen, das
plausibel klingt, ist die Klasse, die in welle-82 achtmal gekippt ist
([`BEO-011`](../observations.md)). Wo eine Zeile behauptet, Gate X trage Regel Y,
muss ein **konstruierter Verstoß** X rot färben — genau das *„Bewusste Brechen"*,
das
[`modul-13`](../../../../.harness/baseline/v5.11.0/regelwerk/modul-13-quality-gates.md#adr-zur-fitness-function)
für die Fitness Function verlangt.

## 2. Vorgehen

1. **Je Regel eine Zeile** in einer Tabelle: Regel · tragender Gate-Lauf
   (Target + Befund-Code) · Beleg-Form · Verdikt.
2. **Vier Verdikte, und die Unterschiede sind tragend** — sie zu verwischen
   hieße aufrunden:
   - **gedeckt** — ein Repo-Gate trägt die **ganze** Regel.
   - **teilgedeckt** — ein Repo-Gate trägt einen **Teil**; der Rest ist Urteil.
     §3.3 etwa koppelt Lifecycle und Roadmap über `planning-check`, sagt aber
     auch etwas über **Commit-Zerlegung**, das kein Gate sieht.
   - **werkzeug-lokal** — ein Wächter trägt die Regel, aber **außerhalb** der
     Gates: kein `make`-Target und keine CI rufen ihn, ein Lauf mit einem
     anderen Werkzeug ist ungebunden. Das ist etwas anderes als *teilgedeckt*
     und darf nicht darunter verschwinden.
   - **einseitig** — kein Wächter, nirgends.
3. **Belegen, nicht zuordnen:** je *gedeckt*-Zeile ein konstruierter Verstoß,
   der das genannte Gate rot färbt, und der Rückbau grün.
4. **Die einseitigen ausweisen** — im Regeltext selbst, in der Form, die §3.7
   und §3.9 schon führen (*Auflösungs-Trigger*), ohne Forensik.
5. **Schneiden:** was baubar ist, wird als Slice benannt; was nicht baubar ist,
   bleibt ausgewiesen. Beides kommt in die Welle-Datei §4.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Gates bauen.** Der Zensus misst und schneidet; das Bauen ist Folge.
- **Kein Heuristik-Wächter**, um eine Zeile von *einseitig* auf *gedeckt* zu
  heben. Ein behauptetes Gate ist schlechter als eine ausgewiesene Lücke.
- **Nicht die Regeln ändern.** Wer beim Zensus merkt, dass eine Regel unscharf
  ist, notiert es — geändert wird sie in einem eigenen Slice.

## 4. Definition of Done

- [x] **Elf** Zeilen: neun Abschnitte in §3 tragen **zehn** Regeln (§3.7 bündelt
      zwei), dazu die Botschafts-Regel aus §5. Keine Regel ohne Zeile. Die
      Grenze zum übrigen §5 ist **gesetzt und benannt**, nicht verschwiegen
      (Welle §6).
- [x] Jede Zeile mit Gate-Deckung trägt einen **konstruierten Verstoß mit rotem
      Exit** — vier Proben gefahren, alle Exit-Codes gelesen.
- [x] Alle sechs einseitigen Regeln sind im Regeltext ausgewiesen; fünf
      Ausweisungen sind neu, §3.8/§3.9 trugen ihre schon.
- [x] Der Schnitt steht in der Welle §4: **zwei** baubare Regeln als Slices
      ([slice-134](../in-progress/slice-134-nolintlint.md),
      [slice-135](../open/slice-135-uses-pin-sensor.md)), der Rest als
      ausgewiesen — ihre Durchsetzung wäre ein Heuristik-Wächter.
- [x] `make gates` Exit 0 (neun Glieder, 474 Dateien, 0 Befunde); unabhängiger
      Review ([Report](../../../reviews/2026-08-23-slice-132-hard-rule-zensus-review.md)),
      blockierend mit **drei HIGH**, alle Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Eine Vollständigkeits-Aussage über eine ganze Datei ist die Form, die in
  welle-82 achtmal gekippt ist.** — **Ausgang:** *eingetreten, an der
  Auflösung.* Die Regeln waren vollständig, die **Zählung** nicht: der Commit
  spaltete §3.7 in zwei Regeln und sprach im selben Atemzug von „zehn Hard
  Rules"; ein Absatz überschrieb sechs Kennungen mit „(5)". Neun Abschnitte,
  zehn Regeln, elf Zeilen — drei verschiedene Größen, und ich habe sie vermischt.
- **„Gate genannt" ist nicht „Regel getragen".** — **Ausgang:** *eingetreten,
  zweimal, im selben Slice.* Erstens: die `//nolint`-Probe nannte den falschen
  Linter, der Gate blieb rot — aber aus einem anderen Grund, und beinahe hätte
  das als „Regel greift" gezählt. Zweitens: §3.4 stand als *gedeckt* da, obwohl
  meine Probe nur **eine** seiner zwei Aussagen bricht; die Sprachfreiheit der
  Architektur-Sicht prüft nichts. Beide Male war der Fehler nicht die Messung,
  sondern was ich aus ihr schloss.
- **Der Zensus kann den Zuschnitt der Welle umwerfen.** — **Ausgang:** *nicht
  eingetreten.* Er hat den Zuschnitt **geschärft**, nicht umgeworfen: zwei
  baubare Slices dazu, und die Produkt-Welle bleibt so ungeschnitten wie zuvor —
  die Freshness-Frage hat der Zensus gar nicht berührt.

## 6. Trigger

**Start** (`open` → `in-progress`): [welle-84](../welle-84-durchsetzung.md)
eröffnet, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der Zensus ergibt, dass eine
Hard Rule vor der Durchsetzung **umformuliert** werden muss — dann ist das eine
Auftraggeber-Frage und kein Nachzug.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF), Gate-Landschaft (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) ist zentral — dieser Slice produziert genau
  eine Vollständigkeits-Aussage. [`BEO-012`](../observations.md) ebenso, denn
  jede Zeile zitiert eine Regel und behauptet eine Reichweite; das
  Geltungs-Feld ist zu lesen, nicht der Titel. [`BEO-007`](../observations.md)
  für jeden Beleg-Lauf: der Exit gehört gelesen, nicht hinter eine Pipe.

Slice-ID: slice-132. Betroffene IDs: — (Harness-Dateien; keine Anforderung,
keine ADR, keine Adaption). Module: Harness-Dateien, Gate-Landschaft.
Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Bestandsaufnahme an eigenen Dateien.

## 9. Closure-Notiz (nach `done/`)

### Der Zensus

| Regel | Verdikt | Tragender Lauf | Beleg |
|---|---|---|---|
| §3.1 Docker/make-only | **werkzeug-lokal** | Tool-Call-Wächter `.claude/hooks/pretooluse-command-guard.sh` | kein `make`-Target und keine CI rufen ihn (0 Treffer); er nennt seine Grenze selbst — *Stolperdraht, keine Sandbox* |
| §3.2 Suppression-Verbot | **einseitig** | — | konstruiert: echter Verstoß + `//nolint:unused,gochecknoglobals` ⇒ `make lint` **Exit 0**; `nolintlint` fehlt im Profil |
| §3.3 `git mv` = zwei Commits | **teilgedeckt** | `make planning-check` | konstruiert: Ruhe-Marker bei belegtem `in-progress/` ⇒ `planning-drift`, Exit 2. Die Commit-Zerlegung sieht kein Gate |
| §3.4 Architektur sprachfrei; Straten nie abwärts | **teilgedeckt** | `make doc-check` (Modul `matrix`) | konstruiert: `slice-999` in `spec/spezifikation.md` ⇒ `matrix-forbidden`, Exit 2. Die **Sprachfreiheit der Sicht** prüft nichts |
| §3.5 ADRs immutable | **gedeckt** | `make adr-check` (Modul `vcs`) | konstruiert: Kern-Satz in [ADR-0005](../../adr/0005-modul-layout-hexagon-ordner.md), gestagt ⇒ `core-drift-vcs`, Exit 2 |
| §3.6 Gates nicht ohne ADR lockern | **einseitig** | — | die Regel gilt einem **Akt**; es gibt keinen ruhenden Zustand, den ein Sensor lesen könnte |
| §3.7a Kommentar-Klassen | **einseitig** | — | Urteil, kein `grep`; zwei HIGH-Anker im Reviewer-Skill |
| §3.7b Zustandsfelder | **einseitig** | — | dito |
| §3.8 Modul-Zusagen auf der Ziel-Achse | **einseitig** | — | in welle-83 so zugezogen und dort benannt |
| §3.9 Action-Referenzen SHA-gepinnt | **einseitig** | — | **baubar** — als einzige mit auflösendem Trigger statt *permanent* |
| §5 Botschafts-Regel | **teilgedeckt** | `make trace-check` (Modul `commits`) | konstruiert: Botschaft ohne Kennung ⇒ `commit-untraceable`, Exit 2 — das deckt aber eine **andere** §5-Zeile; die **Reichweite** eines Schlusses sieht kein Gate |

**Eins gedeckt, drei teilgedeckt, eins werkzeug-lokal, sechs einseitig.**

### Was der Zensus wert war

**Der schwerste Fund ist keine Lücke, sondern ein wirkender Umgehungspfad.**
§3.2 verbietet Inline-Suppressions — und die verbotene Direktive **funktioniert**:
`make lint` läuft mit Exit 0 über einen echten Verstoß, wenn ein `//nolint`
darüber steht. Alle anderen einseitigen Regeln sind ungewacht; diese eine ist
aktiv umgehbar, und nichts meldet die Umgehung. Ein Zensus, der nur gezählt
hätte, hätte sie nicht gefunden — sichtbar wurde sie erst durch den
konstruierten Verstoß.

**Zweimal habe ich aus einer richtigen Messung den falschen Schluss gezogen.**
Die erste `//nolint`-Probe nannte den falschen Linter; der Gate blieb rot, aber
aus einem anderen Grund — beinahe hätte das als *Regel greift* gezählt. Und §3.4
stand als *gedeckt* da, obwohl die Probe nur eine seiner zwei Aussagen bricht.
Beide Male stimmte der Exit-Code, und beide Male war der Schluss daraus zu weit.
Genau das steht als Risiko in §5 dieses Slice, geschrieben bevor es eintrat.

**Die Zählung war die dritte Falle.** Neun Abschnitte, zehn Regeln, elf Zeilen —
drei verschiedene Größen. Ich habe sie in einem Satz vermischt und dabei
gleichzeitig sechs Kennungen mit „(5)" überschrieben. Die Auflösung eines
Abschnitts in seine Regeln ist die slice-127-Lehre; sie anzuwenden **und** dann
nach Abschnitten weiterzuzählen ist ihre halbe Anwendung.

**Und eine Grenze, die nicht dastand, sah aus wie eine Auslassung.** Der Zensus
nimmt aus §5 nur die Botschafts-Regel. Das ist begründet — `modul-09`s
Zwei-Quadranten-Satz gilt der **Hard Rule**, und die übrigen §5-Zeilen sind
Konventionen —, aber der Grund stand nirgends. Er steht jetzt in §6 der Welle,
und ein eigener Zensus über §5 ist damit **benannt** statt vergessen.

**Offen und benannt:** Sechs Regeln bleiben einseitig, zwei davon werden in
dieser Welle gebaut. Für die vier übrigen ist die Einseitigkeit kein Mangel,
sondern das Ergebnis: ihre Durchsetzung wäre ein Heuristik-Wächter, und `modul-13`
sagt dazu, dass ein abgeschalteter Wächter schlechter ist als ein löchriger.
