# Slice slice-062: Handbuch-E2E-Beispiel-Verankerung (Kommando/Ausgabe ↔ echtes Binary)

**Status:** next (Backlog; ausgegliedert aus slice-061 per Nutzer-Entscheid
2026-07-04 „A jetzt, B als slice-062").

**Welle:** noch nicht terminiert (Folge von welle-50).

**Bezug:** Verifikations-Mechanik gegen bestehende Verträge
([`DC-FA-CLI-003`](../../../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)
Exit-Codes;
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
Image-Verhalten). **Kein Change Request** (kein neuer Vertrag), **kein ADR**
(E2E-Test-Erweiterung im bestehenden Schnitt). Schwester-Slice zu
[`slice-061`](../done/slice-061-doc-config-beispiel-verifikation.md)
(dort Dimension A: Config-Fragmente ↔ Parse; hier Dimension B: Kommando-/
Ausgabe-Beispiele ↔ Verhalten).

**Autor:** pt9912. **Datum:** 2026-07-04.

---

## 1. Ziel

Das Benutzerhandbuch führt **Kommando-Aufrufe mit dokumentierter Ausgabe/
Exit-Code** (` ```bash ` + ` ```text `) — sauberes Repo ⇒ Exit 0 +
„0 Befund(e)"; kaputter Link ⇒ Befund-Zeile `Datei:Zeile  Ziel  Grund` +
Exit 1; `--doctor` ⇒ gruppierte Diagnose. Diese Verhaltensbehauptungen werden
nie gegen das echte Binary geprüft; driftet das CLI-Verhalten, bleibt das
Handbuch still falsch. **Neu:** repräsentative Kommando-Beispiele werden als
E2E-Fälle gegen das reale Binary/Image über Fixtures verankert (Anregung des
Auftraggebers: „von E2E-Tests könnte ein Handbuch profitieren, da ja im
Handbuch Beispiele aufgeführt sind").

## 2. Entscheidungen (Entwurf — bei Aufnahme schärfen)

- **Vorhandene E2E-Infrastruktur nutzen** statt neuen Runner bauen:
  `tools/image-test.sh` ([`make image-test`](../../../../harness/README.md#sensors-feedback-gates))
  fährt schon nativ == Container mit Fixtures + Exit-Code-/Ausgabe-Prüfung;
  `cli_acceptance_test.go` fährt CLI-E2E gegen `MemFS`/Temp-Repos.
- **Nur Beispiele mit prüfbarer Verhaltensbehauptung** verankern: je ein
  Fixture, das die Prämisse herstellt (sauberes Repo / kaputter Link /
  `--doctor`-Fall), der dokumentierte Aufruf (die **Flags**, nicht der
  wörtliche `docker run …@sha256`-Pull), Assertion auf Exit-Code + Ausgabe-
  **Form** (Befund-Zeilen-Schema, „N Befund(e)"-Zeile, Diagnose-Kopf).
- **Auf Form prüfen, nicht auf wörtliche Zeilen:** Regex/Schema statt
  Datei-Zahlen/Versions-Pins — sonst bricht der Test bei jeder harmlosen
  Doku-/Release-Änderung (Wartungsfalle statt Wert).
- **Nicht-replaybare Beispiele markieren** (externer Zustand, konkrete
  Digests, Netz) — gleicher Opt-out-Marker-Ansatz wie slice-061, mit Grund;
  kein stiller Ausschluss (welche Beispiele E2E-verankert sind, welche nicht
  und warum).
- **Determinismus/netzlos:** lokal gebautes Image (wie `image-test`), kein
  Pull.

## 3. Definition of Done (Entwurf)

- [ ] repräsentative Handbuch-Kommando-Beispiele mit Verhaltensbehauptung als
  E2E-Fälle (Fixture + Flags + Exit-Code-/Ausgabe-Form-Assertion), eingereiht
  in `cli_acceptance_test.go` bzw. `tools/image-test.sh`.
- [ ] Beispiel-Auswahl dokumentiert (E2E-verankert vs. begründet ausgenommen).
- [ ] Fail-closed-Guards mutations-verifiziert; adversariale Probe
  (dokumentierte Ausgabe künstlich verletzt ⇒ rot).
- [ ] `make gates`/`make ci` grün; unabhängiges Review; Closure-Move + Body.
  **Kein Produkt-Code, kein Release** (Test-Infra).

## 4. Risiken / offene Punkte

- **Ausgabe-Matching-Stabilität** (s. §2): Form statt Wortlaut — der
  Kern-Entscheid, damit der Test kein Wartungsklotz wird.
- **Beispiel-Auswahl:** nicht jedes der ~15 Kommando-Beispiele lohnt einen
  E2E-Fall; die mit einer eindeutigen, stabilen Verhaltensbehauptung zuerst.
- **Abgrenzung zu bestehenden Tests:** `cli_acceptance_test.go`/`image-test`
  decken viele Verhalten schon ab — Dimension B ergänzt gezielt die
  **im Handbuch gezeigten** Fälle, ohne Bestehendes zu duplizieren.

## 5. Trigger

Ausgegliedert aus [`slice-061`](../done/slice-061-doc-config-beispiel-verifikation.md)
per Nutzer-Entscheid 2026-07-04 (Stufung „A jetzt / B als slice-062"); die
E2E-Anregung selbst kam vom Auftraggeber in derselben Sitzung. Aufnahme in eine
Welle, wenn slice-061 (Dimension A) abgeschlossen ist.

## 6. Sub-Area-Modus-Begründung

GF (E2E-Test-Erweiterung im bestehenden Schnitt). Kein neuer Adapter, keine
BF-Sub-Area.

## 7. Closure-Notiz (nach `done/`)

*(bei Closure zu füllen.)*
