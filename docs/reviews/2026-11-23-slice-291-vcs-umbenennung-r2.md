# Review-Report: slice-291 — `vcs` folgt einer Umbenennung, Fix-Commits nach R1 (R2) — 2026-11-23

**Review-Art:** Design (Diff gegen Slice-Plan, ADR und Hard Rules, Maintainability;
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Nicht gegen die DoD — die
prüft danach der Verifier.

**Gegenstand:** `slice-291`, Commit-Range `2fa85f5d..746d1863` (acht Commits, nur lokal auf `main`):
`3a2d7e3f` und `6295c2d4` ([ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md), Architect,
`Proposed`), `ee9a6a96` (Plan-Änderung), `4317dc0a` und `18157ed7` (Lastenheft
[`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
`18157ed7` auch Spezifikation), `4c446050` (Spezifikation
§[`DC-FA-VCS-001.a`](../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 und 4, Schema-Zeile), `6ccf5ac5` und `746d1863` (Planungs-Stellen, Planner). Gelesen wurde
nicht nur der Diff, sondern Absatz, Grenzen, alle 24 Kriterien, Spezifikation Schritt 2 und 4 und die
ADR im Stand `746d1863` nebeneinander.

**Skill:** `.harness/skills/reviewer.md` @ 1.20.0 (`d1c37c14`).
**Modell:** claude-opus-5-5 · **Datum:** 2026-11-23.

**Eingangs-Kontext:**

- Report R1 zu `slice-291` ([`2026-11-23-slice-291-vcs-umbenennung-r1.md`](2026-11-23-slice-291-vcs-umbenennung-r1.md)),
  Befunde M-1 bis M-4 und I-1 bis I-3
- Plan `slice-291` §1, §2, §3 samt beiden Plan-Änderungen und Plan-Nachtrag, §6;
  Plan `slice-292` (Testfall-Liste), Welle `welle-96` §Out-of-Scope
- [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) (Entscheidung 1–14, Alternativen, Konsequenzen,
  Fitness Function, Geschichte); [ADR-0103](../plan/adr/0103-adr-gate-laesst-pfad-nachzug-durch.md)
- [`AGENTS.md`](../../AGENTS.md) §3.4, §3.5, §3.8, §5 (Regeln 13, 15), §6 Schritt 4;
  [`MR-032`](../../harness/conventions.md#mr-032) und der Kopf des Lastenhefts;
  `v6.18.1` · `regelwerk/grundlagen-source-precedence.md` §Wann die CR-Pflicht beginnt;
  `v6.18.1` · `regelwerk/modul-08-agentenrollen.md` §Rollen-Sequenz für einen Slice
- Gelesen am Gegenstand: `internal/adapter/driven/git/git.go` (`AllPaths` mit dem Vorfahren-Walk der
  Basis, Index-Zweig im gestagten Modus, `ancestors`, `isAncestor`), `internal/hexagon/core/rules/vcs.go`
  und `vcs_link_targets.go` (eine Auflösungs-Funktion für beide Fassungen)

---

## Stand der R1-Befunde

| R1 | Stand | Beleg |
|---|---|---|
| M-1 (Kriterien ohne Bedingung auf den zweiten Schlüssel) | **geschlossen** | Jedes Negativ-Kriterium, das `core-drift-vcs` an einen Link-Nachzug bindet, trägt „ohne `vcs.ignore-link-targets`“ (Lastenheft Z. 2365, 2366, 2368–2373); „Allgemein“ steht unter seinem Given (Z. 2374); Z. 2376 bindet die Negativ-Kriterien für beide Schlüssel an verschiedene Dateinamen. Ungebunden bleiben nur Kriterien, die mit beiden Schlüsseln unverändert gelten (Z. 2364, 2367, 2378, 2379). Rest: Finding M-3 unten (der „sonst“-Zweig von Z. 2376) |
| M-2 (Grenze „beide Schlüssel“ ohne Bedingung) | **geschlossen** | Lastenheft Z. 2322–2327 und Spezifikation Z. 2171–2177 nennen „wenn die Umbenennung den Dateinamen wechselt“ und den Verzeichnis-Wechsel; Welle-Plan, Plan `slice-291` Entscheid 10 und `slice-292` ziehen nach (`6ccf5ac5`) |
| M-3 (Neubelegung eines Umbenennungs-Pfads) | **geschlossen** für beide R1-Folgen | Folge (a): die Linie endet mit „gelöscht“, A hat keine Umbenennung (Lastenheft Z. 2281–2283, Spezifikation Z. 2086–2089, ADR Entscheidung 11), Kriterium Z. 2372. Folge (b): Quelle nur ein BASE-Pfad, der Zwischenpfad nimmt nicht teil (Z. 2279–2281, Spezifikation Z. 2149–2151, ADR Entscheidung 9), Kriterium Z. 2373; Merge-Fall Z. 2371. Neu eröffnet die Schärfung zwei Stellen: M-1 und M-2 unten |
| M-4 (Zusage übersieht den fail-closed-Fall) | **geschlossen**, Rest LOW | Spezifikation Z. 2159–2163, ADR Z. 116–119 und Alternative K sagen „kein Befund“ und benennen Exit 2. Das exklusive „nur“ in Lastenheft Z. 2374 ist enger als die fail-closed-Regel (L-1) |
| I-1 (Schritt-Definition) | **geschlossen** | Lastenheft Z. 2268–2271 und Historie-Zeile 0.114.0 sagen „Commit mit genau einem Eltern-Commit“, gleich Spezifikation Z. 2059–2061 und ADR Entscheidung 2 |
| I-2 (`Schärft:` nennt Schritt 1) | **geschlossen** | ADR Z. 25: Schritt 2 und 4 |
| I-3 (Bestätigung neben angeglichenem Text) | **geschlossen** | ADR Z. 76–78 trennt Bestätigung, Angleichung ohne erneute Bestätigung und Schärfung nach Review; die Geschichte trägt beide Schärfungen. Ob die Schärfung eine Bestätigung braucht, bleibt eine Auftraggeber-Frage; die Begründung „fällt nur in Richtung Drift“ dafür trägt nur ohne `vcs.ignore-link-targets` (L-3) |

---

## Findings

### M-1 — Ein Pfad, der in einer Kette gelöscht wird und danach Ende einer anderen Umbenennung ist: die Spezifikation schließt ihn aus, Lastenheft-Absatz und ADR nehmen ihn auf

- **Kategorie:** MEDIUM
- **Quelle:** [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  Absatz „Was gleichgesetzt wird“; [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 9
- **Pfad:** `spec/spezifikation.md` · „Ein Pfad, dessen Kette mit „gelöscht“ endet,“ (Z. 2153–2156);
  `spec/lastenheft.md` · „ein Ziel, das auf A oder“ (Z. 2292–2297) und „Ein Pfad, der zugleich Quelle der einen und Ende einer“ (Z. 2290–2291);
  `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „dieses Ende), oder wenn er ein solches Ende ist.“ (Z. 112)
- **Befund:** Eingabe: A liegt in BASE, ein Schritt löscht A, ein späterer benennt C (aus BASE) nach A
  um, eine immutable Datei zieht einen Link von C auf A nach, nur `vcs.follow-renames`. Nach Lastenheft
  und ADR ist C → A eine Umbenennung, A ist ihr Ende und nicht zugleich Quelle (A hat keine
  Umbenennung) — beide Ziele werden zu A, kein Befund. Nach der Spezifikation nimmt A nicht teil, weil
  seine Kette mit „gelöscht“ endet, „auch wenn er neu angelegt wird“: C wird zu Marke und A, A bleibt
  roh — `core-drift-vcs`. Dieselbe Doppelrolle hat der Zwischenpfad (A → B → D, danach C → B): alle drei
  Straten sagen „weder Quelle noch Ende“ bzw. „nimmt nicht teil“ und zugleich „oder ein solches Ende“.
  Beide Lesarten fallen in die sichere Richtung; widersprüchlich ist, was für dieselbe Eingabe
  zugesagt ist.
- **Failure-Szenario:** `slice-292` implementiert nach Schritt 4; ein Test aus dem Lastenheft-Absatz
  („ein Ziel auf C oder A gilt als Ziel auf A“) für das Ersetzen eines gelöschten Dokuments durch eine
  Umbenennung läuft gegen spezifikationskonformen Code rot, oder der Code folgt dem Lastenheft und der
  Spezifikations-Satz ist falsch. Ein Adopter, der eine alte Datei löscht und eine andere an ihren Namen
  verschiebt, bekommt den Nachzug als Drift, obwohl der Absatz ihn zusagt.
- **Warum nicht HIGH:** Kein Code, und keine Lesart ergibt ein Grün, das es ohne den Schlüssel nicht
  gibt.
- **Verifizierbar:** ja, ab `slice-292` — Fixture „A löschen, C → A, Nachzug C → A“ zeigt, welcher Satz
  gilt.
- **Klasse:** pfad-mit-zwei-rollen-in-der-abbildung
- **Träger:** Architect (welche Rolle sticht), danach Implementer

### M-2 — „Löschen und neu anlegen“ steht ohne die Schritt-Bedingung; im selben Commit ist es ein Überschreiben an Ort und Stelle

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 13; Prüffrage 18;
  [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 11
- **Pfad:** `spec/lastenheft.md` · „Löscht die Seite den Pfad und“ (Z. 2340–2341), „B in einem späteren Commit gelöscht und als **andere** Datei neu angelegt“ (Z. 2372),
  Historie-Zeile 0.114.0 · „Löschen und Neuanlegen auf einer Seite ist Mehrdeutigkeit“ (Z. 4750);
  `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „Löscht die andere Seite einen Pfad, den diese umbenennt, und legt ihn neu an“ (Z. 202),
  „Merge mit einer Seite, die den Pfad löscht und neu anlegt“ (Z. 221), „Ende gelöscht und als andere Datei neu angelegt“ (Z. 222)
- **Befund:** Spezifikation (Z. 2086–2094) und ADR Entscheidung 11 knüpfen „gelöscht“ an einen
  **Schritt**, in dem der Pfad verschwindet; „bleibt der Pfad in einem Schritt bestehen und ändert nur
  seinen Inhalt, läuft die Linie weiter“. Ein `git rm` samt Neuanlage im selben Commit hinterlässt in den
  Bäumen einen Pfad, der besteht und seinen Inhalt ändert. Die Merge-Grenze des Lastenhefts, das
  Kriterium Z. 2372 (die Präposition „in einem späteren Commit“ trägt beide Verben), die Historie-Zeile
  und die ADR-Stellen unter Konsequenzen und Fitness Function sagen „löscht und legt neu an ⇒
  Mehrdeutigkeit bzw. Befund“ ohne diese Bedingung. Nur Kriterium Z. 2371 trägt sie („in einem Commit
  löscht und in einem späteren neu anlegt“).
- **Failure-Szenario:** Ein Adopter ersetzt eine umbenannte Datei mit `git rm` und neuer Datei in einem
  Commit und verlässt sich auf die benannte Grenze („Löschen und Neuanlegen ist Mehrdeutigkeit, der
  Nachzug bleibt Drift“); das Modul folgt der Umbenennung, ein Link auf das neue Dokument passiert das
  ADR-Gate. Oder das Fixture zu Z. 2372 wird in einem Commit gebaut und läuft gegen
  spezifikationskonformen Code rot.
- **Warum nicht HIGH:** Das Grün ist die benannte Grenze „umbenennen, dann umschreiben“ (Z. 2313–2318);
  falsch ist, dass der Text es für diese Eingabe bestreitet.
- **Verifizierbar:** ja, ab `slice-292` — Fixture mit Löschung und Neuanlage im selben Commit liefert
  nach der Spezifikation eine Umbenennung.
- **Klasse:** grenz-aussage-ohne-ihre-bedingung (zweites Auftreten in diesem Slice nach R1 M-2)
- **Träger:** Implementer (Lastenheft), Architect (ADR, `Proposed`)

### M-3 — Mit beiden Schlüsseln ist ein Umbiegen auf eine umbenannte Quelle grün, das mit `vcs.ignore-link-targets` allein Drift ist; drei neue Sätze bestreiten es

- **Kategorie:** MEDIUM
- **Quelle:** [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  Kriterium „Boundary (Pfad-Nachzug, ohne `vcs.follow-renames`)“; Prüffrage 18
- **Pfad:** `spec/lastenheft.md` · „sonst gilt die Grenze „gleicher Name, andere Datei“ wie mit `vcs.ignore-link-targets` allein — kein Befund“ (Z. 2376),
  „so grün oder rot wie ohne `vcs.follow-renames`, nie grüner“ (Z. 2331),
  „strenger als `vcs.ignore-link-targets` allein, und nur“ (Z. 2324–2325);
  `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „`ignore-link-targets` allein zurück — so grün oder rot wie ohne `follow-renames`, nie grüner“ (Z. 208)
- **Befund:** Eingabe: beide Schlüssel; ein Link wird von `d1/a.md` (ohne Umbenennung, oder in der Range
  gelöscht) auf `d2/c.md` umgebogen, und `d2/c.md` wird in der Range nach `d3/a.md` umbenannt. Nach der
  Gleichsetzung tragen beide Ziele `a.md` — kein Befund. Mit `vcs.ignore-link-targets` allein stehen
  `a.md` und `c.md` gegeneinander — `core-drift-vcs` (Kriterium Z. 2363: „Datei mit anderem Namen“).
  Z. 2376 nennt das Ergebnis „wie mit `vcs.ignore-link-targets` allein“, Z. 2331 und ADR Z. 208 sagen
  für den gelöschten Pfad „nie grüner“, die Grenze Z. 2322–2327 kennt mit beiden Schlüsseln nur die
  strengere Richtung. Der Fall gehört zur Grenze „gleicher Name, andere Datei“, erreicht über den Namen
  des Endes; benannt ist er nicht.
- **Failure-Szenario:** Ein Adopter mit `vcs.ignore-link-targets` schaltet `vcs.follow-renames` dazu,
  weil die Anforderung sagt, mit beiden werde es nur strenger; ein Link in einer `Accepted`-ADR wird auf
  ein anderes Dokument umgebogen, das in derselben Range den alten Dateinamen bekommt, und
  `make adr-check` bleibt grün, wo es vorher rot war. In `slice-292` widerspricht ein Test zu Z. 2376
  („wie mit `ignore-link-targets` allein“) dem Ergebnis der Normalform.
- **Warum nicht HIGH:** Das neue Ziel ist der alte Pfad einer umbenannten Datei und fehlt in HEAD; den
  toten Link meldet das Modul `links`, wo es mitläuft. Kein Code.
- **Verifizierbar:** ja, ab `slice-292` — Fixture wie oben, einmal mit beiden Schlüsseln, einmal mit
  `vcs.ignore-link-targets` allein.
- **Klasse:** zusage-nie-gruener-ohne-den-paar-fall
- **Träger:** Implementer (Lastenheft), Architect (ADR §Konsequenzen)

### L-1 — „nur im Fall fail-closed (Zwischen-Commit)“ ist enger als die fail-closed-Regel

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 13 („Jedes Objekt, das
  sie dafür braucht und nicht lesen kann, ist ein Fehler“)
- **Pfad:** `spec/lastenheft.md` · „mit ihm rot ohne Befund wird ein Lauf nur im Fall „fail-closed (Zwischen-Commit)““ (Z. 2374)
- **Befund:** Das Kriterium „fail-closed (Zwischen-Commit)“ (Z. 2381) gilt einem Commit oder Baum
  „zwischen BASE und HEAD“. Die Regel selbst (Z. 2305–2306, Spezifikation Z. 2101–2107) gilt jedem
  Objekt, das die Erkennung braucht — auch dem Baum des Eltern-Commits eines Schritts, der vor BASE liegt
  (die Seite eines Merges, die vor BASE abzweigt). `AllPaths` liest ohne den Schlüssel von diesen
  Vorfahren nur die Commits, nicht die Bäume. Ist ein solcher Baum unlesbar, endet der Lauf mit Exit 2
  außerhalb des Falls, den das „nur“ als einzigen nennt.
- **Verifizierbar:** nein
- **Klasse:** zusage-uebersieht-eigenen-fail-closed-fall (Rest von R1 M-4)

### L-2 — Die Historie-Zeile 0.114.0 ist fortgeschrieben statt um eine Zeile ergänzt

- **Kategorie:** LOW
- **Quelle:** [`MR-032`](../../harness/conventions.md#mr-032)
- **Pfad:** `spec/lastenheft.md` · „neunzehn Kriterien kommen hinzu“ (Z. 4750) und Kopf · „**Version:** 0.114.0“ (Z. 3)
- **Befund:** `4317dc0a` und `18157ed7` ändern Anforderungs-Aussagen (Kette, Teilnahme, Merge-Grenze,
  vier neue Kriterien), heben die Version aber nicht und schreiben die Zeile 0.114.0 um. Damit tragen
  `f0622065` und `746d1863` verschiedene Dokumente unter derselben Version. Der R1-Report, der nicht
  überschrieben wird, nennt „Version 0.114.0“ mit fünfzehn Kriterien, die Zeile 0.114.0 sagt jetzt
  neunzehn.
- **Verifizierbar:** nein
- **Klasse:** historie-zeile-ueberschrieben-statt-ergaenzt

### L-3 — Zwei Botschaften behaupten, die Schärfung wirke nur in Richtung Drift; mit beiden Schlüsseln stimmt das nicht

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §5 Regel 15
- **Pfad:** Commit `3a2d7e3f` · „Die Regel nimmt nur Gleichsetzungen zurueck.“; Commit `6295c2d4` · „Faellt nur in Richtung Drift.“
- **Befund:** Mit `vcs.follow-renames` allein stimmt beides: ein Pfad, der nicht mehr teilnimmt, bleibt
  roh. Mit beiden Schlüsseln fällt er auf seinen eigenen Dateinamen zurück statt auf den des Endes. Ein
  Link von ihm auf eine andere Datei mit seinem Dateinamen war nach dem Text vor der Schärfung Drift
  (Name des Endes gegen den eigenen) und ist danach grün — so grün wie mit `vcs.ignore-link-targets`
  allein, aber grüner als der bestätigte Text. Die Botschaften ziehen ihren Schluss über die eine
  Schalter-Kombination hinaus, an der er geprüft ist; die ADR-Einleitung (Z. 76–78) lässt die erneute
  Bestätigung offen, und das ist die Begründung, auf die sie sich stützen würde.
- **Verifizierbar:** nein
- **Klasse:** botschaft-ueber-messung (`BEO-ALL/commit-message-overclaims-work`, zweite Richtung)

### L-4 — Die Plan-Änderung nennt drei neue Kriterien; die vierte und die neue Merge-Grenze stehen in keinem Plan-Text

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §6 Schritt 4; `v6.18.1` · `regelwerk/modul-08-agentenrollen.md`
  §Rollen-Sequenz für einen Slice („keine Rolle springt rückwärts … ohne Übergabe-Artefakt“)
- **Pfad:** `docs/plan/planning/in-progress/slice-291-vcs-umbenennung-anforderung-und-adr.md` · „Neu kommen drei Kriterien hinzu“ (Z. 288)
- **Befund:** `18157ed7` fügt „Negative (Merge-Seite löscht und legt neu an)“ hinzu und schränkt die
  Merge-Grenze auf das Überschreiben an Ort und Stelle ein; davor steht kein Plan-Nachtrag, der Plan
  zählt weiter drei. Die „offene Frage des Implementers“, auf die `6295c2d4` antwortet, liegt in keinem
  Artefakt des Repos. Die Abgrenzung (§1) ist nicht berührt — beides liegt in schon geführten Abschnitten.
- **Verifizierbar:** nein
- **Klasse:** plan-aenderung-zaehlt-den-inhalt-nicht

### L-5 — „Ein Nachzug auf eine Datei, die nur noch in BASE existiert, geht durch“ steht jetzt hinter dem Satz zu `vcs.follow-renames`

- **Kategorie:** LOW
- **Quelle:** [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) Entscheidung 10
- **Pfad:** `spec/spezifikation.md` · „Ein Nachzug auf eine“ … „Datei, die nur noch in BASE existiert, geht durch“ (Z. 2201–2203)
- **Befund:** Der Satz stammt aus der Fassung ohne `vcs.follow-renames` und hängt an „denselben Namen“
  von `vcs.ignore-link-targets`. Seit `a80cf8f6` folgt er dem Satz „Mit ihm gilt …“ und liest sich als
  allgemein. Mit `vcs.follow-renames` allein bleibt ein solches Ziel roh oder wird zum Ende — der
  Nachzug ist Drift; mit beiden Schlüsseln ebenso, wenn die Datei eine umbenannte Quelle mit anderem
  Endnamen ist. Der Satz lag schon im R1-Gegenstand; R1 hat ihn nicht gemeldet.
- **Verifizierbar:** nein
- **Klasse:** grenz-satz-verliert-seinen-schluessel

### I-1 — „Unberührt ist nur …“ lässt eine Linie ohne Klasse, auf der der Pfad nie liegt

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `spec/spezifikation.md` · „Unberührt ist nur eine Linie, auf der der“ (Z. 2093–2094); ebenso Lastenheft Z. 2276–2278, ADR Z. 135–136
- **Befund:** Zweigt eine Merge-Seite vor BASE ab, an einem Commit ohne den Pfad, und berührt sie ihn
  nie, benennt sie ihn weder um noch löscht sie ihn, und „unberührt“ ist sie nach der neuen Definition
  auch nicht. Welche Linie sie ist, sagt keine der drei Stellen. Beide denkbaren Antworten (kein Ende;
  „gelöscht“) fallen nicht ins Grüne.
- **Verifizierbar:** nein
- **Klasse:** linien-klassen-nicht-erschoepfend

### I-2 — Die ADR zählt „gelöscht“ als Ende, Lastenheft und Spezifikation als „ohne Ende“

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Pfad:** `docs/plan/adr/0115-vcs-folgt-umbenennung.md` · „das zählt als eigenes Ende, das in HEAD nie“ (Z. 130–131) und „Genau ein Ende, das in HEAD existiert, ist eine Umbenennung“ (Z. 127)
- **Befund:** Lastenheft (Z. 2282) und Spezifikation (Z. 2088) sagen „endet ohne Ende“. In der ADR wird
  „gelöscht“ ein eigenes Ende; „genau ein Ende, das in HEAD existiert“ lässt sich dann auch als „genau
  eines der Enden existiert in HEAD“ lesen, was für {B, gelöscht} eine Umbenennung ergäbe — der nächste
  Satz der ADR sagt Mehrdeutigkeit. Das Ergebnis der gemeinten Lesart ist überall gleich.
- **Verifizierbar:** nein
- **Klasse:** stratum-definition-ungleich-formuliert

### I-3 — Testfall „Tausch in einem Commit“ begründet sich mit einer Regel, die er nicht fährt (außerhalb der Range)

- **Kategorie:** INFO
- **Quelle:** Prüffrage 20
- **Pfad:** `docs/plan/planning/open/slice-292-vcs-umbenennung-modulkern.md` · „Tausch A ↔ B in einem Commit ⇒ kein Pfad nimmt teil (Quelle und Ende zugleich).“ (Z. 148)
- **Befund:** In einem Commit bestehen A und B vor und nach dem Schritt; nach Spezifikation Z. 2073–2077
  ist das gar keine Umbenennung, die Regel „zugleich Quelle und Ende“ greift erst beim Tausch über einen
  Zwischennamen in mehreren Commits (ADR §Konsequenzen, Z. 188). Der Testfall bliebe grün, wenn die
  Regel fehlte. Nicht im R2-Diff; gemeldet an den Planner von `slice-292`.
- **Verifizierbar:** nein
- **Klasse:** test-begruendung-faehrt-ihre-regel-nicht

---

## Zur Historie-Form (offene Frage)

**Urteil: nach [`MR-032`](../../harness/conventions.md#mr-032) nicht gedeckt; die Form wäre eine eigene
Zeile mit eigener Version (L-2).** MR-032 gilt „jede[r] Änderung am Dokument, solange sein `**Status:**`
unter `Accepted` liegt“, und seine Adaption ist gerade, die Freiheit des Kanons — ändern ohne Bump und
ohne Historie-Zeile — **nicht** zu nutzen. Eine Ausnahme für unveröffentlichte Fassungen kennt weder
MR-032 noch der Kopf des Lastenhefts; der Kopf begründet die Strenge damit, dass die Historie „der
einzige Ort“ ist, an dem eine Anforderungs-Änderung ihren Beleg hinterlässt. „Unveröffentlicht“ ist
zudem kein Zustand, den das Repo führt: `f0622065` ist ein Commit auf `main` und reist mit demselben
Push wie die Fix-Commits — danach gibt es drei verschiedene Dokumente, die 0.114.0 heißen, und ein
eingefrorener Report (R1), der 0.114.0 eine Kriterienzahl zuschreibt, die die Zeile nicht mehr trägt.
Welche Stelle steigt, entscheidet MR-032 nicht; das bleibt beim Implementer. Für die
Spezifikation gilt MR-032 nicht — ihre Datumszeile fortzuschreiben ist ohne Befund.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Widerspruchsfreiheit „Überschreiben an Ort und Stelle“ | geprüft. Lastenheft Z. 2284–2285, 2313–2318, 2338–2340, Spezifikation Z. 2089–2091 und 2199–2201, ADR Entscheidung 11 und §Konsequenzen sagen dasselbe: Pfad bleibt, Inhalt ändert sich, die Linie läuft weiter. Rest M-2 |
| Widerspruchsfreiheit „Tausch“ | geprüft, ohne Befund im Vertrag. In einem Commit keine Umbenennung (beide Pfade bestehen), über einen Zwischennamen Quelle und Ende zugleich; Lastenheft Z. 2319, 2370, Spezifikation Z. 2151–2153, ADR §Konsequenzen gleich. Die neue Zwischenpfad-Regel ändert das Ergebnis nicht. Rest I-3 (Plan `slice-292`) |
| Widerspruchsfreiheit „Zwischenpfad“, „Neubelegung linear“ | geprüft. Beide R1-Folgen ergeben nach allen drei Straten Drift (Kriterien Z. 2372, 2373). Rest M-1 (Doppelrolle mit einem Ende) |
| Widerspruchsfreiheit „beide Schlüssel“ | geprüft. Grenze, Kriterium Z. 2375 und Z. 2376, Spezifikation Z. 2145–2147 und 2171–2177 tragen die Dateinamen-Bedingung. Z. 2376 bindet nur die Kriterien „oben“; Z. 2377 und Z. 2380 stehen darunter, sind aber selbst „ohne `vcs.ignore-link-targets`“ gefasst und widersprechen nichts. Rest M-3 |
| Richtung ohne `vcs.ignore-link-targets` | geprüft, ohne Befund. Die Normierung ist eine Funktion des Ziels auf beiden Seiten; gleicher Rohtext bleibt gleich, und eine entzogene Teilnahme lässt ein Ziel roh — keine der Schärfungen (BASE-Quelle, Zwischenpfad, „gelöscht“, Merge-Seite) erzeugt ein Grün, das es ohne den Schlüssel nicht gibt, und keinen Befund, den es ohne ihn nicht gibt. Exit 2 ist die einzige benannte Ausnahme; ihr Umfang ist L-1 |
| Richtung mit beiden Schlüsseln | geprüft. Die Schärfungen fallen auf `vcs.ignore-link-targets` allein zurück, nie darüber hinaus (Rest L-3 zur Botschaft); der vorbestehende Paar-Fall M-3 ist grüner als `vcs.ignore-link-targets` allein |
| Referenz-Richtung ([`AGENTS.md`](../../AGENTS.md) §3.4) | geprüft, ohne Befund. In den hinzugefügten Zeilen beider Spec-Straten kein `ADR-`, kein `slice-`, kein `welle-`, kein Commit-Hash, kein Link nach `docs/plan/adr/`, kein „Review“ (grep über die `+`-Zeilen leer); „Begründung in begleitender ADR“ ist die gelebte Form |
| ADR-Form ([`AGENTS.md`](../../AGENTS.md) §3.5) | geprüft, ohne Befund. `Proposed`, nicht immutabel; zwei Geschichte-Zeilen mit Provenance-Marker auf `slice-291` als Herkunft, nicht als Begründung; Alternative K, Konsequenzen und Fitness Function nachgezogen |
| Grenz-Aussagen gegen den Code (Regel 13) | geprüft, ohne Befund. `AllPaths` liest ohne Schlüssel die Vorfahren der Basis (Commits) und die beiden Bäume bzw. jeden Index-Eintrag, keinen Commit zwischen BASE und HEAD — Kriterium Z. 2381/2382 und ADR Entscheidung 13 („im Index nimmt `AllPaths` jeden Eintrag“) stimmen. Eine Auflösungs-Funktion für beide Fassungen (`vcs.go`, `vcs_link_targets.go`) trägt die Seiten-Gleichheit |
| Kriterien-Zahl | geprüft, ohne Befund. 24 Kriterien (Lastenheft Z. 2359–2382), davon fünf aus dem Bestand und neunzehn neu (Z. 2364–2382) — wie die Historie-Zeile sagt; Zwischenstand „achtzehn“ in `4317dc0a` stimmte für jenen Commit |
| Abgrenzung des Plans (§1) | geprüft, ohne Befund. Der Diff berührt ADR, Lastenheft, Spezifikation und drei Planungs-Dokumente; kein Code, kein `.d-check.yml`, kein README/CHANGELOG/Handbuch |
| Commit-Disziplin: Kennungen | geprüft, ohne Befund. Jede Botschaft nennt `slice-291` oder `slice-292` und [`DC-FA-VCS-001`](../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in); die ADR-Commits und `746d1863` auch [ADR-0115](../plan/adr/0115-vcs-folgt-umbenennung.md) |
| Commit-Disziplin: Reihenfolge | geprüft. ADR (`3a2d7e3f`) vor der Plan-Änderung (`ee9a6a96`), diese vor dem Spec-Inhalt (`4317dc0a`, `4c446050`); Architect-Antwort `6295c2d4` vor `18157ed7`. Rest L-4 |
| Commit-Disziplin: Botschaften (Regel 15) | geprüft. `4317dc0a`, `4c446050`, `6ccf5ac5`, `18157ed7` und `746d1863` decken sich mit ihrem Diff; Rest L-3 |
| R1-Treffer in Planungs-Dokumenten | geprüft, ohne Befund. Welle `welle-96` Z. 128–131, Plan `slice-291` Entscheid 10, `slice-292` Z. 149–151 tragen die Dateinamen-Bedingung; `slice-292` führt die drei Neubelegungs-Fälle als Testfälle |

## Kategorie-Summary

| Kategorie | Anzahl | IDs |
|---|---|---|
| HIGH | 0 | — |
| MEDIUM | 3 | M-1, M-2, M-3 |
| LOW | 5 | L-1, L-2, L-3, L-4, L-5 |
| INFO | 3 | I-1, I-2, I-3 |

## Verdikt

**Nachbessern vor Closure.** Kein HIGH. Alle sieben R1-Befunde sind geschlossen. M-4 bleibt als L-1
übrig, I-3 hängt an L-3. Die drei MEDIUM liegen in dem, was die Schärfung neu
eröffnet oder neu behauptet: eine Doppelrolle, für die Spezifikation und Lastenheft/ADR Verschiedenes
zusagen (M-1); eine Grenze, der die Schritt-Bedingung fehlt, die Spezifikation und ADR-Entscheidung
tragen (M-2, dieselbe Klasse wie R1 M-2); und drei neue Sätze, die mit beiden Schlüsseln „nie grüner“
zusagen, wo der Paar-Fall grüner ist (M-3). Ohne `vcs.ignore-link-targets` öffnet keine Schärfung ein
Grün und keinen Befund, den es ohne den Schlüssel nicht gibt. M-1 und M-3 berühren Entscheidung 9 und
die Konsequenzen der ADR und gehen zuerst an den Architect. Plan-Risiko 1 bleibt **eingetreten** (R1);
die Klasse `grenz-aussage-ohne-ihre-bedingung` steht in diesem Slice bei 2×. Zur Historie-Form siehe
den eigenen Abschnitt oben: nach MR-032 eine eigene Zeile.
