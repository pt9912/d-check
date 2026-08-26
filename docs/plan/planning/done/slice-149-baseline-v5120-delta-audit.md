# Slice slice-149: Delta-Audit v5.12.0 und Freshness-Audit aller Adaptionen — Etappe B

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md).

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
      ([slice-150](../open/slice-150-pin-gebundene-zitate.md),
      [slice-151](../open/slice-151-urteilsfreie-haelfte-voll.md)), **zwei
      belegt folgenlos** mit Fundstelle (§9).
- [x] **Alle 16** aktiven Adaptionen gefragt, die Reduktion **gemessen** statt
      angenommen: je Eintrag geprüft, ob sein `Ersetzt-Baseline-Regel`-Ziel
      unter den fünf geänderten Dateien liegt (4 von 16), und bei diesen vier,
      ob die neuen Zeilen in **seinem** Abschnitt landeten (1 von 4).
- [x] [`MR-033`](../../../../harness/conventions.md#mr-033) **bleibt** — und die
      Begründung stand schon im eigenen Speicher (§9).
- [x] Die Differenz zur vollen urteilsfreien Hälfte ist benannt und als
      [slice-151](../open/slice-151-urteilsfreie-haelfte-voll.md) geschnitten,
      nicht hier gebaut — das Audit schneidet (§3).
- [x] `make gates` Exit 0 (zehn Glieder), `make fullbuild` Exit 0.
- [ ] **Unabhängiger Review — ausstehend.** Der Haken stand hier angekreuzt,
      bevor ein Review gelaufen war; die Berichtigung ist Teil der Akte, nicht
      der Formulierung.

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
| 1 — urteilsfreie Hälfte benannt | `modul-05`, +11 Zeilen | **Handlung:** [slice-151](../open/slice-151-urteilsfreie-haelfte-voll.md) — der Kanon nennt *dass ein Ausgang dasteht* **und welcher der drei*; unsere Regel prüft den häufigsten Auslöser |
| 2 — *Modul-Pfad* = **Code**-Modul-Pfad | `modul-03` + `AGENTS.template.md` | **Handlung:** [slice-150](../open/slice-150-pin-gebundene-zitate.md) — nicht wegen der Lesart (die bestätigt uns), sondern wegen des **Zitats**, das mitgehoben wurde |
| 3 — das Rot muss von *dieser* Regel kommen | `modul-13`, +8 Zeilen | **folgenlos für den Regeltext:** der Kanon trägt es, und ihn zu duplizieren verstößt gegen [`AGENTS.md`](../../../../AGENTS.md) §1. **Aber** die Klasse ist hier dreimal aufgetreten und hatte keine Registerzeile — [`BEO-017`](../observations.md) angelegt |
| 4 — die Reichweitenfrage als Frage | `grundlagen-source-precedence`, +6 Zeilen | **folgenlos:** die Klasse führt [`BEO-012`](../observations.md) bei Zähler 4, und [slice-147](../open/slice-147-reviewer-anker-reichweite.md) liegt für die Feedforward-Hälfte bereits in `open/` |

### Freshness-Audit — 16 gefragt, einer betroffen

Die Verengung ist **gemessen**, nicht geraten, und läuft über zwei Schritte:
zeigt der `Ersetzt-Baseline-Regel`-Zeiger in eine der **fünf** geänderten
Dateien (**4 von 16**: [`MR-007`](../../../../harness/conventions.md#mr-007), [`MR-013`](../../../../harness/conventions.md#mr-013), [`MR-032`](../../../../harness/conventions.md#mr-032), [`MR-033`](../../../../harness/conventions.md#mr-033)), und landeten die
neuen Zeilen im **Abschnitt** dieses Eintrags (**1 von 4**)? [`MR-007`](../../../../harness/conventions.md#mr-007) zeigt auf
`modul-13` §Hard Rule, die Änderung liegt in §Fitness Function; [`MR-013`](../../../../harness/conventions.md#mr-013) zeigt
auf `modul-05` §Lifecycle als State Machine, die Änderung in §Offene Risiken;
[`MR-032`](../../../../harness/conventions.md#mr-032) zeigt auf den Absatz *Wann die CR-Pflicht beginnt*, der unverändert
ist. Gegengeprüft mit einer **zweiten** Methode — einer normalisierten
Zitat-Probe über alle Einträge —, die dasselbe Ergebnis liefert.

**[`MR-033`](../../../../harness/conventions.md#mr-033) bleibt, und die Begründung stand schon im eigenen Speicher.** Die
Frage war, ob unser Verbot nach der Klärung noch eine Verschärfung ist oder nur
eine nicht ausgeübte Erlaubnis — der Kanon sagt jetzt ausdrücklich *„Die
Erlaubnis ist keine Pflicht"*. [`MR-032`](../../../../harness/conventions.md#mr-032)
hat genau diese Abgrenzung schon einmal geprüft und **verworfen**, mit Verweis
auf [`MR-031`](../../../../harness/conventions.md#mr-031): *wer verschärft,
weicht ab, auch wenn er nur mehr verlangt.* Der Eintrag bleibt also — und seine
Prämisse ist jetzt **bestätigt** statt angenommen, was ihn stärker macht als
vorher.

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
[slice-150](../open/slice-150-pin-gebundene-zitate.md). Die bequeme Antwort
wäre die zweite; genau darum bekommt sie einen eigenen Slice statt eines
Nebensatzes hier.
