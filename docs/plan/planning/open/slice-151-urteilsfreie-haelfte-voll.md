# Slice slice-151: Die urteilsfreie Hälfte, so weit wie der Kanon sie benennt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-85-baseline-v5120-migration](../welle-85-baseline-v5120-migration.md)
— **Etappe C**, geschnitten vom Delta-Audit in
[slice-149](../done/slice-149-baseline-v5120-delta-audit.md).

**Bezug:** [slice-143](../done/slice-143-structure-abschnitts-skopus.md) (die
heutige Deckung), [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md),
[`BEO-015`](../observations.md),
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).

**Berührte Spec-Stellen:** `spec/lastenheft.md` — **falls** die Messung ein
Produkt-Delta ergibt; Bump und Historie dann nach
[`MR-032`](../../../../harness/conventions.md#mr-032).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Das Regelwerk benennt seit `v5.12.0` die urteilsfreie Hälfte der
Drei-Ausgänge-Regel ausdrücklich: urteilsfrei ist, **dass** zu jedem notierten
Risiko ein Ausgang dasteht **und welcher der drei** es ist — die drei sind eine
geschlossene Menge, kein Freitext. Es schließt mit: *„Welches Werkzeug die
urteilsfreie Hälfte prüft, ist Repo-Entscheidung; dass sie eine hat, ist es
nicht."*

Wir haben eine — sie prüft den **häufigsten Auslöser**, den stehengebliebenen
Vorlagen-Platzhalter. Zwei Fälle deckt sie nicht: ein Risiko **ganz ohne**
Ausgang, und ein Ausgang als **Freitext** statt einer der drei Formen. Der
zweite ist genau die Gestalt, in der
[`BEO-015`](../observations.md) auftrat — der erfundene vierte Ausgang.

**Die erste Frage ist eine Messung, keine Konstruktion:** Trägt `structure` die
Aussage *„jedes Risiko in §5 hat einen Ausgang"* überhaupt? Die Bedingungen
wirken auf den **Abschnitts-Text**, nicht je Listen-Eintrag — eine Korrelation
Risiko ↔ Ausgang ist damit womöglich nicht ausdrückbar.

## 2. Vorgehen

1. **Messen, was ausdrückbar ist**, bevor irgendetwas entschieden wird:
   `require-pattern`, `require-all`, `forbid-pattern` gegen die drei Formen —
   und ehrlich benennen, was davon *je Abschnitt* statt *je Risiko* wirkt.
2. **Am Bestand messen**, was eine kandidierende Regel im heutigen `done/`
   melden würde. Die Prüfmenge ist der Bestand von `done/` **zum Zeitpunkt der
   Messung** — beim Anlegen dieses Slice 142 Slice-Dateien; die Zahl ist zu
   messen, nicht aus einem älteren Slice zu übernehmen.
3. Reicht `structure` nicht, ist die Frage ein **Produkt-Delta** — und dann eine
   eigene Anforderung mit ADR, nicht ein Anhängsel.
4. Bewusstes Brechen je gedeckter Form, **Ursache gelesen** — nicht nur der
   Exit (Regelwerk `modul-13`, Schritt 6, seit `v5.12.0` in beiden Hälften).
5. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Urteils-Prüfung.** Ob ein eingetragener Ausgang inhaltlich **trägt**,
  bleibt Urteil — der Kanon sagt das ausdrücklich, und
  [ADR-0059](../../adr/0059-closure-waechter-weicht-structure-regel.md) hat die
  Grenze schon gezogen.
- **Kein Konformitätsdruck.** Der Kurs hat die Ausweitung ausdrücklich als
  unsere Entscheidung und **kein** Konformitätsthema bezeichnet. Ein „wir
  müssen" wäre hier falsch.
- **Keine zweite Mechanik.** Was entsteht, ersetzt oder erweitert die
  bestehende Regel; es tritt nicht daneben.

## 4. Definition of Done

- [ ] Die Ausdrückbarkeits-Frage ist **gemessen** beantwortet, nicht
      angenommen — mit der Grenze *je Abschnitt vs. je Risiko*.
- [ ] Der Bestand ist gemessen; jede Fundstelle geräumt oder ausgewiesen.
- [ ] Bei Umsetzung: je gedeckter Form ein konstruierter Verstoß mit **gelesener
      Ursache**, Rückbau je grün.
- [ ] Bei Nicht-Umsetzung: die Entscheidung steht **mit ihrer Messung** in der
      Closure-Notiz, und [`BEO-015`](../observations.md) trägt sie.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Eine Regel je Abschnitt kann eine Aussage je Risiko nicht treffen.** Fällt
  die Messung so aus, ist die ehrliche Antwort ein Produkt-Delta oder ein
  Verzicht — nicht eine Regel, die weniger prüft und mehr verspricht. —
  **Ausgang:** *(bei Closure)*
- **Der Bestand ist gewachsen und uneinheitlich.** 139 Slices tragen ihre
  Ausgänge in Prosa; eine Form-Prüfung könnte breit rot laufen und damit einen
  Retrofit erzwingen, den niemand beschlossen hat. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Messung ein Produkt-Delta
verlangt, das eine eigene Anforderung braucht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Produkt-Module (GF), Harness-Regeltext (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-015`](../observations.md) ist der Anlass;
  [`BEO-011`](../observations.md) für jede Aussage darüber, was die Regel
  „vollständig" abdecke.

Slice-ID: slice-151. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in).
Module: `structure`, `planning`. Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Messung vor Konstruktion an bestehender
Modul-Mechanik.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
