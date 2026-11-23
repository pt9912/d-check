# Architect-Verdikt — welle-96 (`vcs` folgt einer Umbenennung), Plan gegen Entscheidungslage

**Rolle:** Architect (Modul 8), frischer Kontext, **vor** Code. Rolleninhaber: pt9912.
**Gegenstand:** [welle-96](../plan/planning/welle-96-vcs-folgt-umbenennung.md) und ihre zwei Slice-Pläne
[slice-291](../plan/planning/next/slice-291-vcs-umbenennung-anforderung-und-adr.md) (Entscheide 1–8 in §3,
vom Auftraggeber am 2026-11-23 mit den Empfehlungen bestätigt) und
[slice-292](../plan/planning/open/slice-292-vcs-umbenennung-modulkern.md).
**Anlass:** [Issue #18](https://github.com/pt9912/d-check/issues/18) — ein Link-Nachzug in einer
`Accepted`-ADR auf eine in derselben Range umbenannte Datei meldet `core-drift-vcs`, auch mit
`vcs.ignore-link-targets`.
**Datum:** 2026-11-23.
**Geändert wurde kein Plan und keine Spec-Datei.** Mit diesem Verdikt entsteht die ADR als
[ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) (Status `Proposed`) samt Index-Zeile; sie trägt die
Entscheide 1–8 wie bestätigt und die Folge-Entscheide dieses Verdikts als Entscheidung 9–14 (Empfehlung).

## Gesamturteil

Die Richtung trägt: eigener opt-in-Schalter, Identität aus der Blob-Gleichheit je Commit, Kette über die
Commits der Range, Mehrdeutigkeit als Drift, ohne Schalter byte-identisch, lesend und netzlos. **Kein
bestätigter Entscheid widerspricht dem Lastenheft oder einer `Accepted`-ADR in der Sache.** Das Lastenheft
muss aber an einer Stelle mehr geändert werden, als der Plan führt (F-1). **Vier Entscheidungen fehlen**,
die der Code sonst stillschweigend träfe (F-2 Seite der Gleichsetzung, F-3 Zusammenspiel der Schalter, F-4
Topologie der Kette, F-5 Basis kein Vorfahr), dazu der Vertrag der neuen Port-Methode (F-6). **Drei Stellen
behaupten mehr, als der Gegenstand hergibt** (F-7 Begründung von Entscheid 1, F-8 Determinismus-Test, F-10
README-Spiegel). Entscheid 7 ist eine Lockerung, die ADR-0115 trägt; es fehlen nur Reihenfolge und zwei
Zeiger (F-9). Alle Folge-Entscheidungen stehen in der **einen** ADR; eine zweite ADR, ein
Beobachtungs-Eintrag oder ein Folge-Slice ist für keinen Befund nötig.

Verdikt-Kürzel (Modul 8, Konflikt-Pfad): **V1** = Entscheidung gilt, Plan hat falsch behauptet/ausgelassen ·
**V2** = Folge-Entscheidung nötig · **V3** = Lockerung legitim, aber undokumentiert.

## Befunde

### F-1 — Das Lastenheft wird an mehr Stellen geändert, als der Plan als berührt führt — **V1**

- slice-291 nennt als berührt Absatz „Pfad-Nachzug (opt-in)“, Akzeptanzkriterien, Version und Historie
  ([slice-291](../plan/planning/next/slice-291-vcs-umbenennung-anforderung-und-adr.md) Z. 21–25). Nicht
  geführt ist der Out-of-Scope-Satz von
  [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in):
  „Umbenennungs-Erkennung über Inhalts-Ähnlichkeit statt Pfad (eine immutable Datei behält ihren Pfad)“
  ([lastenheft.md](../../spec/lastenheft.md) Z. 2299–2300). Entscheid 1 widerspricht ihm nicht — Gleichheit
  ist keine Ähnlichkeit, und die Klammer meint die immutable Datei —, aber eine Anforderung, die
  Umbenennungen erkennt, neben einem Ausschluss, der Umbenennungs-Erkennung ausschließt, liest sich als
  Widerspruch im obersten Rang. Der Satz ist auf „für die immutable Datei selbst; für Link-Ziele zählt nur
  Inhalts-Gleichheit, nie Ähnlichkeit“ zu schärfen.
- Ebenso zu qualifizieren sind die Beschreibung Z. 2242–2245 und Z. 2251–2252 („denselben Namen trägt“) und in
  der Spezifikation Schritt 4 Z. 2076–2077, Z. 2096–2097 und die Fail-safe-Liste Z. 2101–2103
  („Umbenennung beim Umzug“) — Schritt 4 führt der Plan, die Stellen nicht.
- Das Kriterium „Boundary (Pfad-Nachzug)“ (Z. 2272) ist unter `ignore-link-targets` allein weiter richtig,
  weil Entscheid 4 den neuen Schalter unabhängig macht. Die kleinste Änderung ist, es auf „ohne
  `vcs.follow-renames`“ einzuschränken und ein neues Kriterium daneben zu stellen — nicht es umzuschreiben.
  Damit bleiben die Tests von `ignore-link-targets` unverändert gültig (slice-292 Risiko 1).
- **Verdikt:** Lastenheft gilt; der Plan lässt Stellen aus. Kein Widerspruch eines Entscheids.
- **Plan-Änderung:** slice-291 Kopf „Berührte Spec-Stellen“ und DoD-Punkt 1 um den Out-of-Scope-Satz und
  die Beschreibungs-Sätze ergänzen; DoD-Punkt 1 „Kriterium widerspricht nicht mehr“ als „eingeschränkt auf
  ohne `vcs.follow-renames`“ fassen.

### F-2 — Entscheid 5 setzt nur auf der BASE-Seite gleich; das macht unveränderte ADRs rot — **V2** (Entscheidung 9)

- Entscheid 5: „Ein Link-Ziel, das in BASE auf den Pfad A auflöst, gilt als Ziel auf B“
  ([slice-291](../plan/planning/next/slice-291-vcs-umbenennung-anforderung-und-adr.md) Z. 153–155).
- Am Gegenstand: das Modul vergleicht den Core **jeder** geschützten Datei, die in beiden Ständen liegt,
  nicht nur der geänderten ([vcs.go](../../internal/hexagon/core/rules/vcs.go) Z. 58–66), mit **einer**
  Auflösungs-Funktion für beide Fassungen (Z. 130–132). Wird nur die BASE-Fassung gleichgesetzt, ergibt
  eine unveränderte ADR mit einem Link auf eine umbenannte Datei zwei verschiedene Cores — ein Befund an
  einer Datei, die die Range nicht ändert. Das Lastenheft prüft „jede in der Range **geänderte** Datei“
  ([lastenheft.md](../../spec/lastenheft.md) Z. 2214); die einseitige Lesart verletzte Rang 1.
- **Verdikt:** Folge-Entscheid. Empfehlung: die Gleichsetzung ist **eine** Funktion des aufgelösten Pfads,
  auf beide Fassungen gleich angewandt, gegen die Vereinigung beider Pfad-Bäume wie heute
  ([vcs_link_targets.go](../../internal/hexagon/core/rules/vcs_link_targets.go) Z. 65–93). Teil hat ein Pfad,
  der Quelle einer eindeutigen Umbenennung mit einem Ende in HEAD ist, oder ein solches Ende; ein Pfad, der
  zugleich Quelle und Ende ist, nicht (Tausch zweier Namen bleibt Drift). Invariante: gleicher Rohtext ⇒
  gleicher Core.
- **Plan-Änderung:** slice-291 Entscheid 5 auf diese Fassung bringen (Neu-Bestätigung); slice-292 §3
  Testfälle „unveränderter Link auf eine umbenannte Datei in einer ADR, die nur einen erlaubten Anhang
  bekommt ⇒ kein Befund“ und „Tausch A ↔ B ⇒ Befund“.

### F-3 — Zusammenspiel von `follow-renames` und `ignore-link-targets`: unbestimmt — **V2** (Entscheidung 10), und die Welle behauptet „nicht geschärft“ — **V1**

- Entscheid 4 nennt die Schalter unabhängig, Entscheid 5 sagt, ein nicht teilnehmendes Ziel werde „wie
  bisher“ behandelt. Offen bleibt die Normalform eines teilnehmenden Ziels: voller Endpfad oder Dateiname
  des Endes, und was gilt, wenn beide an sind. Heute erzeugt die Normierung Marke, Dateiname, Anker
  (vcs_link_targets.go Z. 91), und ohne `ignore-link-targets` gibt es gar keine Auflösung (vcs.go Z. 53–56).
- **Verdikt:** Folge-Entscheid. Empfehlung: nur `follow-renames` ⇒ Marke, repo-relativer Endpfad, Anker;
  beide ⇒ erst gleichsetzen, dann die bestehende Normierung auf den Dateinamen des Endes. So bleiben die
  Filter von `ignore-link-targets` unangetastet (slice-292 Risiko 1). Folge, die zu benennen ist: mit beiden
  Schaltern sind zwei Ziele gleichen Namens, deren Enden verschieden heißen, Drift — etwa ein Nachzug von
  einer umbenannten Datei auf eine **andere** Datei mit dem alten Namen. Das ist strenger als
  `ignore-link-targets` allein und schließt einen Teil der Grenze
  [`BEO-ALL/pfad-nachzug-gleicher-name-andere-datei`](../plan/planning/observations/BEO-ALL/pfad-nachzug-gleicher-name-andere-datei/observation.md).
  Opt-in und in die sichere Richtung — aber nicht, was die Welle sagt: „wird nicht geschärft“
  ([welle-96](../plan/planning/welle-96-vcs-folgt-umbenennung.md) Z. 120–124).
- **Plan-Änderung:** Welle §6 zweiter Punkt: „bleibt für Ziele ohne Umbenennung; mit beiden Schaltern wird
  der Fall ‚umbenannte Quelle, gleichnamige andere Datei‘ Drift (ADR-0115 Konsequenzen)“; slice-292 §3
  Testfall dazu.

### F-4 — Topologie der Kette: Entscheidung fehlt — **V2** (Entscheidung 11)

- Entscheid 2 verkettet A → B und B → C zu A → C, sagt aber nicht, **in welcher Ordnung**. Die Commit-Menge
  liefert der Port heute über `CommitMessages`, nach Commit-SHA sortiert
  ([git.go](../../internal/adapter/driven/git/git.go) Z. 187–228, Sortierung Z. 226) — deterministisch, aber
  nicht topologisch. Verkettet der Code in dieser Reihenfolge, hängt das Ergebnis an Hash-Werten; verkettet
  er ordnungsfrei über alle Kanten, setzt er Umbenennungen paralleler Linien fort (X → Y auf der einen,
  Y → Z auf der anderen) und erzeugt ein falsches Grün.
- Ein HEAD, der selbst ein Merge-Commit ist (etwa der Test-Merge eines Pull-Request-Laufs), verlangt, dass
  ein Merge die Umbenennung **einer** Seite weiterträgt, wenn die andere den Pfad unberührt lässt; sonst
  folgte der Nachzug in diesem Lauf nie.
- **Verdikt:** Folge-Entscheid. Empfehlung: eine Umbenennung setzt eine frühere nur fort, wenn deren Commit
  ein Vorfahr ist; je Pfad zählen die Enden aller Linien, die ihn umbenennen, eine unberührte Linie trägt
  kein Ende bei; genau ein Ende, das in HEAD existiert, ist eine Umbenennung. Das Ergebnis hängt nicht von
  der Aufzählungsreihenfolge ab ([`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Plan-Änderung:** slice-291 DoD-Punkt 2 (Verfahren): Ordnung der Kette; slice-292 §3 Testfälle
  „Umbenennungen auf zwei parallelen Linien setzen sich nicht fort“, „HEAD ist ein Merge, eine Seite
  benennt um ⇒ folgt“, „Ende der Kette in HEAD gelöscht ⇒ Befund“.

### F-5 — Basis kein Vorfahr: der Plan lässt zwei Antworten offen, und Exit 2 widerspräche ADR-0109 im eigenen Profil — **V2** (Entscheidung 12)

- slice-291 Risiko 3 erwartet „fail-closed, Exit 2“ (Z. 240–245), slice-292 §3 lässt „Exit 2 oder benannte
  Menge“ offen (Z. 137–138). Heute verlangt `vcs` **keine** Abstammung (git.go `AllPaths` Z. 107–133, nur
  der Leerfall); `LineChanges` verlangt sie (Z. 391–408).
- Exit 2 wäre im eigenen Profil (Entscheid 7) ein Rückfall hinter
  [ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md): nach einem Force-Push prüft
  die CI `origin/<default>..HEAD` ([ci-range.sh](../../tools/harness/ci-range.sh) Z. 54); ist der
  Default-Branch nach dem Rebase weitergezogen, ist diese Basis kein Vorfahr, und der Push-Lauf würde rot,
  ohne dass eine Eigenschaft verletzt ist — genau der Zustand, den ADR-0109 beseitigt hat.
- Die Kette ab der Merge-Basis zu lesen, ergäbe ein falsches Grün: belegt die BASE-Seite einen Pfad A neu
  (C → A) und benennt die HEAD-Seite das alte A nach B um, gilt ein Link auf das neue A als nachgezogen.
- **Verdikt:** Folge-Entscheid. Empfehlung: ist BASE kein Vorfahr von HEAD, trägt die Range **keine**
  Umbenennung; das Modul urteilt wie ohne Schalter — der strengere, bisherige Ausgang, kein stilles Grün.
  Die Anforderung nennt es als Grenze. Der Vorfahren-Walk ist kein neuer Bedarf: `AllPaths` wandert ihn
  schon für den Leerfall (DC-FA-VCS-002, Grenze shallow-Klon, [lastenheft.md](../../spec/lastenheft.md)
  Z. 2322–2324).
- **Plan-Änderung:** slice-291 Risiko 3: Ausgang *entfallen*, wenn die Spezifikation die Commit-Menge und
  „Basis kein Vorfahr ⇒ keine Umbenennung“ festschreibt; slice-292 Z. 137–138 auf dieses Verhalten.

### F-6 — Vertrag der neuen Port-Methode: fail-closed erben, Bibliothek nicht nutzen — **V2** (Entscheidung 13)

- Am Gegenstand gemessen (go-git v5.19.3, im Modul-Cache des Test-Images gelesen, kein Host-Go): die
  Vorbelegung der Diff-Optionen ist `DetectRenames: true`, `RenameScore: 60` (difftree.go Z. 49–54); mit
  `OnlyExactRenames` wählt die exakte Erkennung unter mehreren gelöschten Kandidaten gleichen Blobs den mit
  dem ähnlichsten Namen und verlangt gleichen Modus (rename.go Z. 68–111). Das weicht von Entscheid 3
  (Mehrdeutigkeit ⇒ keine Umbenennung) und Entscheid 1 (nur der Blob) ab. Der Diff-Walker derselben
  Bibliothek las ein unlesbares Unterbaum-Objekt früher still als Ende (Historie
  [spezifikation.md](../../spec/spezifikation.md) Z. 4269). slice-292 §3 („kein Diff-Walker“, Z. 103) liegt
  damit richtig.
- Offen ist, was die Methode lesen **muss**: Entscheid 6 sagt „dieselbe Politik wie `AllPaths`“, aber
  `AllPaths` liest zwei ganze Bäume, die Methode viele Schritte.
- Blob-Kennungen stehen in Tree- und Index-Einträgen (git.go Z. 159–175, Z. 336–340); Gleichheit ist ohne
  Lesen des Blobs entscheidbar. Die Rückführung „go-git liefert die Blob-Kennungen nicht“ (slice-292
  Z. 161–163) ist damit kaum erreichbar — darf stehen bleiben.
- **Verdikt:** Folge-Entscheid. Empfehlung: nur bei aktivem Schalter gerufen; liefert die Abbildung Quelle →
  eindeutiges Ende (Kette und Mehrdeutigkeit im Adapter, Gleichsetzen im Kern — wie slice-292 DoD-Punkt 1
  es schon schneidet); liest jeden Commit der Menge und jeden Teilbaum, der sich zwischen Commit und
  Eltern-Commit unterscheidet, ein Teilbaum gleicher Kennung muss nicht gelesen werden; jedes dafür nötige,
  unlesbare Objekt ⇒ Fehler, auch wenn BASE und HEAD lesbar sind. Teil haben Datei, ausführbare Datei,
  Symlink; gleich sind gleiche Blob-Kennung **und** gleiche Art, das Ausführbar-Bit zählt nicht.
- **Plan-Änderung:** slice-292 §3 Zeile Port/Adapter: diesen Vertrag zitieren; Testfälle „Ausführbar-Bit
  beim Umbenennen geändert ⇒ folgt“, „Symlink und Datei mit gleichem Blob ⇒ keine Umbenennung“.

### F-7 — Die Begründung von Entscheid 1 behauptet mehr, als Entscheid 2 zulässt — **V1**

- Entscheid 1 verwirft die Schwelle, weil sie „umbenennen und umschreiben“ durchließe, „und genau das ist
  das Umbiegen auf ein anderes Dokument“ (slice-291 Z. 127–131). Entscheid 2 lässt aber denselben Vorgang
  in **zwei** Commits durch: reiner Move, dann Umschrift der umbenannten Datei — mit Absicht (Z. 138–142).
- Die Entscheidung trägt trotzdem, auf den anderen Gründen: ein Parameter, den niemand entschieden hat;
  Determinismus; ein reiner Move-Commit nach [AGENTS.md](../../AGENTS.md) §3.3 ist genau die erkannte
  Form. Und die Zusage ist Pfad-Identität entlang der Historie, nicht Inhalts-Identität des Ziels — auch ein
  unveränderter Link auf eine umgeschriebene Datei ist heute keine Drift.
- **Verdikt:** Entscheid gilt; die Begründung ist zu korrigieren. **Akzeptiertes Negativ** (in ADR-0115
  Konsequenzen): wer eine Umbenennung vortäuscht, um umzubiegen, sieht das Gate nicht; das trägt der Review.
- **Plan-Änderung:** slice-291 Entscheid 1, Begründung auf die tragenden Gründe; die Welle-Gegenprobe
  (Z. 70–71, ein Commit) bleibt richtig.

### F-8 — Determinismus: „zweimal“ ist weniger als die Messmethode von Rang 1 — **V1**

- slice-292 Z. 144: „Zweimal derselbe Lauf ⇒ byte-gleiche Ausgabe.“ Die Messmethode von
  [`DC-QA-02`](../../spec/lastenheft.md#dc-qa-02--determinismus) ist „denselben Lauf 10× … Hashes“
  ([lastenheft.md](../../spec/lastenheft.md) Z. 4597). Die neue Stelle baut Mengen über Blob-Kennungen und
  Pfade; zwei Läufe fangen eine Abhängigkeit von der Aufzählungsreihenfolge selten.
- **Plan-Änderung:** slice-292 Z. 144: zehn Läufe über ein Fixture mit Kopien und verzweigter Kette.

### F-9 — Entscheid 7 lockert das eigene Gate; ADR-0115 trägt es, es fehlen Reihenfolge und zwei Zeiger — **V1** (Entscheidung 14)

- Ja, es ist eine Lockerung einer Prüfregel ([AGENTS.md](../../AGENTS.md) §3.6): ein Nachzug auf eine
  umbenannte Datei unter neuem Namen war im eigenen `make adr-check` rot und wird grün. ADR-0115 trägt sie
  ausdrücklich (Entscheidung 7) — kein V3.
- Sie **löst ADR-0103 nicht ab**: deren Entscheidungen 1–4 (Schlüssel, eigenes Profil, was Befund bleibt,
  Nennung in §3.5) bleiben wahr, die Wirkung von `ignore-link-targets` ändert sich nicht
  ([ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) §Entscheidung). Ein neuer Schlüssel
  daneben erweitert; kein `supersedes` nötig.
- Zwei Lücken: (1) Entscheid 8 hält die ADR `Proposed` bis zur Closure von slice-292, der Schalter im Profil
  kommt in slice-292 — ohne Reihenfolge liefe die Lockerung auf `main` mit einer nicht angenommenen ADR.
  (2) [AGENTS.md](../../AGENTS.md) §3.5 und [harness/README.md](../../harness/README.md) §Traceability rules
  nennen ADR-0103 als Träger des Pfad-Nachzugs; die Spiegel-Suche nach dem Schlüssel (slice-292 Z. 120–123)
  findet sie nicht.
- **Plan-Änderung:** slice-292 DoD-Punkt 3: „`.d-check.yml` schaltet den Schlüssel im selben Commit ein,
  der ADR-0115 nach `Accepted` hebt“; §3 um AGENTS.md §3.5 und harness/README.md (Zeiger auf ADR-0115 neben
  ADR-0103) ergänzen.

### F-10 — README-Spiegel im Feature-Commit: die Begründung trägt nicht, Regel 17 verlangt Release-Prep — **V1**

- slice-292 §3 Z. 109 nimmt `README.md` und `README.de.md` in den Feature-Commit, weil sie „die Grenze
  ‚anderer Name bleibt Drift‘“ nennten. Am Gegenstand nennen sie sie nicht (README.de.md Z. 89–94,
  README.md Z. 93: nur „reiner Pfad-Nachzug … keine Drift“, weiter wahr). README-Prosa ist Release-Prep
  ([AGENTS.md](../../AGENTS.md) §5 Zeile 17); einen Spiegel-Test für diese Zeile gibt es nicht.
- Die Grenze nennen das Handbuch ([benutzerhandbuch.md](../../docs/user/benutzerhandbuch.md) Z. 2155–2159)
  und die Sensor-Datei ([adr-check.md](../../harness/sensors/adr-check.md) Z. 10–15) — die bleiben im Plan.
- **Plan-Änderung:** slice-292 §3 Z. 109 ohne die READMEs; sie gehen in die Release-Prep.

## Bestätigt, mit akzeptiertem Negativ

- **`--staged`-Grenze (slice-291 Risiko 2).** Einen Weg, den früheren Move-Commit mitzulesen, bestätigt der
  Architect **nicht**: im gestagten Modus ist BASE per Definition HEAD (spezifikation.md Z. 1990); den
  Commit HEAD zusätzlich zu lesen wäre willkürlich und fände einen Move zwei Commits zurück trotzdem nicht.
  Die CI-Range sieht beide Commits, und das Handbuch gibt den Weg schon an (Z. 2160–2161, Verweise im
  Move-Commit; [`BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise`](../plan/planning/observations/BEO-ALL/umzug-mit-adr-nachzug-braucht-mitreisende-verweise/observation.md)).
  **Akzeptiertes Negativ**, als Alternative J in ADR-0115 verglichen. Der Ausgang des Risikos folgt dem Plan
  (*weiter offen*, Grenze in der Anforderung).
- **Gestagter und Range-Modus geben verschiedene Antworten — aber auf verschiedene Paare.** Für dasselbe
  BASE/HEAD-Paar ist der Schritt derselbe; die Kopplungs-Aussage im pre-commit-Hook bleibt wahr.
- **fail-closed und „nicht auflösbare Range“.** Unverändert; neu ist nur ein Exit 2 für ein unlesbares
  Objekt in einem **Zwischen**-Commit und nur mit Schalter (Entscheid 6) — in der Anforderung zu nennen,
  steht so im Plan. Kein neuer Bedarf an Historien-Tiefe (F-5).
- **[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).**
  Gelesen wird nur `.git`, nichts geschrieben, kein Netz; die Messmethode nimmt `vcs` ohnehin aus.
- **Entscheid 2 „dieselbe Menge wie `commits`“.** Stimmt am Gegenstand (git.go Z. 199–225: Ausschluss der
  Basis-Vorfahren, höchstens ein Elternteil).
- **Entscheid 4.** Der Grund gegen die Erweiterung von `ignore-link-targets` trägt (§3.6); als Alternative B
  geführt.
- **Entscheid 8.** `Proposed` bis zum umsetzenden Slice, wie ADR-0114.
- **Anlass gezählt.** Das Kommando aus slice-291 Z. 114–118 liefert am 2026-11-23 erneut 107 und 59.
- **Ein einziger Test-Doppelgänger des Ports.** `LineChanges` implementieren nur git.go und
  `internal/hexagon/core/rules/vcs_test.go`; die Planzeile trägt.
- **Merge-Seite, Tausch, Kopien.** Drei Grenzen in die sichere oder in eine konfliktbehaftete Richtung,
  in ADR-0115 unter Konsequenzen benannt; keine eigene Folgepflicht.
- **Die Welle-Probe (Z. 66–71)** trägt so, wie sie steht; im Wegwerf-Klon ist der Schalter für „ohne“ aus
  dem Profil zu nehmen, sobald Entscheid 7 umgesetzt ist.

## Nicht geprüft

- Kein Code gefahren, keine git-Fixtures gebaut; die Aussagen über das Verhalten stehen auf gelesenem
  Code.
- Die Bibliothek nur in den zwei zitierten Dateien gelesen; ob ihr Diff-Walker in v5.19.3 ein unlesbares
  Objekt noch still überspringt, nicht nachgemessen — die Empfehlung meidet ihn ohnehin.
- Dass ein Pull-Request-Lauf einen Test-Merge als HEAD prüft, ist nicht am Gegenstand gemessen; F-4 trägt
  auch ohne diese Annahme (jeder Merge als HEAD).
- Die Laufzeit des Commit-Walks (slice-292 Risiko 3) nicht gemessen.
- Gates: nur `make doc-check` über dieses Verdikt, ADR-0115 und den Index.

## Zusammenfassung

| Befund | Verdikt | Träger der Korrektur |
|---|---|---|
| F-1 Out-of-Scope-Satz und Beschreibung nicht als berührt geführt; Kriterium einschränken statt ersetzen | V1 | Plan 291 |
| F-2 Gleichsetzung nur auf BASE macht unveränderte ADRs rot | V2 | ADR-0115 Entscheidung 9; Plan 291/292 |
| F-3 Zusammenspiel der Schalter; Welle „nicht geschärft“ | V2 + V1 | ADR-0115 Entscheidung 10; Welle §6; Plan 292 |
| F-4 Topologie der Kette | V2 | ADR-0115 Entscheidung 11; Plan 291/292 |
| F-5 Basis kein Vorfahr | V2 | ADR-0115 Entscheidung 12; Plan 291/292 |
| F-6 Port-Vertrag, Bibliothek nicht nutzen | V2 | ADR-0115 Entscheidung 13; Plan 292 |
| F-7 Begründung von Entscheid 1 | V1 | Plan 291; ADR-0115 Konsequenzen |
| F-8 zehn statt zwei Läufe | V1 | Plan 292 |
| F-9 Lockerung getragen; Reihenfolge und Zeiger | V1 | ADR-0115 Entscheidung 14; Plan 292 |
| F-10 READMEs in die Release-Prep | V1 | Plan 292 |

Kein Befund fällt unter V3: die einzige Lockerung (Entscheid 7) trägt ADR-0115.

**Braucht eine Entscheidung des Auftraggebers:** Entscheidung 9–14 als Empfehlung; davon ändert 9 den
Wortlaut des bestätigten Entscheids 5 (Neu-Bestätigung), 10 widerspricht der Welle-Abgrenzung §6 in einem
Fall, 12 ersetzt den im Plan erwarteten Exit 2. Die Schärfung des Out-of-Scope-Satzes (F-1) ist eine
Rang-1-Änderung außerhalb der geführten Stellen und geht mit slice-291 zur Bestätigung. F-7 braucht keine
Entscheidung, nur Kenntnis: Entscheid 1 steht auf anderen Gründen als den bestätigten.

## Plan-Änderungen für den Planner

1. slice-291 Kopf „Berührte Spec-Stellen“ und DoD-Punkt 1: Out-of-Scope-Satz von DC-FA-VCS-001
   (Lastenheft Z. 2299–2300) und die Beschreibungs-Sätze Z. 2242–2245, 2251–2252; Kriterium „Boundary
   (Pfad-Nachzug)“ auf „ohne `vcs.follow-renames`“ einschränken, neues Kriterium daneben (F-1).
2. slice-291 DoD-Punkt 2: in Spezifikation Schritt 4 auch Z. 2076–2077, 2096–2097, 2101–2103 (F-1).
3. slice-291 §3 Entscheid 5 auf die seiten-gleiche Fassung (ADR-0115 Entscheidung 9), neu bestätigen (F-2).
4. slice-291 §3 Entscheidungen 9–14 aus ADR-0115 zur Bestätigung aufnehmen; §4 Rückführung
   `in-progress → open` um diese Punkte ergänzen (F-2 bis F-6, F-9).
5. slice-291 §3 Entscheid 1: Begründung auf die tragenden Gründe (F-7).
6. slice-291 Risiko 3: Ausgang an „Basis kein Vorfahr ⇒ keine Umbenennung“ binden (F-5).
7. welle-96 §6 zweiter Punkt: der eine Fall, den beide Schalter schärfen (F-3).
8. slice-292 §3 Testfälle ergänzen: unveränderter Link auf umbenannte Datei mit erlaubtem Anhang; Tausch
   A ↔ B; beide Schalter mit gleichnamiger anderer Datei; parallele Linien; Merge als HEAD mit einer
   umbenennenden Seite; Ende in HEAD gelöscht; Ausführbar-Bit; Symlink gegen Datei; Z. 137–138 „Basis kein
   Vorfahr ⇒ wie ohne Schalter“; Z. 144 zehn Läufe (F-2 bis F-6, F-8).
9. slice-292 §3 Zeile Port/Adapter: Vertrag aus ADR-0115 Entscheidung 13 (F-6).
10. slice-292 DoD-Punkt 3 und §3: Schalter im Profil im selben Commit wie ADR-0115 `Accepted`; AGENTS.md §3.5
    und harness/README.md §Traceability rules nennen ADR-0115 neben ADR-0103; READMEs aus §3 in die
    Release-Prep (F-9, F-10).
