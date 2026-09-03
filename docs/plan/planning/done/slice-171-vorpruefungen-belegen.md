# Slice slice-171: Die Vorprüfungen belegen ihre Regel, statt sie zu behaupten

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in)
(das Werkzeug, das den Beleg prüft);
[ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
(seine Entscheidung);
[`MR-051`](../../../../harness/conventions.md#mr-051)
(die Bump-Prozedur, deren Geltungsbereich dieser Slice verschiebt);
[`MR-031`](../../../../harness/conventions.md#mr-031)
(die Präzedenz: ein Vorab-Schritt wird benannt statt vorausgesetzt).

**Berührte Spec-Stellen:** — (Planning-Form und Konvention; keine
Produkt-Anforderung berührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-29.

---

## 1. Ziel

**Ein Vorprüfungs-Block behauptet heute, dass gelesen wurde — und niemand kann
es nachprüfen.**

Die drei Blöcke in §7 jedes Slice-Plans (Sub-Area · Beobachtungen · Nachtlauf)
sind **Deklarationen**: Der Autor schreibt hin, dass er geprüft hat. Ob er die
Regel, die den Schritt vorschreibt, überhaupt gelesen hat, steht nirgends und
ist an nichts gebunden.

**Der Anlassfall ist gemessen und liegt in diesem Repo.** In der Sitzung, die zu
diesem Slice führte, liefen drei Slices (168, 169, 170) durch, deren
Vorprüfungs-Blöcke vollständig ausgefüllt waren — während der zuständige
Zyklus-Abschnitt des Regelwerks (`modul-01`, `-05`, `-06`, `-08`) **nicht
gelesen** war. Die Folge: Review und Verifikation fielen aus, das
Beobachtungs-Register wurde nie fortgeschrieben, und drei Slices gingen mit
offenem DoD-Haken nach `done/`. Kein Block war falsch ausgefüllt — sie waren
alle nur *Selbstauskunft*.

**Das Werkzeug für den Beleg liegt im Repo.** `citations` prüft eine
`d-check:cite`-Direktive **wortgleich** gegen die Quell-Spanne, fail-closed, im
inneren Loop ([ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)).
Die **Skills** nutzen es bereits so — `reviewer.md` zitiert `modul-10` mit
Spanne. Was für einen Skill gilt, gilt für einen Vorprüfungs-Block: Wer die
Regel zitieren muss, hat die Datei geöffnet.

**Was der Beleg nicht leistet, und das gehört vorweg:** Ein Zitat belegt
**Zugriff**, nicht **Verständnis**. Es lässt sich kopieren. Es hätte den
Anlassfall trotzdem verhindert — die Spanne stammt aus `modul-05`, und wer sie
korrekt setzt, hat die Datei offen, in der zwei Abschnitte weiter der Zyklus
steht.

## 2. Vorgehen

1. **Die zwei kanonischen Blöcke bekommen je eine Direktive.** Sub-Area-Wahl und
   Beobachtungs-Sichtung stehen beide in `modul-05` §Zwei Schritte vor der
   Modus-Begründung; jeder Block zitiert die Zeilen, die ihn vorschreiben.
2. **Der dritte Block bekommt keine — und das ist der Entscheid, nicht das
   Versäumnis.** Der Nachtlauf-Schritt ist eine Repo-Adaption
   ([`MR-053`](../../../../harness/conventions.md#mr-053)); sein Ziel liegt
   **nicht** unter `.harness/baseline/`. Eine Direktive auf ein repo-eigenes
   Ziel meldet bei **jeder** Änderung des Ziels, nicht nur bei einer
   inhaltlichen — der Preis stünde in keinem Verhältnis.
3. **Der Geltungsbereich von [`MR-051`](../../../../harness/conventions.md#mr-051)
   verschiebt sich, und das braucht einen Nachtrag.** Er sagt heute wörtlich
   *„Nicht `done/` … dort steht keine"*. Sobald ein Slice mit Direktiven
   schließt, steht dort eine — und jeder Baseline-Bump machte eingefrorene
   Belege rot. **Lösung:** `citations.scope` nimmt die drei eingefrorenen
   Verzeichnisse aus. Für `done/` und `conventions/done/` zählt der Beleg zum
   Zeitpunkt seiner Prüfung; `docs/reviews/` hat gar keine Live-Phase und ist
   von Geburt an Lauf-Beleg. **Der Preis:** dort entfällt auch der
   fail-closed-Pfad — vorher nahm eine malformte Direktive den Lauf mit.
4. **Die Konvention wird benannt, nicht nur vorgemacht.** Ein Slice, der die
   Direktiven trägt, ohne dass die Regel irgendwo steht, ist ein Einzelfall —
   genau die Bauform, die dieses Repo als
   [`BEO-011`](../observations.md) führt.
5. **Dieser Slice wendet die Regel auf sich selbst an:** sein eigener §7 trägt
   die Direktiven, und `make doc-check` prüft sie. Der Beleg ist damit nicht
   behauptet, sondern gefahren.
6. `make gates`; **Review** und **Verifikation** als getrennte Läufe; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Lese-Nachweis für das ganze Regelwerk.** Eine Direktive belegt die
  zitierte Stelle, sonst nichts. Wer daraus „der Autor kennt das Regelwerk"
  liest, wiederholt genau den Fehler, den dieser Slice adressiert.
- **Keine Selbstauskunft-Checkliste** („bestätige, dass du X gelesen hast") —
  das wäre eine Deklaration mehr, dieselbe Klasse wie der Block, den sie
  ersetzen soll.
- **Kein Sensor auf offene DoD-Haken.** Die zweite Lücke derselben Sitzung ist
  real und eigenständig; sie gehört in einen eigenen Slice, nicht hier
  angehängt.
- **Keine Umschichtung von `AGENTS.md`.** Der Befund (511 Zeilen, Duplikation
  gegen `harness/README.md`, fehlender Zyklus) steht, ist aber ein eigener
  Schnitt.
- **Keine rückwirkende Nachrüstung der `done/`-Slices.** Sie sind Belege ihrer
  Zeit; eine nachträglich eingesetzte Direktive belegt nichts.

## 4. Definition of Done

- [x] Die zwei kanonischen Vorprüfungs-Blöcke tragen je eine
      `d-check:cite`-Direktive mit wörtlichem Zitat; `make doc-check` prüft sie
      grün, und ein manipuliertes Zitat wird rot **gemessen**, nicht behauptet.
- [x] Der dritte Block trägt **keine** Direktive, und der Grund steht im Slice —
      nicht als Auslassung, sondern als benannter Entscheid.
- [x] `citations.scope` nimmt `done/` aus; der Nachtrag zu
      [`MR-051`](../../../../harness/conventions.md#mr-051) hält fest, dass seine
      Geltungsbereich-Aussage absehbar unwahr wird — heute steht dort noch keine
      wirksame Direktive — und warum der Beleg trotzdem trägt.
- [x] Die Konvention steht an **einer** Stelle (nicht in `AGENTS.md`
      **und** `harness/README.md`) und ist von dort verlinkt.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**;
      **Verifikation** gegen DoD/Spec — beide in eigenen Kontexten.

## 5. Abnahme-Punkte / Risiken

- **Ein Zitat belegt Zugriff, nicht Verständnis.** Es ist kopierbar, und ein
  Autor, der die Spanne aus einem anderen Slice übernimmt, hat nichts gelesen.
  Der Beleg ist damit schwächer, als er aussieht — und ein Beleg, der stärker
  aussieht, als er ist, ist die Klasse, gegen die dieser Slice gebaut ist. —
  **Ausgang:** eingetreten, und der Review hat es gefunden. Die zweite
  Direktive dieses Slice war wortgleich, grün und belegte trotzdem nicht die
  Regel, die ihre Anrede behauptete — sie zeigte auf die Notier-Nebenregel
  statt auf den vorschreibenden Satz. Der Beleg sah stärker aus, als er war.
  **Antwort:** [`MR-054`](../../../../harness/conventions.md#mr-054) verlangt
  jetzt ausdrücklich die **vorschreibende** Zeile, und der Fall steht dort als
  eigener Absatz.
- **Die Direktive bindet die Planung an eine Zeilennummer.** Jeder
  Baseline-Bump verschiebt sie ([`MR-051`](../../../../harness/conventions.md#mr-051)),
  und ab diesem Slice trifft das **jeden neuen Slice-Plan**, nicht mehr nur die
  Skills. Der Bump wird dadurch teurer. — **Ausgang:** weiter offen — der
  Preis fällt erst beim nächsten Bump an und ist bis dahin nicht messbar.
  [`MR-054`](../../../../harness/conventions.md#mr-054) trägt dafür einen
  Beobachtungs-Trigger: melden die Slice-Pläne mehr Spannen-Brüche als die
  Skills, ist über eine stabilere Adresse zu entscheiden (Abschnitts-Anker
  statt Zeilenspanne).
- **`citations` ist fail-closed im inneren Loop.** Eine malformte Direktive in
  einem Slice-Plan beendet den ganzen Lauf mit Exit 2, auch im `pre-commit`-Hook
  — für jeden künftigen Slice-Autor. — **Ausgang:** weiter offen für die
  **lebenden** Pläne; dort bleibt der Pfad, wie er ist. Für die drei
  eingefrorenen Verzeichnisse hat sich die Richtung **umgekehrt** und ist
  gemessen: dort entfällt der fail-closed-Pfad ganz, eine malformte Direktive
  nimmt den Lauf nicht mehr mit. Das ist der Preis des Ausschlusses, und er
  steht in [`MR-054`](../../../../harness/conventions.md#mr-054).
- **Die Regel wird aus EINEM Anlassfall gezogen** ([`BEO-011`](../observations.md)):
  drei Slices einer Sitzung. Ob Vorprüfungs-Blöcke *generell* zu schwach sind
  oder ob dieser Lauf ein Ausreißer war, ist nicht gemessen. —
  **Ausgang:** weiter offen. Der **Anlass** stammt aus einer Sitzung, die
  **Regel** aus dem Bestand — 13 Marker-Treffer, 455 Dateien, 0 Befunde, und
  die Zahlen zu offenen DoD-Haken und Report-Deckung in
  [welle-86](welle-86-closure-uebergang-durchsetzen.md). Ungemessen bleibt
  die eigentliche Frage: ob Vorprüfungs-Blöcke *generell* zu schwach sind.
  Sie wird es erst, wenn mehrere Pläne unter der neuen Form gelaufen sind.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei (slice-170 geschlossen).

**Rückführungen:** `in-progress` → `open`, falls sich zeigt, dass der
`done/`-Scope-Ausschluss andere Direktiven ungeprüft ließe, die dort heute
schon stehen — dann ist der Befund ein anderer, und die Reihenfolge kehrt sich
um (erst Bestand klären, dann Konvention).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `docs/plan/planning/` (Planning-Form) und `harness/`
  (Konventionsspeicher). Beide fallen unter den Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration); eine eigene Deklaration führt nur `tools/harness/`.
  Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:213-214 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-08-29, höchste Kennung
  `BEO-021`): [`BEO-011`](../observations.md) — die Regel aus dem **Anlass**
  statt aus dem Bestand: dieser Slice zieht seine Regel aus **einer** Sitzung,
  und das steht als Risiko in §5; [`BEO-012`](../observations.md) — eine Quelle
  über ihren Geltungsbereich hinaus zitiert: genau die Klasse, die eine
  `cite`-Spanne **eingrenzt**, weil sie den Geltungsbereich mitliefert;
  [`BEO-013`](../observations.md) — ein Wächter, der nichts mehr fängt: die
  Direktive muss beim Bump neu angekert werden, sonst meldet sie irgendwann
  nur noch sich selbst. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.12.0/regelwerk/modul-05-planning-harness.md:219-219 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): beide Achsen melden
  **gruen** — `upstream-drift.yml` zuletzt 2026-08-29T07:34:35Z (der geplante
  Lauf dieses Tages, **6,5 Stunden nach dem Cron `0 1 * * *`**),
  `image-scan.yml` 2026-08-28T15:25:09Z. **Dieser Block trägt bewusst keine
  `cite`-Direktive** — sein Ziel ist eine Repo-Adaption, kein
  Baseline-Abschnitt; die Begründung steht in §2 Punkt 2.

Slice-ID: slice-171. Betroffene IDs:
[`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in),
[`MR-051`](../../../../harness/conventions.md#mr-051). Module: `citations`.
Gates: `make gates`, `make doc-check`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas
(`docs/plan/planning/`, `harness/`) fallen unter den Default: Doc führt, Code
folgt. Konventions-Änderung plus eine Konfigurations-Zeile; kein Fremdsystem,
keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)

**Geliefert.** Die beiden kanonischen Vorprüfungs-Blöcke tragen je eine
`d-check:cite`-Direktive auf die Regelwerk-Zeile, die den Schritt vorschreibt;
[`MR-054`](../../../../harness/conventions.md#mr-054) trägt die Konvention an
einer Stelle, [`AGENTS.md`](../../../../AGENTS.md) §5 und
[`harness/README.md`](../../../../harness/README.md) stellen sie zu, und
`citations.scope` nimmt die drei eingefrorenen Verzeichnisse aus.

**Der Beleg ist gefahren, nicht behauptet.** `make gates` grün, Exit 0, 591
Dateien, 0 Befunde. Die Bruch-Probe zweimal: ein manipuliertes Zitat („sichten"
→ „prüfen") liefert `citation-mismatch` und Exit 1, eine um eine Zeile
verschobene Spanne (`219` → `220`) ebenfalls; unverändert Exit 0. Die
Produkt-Gegenprobe über die drei ausgeschlossenen Verzeichnisse — 455 Dateien,
`modules: [citations]`, ohne jeden `ignore` — meldet 0 Befunde: der Ausschluss
legt **rückwirkend nichts** still, er wirkt nur vorwärts.

**Der Slice hat seine eigene Regel in einem von zwei Fällen verletzt, und das
ist sein wichtigstes Ergebnis.** Die zweite Direktive war wortgleich und grün
und zeigte trotzdem auf die Notier-Nebenregel statt auf den vorschreibenden
Satz — ein Beleg, der stärker aussah, als er war, also genau die Klasse, gegen
die der Slice gebaut ist. Gefunden hat es der **unabhängige Review**, nicht der
Schreibende; [`MR-054`](../../../../harness/conventions.md#mr-054) verlangt seither ausdrücklich die vorschreibende Zeile.
Zweiter Befund derselben Familie: die Regel galt, bevor sie zugestellt war —
`AGENTS.md` verlinkte für die dritte Vorprüfung [`MR-053`](../../../../harness/conventions.md#mr-053) und schwieg zu den
beiden anderen.

**Review und Verifikation liefen in eigenen Kontexten** und liegen als
[Review-Report](../../../reviews/2026-08-29-slice-171-vorpruefungen-belegen-review.md)
und
[Verifikations-Report](../../../reviews/2026-08-29-slice-171-vorpruefungen-verifikation.md)
vor: 1 HIGH, 4 MEDIUM, 2 LOW, 1 INFO plus sechs Abweichungen. Sieben Befunde
sind eingearbeitet; F-8 bleibt als INFO stehen, weil sein Ziel repo-eigen ist —
derselbe Grund, aus dem der dritte Block keine Direktive trägt. Der Anlassfall
dieses Slice war der **Ausfall** dieser beiden Schritte; sie hier zu fahren war
die Probe darauf, dass die Antwort mehr ist als ein Text.

**Was der Slice ausdrücklich nicht liefert:** einen Sensor. Ein Plan ganz ohne
Direktiven ist grün — die Adaption wirkt über Zustellung und Review. Die
Durchsetzung ist [welle-86](welle-86-closure-uebergang-durchsetzen.md), die
pfad-gebundene Zustellung [slice-176](../done/slice-176-planning-rule-pilot.md).
