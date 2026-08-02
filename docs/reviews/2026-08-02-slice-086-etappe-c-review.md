# Review — slice-086 Etappe C (MR-Bereinigung + Datei-Migration)

- **Review-Art:** unabhängiges Frischkontext-Review (adversarial), Gegenstand
  Etappe C der Baseline-Migration `v1.4.0` → `v5.0.0` (Umbau
  `harness/conventions.md` auf Index + Datei-je-MR).
- **Skill:** `.harness/skills/reviewer.md` v1.2.0.
- **Modell:** claude-opus-4-8 (Opus 4.8).
- **Datum:** 2026-08-02.
- **Stand:** HEAD `143bf7b` (Working-Tree clean), slice-086 committet über die
  Kette `258f46b`…`143bf7b`.
- **Gates:** `make gates` grün (299 Dateien, 0 Befunde; Coverage 94,20 %);
  `make adr-check` grün (vcs über `HEAD~1..HEAD`, 0 Befunde). Zusatzprüfung: über
  den **vollen** slice-086-Bereich (`258f46b^..HEAD`) ist **keine** Datei unter
  `docs/plan/adr/` geändert.
- **Prüf-Achsen:** Template-Konformität · Link-/Anker-Vollständigkeit ·
  Klassifikation aktiv/aufgelöst · `Ersetzt-Baseline-Regel` · `aufgelöst durch` ·
  §Anforderungs-Anlege-Prozess-Löschung. Verifiziert gegen Quelle
  (vendored `.harness/baseline/v5.0.0/`), Template und Gate — nicht gegen
  Zusammenfassung.

---

## Findings

### F-1 · LOW · MR-021 · `harness/conventions/MR-021-vendored-verweise-pin-gebunden.md:17`

- **befund:** Der **aktive** Eintrag illustriert den „konkreten Pin" der
  Live-Doku-Links mit `…/v1.4.0/…`, obwohl der Baseline-Pin auf `v5.0.0` steht und
  sämtliche Live-Links (u. a. alle acht aktiven MR-Dateien) `…/v5.0.0/…` tragen.
  slice-086 Vorgehen 7 („Prosa entpinnen … die interne `v1.4.0`-Prosa … auf
  `v5.0.0` ziehen") und der Auflösungs-Trigger von MR-023 nennen genau diese Stelle
  (`…/v1.4.0/…` in der Pin-Bindungs-Adaption) als zu aktualisieren; sie ist nicht
  gezogen. Ein aktiver Eintrag (jeder Lauf liest ihn) zeigt damit ein Pin-Beispiel,
  das dem aktuellen Baseline-Stand widerspricht.
- **verifizierbar:** kein Gate greift (Prosa-Beispiel unterhalb der
  `versions`-Pin-Erkennung); per Datei-Lesung belegt.

### F-2 · LOW · MR-006 · `harness/conventions/MR-006-referenzrichtung-matrix.md:28`

- **befund:** Der Eintrag steht in `### Aktive Adaptionen`, sein
  Auflösungs-Trigger benennt aber nur den **bereits eingetretenen** Alt-Trigger
  („Eingetreten mit Baseline `v1.3.0`") mit Ausgang „wird zur reinen
  Baseline-Konformität". Der eigentliche Grund für die Aktiv-Setzung — die in
  dieser Etappe ergänzte, **permanente** C-4-Matrix-Scope-Grenze — steht nicht im
  Trigger. Wendet ein Agent das Split-Kriterium (Trigger eingetreten / Baseline
  eingeholt ⇒ aufgelöst) auf den **deklarierten** Trigger an, folgert er
  fälschlich `done/`; parallel klassifizieren MR-001/MR-010/MR-014 unter genau
  diesem Kriterium als aufgelöst. Die Aktiv-**Platzierung** ist wegen C-4
  vertretbar; das Trigger-**Feld** trägt sie nicht.
- **verifizierbar:** kein Gate greift; per Vergleich Trigger-Text ↔
  Verzeichnis-Position belegt.

### F-3 · LOW · MR-023 · `harness/conventions/MR-023-baseline-v500.md:66`

- **befund:** Der **aktive** Eintrag deklariert als Auflösungs-Trigger, dass „die
  Konventionsspeicher-Migration (spätere Etappe) diese MR in die Datei-je-MR-Form
  fasst und die Prosa … vollständig auf das v5.0.0-Layout angleicht". Diese
  Migration ist die vorliegende Etappe C: die Datei-je-MR-Form ist hergestellt
  (`harness/conventions/MR-023-baseline-v500.md`), der Eintrag bleibt aber per
  Abnahme-Punkt 2 aktiv. Der deklarierte Trigger beschreibt damit ein
  ausgeführtes Ereignis als noch ausstehend; sein zweiter Teil (Prosa-Angleich)
  ist zudem tatsächlich unvollständig (F-1). Ein Agent, der dem Trigger folgt,
  erwartet MR-023 in `done/`.
- **verifizierbar:** kein Gate greift; per Vergleich Trigger-Text ↔
  Verzeichnis-Position + F-1 belegt.

### F-4 · LOW · MR-014 (Provenance) · `harness/conventions/done/MR-014-slice-adr-haus-stil.md:4`

- **befund:** Das `aufgelöst durch`-Feld von MR-014 attribuiert „ADR-Alternativen-
  Tabelle wurde Default" dem `Baseline-Stand v4.0.0`. Der aufgelöste Eintrag
  MR-016 (`harness/conventions/done/MR-016-baseline-pin-hebung-2.md:40`) belegt am
  Tag verifiziert, dass dieselbe ADR-„Verglichene Alternativen"-Tabelle **mit
  `v1.4.0`** Baseline-Default wurde. Zwei aufgelöste Einträge nennen für dieselbe
  Tatsache verschiedene Baseline-Stände; die Auflösungs-Provenance ist in sich
  widersprüchlich. (Betrifft nur `done/`-Historie, nicht den aktiven Lesepfad.)
- **verifizierbar:** kein Gate greift; per Zitat-Vergleich MR-014:4 ↔ MR-016:40
  belegt.

---

## Negativbefunde (geprüft, ohne Befund)

- **Template-Konformität (Abschnitte).** `harness/conventions.md` deckt sich
  Abschnitt-für-Abschnitt mit `.harness/baseline/v5.0.0/templates/harness/conventions.template.md`
  (Purpose · Baseline · Adoptierte Konventions-Quellen · Adaptions-Block · Zusatz-
  klassen · Modus). Einzig das als „(optional)" markierte §Glossar fehlt — konform.
  Kein zusätzlicher Abschnitt. Kein duplizierter Baseline-Normtext (nur Pointer).
- **Adaptions-Block-Form.** `### MR-000` inline mit `Ersetzt-Baseline-Regel: —`;
  `### Aktive Adaptionen` mit vier Spalten (MR/Titel/Geltungsbereich/Ersetzt-
  Baseline-Regel), 8 Zeilen; `### Aufgelöste Adaptionen` mit zwei Spalten
  (MR/aufgelöst durch), 15 Zeilen. Kein Adaptions-Rumpf mehr im Index außer MR-000.
- **Anker-Vollständigkeit (load-bearing).** 24 Definitions-Anker vorhanden: 23
  Voll-Slug-`<a id>` (MR-001…MR-023) in die Index-Zeilen gefaltet + `### MR-000`
  via Überschrift. Repo-weite Set-Differenz aller referenzierten
  `conventions.md#mr-…`-Slugs gegen die definierten Anker ist **leer** in beide
  Richtungen — kein verwaister Slug, kein toter Anker. Das einzige `#mr-xxx` ist
  Prosa-Inline-Code in `docs/reviews/2026-08-01-slice-083-v500-migration-review.md`
  (kein Live-Link). Über den vollen slice-086-Bereich ist **keine** ADR-Datei
  geändert; die 12 immutablen ADR-Links lösen ohne Retarget/ADR-Edit auf.
- **Klassifikation aktiv/aufgelöst.** Alle 8 aktiven + 15 aufgelösten geprüft.
  Diskriminator konsistent angewandt: permanentes repo-spezifisches Regel-Leben ⇒
  aktiv (MR-004/005/007/013/015/021, C-4 in MR-006, Vendoring-Record MR-023),
  Baseline eingeholt/Nachfolger ⇒ aufgelöst. Die vom Prompt hervorgehobene Trias
  ist stimmig: MR-007 aktiv (Dogfooding permanent, nicht Baseline-Default), MR-010
  aufgelöst (9-Rang-Liste ist Baseline-Default), MR-015 aktiv (Routen-statt-
  Spiegeln als permanente Pointer-Disziplin). Abweichungen nur im Trigger-Text
  zweier aktiver Einträge (F-2/F-3), nicht in der Verzeichnis-Position.
- **`Ersetzt-Baseline-Regel` je aktiver Adaption.** Alle acht Anker-Links lösen in
  den vendored Baum auf (`grundlagen-durchsetzungsschicht.md` §Drei Bindepunkte /
  §Vier Design-Eigenschaften; `grundlagen-referenz-richtung.md` §Referenz-Richtung
  (SDP); `modul-13-quality-gates.md` §Hard Rule (Doku-Disziplin);
  `modul-05-planning-harness.md` §Lifecycle als State Machine;
  `grundlagen-harness-dateien.md` §Template-Schichtung / §Verzeichniskonvention —
  Überschriften am Quelltext bestätigt). Semantischer Sitz jeweils tragfähig
  (MR-023/§Template-Schichtung deckt „beide Bäume vendored"; MR-015/§Template-
  Schichtung deckt „Zeiger statt Zitat"). Kein aktiver Eintrag ohne benennbare
  Regel — kein verkappter Fork.
- **`aufgelöst durch` je aufgelöster Adaption.** Die vier gerade korrigierten
  Werte plausibel: MR-002/MR-008 → MR-000 (ID-Schema-Deklaration ist Teil der
  v5.0.0-Baseline-Aussage; d-checks MR-000 trägt sie), MR-010 → Baseline-Stand
  v5.0.0 (9-Rang-Default), MR-014 → Baseline-Stand (Form ist Wahl) — Versions-
  Detail siehe F-4. Die Supersede-Ketten MR-009→MR-010, MR-011→MR-012→MR-016→MR-023
  und MR-003→MR-007 sind schlüssig.
- **§Anforderungs-Anlege-Prozess (Baseline-Duplikat modul-03-spec).** Abschnitt
  aus `harness/conventions.md` entfernt; **kein** Live-Markdown-Link auf
  `conventions.md#anforderungs-anlege-prozess` mehr im Repo (alle Rest-Nennungen
  sind Prosa/Inline-Code). Die vier vormals betroffenen Live-Verweise
  (`docs/plan/planning/done/slice-058-arch-check-via-a-check.md`,
  `…/done/slice-065-suggest-ai-harness-modulset.md`,
  `docs/reviews/2026-07-17-slice-071-implementation-r4.md` ×2) sind auf
  `AGENTS.md#5-dokumentations-regeln` retargetet; die Überschrift existiert
  (`AGENTS.md:167`), und §5 verweist kanonisch auf `modul-03-spec` — keine
  Duplikat-Wiederkehr.
- **`.d-check.yml` ids-Target.** `ids`-Muster `MR-\d{3}` zeigt auf das Verzeichnis
  `harness/conventions/` (Umstellung Einzeldatei → Verzeichnis vollzogen); die
  Linkpflicht ist repo-weit grün.
- **Sektions-Anker.** Die live verlinkten `conventions.md#…`-Sektions-Anker
  (`baseline`, `adoptierte-konventions-quellen`, `modus-deklaration-pro-sub-area`)
  existieren weiterhin.

## Beobachtungen (INFO — kein Handlungszwang)

- **Feld-Modell der `done/`-Einträge.** Aufgelöste Einträge lassen das
  Template-Pflichtfeld `Ersetzt-Baseline-Regel` weg und führen stattdessen
  `Aufgelöst durch:`; aktive Ablöser (z. B. MR-007) tragen kein `Löst auf` /
  `Ausgelöst durch Baseline-Stand`, sondern dokumentieren die Beziehung in Prosa +
  Index. Für historische bzw. slice-getriebene Auflösungen (Trigger ≠ Baseline-
  Bump) vertretbar, aber ein wörtlicher Abstand zum Einzel-Eintrag-Template.
- **MR-018 vs. Etappe D.** Die Template-Freiheit ist in Etappe C aufgelöst, die
  co-located Vorlagen sind Etappe-D-Arbeit; im Fenster zwischen C und D trägt
  d-check weder die Template-Freiheits-Deklaration noch die Vorlagen. Durch die
  A→D-Reihenfolge bewusst in Kauf genommen.

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 2 |

## Verdikt

**Abnahmereif.** Der load-bearing Kern — Index + Datei-je-MR, der in die
Index-Zeilen gefaltete Voll-Slug-Anker-Block und die dadurch ohne ADR-Edit
erhaltenen `conventions.md#mr-…`-Links — ist korrekt: alle referenzierten Slugs
lösen auf, keine `Accepted`-ADR wurde berührt, `make gates` und `make adr-check`
sind grün. Es bleiben ausschließlich vier LOW-Befunde (drei stale/inkonsistente
Trigger-/Prosa-Stellen an sonst korrekt platzierten Einträgen, eine Versions-
Provenance-Diskrepanz in der `done/`-Historie) und zwei INFO-Beobachtungen; keiner
blockiert die Abnahme. F-1 (aktives `…/v1.4.0/…`-Beispiel) und F-2/F-3 (Trigger-
Text ↔ Aktiv-Platzierung) sind die lohnendsten Nacharbeiten, da sie im aktiven
Lesepfad liegen.
