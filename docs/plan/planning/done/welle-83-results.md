# Welle 83 — Acht Wellen Rückstand aufgeholt, und dreimal eine Quelle für mehr zitiert, als sie sagt — Closure-Notiz

**Welle:** welle-83-baseline-v5110-migration
**Abschluss:** 2026-08-23
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief · Steering-Loop · Register-Lese-Schritt.

Der Baseline-Pin ist von `v5.9.0` (Kurs-Welle 86) auf `v5.11.0` (Kurs-Welle 94)
gehoben, acht Kurs-Wellen sind einzeln beantwortet, und die zwei Wellen mit
Handlungs-Antwort sind abgearbeitet:

- [slice-128](welle-83/slice-128-baseline-v5110-vendoring.md) — **Etappe A:** Bundle
  `v5.11.0` netzlos vendored und gegen `SHA256SUMS` verifiziert, Pin gehoben als
  [`MR-030`](../../../../harness/conventions.md#mr-030), 38 pin-gebundene
  Verweise retargetet, Alt-Baum entfernt. Der gemessene Delta: 52 Dateien, 30
  unterschiedlich, 23 davon reine Versions-Stempel — **fünf** mit echtem
  Regel-Inhalt.
- [slice-129](welle-83/slice-129-baseline-v5110-delta-audit.md) — **Etappe B:** je
  Kurs-Welle 87–94 **eine** Antwort mit Beleg. Zwei verlangen Handlung, sechs
  sind **gemessen** folgenlos (der Regelwerks-Diff berührt nur fünf Dateien),
  keine ist vermutet. Der Vollständigkeits-Zensus aus Kurs-Welle 94 ergab
  **fünf** Fundorte statt des einen bekannten.
- [slice-130](welle-83/slice-130-lastenheft-historie-form.md) — **Etappe C-1:** die
  Lastenheft-Historie trägt die kanonische vierte Spalte `Verweis`, alle 95
  Bestandszeilen `—`; die Spezifikations-Historie bleibt bei zwei Spalten. Die
  eigene Strenge — Bump und Historie schon vor `Accepted`, was der Kanon nicht
  verlangt — ist deklariert und als
  [`MR-032`](../../../../harness/conventions.md#mr-032) geführt.
- [slice-131](welle-83/slice-131-reviewer-skill-waisen.md) — **Etappe C-2:** ein Zensus
  über 18 Aussagen außerhalb der Rangliste; **drei** Waisen nach `AGENTS.md`
  umgezogen (§3.8 Modul-Zusagen auf der Ziel-Achse, §3.9 SHA-gepinnte
  Action-Referenzen, §5 die Reichweite eines Schlusses), ein eigener Befund
  widerlegt, Reviewer-Skill auf 1.10.0.

Dazu, außerhalb der Welle und von ihr entblockt:
[slice-127](slice-127-claude-md-pointer.md) — `CLAUDE.md` steht auf vier Zeilen
und verweist; zwei Hard Rules, die dort allein standen, sind nach `AGENTS.md`
umgezogen.

## Was hat funktioniert?

**Die Reihenfolge „erst pinnen, dann dem Kanon folgen".** Sie stand im
Wellendokument als nicht verhandelbar, und sie hat getragen: kein Slice hat eine
Kurs-Welle vorweggenommen, und das Audit las den neuen Baum, nicht den alten.

**Das Delta-Audit als Messung statt als Vermutung.** „Sechs Wellen folgenlos"
ist hier kein Urteil, sondern eine Ableitung aus dem Regelwerks-Diff — fünf
berührte Dateien, und Wellen ohne Diff können nichts fordern. Das ist die Form,
die eine Migration auditierbar macht.

**Der Kurs hat unseren Konsumenten-CR angenommen.** Kurs-Welle 94 nennt ihn im
CHANGELOG als Auslöser; die Vollständigkeits-Zusage, die daraus entstand, gilt
seit dem Bump für uns selbst — und hat in dieser Welle vier Slices getragen.
Zwei CR-Punkte wurden abgelehnt, einer davon zu Recht mit dem Satz, den wir uns
merken sollten: *ohne Baubarkeit wäre sie ein behauptetes Gate.*

**Vier unabhängige Reviews, alle blockierend, alle eingearbeitet.** Was sie
fanden, waren fast durchweg **Belegketten**, nicht Ausführungsfehler: eine
übersehene fünfte Fundstelle, ein Selbstwiderspruch gegen einen zwei Stunden
alten Eintrag, ein Zitat außerhalb seines Geltungsbereichs, ein unbelegtes
Kriterium. Das ist die Ebene, auf der dieses Repo seine Fehler macht.

## Was ging anders als geplant?

**Die Lehre dieser Welle ist eine einzige Klasse, dreimal begangen: eine Quelle
für mehr zitiert, als sie sagt.**

1. *„Der Widerspruch gehört benannt"* — im Kanon nur für den Fall
   MR-gegen-neue-Baseline, von mir als universale Konfliktregel geführt.
2. [ADR-0047](../../adr/0047-matrix-spec-historie-nicht-provenance-exempt.md) sagte *„die Verweis-Spalte streichen"*, und
   [`.d-check.yml`](../../../../.d-check.yml) trug diesen Satz als Begründung
   einer live geschalteten `matrix`-Verengung. slice-130 machte ihn falsch. Die
   Auflösung brauchte keine neue Entscheidung: der Baseline-Text trägt für die
   Historie **dieselbe** Begründung und ändert nur das Mittel.
3. [`MR-015`](../../../../harness/conventions.md#mr-015) wurde als Wächter gegen „`AGENTS.md` wird zur Sammelstelle"
   zitiert; sein `Geltungsbereich` nennt `AGENTS.md` §1 und die Pin-Einträge.
   Nur sein **Titel** klingt allgemein.

Dazu zweimal die Gegenrichtung derselben Klasse — eine Aussage über den Kanon,
die aus meiner Konstruktion stammte statt aus seinem Text: die Begründung, warum
die Spezifikation bei zwei Spalten bleibt, und das Trenn-Kriterium des Zensus in
slice-131. Beide hielten der Nachprüfung stand oder wurden berichtigt; **keine
war belegt, als sie geschrieben wurde.**

**Zwei Slice-Prämissen waren falsch, nicht ihre Ausführung.** slice-127 wollte
`CLAUDE.md` kürzen und stützte sich auf *„jede Zeile steht schon woanders"* —
falsch, zwei Hard Rules standen nirgends sonst. slice-131 nannte den
Reviewer-Skill den *schwersten* Fundort — falsch, `modul-10` §Ziel-Form **weist**
ihm die repo-konkrete Klassifikation zu, und alle sieben HIGH-Anker buchstabieren
aus. Beide Male hätte die Ausführung des Plans Schaden angerichtet, und kein Gate
hätte es gemeldet. Gefunden hat es einmal ein unabhängiger Review
(slice-127) und einmal die Gegenprobe am Kanon im Slice selbst (slice-131) —
nicht der Review, der danach lief.

**Ein Slice hat seine eigene Regel im selben Commit gebrochen.** slice-130
änderte das Lastenheft ohne Bump und ohne Historie-Zeile, während er in
derselben Datei deklarierte, dass dieses Repo beides seit der ersten Fassung
führt. Kein Gate kann das melden — die Kopplung „Lastenheft geändert ⇒
Historie-Zeile" ist nirgends mechanisiert.

**Ein Herkunfts-Anker war zweimal falsch, bevor er richtig war.** Erst datierte
er einen Slice, dann einen Baseline-Bump — und richtig war, **gar keinen** zu
setzen: `grundlagen-traceability.md` verbietet das Nachrüsten, *„der leere
Zustand **ist** die ehrliche Information."* Dieselbe Frage stellte sich in
slice-131 neu und wurde dort auf Anhieb getrennt beantwortet: `welle-73` und
`welle-82` sind belegte Steering-Loop-Ursprünge, §3.9 bekommt **keinen** Anker.

## Steering-Loop-Einträge

- **`AGENTS.md` hat drei Zuzüge bekommen** — §3.8, §3.9 und eine §5-Zeile. Alle
  drei kommen aus dem Steering Loop bzw. aus dem Vollständigkeits-Zensus, alle
  drei standen zuvor nirgendwo gerankt, und für alle drei ist belegt, dass sie
  sich im Repo nicht duplizieren.
- **Reviewer-Skill 1.10.0:** keine neue Kategorie, sondern zwei Anker, die jetzt
  auf ihr geranktes Zuhause **verweisen** statt die Regel allein zu tragen —
  plus zwei bis dahin quellenlose HIGH-Anker, die ihre Fundstelle nun im Text
  führen. Nachgetragen wurde außerdem eine Pflicht-Eingabe, die der Kanon
  verlangt und der Block nicht nannte.
- **Nicht** mechanisiert: keine der drei neuen Regeln hat ein Gate. §3.8 und die
  §5-Zeile sind Urteilsfragen und bleiben es; §3.9 ist der einzige mit einem
  **auflösenden** Trigger statt *permanent* — ein Sensor auf `uses:`-Pins löst
  sie ab. Nach `modul-09` sind alle drei damit **halb durchgesetzt**; das steht
  hier, statt mit einem Heuristik-Wächter übertüncht zu werden.

## Beobachtungs-Register (Zeiger)

Gelesen zur Closure ([`observations.md`](../observations.md)):

- **[`BEO-012`](../observations.md) neu angelegt**, Zähler **3**, Schwelle
  erreicht — die Klasse dieser Welle: *eine Quelle wird über ihren
  Geltungsbereich hinaus zitiert.* Sie ist **nicht** dieselbe wie
  [`BEO-011`](../observations.md), auch wenn sie danebensteht: dort wird eine
  Regel aus dem Anlass **gebildet**, hier eine bestehende **gedehnt**. Und weil
  ein Zitat wie ein Beleg aussieht, ist es schwerer zu bemerken als eine
  unbelegte Behauptung. Die Trennung ist an der **Antwort** geprüft, nicht am
  Gefühl: `BEO-011` verlangt *zähle den Bestand*, `BEO-012` verlangt *lies das
  Geltungs-Feld, nicht den Titel*.
- **[`BEO-011`](../observations.md) hat zweimal getroffen** — beide Male als
  Aussage über den Kanon, die aus meiner Konstruktion stammte. Zähler
  unverändert bei 3; die Schwelle ist seit welle-82 erreicht.
- **[`BEO-004`](../observations.md) und [`BEO-009`](../observations.md) haben
  ihr geranktes Zuhause bekommen.** Beide Stand-Zellen sagten *„verkörpert im
  Reviewer-Skill"*; das war bis slice-131 der einzige Ort und damit eine Waise.
  Die Zellen sind nachgezogen — der Skill klassifiziert weiter, die Regel steht
  in `AGENTS.md`.
- **[`BEO-010`](../observations.md) bleibt bei 1 und bleibt gate-blind.** Der
  Zensus hat seine `Makefile`-Hälfte geprüft und **entlastet**: *„spiegelt die
  Modul-Liste; wächst die dort, hier nachziehen"* ist die erlaubte
  Kommentar-Klasse **Kopplung**, keine Waise. Die Lücke ist der fehlende Sensor,
  nicht der Kommentar.
- **[`BEO-002`](../observations.md)** traf in slice-128: die Klassen-2-Hebung
  benutzte eine handgeschriebene Dateiliste statt eines Musters und ließ die
  fünfte Fundstelle stehen — dieselbe Datei, dieselbe Zeile wie zweimal zuvor.
  Der Review hat sie gefunden. Zähler unverändert.
- Alle übrigen Einträge unverändert; keine Streichung.

## Folge-Slices

- **Die Durchsetzungs-Welle** (Auftraggeber-Entscheid 2026-08-23): Closure-
  Bedingung ist, dass jede Hard Rule in `AGENTS.md` eine Feedback-Hälfte hat
  oder als einseitig ausgewiesen ist. **Zensus zuerst**, dann Slices. Die drei
  Zuzüge dieser Welle sind ihre ersten Kandidaten — §3.9 der aussichtsreichste,
  weil er als einziger mechanisierbar ist.
- **Die Ränder, die aus welle-82 offen geblieben sind**, bleiben offen: der
  `diagrams.scope`-Rückbau, die 3×-Form von [`BEO-008`](../observations.md) und
  die Frage, ob `citations` eine feine Ventil-Achse bekommt.

## Verifikation

- `make gates` nach jedem Slice Exit 0; `make fullbuild` zu jeder Closure Exit 0
  — Exit-Codes explizit gelesen, nie hinter einer Pipe
  ([`BEO-007`](../observations.md)). Letzter Lauf: 48 Anforderungen, 0 Waisen;
  Closure-Profil 430 Dateien, 0 Befunde.
- **Kein Release.** Das Audit hat keine Produkt-Konsequenz gefunden: der Bump
  berührt Modul-Verhalten und Konfiguration nicht. Das war die im
  Wellendokument gesetzte Bedingung, und sie ist gemessen, nicht angenommen.
- Das Bundle ist **netzlos** verifiziert (`SHA256SUMS` gegen das Release-Asset),
  ohne Handanlegen an den Bäumen.
- **Vier unabhängige Reviews** in dieser Welle, jeder mit blockierendem
  Verdikt: **0 HIGH, 13 MEDIUM, 6 LOW, 1 INFO** — aus den Kategorie-Summaries
  der vier Reporte gezählt, nicht geschätzt. Alle eingearbeitet, jede
  Einarbeitung nachgemessen. Die zwei Reviews zu
  [slice-127](slice-127-claude-md-pointer.md) sind hier **nicht** mitgezählt:
  der Slice lief neben der Welle, nicht in ihr.
