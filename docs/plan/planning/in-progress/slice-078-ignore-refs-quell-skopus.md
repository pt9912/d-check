# Slice slice-078: `ignore-refs` mit Quell-Skopus (`in:`) — Referenz-Ventil über beide Achsen

**Status:** in-progress (welle-61). **§4-Vorfrage entschieden** (Auftraggeber,
2026-07-18): Option (a) — das Ventil wird als **neues geteiltes Bereichskürzel** in
Lastenheft §3 deklariert (Ziel-Achsen-Pendant zu
[`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)),
`codepaths.ignore-refs` bleibt Alias. Doc-first: Lastenheft → ADR → Spezifikation
führen, Code folgt.

**Bezug:** **Change Request** eines Konsumenten (`ai-harness-course`), eingereicht
2026-07-17, Design nach zwei Rückfragen verfeinert (§2.1). Betrifft das Ventil
`ignore-refs`, das heute **in**
[`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
wohnt und modul-lokal zu `codepaths` ist; der CR verlangt es für
[`DC-FA-LINK-001`](../../../../spec/lastenheft.md#dc-fa-link-001--lokale-link--und-bildreferenzen-modul-links)
und
[`DC-FA-ANCH-001`](../../../../spec/lastenheft.md#dc-fa-anch-001--heading-anker-validierung-modul-anchors)
mit. **Noch kein ADR und noch keine Lastenheft-Änderung** — beides folgt, wenn die
Vorfrage aus §4 entschieden ist. IDs werden nicht ad hoc vergeben
([`AGENTS.md`](../../../../AGENTS.md) §5).

**Autor:** pt9912 (CR: Konsument `ai-harness-course`). **Datum:** 2026-07-17.

---

## 1. Ziel

Ein Verzeichnis mit **Template-Dateien** enthält zwei Klassen von Referenzen, die
d-check heute nicht auseinanderhalten kann:

| Klasse | im Quell-Repo | im Ziel-Repo |
|---|---|---|
| **Ziel-Repo-Pfade** (Template-Platzhalter) | lösen **nicht** auf (das Template liegt nicht an Zielposition) | real |
| **Kurs-/Doku-Verweise** | lösen **auf**, werden beim Release-Bau auf eine getaggte Blob-URL gepinnt | verschwinden (absolute URL) |

Weil die erste Klasse massenhaft `target-missing` erzeugt, steht beim Konsumenten
das **ganze Verzeichnis** in `scan.ignore` — „by-design symbolisch". **Damit ist
auch die zweite Klasse ungeprüft**, und genau die braucht das Gate: ihre Auflösung
wird beim Release **unveränderlich in ein ausgeliefertes Artefakt eingefroren**.
Wer eine Überschrift umbenennt und das Template vergisst, merkt es nie; der tote
Anker wird mit ausgeliefert.

Gemessen im Konsumenten-Repo (Template-Verzeichnis aus `scan.ignore` entfernt,
gegen ein digest-gepinntes d-check-Image):

```text
42 Findings — 37 target-missing, 5 codepath-missing, 0 anchor
davon:  8  echte Platzhalter   (CO-<NNN>-<titel>.md, #<anker>, NNNN, …)
       34  im Ziel-Repo real   (../spec/lastenheft.md, ../AGENTS.md,
                                harness/conventions.md, in-progress/roadmap.md, …)
```

Die 63 auflösenden Verweise (39× in den Kurs, 24× Template-intern) sind aktuell
alle intakt, Pfad **und** Anker — die Prüfung wäre heute grün. Sie fehlt trotzdem,
und das ist der Punkt: **der Status quo ist nicht verifiziert, sondern ungeprüft.**

## 2. Entscheidungen / Regel

Fehlt: „Referenz auf **Y** nicht prüfen, **wenn sie in X steht**."

- **Kein neuer Key daneben, sondern `in:` an `ignore-refs`.** Der Schlüssel nimmt
  einen Glob auf den Pfad der **Quelldatei** (die, in der die Referenz steht).
- **Zwei Felder statt Negations-Syntax:** `refs` (Globs auf das **aufgelöste
  Ziel**) und `keep` (Ausnahmen). Semantik: **ignorieren, wenn `refs` matcht und
  `keep` nicht** — `keep` gewinnt unbedingt und **reihenfolge-unabhängig**.
- **Ziel-Globs matchen den aufgelösten Pfad**, so wie das jeweilige Modul heute
  schon auflöst — keine neue Auflösungs-Semantik. Das ist die Ventil-Parität, die
  [ADR-0030](../../adr/0030-tracked-referenz-ziele.md) bereits entschieden hat
  (Befund-`target` = aufgelöster Pfad).
- **Additiv:** ohne den Key ändert sich nichts. Gilt für `links`, `anchors`,
  `codepaths`.

Form (Skizze, nicht normativ — die Normierung folgt mit der Spezifikation):

```yaml
ignore-refs:
  - in: "lab/templates/**"        # Glob auf die QUELLDATEI
    refs: ["lab/templates/**", "tools/check_*.py"]
    keep: ["lab/templates/**/*.template.md", "lab/templates/README.md"]
```

### 2.1 Die zwei Korrekturen am eingereichten CR

Der CR schlug ursprünglich einen eigenen Key `ignore-refs-in` mit
`!`-Negation und gitignore-Semantik vor. Beides wurde auf Rückfrage vom
Konsumenten selbst korrigiert — die Korrekturen sind der Grund, warum dieser
Slice die Fassung von unten trägt und nicht die eingereichte.

**Frage 1 — Ist die `!`-Negation als zweiter Glob-Dialekt gewollt, oder reichen
zwei Felder ohne neue Sprache?**

> Nein, `!` war Verpackung — zwei Felder reichen und sind besser. Der Ausschlag
> ist, dass die Messungen die Zwei-Feld-Semantik **schon implementiert haben**:
> die entscheidende Zeile der Simulation war `ignore if match(refs) and not
> match(keep)` — also `refs ∧ ¬keep`, `keep` gewinnt unbedingt und
> reihenfolge-unabhängig. Das ist gerade **nicht** gitignore-Last-Match. Die
> Zahlen (38 ignoriert / 0 blind / 63 geprüft) belegen damit `ignore`/`keep`,
> nicht `!`; der CR-Text hatte dieselbe Semantik nur in `!`-Syntax gegossen und
> dann fälschlich gitignore-Ordnung drangeschrieben — er war an der Stelle
> **schlicht falsch beschriftet**.
>
> Der Preis ist real, aber hier null: zwei Felder können nicht alternieren
> (`ignore` → `keep` → wieder `ignore`); ordnungsbasierte Globs können das.
> Gemessen: **alle 24** von `keep` zurückgeholten Ziele existieren, **kein
> einziger** Fall braucht ein Re-Ignore. Für den Verlust an Ausdruckskraft gibt es
> hier also keinen Anwendungsfall — und er kauft Ordnungsunabhängigkeit in einer
> YAML-Liste, wo Leser Reihenfolge ohnehin nicht als bedeutungstragend lesen.

**Frage 2 — Trägt die vierte Ventil-Achse, oder lässt sich `ignore-refs`
skopieren statt einen Key danebenzustellen?**

> Richtig — es ist **keine vierte Achse, sondern das Kreuzprodukt der zwei
> vorhandenen**. Die Ventile heute sind zwei Achsen plus ein Notausgang:
> `scan.ignore` skopiert über die **Quelle** (Datei X gar nicht scannen),
> `ignore-refs` über das **Ziel** (Referenz auf Y nicht prüfen), der Zeilen-Marker
> ist die lokale Ausnahme. Was fehlt, ist „Referenz auf Y nicht prüfen, wenn sie
> in X steht" — die **Kombination** beider Achsen, keine neue. Ein
> `ignore-refs-in` danebenzustellen hätte die Ziel-Achse **dupliziert**; ein `in:`
> an `ignore-refs` ist dieselbe Achse mit Skopus.
>
> Eine Bedingung hängt daran: `ignore-refs` sitzt heute unter `codepaths:` und ist
> damit **modul-lokal**. 37 der 42 Findings sind `links`, 5 `codepaths` — das
> heutige `ignore-refs` **erreicht meinen Fall also gar nicht**. Skopieren hilft
> nur, wenn der Key zugleich **nach oben wandert** und für
> `links`/`anchors`/`codepaths` gilt. `codepaths.ignore-refs` könnte als **Alias**
> bestehen bleiben, dann ist es **kein Config-Bruch**.

### 2.2 Warum die vorhandenen Mechanismen nicht reichen

- **`scan.ignore`** (Status quo) ist alles-oder-nichts pro Datei. Es opfert die
  prüfbare Klasse, um die symbolische loszuwerden.
- **Der Zeilen-Marker `d-check:ignore`** ist hier **aktiv schädlich**:
  Template-Dateien werden vom Adopter in sein Repo **kopiert**, und die Marker
  reisen mit. An Zielposition sind 34 der 42 Referenzen **gültig und prüfbar**;
  der Marker unterdrückt dort dauerhaft eine echte Prüfung. Der Konsument
  dokumentiert Adoptern sogar ausdrücklich, Marker beim Ausfüllen **zu behalten**.
- **`codepaths.ignore-refs`** ist global (nicht skopiert) und modul-lokal: die
  Skelett-Pfade dort einzutragen machte dieselben Pfade **repo-weit** blind, und
  es deckt `links`/`anchors` gar nicht ab — also 37 der 42 Findings nicht.

### 2.3 Warum `keep` konstitutiv ist

Das Template-Verzeichnis mischt beide Klassen. Gegen den realen Bestand
durchgerechnet (Markdown-Links):

| Musterset | ignoriert | davon real (= fälschlich blind) | geprüft | ERROR |
|---|---|---|---|---|
| nur `refs` (ohne `keep`) | 62 | **24** | 39 | 0 |
| `refs` **∧ ¬**`keep` | 38 | **0** | **63** | 0 |

Ohne `keep` kostet das Feature 24 real prüfbare Verweise — es tauscht nur eine
Blindstelle gegen eine andere. Mit `keep` ist der Schnitt exakt: alles Symbolische
ignoriert, **nichts** Reales blind, Gate grün.

## 3. Definition of Done

- [x] **CR erfasst** samt Messungen und den zwei Korrekturen (§2.1).
- [x] **Vorfrage entschieden** (§4, Auftraggeber 2026-07-18): Option (a) — neues
  geteiltes Bereichskürzel, nicht Änderung der drei bestehenden Anforderungen.
- [x] **Lastenheft-CR:** [`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus)
  nach Schema (Bereich `REF` deklariert), 6 Akzeptanzkriterien
  (Happy/Boundary-`keep`/Negative-Tippfehler/Skopus-Isolation/Alias/Regression),
  Out-of-Scope, Version 0.46.0→0.47.0, Historie-Zeile; das modul-lokale
  `codepaths.ignore-refs` zum Zeiger + Alias refaktoriert, die Module
  `links`/`anchors`/`codepaths` an das Ventil angebunden.
- [x] **ADR:** [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md)
  (Proposed) — Zwei-Feld-Semantik gegen `!`-Negation/gitignore-Ordnung, `keep`
  konstitutiv, `in:` als Skopus (keine 4. Achse), Alias-Pfad + Deprecation-Frage,
  Wirkung (Existenz/Escape/Anker, Symlink bleibt).
- [x] **Spezifikation-`.a`:** [`DC-FA-REF-001.a`](../../../../spec/spezifikation.md#dc-fa-ref-001a--geteiltes-referenz-ventil-ignore-refs)
  — Match-Prädikat, Wirkung, Alias-Semantik, Achsen-Präzedenz gegen
  `scan.ignore`/Zeilen-Marker, leeres `refs` inert; §2-Schema-Keys
  (`ignore-refs[].in`/`refs`/`keep`).
- [ ] **Tests, aus den Abnahmekriterien des CR:**
  - ohne den Key verhält sich d-check unverändert (Regression);
  - **`keep` wirkt** — die 24 im Template-Baum real auflösenden Verweise bleiben
    scharf. **Ohne `keep`-Support ist der CR nicht abgenommen**;
  - **Tippfehler-Test (die eigentliche Regression):** ein verfälschter Pfad in
    einer Template-Datei ⇒ **ERROR**. Eine „ignoriere, was nicht auflöst"-
    Heuristik besteht diesen Test **nicht** — deshalb Muster statt Heuristik;
  - **Anker-Test:** ein verfälschter Anker ⇒ ERROR;
  - **Skopus wirkt nur im Quell-Glob:** dieselben Ziel-Muster bleiben außerhalb
    voll geprüft;
  - `codepaths.ignore-refs` als Alias bleibt grün (kein Config-Bruch).
- [ ] **Realdatenbeleg** gegen das Konsumenten-Repo: 0 Findings bei **63
  tatsächlich geprüften** Verweisen — nicht durch Wegschauen (38 ignoriert, davon
  0 real existierend).
- [ ] **Nutzerdoku** (Handbuch §5/§6 Ventil-Achsen) + CHANGELOG mit dem
  Release-Prep.
- [ ] **Qualität:** unabhängiger, kontext-getrennter Review **vor** dem Release;
  `make gates`/`make ci` grün.

## 4. Risiken / offene Punkte

- **Vorfrage ENTSCHIEDEN (Auftraggeber, 2026-07-18): Option (a).** Das Ventil wird als
  **neues geteiltes Bereichskürzel** in Lastenheft §3 deklariert — das
  Ziel-Achsen-Pendant zu
  [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)
  der Quell-Achse; `links`/`anchors`/`codepaths` verweisen darauf,
  `codepaths.ignore-refs` bleibt Alias (kein Config-Bruch). Damit entfällt die
  Verdreifachung der Ventil-Spezifikation (Option b). Die konkrete Kennung vergibt
  der Lastenheft-CR nach
  [`MR-002`](../../../../harness/conventions.md#mr-002--id-schema-mit-bereichskürzeln-ab-initialer-fassung).
  _Ursprüngliche Vorfrage (zur Nachvollziehbarkeit):_ `ignore-refs` steht heute **in**
  [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
  und ist modul-lokal, der CR verlangt es für drei Module — **(a)** neues Kürzel vs.
  **(b)** Änderung der drei bestehenden Anforderungen.
- **Config-Fläche.** Auch wenn es das Kreuzprodukt der zwei vorhandenen Achsen ist
  und keine dritte: die Fläche wächst (`scan.ignore`, Zeilen-Marker,
  `*-exempt-paths`, `ignore-refs` + `in:`/`keep`). Die Handbuch-Doku muss die
  Achsen **gegeneinander** erklären, nicht nur nebeneinander.
- **Alias-Pfad ist ein Vertrag.** `codepaths.ignore-refs` bestehen zu lassen
  vermeidet den Bruch, verdoppelt aber die Config-Oberfläche dauerhaft. Ob der
  Alias eine Deprecation-Frist bekommt, ist Teil der ADR — nicht dieses Slices.
- **Der Konsument ist heute blind, nicht rot.** Die 63 Verweise sind intakt; das
  Gate fehlt, es schlägt nicht fehl. Dringlichkeit daher niedriger als
  [slice-075](../done/slice-075-komma-kurzform-fail-closed.md) (verfälscht produktiv
  verdrahtete Zahlen) — aber die Blindheit wird beim Release **eingefroren**, und
  das ist der Grund, warum sie nicht beliebig warten kann.
- **Glob-Syntax trägt:** `matchGlob` löst `**` segmentweise auf
  (`internal/hexagon/core/rules/paths.go`) — die vorgeschlagenen Muster
  funktionieren ohne neue Glob-Mechanik.

## 5. Trigger

Konsumenten-CR `ai-harness-course` (2026-07-17), getestet gegen ein
digest-gepinntes d-check-Image. Der Konsument hat seine generischen
Referenzprüfungen 2026-06-12 bewusst von einem eigenen JS-Checker auf d-check
migriert; ein Sonderprüfer nur für dieses Verzeichnis würde diese Konsolidierung
rückgängig machen. Deshalb der CR statt eines lokalen Workarounds.

## 6. Sub-Area-Modus-Begründung

GF (Repo-Default): Die Spezifikation führt, der Code folgt. Hier zusätzlich
verschärft — das Lastenheft führt: der CR ändert einen **Vertrag** (ein Ventil
wandert von einem Modul zu dreien), und dieser Slice hält deshalb an §4 an, statt
Code zu schreiben.

## 7. Closure-Notiz (nach `done/`)

_Ausstehend._
