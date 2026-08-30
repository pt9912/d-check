# ADR-0077: `spans` am Closure-Bindepunkt — die Änderung bleibt, die Begründung trägt anders

**Status:** Accepted

**Datum:** 2026-08-30

**Autor:** pt9912

**Supersedes:** [ADR-0076](0076-spans-am-closure-bindepunkt.md)

**Bezug:**
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md) (der
Closure-Bindepunkt und sein zweites Prüf-Profil),
[ADR-0059](0059-closure-waechter-weicht-structure-regel.md) (der ausgewiesene
Preis der Bereinigung),
[`DC-FA-SPAN-001`](../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in)
(das Modul),
[`DC-FA-CLI-012`](../../../spec/lastenheft.md#dc-fa-cli-012--konfigurations-pfad-überschreiben)
(der Profil-Pfad)

**Schärft:** —. Wie die abgelöste Entscheidung ändert diese **keine**
Produkt-Zusage.

**Regeln:** Baseline-Regelwerk
[`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../../.harness/baseline/v5.12.0/regelwerk/modul-04-adrs.md).

---

## Kontext

[ADR-0076](0076-spans-am-closure-bindepunkt.md) hat `spans` an den
Closure-Bindepunkt gestellt. **Die Änderung ist richtig, ihre Begründung war
es nicht** — ein unabhängiger Review hat drei ihrer Aussagen gemessen und
widerlegt. Die tragende darunter:

> *„Der Preis bleibt: der Defekt fällt erst bei `make fullbuild` auf."*
> — ADR-0076 §Verglichene Alternativen

**Gemessen, an einem Klon des Repos mit einem angehängten offenen Fence:**

| Lauf | Ergebnis |
|---|---|
| Hauptprofil (`make doc-check` — in `gates` **und** im `pre-commit`-Hook) | `fence-unclosed`, 608 Dateien, **Exit 1** |
| Closure-Bindepunkt | **derselbe** Befund, 546 Dateien, Exit 1 |

`spans` steht seit langem in `modules:` der [`.d-check.yml`](../../../.d-check.yml),
`doc-check` ist Teil von `gates`, und der `pre-commit`-Hook ruft es. Die
Scan-Menge des Bindepunkts ist zudem eine **Teilmenge**: das Closure-Profil
führt keinen `scan:`-Block, also gelten die Default-Wurzeln (546 ⊂ 608).

**Damit ist der Zuwachs an gefundenen Defekten null.** Es gibt heute keine
Datei, für die der Bindepunkt einen `spans`-Befund meldet, den `gates` nicht
schon beim Commit gemeldet hätte. ADR-0076 §Kontext nennt den vergessenen
Schluss-Fence *„die einzige, die heute niemand sieht"* — das ist falsch.

**Der Slice-Plan hatte es genauer.** Er schrieb: *„Der Closure-Bindepunkt kann
zwei Defekte nicht sehen, die das Produkt **längst findet**."* Auf dem Weg in
die ADR ist die Einschränkung verlorengegangen, und aus „der Bindepunkt sieht
sie nicht" wurde „niemand sieht sie".

## Entscheidung

**1. Die Änderung bleibt — mit der Begründung, die trägt.** `spans` läuft am
Closure-Bindepunkt, weil **ein Bindepunkt die Deckung eines fremden Profils
nicht unterstellen darf**. Das ist [`AGENTS.md`](../../../AGENTS.md) §3.8 in
Reinform: *ein Modul verspricht nur über das, was es scannt* — und ein
Prüf-Profil, dessen Zusage von einem **anderen** Profil abhängt, verspricht
über etwas, das es nicht kontrolliert. `make verify-closure-notes` ist ein
eigener Bindepunkt mit eigener Config; wer ihn allein fährt — ein Adopter, der
das Closure-Profil übernimmt, oder dieses Repo nach einer künftigen Änderung an
`.d-check.yml` —, bekam bisher ein stilles Grün.

**2. Der Zuwachs an gefundenen Defekten ist null, und das steht hier.** Diese
Entscheidung kauft **keine** neue Deckung, sondern die **Unabhängigkeit** einer
bestehenden. Wer sie als Fund-Zuwachs verkauft, verkauft sie falsch.

**3. Die drei widerlegten Aussagen aus ADR-0076 gelten nicht mehr.** Im
Einzelnen (Belege in deren `## Geschichte`):
   - *„die einzige, die heute niemand sieht"* — `gates` sieht sie, früher.
   - *„`make gate-consistency` hält dafür den `##`-Hilfetext gegen die Doku"* —
     das Modul `targets` vergleicht Target-**Namen**, keine Beschreibungstexte;
     gemessen bleibt der Lauf grün, wenn man den Hilfetext verfälscht **oder**
     `--enable spans` aus dem Rezept entfernt. **Kein Sensor** hält die vier
     Deklarations-Flächen gegen das Rezept.
   - *„`non-empty: true` … macht denselben Fall laut (gemessen: `section-empty`)"* —
     für den Fall, auf dem der Entscheid steht (Fence am **Dateiende**), wird sie
     **nicht** laut. Die verworfene Alternative ist schwächer, als ADR-0076 ihr
     zugestand; die Verwerfung wird dadurch richtiger.

**4. Die Verortung im Makefile-Rezept bleibt** — aber ohne die Begründung
„ein Ort". Der Kopfkommentar des Profils hatte im selben Commit einen
**zweiten** angelegt (die Modul-Aufzählung in Prosa); er ist entfernt. Dass
kein Sensor Rezept und Deklarations-Flächen gegeneinander hält, ist die
benannte Lücke aus Entscheidung 3 und **keine** Zusage mehr.

**5. Das Modul führt einen dritten Grund-Code**, `span-nested-link`, den keine
Deklarations-Fläche nannte. Er macht den Bindepunkt ebenfalls rot und gehört
in die Risiko-Aussage: *„ein drittes Modul ist ein drittes, das rot werden
kann"* gilt für **drei** Codes, nicht zwei.

## Verglichene Alternativen

**Nur `## Geschichte` an ADR-0076.** Verworfen: die widerlegte Aussage ist
nicht ein falsches Detail in einer richtigen Entscheidung, sondern ihr
**Grund**. Wer den Kern liest, läse die falsche Begründung zuerst und die
Korrektur im Anhang. [`AGENTS.md`](../../../AGENTS.md) §3.5 sieht für diesen
Fall die Folge-ADR vor.

**Die Änderung zurücknehmen.** Ernsthaft erwogen — bei null Fund-Zuwachs liegt
der Verdacht der Zeremonie nahe. Verworfen, weil das §3.8-Argument aus
Entscheidung 1 unabhängig davon trägt: ein Bindepunkt, dessen Zusage an einem
fremden Profil hängt, ist genau die Kopplung, die dieses Repo sonst überall
auflöst. Der Preis der Rücknahme wäre, diese Kopplung wieder einzugehen.

**Den Bindepunkt in `gates` ziehen.** Weiterhin verworfen — das ist die
Bindepunkt-Trennung aus
[ADR-0048](0048-closure-note-struktur-im-planning-modul.md). Der in ADR-0076
dafür genannte Preis („der Defekt fällt erst bei `fullbuild` auf") existiert
allerdings nicht: `gates` fängt ihn ohnehin.

## Konsequenzen

**Positiv.** Der Closure-Bindepunkt ist für diese Defekt-Klasse
**selbstgenügsam**. Verlässt `spans` je das Hauptprofil, oder übernimmt ein
Adopter das Closure-Profil allein, bleibt die Zusage intakt. `spans` prüft
dabei die **ganze** Datei, nicht nur den Closure-Abschnitt.

**Negativ.** Ein drittes Modul am Bindepunkt ist ein drittes, das rot werden
kann — über **drei** Grund-Codes (`fence-unclosed`, `span-unclosed`,
`span-nested-link`). Bestands-Rauschen heute: null.

**Was diese Entscheidung nicht liefert.** Keinen zusätzlichen Fund, keine
Deckung des **wohlgeformten** Spans, der Prosa umschließt (kein Defekt, sondern
Code), und keinen Sensor über die Deklarations-Flächen.

## Fitness Function (falls maschinell prüfbar)

`make verify-closure-notes` — der Bindepunkt selbst; er meldet den offenen
Fence, gemessen an einem echten `done/`-Slice: unverändert grün (Exit 0), mit
angehängtem offenem Fence rot (`fence-unclosed`, Exit 1), während
`planning`+`structure` allein bei 0 Befunden bleiben.

**Kein Sensor** hält, dass das Rezept `--enable spans` führt, und **keiner**
hält die vier Deklarations-Flächen dagegen. Das ist gemessen (Mutation an
Hilfetext und Rezept, beide Male grün) und benannt, nicht zugesagt.

## Re-Evaluierungs-Trigger

**Wenn `spans` das Hauptprofil verlässt.** Dann wird aus der Unabhängigkeit
eine Deckung, und Entscheidung 2 ist neu zu lesen.

**Wenn ein Sensor entsteht, der Rezept und Deklarations-Flächen gegeneinander
hält.** Dann verliert die Lücke aus §Fitness Function ihren Gegenstand.

**Wenn `spans` am Bindepunkt Bestands-Befunde erzeugt, die niemand einplanen
wollte.** Dann ist zuerst der Bestand die Frage und nicht das Gate.
