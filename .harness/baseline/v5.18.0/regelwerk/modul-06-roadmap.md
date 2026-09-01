## Modul 6 — Roadmap Engineering

<!-- Quelle: [02-planung/modul-06-roadmap.md](https://github.com/pt9912/ai-harness-course/blob/v5.18.0/kurs/de/02-planung/modul-06-roadmap.md) -->

### Kernidee (Modul 6)

Eine Roadmap ist eine Reihenfolge von Wellen, keine Reihenfolge von
Terminen. Termine sind eine Folge der Wellen, nicht ihr Treiber.

Konkret ist eine Roadmap eine geordnete Folge von **Wellen**: jede Welle
bündelt Slices, schließt durch einen *beobachtbaren Trigger* — nicht durch
ein Datum — und hinterlässt eine Closure-Notiz. Ein Termin darf als
Schätzung *erscheinen*, triggert aber nie; er ist Output der
Wellen-Reihenfolge, nicht ihr Treiber. Die fünf Abschnitte unten sind die
Form, die Regeln der Inhalt.

### Roadmap-Regeln (Modul 6)

- Ein Welle-Eintrag braucht minimal drei Bestandteile: Slice-IDs (Inhalt) · Trigger als beobachtbare Bedingung (kein Datum) · Closure-Kriterien (z. B. Replay grün, alle Slices in `done/`). Datum darf *erwähnt* werden (Prognose), darf aber nie Trigger sein — sonst kappt die Welle halbfertige Slices am Kalendertag und das Auditierbarkeits-Versprechen bricht.
- Ein Trigger ist beobachtbar dann, wenn ein *anderer* Mensch ohne Rückfrage sagen kann, ob er eingetreten ist. "Sobald wir Zeit haben" scheitert daran; "slice-024 in `done/`" besteht. Beispiele für beobachtbare Trigger: "slice-024 liegt in `done/`" · "Replay-Lauf gegen Golden Set grün" · "Carveout `CO-007` aufgelöst".
- **Der Start-Trigger darf kein Ergebnis dieser Welle sein** — Beobachtbarkeit allein genügt nicht. "Alle Slices in `done/`" ist beobachtbar *und* ein Ergebnis: als Closure-Trigger richtig, als Start-Trigger zirkulär. Zwei Prüfungen, nicht eine. Test: Steht der Trigger in der Slice-Liste *dieser* Welle, ist er falsch platziert.
- Welle 30 % über Schätzung — Diagnose vor Aktion: liegt es an Slice-Größe (→ neu schneiden), an Reihenfolge (→ neu planen), oder an unerwarteter Komplexität (→ Carveout)? 30 % früh können ein Steering-Loop-Signal sein (Slice-Sizing-Regel schärfen), 30 % spät (vor Welle-Closure) eher Carveout.

### Welle ≠ Meilenstein ≠ Release (Modul 6)

- **Welle** = Bündel paralleler/serialisierter Slices mit Closure-Kriterien. Eine Welle endet *durch* Closure-Kriterien.
- **Meilenstein** = extern beobachtbarer Zustand (Release, Audit-Punkt). Ein Meilenstein endet durch *Datum oder externe Bestätigung* — und genau deshalb leitet sich der Meilenstein aus Wellen ab, nicht umgekehrt.
- **Release** — Trigger: ein Artefakt verlässt das Repo in eine Umgebung (Tag + Staging). Ein Release kann mehrere Wellen umfassen, der Meilenstein liegt *neben* der Welle (externe Bestätigung), die Welle endet *durch* Closure.

### Wann Arbeit eine Welle braucht (Modul 6)

- **Eine Welle liegt vor, wenn es eine beobachtbare Closure-Bedingung gibt, die mehr beobachtet, als die DoDs ihrer Slices schon belegen.** Das Kriterium liegt nicht in der Größe der Arbeit. Ein Trigger, der nichts beobachtet, was die Slices nicht ohnehin belegen, ist Zeremonie.
- Die kanonische Form dieses *Mehr*: alle Slices in `done/` **und** `make gates` grün **und** der Replay-Lauf grün — die beiden Gate-Bedingungen sind repo-weit und stehen in keiner einzelnen DoD.
- Fehlt dieses Mehr, gibt es keine Welle. Bei einem einzelnen Slice ist das der Regelfall: Sein Closure-Trigger würde die eigene DoD abschreiben. Solche Arbeit läuft **ohne Welle** — typisch für Reaktives (Sensor hat gefeuert, Pin ist veraltet, Meldung liegt vor), aber nicht darauf beschränkt: auch eine neue Fähigkeit kann ein einzelner Slice sein. Umgekehrt bleibt ein Ein-Slice-Bündel eine Welle, wenn sein Trigger repo-weite Belege fordert, die der Slice allein nicht liefert.
- **Wellenlose Arbeit erscheint nicht in der Roadmap** — weder beim Start noch beim Abschluss. Ihr Zustand ist die Verzeichnis-Position (Modul 5); `ls docs/plan/planning/in-progress/` beantwortet "was läuft gerade" autoritativ und ohne Pflegeaufwand — gelesen auf dem **Hauptzweig**: Der Übergang hierher landet dort, vor der Arbeit ([Modul 5 §Lifecycle als State Machine](modul-05-planning-harness.md#lifecycle-als-state-machine)). Ein Eintrag daneben wäre eine zweite Quelle für denselben Zustand, und die altert.
- Die Belege eines geschlossenen wellenlosen Slice stehen in seiner Datei und in git; das Closure-Log der Roadmap ist für Wellen.
- **Vorwärts-Blick:** „was kommt als Nächstes" beantwortet `next/` (*priorisiert/eingeplant*, Modul 5); der Übergang `open→next` ist die Priorisierungs-Entscheidung. Wellenlose Arbeit steht dabei nicht schlechter da als wellengebundene: Eine Reihenfolge *einzelner Slices* kennt der Harness nicht. Die Roadmap ordnet **Wellen**; die Spalte *Wichtigste Slices* nennt Inhalt, keinen Rang, und innerhalb einer Welle sind die Slices ein Bündel, das gemeinsam schließt. Eine Rangliste neben der Roadmap wäre eine Sortierung, die es für Slices nie gab — und die zweite Quelle, die die Regel oben vermeidet.
- **Wellenlos heißt nicht wächterlos.** Der Slice schreibt seine Closure-Notiz §7 wie jeder andere und trägt seine Beobachtungen ins **Beobachtungs-Register** ein — der Zähler unterscheidet nicht nach Welle-Zugehörigkeit. Damit zählt der Steering Loop weiter vollständig und offene Risiken finden ihren Ausgang.
- **Was offen bleibt:** Die **Carveout-Frist** misst in Wellen („seit > 2 Wellen aktiv", Modul 7). Wer lange wellenlos arbeitet, dehnt sie — ein Carveout steht dann bei „0 Wellen aktiv", obwohl Monate vergangen sind. Einen wellenlosen Ersatz-Träger gibt es dafür nicht.
- **Was der wellenlose Betrieb selbst auslöst:** alles, was am Slice hängt — mit benanntem Moment, sonst ist es ein Trigger ohne Wächter. **Achse zuerst:** *Wellenlos* ist eine Eigenschaft des **Repos** (Wellen+Slices, oder nur Slices), nicht des einzelnen Slice. Das Kopf-Feld `**Welle:**` eines Slice-Plans sagt nur, ob dieser Slice in ein Bündel gehört; daraus folgt für die Vorgänge unten nichts. Ein Repo mit Wellen hat eine Welle-Closure, und die liest und prüft alles, was seit der letzten Welle in `done/` liegt — **auch Slices ohne Wellen-Zugehörigkeit**.

  | Vorgang | Träger im Repo **ohne** Wellen | Wann |
  |---|---|---|
  | **Zähler** | Slice-Closure §7 | vor dem `git mv` nach `done/` |
  | **Lese-Schritt** (was hat 3× erreicht → Ausgang zuweisen) | Slice-Closure §7 | vor dem `git mv`; Anker `seit slice-<NNN>` statt `seit welle-<NN>` |
  | **Sichtungs-Schritt** (offene Beobachtungen unter der Schwelle) | Slice-**Planung**, §8 *Vorgelagert — offene Beobachtungen sichten* | beim Anlegen jedes Slice, unabhängig vom Sub-Area-Modus |
  | **Trigger-Audit** (Carveout · Bootstrap-aware Gate · ADR) | Slice-Closure | bei jeder Closure, zusammen mit dem Lese-Schritt |
  | **Alle drei Paarungen** (a/b/c aus Closure-Schritt 3) | Slice-Closure | **nach** dem `git mv` — sie suchen in `done/` |

  Ohne den Lese-Schritt bliebe der einzige Fall ungeprüft, in dem `seit slice-<NNN>` entsteht. Ohne den Sichtungs-Schritt hätte alles *unter* der Schwelle keinen Leser: In einem Repo mit Wellen trägt ihn die Wellen-Eröffnung Schritt 2 — ohne Wellen-Betrieb findet die nicht statt.
- **Einzige Berührung mit der Roadmap:** Liefert wellenlose Arbeit den letzten Beleg eines Meilensteins, bleibt die Spalte `Welle(n)` leer (`—`) und der Beleg steht als Slice-ID daneben — Beleg für eine externe Bedingung, nicht Zustand.

### Roadmap-Struktur: fünf Abschnitte (Modul 6)

Die Form liefert die Vorlage
[`templates/docs/plan/planning/roadmap.template.md`](../templates/docs/plan/planning/roadmap.template.md)
— fünf Abschnitte: *Offene Wellen · Nächste Wellen · Meilensteine ·
Abgeschlossene Wellen · Historische Trigger-Verschiebungen*. Operative Lesart:

- **Offene Wellen** — *derivativ*: Der Zustand sind die flachen Welle-Dateien, und woran gerade gearbeitet wird, sagt das `Welle:`-Feld der Slices in `in-progress/` ([Modul 5](modul-05-planning-harness.md#lifecycle-als-state-machine)). Ziel, Trigger und Closure-Kriterien stehen in der Welle-Datei, nicht hier. Der Abschnitt trägt **zwei unabhängige Aussagen**: Die *Liste* folgt den Dateien (ein Zeiger je offener Welle-Datei). Der Ruhe-Marker *Nichts in Arbeit* folgt dem Anspruch — er steht genau dann, wenn `in-progress/` keinen Slice trägt, **zusätzlich zur Liste, nicht an ihrer Stelle**; beides zugleich ist der Normalfall direkt nach der Wellen-Eröffnung (Welle eröffnet, noch nicht beansprucht). Zwei Aussagen, zwei Wächter — wer die Kopplung mechanisiert, muss wissen, *welche* Hälfte sein Sensor prüft, sonst hält er einen halben Wächter für einen ganzen. Die Marker-Hälfte ist die **deklarierte Redundanz**: Ein Doku-Sensor hält den Marker gegen das Verzeichnis, und zwar in **beide** Richtungen — ein fehlender Marker bei leerem `in-progress/` und ein stehengebliebener Marker bei beanspruchtem Slice sind derselbe Defekt. Die Listen-Hälfte ist kein Marker-Vergleich, sondern eine **Bijektion**: die im Abschnitt genannten Wellen-Kennungen gegen die flachen Welle-Dateien, ebenfalls in beide Richtungen — ein Zeiger ohne Datei und eine Datei ohne Zeiger sind derselbe Defekt. Sie hat eine Vorbedingung, die der Marker nicht hat: Der Sensor muss das **Kardinalitäts-Modell** kennen. Ein Wächter, der den Abschnitt gegen *genau eine* Datei hält (Ein-Wellen-Betrieb), meldet unter *Offene Wellen* legitime Zustände als Drift — zwei offene Wellen, oder eine Welle eröffnet und nichts beansprucht (Zeiger und Marker nebeneinander). Der Ruhe-Marker geht in die Bijektion **nicht** ein; er bleibt Sache des Marker-Wächters. Wer eine Hälfte ungewächtert lässt, benennt die Lücke — bekannt ist sie zulässig, verschwiegen nicht. Das *Geplante Ende* in der Welle-Datei ist Schätzung, kein Closure-Kriterium: kippt sie, kippt sie als Schätzung.
- **Nächste Wellen** — die geordnete Vorschau; jede Zeile trägt Welle, Trigger (die Abhängigkeit als beobachtbare Bedingung), wichtigste Slices und geschätzten Aufwand (S/M/L, kein Termin). Eine Welle, die ohne fertige Vorgängerin nicht starten kann, ist eine Phantom-Welle — die Abhängigkeit steht explizit in der `Trigger`-Spalte und als gerichtete Kante im Abhängigkeitsgraphen.
- **Meilensteine** — extern beobachtbare Zustände, orthogonal zur Welle: die Welle endet *durch* Closure-Kriterien (intern), der Meilenstein durch externe Bestätigung (Audit, Release, Kunde). Der Meilenstein liegt *neben* der Welle, nicht in ihr; ein Audit-*Termin* ist Anhang im Meilenstein-Eintrag, nie Trigger der Welle. Ist das externe Datum unverrückbar, aber die Closure-Trigger unerreichbar, ist die richtige Antwort ein *Carveout* (Modul 7), kein halbfertiges `done/`. Ein erreichter Meilenstein bleibt in der Tabelle: `Status` sagt *erreicht* mit Datum und Beleg.
- **Abgeschlossene Wellen** — das Closure-Log (ruhender Audit-Bestand): welche Welle wann geschlossen wurde, mit Zeiger auf ihre `done/welle-NN-results.md`. Es sagt, *was* geschlossen ist — Welle, Datum, Zeiger auf die Ergebnis-Notiz — und ist das einzige Closure-Log der Roadmap.
- **Historische Trigger-Verschiebungen** — das Drift-Log (Bewegungs-Signal): jede Umplanung mit Datum, Änderung, Grund. Wer es leer hat, hat eine starre Roadmap; wer *jeden* Eintrag voll hat, eine treibende. Closure-Log und Drift-Log zusammen machen die Vergangenheit der Roadmap auditierbar. Das Drift-Log sagt, *was umgeplant* wurde — ein Trigger verschoben, präzisiert oder ersetzt, ein Slice oder eine Welle umgehängt — und sonst nichts: Eine Schließung ist keine Umplanung, ein erreichter Meilenstein auch nicht; für den sagt die `Status`-Spalte der Meilenstein-Tabelle *erreicht* mit Datum und Beleg. Wer Schließungen oder Meilensteine ins Drift-Log schreibt, führt ein zweites Closure-Log, und zwei Logs driften. Jede `Stand`-/`Status`-Zelle — in der Roadmap wie im Beobachtungs-Register — trägt den Zustand und den Beleg als auflösbaren Anker, nie die Chronik ([`grundlagen-harness-dateien.md` §Was ein Kommentar trägt](grundlagen-harness-dateien.md#was-ein-kommentar-trägt--code-konfiguration-skripte), *Dieselbe Regel für Zustandsfelder*).

**Wird ein Closure-Trigger doch als Datum geschrieben** und der Kalendertag
erreicht, bevor die Slices grün sind, gibt es drei mögliche Antworten:

| Antwort | Diagnose |
|---|---|
| Welle wird trotzdem geschlossen, `slice-019` wandert in `welle-4`. | Datum hat Closure überschrieben — der Audit fällt durch, weil `slice-019` nicht belegt ist. Trigger-Disziplin ist Theorie geblieben. |
| Welle bleibt offen, das Datum wird verschoben. | Trigger-Disziplin wirkt, aber die Roadmap-Drift-Tabelle muss den Eintrag bekommen — sonst ist die Verschiebung still. |
| Carveout `CO-009` für die fehlende Latenz, Welle schließt mit Carveout. | Sauber: das Versprechen wird offen reduziert, Folge-Slice ist verdrahtet, Audit weiß, was er ansieht. |

*Eine Roadmap ist nicht „wann?", sondern „in welcher Reihenfolge wovon?"*.

### Das Beobachtungs-Register (Modul 6)

Der Zähler des Steering Loops liegt als **stehende Datei** flach im
Planning-Layout, neben den offenen Wellen: `docs/plan/planning/observations.md`
(Ziel-Form [`../templates/docs/plan/planning/observations.template.md`](../templates/docs/plan/planning/observations.template.md)).

- **Warum stehend:** Eine von Closure zu Closure übernommene Sektion hängt an
  einer ungebrochenen Kette — vergessene Übernahme setzt den Zähler auf null,
  die erste Welle braucht eine Sonderregel, ohne Welle gibt es keinen Träger.
  Der feste Ort streicht alle drei Fälle.
- **Form:** `| Kennung | Beobachtung | Sub-Area | Zähler | Belege | Stand |` — die Spalte `Stand` trägt den Zustand und den Beleg als auflösbaren Anker (*„verkörpert in `AGENTS.md` §2.7 (`seit welle-1`)"*), keine Chronik; zweite Sektion *Gestrichene Einträge* für das, was nicht mehr auftreten kann — mit Begründung. Wer eine Zeile still löscht, macht sie ununterscheidbar von einer, die es nie gab.
- **Ist nichts offen**, trägt die Tabelle `— keine —` und bleibt stehen. Die
  leere Liste **ist** die Aussage — und sie ist die, mit der jedes Repo
  anfängt.
- **Die Sub-Area-Spalte** trägt die Sub-Area, deren Konventions-Härte oder
  Inventur-Linie die Beobachtung betrifft — **nicht** die, in deren Verzeichnis
  sie aufgefallen ist. Dieselbe Berührungs-Frage wie beim §8-Block des
  Slice-Plans ([`grundlagen-bootstrap.md` §Was ist eine Sub-Area?](grundlagen-bootstrap.md#was-ist-eine-sub-area)),
  rückwärts gestellt. Steht in der Spalte ein Name, den die Modus-Deklaration
  in `harness/conventions.md` nicht führt, ist entweder die Zuordnung falsch
  oder die Deklaration unvollständig.
- **`BEO-<NNN>` ersetzt die Namens-Disziplin.** Erstauftreten benennt und
  vergibt die Kennung; Wiederauftreten zitiert sie und erhöht den Zähler. Ohne
  Kennung zählt eine Umformulierung als zweite Beobachtung, und keine erreicht
  je 3×. Das Register ist zugleich die Vergabestelle.
- **Wer schreibt, wer liest:** Eingetragen wird bei der **Slice-Closure** — neuer Eintrag mit neuer `BEO-<NNN>` oder Zähler erhöhen und Beleg ergänzen. Das macht den Zähler von der Welle unabhängig: Er läuft mit jedem geschlossenen Slice. **Gelesen wird an zwei Stellen:** die **Welle-Closure** liest, was 3× erreicht hat (*Lese-Schritt*, verkörpert); die **Slice-Planung** liest, was darunter steht (*Sichtungs-Schritt*, §8 des Slice-Plans — [Modul 5](modul-05-planning-harness.md#zwei-schritte-vor-der-modus-begründung)). Wer nur den ersten kennt, sieht alles unter 3× nie wieder an.
- **Mensch urteilt, Maschine prüft Deckung.** Das Urteil *ist das dieselbe Beobachtung?* fällt beim Schreiben (Kennung vergeben oder zitieren). Maschinell entscheidbar ist nur die Deckung: ob eine in `done/` zitierte `BEO-<NNN>` eine Registerzeile hat **und ob jede Registerzeile mindestens einen Beleg trägt** — die maschinelle Hälfte der Register-Paarung (c). *Nicht* geprüft wird die Umkehrung „jede Zeile ist irgendwo zitiert": Die allermeisten stehen unter der Schwelle und sind nirgends zitiert. Muster: schreiben → committen → Gate prüft. Welches Werkzeug, ist Repo-Entscheidung.
- **Der Beleg ist formgebunden** — drei Prüfungen ohne Urteil: **Form** (die Kennung eines abgeschlossenen **Vorgangs**, kein Freitext) · **Anzahl** (so viele wie der Zähler) · **Lage** (führt das Repo die genannte Datei, liegt sie dort, wo ihre Klasse abgeschlossen wird — für den Regelfall Slice also in `done/`) — diese dritte Prüfung läuft **nach** dem `git mv`, zusammen mit der Register-Paarung (c): Der Beleg wird vor dem `mv` geschrieben, der `mv` ist ein eigener Commit; auf dem Schreib-Commit läge die Datei noch nicht in `done/`, ein Sensor dort meldete bei jeder korrekten Closure rot. Form und Anzahl prüfen den Registereintrag und sind davon unabhängig. **Grenze:** Die *Existenz* der Datei wird nicht verlangt — ein Repo darf Slices führen, die es nicht als Plan-Datei ablegt; ein Sensor, der sie einforderte, liefe auf jedem gewachsenen Repo rot. Ein erfundenes `slice-999` bleibt damit unentdeckt — und die Grenze gilt für jede Vorgangs-Klasse gleich: Eine erfundene Welle- oder Review-Kennung bleibt aus demselben Grund unentdeckt. Das ist die Grenze der Deklaration und gehört benannt.
- **Ein Vorgang zählt einmal — und was keinen hat, zählt gar nicht.** Der Regelfall
  eines Belegs ist die Slice-Kennung; auch eine Welle und ein Review-Report sind
  abgeschlossene Vorgänge und taugen als Beleg. Zwei Funde **im selben** Vorgang
  sind dagegen *eine* Gelegenheit, kein zweites Auftreten: Der Zähler misst Wiederholung
  über Vorgänge hinweg, nicht die Zahl der Funde. Ein Vorkommen **ohne**
  abgeschlossenen Vorgang bekommt keinen Beleg und bewegt den Zähler nicht; es
  gehört trotzdem in den Eintrag — *benannt, nicht gezählt.*
- **Bei 3×** wandert der Eintrag in die Steering-Loop-Einträge der laufenden
  Welle-Closure und wird zur verkörperten Regel (mit Herkunfts-Anker) — ohne
  Wellen-Betrieb beim Lese-Schritt, den dann die Slice-Closure selbst auslöst, Anker `seit slice-<NNN>`. Die
  Zeile bleibt im Register mit Vermerk stehen; gestrichen wird nur mit
  Begründung, warum die Beobachtung nicht mehr auftreten kann.
- **Der Stand wird dabei zu einem von drei Ausgängen** — dieselbe geschlossene Menge wie
  beim offenen Risiko ([Modul 5](modul-05-planning-harness.md#offene-risiken-werden-bei-closure-aufgelöst)):

  | Ausgang | Wann | Wohin |
  |---|---|---|
  | **verkörpert** | die Regel steht | Zielort **und** Herkunfts-Anker (`seit welle-<NN>` bzw. `seit slice-<NNN>`) |
  | **geplant** | die Regel ist beschlossen, aber noch nicht geschrieben | Kennung des Slice oder der Welle, die sie schreibt |
  | **gestrichen** | die Beobachtung kann nicht mehr auftreten | §Gestrichene Einträge, **mit Begründung** |

  Zugewiesen wird der Ausgang vom **Lese-Schritt**; zwischen dem Beleg, der den
  Zähler auf 3 hebt, und diesem Schritt steht der Eintrag noch `offen` — das ist
  zulässig und vorübergehend. **Nicht** zulässig ist ein Eintrag, der eine
  Closure ohne Ausgang übersteht. Unterhalb der Schwelle ist `offen` der Stand —
  dort ist er kein Ausgang, sondern der Normalzustand. **Nur zwei der drei
  hängen an der Schwelle:** *verkörpert* und *geplant* sind ihre Antwort.
  *Gestrichen* ist an sie nicht gebunden — fällt die Ursache weg, bevor der
  Zähler 3 erreicht, wandert die Zeile mit Begründung in §Gestrichene Einträge. *Geplant* ist ein Ausgang **mit Kennung**, kein Vorsatz.
  **Urteilsfrei** ist, *dass* ein Eintrag ab 3× einen Ausgang trägt, *welcher
  der drei* es ist, und ob die genannte Kennung im Repo auflöst: Die drei sind
  eine geschlossene Menge, kein Freitext. **Urteil** bleibt, ob der Ausgang
  trägt.

### Wellen-Closure-Prozedur (Modul 6)

Modul 5 gibt den *Slice*-Zyklus als Zustandsmaschine vor (`open/` →
`next/` → `in-progress/` → `done/`). Die *Welle* liegt eine Ebene
darüber: Sie schließt nicht durch einen einzelnen Slice-Übergang, sondern
durch einen geordneten Ablauf, der alle ihre Slices bündelt. Sechs
Schritte — jeder hinterlässt einen Beleg, keiner ein Datum:

**Eröffnung — drei Schritte.** (1) Welle-Ziel, **Out-of-Scope** und
Closure-Trigger festlegen: beobachtbare Bedingung, kein Datum; erst danach
Slices zuordnen. Out-of-Scope gehört dazu — dieselbe Disziplin wie im
Lastenheft (Modul 3) und im Slice-Plan (Modul 9); was nicht ausdrücklich
ausgeschlossen ist, dehnt die Welle, bis der Closure-Trigger unerreichbar
wird. (2) **Offene Beobachtungen sichten** —
Das Register `docs/plan/planning/observations.md` durchgehen: Betrifft eine davon die
Sub-Areas dieser Welle, gehört sie in die Slice-Planung (Risiko im
betroffenen Slice) oder, bei Erreichen von 3×, als eigener Slice, der die
Lücke schließt. **Bei der ersten Welle entfällt dieser Schritt nicht** — das Register existiert
ab Repo-Beginn und ist durch die bis dahin geschlossenen Slices bereits gefüllt,
auch die wellenlosen. Ist es leer, ist das die Antwort und wird notiert. Das
ist der Schritt, der das Register auf der Planungsseite konsumiert; ohne ihn
bleibt es dort ohne Leser. **Ohne Wellen-Betrieb trägt ihn die Slice-Planung
selbst** (§8 des Slice-Plans, Block *Vorgelagert — offene Beobachtungen
sichten*, unabhängig vom Sub-Area-Modus). (3) Welle-Datei flach
anlegen (`docs/plan/planning/<welle-id>.md`, Ziel-Form
[`../templates/docs/plan/planning/welle.template.md`](../templates/docs/plan/planning/welle.template.md))
— ihre Zeile verlässt *Nächste Wellen*, unter *Offene Wellen* steht der
Zeiger auf die Datei.

**Nicht** in den Lauf-Kontext: `done/` wird dem Implementer-Agenten
nicht geladen. Schritt 2 ist Planungs-Leistung — was die Schwelle
erreicht hat, ist in AGENTS.md, Gates und Skills verkörpert und wirkt
dort automatisch (Modul-0-Prinzip).

**Wer führt die Schritte aus?** Die Eröffnung ist Planner-Arbeit, die Closure
hat ihre Übergaben in drei Zügen — Träger und Übergabe-Artefakt für **jeden**
ihrer Schritte in
[Modul 8 §Rollen-Sequenz für eine Welle](modul-08-agentenrollen.md#rollen-sequenz-für-eine-welle).

**Closure — sechs Schritte.**

1. **Trigger prüfen.** Alle Slices der Welle liegen in `done/`,
   `make gates` und der Replay-Lauf sind grün. Das ist die *beobachtbare*
   Closure-Bedingung aus der Welle-Definition — nicht der Kalendertag.
2. **Trigger-Audit der Welle.** Drei Artefaktklassen tragen einen Trigger,
   alle drei werden geprüft: **Carveout** (Auflösungs-Trigger → aufgelöst ·
   verlängert mit Folge-Slice · permanent, Modul 7) · **bootstrap-aware Gate**
   (Hochschalt-Trigger → Stufe hochschalten, oder Carveout eröffnen wenn die
   neue Schwelle rot ist, Modul 13) · **ADR** (Re-Evaluierungs-Trigger →
   bestätigen oder Folge-ADR mit `supersedes`, Modul 4). Eine Welle darf *mit*
   dokumentiertem Carveout schließen — aber nie mit einem stillen roten Gate,
   einer stehengebliebenen Reifestufe oder einer Entscheidung, deren
   Re-Evaluierungs-Bedingung vor drei Wellen eintrat. **Ein Trigger ohne
   Wächter ist eine Absichtserklärung mit Verfallsdatum.**
3. **Welle nach `done/` schließen.** Grundlage ist das **Beobachtungs-Register**,
   nicht die einzelnen Closure-Notizen: Dort steht der Zähler bereits,
   fortgeschrieben von jeder Slice-Closure. Die Welle-Closure ist der
   Lese-Schritt: Welche Einträge haben **3×** erreicht? Die bekommen ihren
   Ausgang — im Regelfall *verkörpert* — und werden zu
   *Steering-Loop-Einträgen*; im Register bleibt die Zeile mit
   dem Vermerk stehen, wohin sie ging. Was darunter liegt, bleibt offen und
   wartet. **Ohne diesen Lese-Schritt ist das Register write-only** — gezählt würde
   weiter, aber nichts würde je zur Regel.
   Closure-Notiz `done/welle-NN-results.md`
   schreiben (*was gelernt wurde*: geliefert · was funktionierte · was anders
   lief · **Steering-Loop-Einträge** (geschärfte Regel / neuer Sensor /
   benannte Spec-Lücke) · Zeiger aufs **Beobachtungs-Register** ·
   Folge-Slices (*derivativ* — der Folge-Slice selbst ist eine Datei in `open/`) ·
   Verifikation aus
   Schritt 1). Ziel-Form:
   [`../templates/docs/plan/planning/welle-results.template.md`](../templates/docs/plan/planning/welle-results.template.md).
   Ohne Lerneintrag ist die Welle nicht „fertig", nur „weg"
   (Modul 1). **Und die Welle-Plan-Datei wandert per `git mv` von flach nach
   `done/`** — neben ihre Ergebnis-Notiz; der Zustand ist die
   Verzeichnis-Position, kein `Status`-Feld (wie beim Slice). Offene Wellen
   flach, geschlossene in `done/`, die Roadmap bleibt Sequenzierungs-Autorität.
   **Zum Schluss alle drei Paarungen prüfen** — erst jetzt, weil sie die gerade
   entstandenen Einträge prüfen; in Schritt 2 gäbe es sie noch nicht. (Die
   Closure-Notiz wird direkt nach `done/` geschrieben; der `git mv` oben
   betrifft die Welle-*Plan*-Datei, die keine Paarung trägt.)
   (a) **Anker-Paarung** — ausgelöst durch das Pflichtfeld `liegt in <Zielort>`,
   **innerhalb dieser Sektion** und nicht durch die Semantik des Eintrags
   (der Trigger-Sprachgebrauch „`slice-024` liegt in `done/`“ aus
   §Roadmap-Regeln löst also nichts aus): Wo das Feld steht, existiert der
   Zielort und trägt `seit welle-<NN>` bzw. `seit slice-<NNN>`. Ein Eintrag
   **ohne** dieses Feld ist *gezählt, nicht verkörpert* und kein Gegenstand der
   Paarung. Die **benannte Spec-Lücke** ist der eine Fall, der ohne Feld
   trotzdem verkörpert ist — in einer versionierten Spec statt an einem
   Zielort; ihr Gegenstück ist die `LH-*`-ID
   ([`grundlagen-traceability.md` §Herkunfts-Anker](grundlagen-traceability.md#herkunfts-anker));
   (b) **Folge-Slice-Paarung** — jeder genannte Folge-Slice existiert als Datei
   **im Planning-Lifecycle** (`open/`, `next/`, `in-progress/`, `done/`), nicht
   nur in `open/`: bis zur Prüfung kann er weitergewandert sein.
   (c) **Register-Paarung** — zwei Hälften: jede in einer Closure-Notiz oder
   einem Risiko-Ausgang genannte `BEO-<NNN>` existiert als Zeile im
   Beobachtungs-Register, **und** jede Registerzeile trägt mindestens einen
   Beleg. *Nicht* geprüft wird die Umkehrung „jede Zeile ist irgendwo zitiert".
   Rot heißt in allen drei Fällen: etwas wurde
   versprochen und nicht angelegt.
4. **Zeitdokumente der Welle archivieren.** Die Slice-Dateien, die sie
   einsammelt, ihr eigener Plan und die Review-Reports dieser Slices wandern
   in ein unveränderliches Archiv `done/<welle-id>/archiv.zip`.

   Die **Ergebnisnotiz** bleibt vollständig und flach. Slice-Dateien und
   Welle-Plan bleiben als **gekürzter Stub** im Wellen-Verzeichnis:
   Überschrift, Archiv-Zeiger, Zustand, und die Kennungen, die den Vorgang
   überlebt haben Der Stub des **Welle-Plans**
   trägt statt der Kennungen den Zeiger auf seine Ergebnisnotiz und die Zahl
   der archivierten Vorgänge — die Zahl, gegen die sich die Vollständigkeit
   des Archivs abzählen lässt. Review-Reports bekommen keinen Stub; sie haben keine
   Identität jenseits ihres Slice. Ziel-Form:
   [`archiv-stub-slice.template.md`](../templates/docs/plan/planning/archiv-stub-slice.template.md)
   und [`archiv-stub-welle.template.md`](../templates/docs/plan/planning/archiv-stub-welle.template.md).

   **Eingesammelt wird nach der Welle, nicht nach dem Verzeichnis:** die
   Slices, deren `Welle:` diese Welle nennt, **und** die wellenlosen, die seit
   der letzten Closure geschlossen wurden. Slices einer noch **offenen** Welle
   bleiben liegen. Diese Auswahl gehört in die Operation, nicht in ihren
   Aufrufer.

   Der Stub hält die Verzeichnis-Position als Zustand und lässt eingehende
   Verweise gültig. Der Umzug ändert Pfade; die Operation zieht die Verweise
   nach — in **beiden** Formen, mit Verzeichnis-Präfix und
   geschwister-relativ. **Kein Zwang zum Nachrüsten — und kein Verbot:**
   Wellen, die vor der Einführung schlossen, müssen nicht archiviert werden;
   ein Repo bleibt ohne das konform. Wer den Altbestand loswerden will, führt
   die Archivierung als eigenen Vorgang aus, je geschlossener Welle einmal.
   Sie braucht dafür eine Entscheidung, die die laufende Regel nicht liefert:
   welche Welle die Slices einsammelt, die keiner angehören. Das Repo benennt
   die Zuordnung — die chronologisch nächste geschlossene Welle oder ein
   einzelnes Sammel-Archiv für den Bestand vor der Einführung.

   **Urteilsfrei** ist die **Form** des Stubs: dass er den Archiv-Zeiger trägt
   und die Abschnitte des vollen Plans **nicht** mehr. Zwei Bedingungen, kein
   Urteil — **die zweite ist die wichtigere**, denn ein Stub, der nur den
   Zeiger trägt und den Text behält, wäre die Archivierung, die es nicht gab.

   **Drei Grenzen:** Geprüft ist die Form, nicht die Länge. Ob das Archiv
   vollständig ist, bezeugt nur der Archivierungs-Commit — deshalb gehört die
   Operation in ein Werkzeug und nicht in Handarbeit. Und in einem Repo **ohne**
   Wellen-Betrieb fehlt der Auslöser ganz; die Frage bleibt offen.

   **Vor der ersten Archivierung ist der Geltungsbereich der vorhandenen
   Sensoren zu prüfen.** Ein Sensor, der auf `done/*.md` keilt, sieht die
   archivierten Stubs im Unterverzeichnis nicht mehr und bleibt grün, ohne
   noch etwas zu prüfen. Wer archiviert, zieht den Geltungsbereich mit — oder
   benennt, dass die Zusage für Stubs nicht mehr gilt.

5. **Wave-Self-Close-Commit.** Ein einzelner, beobachtbarer Commit
   markiert den Abschluss — der Audit sieht *einen* Punkt, an dem die
   Welle schloss, statt eines verstreuten Verschwindens.
6. **Roadmap fortschreiben.** Die Welle bekommt ihre Zeile in der Tabelle
   *Abgeschlossene Wellen* (mit Zeiger auf ihre Closure-Notiz), ihr Zeiger
   verlässt *Offene Wellen*. **Befördert wird niemand**: Welche Wellen offen
   sind, sagen die flachen Dateien; woran gearbeitet wird, das `Welle:`-Feld
   der Slices in `in-progress/`. Löste ein Trigger eine Umplanung aus, bekommt
   die *Historische Trigger-Verschiebungen*-Tabelle ihren Eintrag.

Erst wenn alle sechs Belege vorliegen, ist die Welle *auditierbar*
geschlossen.

### Regeln gegen typische Fehlannahmen (Modul 6)

- **Gegen "Roadmap ist eine Datumsleiste":** Datum ist Output, nicht Input. Wer Datumsleisten plant, plant Wunschdenken.
- **Gegen "Burndown ist Fortschritt":** Burndown ist *Tempo*. Fortschritt ist, ob die Welle das verspricht, was sie sollte.
- **Gegen "Eine Roadmap ist statisch":** Eine Roadmap, die nach drei Wellen nicht angepasst wurde, hat den Steering Loop nicht durchlaufen.
- **Gegen "Welle = Sprint":** Ein Sprint endet durch *Datum* (zwei Wochen sind um). Eine Welle endet durch *Closure-Kriterien* (alle ihre Slices in `done/`, Replay-Lauf grün, Closure-Einträge geschrieben). Wer Wellen wie Sprints schneidet, kappt halbfertige Slices am Datum — und produziert genau die Auditierbarkeits-Lücke, die der Harness verhindern soll.
- **Gegen "Trigger = Datum":** Ein Trigger ist eine *beobachtbare Bedingung* ("slice-024 liegt in `done/`", "Replay-Lauf gegen Golden Set grün", "Carveout `CO-007` aufgelöst"). Ein Datum ist kein Trigger, sondern eine Prognose. Wenn das einzige Trigger-Kriterium ein Kalendertag ist, plant die Roadmap nicht — sie hofft.
