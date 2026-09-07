# Review-Report — slice-208, Runde 2 (Form-Migration und die drei Handlungen, DoD 2 und 3)

**Review-Art:** Code/Design — geprüft werden die **Form-Migration** und die **drei aus dem Delta-Audit abgeleiteten Handlungen** (R4 · R5/`MR-068` · R7/`MR-069`, dazu R6 im Reviewer-Skill) gegen den Slice-Plan, den Kanon `v6.5.0` und den eigenen Bestand. **Nicht** geprüft: der Delta-Audit als solcher — das war Runde 1; hier werden nur die **Korrekturen** an ihm nachgemessen.
**Gegenstand:** `33c24e1` (Audit regel-weise, Runde-1-Befunde) · `8be0950` (Form-Migration, `MR-068`, `MR-069`, DoD 2/3). Diff-Range `3d15de7..HEAD`.
**Skill:** `.harness/skills/reviewer.md` v1.13.0 @ `8be0950` — **in der Fassung, die dieser Slice selbst geändert hat** (neue §Zitier-Form); sie ist Teil des Gegenstands.
**Modell-ID:** claude-opus-5[1m] · **Datum:** 2026-09-07
**Eingangs-Kontext:** der Slice-Plan slice-208 (DoD 2/3, §3, §5, §2-Audit); der Report von Runde 1; `AGENTS.md` §3.5/§3.6/§3.7/§4/§5/§6; `harness/README.md`; `harness/conventions.md` samt `MR-000`, `MR-021`, `MR-039`, `MR-049`, `MR-053`, `MR-054`, `MR-056`, `MR-066`, `MR-067`; die Vorlage `MR-NNN-titel.template.md`; `.d-check.yml` (`trace`, `ignore-refs`) und `.d-check.closure.yml`; `DC-FA-REF-001`, `DC-FA-STRUCT-001`, `DC-FA-CLI-011`, `DC-FA-COV-001`; Baseline `v6.5.0` · `regelwerk/grundlagen-traceability.md`, `regelwerk/grundlagen-harness-dateien.md`, `regelwerk/modul-05-planning-harness.md`, `regelwerk/modul-09-implementierung.md` sowie `templates/docs/plan/planning/slice.template.md`, `templates/docs/reviews/review-report.template.md`, `templates/docs/plan/adr/NNNN-titel.template.md`.
**Eigene Läufe:** `make gates` (zehn Gates grün, **693** Dateien / 0 Befunde, Coverage 94,60 %) · `make verify-closure-notes` (592 Dateien / 0 Befunde) · **vier eigene Bruch-Tests** am Closure-Profil (Haus-Form · Baseline-Form · Nicht-Treffer · Doppel-Treffer) · **ein eigener Bruch-Test** am `ignore-refs`-Ventil über `make doc-check` · repo-weite Nachmessung der Ventil-Zahlen (25 Einträge / zehn Tags / 28 Dateien, klassenweise) · Existenz-Prüfung aller `in:`-Skopen · Zählung der Haus-Form-`done/`-Slices und ihrer Ausschluss-Abschnitte. Arbeitsbaum nach jedem Bruch-Test über `git checkout --` wiederhergestellt; `git status` am Ende leer.

**Zitier-Form.** Dieser Report ist ein einfrierendes Artefakt und folgt der Form, die dieser Slice gerade eingeführt hat: Baseline-Stellen als Tag + Pfad in Inline-Code, Slices als Kennung statt als Lifecycle-Pfad, Gates als `make <target>`-Token.

---

## Was gemessen reproduziert (vorweg, damit die Findings nicht als Gesamturteil gelesen werden)

| Behauptung | Quelle | Nachmessung |
|---|---|---|
| Bruch-Test je Richtung meldet `section-forbidden` | Commit-Botschaft `8be0950` | **reproduziert** — beide Richtungen, beide auf Zeile 170 |
| „zehn Gates, 693 Dateien, 0 Befunde" | Commit-Botschaft `8be0950` | **reproduziert** exakt |
| 25 Ventil-Einträge über zehn entfernte Tags | `MR-069` | **reproduziert** exakt |
| 28 Dateien: 18 `Accepted`-ADRs, 6 aufgelöste `MR`, 3 `done/`-Slices, 1 CR | `MR-069` | **reproduziert** exakt, Summe und jede Klasse |
| acht `done/`-Slices führen den Ausschluss-Abschnitt, zwei nennen einen echten Folge-Slice | Audit §2 (R2, nach Runde 1 korrigiert) | **reproduziert** exakt (202→203, 207→208) |
| neun Haus-Form-Abschnitte gegen acht der Baseline | Audit §5 | **reproduziert** |
| `DC-FA-CLI-011` §Out-of-Scope: *„eine bloße ADR-Referenz ohne Slice/Coverage deckt weiterhin **nicht** ab"* | `MR-068` | **wortgleich** |
| Kanon-Zitat *„ein Ausnahme-Ventil im Prüfbereich, also eine Gate-Senkung mit eigener Begründungslast"* | `MR-069` | **wortgleich** |

---

## Findings

### F-1 · HIGH · `MR-069`s zentrale Behauptung über die Reichweite des Ventils ist durch Messung widerlegt

- **quelle:** `DC-FA-REF-001`
- **pfad:** `harness/conventions/MR-069-ignore-refs-ist-eine-deklarierte-gate-senkung.md` (Abschnitt *Adaption*, dritter Absatz)
- **befund:** Der Eintrag schreibt: *„im Prüfbereich hinter dem Ventil meldet `links` nichts mehr, auch nicht über einen Defekt, der nichts mit der Baseline zu tun hat."* Das ist der Satz, mit dem er die Charakterisierung als **Gate-Senkung** trägt. Gemessen ist das Gegenteil: ein in eine vom `in:`-Glob `docs/plan/adr/**` gedeckte Datei eingesetzter, **nicht** baseline-bezogener Link meldet weiterhin — `docs/plan/adr/0069-…:155 0999-gibt-es-nicht-r2-bruchtest.md target-missing`, `make doc-check` Exit ≠ 0. `DC-FA-REF-001` sagt es auch: *„das Ventil unterdrückt **nur** die Auflösungs-Klasse des genannten Ziels … keine anderen Befunde"*. Der Satz verwechselt den **Quell**-Skopus (`in:`) mit dem **Ziel**-Skopus (`refs:`); alle 25 Einträge tragen ein enges `refs: [".harness/baseline/<tag>/**"]`. Der Eintrag beschreibt damit eine Senkung, die breiter ist als die, die er deklariert.
- **verifizierbar:** ja — `make doc-check` nach Einsetzen eines beliebigen toten Nicht-Baseline-Links in eine der gedeckten ADRs.
- **klasse:** `ventil-reichweite-ueber-refs-skopus-hinaus-behauptet`

### F-2 · HIGH · Die Migration ist unvollständig: sechs lebende Träger nennen weiter die Haus-Form-Abschnitte, darunter die Deklaration der gerade geänderten Gate-Regel

- **quelle:** DoD (3) („Die Konventions-Einträge sind nachgezogen"); `AGENTS.md` §3.7 (Zustandsfeld); `BEO-ALL/semantic-change-body-only-edges-stale`
- **pfad:** `AGENTS.md:436` · `harness/conventions.md:128` · `harness/conventions/MR-049-ausgangs-wortschatz.md:12-13` (dazu Z. 27, Z. 50) · `harness/conventions/MR-053-dritte-vorpruefung-nachtlauf.md:9` und `harness/conventions.md:132` · `harness/conventions/MR-054-vorpruefungen-belegen-ihre-regel.md:11-12` · `docs/plan/planning/observations/README.md:18` · `AGENTS.md:538`
- **befund:** Die Regel im Closure-Profil deckt jetzt beide Titel und die offene Ziffer; ihre **Deklarations-Träger** tun es nicht. `AGENTS.md` §4 sagt weiterhin *„der **erfundene** als Komplement-`forbid-pattern` über **§5**"*, die Index-Zeile zu `MR-049` *„`.d-check.closure.yml`, **§5** der `done/`-Slices"*, und `MR-049`s Geltungsbereich *„Abschnitt **§5** der Slice-Dateien unter `docs/plan/planning/done/`"* — für jeden Slice in der adoptierten Form ist es §6. `MR-053` und `MR-054` skopieren auf einen *„§Vorgelagert"*, den die adoptierte Form nicht mehr hat (die zwei Blöcke sind unbedingter Kopf von §8, `v6.5.0` · `templates/docs/plan/planning/slice.template.md`), und `observations/README.md` nennt *„§7/§8 jedes Slice-Plans"* für den Sichtungs-Schritt, der jetzt allein in §8 liegt. `AGENTS.md` §5 nennt die Zielsektion weiter *„Sub-Area-Modus-Begründung"* statt *„Sub-Area-Prüfungen und Modus-Begründung"* und stellt die drei Vorprüfungen *„**vor**"* sie statt in ihren Kopf. Dazu die Gegenrichtung: der neue Schritt 4 in `AGENTS.md` §6 verweist auf *„§1 des Plans"* — beide lebenden Pläne dieses Repos (slice-208 in `in-progress/`, slice-209 in `open/`) schließen in **§3** aus; die Regel ist am Tag ihrer Einführung für den gesamten lebenden Bestand falsch adressiert. §3 des Plans nimmt nur die sechs (tatsächlich acht, s. F-3) `done/`-Slices aus, nicht diese Träger.
- **verifizierbar:** ja — Textvergleich der genannten Zeilen gegen `v6.5.0` · `templates/docs/plan/planning/slice.template.md`; kein Gate hält das.
- **klasse:** `semantik-geaendert-deklarations-traeger-stehen-gelassen`

### F-3 · MEDIUM · Die von Runde 1 widerlegte „Sechs" steht an vier Stellen des Plans und in der Commit-Botschaft weiter

- **quelle:** `AGENTS.md` §5 (*„behauptet nicht mehr, als die Arbeit trägt"*); `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:41` (§2), `:166` (§3), `:221` (§5), `:326` (§7); Commit-Botschaft `8be0950`, letzter Absatz
- **befund:** Runde 1 hat gemessen, dass **acht** `done/`-Slices die Haus-Form tragen (F-6 dort), und der Implementer hat genau **eine** Fundstelle nachgezogen — den R2-Absatz in §2 (*„acht geschlossene Slices führen den Abschnitt"*). Die vier übrigen Vorkommen von „sechs" stehen unverändert, darunter **§3 Ausdrücklich NICHT** („Die sechs Slices in Haus-Form bleiben, wie sie sind") — also die Stelle, die den Umfang des Ausschlusses normativ festlegt. Der Plan widerspricht sich damit in derselben Datei: acht Slices führen den Abschnitt, sechs sollen ausgenommen sein. Die Commit-Botschaft trägt die Zahl ein zweites Mal nach außen. Nachgemessen: acht Dateien tragen **beide** Haus-Form-Überschriften (slice-200 bis slice-207).
- **verifizierbar:** ja — `grep -rl '## 3\. Ausdrücklich NICHT' docs/plan/planning/done/` liefert acht.
- **klasse:** `review-befund-nur-an-einer-fundstelle-eingearbeitet`

### F-4 · MEDIUM · `MR-068` deklariert eine Abweichung, die dieses Repo nicht praktiziert, und sein Geltungsbereich nennt einen Konfigurationsschlüssel, den es nicht gibt

- **quelle:** `MR-000` (Adaption = Abweichung **dieses Repos**); `DC-FA-COV-001`; `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`
- **pfad:** `harness/conventions/MR-068-trace-coverage-als-dritte-referenzklasse.md` (Geltungsbereich, Adaption)
- **befund:** Der Geltungsbereich lautet *„`trace.coverage` in `.d-check.yml`"*. Der `trace:`-Block dieser Datei führt `requirements`, `adrs` und `slices` — **kein** `coverage`. Der Kanon-Auslöser lautet *„Ein Repo, das das anders schneidet … deklariert sie"* (`v6.5.0` · `regelwerk/grundlagen-traceability.md` §Die zweite Richtung: Anforderung → Beleg); dieses Repo schneidet **nicht** anders — seine RTM entlastet ausschließlich über Slices, exakt nach dem Kanon-Vorschlag. Was abweicht, ist eine **Produkt-Fähigkeit** (`DC-FA-COV-001`, strikt opt-in, default-aus byte-identisch), nicht die Konfiguration dieses Repos. Der Eintrag überträgt eine Regel über die eigene Konfiguration auf den Funktionsumfang des eigenen Werkzeugs — dieselbe Verwechslung von eigener und fremder Menge, die Runde 1 als F-1 an derselben Regel gefunden hat. Die Adaption selbst ist konditional formuliert (*„Ist sie konfiguriert"*) und damit ehrlich; der Geltungsbereich ist es nicht.
- **verifizierbar:** ja — `grep -n 'coverage' .d-check.yml` liefert nur einen Treffer, und der betrifft `make mention-coverage`.
- **klasse:** `produkt-faehigkeit-als-repo-abweichung-deklariert`

### F-5 · MEDIUM · Die „Auflösung nach vorn" von `MR-069` erreicht die Klasse nicht, die das Ventil wachsen lässt

- **quelle:** `MR-069` (Auflösungs-Trigger, Grenze); eigene Messung
- **pfad:** `harness/conventions/MR-069-ignore-refs-ist-eine-deklarierte-gate-senkung.md` (Grenze/Auflösungs-Trigger) gegen `.harness/skills/reviewer.md` (Ablage/Zitier-Form)
- **befund:** Der Trigger lautet *„kein neuer Eintrag über drei aufeinanderfolgende Pin-Hebungen"*, und die Auflösung soll die Kennung-statt-Adresse-Form leisten. Die **einzige** Handlung, die dieser Slice dafür setzt, steht im Reviewer-Skill und regiert **Review-Reports**. Von den 28 gemessenen Dateien hinter dem Ventil ist **kein einziger** ein Review-Report; 18 sind `Accepted`-ADRs, 6 aufgelöste `MR`-Einträge, 3 `done/`-Slices, einer ein CR. Die dominierende Klasse — ADRs — hat keinen Träger bekommen, obwohl die gelebte ADR-Praxis dieses Repos die `Regeln:`-Zeile weiter als **Link** in den vendorten Baum schreibt (ADR-0084, ADR-0069) und die Baseline-Vorlage `v6.5.0` · `templates/docs/plan/adr/NNNN-titel.template.md` bereits auf die Inline-Code-Form ohne Link umgestellt hat. Der Trigger kann so nicht feuern; der Eintrag benennt diese Lücke nicht.
- **verifizierbar:** ja — Klassen-Zählung der 28 Dateien; `grep -n 'Regeln:' docs/plan/adr/0084-*.md` zeigt die Link-Form.
- **klasse:** `handlung-adressiert-nicht-die-gemessene-ursache`

### F-6 · MEDIUM · Sieben der 25 Ventil-Einträge zeigen auf Quelldateien, die es nicht mehr gibt — `MR-069` führt sie als wirksame Skopen

- **quelle:** `MR-069` (*„jeder Eintrag so eng wie möglich zu skopieren"*); Präzedenz: der Review zu slice-199 (verwaiste `ignore-refs`)
- **pfad:** `.d-check.yml` §`ignore-refs` — `docs/plan/planning/done/slice-080-…`, `…/slice-081-…`, `…/slice-109-…`, `…/slice-183-…`, `docs/reviews/2026-08-02-slice-091-…`, `docs/reviews/2026-08-09-backlog-schnitt-review.md`, `docs/reviews/2026-08-09-slice-096-…`
- **befund:** Die sieben datei-skopierten Einträge nennen Pfade, die nach den Archivierungen (Wellen-Archiv bzw. `-slice=`/`-review=`) nicht mehr existieren — slice-183 liegt heute unter `done/wellenlos/`, der backlog-schnitt-Report unter `docs/reviews/archiv/`. Sie decken nichts. `MR-069` ist genau der Eintrag, der *„eine Aussage über das **Ganze**"* sein will, und stellt die 25 als in Kraft stehende Skopen dar; 28 % davon sind tote Konfiguration. Bemerkenswert: der Kommentar über dem `v6.3.1`-Eintrag beschreibt, dass zwei **Glob**-Einträge mit slice-199 entfernt wurden — der datei-skopierte Eintrag auf dieselbe archivierte Datei blieb stehen. Der Bestand ist älter als dieser Slice; die Behauptung über ihn ist neu.
- **verifizierbar:** ja — Existenz-Test aller `in:`-Skopen; kein Gate meldet einen verwaisten `ignore-refs`-Eintrag.
- **klasse:** `ventil-eintrag-ueberlebt-sein-ziel`

### F-7 · MEDIUM · Der Reviewer-Skill bekommt eine neue Regel, aber keine neue Version — und Reports identifizieren ihn über die Version

- **quelle:** `.harness/skills/reviewer.md` §Ablage (*„**Skill** (`reviewer.md` @ Version/Commit)"*); `AGENTS.md` §3.7 (Zustandsfeld)
- **pfad:** `.harness/skills/reviewer.md:3`
- **befund:** Der Kopf steht unverändert auf *„**Version:** 1.13.0 · **Datum:** 2026-08-28"*, während `8be0950` eine ganze Regel-Sektion (§Zitier-Form) hinzufügt. Die gelebte Praxis dieses Repos bumpt bei jeder inhaltlichen Regel-Erweiterung (1.10.0, 1.11.0, 1.12.0, 1.13.0 tragen je eine); unverändert bleibt die Version bisher nur bei rein mechanischen Pin-Nachzügen. Damit bezeichnen zwei verschiedene Skill-Inhalte dieselbe Version: der Report von Runde 1 nennt *„v1.13.0 @ `3d15de7`"* (ohne §Zitier-Form), dieser Report müsste *„v1.13.0 @ `8be0950`"* nennen (mit ihr). Das Versionsfeld trägt die Unterscheidung nicht mehr, nur noch der Commit.
- **verifizierbar:** nein (kein Gate) — nachweisbar über `git log -p -- .harness/skills/reviewer.md`.
- **klasse:** `regel-erweiterung-ohne-versions-bump`

### F-8 · MEDIUM · Die §Zitier-Form gibt die Ausnahme für das `pfad`-Feld weiter wieder als der Kanon

- **quelle:** `v6.5.0` · `templates/docs/reviews/review-report.template.md` §Zitier-Form; Reviewer-Skill-Prüffrage 9
- **pfad:** `.harness/skills/reviewer.md:221-222`
- **befund:** Der Skill schreibt: *„Das `pfad`-Feld eines Findings bleibt davon **unberührt**: Es benennt die Fundstelle im geprüften Stand und ist damit selbst Messung, nicht Verweis."* Die Quelle sagt weniger: *„Ein `pfad`-Feld auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den Stand des Laufs und **darf ihn festhalten**"* — und führt als Beispiel genau die gebundene Form vor (`` `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> ``, ausdrücklich *„diese Zeile ist selbst ein Beispiel der Form"*). Die Ausnahme des Kanons betrifft die **Kennung-statt-Adresse**-Hälfte, nicht die Form: der Tag bleibt. „Unberührt" lässt einen Leser schließen, das `pfad`-Feld dürfe in den vendorten Baum **verlinken** — also genau das, wogegen die Regel gebaut ist.
- **verifizierbar:** nein (kein Gate) — Textvergleich mit der Vorlage.
- **klasse:** `ausnahme-breiter-wiedergegeben-als-die-quelle`

### F-9 · MEDIUM · `MR-068` und `MR-069` tragen im Pflichtfeld `Ersetzt-Baseline-Regel` einen Datei-Link ohne Abschnitts-Anker

- **quelle:** `v6.5.0` · `templates/harness/conventions/MR-NNN-titel.template.md` (*„als Link mit Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel"*)
- **pfad:** `harness/conventions/MR-068-trace-coverage-als-dritte-referenzklasse.md:4-9` · `harness/conventions/MR-069-ignore-refs-ist-eine-deklarierte-gate-senkung.md:4-8`
- **befund:** Beide Einträge verlinken die **Datei** (`…/grundlagen-traceability.md` bzw. `…/grundlagen-harness-dateien.md`) und tragen den Abschnitt nur im Link-Text. Bei `MR-068` existiert der Anker und fehlt nur (`#### Die zweite Richtung: Anforderung → Beleg` ist eine Überschrift). Bei `MR-069` ist die genannte Stelle *„§Ein einfrierendes Artefakt nennt ein prozess-bewegtes bei seiner Kennung"* **überhaupt keine Überschrift**, sondern ein fetter Absatz — ein Anker existiert nicht, und die Vorlage-Bedingung ist an dieser Stelle nicht erfüllbar, ohne das zu sagen. Runde 1 hat dieselbe Feld-Klasse bereits an `MR-067` beanstandet; die Form ist mit den zwei Neuzugängen nicht besser geworden. Kein Gate fängt das: ein Link ohne Anker ist für `links`/`anchors` gültig.
- **verifizierbar:** nein (kein Gate) — Vergleich mit der Vorlage und der Überschriften-Liste der Zieldatei.
- **klasse:** `pflichtfeld-formal-erfuellt-inhaltlich-nicht-adressiert`

### F-10 · MEDIUM · Der neue Schritt 4 kopiert den Kanon-Absatz nahezu wörtlich in `AGENTS.md`, unmarkiert und ungebunden

- **quelle:** `AGENTS.md` §1 (*„trägt Hard Rules und Pointer … und **dupliziert deren Inhalt nicht** — sonst entsteht Drift"*); `MR-039`
- **pfad:** `AGENTS.md:597-602`
- **befund:** Der Zusatz gibt den Kanon-Absatz (`v6.5.0` · `regelwerk/modul-09-implementierung.md` §Minimal Agent Workflow) in vier Zeilen fast Wort für Wort wieder — *„er erfindet die Abgrenzung nicht neu und weitet sie nicht stillschweigend … ist das eine Plan-Änderung und gehört vor den Code, nicht in den Bericht danach"*. Er ist **kein** ausgezeichnetes Zitat, trägt keine `d-check:cite`-Direktive und ist damit für `citations` unsichtbar: ändert die Baseline den Wortlaut, driftet die Kopie still. §1 verlangt an dieser Stelle einen Pointer statt einer Kopie; die vorhandene Präzedenz (die acht Workflow-Schritte selbst sind eine adoptierte Kopie) deckt die Schritt-Liste, nicht die begleitende Begründungsprosa. Der Inhalt ist korrekt wiedergegeben — geprüft, keine Über- oder Unterdehnung.
- **verifizierbar:** nein (kein Gate) — `citations` prüft nur ausgezeichnete Spannen.
- **klasse:** `kanon-prosa-unmarkiert-in-briefing-kopiert`

### F-11 · LOW · Der Audit sagt im Präsens, der Skill trage die Zitier-Form nicht — seit demselben Commit tut er es

- **quelle:** `AGENTS.md` §3.7 (Zustandsfeld nennt Zustand, nicht Chronik)
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:135-136` (§2, R6-Absatz) und die R6-Zeile der Audit-Tabelle
- **befund:** Der Absatz schließt *„Der Skill **trägt die Zitier-Form nicht**; sie gehört hinein."* Nach `8be0950` trägt er sie. Die Antwort-Spalte der R6-Zeile steht weiterhin auf *„übernommen, template-forward"* und nennt die Handlung nicht, während R4 und R7 als *„übernommen, mit Handlung"* geführt werden — drei Handlungen sind gefahren, die Tabelle weist zwei aus.
- **verifizierbar:** nein — Textvergleich Plan gegen Diff.
- **klasse:** `plan-aussage-nach-der-handlung-nicht-nachgezogen`

### F-12 · LOW · Index-Zeile und Eintrag widersprechen sich im Feld `Ersetzt-Baseline-Regel`

- **quelle:** `harness/conventions.md` §Adaptions-Block (die Index-Spalten stehen dort, *„damit ein Agent ohne Öffnen entscheiden kann"*)
- **pfad:** `harness/conventions.md:144` gegen `harness/conventions/MR-069-…md:4`
- **befund:** Die Index-Zeile führt in der Spalte *Ersetzt-Baseline-Regel* den Wert *„keine — Einlösung der Begründungslast aus …"*; die Eintrags-Datei füllt dasselbe Feld mit einem Link auf eine Baseline-Stelle und schreibt daneben *„keine Abweichung, sondern die **Einlösung**"*. Wer nur den Index liest, sieht ein leeres Feld; wer die Datei öffnet, einen gefüllten. Genau diese Entscheidung soll der Index ohne Öffnen tragen.
- **verifizierbar:** nein.
- **klasse:** `index-zeile-und-eintrag-fuellen-dasselbe-feld-verschieden`

### F-13 · LOW · Die beiden neuen Index-Zeilen tragen nur den Kurz-Anker, während `harness/conventions.md` von jeder Zeile einen Voll-Slug-Anker behauptet

- **quelle:** `harness/conventions.md` §Adaptions-Block (*„Je Index-Zeile steht ein Voll-Slug-`<a id>`"*)
- **pfad:** `harness/conventions.md:143-144`
- **befund:** `MR-068` und `MR-069` führen ausschließlich `<a id="mr-068">` bzw. `<a id="mr-069">`; jede andere Zeile — auch die unmittelbar davor stehende `MR-067` — trägt zusätzlich den Voll-Slug. Die Begründung im Fließtext (*Migrations-Schuld* für eingefrorene Verweise) trägt die Auslassung sachlich, aber der Satz darüber ist als Aussage über **jede** Zeile formuliert und stimmt damit nicht mehr.
- **verifizierbar:** nein.
- **klasse:** `selbstbeschreibung-der-datei-durch-neuzugang-falsch-geworden`

### F-14 · LOW · „Drei Regeln im Closure-Profil, die einen Abschnitt adressieren" — es sind fünf

- **quelle:** `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`
- **pfad:** `docs/plan/planning/in-progress/slice-208-v650-regel-adoption.md:241-244` (§5, letzter Risiko-Punkt)
- **befund:** Der Absatz stützt die Aussage *„Gate-seitig ist es weniger, als es aussieht"* auf eine Zählung: *„Von **drei** Regeln im Closure-Profil, die einen Abschnitt adressieren, keilt genau eine auf einen wörtlichen Titel."* `.d-check.closure.yml` trägt **fünf** `structure`-Regeln mit Abschnitts-Selektor (Closure-Notiz · DoD non-empty · H1-Platzhalter · Risiko-Ausgang · DoD-Haken). Die Schlussfolgerung bleibt richtig — die beiden nicht gezählten Regeln sind form-agnostisch (`^#{2,3} .*Closure-Notiz`, `^# `) —, die Zahl, mit der sie belegt wird, misst aber nicht ihren Gegenstand. Der Absatz stammt aus der Fassung vor Runde 1 und ist beim Nachziehen nicht mitgeprüft worden.
- **verifizierbar:** ja — Zählung der `section`/`section-pattern`-Schlüssel in `.d-check.closure.yml`.
- **klasse:** `zaehlung-misst-teilmenge-des-behaupteten-gegenstands`

### F-15 · INFO · Die Doppel-Form macht den Abschnitts-Selektor mehrdeutig, und der Kommentar sagt es nicht

- **quelle:** `DC-FA-STRUCT-001` (`sections: one` ⇒ `section-ambiguous` + Abbruch für diese Datei)
- **pfad:** `.d-check.closure.yml:185-192`
- **befund:** Gemessen: ein `done/`-Slice, der **beide** Titel trägt, meldet `section-ambiguous` und die `forbid-pattern`-Prüfung läuft für diese Datei **nicht** — der eingebaute Freitext-Ausgang blieb im Test unentdeckt (1 Befund statt 2). Das ist laut und damit kein stiller Grün-Pfad, aber es ist ein Zustand, den erst diese Änderung möglich macht: solange der Selektor wörtlich war, konnte er nicht zwei verschiedene Titel derselben Datei treffen. Der neue Kommentar über der Regel begründet die Doppel-Form ausführlich und erwähnt die Kardinalität nicht. Zur Einordnung mitgemessen: ein Titel, der **keine** der beiden Formen trifft, meldet weiterhin `section-missing` je Datei — die Blindheits-Sorge aus §5 des Plans ist damit auch in dieser Richtung ausgeräumt.
- **verifizierbar:** ja — Doppel-Überschrift einsetzen, `make verify-closure-notes`.
- **klasse:** `doppelform-selektor-erhoeht-mehrdeutigkeits-flaeche`

---

## Negativbefunde (geprüft, ohne Befund)

- **DoD (2), Bruch-Test je Richtung:** selbst gefahren, beide Richtungen reproduzieren — `## 5. Abnahme-Punkte / Risiken` und `## 6. Risiken und offene Punkte` melden je `section-forbidden` auf derselben Zeile (170). Die Behauptung der Commit-Botschaft trägt.
- **Falsch-Positiv-Richtung des neuen Musters:** geprüft. Das Muster ist auf `^…$` verankert und auf die Datei-Klasse `docs/plan/planning/done/slice-*.md` beschränkt; die offene Ziffer trifft auch `## 12. Risiken und offene Punkte`, was innerhalb dieser Klasse kein anderer Abschnittstyp ist. Kein konstruierbarer Fall, in dem das Muster einen Nicht-Risiko-Abschnitt greift, ohne dass die Datei bereits fehlerhaft nummeriert wäre.
- **`exempt-paths` der Regel:** unverändert (`slice-0??-*`, `slice-1[0-3]?-*`) und decken nach der Umstellung dieselbe Menge — die Ausnahme ist über Ziffern definiert, nicht über den Titel; keiner der ausgenommenen Slices wandert durch die Muster-Erweiterung in die Prüfmenge oder aus ihr heraus.
- **Kein stiller Grün-Pfad durch die Umstellung:** gemessen — ein Titel außerhalb beider Formen liefert `section-missing` je Datei (nicht nur bei leerer Kandidatenmenge).
- **Die übrigen Abschnitts-Regeln des Closure-Profils:** geprüft, alle form-agnostisch — `^#{2,3} .*Closure-Notiz`, `^## [0-9]+\. Definition of Done`, `^#{2,3} [0-9]+\. Definition of Done`, `^# `. Auch `planning.closure.heading-pattern` steht auf dem Default. Keine weitere Regel keilt auf einen Haus-Form-Titel.
- **`.d-check.yml` (Haupt-Profil):** geprüft — keine `structure`-Regel adressiert einen Slice-Abschnitt; die wörtlichen `section:`-Schlüssel dort betreffen Spec-Straten, ADR-Index, Roadmap, Handbuch und den Konventions-Index.
- **Zahlen von `MR-069`:** vollständig nachgemessen und exakt reproduziert (25 Einträge · zehn Tags · 28 Dateien · 18/6/3/1). Kein Mischfall (keine Datei verlinkt entfernten **und** gepinnten Baum).
- **Zitate:** beide Kanon-Zitate in `MR-068`/`MR-069` und das `DC-FA-CLI-011`-Zitat sind **wortgleich** gegen ihre Quelle geprüft.
- **`MR-068`s Kern-Aussage zur ADR-Spalte:** trägt. `DC-FA-CLI-011` §Out-of-Scope sagt es wörtlich; die Korrektur aus Runde 1 ist inhaltlich richtig eingearbeitet, und `.d-check.yml` konfiguriert `adrs` tatsächlich nur als Anzeigequelle.
- **R4 gegen `modul-09`:** die Regel ist **korrekt** wiedergegeben — weder mehr noch weniger als der Kanon sagt (Plan-Ausgabe nennt Out-of-Scope; Mitnahme ist Plan-Änderung und gehört vor den Code). Beanstandet ist nur die **Form** der Übernahme (F-10) und die Abschnitts-Adresse (F-2).
- **`harness/README.md` Schritt 4:** korrekt als Kurzform mit Pointer auf `AGENTS.md` §6 gebaut — keine zweite Quelle.
- **Die drei Formen der §Zitier-Form gegen den Kanon:** vollständig und richtig (Slice-Kennung · `make <target>`-Token · Tag + Pfad in Inline-Code). Beanstandet ist nur die Ausnahme-Hälfte (F-8).
- **Der Reviewer-Skill als lebendes Artefakt:** geprüft — seine eigenen Links in den vendorten Baum sind **kein** Verstoß gegen die neue Regel; der Kanon verlangt vom lebenden Artefakt ausdrücklich den Link. Die neue Sektion steht in §Ablage und adressiert erkennbar den Report, nicht sich selbst.
- **Lokale Slice-Vorlage:** existiert **nicht** — dieses Repo folgt direkt der vendorten `slice.template.md`. Die DoD-(2)-Hälfte *„die Vorlage folgt der Baseline-Form"* ist damit ohne Änderung erfüllt; sie ist allerdings auch nirgends als gemessen festgehalten.
- **slice-209 in `open/` trägt die Haus-Form:** geprüft und **kein** Defekt im Sinne von §3 — ein Plan in `open/` ist noch nicht beansprucht, und der Slice schließt ein Retrofit aus. Er ist aber der zweite Beleg für F-2: `AGENTS.md` §6 Schritt 4 verweist ab sofort auf ein „§1", das dieser Plan nicht hat.
- **Gate-Läufe:** `make gates` (zehn Gates, 693 Dateien, 0 Befunde) und `make verify-closure-notes` (592 Dateien, 0 Befunde) selbst gefahren und grün. Arbeitsbaum nach allen fünf Bruch-Tests wiederhergestellt, `git status` leer.
- **`AGENTS.md` §3.5/§3.6:** nicht berührt — keine `Accepted`-ADR angefasst, keine Schwelle gesenkt. Die einzige Gate-Berührung ist eine **Erweiterung** der Prüfmenge (F-1 beanstandet die Beschreibung einer bestehenden Senkung, nicht eine neue).

---

## Kategorie-Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 8 |
| LOW | 4 |
| INFO | 1 |
| **gesamt** | **15** |

---

## Verdikt

**Blockierend.** Zwei HIGH und acht MEDIUM.

Die **Gate-Arbeit selbst trägt** und ist der stärkste Teil dieses Slice: der Bruch-Test je Richtung reproduziert, die Nicht-Treffer-Richtung meldet weiterhin `section-missing`, die `exempt-paths` decken unverändert dieselbe Menge, und keine zweite Regel des Profils keilte auf einen Haus-Form-Titel. Auch die Zahlen von `MR-069` und die R2-Korrektur des Audits reproduzieren exakt — das ist mehr gemessene Substanz, als der Vorgänger-Stand hatte.

Blockierend sind zwei Dinge. **Erstens** beschreibt `MR-069` die Senkung, die er deklariert, breiter als sie ist; der Satz, der ihn trägt, ist durch einen Bruch-Test widerlegt und widerspricht `DC-FA-REF-001`. Ein Eintrag, dessen Zweck die ehrliche Benennung einer Gate-Senkung ist, darf sie nicht falsch bemessen. **Zweitens** ist die Migration nicht vollständig: die Regel im Profil kennt beide Formen, ihre Deklarations-Träger nicht — `AGENTS.md` §4, die Index-Zeile und der Eintrag zu `MR-049` nennen weiter „§5", `MR-053`/`MR-054` und `observations/README.md` einen „§Vorgelagert"/„§7", den die adoptierte Form nicht mehr hat, und der frisch geschriebene Schritt 4 verweist auf ein „§1", das kein lebender Plan dieses Repos führt. DoD (3) ist mit „Die Konventions-Einträge sind nachgezogen" abgehakt.

**Ist die Migration vollständig? Nein** — gate-seitig ja, deklarations-seitig nein. Genau das ist die Klasse, die der Plan in §7 als *„der zentrale"* Register-Eintrag benannt und deren Gegenmittel er sich selbst verordnet hat (*Spiegel vor dem Editieren auflisten*). Der Spiegel-Kreis wurde um die Konfiguration gezogen und nicht um ihre Deklaration.

Auffällig als Muster, nicht als einzelnes Finding: **drei** der Befunde (F-3, F-4, F-11) sind Stellen, an denen eine Korrektur aus Runde 1 an genau einer Fundstelle eingearbeitet wurde und ihre Geschwister stehen blieben. Das ist ein Steering-Loop-Signal in Richtung `BEO-ALL/semantic-change-body-only-edges-stale` und `BEO-ALL/eigene-menge-gemessen-fremde-behauptet`; beide sind im Plan gesichtet, und beide haben mit diesem Stand einen weiteren Beleg verdient.
