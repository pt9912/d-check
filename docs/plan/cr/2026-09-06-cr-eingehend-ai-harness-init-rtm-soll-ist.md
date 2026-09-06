# Eingehender Change Request — RTM deckt die Soll/Ist-Achse nicht

**Absender:** Adopter `ai-harness-init` · **Eingegangen:** 2026-09-06
**Richtung:** eingehend — dieses Repo ist der **Empfänger**, nicht der Bittsteller.
**Ziel-Dokument:** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
**Berührt:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) (RTM), [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in) (`trace.coverage`)
**Stand:** **entschieden am 2026-09-06** — Vorschlag **A angenommen**, Vorschlag **B
zurückgestellt**. Umgesetzt als Anforderung
[`DC-FA-MENT-001`](../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
mit [ADR-0084](../adr/0084-mentions-eigenes-modul.md); Begründung unten.
**Eingriff am Wortlaut:** keiner am Text — die vier genannten `DC-*`-Kennungen sind
zusaetzlich als Link auf ihre Definition gesetzt (Kennungs-Linkpflicht dieses
Repos). Der gelesene Wortlaut ist unveraendert.

**Ablage-Hinweis.** [`MR-035`](../../../harness/conventions.md#mr-035) regelt
**ausgehende** CRs, [`MR-036`](../../../harness/conventions.md#mr-036) die
**Antworten** darauf. Ein **eingehender** CR ist eine dritte Klasse, für die
dieses Repo keinen deklarierten Ort führt und der Kanon auch in `v6.3.1`
keinen benennt — er kennt den eingehenden CR als *externen Vorgang*,
ausdrücklich als „bewusst kein Harness-Konstrukt". Die Datei liegt hier nach
derselben Begründung, die
[`MR-035`](../../../harness/conventions.md#mr-035) trägt: Der vorige
Konsumenten-CR ging verloren, und mit ihm die Frage, was genau gebeten und mit
welcher Begründung entschieden wurde. Ob daraus eine Adaption wird, ist offen.

---

## Wortlaut (unverändert übernommen)

> **Betreff:** Die RTM deckt die Soll/Ist-Achse nicht — Abdeckung durch ein
> Ist-Dokument, das keine Kennungen trägt
>
> **Ziel-Dokument:** `spec/lastenheft.md` (Version 0.84.0, Status Draft)
> **Berührt:** [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) (RTM), [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in) (trace.coverage)
> **Einreicher:** ai-harness-init (Adopter), 2026-09-06
>
> ### Anlass, gemessen
>
> Ein Soll/Ist-Vergleich in einem Adopter-Repo brachte eine Doku-Lücke zutage,
> die die RTM nicht sehen kann. Das Lastenheft fordert eine Erfassungsschicht
> (LH-FA-10); der Bootstrap emittiert sie nachweislich:
>
> ```
> grep -n 'erfassung-feldliste' internal/emit/fieldlist.go
> # 17:const FieldListPath = "harness/erfassung-feldliste.md"
> ```
>
> Das Benutzerhandbuch — das einzige Dokument, das den Ist-Zustand beschreibt —
> erwähnt sie mit keinem Wort:
>
> ```
> grep -icE 'span|telemetri|erfassung|aufzeichn|feldliste' docs/user/benutzerhandbuch.md   # 0
> ```
>
> Die RTM meldet für dieselbe Anforderung `ok`. Zu Recht: Sie misst
> **verfolgt** (ADRs/Slices referenzieren sie), nicht **beschrieben**.
>
> ### Warum `trace.coverage` das nicht trägt
>
> [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in) zählt eine Anforderung als gedeckt, „wenn ihre Kennung im
> (abschnitts-gefilterten) Text vorkommt". Diese Annahme hält für ADRs, Slices
> und kuratierte Matrizen — für ein Benutzerhandbuch hält sie nicht, und zwar
> aus einem guten Grund: Es adressiert Adopter, die die internen Kennungen des
> Lieferanten nicht kennen. Kennungen hineinzuschreiben, damit ein Werkzeug sie
> findet, verschlechtert das Dokument für seine tatsächliche Zielgruppe.
>
> Beleg, dass das kein Einzelfall ist: d-checks eigenes Handbuch trägt 5
> DC-Nennungen, unseres 0 — in beiden Fällen weit unterhalb der
> Anforderungs-Menge. Die Konvention „Ist-Dokument zitiert IDs" existiert in
> der Praxis nicht.
>
> ### Was fehlt
>
> Eine **Bestands-Gegenprobe**: nicht „nennt das Dokument die Anforderung?",
> sondern „gibt es ein geführtes Artefakt, das im Ist-Dokument nicht vorkommt?"
>
> Das ist die Richtung, die den Fund oben produziert hat, und sie ist mechanisch
> und urteilsfrei: Die Artefakt-Menge steht fest (was das Repo führt bzw.
> emittiert), das Ist-Dokument ist Text, die Frage ist Mengendifferenz. Die
> Gegenrichtung (Stichwort-Suche „kommt LH-FA-10 sinngemäß vor?") ist dagegen
> ein Urteil — meine 0 oben entstand mit selbstgewählten Stichwörtern, und dass
> 0 Treffer „nicht dokumentiert" heißt, kann kein Muster entscheiden.
>
> ### Vorschlag
>
> Zwei Schnitte, der erste ist der tragende:
>
> **A — Bestands-Gegenprobe (neue Achse).** Eine opt-in Quelle, die eine
> Artefakt-Menge (Pfad-Globs) gegen ein Ist-Dokument hält und meldet, welche
> Mitglieder der Menge dort nicht vorkommen. Ausgabe als Bericht, kein Verdikt.
>
> **B — Zuordnung ohne Kennung im Ziel (kleiner, ergänzend).** `trace.coverage`
> bekommt eine Quell-Variante, deren Zuordnung extern in einer Mapping-Datei
> steht (Kennung → Beleg-Stelle), sodass das Ist-Dokument selbst keine
> Kennungen tragen muss. Löst den deklarierten Fall, findet aber keine
> unbekannten Lücken — deshalb nachrangig.
>
> ### Akzeptanzkriterien (für A)
>
> - **Happy Path:** Given eine konfigurierte Artefakt-Menge und ein
>   Ist-Dokument, when der Lauf fährt, then listet er jedes Mitglied, das im
>   Dokument nicht vorkommt — mit Pfad und der Zeile, an der die Menge es führt.
> - **Boundary (leere Differenz):** Ist jedes Mitglied erwähnt, sagt der Lauf
>   das mit seiner Bezugsmenge („N von N"), nicht bloß „0".
> - **Negativ (kein Gate über dem Urteil):** Der Lauf trägt keinen Exit-Code
>   über die Frage, ob eine Anforderung sinngemäß beschrieben ist. Nur die
>   Mengendifferenz ist mechanisch; alles darüber bleibt Bericht.
> - Determinismus/Seiteneffektfreiheit wie [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix) (read-only, stdout,
>   kein Dokument erzeugt).
>
> ### Abgrenzung
>
> - Kein Ersatz für die RTM — andere Frage, andere Achse: RTM misst
>   **verfolgt**, dies misst **beschrieben**.
> - Kein Vollständigkeits-Verdikt über Prosa. Ein Dokument kann eine Fähigkeit
>   korrekt beschreiben, ohne ein Artefakt beim Namen zu nennen; die Differenz
>   ist ein **Hinweis**, kein Befund.
> - Keine Änderung an `trace.coverage`s bestehender Semantik (Vorschlag B wäre
>   additiv).
>
> ### Nutzen über den Einreicher hinaus
>
> Jedes Repo dieser Baseline führt ein Soll-Dokument (`lastenheft.md`) und ein
> Ist-Dokument (`docs/user/…`), und die Traceability-Achse liegt bereits in
> d-check. Die Lücke ist damit keine Eigenheit eines Adopters, sondern eine
> strukturelle Lücke der Kette Spec → Implementierung → Beschreibung.

---

## Entscheid (2026-09-06)

**Vorschlag A angenommen**, umgesetzt als
[`DC-FA-MENT-001`](../../../spec/lastenheft.md#dc-fa-ment-001--erwähnungs-deckung-einer-artefakt-menge-modul-mentions-opt-in)
(Modul `mentions`, opt-in) mit [ADR-0084](../adr/0084-mentions-eigenes-modul.md).
**Vorschlag B zurückgestellt**, nicht abgelehnt — der CR nennt ihn selbst
nachrangig, und er löst nur den **deklarierten** Fall: Er findet, was jemand
schon als Kante gepflegt hat, und damit gerade nicht die unbekannte Lücke, die
den Anlass produziert hat.

### Was den Ausschlag gab — und es steht nicht im CR

**Die Prämisse des CR trifft zu, aber nicht in der Schärfe, in der sie hier
zuerst stand** (unabhängiger Review, M-1): d-checks **Handbuch** trägt fünf
`DC-*`-Nennungen — über die ganze Ist-**Menge** `docs/user/` sind es **35
Vorkommen** und **24 verschiedene** Kennungen bei 52 Anforderungen. Die
Konvention *„das Ist-Dokument zitiert Kennungen"* existiert hier also
**teilweise**, nicht gar nicht; der erste Wortlaut maß gegen ein Dokument und
schloss über das Repo — genau der Fehler, den die Anforderung eine Zeile
weiter benennt.

Getragen hat die Entscheidung ohnehin etwas anderes: die **Inventur**. Zwölf
Repos unter demselben Wurzelverzeichnis führen das Paar Soll
(`spec/lastenheft.md`) und Ist (`docs/user/`). **Neun** davon führen
kennungstragende Überschriften, in **vier** verschiedenen Schemata (`LH` in
sechs Repos, dazu `HSM`, `DC`, `AC`); **drei** führen auf dieser Ebene gar
keine, im Rumpf dagegen eigene Familien. Genau daraus folgt die Bauform —
eine kennungsbasierte Lösung müsste jedes kennen, eine pfadbasierte keines.
Das ist der eigentliche Grund, warum sich die Achse zu bauen lohnt, und er ist
stärker als der Anlass. **Die Zahlen sind zweimal berichtigt worden** — die
erste Fassung nannte zehn Repos und ein Schema, das es nicht gibt; die zweite
zählte mit einer Methode, die deutsche Komposita für Kennungen hielt. Beide
Berichtigungen stehen in
[ADR-0084](../adr/0084-mentions-eigenes-modul.md) §Geschichte; die Richtung
des Arguments wird dadurch stärker, nicht schwächer.

### Fünf Abweichungen vom CR-Wortlaut, jede gemessen

| CR sagt | Umgesetzt | Warum |
|---|---|---|
| „ein Ist-Dokument" | eine **Menge** von Dokumenten | Ein Fremd-Repo führt acht Dateien unter `docs/user/`; wer gegen eine prüft, misst an sieben vorbei. |
| (nichts über die leere Prüfmenge) | **fail-closed auf beiden Seiten** | Die erste Stichprobe dieses Entscheids meldete „nichts gefunden", nachdem ihr Glob keine Datei getroffen hatte — eine Zusage über null Mitglieder. |
| „die Differenz ist ein **Hinweis**, kein Befund" | **ein Befund ist ein Befund** (Exit 1) | Das Argument trifft zu, zieht aber die falsche Konsequenz: Die Unschärfe sitzt in der **Menge**, nicht im Exit-Code. Gemessen: `tools/**/*.sh` gegen `docs/user/` liefert elf Funde und **keinen** Mangel — das ist eine falsch gewählte Menge. Ein Berichts-Modus machte sie dauerhaft bequem; das Verdikt erzwingt die Kuratierung. Ausgeschrieben in [ADR-0084](../adr/0084-mentions-eigenes-modul.md) §Entscheidung (b). |
| „Ist jedes Mitglied erwähnt, sagt der Lauf das mit seiner Bezugsmenge (`N von N`), nicht bloß `0`" | **übernommen** | In der ersten Fassung dieses Entscheids überhört (unabhängiger Review, M-5). Die Bitte ist gut und kostet nichts: Die Zusammenfassung nennt `N von M Artefakt(en) erwähnt`, über `D` Dokument(e). Sie steht dort und **nicht** im `message`-Feld eines Befundes — dessen Wortlaut sagt [`DC-FA-CLI-004`](../../../spec/lastenheft.md#dc-fa-cli-004--ausgabeformate) ausdrücklich nicht zu. |
| Fundzeile „mit Pfad **und der Zeile, an der die Menge es führt**" | **nicht umsetzbar**, als Grenze aufgenommen | Ebenfalls zunächst überhört (M-5). Eine Soll-Menge aus Pfad-Globs führt ihre Mitglieder **an keiner Zeile**; eine erfundene wäre falsch. Der Befund trägt deshalb den Artefakt-Pfad und den Platzhalter `1` als Zeile, und Out-of-Scope (5) der Anforderung sagt das ausdrücklich. Wer die Herkunft eines Mitglieds sehen will, liest das Glob in der Konfiguration. |

### Was hier **nicht** entschieden ist

Die **Implementierung** — Modul, Grund-Codes, Konfigurations-Schema — liegt in
einem Folge-Slice. Diese Ablage trägt den Entscheid, nicht die Lieferung.
