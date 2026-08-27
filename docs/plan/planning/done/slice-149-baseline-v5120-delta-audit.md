# Slice slice-149: Delta-Audit v5.12.0 und Freshness-Audit aller Adaptionen — Etappe B

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../done/welle-85-baseline-v5120-migration.md).

**Bezug:** [`modul-02-harness-bootstrap.md` §Freshness-Audit](../../../../.harness/baseline/v5.12.0/regelwerk/modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2);
[Antwort des Kurses](../../cr/2026-08-26-antwort-regelwerk-v5110.md);
[`MR-033`](../../../../harness/conventions.md#mr-033) (die eine Adaption, deren
Grundlage sich nachweislich geändert hat);
[slice-143](../done/slice-143-structure-abschnitts-skopus.md) (unsere Deckung
der urteilsfreien Hälfte).

**Berührte Spec-Stellen:** — vermutlich keine; ergibt das Audit ein
Produkt-Delta, wächst es in einem **geschnittenen** Slice, nicht hier.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Zwei Audits, die sich eine Frage teilen — *was heißt der neue Stand für uns?*

**Delta-Audit:** je Kurs-Welle (95–98, eine je CR-Punkt) eine Antwort:
Handlung mit Slice-Kennung oder **belegt** folgenlos. Zwei Erwartungen stehen
schon (Wellendokument §1) — sie sind Hypothesen, keine Belege.

**Freshness-Audit:** je **aktiver** Adaption eine Antwort auf die Frage des
Kanons — *Regelt die neue Fassung das, wofür dieser Eintrag angelegt wurde?*
Keiner bleibt ungefragt, auch nicht die offensichtlich unberührten.

## 2. Vorgehen

1. **Delta zuerst, Adaptionen danach** — eine Adaption gegen einen
   ungelesenen Stand zu prüfen, ist Raten.
2. Je Kurs-Welle: den **vendorten** Text lesen, nicht die Antwort des Kurses.
   Sie ist Erwartung; der Beleg ist die Datei.
3. **[`MR-033`](../../../../harness/conventions.md#mr-033) ist der benannte
   Prüffall.** Der Kanon sagt jetzt *„Die Erlaubnis ist keine Pflicht"* — die
   Frage ist, ob unser **Verbot** danach noch eine Verschärfung ist oder eine
   nicht ausgeübte Erlaubnis. Beide Antworten sind zulässig; die Begründung
   entscheidet, nicht die Bequemlichkeit.
4. **Die urteilsfreie Hälfte gegenrechnen:** Der neue `modul-05`-Text nennt
   *dass ein Ausgang dasteht* **und** *welcher der drei*. Unsere Regel prüft den
   häufigsten Auslöser. Die Differenz gehört benannt und die Entscheidung
   getroffen — als **eigener** Slice, falls sie fällt.
5. Alles Geschnittene bekommt eine Kennung, eine Zeile in §4 des
   Wellendokuments und eine Drift-Log-Zeile.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Umsetzung der geschnittenen Arbeit.** Das Audit schneidet, es baut
  nicht.
- **Keine Antwort ohne gelesene Stelle.** „Vermutlich folgenlos" ist kein
  Ergebnis; folgenlos **mit Fundstelle** ist eines.
- **Keine Änderung an eingefrorenen Dokumenten.**

## 4. Definition of Done

- [x] Vier Kurs-Wellen, vier Antworten: **zwei mit Handlung**
      ([slice-150](../done/slice-150-pin-gebundene-zitate.md),
      [slice-151](../in-progress/slice-151-urteilsfreie-haelfte-voll.md)), **zwei
      belegt folgenlos** mit Fundstelle (§9).
- [x] **Alle 16** aktiven Adaptionen gefragt, die Reduktion **gemessen** statt
      angenommen: je Eintrag geprüft, ob sein `Ersetzt-Baseline-Regel`-Ziel
      unter den fünf geänderten Dateien liegt (4 von 16), und bei diesen vier,
      ob die neuen Zeilen in **seinem** Abschnitt landeten (1 von 4).
- [x] [`MR-033`](../../../../harness/conventions.md#mr-033) **bleibt** — und die
      Begründung stand schon im eigenen Speicher (§9).
- [x] Die Differenz zur vollen urteilsfreien Hälfte ist benannt und als
      [slice-151](../in-progress/slice-151-urteilsfreie-haelfte-voll.md) geschnitten,
      nicht hier gebaut — das Audit schneidet (§3).
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0.
- [x] Unabhängiger Review
      ([Report](../../../reviews/2026-08-26-slice-149-delta-audit-review.md)),
      blockierend mit **sieben MEDIUM** und vier LOW, alle elf eingearbeitet.
      **Der Haken stand hier schon einmal, bevor der Review gelaufen war** — die
      Berichtigung bleibt in der Akte (§9), weil sie zur Sache gehört und nicht
      zur Formulierung.

## 5. Abnahme-Punkte / Risiken

- **Ein Audit, das nur bestätigt, was die Antwort sagt, ist kein Audit.** Die
  Kurs-Antwort ist eine gute Erwartung und ein schlechter Beleg; sie kann in
  Details anders gelandet sein, als sie ankündigt. — **Ausgang:** *entfallen —
  geprüft und deckungsgleich.* Alle vier Änderungen sind gegen den **vendorten**
  Text gelesen, nicht gegen die Antwort; sie landen wörtlich dort, wo die
  Antwort sie ankündigt, und der Umfang stimmt bis auf die Datei genau (fünf
  Dateien, keine sechste). Punkt 1 ist wie angekündigt **anders encodiert** als
  von uns beantragt — auch das steht so in der Antwort.
- **Der Freshness-Audit ist lang und langweilig, und genau darum wird er
  abgekürzt.** Wer ihn auf die „offensichtlich betroffenen" Einträge verengt,
  hat die Klasse [`BEO-011`](../observations.md) begangen, bevor er anfängt. —
  **Ausgang:** *entfallen — die Verengung ist gemessen statt geraten.* Alle 16
  wurden gefragt; die Reduktion auf einen betroffenen Eintrag läuft über zwei
  mechanische Schritte (Ziel-Datei unter den fünf geänderten; neue Zeilen im
  Abschnitt des Eintrags) und wurde mit einer **zweiten**, unabhängigen Methode
  gegengeprüft — einer normalisierten Zitat-Probe über alle Einträge.

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-148](../done/slice-148-baseline-v5120-vendoring.md)
in `done/` — ohne gehobenen Pin gibt es keinen Stand, gegen den geprüft wird.

**Rückführungen:** `in-progress` → `next`, falls das Audit ein Produkt-Delta
ergibt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Regeltext (GF), Konventionsspeicher (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-011`](../observations.md) für jede Vollständigkeits-Aussage,
  [`BEO-012`](../observations.md) für jedes Zitat aus dem neuen Stand.

Slice-ID: slice-149. Betroffene IDs: — (kein `DC-`-Bezug). Module:
Harness-Regeltext. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Audit-Arbeit an vendortem Fremdtext.

## 9. Closure-Notiz (nach `done/`)

Geliefert: zwei Audits, vier Antworten, ein betroffener Eintrag, zwei
geschnittene Slices.

### Delta-Audit — vier Kurs-Wellen, vier Antworten

| Punkt | Was landete | Antwort für dieses Repo |
|---|---|---|
| 1 — urteilsfreie Hälfte benannt | `modul-05`, +11 Zeilen | **Handlung:** [slice-151](../in-progress/slice-151-urteilsfreie-haelfte-voll.md) — der Kanon nennt *dass ein Ausgang dasteht* **und** *welcher der drei*; unsere Regel prüft den häufigsten Auslöser |
| 2 — *Modul-Pfad* = **Code**-Modul-Pfad | `modul-03` + `AGENTS.template.md` | **Handlung:** [slice-150](../done/slice-150-pin-gebundene-zitate.md) — nicht wegen der Lesart (die bestätigt uns), sondern wegen des **Zitats**, das mitgehoben wurde |
| 3 — das Rot muss von *dieser* Regel kommen | `modul-13`, +8 Zeilen | **folgenlos für den Regeltext:** der Kanon ist eine kanonische Quelle und wird pro Entscheidung gelesen; ihn in [`AGENTS.md`](../../../../AGENTS.md) zu wiederholen, verbietet dessen §1 ausdrücklich für **diese Datei**. Das ist **kein** Argument gegen eine Regelstelle anderswo — die braucht es nur nicht, weil die Klasse kein Regeldefizit ist, sondern ein Prüf-Verhalten. **Registerwürdig ist sie trotzdem:** dreimal aufgetreten, ohne Zeile — [`BEO-017`](../observations.md) angelegt |
| 4 — die Reichweitenfrage als Frage | `grundlagen-source-precedence`, +6 Zeilen | **folgenlos für diese Welle:** die Klasse führt [`BEO-012`](../observations.md) bei Zähler 4. Ihre Feedforward-Hälfte trägt [slice-147](../open/slice-147-reviewer-anker-reichweite.md) — der ist **wellenlos** und steht darum nicht in §4 des Wellendokuments; die Welle schließt nicht auf ihn |

### Freshness-Audit — 16 gefragt, einer betroffen

Die Verengung ist gemessen, und **der Filter hat eine benannte Blindstelle**.
Er fragt, wohin der `Ersetzt-Baseline-Regel`-Zeiger zielt — die Frage des Kanons
lautet aber, ob die neue Fassung das regelt, **wofür der Eintrag angelegt
wurde**. Das ist der Gegenstand, nicht das Zeigerziel. Wo beides auseinanderfällt,
trägt der Filter nichts bei, und das gilt für **vier** Einträge, deren Feld gar
keine Baseline-Datei nennt ([`MR-034`](../../../../harness/conventions.md#mr-034), [`MR-035`](../../../../harness/conventions.md#mr-035),
[`MR-036`](../../../../harness/conventions.md#mr-036), [`MR-037`](../../../../harness/conventions.md#mr-037)) — sie sind einzeln gelesen und
unberührt: ihre Gegenstände sind die `matrix`-Kante ADR → Welle, die Ablage
ausgehender CRs, die Ablage der Antwort und der Pin selbst. Keiner der fünf
geänderten Texte spricht darüber.

Für die übrigen zwölf lief der Filter, und seine **erste** Stufe ist exakt
reproduzierbar: **4 von 16** zeigen in eine der fünf geänderten Dateien —
[`MR-007`](../../../../harness/conventions.md#mr-007), [`MR-013`](../../../../harness/conventions.md#mr-013), [`MR-032`](../../../../harness/conventions.md#mr-032),
[`MR-033`](../../../../harness/conventions.md#mr-033).

Die **zweite** Stufe fragt, ob die neuen Zeilen die genannte Stelle treffen —
und hier gehört die **Granularität** dazu, sonst reproduziert die Vorschrift
das eigene Ergebnis nicht: auf **Abschnitts**-Ebene wären es *zwei von vier*,
weil die sechs neuen Zeilen von `grundlagen-source-precedence.md` in
§Source Precedence landen, den [`MR-032`](../../../../harness/conventions.md#mr-032) zuerst nennt. Erst die
**Absatz**-Ebene entscheidet ihn: sein Feld nennt den Absatz *Wann die CR-Pflicht
beginnt*, und der ist unverändert. [`MR-007`](../../../../harness/conventions.md#mr-007) zeigt auf `modul-13`
§Hard Rule, die Änderung liegt in §Fitness Function; [`MR-013`](../../../../harness/conventions.md#mr-013)
zeigt auf `modul-05` §Lifecycle als State Machine, die Änderung in §Offene
Risiken — beide schon auf Abschnitts-Ebene entschieden. Bleibt **einer**.
Gegengeprüft mit einer zweiten Methode — einer normalisierten Zitat-Probe über
alle Einträge —, die dasselbe Ergebnis liefert.

**[`MR-033`](../../../../harness/conventions.md#mr-033) bleibt, und die Begründung stand schon im eigenen Speicher.** Die
Frage war, ob unser Verbot nach der Klärung noch eine Verschärfung ist oder nur
eine nicht ausgeübte Erlaubnis — der Kanon sagt jetzt ausdrücklich *„Die
Erlaubnis ist keine Pflicht"*. Der Kanon
stellt die Frage, also gehört sie zuerst an ihm beantwortet: `modul-02`
§Freshness-Audit kennt für einen Eintrag die Ausgänge *bleibt unverändert* ·
*wird enger* · *wird abgelöst* · *ist gegenstandslos*. Hier greift der erste,
denn die Erlaubnis besteht fort und unser Verbot steht weiter gegen sie.

**Die Repo-Präzedenz stützt das, sie ersetzt es nicht.**
[`MR-032`](../../../../harness/conventions.md#mr-032) hat die Abgrenzung *„eine Freiheit nicht nutzen ist
keine Abweichung"* schon einmal geprüft und **verworfen**, mit Verweis auf
[`MR-031`](../../../../harness/conventions.md#mr-031). Dessen Satz lautet vollständig — *„Wer **ihn**
verschärft, weicht von einem **Baseline-Default** ab — auch wenn er nur *mehr*
verlangt"* —, und das `ihn` meint dort Schritt 3 des Agenten-Workflows. Die
Übertragung trägt, weil auch hier ein **Baseline-Default** verschärft wird: der
Kanon erlaubt, wir verbieten. Sie trüge **nicht**, wenn man den Satz ohne seinen
Skopus als allgemeines Prinzip führte — genau die Verkürzung, in der
[`MR-032`](../../../../harness/conventions.md#mr-032) ihn zitiert und in der auch dieser Abschnitt ihn zuerst
gebracht hat.

Der Eintrag bleibt also — und seine Prämisse ist jetzt **bestätigt** statt
angenommen, was ihn stärker macht als vorher.

**Ein Fehlalarm gehört in die Akte, weil sein Grund lehrreich ist.** Ein
zeilenweises `grep` meldete, das Zitat von
[`MR-032`](../../../../harness/conventions.md#mr-032) sei aus dem Kanon
verschwunden — es steht dort, nur **umbrochen**. Beinahe hätte ich eine zweite
betroffene Adaption gemeldet. Das ist dieselbe Klasse, die dieser Slice-Strang
gerade als vierte Spiegel-Klasse registriert hat: Prosa gegen Prosa messen ohne
Normalisierung ist keine Messung.

### Was das Audit NICHT getan hat

Es hat geschnitten, nicht gebaut — auch dort nicht, wo die Reparatur naheliegt.
Das veraltete Zitat in [`MR-033`](../../../../harness/conventions.md#mr-033)
steht noch, weil vorher eine **Regelfrage** zu beantworten ist: Ein
`Accepted`-Eintrag wird nie überschrieben, und ob die Zitat-Korrektur darunter
fällt oder unter die Pflege pin-gebundener Referenzen, entscheidet
[slice-150](../done/slice-150-pin-gebundene-zitate.md). Die bequeme Antwort
wäre die zweite; genau darum bekommt sie einen eigenen Slice statt eines
Nebensatzes hier.
