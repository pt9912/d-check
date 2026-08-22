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
2. **Vergabe §2:** fünf Überschriften `### SPEC-NNN — Befund` … fortlaufend;
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

- [x] §2/§3/§4/§6 tragen `SPEC-NNN` fortlaufend; Zählung in der Closure-Notiz
      (gemessen).
- [x] `structure`-Regel scharf im Vergabe-Commit; Messung rot-vorher /
      grün-nachher dokumentiert.
- [x] Alle Verweise auf geänderte §2-Slugs retargetet; `make doc-check` grün.
- [x] [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)-„§Kern"-Messung in der Closure-Notiz.
- [x] `make gates` grün; unabhängiger Review; Closure-Notiz; Register
      gesichtet.

## 5. Abnahme-Punkte / Risiken

- **Slug-Wechsel der §2-Überschriften** bricht eingefrorene Verweise (Accepted-
  ADRs, `done/`-Slices) — Ventil `ignore-refs` bzw. Abwägung `<a id>` gegen das
  Baseline-Verbot. — **Ausgang:** entfallen — **kein** eingefrorener
  ADR-Verweis zeigte auf einen §2-Anker (alle ADR-Anker sind
  `.a`-Verfeinerungen); die Zeilenanker-Frage stellte sich nicht. Alle zwölf
  Verweise sind retargetet — bis auf den Review-Report, dessen Retarget der
  Review als Auflage zurückgenommen hat.
- **`forbid-pattern`-Form** ist ohne Lookahead grob; Falsch-Positiv bei einer
  Überschrift mit anderem Anfangsbuchstaben. — **Ausgang:** **eingetreten,
  anders als gedacht.** Das Risiko war das Falsch-Positiv; der Defekt war das
  Falsch-**Negativ**: mein Muster verlangte Zeilenanfang und genau ein
  Leerzeichen, die Heading-Lexik des Moduls trimmt aber führenden Weissraum
  und nimmt Tab als Trenner — eine eingerückte Sektion entkam still (Review
  F-1). Geheilt und mit dem Produkt belegt; zwei echte Falsch-Positive fing
  die eigene Logik-Probe (nacktes `###`, `### SPEC-123` ohne Titel). <!-- d-check:ignore -->
- **Tabellen-Spalte vorn** verschiebt `table-order`/`table-column`-Regeln?
  (§7-Historie ist nicht betroffen; §3/§4/§6 tragen keine Chronologie-Regel —
  prüfen.) — **Ausgang:** entfallen — gemessen: keine `table-order`/
  `table-column`-Regel liegt auf §3/§4/§6, die Spaltenlage der sechs
  Chronologie-Regeln ist unberührt.

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

**Geliefert:** 66 Kennungen, fortlaufend je Datei und in Dokumentreihenfolge —
§2 fünf Schema-Überschriften (001–005), §3 sieben Defaults (006–012), §4 51
Grund-Codes (013–063), §6 drei externe Verträge (064–066), letztere drei je in
einer neuen ersten Spalte `Kennung`. Zwei Konsumenten halten die Vergabe
lebendig, beide im Vergabe-Commit scharf: die `structure`-Regel auf §2
(Inner Loop) und der bestehende Lockstep-Test, der jetzt Spalte 2 liest und
Spalte 1 auf eine eindeutige Kennung prüft — er war für genau diese Anpassung
fail-closed gebaut und hat sie eingefordert. Die §2-Anker wandern mit den
Überschriften; elf Verweise sind retargetet.

**Review** ([Report](../../../reviews/2026-08-22-slice-114-spec-vergabe-review.md)):
APPROVE mit Auflagen — 0 HIGH, 2 MEDIUM, 2 LOW, 2 INFO, vierzehn
Negativ-Proben. Alle sechs eingearbeitet.

**Was ging anders als geplant — zweimal dieselbe Wurzel:** Ein Wächter, der
eine Lexik prüft, muss *dieselbe* Lexik sprechen wie das Modul, das er
benutzt. Mein Muster verlangte Zeilenanfang und genau ein Leerzeichen;
`parseATXHeading` trimmt beliebigen führenden Weissraum und akzeptiert Tab —
also entkam eine eingerückte Sektion still, und der Wächter behauptete eine
Deckung, die er nicht hatte. Gefunden hat das der Review, nicht ich; meine
eigene Logik-Probe fand danach noch zwei Falsch-Positive derselben Art
(ein nacktes `###` ist keine Überschrift; `### SPEC-123` ohne Titel trägt <!-- d-check:ignore -->
sehr wohl eine Kennung). Die zweite Lehre betrifft das Retargeten: ein
Anker-Rename zieht Verweise nach — aber ein Lauf-Beleg, der ein
Überschrift→Slug-**Paar zitiert**, dokumentiert den Stand von damals; ihn
mitzuziehen macht ihn falsch. Der Review hat den Retarget des Reports als
Auflage zurückgenommen.

- **Steering-Loop-Eintrag:** Sensor ergänzt: der Schema-Abschnitt der
  Spezifikation ist an die Kennungs-Pflicht gebunden — liegt in
  `.d-check.yml §structure`. Kein Auslöser aus dem Register.
- **Beobachtungs-Register (`../observations.md`):** keine neue Beobachtung —
  die Wächter-Lexik-Divergenz ist die bekannte Klasse BEO-003 („eine geteilte
  Lexik driftet an den Rändern, weil jeder Konsument sie selbst vorbereitet"),
  verkörpert in welle-74; hier trat sie erstmals **konfigurationsseitig** auf
  (ein Muster in der Config statt ein Prädikat im Code). Zitiert statt neu
  formuliert; der Zähler bleibt bei 3.
- **Folge-Slices:** [slice-115](../done/slice-115-arc-vergabe-architektur.md)
  (geschlossen) und [slice-116](../open/slice-116-adr-neuzugangs-regel.md)
  (wartet auf 115).
- **Risiken aus §6:** alle mit Ausgang (§5) — eines eingetreten (anders als
  gedacht), zwei entfallen.
- **Drei Paarungen:** Wellen-Slice — die Paarungen prüft die Welle-Closure.

**Messung für den Auflösungs-Trigger der abgelösten Adaption
(Plan-Schritt 5):** [ADR-0012](../../adr/0012-kern-paketschnitt-model-rules-app.md)
nennt als `Schärft:`-Ziel „§Kern" — und „Kern" ist in der Architektur-Sicht
eine **Tabellenzeile ohne Anker**, keine Überschrift. Die ADR zeigt damit auf
etwas, das nicht eindeutig adressierbar ist; sie bleibt immutabel, die Zeile
bekommt ihre Kennung mit slice-115. Der Trigger ist damit gemessen, nicht
behauptet.
