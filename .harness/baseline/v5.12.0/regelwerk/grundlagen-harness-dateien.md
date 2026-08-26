## Die Harness-Dateien und ihre Form
<!-- Quelle: [grundlagen/harness-dateien.md](https://github.com/pt9912/ai-harness-course/blob/v5.12.0/kurs/de/grundlagen/harness-dateien.md) -->

### Verzeichniskonvention

```
spec/                       # Spec-Straten: Vertrag · Technik · Sicht
docs/plan/adr/              # Architecture Decision Records
docs/plan/planning/open/    # geplante, noch nicht gestartete Slices
docs/plan/planning/next/    # priorisiert/eingeplant
docs/plan/planning/in-progress/  # aktive Slices
docs/plan/planning/done/    # abgeschlossene Slices
docs/plan/planning/<welle-id>.md            # offene Wellen, flach (Modul 6)
docs/plan/planning/observations.md          # Beobachtungs-Register: der Steering-Loop-Zähler
docs/plan/planning/reconciliation.md        # Reconciliation-Register: nur im Brownfield-Bootstrap
docs/plan/planning/in-progress/roadmap.md   # Meilensteine, nächste Wellen, Zeiger auf offene
docs/plan/carveouts/        # Ausnahmen mit Plan zur Auflösung
docs/reviews/               # Review-Reports, ein Report pro Lauf (Modul 10)
AGENTS.md                   # maschinell lesbare Projekt-Konventionen für Agenten
harness/README.md           # Einstiegspunkt: Precedence, Guides, Sensors, Safety
harness/conventions.md      # Index: repo-lokale Regeln, Adaptionen, Modus pro Sub-Area
harness/conventions/        # ein MR je Datei; done/ = aufgelöst
.harness/                   # Skills, Tool-Allowlists, Checklisten-Middlewares
```

### Template-Schichtung — was der Rumpf trägt und was der Kommentar

Ein Template wird beim Adoptieren **abgebaut**: Platzhalter ersetzt,
Hinweis-Block entfernt, **alle HTML-Kommentare gelöscht** — bis auf die
`d-check:ignore`-Marker, die Falsch-Positive unterdrücken und bleiben müssen
([`../templates/README.md`](../templates/README.md) §Verwendung, Schritt 5). Was danach
dasteht, ist alles, was der Adopter Wochen später hat. Vier Schichten:

| Schicht              | Inhalt                                                                                                                                       | Überlebt das Adoptieren?                                                            |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| **Regelwerk**        | Der Normtext. **Einzige** Quelle.                                                                                                            | — vendored unter `.harness/baseline/<tag>/regelwerk/`, lebt außerhalb des Artefakts |
| **Rumpf**            | Nur, was das *fertige Artefakt* trägt: Feldnamen, Feldreihenfolge, `<Platzhalter>` — plus **genau ein** Regelwerk-Zeiger pro Pflicht-Sektion | ja                                                                                  |
| **DoD / Checkliste** | Jede Pflicht, die der Ausfüllende **abhaken** muss. Das ist die Prozedur                                                                     | ja                                                                                  |
| **Kommentar**        | Begründung und Bedienhinweis                                                                                                                 | nein                                                                                |

- **Test für den Rumpf:** Liest sich das im veröffentlichten Artefakt als
  *Inhalt* — oder als *Anleitung an jemanden*? Anleitung gehört nie in den
  Rumpf. Dazu kommt: Normtext im Rumpf wird vom Platzhalter-Ersetzen
  zerschossen — aus der allgemeinen Regel wird eine falsche Einzelaussage.
- **Hard Rule:** *Kein Kommentar ist die einzige Fundstelle einer Norm.* Wer
  eine Regel in einen Template-Kommentar schreibt, schreibt sie in den
  Papierkorb des Adopters. Sie gehört ins Regelwerk; im Template steht der
  **Zeiger** darauf, im Rumpf, bei der Sektion, für die sie gilt.
- **Der Zeiger ist kein Zitat.** Ein Template, das den Normtext ausschreibt,
  führt ihn ein zweites Mal — und zwei Fassungen driften.
- **Feedback-Hälfte ist *inferential*, nicht computational:** „Ist dieser Satz
  eine Norm?" ist ein Urteil, kein Match — und Template-Verzeichnisse sind für
  Referenz-Gates bewusst ausgenommen (symbolische Pfade). Die Regel steht
  deshalb als HIGH-Eintrag *Norm nur im Template-Kommentar* im Reviewer-Skill
  (Ziel-Form `../templates/.harness/skills/reviewer.template.md`).
- **Grenze:** Sie hängt damit an einem Review, nicht an einem Lauf. Wer ohne
  Review committet, wird nicht erwischt. Und der Skill oben ist eine
  **Ziel-Form für das adoptierende Repo** — ob dort ein Review mit dieser
  HIGH-Regel tatsächlich läuft, entscheidet der Adopter, nicht diese
  Konvention. Wo er es nicht einrichtet, hat die Hard Rule keinen Träger, und
  das ist kein Sonderfall: Es ist der Auslieferungszustand. Einen *Sensor* zu
  behaupten, wo keiner steht, wäre die Klasse *halluziniertes Gate* (Modul 13)
  — auf die eigene Konvention angewandt.

### Was ein Kommentar trägt — Code, Konfiguration, Skripte

Der Abschnitt oben regelt Kommentare in **Templates**; die werden beim
Adoptieren gelöscht. Ein Kommentar in Code, Konfiguration oder Skript bleibt,
solange die Zeile lebt. Seine Bestimmung ist **positiv**: nicht, was verboten
ist, sondern was er zu tragen hat.

- **Leser-Modell:** Ein Kommentar schreibt an jemanden, der den Code **gleich
  ändert** — nicht an jemanden, der die Entscheidung noch einmal treffen will.
  Der zweite Leser hat die ADR. Der erste beginnt bei Null
  ([`modul-03-spec.md` §Kernidee](modul-03-spec.md#kernidee-modul-3)): Was nicht
  im Kontext steht, war für ihn nie da — und was daneben steht, liest er
  trotzdem mit und bezahlt es mit Kontext.
- **Fünf Klassen.** Ein Kommentar beantwortet, was der Code nicht beantworten
  kann:

| Klasse | Die Frage, die der Code offen lässt |
|---|---|
| **Zusage** | Was garantiert diese Stelle — und was müsste passieren, damit sie bricht? Mit dem Sensor, der es sähe. |
| **Kopplung** | Was muss ich mitändern, wenn ich das hier ändere? |
| **Abgrenzung** | Welche Nachbargröße verwechsle ich hiermit, und warum ist sie es nicht? Sprach- und Plattform-Fallstricke gehören hierher. |
| **Rang-Zeiger** | Wo wohnt die Norm, deren *Umsetzung* das hier ist? |
| **Grenze** | Was leistet diese Stelle ausdrücklich **nicht**? |

- **Die Liste ist abschließend gemeint, nicht abschließend bewiesen.** Wer eine
  sechste Klasse findet, die keine der fünf trägt, erweitert sie; das ist der
  reguläre Weg, nicht die Ausnahme.
- **Adressaten-Test:** Richtet sich der Satz an den, der *ändert*, oder an den,
  der *entscheidet*? Der Entscheider hat ADR und Slice.
- **Zeitform-Test:** Steht der Satz im Indikativ über das, was ist? Konjunktiv
  ist zulässig über den **Bruch** (*„was müsste passieren, damit die Zusage
  fällt"*) und über **künftige Arbeit** (*„das wäre ein eigener Slice"* — die
  Grenze-Klasse). Unzulässig ist er über die **verworfene Alternative**: Die
  ist entschieden, und die Entscheidung steht in der ADR. Die Probe dafür ist
  die Zeitrichtung — zeigt der Konjunktiv nach vorn oder zurück?
- **Hard Rule:** *Ein Kommentar beschreibt, was da ist.* Wer Herkunft nennt,
  nennt sie als **ein** auflösbares Feld — `LH-*`, `ADR-*`, `· seit welle-<NN>`
  ([`grundlagen-traceability.md` §Herkunfts-Anker](grundlagen-traceability.md#herkunfts-anker))
  — und nie als Absatz.
- **Dieselbe Regel für Zustandsfelder.** Ein Feld, das einen *Zustand* trägt —
  `Stand` oder `Status` einer Register-Zeile, die Stand-Zelle eines
  Roadmap-Fadens, die Status-Spalte eines Meilensteins —, ist ein
  Zustands-Artefakt wie der Kommentar, nur im Rumpf. Es nennt den Zustand und
  den Beleg — als auflösbaren Anker (Welle, Commit, Register-Zeile,
  Messdatei), nicht als Chronik, wie der Zustand zustande kam; eine Begründung
  des Zustands (*gestrichen — tritt nicht mehr auf, weil …*) ist Zustand,
  nicht Chronik. Was sonst in der Zelle stand, hat seine Orte: Behauptung und
  vorgeschlagene Handlung beim Vorhaben selbst, die Schließung im Closure-Log
  (im Wellen-Betrieb *Abgeschlossene Wellen*; sonst dort, wo das Repo
  Schließungen führt — die Slice-Datei in `done/` und `git`, ohne
  Slice-Lifecycle etwa ein `CHANGELOG`), die Umplanung im Drift-Log, und was
  keines davon ist, hält `git`. Ein zweites Log neben dem Closure-Log — ein Drift-Log, das
  Schließungen und erreichte Meilensteine protokolliert — ist eine Kopie, und
  Kopien driften ([`grundlagen-source-precedence.md` §Source Precedence](grundlagen-source-precedence.md#source-precedence)).
  Die zwei Tests oben gelten unverändert: Adressat ist, wer den Zustand
  liest, um zu handeln — nicht, wer die Geschichte hören will; die Zeitform
  ist der Indikativ über das, was ist. Wie das in der Roadmap aussieht, sagt
  [`modul-06-roadmap.md` §Roadmap-Struktur](modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6).
- **Die Kopfzeile eines lebenden Registers ist derselbe Fall.** Roadmap,
  Beobachtungs- und Reconciliation-Register tragen keine Zeile
  `Status: Aktiv. Letzte Änderung: <Datum>`: *Aktiv* ist kein Zustand, den
  ein Register je wechselt, und ein Datum, das niemand pflegt, behauptet
  einen. Der Zustand eines Registers ist sein Inhalt (die Zeilen mit ihren
  Belegen, das Drift-Log mit seinen Daten), sein Änderungsdatum hält `git`.
  Ein Datum, das ein benannter Trigger pflegt, ist kein solcher Kopf — der
  Unterschied ist der Trigger, nicht die Zeile: Der Wellen-Stand des
  Regelwerks wird mit jeder Welle gesetzt, und im Sicht-Stratum der Spec ist
  `Letzte Änderung` der bewusste **Frische-Marker** ([`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](modul-03-spec.md#ziel-form-architektur-sicht))
  — die Aussage, wann die Sicht zuletzt gegen den Code gehalten wurde, und
  kein Inhalt darunter trägt sie.
- **Emittierte Artefakte tragen keinen Anker.** Erzeugt ein Werkzeug Code,
  Konfiguration oder Skripte *in ein anderes Repo*, reist der Erzeuger-Kontext
  nicht mit: Eine Slice- oder Befund-Nummer des Erzeugers löst dort in **null**
  Hops auf und ist dann kein Anker, sondern ein toter Verweis. Dieselbe Regel
  hängt beim Regelwerk-Split die Deixis um — *keine Verweise auf Material, das
  nicht mitreist*.
- **Drei Klassen fallen heraus, weil sie keine der fünf tragen:**
  **Deliberation** (der Konjunktiv über die verworfene Alternative — *„Ohne
  dieses Feld behauptete die Ausgabe …"*; ihr Ort ist die ADR oder §3/§6 des
  Slice), **Herkunfts-Prosa** (beschreibt abwesenden Text oder Code — *„die
  frühere Zusage wurde ersetzt"*; ihr Ort ist `git`) und
  **Ersetzungs-Trümmer** (eine Teilersetzung lässt den Rest des alten Satzes
  stehen, der Kommentar bricht mitten im Satz ab; sie gehört nirgendwohin).
- **Feedback-Hälfte:** Von den dreien ist genau eine ein Match —
  **Ersetzungs-Trümmer** ist syntaktisch erkennbar; gebaut ist der Sensor dafür
  nicht. **Deliberation** und **Herkunfts-Prosa** sind Urteile. Träger aller
  drei ist das Briefing (Ziel-Form
  `../templates/AGENTS.template.md` §3) plus der HIGH-Eintrag
  *Kommentar trägt keine der fünf Klassen* im Reviewer-Skill (Ziel-Form
  `../templates/.harness/skills/reviewer.template.md`). Für Zustandsfelder
  ebenso: *„ist das eine Chronik?"* ist ein Urteil — Träger sind das Briefing
  (§3.7) und der HIGH-Eintrag *Zustandsfeld trägt Chronik* im Reviewer-Skill.
- **Grenze:** Einen *Sensor* zu behaupten, wo keiner steht, wäre die Klasse
  *halluziniertes Gate* (Modul 13). Gemessen: Auf ausdrückliche Anweisung
  korrigierte ein Agent sieben Kommentare — fünf sauber, einen reproduzierte er
  im selben Diff, einen ließ er als Trümmer stehen.

### harness/README.md als Einstiegspunkt

Pro Repo bündelt eine einzige Datei alles, was ein Agent oder ein neuer
Mensch zuerst lesen muss. Pflichtgliederung:

```
# Harness

## Purpose                  # ein Absatz, was diese Datei ist (und was nicht)
## Source precedence        # die obige Tabelle, repo-spezifisch
## Guides                   # Tabelle der Feedforward-Quellen
## Sensors                  # Tabelle der Feedback-Gates (nur real existierende!)
## Traceability rules       # Welche IDs müssen in Commits/PRs auftauchen?
## Safety and scope boundaries  # repo-spezifische Hard Rules
## Minimal agent workflow   # der 8-Schritt-Pfad (siehe Modul 9)
## Leseordnung              # für den neuen Menschen: was zuerst, was bei Bedarf
```

**Die Leseordnung ist die Menschen-Hälfte des Einstiegs.** Die sieben
Sektionen darüber sind eine Referenzfläche — richtig für den Lauf, der
nachschlägt; ein neuer Mensch braucht eine **Reihenfolge**: was zuerst, was
bei Bedarf. Drei bis fünf geordnete Zeiger genügen; eine Leseordnung, die
alles nennt, ist keine.

Wichtig: Die Sensors-Tabelle darf keine Befehle behaupten, die es im Repo
nicht gibt. Halluzinierte Gates sind die häufigste Form von Harness-Lüge.

Die Sensors-Tabelle trägt **keinen Lauf-Status** ("grün"/"rot"):
Lauf-Wahrheit pro Commit lebt in CI (Badges/Dashboard), also in höher
rangierten Quellen, nicht in `harness/README.md` (unterster Rang). Strukturell
rote Gates werden als Carveout in `docs/plan/carveouts/` dokumentiert
(Modul 7); die Bindung-Spalte der Tabelle (`Target | Vertrag | Bindung`)
verweist auf die `CO-<NNN>`-ID, die Begründung lebt im Carveout, nicht
hier. Damit ist "rot dokumentieren, nicht verstecken" ortsdiszipliniert:
es geschieht im Carveout-Index, nicht in einer Status-Spalte, die sich
selbst grünfärben kann.

Die Bindung-Spalte trägt vier **kanonische Klassen**:

- **ADR-Bindung** (`ADR-<NNNN>`) — Gate setzt eine Architektur-Entscheidung
  durch.
- **Carveout-Bindung** (`CO-<NNN>`) — Gate bewusst geschwächt, mit
  Auflösungs-Trigger und Folge-Slice (Modul 7).
- **Kalibrierungs-Bindung** (`Schwelle X %, M<n> → Y %`) — bewegliche
  Eichung mit Meilenstein-Schaltplan.
- **Reproduzierbarkeits-Bindung** (Image-Hash, Toolchain-Pin) — Gate
  hängt an bit-identischem Artefakt (Modul 14).

Repos können **weitere Klassen** einführen — etwa Anforderungs-Bindung
(`LH-…`), Compliance-Bindung (Regulatorik-Artikel) oder
Modell-Version-Bindung (für KI-Evals). Diese werden im **repo-lokalen
Konventionsdokument** deklariert (Default-Pfad `harness/conventions.md`,
Form projektabhängig), damit ein Reviewer sie als legitim erkennt und
nicht als Tippfehler abtut. Eine Bindung ohne Deklaration ist eine
stille Setzung — und damit eine Harness-Lüge in derselben Klasse wie
ein halluziniertes Gate.

### harness/conventions.md als Konventionsspeicher

`harness/conventions.md` trägt die **repo-lokalen Strukturregeln** und
Adaptionen ggü. der adoptierten Baseline (Kurs, interner Standard,
Industrie-Norm). Sie ist **Pflicht** (Existenz), ihre Form (Einzeldatei
vs. Verzeichnis, ADR-artig vs. Prosa) ist **Wahl** — projektabhängig
nach Projektgröße, Adaptions-Frequenz, Audit-Tiefe.

Pflichtgliederung (Default-Form als Einzeldatei):

| Abschnitt                                     | Inhalt                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| --------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Purpose                                       | was die Datei trägt, was nicht                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| Baseline                                      | welche Konvention adoptiert, mit Stand/Version                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| Adoptierte Konventions-Quellen                | Pointer extern (Kurs/Standard) und in-Repo (Templates)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Adaptions-Block | **Index** der Abweichungen ggü. Baseline, nicht die Einträge selbst: `MR-000` (Adoptions-Erklärung) plus je eine Tabellenzeile pro Adaption. Pflichtfelder eines Eintrags: Datum, Geltungsbereich, `Ersetzt-Baseline-Regel`, Adaption, Begründung, Auflösungs-Trigger oder "permanent"). Löst ein Eintrag einen früheren **ab**, nennt er zusätzlich *Löst auf* und *Ausgelöst durch Baseline-Stand*; *schärft* er ihn nur (der alte gilt weiter, die Regel wird **strenger**), steht das im Titel — `(schärft MR-<NNN>)`. Verliert ein Eintrag durch die Baseline dagegen einen *Teil seines Geltungsbereichs*, ist das eine **Ablösung** mit engerem Nachfolger, keine Schärfung. Einträge werden nie überschrieben. |
| Zusatzklassen-Deklaration für Sensors-Bindung | repo-spezifische Bindung-Klassen jenseits der vier kanonischen (`LH-…`, Compliance, Modell-Version)                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| Modus-Deklaration pro Sub-Area                | Greenfield · Brownfield (mit Konvergenz-Auftrag) · Hybrid                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| Glossar (optional)                            | repo-spezifische Begriffe, die nicht im Regelwerk-Glossar stehen                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |

**Ein Eintrag je Datei — und der Grund ist der Kontext des Agenten.**
Die Einträge selbst leben unter `harness/conventions/MR-<NNN>-<titel>.md`;
ist der Auflösungs-Trigger eingetreten, wandert die Datei nach
`conventions/done/`. Der Zustand ist die **Verzeichnis-Position**, kein
Status-Feld — dieselbe Lifecycle-Form wie bei Slices
([`modul-05-planning-harness.md` §Lifecycle als State Machine](modul-05-planning-harness.md#lifecycle-als-state-machine)).

Der Schnitt folgt nicht der Ästhetik, sondern dem Lesepfad: `conventions.md`
liest **jeder** Agentenlauf. Steht der volle Text aller Adaptionen darin,
wächst der Pflichtanteil des Kontexts mit jeder Adaption — und trägt bald
mehrheitlich Einträge, die *aufgelöst* sind und trotzdem gelesen werden. Das
ist nicht nur Kontext-Kosten, sondern ein Korrektheits-Risiko: Ein
aufgelöster Eintrag liest sich wie ein geltender. Mit Index plus Dateien
zahlt jeder Lauf **eine Zeile pro aktiver Adaption**; geöffnet wird nur, was
den eigenen Geltungsbereich trifft.

Die Form bleibt Wahl: Ein Repo mit zwei permanenten Adaptionen darf sie
inline führen. Der **Default** ist die Verzeichnis-Form, weil sie mit der
Adaptions-Zahl nicht mitwächst.

**Von außen wird der Index adressiert, nicht die Eintrags-Datei — und dafür
braucht die Index-Zeile einen expliziten Anker.** Wer aus `AGENTS.md`, einem
Slice oder einer ADR auf eine Adaption zeigt, verlinkt
`harness/conventions.md#mr-<NNN>`, nicht `conventions/MR-<NNN>-<titel>.md`. Der
Grund ist der Lifecycle: Die Eintrags-Datei wandert bei Auflösung per `git mv`
nach `conventions/done/` — ein Pfad-Link darauf bricht **genau in dem Moment,
in dem die Adaption sich auflöst**. Die Index-Zeile wandert dabei nur von einer
Tabelle in die andere, innerhalb derselben Datei; der Anker reist mit. Es ist
dasselbe Argument, mit dem die Slice-ID ein Token bleibt statt ein Pfad
([`modul-05-planning-harness.md` §Lifecycle als State Machine](modul-05-planning-harness.md#lifecycle-als-state-machine)).

Ein solcher Anker muss **explizit** gesetzt werden: Eine Tabellenzeile bekommt
keinen automatischen Anker
([`grundlagen-source-precedence.md` §ID-Schema als Klammer](grundlagen-source-precedence.md#id-schema-als-klammer)),
also trägt die `MR`-Zelle ein `<a id="mr-<NNN>"></a>`. Adressiert wird die
**Kennung**, nicht der Titel — ein Anker, der den Titel enthält, bricht bei der
ersten Umformulierung und wäre damit genau die instabile Adresse, die er
ersetzen soll. Wandert ein Repo von der Inline- in die Verzeichnis-Form,
verschwinden mit den `### MR-<NNN> — <Titel>`-Überschriften auch deren
Auto-Anker: Dann trägt die Index-Zeile den alten Überschriften-Slug **zusätzlich**,
sonst rottet jeder bereits veröffentlichte Verweis.

Nebeneffekt, kein Selbstzweck: Ein Eintrag je Datei ist auch die einzige
Form, in der die Append-only-Disziplin *prüfbar* wird — eine wachsende
Sammeldatei lässt sich nicht gegen Core-Drift pinnen, eine akzeptierte
Einzeldatei schon.

Wichtig: `harness/conventions.md` dupliziert keinen Baseline-Text — sie
verweist und ergänzt. Eine Kopie ginge gegen die Baseline in Drift,
sobald letztere sich weiterentwickelt. Zwei Quellen derselben
Konvention sind dasselbe Drift-Risiko, das die Source-Precedence-Regel
für Spec/ADR adressiert — hier in der Form-Ebene.

Vorlagen:
[`templates/harness/conventions.template.md`](../templates/harness/conventions.template.md)
(Index) und
[`templates/harness/conventions/MR-NNN-titel.template.md`](../templates/harness/conventions/MR-NNN-titel.template.md)
(ein Eintrag).

### Jedes Artefakt hat einen Konsumenten

Wer dem Harness ein Artefakt hinzufügt — eine Sektion, eine Liste, eine Notiz
—, benennt, **wer es liest und wann**. Findet sich kein Leser, ist es Ablage,
keine Steuerung, und gehört nicht angelegt.

- **Derivative Artefakte** (ADR-Index, Carveout-Index, *Folge-Slices* in der
  Closure-Notiz) brauchen keinen eigenen Leser, wohl aber eine **Deckung**:
  das Original muss existieren. Als *derivativ* kennzeichnen, sonst schlägt
  die Probe falschen Alarm.
- **Lauf-Belege** (Review-Report, Verifikations-Belege) haben ihren Konsumenten
  im Vorgang selbst und danach im Audit; über Läufe hinweg werden sie nicht
  wieder gelesen und müssen es nicht.

**Einordnung — und ihre Grenze.** Die Regel ist *inferential feedforward* und
greift zur **Entwurfszeit**: wenn jemand den Harness *erweitert*, nicht wenn er
ihn *betreibt*. Sie ist ausdrücklich **kein Prüfpunkt der Closure-Prozedur** —
dort spräche sie in den meisten Wellen auf nichts an, würde nach der dritten
Welle übersprungen und wäre danach eine
[Harness-Lüge](grundlagen-begriffe.md#kernbegriffe).
Der häufige Fall ist gedeckt — eine Beobachtung, die die Schwelle erreicht und
zur Regel wird, hat ihren Leser automatisch (die verkörperte Form wird in jedem
Lauf gelesen), und die Anker-Paarung prüft deterministisch, dass sie landete.
Die Regel sagt **nicht**, ob ein genannter Konsument den Inhalt auch nutzt;
„wird beim Audit gelesen" ist gültig und zugleich die schwächste Antwort.
