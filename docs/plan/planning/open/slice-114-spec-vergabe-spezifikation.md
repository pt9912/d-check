# Slice slice-114: `SPEC-*`-Vergabe in der Spezifikation + `structure`-Wächter für §2

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-80-struktur-ids](../welle-80-struktur-ids.md) (zugeordnet bei
der Eröffnung).

**Bezug:**
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in)
(die §2-Invariante läuft als `structure`-Regel),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)
(jede Nennung außerhalb der Zieldatei ist linkpflichtig — seit slice-113),
Baseline-Template `spezifikation.template.md` §2 (Überschrift-Kennung
`### SPEC-NNN — …`; Tabellen-Kennung je Zeile in §3/§4/§6), Entscheide D1
(Vollvergabe) und D2 (eine Kennung je §4-Zeile).

**Berührte Spec-Stellen:** `spezifikation.md` §2 (fünf kennungslose
Sektionen), §3 Defaults, §4 Grund-Codes, §6 externe Verträge, §7 Historie — der
Verweis zeigt aufwärts, die Spec nennt diesen Slice nicht.

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

`spec/spezifikation.md` trägt `SPEC-<NNN>` fortlaufend je Datei: §2 — jede
kennungslose `###`-Sektion (heute fünf: Befund, JSON-Ausgabe, JSON-Diagnose,
YAML-Ausgabe, `.d-check.yml`) wird zu `### SPEC-NNN — <Titel>`; §3 Defaults,
§4 Grund-Codes und §6 externe Verträge bekommen eine **erste Spalte**
`Kennung` je Zeile (die `.a`-Verfeinerungen in §1 tragen weiter ihre
Anforderungs-Kennung — Baseline: „Ein Abschnitt, der eine einzelne Anforderung
verfeinert, trägt die Verfeinerung; alles Übrige trägt `SPEC-*`"). Im **selben
Commit** wird der `structure`-Wächter scharf: „jede `###`-Überschrift unter
§2 trägt eine Kennung" — vorher am Bestand gemessen **rot** (fünf Befunde),
danach null. Zahlen werden gemessen, nicht übernommen.

## 2. Vorgehen

1. **Inventur mit dem Produkt/`grep`:** `###`-Sektionen je `##`-Teil, Zeilen
   je Tabelle in §3/§4/§6, vorhandene Verweise auf die §2-Slugs (Anker ändern
   sich mit der Überschrift! — `anchors` misst es, jeder Verweis wird
   retargetet, der alte Slug ggf. als `<a id>` nur dann, wenn ein
   **eingefrorener** Verweis (Accepted-ADR) darauf zeigt — dann ist das
   Zeilenanker-Verbot der Baseline gegen die ADR-Immutabilität abzuwägen und
   der Ausgang in der Closure-Notiz festgehalten).
2. **Vergabe §2:** fünf Überschriften `### SPEC-001 — Befund` … fortlaufend;
   Verweise im Baum retargeten (Handbuch, ADRs **nur** Proposed, Slices
   `done/` über `ignore-refs`-Ventil, falls nötig).
3. **Vergabe §3/§4/§6:** neue Spalte `Kennung` links, `SPEC-NNN` je Zeile,
   Zählung setzt §2 fort; §5 (Metriken) bleibt leer.
4. **`structure`-Regel** in `.d-check.yml`: `files: spec/spezifikation.md`,
   `section: '## 2. Datenstrukturen und Schemas'`, `forbid-pattern` auf
   `###`-Zeilen, die weder `SPEC-` noch eine Anforderungs-Verfeinerung tragen
   (RE2 ohne Lookahead — Zeichenklassen-Form; vorher gemessen: rot am
   Bestand = 5, grün nach Vergabe = 0). Ist die Invariante mit den heutigen
   Schlüsseln nicht präzise ausdrückbar, wird das als Grenze benannt und der
   neue Schlüssel als CR-Kandidat notiert — **kein** Produkt-Code in dieser
   Welle.
5. **[ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-Kandidat messen:** zeigt `Schärft:` „§Kern" auf eine Zeile ohne
   Anker? Befund in die Closure-Notiz (Beleg für den Auflösungs-Trigger).
6. **Spiegel** ([`MR-025`](../../../../harness/conventions.md#mr-025),
   `grep` nach den alten §2-Slugs und nach „SPEC"): Spezifikation §7-Historie,
   Handbuch (Verweise auf §2-Anker), `harness/README.md`, ADR-Index-Konvention
   (Beispiel), `.d-check.yml`-Kommentar.
7. Unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine `ARC-*`** (slice-115), **keine ADR-Nachzüge** (slice-116).
- **Keine Änderung am Lastenheft** (Struktur-IDs sind Form, nicht
  Anforderung).
- **Kein Produkt-Code** — auch kein neuer `structure`-Schlüssel.

## 4. Definition of Done

- [ ] §2/§3/§4/§6 tragen `SPEC-NNN` fortlaufend; Zählung in der Closure-Notiz
      (gemessen).
- [ ] `structure`-Regel scharf im Vergabe-Commit; Messung rot-vorher /
      grün-nachher dokumentiert.
- [ ] Alle Verweise auf geänderte §2-Slugs retargetet; `make doc-check` grün.
- [ ] [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-„§Kern"-Messung in der Closure-Notiz.
- [ ] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Slug-Wechsel der §2-Überschriften** bricht eingefrorene Verweise (Accepted-
  ADRs, `done/`-Slices) — Ventil `ignore-refs` bzw. Abwägung `<a id>` gegen das
  Baseline-Verbot. — **Ausgang:** *(bei Closure)*
- **`forbid-pattern`-Form** ist ohne Lookahead grob; Falsch-Positiv bei einer
  Überschrift mit anderem Anfangsbuchstaben. — **Ausgang:** *(bei Closure)*
- **Tabellen-Spalte vorn** verschiebt `table-order`/`table-column`-Regeln?
  (§7-Historie ist nicht betroffen; §3/§4/§6 tragen keine Chronologie-Regel —
  prüfen.) — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): slice-113 in `done/` (die Muster und die
Konvention stehen).

**Rückführungen:** `in-progress` → `next`, falls die Slug-Wechsel mehr
eingefrorene Verweise treffen, als `ignore-refs` sauber trägt (dann Entscheid
über Zeilenanker vor der Vergabe).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Spec-Stratum Technik (`spec/spezifikation.md`, GF),
  Prüf-Profil (GF), Handbuch-Verweise (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): BEO-002 ([`MR-025`](../../../../harness/conventions.md#mr-025)
  — Spiegel per `grep` nach dem alten Wortlaut/Slug), BEO-006/009
  Arbeitsregeln; BEO-008 nicht einschlägig (kein Pin).

Slice-ID: slice-114. Betroffene IDs:
[`DC-FA-STRUCT-001`](../../../../spec/lastenheft.md#dc-fa-struct-001--struktur-invarianten-innerhalb-eines-dokuments-modul-structure-opt-in),
[`DC-FA-ID-001`](../../../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids).
Module: Spezifikation, Prüf-Profil (`structure`), Handbuch (Verweise). Gates:
`make doc-check` (eng), `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Vergabe im eigenen Spec-Stratum nach
Baseline-Form; kein Legacy-Import.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
