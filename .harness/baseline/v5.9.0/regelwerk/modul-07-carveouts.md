## Modul 7 — Carveout Management

<!-- Quelle: [02-planung/modul-07-carveouts.md](https://github.com/pt9912/ai-harness-course/blob/v5.9.0/kurs/de/02-planung/modul-07-carveouts.md) -->

### Harness-Einordnung (Modul 7)

Carveout-Pflege ist ein Pfeiler von *Entropy Management* (siehe
[`klassifikation.md`](grundlagen-klassifikation.md)):
ein Carveout-Audit pro Welle verhindert, dass temporäre Ausnahmen zu
permanenten Lügen werden.

### Kernidee (Modul 7)

Jeder temporäre Carveout benötigt einen Plan. Ein Carveout ohne
Auflösungs-Trigger ist ein permanenter Carveout, der lügt.

### Ziel-Form: Carveout

Ein temporärer Carveout folgt der Vorlage
[`templates/docs/plan/carveouts/carveout.template.md`](../templates/docs/plan/carveouts/carveout.template.md)
(Dateikonvention `docs/plan/carveouts/CO-<NNN>-<kurztitel>.md`; ID in
eigener `CO-*`-Reihe, siehe
[`grundlagen-source-precedence.md`](grundlagen-source-precedence.md#id-schema-als-klammer)).
Operative Regeln, die das Template nicht selbst erzwingt:

- **Sechs Pflicht-Header-Felder:** Status · Datum angelegt · Letzte
  Prüfung · betroffenes Gate · Geltungsbereich · **Folge-Slice**. Fehlt
  der Folge-Slice, ist der Carveout *de facto* permanent — dann gehört
  er nicht in `carveouts/`, sondern über den Trichter unten in eine ADR.
- **Auflösungs-Trigger als beobachtbare, messbare Bedingung** — nicht
  „sobald wir Zeit haben", sondern eine Schwelle, die ein anderer Mensch
  ohne Rückfrage als erreicht beurteilen kann (z. B. „`internal/parser/`-
  Coverage ≥ 90 %, geprüft in `make coverage-gate-critical` ohne
  Ausnahmen"). Der Welle-Bezug ist der Roadmap-Anker, die messbare
  Schwelle das, was die CI prüft.
- **Die Gate-Konfiguration nennt die `CO-<NNN>` im Gate-Output** — sonst ist
  die Ausnahme eine stille Senkung ohne Begründung. Ob sie als
  Ausschluss-Liste oder als Verengung der Messung umgesetzt ist, hängt am
  Werkzeug; sichtbar sein muss sie in jedem Fall.
- **Auflösung ist ein `git mv` nach `done/`** (plus Gate-Ausnahme
  zurücknehmen, Gate grün ohne Ausnahmen). Auflösen ohne
  Verschiebung ist eine zweite Lüge: der Carveout wirkt „aufgelöst",
  liegt aber weiter im aktiven Verzeichnis.

<a id="werkzeug-wahl"></a>

### Werkzeug-Wahl bei Diskrepanz: Carveout, BF-Markierung oder ADR (Modul 7)

Bevor eine Diskrepanz als Carveout festgeschrieben wird, prüfe, ob das
Werkzeug *Carveout* überhaupt passt. Drei legitime Werkzeuge, getrennt
durch **zwei sequenzielle Wenn-Dann-Fragen** — Granularität *vor*
Temporalität (das verhindert den Reflex, jede Diskrepanz als Carveout zu
führen):

1. **Granularität — einzelne Diskrepanz oder Cluster?** *Cluster im
   selben Geltungsbereich* (mehrere Ausnahmen auf denselben Pfad/dieselbe
   Sub-Area) oder systemisches *„Code existiert vor Doku"*-Muster →
   **BF-Sub-Area-Markierung mit Graduation-Plan** als Modus-Deklaration
   im Adaptions-Block von `harness/conventions.md` (Mechanik in
   [`grundlagen-bootstrap.md` §Modus pro Sub-Area](grundlagen-bootstrap.md#modus-pro-sub-area-greenfield-vs-brownfield));
   Frage 2 entfällt. *Einzelne Diskrepanz* → Frage 2. **Kein harter
   Schwellwert** für „Cluster" — Faustregel (gemeinsamer Geltungsbereich),
   keine Carveout-Zahl.
2. **Temporalität — Trigger ernst zu erreichen?** *Ja* (absehbarer
   Aufwand, sinnvolles Verhältnis zum Nutzen) → **Carveout** (Ziel-Form
   oben). *Nein* („nichts davon werden wir in absehbarer Zeit tun") →
   **permanent**, übergeführt in eine **ADR** (`Status: Permanent —
   übergeführt in ADR-<NNNN>`).

| Wahl | Symptom-Indikator | Träger | Folge-Artefakt |
|---|---|---|---|
| **Carveout** | Eine konkrete Gate-/Regelausnahme, klar abgrenzbar, mit Folge-Slice und ernst erreichbarem Trigger. | einzelne Diskrepanz | `docs/plan/carveouts/CO-<NNN>-*.md` |
| **BF-Sub-Area-Markierung** | Diskrepanz-Cluster im selben Geltungsbereich, oder generelles *„Code-vor-Doku"*-Muster. | ganze Sub-Area | Modus-Deklaration im Adaptions-Block von `harness/conventions.md`, mit Graduation-Trigger |
| **ADR (permanent)** | Trigger ist ehrlich nie zu erreichen — die Senkung ist Architekturentscheidung, kein Übergang. | dauerhafte Architekturregel | `docs/architecture/ADR-<NNNN>-*.md` |

**Bootstrap-aware Gate** (Modul 13) erscheint hier absichtlich nicht: es
regelt *Gate-Reifestufung des Gates selbst* (z. B. 40 % heute → 70 % bei
M2), nicht *Diskrepanz-Auflösung*. BF-Markierung wirkt eine Ebene höher
als Carveout und ADR — sie kippt den Sub-Area-Kontext, in dem die
Diskrepanz erst entsteht (*Sub-Area-Kontext, kein Closure-Werkzeug*,
siehe [Modul 13 §Bootstrap-aware Gates](modul-13-quality-gates.md#bootstrap-aware-gates)).
Verwechslung der Achsen führt zum „Bootstrap-Schlupfloch": Stufung ohne
Trigger ist Carveout-Wildwuchs, Carveout-Kaskade ohne BF-Markierung ist
verschleierte Sub-Area-BF.

Führt der Trichter *nicht* auf Carveout, ist der Entwurf nicht verloren:
Geltungsbereich, Trigger und Verifikations-Checkliste wandern ins
gewählte Werkzeug (BF-Markierung: Trigger → Graduation-Trigger,
Checkliste → Reconciliation-Akzeptanzkriterien; ADR: Trigger fällt weg,
Checkliste reduziert auf die Architektur-Folgen). Der leere
`CO-<NNN>`-Stub wird gelöscht (Inhalt ganz aufgegangen) oder mit
`Status: Überführt in <Ziel>` nach `done/` verschoben, damit die
Werkzeug-Wahl-Spur im Repo lesbar bleibt. Das ist nicht Aufgabe, das ist
Ehrlichkeit.

**Eine Diskrepanz-Klasse liegt außerhalb dieses Trichters.** Alle drei
Werkzeuge antworten auf ein Missverhältnis zwischen *Code und Doku*. Ein Fund
aus der Stichprobe des Freshness-Audits — eine Baseline-Regel, die nie ins
ausgefüllte Artefakt übernommen wurde — liegt zwischen *adoptierter Norm und
Artefakt*. Punktuell behandelt der Trichter ihn richtig: Übernahme im nächsten
Slice, oder Carveout mit Auflösungs-Trigger. Bei einer **Häufung** greift die
Sub-Area-Antwort dagegen nicht — eine BF-Markierung setzte *Code führt, Doku
folgt*, und darum geht es nicht. Mehrere Funde treffen die `MR-000`-Aussage
selbst; sie wird korrigiert, nicht die Sub-Area umgestuft
([`modul-02-harness-bootstrap.md` §Freshness-Audit](modul-02-harness-bootstrap.md#freshness-audit-der-vendored-baseline-schritt-2)).

### Carveout-Audit-Slice (Modul 7)

Der Carveout-Mechanismus hält nur, wenn ein *zweiter* Mechanismus ihn
auditiert — sonst driften „aktive" Carveouts still in De-facto-Permanenz
(genau die Doku-Drift, die Carveouts verhindern sollen). Pro
Welle-Closure ein Audit-Slice `slice-CO-AUDIT-<welle>`, *bevor* die Welle
nach `done/` wandert (eigenes Präfix laut
[`grundlagen-source-precedence.md` §ID-Schema](grundlagen-source-precedence.md#id-schema-als-klammer) — liefert keinen Code,
nur Doku-Updates). Regeln:

- **DoD (vier Punkte, drei Status-Aktionen + ein Beleg):** jeder aktive
  Carveout trägt ein aktuelles `Letzte Prüfung:`-Datum (≤ heute); jeder
  Carveout mit eingetretenem Trigger ist nach `done/` verschoben; jeder
  seit > 2 Wellen „aktive" Carveout wird *explizit* als weiter-gültig
  bestätigt **oder** in eine ADR überführt; Audit-Bericht als Block in
  `done/welle-NN-results.md`.
- **Drei Status-Übergänge** je Carveout: *aufgelöst* (Trigger eingetreten
  → `git mv` nach `done/`), *permanent* (Trigger nie → ADR),
  *weiterhin aktiv* (Trigger sinnvoll → `Letzte Prüfung:`-Datum
  nachtragen, ggf. Folge-Slice).
- **Rollen (Modul 8):** *Planner* identifiziert die fälligen Carveouts,
  *Architect* entscheidet bei „permanent" über die ADR-Überführung,
  *Implementer* führt `git mv` und Config-Updates aus. Verteilung über
  drei Rollen ist Absicht, kein Defekt — Aufräumen ohne Architect-Blick
  verlängert das Lügen.
- **Freshness-Gate (optional, maschinell):** ein Carveout, dessen letzte
  Prüfung > 90 Tage zurückliegt, ist ein HIGH-Warnsignal — egal, ob der
  Trigger nominell noch gilt. *Beobachtung schlägt Behauptung*: ein nicht
  geprüfter Carveout ist ein nicht existierender Audit.

### Regeln gegen typische Fehlannahmen (Modul 7)

- **Gegen "Carveout = Workaround":** Carveout = *dokumentierter* Workaround mit Trigger. Ohne Trigger ist es eine versteckte Annahme.
- **Gegen "Carveouts gehören ins Issue-Tracker":** Sie gehören ins Repo, neben Spec und ADRs. Tracker können vergessen werden, Repo-Files kommen mit beim Klonen.
- **Gegen "Wenn der Trigger eintritt, lösen wir den Carveout auf":** Realität: er bleibt liegen. Deshalb braucht jeder temporäre Carveout einen *Folge-Slice mit ID*, der das Auflösen plant. Slice schlägt Memo.
- **Gegen "Jede entdeckte Diskrepanz ist ein eigener Carveout":** Carveouts sind für **punktuelle** Ausnahmen mit Folge-Slice. Eine Diskrepanz-**Häufung** in einer Sub-Area (Symptom: mehrere Carveouts mit demselben Geltungsbereich, oder die Diskrepanz folgt aus generellem *"Code existiert vor Doku"*-Muster) gehört nicht in eine Carveout-Kaskade, sondern in eine **BF-Sub-Area-Markierung mit Graduation-Plan** (siehe [Modul 2 §Kernidee](modul-02-harness-bootstrap.md#kernidee-modul-2)). Maßgeblich ist das **Symptom-Muster** (gemeinsamer Geltungsbereich), nicht die Carveout-Zahl; die Wahl, welches Werkzeug bei welchem Symptom greift, leistet [§Werkzeug-Wahl bei Diskrepanz](#werkzeug-wahl).
- **Gegen "Wenn Diskrepanz-Häufung BF-Markierung verlangt, ist auch jede einzelne Diskrepanz eine BF-Markierung wert":** BF-Markierung lohnt sich erst beim **Cluster im selben Geltungsbereich** oder beim systemischen *"Code existiert vor Doku"*-Muster — eine einzelne, gut abgrenzbare Diskrepanz mit klarem Folge-Slice ist und bleibt ein Carveout. Das Frage-Schema in [§Werkzeug-Wahl bei Diskrepanz](#werkzeug-wahl) trennt diese Fälle: Frage 1 leitet einzelne Diskrepanzen explizit auf den Carveout-/ADR-Pfad, nicht auf BF-Markierung.

