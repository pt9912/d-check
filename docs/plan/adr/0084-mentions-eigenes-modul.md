# ADR-0084: Die Erwähnungs-Deckung wird ein eigenes Modul, kein dritter Quell-Modus der Kreuzprüfung

**Status:** Accepted

**Datum:** 2026-09-06

**Autor:** pt9912

**Bezug:** eingehender Change Request des Adopters `ai-harness-init`
([`docs/plan/cr/2026-09-06-cr-eingehend-ai-harness-init-rtm-soll-ist.md`](../cr/2026-09-06-cr-eingehend-ai-harness-init-rtm-soll-ist.md)),
Vorschlag **A** angenommen; [ADR-0031](0031-targets-deklarations-konsistenz-modul.md)
(die nächstverwandte Deklarations-Achse), slice-205 <!-- d-check:status-provenance -->.

**Schärft:**
[`DC-FA-MENT-001`](../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
(neue Anforderung — diese ADR trägt ihre Begründung und die zwei Entscheide,
die sie voraussetzt).

**Regeln:** Baseline-Regelwerk
[`modul-13-quality-gates.md` §Gate-Typ ↔ Fehlerbild](../../../.harness/baseline/v6.3.1/regelwerk/modul-13-quality-gates.md#gate-typ--fehlerbild).

---

## Kontext

Ein Adopter meldet: Die RTM sagt, ob eine Anforderung **verfolgt** ist — sie
findet keinen Fall, in dem ein Artefakt existiert und im Ist-Dokument
**nicht vorkommt**. Sein Beleg ist ein reales Handbuch, dem eine Fähigkeit
fehlte, ohne dass irgendein Gate rot wurde.

Die Prämisse ist auch hier wahr: d-checks eigenes Handbuch trägt **fünf**
`DC-*`-Nennungen bei inzwischen 52 Anforderungen. Die Konvention *„das
Ist-Dokument zitiert Kennungen"* existiert in diesem Repo nicht, also kann
`trace.coverage` die Achse prinzipiell nicht bedienen.

**Die Inventur trägt die Entscheidung, nicht der eine Anlass.** Zehn Repos
unter demselben Wurzelverzeichnis führen das Paar Soll (`spec/lastenheft.md`)
und Ist (`docs/user/`). Ihre Soll-Seiten führen **fünf verschiedene
ID-Schemata**. Das ist das stärkste Argument für die pfadbasierte Bauform und
steht im CR nicht: Eine kennungsbasierte Lösung müsste jedes der fünf kennen;
eine pfadbasierte keines.

**Zwei Entscheide sind vorab zu treffen**, weil beide die Bauform festlegen und
beide einen naheliegenden, falschen Default haben.

## Entscheidung

### (a) Eigenes Modul `mentions`, keine dritte Quell-Form von `trace.cross-consistency`

[`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)
ist der nächste Verwandte und sagt über sich selbst: *„Die einzige neue Logik
ist: eine Sicht invertieren und die Mengen diffen — kein neuer Parser."*
Genau derselbe Satz disqualifiziert die Wiederverwendung: Diese Achse
**braucht** einen neuen Leser — Pfad-Globs gegen den Dateibaum und eine
Textsuche im Dokument. Was bei XREF gespart wurde, fiele hier an.

Drei weitere Gründe, jeder für sich hinreichend:

- **Es gibt keine gemeinsame Schlüsselmenge.** XREF ist strukturell **je
  Anforderungs-Kennung `R`** aufgebaut: `F(R)`, `B(R)`, der Diff läuft über
  `keys(F) ∪ keys(B)`, und `exclude-req`, `forward.req-pattern` sowie
  `backward.req-pattern` adressieren dieselbe Kennung. Die neue Achse hat
  **kein `R`**. Sie einzuhängen hieße, den Schlüssel optional zu machen — und
  XREFs eigene Zusagen (Vakuum-Prüfung, Namensraum-Vorbedingung) sind über
  diesen Schlüssel **definiert**. Sie müssten neu gefasst werden; damit wäre
  eine `Accepted`-Anforderung geändert, die dieser Slice ausdrücklich
  unangetastet lässt.
- **Die verglichenen Elemente sind verschieden.** XREF diffed **Kennungen**,
  die auf beiden Seiten mit **demselben** `design-pattern` extrahiert werden —
  eine ausdrückliche Vorbedingung der Anforderung. Hier gibt es kein Muster zu
  teilen: Die Frage lautet, ob eine Zeichenkette in Fließtext vorkommt.
- **XREF hängt an `--trace`.** Es ist ein Zusatz zum RTM-Lauf und gatet nur
  unter `--require-complete`. Die neue Achse hat mit der RTM nichts zu tun; sie
  steht daneben.

**Auch nicht `targets`** ([ADR-0031](0031-targets-deklarations-konsistenz-modul.md)):
Das ist die nächstverwandte Achse nach *Zweck* — Doku gegen Bestand —, aber
nicht nach *Gegenstand*. `targets` vergleicht gegen **Build-Regeln**, diese
Achse gegen den **Datei-Bestand**. Die Erweiterung hätte ein Modul mit zwei
Bestandsquellen erzeugt.

### (b) Ein Befund ist ein Befund — kein modul-lokaler Berichts-Modus

[`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes) kennt
drei Exit-Codes und keinen advisory-Zustand für ein Modul. Ist `mentions`
aktiviert, führt ein Fund zu Exit 1 wie bei jedem anderen Modul.

Das Muster *Bericht per Default, Verdikt per Konfiguration* existiert — aber
**eine Ebene höher**: `--trace` ist ein **Ausgabe-Modus**, kein Modul, und
XREF erbt seine Advisory-Natur von ihm. Ein modul-**lokaler** Schalter wäre der
erste seiner Art und berührte die Modul-Architektur; der Slice-Plan hat genau
das als Rückführungs-Grund benannt.

**Der Berichts-Hebel ist bereits doppelt vorhanden** und braucht keinen
dritten: Das Modul ist **opt-in**, und es kann als eigener Fokus-Lauf außerhalb
von `gates` fahren — dieselbe Bauform, die `review-coverage`, `trace-check` und
`adr-check` heute leben. Wer berichten will, aktiviert das Modul in einem
eigenen Target; wer gaten will, nimmt es in `gates`. Beides ist eine Zeile im
`Makefile` und keine Produkt-Fläche.

**Das widerspricht dem Einreicher, und der Widerspruch gehört benannt.** Sein
CR sagt ausdrücklich: *„Ein Dokument kann eine Fähigkeit korrekt beschreiben,
ohne ein Artefakt beim Namen zu nennen; die Differenz ist ein **Hinweis**,
kein Befund."* Das Argument trifft zu — nur zieht es die falsche Konsequenz.
Die Unschärfe sitzt in der **Menge**, nicht im Exit-Code: Enthält die
Soll-Menge Mitglieder, deren Nicht-Nennung kein Mangel ist, ist die Menge
falsch gewählt und nicht das Verdikt zu streng. Ein Berichts-Modus machte
**genau diese falsche Wahl dauerhaft bequem** — die elf gemessenen
Nicht-Mängel aus `tools/**/*.sh` blieben als ewige Hinweiszeile stehen, statt
die Menge zu korrigieren. Das Verdikt ist deshalb der Druck, der die Kuratierung
erzwingt; es zu entschärfen nähme der Achse ihre einzige Qualitätssicherung.

### Was aus diesen beiden folgt

Die **Mengen-Wahl** ist damit das einzige verbliebene Urteil — und sie liegt
beim Konfigurierenden, nicht beim Modul. Gemessen am eigenen Bestand: die 22
Regelmodule gegen `docs/user/` liefern **null** Funde; `tools/**/*.sh` gegen
dieselbe Ist-Menge **elf**, von denen **keiner** ein Mangel ist — Harness-
Skripte gehören nicht ins Benutzerhandbuch. Dieselbe Ist-Seite, zwei Soll-
Mengen, und die Differenz ist vollständig eine Frage der Wahl. Deshalb gehört
die **Begründung** zur Menge, wie bei `matrix.exclude-sections` und XREFs
`exclude-req`.

## Verglichene Alternativen

| Alternative | Verworfen, weil |
|---|---|
| **Dritte Quell-Form von `trace.cross-consistency`** | XREF ist über die Anforderungs-Kennung `R` definiert, die diese Achse nicht hat; ihre Vakuum- und Namensraum-Zusagen müssten neu gefasst werden. Damit änderte sich eine `Accepted`-Anforderung, die dieser Slice unangetastet lässt — und der Ersparnis-Satz, mit dem XREF sich rechtfertigt („kein neuer Parser"), gilt hier gerade nicht. |
| **Erweiterung von `targets`** | Nächster Verwandter nach Zweck, nicht nach Gegenstand: `targets` liest Build-Regeln, diese Achse den Dateibaum. Ein Modul mit zwei Bestandsquellen hätte zwei Zusagen unter einem Namen. |
| **Modul-lokaler Berichts-Modus** (`mentions.advisory: true`) | Erster seiner Art; berührte die Modul-Architektur und damit [`DC-FA-CLI-003`](../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes). Der Hebel existiert schon zweifach (opt-in, eigener Fokus-Lauf) — ein dritter ist Produkt-Fläche ohne neue Fähigkeit. |
| **Soll-Menge aus Literalen** statt aus Pfaden (Modulnamen, Build-Targets) | Die beiden ersten Fund-Stichproben benutzten Literale, und das verdeckte, dass die Anforderung über Pfade arbeitet. Die Pfad-Form ist an einem Fremd-Repo gemessen belegbar (4 von 5 Mitgliedern eines Glob-Sets ohne Erwähnung); die Literal-Form braucht ihren **eigenen** Beleg, und für die Makefile-Achse gibt es `targets` bereits. |
| **Dateiname statt Pfad als Erkennungsform** | An sieben Artefakten aus zwei Repos lieferten beide Formen **dasselbe** Ergebnis — kein Fall nannte einen Dateinamen ohne seinen Pfad. Die Wahl fällt deshalb auf das Kollisions-Argument: `README.md` als bloßer Name träfe überall. Die laxere Form bleibt über `mentions.match: basename` erreichbar. |
| **Ist-Seite als ein Dokument**, wie der CR sie beschreibt | Ein Fremd-Repo führt **acht** Dateien unter `docs/user/`. Wer gegen eine prüft, misst an sieben vorbei und meldet Funde, die nur woanders stehen. |

## Konsequenzen

- Ein **22. Regelmodul** — die Modul-Fläche wächst, und mit ihr die Zahl der
  Spiegel, die eine Modul-Aufnahme still mitändert
  ([`BEO-ALL/modulliste-spiegel-ungegated`](../planning/observations/BEO-ALL/modulliste-spiegel-ungegated/observation.md);
  beim Schreiben dieser Anforderung fielen zwei bereits eingetretene
  Drift-Fälle auf).
- Neue Konfigurations-Fläche: `mentions.artifacts`, `mentions.documents`,
  `mentions.match`. Drei Schlüssel, davon zwei ohne Default — fehlen sie, ist
  der Lauf Exit 2, nicht inert.
- Die **Implementierung ist nicht Teil dieses Entscheids** und liegt in einem
  Folge-Slice. Diese ADR legt Bauform und Verdikt-Semantik fest, damit jener
  Slice sie nicht neu verhandelt.
- **Fremdes Rauschen statt eigenem.** d-checks Gegenprobe auf der richtigen
  Menge läuft grün; kalibriert wird die Ausnahme-Klasse am Fremd-Repo. Das
  weicht von der gelebten Praxis ab, jede Regel vor der Aufnahme am **eigenen**
  Bestand zu messen, und ist der ehrlichste Preis dieser Entscheidung.

## Fitness Function (falls maschinell prüfbar)

**Für Entscheid (a) gibt es keine** — dass ein Modul ein eigenes ist und keine
Quell-Form eines anderen, ist eine Struktur-Aussage über den Code, die kein
Sensor als *richtig* prüfen kann.

Für Entscheid (b) prüft `make test` die Verdikt-Semantik über die
Akzeptanzkriterien von
[`DC-FA-MENT-001`](../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in):
ein Fund ⇒ Exit 1, leere Soll- oder Ist-Menge ⇒ Exit 2, kein Block ⇒
byte-identische Ausgabe. **Sie existieren noch nicht** — sie entstehen mit dem
Implementierungs-Slice, und bis dahin ist diese Sektion eine Zusage, kein Beleg.

## Re-Evaluierungs-Trigger

Zwei Bedingungen, jede für sich hinreichend:

1. **Eine zweite Soll-Quell-Form wird gebraucht** (Literale statt Pfade), und
   zwar mit eigenem Beleg. Dann ist zu prüfen, ob `mentions` zwei Quell-Formen
   trägt oder ob die zweite anderswo hingehört — dieselbe Frage, die dieser
   Entscheid für XREF beantwortet, dann für das eigene Modul gestellt.
2. **XREF bekommt eine Quell-Form ohne Anforderungs-Kennung.** Fällt die
   Schlüssel-Bindung dort weg, verliert Grund 1 von (a) seine Grundlage, und
   die Zusammenlegung ist neu zu bewerten.

Ohne eines von beiden: permanent.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-09-06 | Proposed → Accepted |
