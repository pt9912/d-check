# Review — slice-087 Spec-§7-Referenzrichtung (C-3-Nachzug)

- **Review-Art:** unabhängiges, kontext-getrenntes Frischkontext-Review
  (adversarial). Gegenstand: die d-check-eigenen Spec-Straten
  `spec/spezifikation.md` und `spec/lastenheft.md` §7 („Historie")
  v5.0.0-konform (keine Abwärtsverweise auf slice-/ADR-Kennungen, auch
  nicht in der Historie).
- **Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Modell:** claude-opus-4-8 (Opus 4.8).
- **Datum:** 2026-08-02.
- **Range:** `341423a..HEAD` — zwei Commits (`70cf38e` Open, `100ec3f`
  Conformance).
- **Gates (selbst gefahren):** doc-check über das Image
  (`docker run --rm --network none -v "$PWD":/repo:ro d-check:latest`) grün:
  302 Datei(en), 0 Befund(e) — matrix läuft mit der verengten
  `exclude-sections` über §7. `make adr-check RANGE=341423a..HEAD` grün:
  302 Datei(en), 0 Befund(e) (Modul vcs, keine Accepted-ADR-Drift im
  Range). Zusatzprüfung: über `341423a..HEAD` ist unter `docs/plan/adr/`
  einzig `0047-…md` neu; keine bestehende ADR-Datei geändert.
- **Prüf-Achsen:** Vollständigkeit des Putzes · korrekt Behaltenes ·
  Prosa-Integrität · ADR-Korrektheit/Immutabilität · Konfig↔Wirkung ·
  Slice-Doc↔Umsetzung. Verifiziert mit grep, `git diff`/`git show` über
  den Range und den beiden Gate-Läufen — nicht gegen Zusammenfassung.

## Findings

### F-1 (LOW) — hängende Semikola in `spec/spezifikation.md` §7 nach Link-Entfernung

- **Kategorie:** LOW (Maintainability / Doku-Drift in der kanonischen Spec).
- **Quelle:** Slice-Ziel §1 („§7 bleibt lesbare `Datum | Änderung`-Chronik"),
  ADR `0047` §Entscheidung 2.
- **Pfad:** `spec/spezifikation.md:2147` und `spec/spezifikation.md:2149`.
- **Befund:** Beide Zeilen trugen im Ausgangszustand die Formulierung
  `… Lastenheft-CR 0.45.0; Begründung in [ADR-0038](…) (Entscheidung 9)`
  bzw. `… Lastenheft-CR 0.44.1; Begründung in [ADR-0038](…) (Entscheidung 8)`.
  Der Putz entfernte die `Begründung in [ADR-0038](…)`-Klausel (echter
  Abwärtslink) samt Verweis-Zelle, ließ aber das vorausgehende Semikolon
  stehen. Die Zellen enden jetzt auf `… Lastenheft-CR 0.45.0; |` bzw.
  `… Lastenheft-CR 0.44.1; |` — ein hängendes Satzzeichen ohne Folgeklausel.
- **Failure-Szenario:** Ein Leser/Pfleger der §7-Chronik trifft auf einen Satz,
  der mit einem alleinstehenden Semikolon endet, und kann nicht erkennen, ob
  Inhalt verloren ging — genau der „lesbare Chronik"-Anspruch, den der Slice
  setzt. Zum Vergleich behandeln die Nachbarzeilen denselben Fall korrekt:
  `spec/spezifikation.md:2146` behält `… Lastenheft-CR 0.46.0; Begründung in
  der begleitenden ADR |` (Prosa-Form ohne Link, vollständige Klausel), und die
  Zeilen mit satzschließendem Punkt vor der entfernten Klausel enden sauber.
- **Verifizierbar:** nein durch ein Gate — kein Modul flaggt hängende
  Interpunktion (doc-check ist grün). Reproduzierbar per
  `grep -nE '(; )\|$' spec/spezifikation.md` (Zeilen 2147/2149) und
  `git diff 341423a..HEAD -- spec/spezifikation.md` (Vorher-Text).
- **Nicht blockierend**, aber vor Closure trivial zu heilen (Semikolon
  entfernen oder — wie 2146 — auf `; Begründung in der begleitenden ADR`
  angleichen). Betrifft die kanonische Spec.

### F-2 (INFO) — bare Slice-Nummern-Form in `spec/lastenheft.md` §7 verbleibt

- **Kategorie:** INFO (dokumentationswürdige Restspur, außerhalb der
  mechanisierten Regel und der erklärten Slice-Reichweite).
- **Quelle:** v5.0.0-`grundlagen-referenz-richtung` (Spec verweist nie abwärts
  auf Slice); ADR `0047` §Entscheidung 2 (Reichweite: `slice-NNN`-Token +
  echte ADR-Links).
- **Pfad:** `spec/lastenheft.md:2214` (Zeile 0.33.0).
- **Befund:** Die Zelle enthält `… war eine Harness-Ehrlichkeits-Lücke über
  drei Slices (049/052/053)`. Das ist dem Sinn nach ein Abwärts-Slice-Verweis,
  liegt aber nicht in der `slice-NNN`-Token-Gestalt vor (nackte Ziffern hinter
  dem Wort „Slices") und wird darum weder vom Modul matrix erkannt noch von der
  erklärten Putz-Reichweite (Token + Links) erfasst.
- **Failure-Szenario:** Ein späterer Leser deutet `(049/052/053)` als
  Slice-Referenzen und liest §7 als weiterhin abwärts-zeigend. Kein
  Gate-Bruch; die Chronik bleibt gültig.
- **Verifizierbar:** nein (kein Token, kein Link — matrix greift nicht;
  doc-check grün). Bewusst als INFO belassen, da außerhalb des Slice-Scope;
  eine Entscheidung „auch bare Nummern tilgen" wäre ein eigener redaktioneller
  Schritt (ggf. Etappe D / Prozess-Notiz, die der Slice selbst unter §5
  „Zukunft" bereits andenkt).

## Negativbefunde (geprüft, ohne blockierenden Befund)

- **Vollständigkeit des Putzes (Token/Links):** `grep -nE 'slice-[0-9]{3}'`
  über beide Spec-Dateien liefert **null** Treffer (nicht nur in §7);
  `grep -nE 'docs/plan/adr'` liefert in §7 **keinen** Treffer — die restlichen
  Vorkommen liegen sämtlich im Config-/Schema-/AK-Körper (Beispielpfade,
  `trace.adrs.dir`-Defaults), nicht in §7. Kein `../docs/plan/adr/`-Live-Link
  in §7 beider Straten. doc-check (matrix aktiv, §7 nicht mehr exempt) grün.
- **Korrekt behaltenes:** die DC-Selbstlinks (`spezifikation.md#…` in-file,
  `lastenheft.md#…` aufwärts) und die `../harness/conventions.md#mr-…`-Links
  (Ziel = Konventionsdatei, weder ADR- noch Slice-Klasse → matrix-neutral)
  bleiben; keiner ist ein verbotener Abwärtsverweis. Die fiktiven
  Beispiel-Kennungen `ADR-0042`/`ADR-0099` in `spec/lastenheft.md:2244`
  (redaktionell, in Backticks, kein Link, Zeile mit `d-check:ignore`-Marker)
  sind korrekt behalten — matrix trägt kein `adr`-Token (nur die slice-Klasse
  hat ein `token`), und sie sind keine Links; ids exemptet sie via Marker.
- **Verweis-Spalte weg (beide Tabellen):** Header
  `spec/spezifikation.md:2137` = `| Datum | Änderung |`, Separator 2138 =
  `|---|---|`; `spec/lastenheft.md:2193` = `| Version | Datum | Änderung |`,
  Separator 2194 = `|---|---|---|`. Die `Verweis`-Spalte ist aus Kopf-,
  Trenn- und allen Datenzeilen entfernt. Die Sonderzeile
  `spec/lastenheft.md:2244` mit trailing `<!-- d-check:ignore … -->`
  behielt ihre Kommentar-Zelle; die vorige `| — |`-Verweis-Zelle (leer für
  eine redaktionelle Zeile) wurde sauber getilgt. Tabellenstruktur gültig
  (doc-check-Reader akzeptiert die Direktiven-Toleranz-Zeile, 0 Befunde).
- **Prosa-Integrität lastenheft §7:** in `spec/lastenheft.md` §7 gab es keine
  ADR-Prosa-Links (nur Slice-Token in der Verweis-Spalte + „in begleitender
  ADR" als Prosa); entfernt wurde nur die Verweis-Spalte. Keine hängende
  Interpunktion vor Zellenschluss (`grep -nE '(; |, |: |— )\|$'` über §7:
  keine Treffer). Die MR-Links in der Prosa blieben unberührt.
- **ADR-Korrektheit/Immutabilität:** `git diff 341423a..HEAD --
  docs/plan/adr/0022-…md` ist **leer** (byte-identisch); über den Range ist
  keine Accepted-ADR berührt (nur `0047-…md` neu, 83 Zeilen). ADR `0047`
  beschreibt die Entscheidung als Supersede-**Verfeinerung** des Geltungs-
  bereichs (verengt nur die `exclude-sections`-Aussage für die Spec-Straten,
  lässt `0022` unberührt), verweist adr→adr auf `0022`/`0016` (erlaubt, keine
  adr→slice-Regel verletzt), und trägt slice-Token einzig in `## Geschichte`
  (`docs/plan/adr/0047-…md:83`). `make adr-check` grün.
- **Konfig↔Wirkung:** der `.d-check.yml`-Kommentar (matrix-Block) beschreibt
  die Verengung korrekt („Nur die immutable ADR-Geschichte bleibt
  provenance-exempt … Spec-§7 NICHT mehr ausgenommen"); `exclude-sections:
  [Geschichte]` deckt weiter alle ADR-`## Geschichte`-Abschnitte (47 Dateien
  tragen `## Geschichte`, inkl. `0047` selbst — die ADR-interne Zählung „46"
  meint die übrigen). Kein ADR-Bruch (doc-check grün, matrix ohne Befund).
- **Slice-Doc↔Umsetzung:** `docs/plan/planning/in-progress/slice-087-…md`
  (Ziel/Vorgehen/DoD/Abnahme-Punkt 1 „entkoppeln") deckt die Umsetzung ohne
  Widerspruch. Der ADR-Status `Proposed` ist mit dem laufenden Slice konsistent
  (die DoD-Zeile „ADR-0047 Accepted" ist ein Closure-Kriterium, noch nicht
  fällig; die Abnahme der Konformitätsform ist im Slice bereits entschieden).

## Kategorie-Summary

- HIGH: 0
- MEDIUM: 0
- LOW: 1 (F-1 — hängende Semikola in `spec/spezifikation.md:2147`/`:2149`)
- INFO: 1 (F-2 — bare Slice-Nummern-Form in `spec/lastenheft.md:2214`)

## Verdikt

**abnahmereif.** Kein HIGH/MEDIUM. Der Kern-Auftrag ist erfüllt und durch die
selbst gefahrenen Gates belegt: §7 beider Straten ist token- und
abwärtslink-frei (grep + doc-check), die Verweis-Spalte ist in beiden Tabellen
vollständig entfernt, die Tabellenstruktur bleibt gültig, ADR `0022` und alle
Accepted-ADRs sind byte-unberührt (leerer Diff, adr-check grün), und
`exclude-sections: [Geschichte]` schützt die immutable ADR-Geschichte weiter.
Die verbleibende LOW-Nacharbeit (zwei hängende Semikola, F-1) ist redaktionell,
nicht blockierend, aber vor Closure trivial zu heilen — sie berührt die
kanonische Spec und läuft dem eigenen „lesbare Chronik"-Anspruch des Slice
zuwider. F-2 ist bewusst außerhalb des Slice-Scope belassen.
