# Slice slice-190: Das Wellen-Archivierungs-Werkzeug bauen

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-87](welle-87-wellen-archivierung.md).

**Bezug:** [`modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 4](../../../../.harness/baseline/v5.18.0/regelwerk/modul-06-roadmap.md#wellen-closure-prozedur-modul-6)
(die Operation, die dieser Slice als Werkzeug einlöst);
[`archiv-stub-slice.template.md`](../../../../.harness/baseline/v5.18.0/templates/docs/plan/planning/archiv-stub-slice.template.md)
und
[`archiv-stub-welle.template.md`](../../../../.harness/baseline/v5.18.0/templates/docs/plan/planning/archiv-stub-welle.template.md)
(die Ziel-Form der beiden Stub-Arten);
[welle-87](welle-87-wellen-archivierung.md) (Ziel, Closure-Trigger,
Out-of-Scope-Abgrenzung zu diesem Slice).

**Berührte Spec-Stellen:** — Dieser Slice baut kein d-check-Regelmodul und
berührt keine `DC-FA-*`-Anforderung; er baut Planning-Infrastruktur, wie
`tools/harness/*.sh`.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-02.

---

## 1. Ziel

**Dieses Repo hat noch nie eine Welle archiviert, weil es kein Werkzeug dafür
gibt.** [slice-188](../done/slice-188-register-gegen-neuen-kanon.md) hat das
gemessen: kein `archiv.zip`, kein `done/<welle-id>/`-Verzeichnis existiert,
während welle-60 bis welle-85 bereits geschlossen sind und
[welle-86](../welle-86-closure-uebergang-durchsetzen.md) demnächst dieselbe
Pflicht trifft.

**Dieser Slice liefert genau das Werkzeug — nicht seine Anwendung auf den
echten Bestand.** Das ist [welle-87](welle-87-wellen-archivierung.md)s
eigene Trennung: erst bauen und an einem kontrollierten Fixture beweisen,
dann auf den Alt-Bestand loslassen (Folge-Slice).

## 2. Vorgehen

1. **Werkzeug-Form entscheiden, nicht stillschweigend wählen.** Alle
   bisherigen `tools/`-Skripte sind Bash
   ([`tools/harness/*.sh`](../../../../tools/harness)). Diese Operation
   braucht mehr: ZIP-Erzeugung, Markdown-Feld-Parsing (`**Welle:** <id>`),
   und einen repo-weiten, pfad-bewussten Verweis-Nachzug in **beiden** Formen
   (Verzeichnis-Präfix und geschwister-relativ). Ein Bash/`sed`-Ansatz dafür
   ist fehleranfällig (Sonderzeichen, Markdown-Verschachtelung) und liefert
   keine Testbarkeit. **Entscheid: ein Go-Programm**, gebaut via Docker
   (dieselbe Toolchain wie d-check selbst — `AGENTS.md` §3.1 verbietet Host-Go,
   nicht Go via Docker), mit `archive/zip` aus der Stdlib. Das ist **keine**
   vierte Toolchain neben `bash`/`go`/`docker`
   ([`MR-046`](../../../../harness/conventions.md#mr-046)) — es ist dieselbe,
   die dieses Repo schon führt. Ablage: tools/archive-wave/main.go, eigenes
   go.mod **oder** eingebunden ins bestehende Modul — zu entscheiden beim
   Bauen, je nachdem ob Kern-Pakete wiederverwendbar sind (z. B. die
   Markdown-Feld-Erkennung).
2. **Sammeln (nur explizite Wellen-Zugehörigkeit).** Für eine übergebene
   Wellen-ID: alle `docs/plan/planning/done/slice-*.md`, deren
   `**Welle:**`-Feld exakt diese ID nennt. **Wellenlose Alt-Slices bleiben
   außen vor** — ihre Zuordnung ist eine eigene Entscheidung
   ([welle-87](welle-87-wellen-archivierung.md) §3, Folge-Slice).
3. **Review-Reports zuordnen.** `docs/reviews/*.md`, deren Dateiname die
   Slice-Kennung (`slice-<NNN>`) eines gesammelten Slice trägt — 1:N zulässig
   (mehrere Reviews desselben Slice, z. B. `-r1`/`-r2`-Suffixe, gemessener
   Bestandsfall).
4. **`archiv.zip` bauen** (`done/<welle-id>/archiv.zip`): Welle-Plan-Volltext,
   alle gesammelten Slice-Volltexte, alle zugeordneten Review-Report-Volltexte
   — Pfade im Archiv identisch zu ihrer heutigen Repo-relativen Lage, damit
   der Stub-Zeiger `unzip -p done/<welle-id>/archiv.zip <pfad-im-archiv>`
   ohne Übersetzung stimmt.
5. **Stubs erzeugen** nach den beiden Templates, an `done/<welle-id>/` (nicht
   an der alten flachen Stelle — der Stub **wandert**, das ist die Ziel-Form).
   `Hervorgegangen:` bzw. `Archivierte Vorgänge:` sind **abgeleitet, wo
   möglich** (Anzahl der Stubs), sonst `<manuell auszufüllen>` — das Werkzeug
   erfindet keine Kennung, die es nicht sicher lesen kann.
6. **Verweise nachziehen.** Repo-weite Suche nach dem alten Pfad in beiden
   Formen, Ersetzung durch den neuen `done/<welle-id>/…`-Pfad — dieselbe
   Pflicht wie bei jedem Lifecycle-Move
   ([`MR-013`](../../../../harness/conventions.md#mr-013--lifecycle-move-commit-bündelt-gekoppelte-verweise)),
   hier über potenziell viele Dateien statt einer Handvoll.
7. **Sicherer Dry-Run als Default, `-apply` als Opt-in-Flag** (statt eines
   `--dry-run`-Opt-in-Flags — dieselbe Wirkung, umgekehrtes Vorzeichen: ohne
   `-apply` listet das Werkzeug die geplante Operation (welche Dateien,
   welche Verweis-Fixes) und schreibt nichts). Das ist der einzige sichere
   Weg, das Werkzeug gegen den echten Bestand zu **prüfen**, bevor ein
   Folge-Slice es **anwendet**.
8. **Verifikation an einem konstruierten Fixture**, nicht am echten Bestand:
   eine Test-Welle mit zwei Test-Slices und einem Test-Review-Report, in
   einem isolierten Testverzeichnis. Umkehr-Proben
   ([`BEO-023`](../observations.md)) für die drei Kernzusagen: Stub-Form
   (Sammeln-Bedingung umgekehrt ⇒ falscher Slice landet im Archiv), Link-
   Nachzug (ein Verweis, der nicht gefunden werden sollte, bricht die Probe),
   ZIP-Vollständigkeit (ein fehlender Vorgang im Archiv wird erkannt).
9. `make gates`; **Review**; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Anwendung auf den echten Alt-Bestand innerhalb DIESES Slices**
  (welle-60…welle-85) — ein neues, ungetestetes Werkzeug wird erst am Fixture
  bewiesen, nicht direkt am echten Baum ausprobiert. Das ist eine
  Reihenfolge-, keine Zeit-Entscheidung: der Folge-Slice
  ([welle-87](welle-87-wellen-archivierung.md) §4) läuft **im Anschluss,
  in derselben Sitzung** — Ziel ist ein vollständig archivierter Bestand am
  Ende des Durchgangs, nicht ein auf unbestimmte Zeit vertagter Rest.
- **Keine Zuordnung wellenloser Alt-Slices** zu einer archivierenden Welle in
  DIESEM Slice — dieselbe Abgrenzung, derselbe unmittelbare Folge-Slice.
- **`welle-86` wird nicht eingesammelt** — sie bleibt eigenständig
  ([welle-87](welle-87-wellen-archivierung.md) §5).
- **Kein neues `make gates`-Gate**, das die Archivierung erzwingt oder
  prüft. Dieser Slice liefert ein **Werkzeug**, keine Automatisierung im
  CI-Pfad — die Archivierung bleibt ein bewusster, von Hand ausgelöster
  Vorgang (der Kanon selbst verlangt das: „gehört in ein Werkzeug, nicht in
  Handarbeit" heißt *ein Werkzeug*, nicht *ein Zwang*).
- **Kein Umbau der bestehenden `tools/harness/*.sh`-Skripte** auf Go — die
  Entscheidung in Punkt 1 gilt dieser einen Operation, ist kein Präzedenzfall
  für alle `tools/`.

## 4. Definition of Done

- [x] Ein Go-Programm sammelt Slices nach `**Welle:**`-Feld, ordnet
      Review-Reports per Dateiname zu, baut `archiv.zip`, erzeugt beide
      Stub-Arten nach Template, zieht repo-weite Verweise in beiden
      Pfad-Formen nach.
- [x] Ohne `-apply` schreibt das Werkzeug nichts und listet die geplante
      Operation vollständig (Dateien, Verweis-Fixes), mit Test.
- [x] `make archive-wave WELLE=<id>`-Target dokumentiert (Handbuch-Klasse:
      internes Werkzeug, kein d-check-Produktfeature — Dokumentation gehört
      in dieses Slice-Verzeichnis bzw. eine kurze tools/archive-wave/README.md,
      **nicht** ins Benutzerhandbuch, das dem Produkt d-check gilt).
- [x] Verifiziert an einem konstruierten Fixture (Test-Welle, zwei
      Test-Slices, ein Test-Review-Report) — Stub-Form, ZIP-Inhalt,
      Verweis-Nachzug alle geprüft.
- [x] Umkehr-Proben ([`BEO-023`](../observations.md)): je Kernzusage eine
      Mutation, die genau den zugehörigen Test rot macht — Sammeln-Bedingung,
      Link-Nachzug, ZIP-Vollständigkeit.
- [x] `make gates` grün (Exit explizit); **unabhängiger Review**.
- [x] Closure-Notiz mit Lerneintrag; jedes Risiko aus §5 mit Ausgang; die drei
      Paarungen geprüft (bei der Closure von [welle-87](welle-87-wellen-archivierung.md),
      da dieser Slice ihr angehört).

## 5. Abnahme-Punkte / Risiken

- **Der Stub-vs-Volltext-Move ist doppelt destruktiv:** Inhalt wird ersetzt
  UND die Datei wandert in ein neues Verzeichnis, im selben Zug. Die
  Rename-Similarity sinkt damit unter 50 % — der Commit muss den Move
  ausdrücklich als `git mv` deklarieren (analog zur MR-/Wellen-Lifecycle-Move-
  Ausnahme in `AGENTS.md` §3.3, hier auf eine neue Kategorie angewandt, die
  dort noch nicht benannt ist). — **Ausgang: eingetreten → Folge-Slice
  slice-191** (Anwendung auf den echten Alt-Bestand, welle-87 §4) — dort
  entsteht der erste echte Stub-vs-Volltext-Move-Commit dieser Art, und die
  `git mv`-Deklarationspflicht gilt ihm, nicht diesem Slice (der nur das
  Werkzeug baute, ohne es gegen echte Dateien anzuwenden).
- **Review-Report-Zuordnung per Dateiname ist eine Heuristik.** Ein Review
  ohne erkennbare Slice-Kennung im Dateinamen (z. B. reine CR-Antwort-
  Dokumente) bleibt unzugeordnet und landet nicht im Archiv — das ist
  gewollt (nur Slice-Reviews gehören zum Slice-Vorgang), aber die Grenze
  gehört benannt und getestet. — **Ausgang: entfallen** — das Risiko war,
  dass diese Grenze unbenannt/ungetestet bliebe; sie ist jetzt beides:
  benannt in [`tools/archive-wave/README.md`](../../../../tools/archive-wave/README.md)
  §Grenzen, getestet in `collect_test.go:TestCollectReviews` (Fall
  `cr-ohne-slice.md`, unabhängig vom Review bestätigt).
- **Ein neues Go-Programm außerhalb des d-check-Kernmoduls ist ungetesteter
  Boden.** Die Fixture-Verifikation deckt den Entwurfsfall, nicht jede
  Eigenheit des echten, gewachsenen Bestands (uneinheitliche
  Review-Dateinamen, historische `**Welle:**`-Feld-Schreibweisen vor einer
  Konventions-Schärfung). — **Ausgang: eingetreten → Folge-Slice slice-191**
  — genau dessen Aufgabe ist, das Werkzeug gegen den echten Bestand zu
  fahren und dabei auftretende Eigenheiten zu behandeln.
- **Die Regel entsteht für einen Bestand, den dieser Slice nicht anfasst**
  ([`BEO-011`](../observations.md)): das Werkzeug wird gegen ein
  Fixture bewiesen, nicht gegen welle-60…85. Der Folge-Slice kann Eigenheiten
  finden, die hier nicht vorgedacht wurden. — **Ausgang: weiter offen** —
  als weitere Instanz von `BEO-011` registriert (Zähler jetzt 6, Beleg
  dieser Slice); ob das Muster (Regel aus dem Anlass statt aus der Klasse)
  hier tatsächlich zutraf oder durch die bewusste Fixture/Bestand-Trennung
  vermieden wurde, entscheidet sich erst an der Closure von slice-191.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei — `in-progress/` trägt
keinen Slice.

**Rückführungen:** `in-progress` → `next`, falls sich beim Bauen zeigt, dass
die Operation (Sammeln + ZIP + Stubs + Verweis-Nachzug) mehr als drei
Liefer-Punkte trägt und in zwei Werkzeug-Slices zerfällt (z. B. Sammeln+ZIP
getrennt vom Verweis-Nachzug) — dann ist der Schnitt ein anderer.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** `tools/` (das neue Werkzeug) und
  `docs/plan/planning/` (die Konvention, die es umsetzt, sowie die
  Stub-Templates als Ziel-Form). Beide fallen unter den Default `*` =
  **Greenfield** ([`harness/conventions.md`](../../../../harness/conventions.md)
  §Modus-Deklaration) — reine Neu-Infrastruktur, kein Bestand zum
  Nachziehen. Die Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:219-220 -->
  > **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
  > muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

- **Offene Beobachtungen sichten** (Register-Stand 2026-09-02, höchste
  Kennung `BEO-027`): [`BEO-011`](../observations.md) (Zähler 5) — die
  Regel entsteht für ein Fixture, nicht für den Anlass-Bestand; jede
  Verallgemeinerung über welle-60…85 hinaus, die dieser Slice trifft, muss
  aus dem Fixture-Test folgen, nicht aus der Absicht. Keine weiteren
  Treffer — Wellen-Archivierung als Klasse steht noch in keinem Eintrag. Die
  Regel, die diesen Schritt vorschreibt:

  <!-- d-check:cite .harness/baseline/v5.18.0/regelwerk/modul-05-planning-harness.md:225-225 -->
  > **Offene Beobachtungen sichten.**

- **Nachtlauf-Stand lesen** (`make nightly-state`,
  [`MR-053`](../../../../harness/conventions.md#mr-053)): `image-scan.yml`
  **gruen** (jüngster Lauf 2026-09-02T07:56:37Z). `upstream-drift.yml`
  **ROT** — jüngster Lauf 2026-09-02T05:19:44Z, planmäßig: der bekannte,
  informative Fremd-Release-Fund (Go 1.27.0→1.27.1, semgrep
  1.175.0→1.176.0), kein Zitat-Bruch, keine Regression. Ohne Konsequenz für
  diesen Slice — er berührt weder eine Toolchain-Version noch eine
  Zitat-Spanne außerhalb der beiden bereits in §7 geprüften.

Slice-ID: slice-190. Betroffene IDs: keine `DC-FA-*` (Planning-Infrastruktur,
kein d-check-Regelmodul). Module: keins. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — beide berührten Sub-Areas fallen unter
den Default: Doc führt, Code folgt. Ein neues Werkzeug ohne Vorgänger, keine
Reconciliation, kein umzustellender Bestand.

## 9. Closure-Notiz (nach `done/`)

**Geliefert:** `tools/archive-wave/` — eigenständiges Go-Programm mit
eigenem `go.mod`, eigenem `Dockerfile` (deps/test/build/runtime-Stages)
und eigenem `Makefile` (`test`/`build`/`run`/`help`), das Modul 6 Schritt 4
umsetzt: Sammeln nach `**Welle:**`-Feld, Review-Zuordnung per Dateiname,
`archiv.zip`-Bau, beide Stub-Arten nach Template, repo-weiter
Verweis-Nachzug. Sicherer Dry-Run-Default (`-apply`-Opt-in), Mount nur bei
`APPLY=1` beschreibbar. Das Wurzel-Makefile delegiert
(`make archive-wave`/`archive-wave-test`) statt den Docker-Aufruf zu
duplizieren. Verifiziert an einem konstruierten Fixture (`TestFixture_EndToEnd`)
plus drei benannten Umkehr-Proben und einem Dry-Run-Baum-Snapshot-Test.
Ein Smoke-Test gegen den echten Bestand (`welle-85`, Dry-Run, `:ro`-Mount)
lief bereits erfolgreich.

**Was funktioniert hat:** Die Trennung „Werkzeug bauen, nicht anwenden"
(§3) hat den Slice klein gehalten und den Fixture-Beweis vom
Anlass-Bestand ferngehalten — genau die BEO-011-Vorsicht, die §7 vorab
benannte. Die unabhängige Review- und Verifikations-Sequenz (Modul 8) fing
zwei HIGH-Befunde (Slice-Nummer-Litter in Kommentaren gegen AGENTS.md §3.7;
fehlender Test für den Dry-Run-Pfad), die beim Schreiben unbemerkt blieben.

**Was anders lief:** Drei Nutzer-Korrekturen während des Baus haben die
Architektur verbessert, bevor sie ins Review ging: (1) Go-Version blind aus
einer anderen Datei kopiert statt aus dem gepinnten `GO_VERSION` gelesen —
korrigiert auf 1.27.0. (2) Das Werkzeug sollte zunächst d-checks interne
Kern-Pakete wiederverwenden — der Hinweis „andere Repos brauchen das auch"
kehrte das um: eigenständiges Modul, self-contained Pfad-Auflösung. (3) Der
Docker-Aufruf stand zunächst roh in README und Wurzel-Makefile dupliziert,
mit einem durchgehend beschreibbaren Voll-Mount auch im Dry-Run — ein
eigenes lokales `Makefile` kapselt ihn jetzt einmalig, und der Mount ist
nur bei `APPLY=1` beschreibbar.

**Steering-Loop-Einträge:**
- **Geschärfte Prozedur (verkörpert in diesem Slice-Plan, keine gesonderte
  Regel-Datei):** ein Docker-Aufruf, der in mehr als einem Makefile/einer
  README vorkäme, gehört hinter EIN Target, nicht dupliziert — dieselbe
  Drift-Vermeidung wie AGENTS.md §1 sie für Doku-Inhalt schon fordert, hier
  erstmals auf einen Docker-Aufruf angewandt.
- [`BEO-011`](../observations.md) — sechste Instanz registriert (Zähler 6),
  Ausgang offen bis zur Closure von slice-191.
- Zwei bereits im Register geführte Beobachtungen wurden während der
  Vorprüfung (§7) gesichtet, keine neue Instanz ausgelöst.

**Zeiger:** [Beobachtungs-Register](../observations.md). Folge-Slice:
[slice-191](../done/slice-191-alt-bestand-archivieren.md) — Anwendung des
Werkzeugs auf welle-60…welle-85 und Zuordnung der wellenlosen Alt-Slices,
welle-87 §4.
