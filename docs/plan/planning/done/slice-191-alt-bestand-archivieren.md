# Slice slice-191: Den Alt-Bestand archivieren (welle-01…welle-85)

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
echten Bestand. Dieser Slice zieht die Trennung nach.

**Scope-Erweiterung während der Ausführung (Nutzer-Entscheid):** ursprünglich
auf welle-60…85 begrenzt (die Wellen, für die [welle-87](../welle-87-wellen-archivierung.md)
die Nachrüstpflicht maß). Nach deren vollständiger Archivierung zeigte die
Sichtung der verbliebenen `docs/reviews/`-Altlast (158 Dateien), dass die
Wellen-Nummerierung **nicht** erst bei welle-60 beginnt, sondern bei
welle-01 — jeder Slice vor welle-60 trägt bereits ein gültiges
`**Welle:**`-Feld (z. B. `slice-001` → `welle-01-fundament`, `slice-050` →
`welle-39-matrix-richtung`), nur ohne flache Welle-Plan-Datei in `done/`
(dieselbe, bereits gelöste Eigenheit wie welle-60…66). Ein Sammel-Archiv für
„wellenlosen Bestand vor der Einführung" ([Modul 6](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6))
ist damit **nicht** nötig — der Bestand ist durchgängig Welle-zugeordnet.
**Neuer Umfang: welle-01 bis welle-85**, mit demselben Werkzeug, demselben
Wellen-für-Welle-Vorgehen.

**Zwei Bestands-Eigenheiten wurden bereits beim Planen gemessen, nicht
erst beim Ausführen entdeckt** — genau die Vorsicht, die
[slice-190 §5](../done/slice-190-wellen-archiv-werkzeug.md) (Risiko 3,
`BEO-011` sechste Instanz) benannt hat:

1. **welle-60 bis welle-66 haben keinen Welle-Plan.** `ls done/welle-6*`
   zeigt nur `welle-60-results.md` … `welle-66-results.md` — keine
   `welle-6N-<titel>.md`. Diese sechs Wellen wurden **vor** der
   Plan-Datei-Konvention geschlossen und tragen laut
   [slice-088](../done/welle-67/slice-088-etappe-d1-doc-form.md) Zeile 131 eine
   „retroaktiv markierte, minimale Ergebnis-Notiz" als bewusste
   Nutzer-Entscheidung nach — nicht eine vollständige Plan-Rekonstruktion.
   `tools/archive-wave`s `FindWellePlan` verlangt **genau eine** Treffer-Datei
   und schlägt für diese sechs fehl.
2. **Die vier vermeintlich wellenlosen Alt-Slices sind keine.** Gemessen:
   `done/slice-*.md` ohne `**Welle:**`-Feld sind genau vier —
   [slice-077](../done/welle-60/slice-077-stiller-tabellen-uebersprung.md),
   [slice-078](../done/welle-61/slice-078-ignore-refs-quell-skopus.md),
   [slice-079](../done/welle-62/slice-079-zitat-verifikation.md),
   [slice-103](../done/welle-74/slice-103-geteilte-lexik-raender.md). Alle vier
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
   **keine** wellenlosen Alt-Slices mehr vor der Einführung der Regel.
   **Präzisierung, nach Abschluss aller 85 Wellen gemessen:** diese Aussage
   galt nur der alten Lücke „Feld ganz abwesend" (vier Fälle). Ab
   `slice-137` führt dieses Repo eine **zweite, bewusste** Wellenlos-Form —
   das Feld trägt ausdrücklich `— wellenlos` (Baseline-Regelwerk
   [`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wann-arbeit-eine-welle-braucht-modul-6):
   „Wellenlos heißt nicht wächterlos"). Das sind **52** Slices — gemessen
   per Feld-Wert, nicht per Feld-Abwesenheit (`slice-137` bis `slice-189`,
   plus sechs ältere Einzelfälle mit gleichwertiger Prosa:
   `slice-083/095/102/112/121/127`) — durchweg **aktueller** Bestand, nicht
   Vor-Einführungs-Altlast, und ohne Wellen-Zugehörigkeit **by design**: ihr
   Closure-Grund ist ihre eigene DoD, keine Wellen-Closure-Bedingung. Für sie
   gilt keine Archivierungspflicht — Modul 6 Schritt 4 archiviert **Wellen**,
   und eine Welle, die es nie gab, hat nichts zu schließen. Die zugehörigen
   `docs/reviews/`-Reports bleiben aus demselben Grund unangetastet. Die
   welle-87-§3-Bedingung „Zuordnung entschieden" bezieht sich auf den
   **Vor-Einführungs**-Bestand (Modul 6: „…für den Bestand vor der
   Einführung") und ist mit welle-01…85 vollständig erfüllt — ohne
   Sammel-Archiv oder chronologische Näherung.
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

- [x] `**Welle:**`-Feld für slice-077/078/079/103 nachgetragen, mit
      Beleg-Zeile in diesem Slice-Plan (§1) belegt.
- [x] `tools/archive-wave` behandelt eine Welle ohne Plan-Datei
      (welle-60…66) korrekt — Test vorhanden, `make archive-wave-test` grün.
- [x] welle-01 bis welle-85 sind archiviert: `done/<welle-id>/archiv.zip`
      existiert je Welle, Stubs an der Zielform, Review-Reports ohne Stub,
      keine gebrochenen Repo-Verweise (`make doc-check` grün nach jeder
      angewendeten Welle).
- [x] `make gates` **und** `make fullbuild` grün auf dem archivierten
      Bestand (Exit explizit — 50 Requirements, 0 Waisen).
- [x] **Unabhängiger Review und unabhängige Verifikation.**
- [x] Closure-Notiz mit Lerneintrag; jedes Risiko aus §5 mit Ausgang; die
      drei Modul-6-Paarungen geprüft (bei der Closure von welle-87, die
      dieser Slice schließt).

## 5. Abnahme-Punkte / Risiken

- **26 (später 85) Wellen in Folge zu archivieren ist viel Wiederholung —
  eine falsche Annahme im Werkzeug wirkt N-fach, nicht einmal.** —
  **Ausgang: eingetreten, aber durch die Mitigation vollständig
  aufgefangen.** Fünf echte Lücken traten tatsächlich auf
  (Ortsfeste-Verweise-Idiom bricht beim Tiefenwechsel, welle-70;
  Fail-Closed-Guard zu streng für einen legitimen Plan-ohne-Slices-Fall,
  welle-73; `Hervorgegangen:`-Feld trug nie die überlebenden Kennungen,
  sichtbar erst bei `make fullbuild`; der erste Nachbesserungs-Regex strich
  beim Rückbau selbst wieder echte Link-Label-Zitate, selbst entdeckt vor
  dem Commit des Backfills; und — von einer ZWEITEN unabhängigen Prüfung
  gefunden, nicht selbst — eine Shell-Näherung des (korrekten) Go-Codes
  behandelte punkt-suffigierte Kennungen wie [`DC-FA-REQ-001.a`](../../../../spec/spezifikation.md#dc-fa-req-001a--anforderungsquellen-headings-und-tabellen) falsch und
  verschluckte sie ganz, während sechs schon 2026-09-03 früh archivierte
  Wellen (welle-03/08/28/34/47/56, vor dem Inline-Code-Fix) unbemerkt bare
  Fremd-/erfundene Kennungen trugen). Keine der fünf von einem Gate
  zwischen den Wellen gefangen — `make fullbuild`s Vollständigkeits-Check
  bleibt grün, solange irgendein anderer Slice dieselbe Anforderung auch
  zitiert; er sieht keine falsche oder fehlende einzelne Kennung. Die
  Schrittweite (Commit je Welle) hat trotzdem die Blast-Radius begrenzt:
  keine der 85 Wellen musste zurückgenommen werden, und die **endgültige**
  Korrektur lief nicht mehr per Handnäherung, sondern per **temporärem
  Go-Test, der die echte `ExtractSurvivingIDs`-Funktion gegen jedes
  `archiv.zip` aller 85 Wellen fuhr** und jede Abweichung von jedem
  committeten Stub meldete (zehn Funde, alle korrigiert, Re-Lauf: 0
  Abweichungen) — die Lehre aus dem vierten Fund (Shell-Näherung ≠ Produkt)
  unmittelbar angewandt, nicht nur notiert.
- **Der Feld-Nachtrag (§2 Punkt 1) ändert historische `done/`-Dateien.** —
  **Ausgang: entfallen** — die Grenze wurde durch den Akt selbst geprüft und
  bestätigt: der `pre-commit`-Hook (inkl. `adr-check`) akzeptierte den
  Welle-Feld-Nachtrag in vier `done/`-Slices anstandslos, weil er ein
  ANDERES Feld als das AGENTS.md §3.7-geschützte `**Status:**` ist. Kein
  Gate-Widerspruch in über 60 Commits dieses Slices.
- **Der Verweis-Nachzug über N Wellen kann eine Datei mehrfach treffen.** —
  **Ausgang: entfallen** — genau dieser Fall trat wiederholt ein (ein
  Review-Report akkumulierte über vier separate Wellen-Läufe hinweg vier
  `ignore-refs`-Einträge, additiv, ohne Konflikt) und `RewriteRepo` verhielt
  sich dabei in jedem gemessenen Fall korrekt.
- **`welle-87`s eigene Closure hängt an diesem Slice.** — **Ausgang:
  entfallen** — dieser Slice schließt jetzt, welle-87s Closure-Prozedur
  folgt unmittelbar danach im selben Zug.

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

**Geliefert:** `**Welle:**`-Feld für vier historische Slices nachgetragen
(slice-077/078/079/103); `tools/archive-wave` um den Plan-losen-Welle-Fall
(welle-60…66), die Review-Selbstreferenz-Vorschaufilterung, das
Tiefenwechsel-korrekte Feld-Retargeting (`RewriteFieldForMove`), den
gelockerten Fail-Closed-Guard und die Kennungs-Übernahme ins
`Hervorgegangen:`-Feld erweitert — alle fünf an echten Bestands-Eigenheiten
gemessen, nicht spekulativ gebaut. **Alle 85 nummerierten Wellen dieses
Repos (welle-01 bis welle-85) archiviert:** je ein `archiv.zip`, Slice-
und Welle-Stubs an der Zielform, Review-Reports ohne Stub. `make gates`
und `make fullbuild` grün auf dem vollständig archivierten Bestand (50
Requirements, 0 Waisen). 52 Slices ab `slice-137` bestätigt als bewusst
wellenlos (Baseline-Konvention „wellenlos heißt nicht wächterlos") und
damit korrekt außerhalb der Archivierungspflicht. Nach unabhängigem
Review zusätzlich: zwei durch einen zu engen `git add`-Scope verpasste
Verweis-Nachzüge außerhalb von `docs/plan/planning`/`docs/reviews`
nachgeholt ([ADR-0058](../../adr/0058-konfigurations-flaechen-additiv-weiten.md),
[MR-003](../../../../harness/conventions/done/MR-003-vendorter-bootstrap-sensor.md)),
und [`MR-059`](../../../../harness/conventions/MR-059-wellen-archiv-stub-move.md)
plus `AGENTS.md` §3.3 vierte Ausnahme neu geschrieben — die bislang nur
implizit gelebte Ein-Commit-Form der Wellen-Archivierung ist damit
kanonisch benannt, nicht länger nur durch einzelne Commit-Botschaften
behauptet.

**Was funktioniert hat:** Die Wellen-für-Welle-Disziplin (§2 Punkt 3) hat
sich vollständig ausgezahlt — jede der drei während der Ausführung
gefundenen Werkzeug-Lücken wurde durch `make doc-check` nach genau der
Welle sichtbar, die sie auslöste, nie erst Dutzende Wellen später. Der
pre-commit-Hook (`doc-check` + `adr-check`) hat bei über 60 Commits kein
einziges Mal fälschlich blockiert und zweimal echte, unvorhergesehene
Belegbrüche (ADR-Geschichte-Referenzen, Review-Report-Ketten) korrekt
gemeldet.

**Was anders lief:** Der Scope wuchs zweimal während der Ausführung, beide
Male durch Nutzer-Entscheid und beide Male, weil eine Messung mehr zeigte
als der ursprüngliche Plan angenommen hatte: zuerst die Erkenntnis, dass
die Wellen-Nummerierung bei welle-01 beginnt (nicht welle-60), dann die
explizite Bestätigung, auch die docs/reviews/-Altlast „regelkonform" zu
archivieren — was sich als bereits durch die Wellen-Archivierung
eingelöst herausstellte, sobald der Scope korrekt war. Drei echte
Werkzeug-Lücken (siehe §5) waren beim Bauen in slice-190 nicht vorhersehbar,
weil sie nur an gewachsenem, echtem Bestand auftreten — exakt die von
`BEO-011` benannte Grenze einer Fixture-Beweisführung.

**Steering-Loop-Einträge:**
- [`BEO-011`](../observations.md) sechste Instanz (slice-190) — **Ausgang:
  weiter offen, nicht verkörpert (Review-Korrektur).** Die Vorsicht, die
  die Beobachtung einforderte, traf zu (fünf echte, im Fixture nicht
  sichtbare Werkzeug-Lücken) und wurde durch das tatsächlich gefahrene
  Vorgehen — Welle-für-Welle mit Gate-Prüfung zwischen jedem Schritt statt
  Massenanwendung — einzeln und ohne Rückabwicklung aufgefangen. Das ist
  aber eine einmalig sorgfältig gefahrene Prozedur, keine Verkörperung im
  Modul-6-Sinn: es fehlt der Zielort samt Herkunfts-Anker, der die Regel
  „ein fixture-bewiesenes Werkzeug wird am echten Bestand inkrementell mit
  Gate-Prüfung zwischen jeder Einheit ausgerollt" für künftige Läufe
  auffindbar machte. Bleibt im Register offen.
- Neue Prozedur-Lehre (keine gesonderte Regel-Datei, in diesem Closure
  festgehalten): ein Archivierungswerkzeug, das Volltexte durch Stubs
  ersetzt, muss die überlebenden Kennungen aktiv ins Stub-Feld übernehmen
  (`Hervorgegangen:`), sonst bricht `--require-complete` lautlos für jede
  Anforderung, deren einziger Beleg-Slice archiviert wird — kein anderes
  Gate hätte das gefunden.
- Zweite Lehre, aus dem Rückbau selbst: eine textbasierte Kennungs-
  Extraktion muss zwischen einer **echten** Zitat-Form
  (`` [`DC-FA-XXX-001`](ziel) ``) und einer **illustrativen**
  (Inline-Code ohne folgenden Link, z. B. eine erfundene Kennung als
  Beispiel für einen Parsing-Grenzfall) unterscheiden — ein blindes
  Wegwerfen aller Inline-Code-Spannen ist keine sichere Näherung, sondern
  verschluckt echte Belege (gemessen: 3 von 4 Kennungen in slice-073
  verschwanden in der ersten Fassung). Vor jeder Commit dieser Art: die
  Fund-Anzahl gegen eine plausible Grundgesamtheit prüfen, nicht nur „lief
  durch".
- Dritte Lehre, aus der zweiten unabhängigen Prüfung: die eigene
  Shell-Näherung der zweiten Lehre trug denselben Fehler in neuer Form
  (punkt-suffigierte Kennungen fielen durch ihre eigene
  Link-Label-Erhalt-Regel). Aufgelöst nicht durch eine dritte Regex-
  Iteration, sondern durch den Wechsel des Mittels: ein **Go-Test, der die
  echte, bereits getestete Funktion aufruft**, statt einer weiteren
  Handnäherung ihrer Logik — genau AGENTS.md „Produkt vor grep/awk vor
  allem anderen", hier auf eine Mess-**Korrektur** angewandt, nicht nur auf
  eine Mess-**Anzeige**.

**Zeiger:** [Beobachtungs-Register](../observations.md). Diese Slice-
Closure löst zugleich [welle-87](../welle-87-wellen-archivierung.md)s
eigene Closure-Prozedur aus (Modul 6, sechs Schritte) — sie folgt
unmittelbar danach.
