# Welle 78 — Baseline-Regelwerk v5.0.0 → v5.6.0 — Closure-Notiz

**Welle:** welle-78-baseline-v560-migration
**Abschluss:** 2026-08-21
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

Die Baseline ist **vollständig** von v5.0.0 auf **v5.6.0** gehoben — vier
Slices, je mit unabhängigem Frischkontext-Review (zwei davon zusätzlich mit
bestätigender Re-Review):

- **Etappe A** ([slice-106](slice-106-baseline-v560-vendoring.md)): Bundle
  v5.6.0 committet vendored (`--verify` 51 Dateien), v5.0.0-Baum entfernt,
  Pin samt [`MR-026`](../../../../harness/conventions.md#mr-026), lebende
  Verweise pin-gebunden retargetet, fünf eingefrorene Fundstellen
  getombstoned; dabei zwei Defekte im `--check-latest`-Content-Audit gefixt
  (falscher Drift-Alarm, Müll-Verzeichnisse).
- **Etappe B** ([slice-107](slice-107-baseline-v560-delta-audit.md)): der
  Stufen-Audit v5.1.0–v5.6.0 — je Regel eine Antwort (konform · anzupassen ·
  n. a. mit Begründung), sechs „anzupassen"-Findings zu zwei
  Etappe-C-Slices gebündelt.
- **Etappe C-1** ([slice-108](slice-108-roadmap-offene-wellen.md)): die
  Roadmap in der v5.6.0-Form **§Offene Wellen** (derivativ, Ruhe-Marker
  „Nichts in Arbeit" mit Wächter), beide Prüf-Profile per Config,
  [`MR-024`](../../../../harness/conventions.md#mr-024) aufgelöst.
- **Etappe C-2…C-6** ([slice-109](slice-109-v560-konventions-nachzuege.md)):
  Vergabe-Deklaration, [`MR-027`](../../../../harness/conventions.md#mr-027)
  (Struktur-ID-Verzicht als deklarierte Abweichung), Kommentar-Regel-Träger
  (AGENTS §3.7 + `reviewer.md` 1.5.0) samt Räumung von zwanzig
  Token-Fundstellen, 27 Kennungs-Anker im MR-Index, Leseordnung,
  Bestands-Stichprobe `modul-14` (null Widerspruchs-Funde).

**Kein Release:** die Hebung ist Harness, nicht Produkt; Etappe C änderte
Befunds-**Klartexte** (nicht stabilitätszugesagt), keinen Grund-Code, keinen
Exit-Code, keinen Befundsatz — die Strings reisen mit dem nächsten Release.
Das ist die ausdrückliche Lesart des Wellen-Triggers „Release nur, falls
Etappe C Produkt-Verhalten ändert", ehrlich benannt statt still entschieden.

## Was hat funktioniert?

- **Der Widerspruchs-Ausgang der Baseline hat an der eigenen Migration
  gewirkt:** der Struktur-ID-Verzicht wäre im ersten Audit-Entwurf als
  „zulässiger Rückfallweg" durchgerutscht — der Review hat ihn als echte
  Abweichung eingeordnet, und die Antwort war genau die Form, die die
  v5.4.0-Stufe vorschreibt: deklarieren
  ([`MR-027`](../../../../harness/conventions.md#mr-027)) statt still die
  [`MR-000`](../../../../harness/conventions.md#mr-000--baseline-aussage)-Aussage brechen.
- **d-check war der Baseline mehrfach voraus:** Kommentar-Regel,
  Gate-Obermenge-Nachweis, TA-7/Hauptzweig-vor-Arbeit — dieselben Lehren
  sind hier zuerst eingetreten und jetzt Baseline-Default. Der Audit konnte
  sie mit Fundstellen belegen statt behaupten.
- **Die adversarialen Reviews trugen die Welle:** vier Erst-Reviews (zwei
  blockierend), zwei bestätigende Re-Reviews; die schwersten Funde waren
  Meta-Funde — Commit-Schnitt, Audit-Vollständigkeit, Räumungs-Scope — und
  jeder wurde vor dem Push geheilt.

## Was ging anders als geplant?

- **Zwei neue Beobachtungs-Klassen am eigenen Arbeiten**, beide je zweimal
  eingetreten: pfad-selektive Commits tragen still den gestagten Index
  (**BEO-006**, Zähler 2 — das zweite Mal Stunden nach der frischen
  Arbeitsregel), und eine Pipe verschluckt den Gate-Exit (**BEO-007**,
  Zähler 2 — einmal ging ein roter doc-check-Stand bis zum Push durch).
  Beide stehen im Register mit gelebtem Gegenmittel und benannter
  mechanischer Form für 3×.
- **Die Räumung einer frisch installierten Regel folgte der Gewohnheit statt
  dem Geltungsbereich** (16 Go-Fundstellen geräumt, vier in Skripten/Config
  übersehen, Botschaft „alle") — vom Review gefangen, per Neuschnitt geheilt,
  Bestandsgrenze (Test-Kommentare) jetzt deklariert statt still.
- **Sonst nichts:** geplantes Ende 2026-08-23, geschlossen am 2026-08-21;
  der Etappen-Schnitt (A→B→C) trug unverändert durch.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 4 — was aus dieser Welle in die Steuerung
zurückfließt.

- **Wer eine Regel scharfschaltet, läuft ihren Geltungsbereich ab — nicht
  seine Gewohnheits-Teilmenge.** Die Kommentar-Räumung fand alle
  Go-Fundstellen und übersah Skripte und CI-Config, die die Regel wörtlich
  nennt. Der Geltungsbereich steht in der Regel; der grep gehört über ihn,
  nicht über das zuletzt bearbeitete Verzeichnis.
- **Ein Audit-Beleg braucht getrennte Entstehung:** „§4-Zeile vergab die
  Nummer vor der Datei" war nicht belegbar, weil beide im selben Commit
  entstanden. Wer einen Ablauf als Beleg anführen will, braucht Artefakte
  mit eigener Historie.

## Beobachtungs-Register (Zeiger)

Der Lese-Schritt dieser Welle: **nichts hat 3× erreicht** — BEO-006 steht
nach dem zweiten Eintreten bei 2, BEO-007 ist neu bei 2 (beide mit gelebtem
Gegenmittel und benannter Verkörperungs-Form); BEO-002/003/004 bleiben
verkörpert, BEO-001/005 gestrichen. Die Verkörperungen haben auch diese
Welle getragen: [`MR-025`](../../../../harness/conventions.md#mr-025) erzwang
die Spiegel-Listen (C-1: fünf Spiegel; die Review-Funde F-1/F-3 waren
Spiegel-Lücken), der BEO-004-Reviewer-Anker blieb ohne neuen Fund.

## Folge-Slices

- **Keiner.** `open/` ist leer, §Nächste Wellen trägt `— keine —`. Die
  benannten Grenzen (Produkt-Default `## Aktuelle Welle`, der
  `wave-drift`-Zustand „Welle offen ohne Anspruch", die
  Test-Kommentar-Altbestände, der [`MR-027`](../../../../harness/conventions.md#mr-027)-Auflösungs-Trigger) haben je einen
  beobachtbaren Trigger und warten dort — keine ist slice-reif.
- **Upstream-Notiz an den Kurs** (kein d-check-Slice): die
  5-vs-6-Finding-Feld-Drift (modul-10 §Output-Schema gegen das
  Report-Template) besteht in v5.6.0 fort.

## Verifikation

- **Closure-Trigger erfüllt:** alle vier Slices in `done/`; Pin v5.6.0 mit
  `--verify` (51 Dateien) und `--check-latest` beidseitig OK; kein lebender
  Verweis nennt `baseline/v5.0.0` (eingefrorene getombstoned, `make
  doc-check` belegt beides); der Stufen-Audit trägt je Regel eine Antwort;
  kein Release (Lesart oben ehrlich benannt); `make fullbuild` grün,
  Image-Hash `sha256:dbe5c4a0…663b`.
- **Trigger-Audit** über die drei Artefaktklassen: keine offenen Carveouts,
  keine stehengebliebene Gate-Reifestufe, kein eingetretener
  Re-Evaluierungs-Trigger — geprüft auch für die neuen
  [`MR-026`](../../../../harness/conventions.md#mr-026)/[`MR-027`](../../../../harness/conventions.md#mr-027)-Trigger
  und [ADR-0057](../../adr/0057-structure-tabellen-monotonie.md) (keiner
  fällig).
- **Sechs Review-Runden über vier Slices** (vier Erst-Reviews, davon zwei
  blockierend; zwei bestätigende Re-Reviews), alle Auflagen vor dem Push
  eingearbeitet.
