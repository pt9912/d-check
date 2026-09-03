# Slice slice-191: Den Alt-Bestand archivieren (welle-60…welle-85)

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-87](../welle-87-wellen-archivierung.md).

**Bezug:** [`modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 4](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6)
(die Pflicht, die dieser Slice nachholt);
[slice-190](../done/slice-190-wellen-archiv-werkzeug.md) (liefert das
Werkzeug `tools/archive-wave/`, das dieser Slice anwendet);
[welle-87](../welle-87-wellen-archivierung.md) §3/§4 (Closure-Trigger,
Slice-Platzhalter `slice-<NN-B>`, den dieser Slice einlöst).

**Berührte Spec-Stellen:** — Planning-Infrastruktur, keine `DC-FA-*`-Anforderung.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-03.

---

## 1. Ziel

[slice-190](../done/slice-190-wellen-archiv-werkzeug.md) hat das
Archivierungs-Werkzeug gebaut und an einem Fixture bewiesen — nicht am
echten Bestand. Dieser Slice zieht die Trennung nach: **welle-60 bis
welle-85 archivieren**, plus die Zuordnung wellenloser Alt-Slices klären,
die [Modul 6](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6)
bei der Nachrüstung verlangt.

**Zwei Bestands-Eigenheiten wurden bereits beim Planen gemessen, nicht
erst beim Ausführen entdeckt** — genau die Vorsicht, die
[slice-190 §5](../done/slice-190-wellen-archiv-werkzeug.md) (Risiko 3,
`BEO-011` sechste Instanz) benannt hat:

1. **welle-60 bis welle-66 haben keinen Welle-Plan.** `ls done/welle-6*`
   zeigt nur `welle-60-results.md` … `welle-66-results.md` — keine
   `welle-6N-<titel>.md`. Diese sechs Wellen wurden **vor** der
   Plan-Datei-Konvention geschlossen und tragen laut
   [slice-088](../done/slice-088-etappe-d1-doc-form.md) Zeile 131 eine
   „retroaktiv markierte, minimale Ergebnis-Notiz" als bewusste
   Nutzer-Entscheidung nach — nicht eine vollständige Plan-Rekonstruktion.
   `tools/archive-wave`s `FindWellePlan` verlangt **genau eine** Treffer-Datei
   und schlägt für diese sechs fehl.
2. **Die vier vermeintlich wellenlosen Alt-Slices sind keine.** Gemessen:
   `done/slice-*.md` ohne `**Welle:**`-Feld sind genau vier —
   [slice-077](../done/welle-60/slice-077-stiller-tabellen-uebersprung.md),
   [slice-078](../done/welle-61/slice-078-ignore-refs-quell-skopus.md),
   [slice-079](../done/slice-079-zitat-verifikation.md),
   [slice-103](../done/slice-103-geteilte-lexik-raender.md). Alle vier
   nennen ihre Welle **in Prosa** (Status-Zeile bzw. Fließtext): `slice-077`
   → welle-60, `slice-078` → welle-61 (**Status**-Zeile), `slice-079` →
   welle-62, `slice-103` → welle-74 (bestätigt durch `welle-74`s eigene
   Slice-Tabelle, die `slice-103` führt). Es gibt keinen echten
   Zuordnungs-Konflikt — nur ein fehlendes Feld, aus einer Zeit vor der
   `**Welle:**`-Feld-Konvention.

## 2. Vorgehen

1. **`**Welle:**`-Feld für die vier betroffenen Slices nachtragen**, nach
   dem oben gemessenen Beleg (Status-Zeile/Prosa des jeweiligen Slice bzw.
   der Ziel-Welle eigene Slice-Tabelle). Das ist **kein** Rückbau des
   historischen `**Status:**`-Felds (das bleibt unverändert, `AGENTS.md`
   §3.7 Bestands-Ausnahme) — es ergänzt ein bislang fehlendes Feld aus
   bereits im selben Dokument stehenden Fakten. Nach diesem Schritt gibt es
   **keine** wellenlosen Alt-Slices mehr vor der Einführung der Regel — die
   welle-87-§3-Bedingung „Zuordnung entschieden" ist damit für diesen
   gemessenen Bestand erfüllt, ohne Sammel-Archiv oder chronologische
   Näherung zu brauchen.
2. **`tools/archive-wave` um den Plan-losen Fall erweitern.** `FindWellePlan`
   liefert für welle-60…66 **null** Treffer statt eines Fehlers; `Apply`
   überspringt in diesem Fall den Welle-Plan-Eintrag im ZIP und die
   Welle-Stub-Erzeugung (nichts zu ersetzen, wo kein Volltext liegt) und
   archiviert nur Slices + Reviews. Test: eine Fixture-Variante ohne
   Welle-Plan-Datei, die genau dieses Verhalten belegt.
3. **Anwenden, Welle für Welle, mit Kontrolle zwischen jedem Schritt:**
   `make archive-wave WELLE=welle-NN` (Dry-Run lesen) → `APPLY=1` (schreiben)
   → `git status` prüfen → `make gates` → Commit, bevor die nächste Welle
   drankommt. Nicht alle 26 Wellen in einem Rutsch, damit ein Fehlschlag
   früh sichtbar wird und nicht 26 Wellen mitreißt.
4. **`make gates` und `make fullbuild`** auf dem vollständig archivierten
   Bestand (welle-87 §3 verlangt beides explizit).
5. Closure-Notiz; bei dieser Slice-Closure zusätzlich: **welle-87 selbst
   schließen** (sie hat keine weiteren offenen Slices nach diesem) — die
   sechs Closure-Schritte aus
   [Modul 6](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6)
   laufen für welle-87 **nach** diesem Slice, nicht in ihm (Rollen-Trennung
   Planner/Architect, Modul 8) — dieser Slice liefert nur den Trigger-Beleg.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Umbau der Fixture-Tests aus slice-190** über den einen neuen
  Plan-losen Fall hinaus — die bestehende Testsuite bleibt Referenz.
- **`welle-86` wird nicht eingesammelt** (welle-87 §5) — sie bleibt
  eigenständig und archiviert sich bei ihrer eigenen Closure.
- **Kein Sammel-Archiv für vor-Einführungs-Bestand** — §1 misst, dass diese
  Ausweich-Option nicht gebraucht wird; alle vier Kandidaten lösen sich per
  Feld-Nachtrag in eine bestehende Welle auf.

## 4. Definition of Done

- [ ] `**Welle:**`-Feld für slice-077/078/079/103 nachgetragen, mit
      Beleg-Zeile in diesem Slice-Plan (§1) belegt.
- [ ] `tools/archive-wave` behandelt eine Welle ohne Plan-Datei
      (welle-60…66) korrekt — Test vorhanden, `make archive-wave-test` grün.
- [ ] welle-60 bis welle-85 sind archiviert: `done/<welle-id>/archiv.zip`
      existiert je Welle, Stubs an der Zielform, Review-Reports ohne Stub,
      keine gebrochenen Repo-Verweise (`make doc-check` grün nach jeder
      angewendeten Welle).
- [ ] `make gates` **und** `make fullbuild` grün auf dem archivierten
      Bestand (Exit explizit).
- [ ] **Unabhängiger Review und unabhängige Verifikation.**
- [ ] Closure-Notiz mit Lerneintrag; jedes Risiko aus §5 mit Ausgang; die
      drei Modul-6-Paarungen geprüft (bei der Closure von welle-87, die
      dieser Slice schließt).

## 5. Abnahme-Punkte / Risiken

- **26 Wellen in Folge zu archivieren ist viel Wiederholung — eine falsche
  Annahme im Werkzeug wirkt 26-fach, nicht einmal.** Mitigiert durch
  Schrittweite (§2 Punkt 3: Commit je Welle, nicht ein Sammel-Commit) —
  bricht es bei welle-70, sind 60–69 bereits sicher committet. —
  **Ausgang:** *(bei Closure)*
- **Der Feld-Nachtrag (§2 Punkt 1) ändert historische `done/`-Dateien.**
  `AGENTS.md` §3.7 schützt das `**Status:**`-Feld dieser Slices ausdrücklich
  als eingefrorenen Lauf-Beleg — ein `**Welle:**`-Nachtrag ist ein anderes
  Feld, aber die Grenze („was darf an einem historischen Slice noch
  ergänzt werden") ist an diesem Bestand vorher nicht geprüft worden. —
  **Ausgang:** *(bei Closure)*
- **Der Verweis-Nachzug über 26 Wellen kann eine Datei mehrfach treffen**
  (ein Dokument, das auf Slices aus zwei verschiedenen Wellen verweist).
  `RewriteRepo` verarbeitet das move-für-move korrekt, aber die
  Reihenfolge der 26 Anwendungen ist neu für das Werkzeug — bei slice-190
  nur an einer einzelnen Welle geprüft. — **Ausgang:** *(bei Closure)*
- **`welle-87`s eigene Closure hängt an diesem Slice** — verzögert sich
  dieser Slice, verzögert sich die ganze Welle, und `welle-86` bekommt bei
  ihrer eigenen Closure weiterhin kein bewiesenes Werkzeug gegen einen
  Mehrwellen-Lauf (nur gegen das Fixture aus slice-190). —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): [slice-190](../done/slice-190-wellen-archiv-werkzeug.md)
liegt in `done/`, WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls sich beim Ausführen zeigt,
dass die Plan-lose-Welle-Erweiterung (Punkt 2) und die Massen-Anwendung
(Punkt 3) getrennte Slices sein sollten, weil die Erweiterung selbst
unerwartet groß wird.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `tools/archive-wave` (Erweiterung des in slice-190
  gebauten Werkzeugs) und `docs/plan/planning/` (der archivierte Bestand
  selbst). Beide Default `*` = **Greenfield**
  ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration) — Erweiterung einer erst in slice-190 entstandenen
  Infrastruktur, kein gewachsener Bestand mit eigener Konventions-Historie.

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-09-03, höchste
  Kennung `BEO-027`): [`BEO-011`](../observations.md) (Zähler jetzt 6,
  sechste Instanz durch slice-190 selbst registriert, Ausgang an dieser
  Closure fällig) — genau der Fall, den dieser Slice auflöst: wurde die
  Fixture-Beweisführung von slice-190 durch echte Bestands-Eigenheiten
  widerlegt (Plan-lose Wellen, unter-getaggte Slices), oder hat sie
  getragen? Beide oben in §1 gemessenen Eigenheiten sind Belege für
  „getragen, mit einer vorhergesehenen und einer neuen Eigenheit" — kein
  unvorhergesehener Bruch. Keine weiteren Treffer.

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `image-scan.yml`
  **grün** (jüngster Lauf 2026-09-02T07:56:37Z). `upstream-drift.yml`
  **ROT** — jüngster Lauf 2026-09-02T05:19:44Z, derselbe planmäßige Fund wie
  bei [slice-190](../done/slice-190-wellen-archiv-werkzeug.md) §7 (Go
  1.27.0→1.27.1, semgrep 1.175.0→1.176.0), kein Zitat-Bruch, keine
  Regression. Ohne Konsequenz für diesen Slice.

Slice-ID: slice-191. Betroffene IDs: keine `DC-FA-*`. Module: keins.
Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Erweiterung einer Infrastruktur ohne
eigene Konventions-Historie, kein Reconciliation-Aufwand.

## 9. Closure-Notiz (nach `done/`)
