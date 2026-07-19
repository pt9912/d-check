# Slice slice-072: Handbuch — Aufgabenorientierung der §4-Kapitel

**Status:** Done (welle-65, abgeschlossen 2026-07-19).

**Welle:** welle-65-handbuch-aufgaben (Trigger: WIP-Slot frei nach welle-64; Nutzer-Aufnahme 2026-07-19).

**Bezug:** Redaktionelle Korrektur des
[Benutzerhandbuchs](../../../user/benutzerhandbuch.md) gegen den
[Benutzerhandbuch-Standard](../../../user/benutzerhandbuch-standard.md) §2
(„Aufgaben statt Funktionen beschreiben") und §5 („Schritt-für-Schritt-
Anleitungen"). **Kein Change Request** (kein Verhaltens- oder Vertragsdelta —
das Produkt bleibt unberührt), **kein ADR** (keine neue Architekturentscheidung),
**kein Release** (nur abgeleitete Nutzer-Dokumentation). Betrifft ausschließlich
`docs/user/benutzerhandbuch.md`.

**Autor:** pt9912. **Datum:** 2026-07-17.

---

## 1. Ziel

§4 des Benutzerhandbuchs ist stellenweise **funktions-** statt aufgabenorientiert:
Blöcke tragen Config-Schlüssel als Titel, begründen eine Fähigkeit aus d-checks
**Produkt-Innensicht** (welches Modul was tut oder nicht tut) und zählen Felder
auf, wo der Leser eine Handlungsfolge braucht. Der Standard verlangt das
Gegenteil: *„Schreibe nicht, was die Software kann. Schreibe, wie der Nutzer damit
seine Arbeit erledigt."*

Der Slice zieht die betroffenen §4-Blöcke auf die Ziel-zuerst-Form nach —
Ausgangslage → Ziel → Vorgehen → Ergebnis — und räumt die dabei sichtbare
Duplikation zwischen §4 (Aufgabe) und §5 (Referenz) auf.

**Nicht** Gegenstand: §5 und §6. Sie sind laut Standard §3 (Punkt 5
„Konfiguration") legitime **Referenz**-Kapitel; eine Feld-/Optionsreferenz ist
dort korrekt. Sie werden nur dort angefasst, wo §4-Inhalt zu ihnen wandert.

## 2. Entscheidungen / Regel

- **Schablone ist der bestehende Bestand, nicht eine neue Erfindung.** §4.10
  (`--repair`) und der §4.12-Kreuzverweis-Block („Prüfen, ob Ihre RTM-Tabelle noch
  zum Design passt", slice-071) erfüllen den Standard bereits: Ausgangslage → Ziel
  → nummerierte Schritte → Ergebnis-Deutung → Fallstricke als Hinweis. Sie sind
  die Vorlage; es wird kein neuer Stil eingeführt.
- **Titel benennen ein Nutzerziel, das Flag/der Schlüssel steht in Klammern**
  (gelebtes Muster von §4.9/§4.10).
- **Referenz gehört nach §5, Aufgabe nach §4.** Wo §4 heute Grammatik-, Feld- oder
  Fail-closed-Enumerationen trägt, die §5 bereits führt, wird §4 gekürzt und
  verweist; **kein Text wird verdoppelt**.
- **Kompatibilitäts-/Versionsprosa gehört nicht in eine Aufgabe** — Fehlerbilder
  nach §7, Versionsaussagen nach §11.
- **Die verifizierten Beispiele bleiben verifiziert.** Jede Umstellung hält die
  beiden Handbuch-Harnesse grün: Config-Beispiele gegen den Validator
  (slice-061) und Ausgabeblöcke E2E gegen das echte Binary (slice-062). Ein
  Ausgabeblock verliert weder seinen Anker noch seine Form-Token.

## 3. Definition of Done

**Frisches §4-Audit (2026-07-19, [Report](../../../reviews/2026-07-19-slice-072-paragraph4-audit.md)):**
die B-Items sind gegen das *aktuelle* §4 re-validiert (dieser Plan war vom
2026-07-17, §4 wuchs seither). Ergebnis: **B-8 bereits erledigt**, **B-7
niedrig/teilweise**, **B-1…B-6 offen** (B-3/B-6 schlimmer geworden), plus
**N-1…N-4** (neue Funktions-Referenzen, die die Post-Audit-Slices 074–077 an
§4.12 anhängten — exakt die in §4 genannte Ursache). Zeilen unten = Stand
2026-07-19.

- [x] **B-1/B-2 (§4.7 `matrix`, Z. 333–382):** `order`/`direction` und `token` auf Ziel-Titel + Lesersituation umstellen (heute produkt-innensichtig: „Trägt eine Klasse `order` …" / „`matrix` sieht standardmäßig nur Markdown-Links").
- [x] **B-3 (§4.12 „Unterstützte Definitionssyntax", Z. 595–662 **vor** dem „Vorgehen" Z. 664 — HIGH, auf ~67 Z. gewachsen):** Grammatik-/Feldreferenz kürzen; die Tabellen-Binding-Regeln leben nur noch in §5 (Z. 1395–1408), §4 verweist. **Inkl. N-1 „Tabellengrenze" (Z. 648–653) + N-2 „Direktiven-Marker" (Z. 655–662)** — Referenz-Prosa nach §5/§7.
- [x] **B-4 (§4.12 Modalität, Z. 776–805):** als Aufgabe formulieren („Nur MUSS-Anforderungen gaten, SOLLTE/KANN advisory"); Klassifikations-Mechanik + Fail-closed-Liste nach §5 (Z. 1454–1473).
- [x] **B-5 (§4.12-Versionsgrenze — proliferiert):** der „Bis v0.42.0 …"-Block (Z. 699–718) als Fehlerbild nach §7, Versionsaussage nach §11; **plus die ≥5 verstreuten inline „Bis vX/Ab vX"-Notizen (Z. 618, 645–646, 652–653, 661–662, 774, 807) + N-3 „Zugesagte Range-Notationen" (Z. 765–774) + N-4** — Versions-/Kompatibilitäts-Prosa raus aus der Aufgabe.
- [x] **B-6 (§4.12 strukturell, Z. 585–915 ≈ 330 Z., ~8 Themen — HIGH):** auftrennen — RTM erzeugen/gaten bleibt §4.12, Coverage/Modalität/Kreuzverweis werden eigene §4.x mit eigener Ziel/Vorgehen/Ergebnis-Klammer; doppelte WAISE-Definition (Z. 589–591 / 671–673) entfernen. **Korrektur zum Alt-Audit:** die „Kreuzverweis-Konsistenz"-H2 (Z. 883) und „F-1" (Z. 601) sind **Fence-Beispielausgaben**, keine echten promovierten Überschriften — der strukturelle Kern gilt dennoch.
- [x] **B-8 (bereits erledigt):** die §4.x-Titel folgen schon dem Muster „Ziel + Schlüssel in Klammern" (§4.5/§4.11/…) — keine Arbeit nötig.
- [x] **B-7 (§4.4-Öffnung, niedrig/teilweise):** der Ziel/Voraussetzung-Rahmen ist gut (war schon zur Audit-Zeit so); nur die Fähigkeits-Inventar-Passage **vor** der Entscheidungsfrage (frisches vs. bestehendes Repo) noch straffen — kleine Severity.
- [x] **Beobachtungen (Ermessen):** N-5 (die neuen Module `citations`/`sources` haben keine eigene §4-Aufgabe — ggf. je eine schlanke Aufgabe), N-6 (leichte §4.9/§4.11-Überlappung bei `--json`/`--yaml`).
- [x] **Kein Regress:** beide Handbuch-Test-Harnesse (`docexamples_test`/`handbook_examples_test`) grün; kein Ausgabeblock/Config-Beispiel verliert Anker/Form-Token; Handbuch-Version + §11-Zeile nachgezogen; `make gates` grün.
- [x] **Abschluss-Gegenprobe:** ein erneutes §4-Audit gegen den Standard meldet keine offenen HIGH/MEDIUM mehr — und benennt ausdrücklich, was bewusst Referenz bleibt (§5/§6, die §4.9-Vergleichstabelle, §4.10 „Was `--repair` behebt", der §4.12-Kreuzverweis-Block Z. 841–914 als Vorbild).

## 4. Risiken / offene Punkte

- **Kein Gate erzwingt Aufgabenorientierung.** Sie ist eine Erkenntnis-, keine
  Laufzeit-Eigenschaft; die Disziplin liegt beim Autor und beim Review. Dieser
  Slice räumt den Bestand auf, verhindert aber den Rückfall nicht (siehe unten).
- **Die eigentliche Ursache ist strukturell, nicht redaktionell:** §4.12 wuchs
  über die Slices 066/067/068/070/071, weil **jeder** Slice seine Fähigkeit an die
  bestehende Aufgabe anhängte, statt eine Aufgabe zu schreiben. Ohne eine Regel am
  Release-Prep-Punkt (analog zur Prosa-Currency-Liste in
  [`releasing.md`](../../../user/releasing.md) §Release-Prep) erodiert §4 erneut.
  **Offener Designpunkt:** ob dieser Slice die Regel gleich mit in die
  Release-Prep-Checkliste schreibt („neuer Handbuch-Abschnitt = eigene Aufgabe mit
  Ziel/Vorgehen/Ergebnis, keine Anhängung") — das wäre die billigste dauerhafte
  Sicherung und wäre konsistent damit, wie die Currency-Blindflecken dort schon
  behandelt werden.
- **Auftrennung von §4.12 (B-6) berührt die Inhalts-Navigation** (§Inhalt,
  Querverweise „siehe §4.12" aus §5/§6/operations.md). Die Verweise sind
  gate-geschützt (`anchors`), ein vergessener Verweis fällt also auf — aber die
  Umnummerierung ist Handarbeit.
- **Grenze zwischen Aufgabe und Referenz ist Ermessen.** Der Slice zieht sie
  bewusst so: Was der Leser braucht, um **einmal** ans Ziel zu kommen, gehört nach
  §4; was er **nachschlägt**, nach §5. Grenzfälle werden zugunsten von §5
  entschieden (kein Text zweimal).

## 5. Trigger

Nutzer-Frage 2026-07-17 während des slice-071-Release-Prep: ob der neu
geschriebene §4.12-Block den Benutzerhandbuch-Standard §2 einhält. Er tat es
nicht — er begründete die Fähigkeit aus d-checks Modul-Landschaft heraus
(„`matrix` prüft die Richtung, `trace.coverage` die Abdeckung, keines vergleicht
die Mengen") statt aus der Lesersituation. Der Block wurde in slice-071
nachgezogen (Commit `d3a8757`); die anschließende Frage „gibt es weitere solche
Stellen?" ergab ein Audit des gesamten Handbuchs mit sieben Befunden (B-1…B-8) —
dieser Slice arbeitet sie ab.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Der Standard führt, die Doku folgt. `docs/user/` ist kein
Brownfield-Bestand mit eigener Spec — die Zielform steht im
Benutzerhandbuch-Standard und ist im Bestand (§4.10, §4.12-Kreuzverweis) bereits
zweimal gelebt; der Slice zieht den Rest nach.

## 7. Closure-Notiz (nach `done/`)

**Abgeschlossen 2026-07-19** (welle-65, reine Doku — kein Release/ADR, Handbuch-Version 1.42).

**Umsetzung.** §4 des Benutzerhandbuchs auf Aufgabenorientierung nachgezogen
(Benutzerhandbuch-Standard §2/§5). Kern: der ~330-Zeilen-§4.12-`--trace`-Monolith
in vier aufgabenorientierte Sektionen aufgetrennt — §4.12 RTM · §4.13 Coverage
(`trace.coverage`) · §4.14 Modalität (`trace.requirements.modality`) · §4.15
Kreuzverweis (`trace.cross-consistency`); `--print-mk` → §4.16. Tabellen-Grammatik,
Tabellengrenze, Direktiven-Marker, Komma-Kurzform und Migration → §5; das
„0 Anforderungen"-Fehlerbild → §7; verstreute „Bis vX/Ab vX"-Prosa raus aus §4
(§11 führt die Historie); doppelte WAISE-Definition entfernt. §4.7 (`order`/
`direction`/`token`) lesersituations-orientiert (B-1/B-2); §4.4 Entscheidung vor
Modul-Inventar (B-7); `citations`-§5-Titel auf Task-Form, Parität zu `pins`/
`sources` (N-5); §4.9/§4.11-Overlap accept-by-design (N-6).

**Commit-Kette (10).** `0442d3d` (B-1/B-2) · `c17c34b` (Cluster B-3…B-6 + N-1…N-4)
· `f6bb802` (Cluster-Review) · `0cc66cf` (DoD B-3…B-6) · `790299f` (Version 1.42)
· `7cf1ba9` (B-7 + N-5) · `58cd2c0` (Gegenprobe) · `53bb6fe` (DoD-Rest) ·
Closure-Move + dieser Body-Commit.

**Review-Verdikte.** Zwei unabhängige, kontext-getrennte, adversariale
Reviewer-Subagenten: der Cluster-Review des §4.12-Umbaus =
[**ACCEPT**](../../../reviews/2026-07-19-slice-072-paragraph4-cluster-review.md)
(F-1 §11-Verweise und F-2 Bestandsinvariante-Satz eingearbeitet, F-3
accept-by-design); die Abschluss-Gegenprobe des gesamten §4 =
[**BESTANDEN**](../../../reviews/2026-07-19-slice-072-paragraph4-gegenprobe.md)
(12 der 14 Ursprungsbefunde erledigt, 2 bewusst Referenz, 0 offen; 2 neue INFO,
kein HIGH/MEDIUM).

**Kein Regress.** Beide Handbuch-Test-Harnesse (`docexamples_test`/
`handbook_examples_test`) durchgängig grün — alle 11 vom Trace-Parsergrenzen-Test
verlangten Teilstrings erhalten, die drei verankerten Ausgabeblöcke intakt, kein
Config-Beispiel/Ausgabeblock verlor Anker/Form-Token. `make gates` grün (265/0).

**Lerneintrag.** (1) „Bewusst separat" braucht ehrliche Attribution: die
Chunk-Grenze war Auftraggeber-Wahl, aber „§11 einfrieren / keine Version" waren
Autoren-Calls — die Umnummerierung brach die `--print-mk`-§-Verweise selbst, das
ist kein frozen-Altbestand, sondern zu fixen. (2) Subagent-Review-Reports tragen
`span-unclosed`, wenn sie inline Triple-Backtick-Fence-Typen in Prosa nennen
(doc-check scannt `docs/reviews/`) → im Reviewer-Prompt Single-Backtick vorgeben.
(3) Die eigentliche Ursache (§4.12 erodierte, weil jeder Folge-Slice seine
Fähigkeit anhängte statt eine Aufgabe zu schreiben) bleibt **strukturell offen** —
der in §4 notierte Designpunkt (Release-Prep-Regel „neuer Handbuch-Abschnitt =
eigene Aufgabe mit Ziel/Vorgehen/Ergebnis") war **kein DoD-Punkt** dieses Slice
(nur eine §4-Risiko-Notiz) und blieb daher offen — er ist als eigener Folge-Slice
[slice-082](../done/slice-082-release-prep-aufgabenregel.md) (welle-66)
aufgenommen, die Entscheidung über die konkrete Umsetzung dort.
