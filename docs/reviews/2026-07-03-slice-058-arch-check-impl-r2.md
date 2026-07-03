# Review-Report — slice-058 (arch-check via a-check) — Implementierungs-Review (R2)

**Datum:** 2026-07-03
**Reviewer-Rolle:** unabhängig/adversarial, Fokus **Implementierung**
(Commit `640001c`, Range `4d753c2..640001c`).
**Gegenstand:** Umstellung von `make arch-check` auf das digest-gepinnte
a-check-Image — `.a-check.yml` + `a-check.mk` (neu), Makefile-/Dockerfile-Umbau,
`git rm tools/arch-check.sh` + fünfter `codepaths.ignore-refs`-Tombstone,
Doku-Currency (AGENTS §4, Sensors-Tabelle, `.golangci.yml`) — geprüft gegen die
**Paritäts-Referenz** `git show 4d753c2:tools/arch-check.sh` (R1–R6), den
Slice-Plan (§2 bindend), [ADR-0029](../plan/adr/0029-arch-check-via-a-check.md)
und den realen a-check-v0.8.0-Stand (Schwester-Repo-Checkout: `rules.go`,
`extract.go`, `config.go`, `cli.go`, Lastenheft `AC-FA-RULE-003`/`AC-FA-CONF-001`).
**Baseline:** `.harness/skills/reviewer.md` v1.2.0, `AGENTS.md` §3,
Plan-Review R1 (`2026-07-03-slice-058-arch-check-plan-r1.md`) als Vorgeschichte.
**NICHT geprüft:** DoD-Abhakung (Verifikations-Rolle), a-check-Repo-interne
Qualität (eigener Harness drüben), CI-Pipeline-Lauf (nur lokale Gate-Läufe).

## Verifikations-Proben (eigene, adversarial — Baum nach jeder Probe restauriert)

Alle Läufe über `make arch-check` (Make-Verdrahtung mitgeprüft, nicht rohes
`docker run`); Abschluss-Zustand: `git status` leer, sauberer Baum grün.

1. **CLEAN:** sauberer Baum → grün (`gesamt: 0 Befund(e)`, Exit 0).
2. **net-Familie im Kern:** `net/textproto` **und** `net/http/httputil` in
   `internal/hexagon/core/model/config.go` → **GRÜN** (Alt-Skript R1: rot) —
   Rest-Delta real, siehe MEDIUM-1.
3. **R1-yaml-Zweig:** `gopkg.in/yaml.v3` in `core/model` → rot,
   `core-impurity: Kern importiert gopkg.in/yaml.v3` (Zweig war in der
   Implementierer-Matrix ungeprobt; jetzt verriegelt).
4. **R6a-Zweitzweig:** `core/model` importiert `core/app` → rot,
   `core-impurity` (Matrix probte nur model→rules; jetzt verriegelt).
5. **R2b-forbid-Flag:** go-git-Import in `internal/adapter/driving/cli/cli.go`
   → rot, `tech-leak … (composition_root: forbid)` (die CLI-forbid-Proben der
   Matrix deckten nur `net/http` und yaml; jetzt verriegelt).
6. **Over-Coverage-Beleg:** `port/driven/filesystem.go` importiert `model` →
   rot, `wrong-direction: ports -> model` (Alt-Skript: erlaubt) — siehe INFO-1.
7. **Substring-Falle:** fiktiver Adapter `internal/adapter/driven/gitlab/`
   importiert go-git → **GRÜN** (Alt-Skript: rot, `rel != …/git`) — siehe
   MEDIUM-2.
8. **Fail-closed (fehlende Config):** `.a-check.yml` temporär entfernt →
   `a-check: open /src/.a-check.yml: no such file or directory`, make rot
   (Fehler 2).
9. **Fail-closed (kaputte Config):** unbekannter Key angehängt →
   `field unbekannter_key not found in type config.yamlConfig`, make rot
   (Fehler 2; `KnownFields(true)` im a-check-Config-Adapter quellen-belegt).
10. **gate-consistency:** Lauf grün („Doku ↔ Makefile konsistent …
    Selbsttest gefeuert"); zusätzlich per Grep: kein `make a-check` in einer
    Doku-Tabelle, `arch-check` bleibt in der Datei `Makefile` definiert.

Quellen-Proben (Lektüre statt Lauf): a-check `internal/cli/cli.go` (Exit 0/1/2;
das dort generierte `--print-mk`-Fragment ist mit dem eingecheckten `a-check.mk`
in Target **und** Recipe identisch, inkl. Digest-Pin `a1c9c4d6…` = offizieller
v0.8.0-Release-Pin — Divergenz nur in Kommentaren); `rules.go`
(`matchTech` = Erst-Treffer in Deklarationsreihenfolge; `inAdapter` =
`strings.Contains`-Substring; `composition_root: forbid` hält `tech-leak` in
der Root aktiv); `extract.go` (`exclude` wirkt **vor** der Extraktion; WalkDir
skippt nur `.git`); Import-Graph aller Nicht-Test-Go-Dateien des Repos gegen
`layers`/`edges` abgeglichen.

---

## Findings

### MEDIUM-1 — R1-`net`-Familie nicht enumeriert: `net/*`-Subpakete (außer `net/http`) im Kern jetzt grün; das dokumentierte Rest-Delta ist per Config schließbar und seine Begründung trägt nicht

- **kategorie:** MEDIUM
- **quelle:** slice-058 §2 (bindend: „`net`-Familie regex-verankert
  enumeriert"), ADR-0029 (§Entscheidung, gleicher Wortlaut), ADR-0005 R1;
  `DC-QA-04`-Analogon (Erkennungs-Differenz zur Alt-Mechanik) im Gate-Pfad
- **pfad:** `.a-check.yml:59-64` (Rest-Delta-Kommentar + tech-Liste)
- **befund:** Das Alt-Skript verbot im Kern `net` **und** `net/*` (einzige
  Ausnahme `net/url`); die Config trägt nur `^net$` und `^net/http$`. Probe 2:
  `net/textproto` **und** `net/http/httputil` in `core/model` bleiben grün —
  das Skript wäre rot. Das Delta ist im Kommentar deklariert, aber (a) §2 des
  Plans verlangt die Enumeration der Familie ausdrücklich, (b) die Begründung
  „RE2 kennt kein Lookahead" ist nicht tragfähig: eine Enumeration verbotener
  Subpakete braucht kein Lookahead, und `matchTech` ist Erst-Treffer in
  Deklarationsreihenfolge (a-check `rules.go`, ADR-0015 dort) — eine
  vorangestellte `^net/url$`-Zeile kann die Ausnahme vor einem breiten
  `^net/`-Muster tragen, (c) die Plan-eigene Regel lässt dokumentierte
  Rest-Deltas nur für **nicht** per Config schließbare Fälle zu (sonst CR ans
  a-check-Lastenheft) — dieses ist schließbar. Keine HIGH-Einstufung, weil das
  Delta offen deklariert ist (kein **stilles** Grün) und die betroffene
  Import-Klasse im Ist-Baum nicht vorkommt.
- **verifizierbar:** ja — Probe 2 reproduziert es ohne Mutation der Config;
  nach einer Korrektur muss dieselbe Injektion rot laufen.

### MEDIUM-2 — `tech.adapter` matcht per Substring: ein künftiger Adapter mit präfix-teilendem Namen erbt die Tech-Kapsel still (go-git in `driven/gitlab` grün, Skript rot)

- **kategorie:** MEDIUM (LOW-Anker „latente Wartungsfalle, zündet bei
  künftigem Edit" + eine Stufe Kontext-Eskalation: Beobachtung liegt im
  Gate-Pfad)
- **quelle:** a-check `rules.go` (`inAdapter` → `strings.Contains`, im
  Gegensatz zum segment-bewussten `segIndex` der Layer-Auflösung);
  Paritäts-Referenz `tools/arch-check.sh` (exakte Paket-Gleichheit
  `rel != internal/adapter/driven/git`)
- **pfad:** `.a-check.yml:45-64` (alle `adapter:`-Werte ohne
  Segment-Abschluss)
- **befund:** Die `adapter:`-Pfade werden als Substring gegen den Dateipfad
  geprüft. Probe 7: ein fiktiver Adapter
  `internal/adapter/driven/gitlab/gitlab.go` mit go-git-Import bleibt
  **grün**, weil der Pfad den Substring `internal/adapter/driven/git`
  enthält — das Skript zog hier rot. Dieselbe Klasse trifft `fs`
  (`fsnotify`…), `configyaml`, `report` und `httpcheck` bei
  präfix-teilenden Namen (realistisch: `git`/`github`/`gitlab`). Zündet erst
  bei einem künftigen Adapter, ist dann aber ein unbemerktes Loch genau in der
  Kapsel-Regel; keine Probe der Matrix (und kein heutiger Baum-Zustand) fängt
  es. Kehrseite derselben Mechanik, policy-konform und nur der Vollständigkeit
  halber notiert: Sub-Pakete eines Adapters (z. B. `driven/git/cache/`) zählen
  neu zur Kapsel — das Skript verlangte exakte Paket-Gleichheit und wäre dort
  rot gewesen; die „eine Tür"-Lesart von ADR-0024 deckt das.
- **verifizierbar:** ja — Probe 7 (Verzeichnis + eine Datei, kein
  Config-Edit); nach einer Härtung muss dieselbe Injektion rot laufen.

### LOW-1 — Proben-Matrix zählt nicht vollständig „je Verbotszweig": R1-yaml, R6a model→app und das forbid-Flag des go-git-Eintrags blieben ungeprobt; die yaml@report-Allow-Gegenprobe existiert nur implizit

- **kategorie:** LOW
- **quelle:** slice-058 §2/DoD (Proben-Zählung „je Verbotszweig, nicht je
  Regel-Nummer", slice-057-R3-Lehre; Allow-Gegenproben „`net/url` im Kern,
  yaml im report-Adapter" ausdrücklich benannt)
- **pfad:** `docs/plan/planning/in-progress/slice-058-arch-check-via-a-check.md:119-123`
  (der Beleg-Anspruch, den die Matrix einlösen soll); Proben-Log des
  Implementierers (13 Proben)
- **befund:** Die Matrix splittet R2a/R2b korrekt und probt zwei der drei
  `composition_root: forbid`-Einträge — aber der Skript-R1-Fall hat drei
  Verbotszweig-Klassen (I/O-Stdlib, Adapter-Import, yaml) und R6-model zwei
  Ziele (rules, app); ungeprobt blieben R1-yaml-im-Kern, model→app und das
  forbid-Flag des go-git-Eintrags (Flag ist Per-Eintrag-Daten, nicht geteilte
  Mechanik). Die yaml-im-report-Allow-Gegenprobe der DoD fehlt als explizite
  Probe; sie ist nur implizit durch den grünen CLEAN-Lauf gedeckt
  (`report.go` importiert yaml real). Meine Nach-Proben 3–5 liefen alle rot
  mit korrektem Grund-Code — es besteht also **kein** reales Gate-Loch, der
  Befund betrifft die Beleg-Vollständigkeit gegenüber der Plan-eigenen
  Zählregel.
- **verifizierbar:** ja — Proben 3–5 dieses Reports.

### INFO-1 — Richtungs-Wechsel Blacklist→Allowlist: die `edges` verbieten mehr als R1–R6; die Über-Deckungs-Richtung ist nirgends als Delta benannt

- **kategorie:** INFO
- **quelle:** Maintainability; ADR-0029 (deklariert „edges als
  Richtungs-Allowlist", listet Rest-Deltas aber nur in der
  Unter-Deckungs-Richtung)
- **pfad:** `.a-check.yml:30-38`
- **befund:** Probe 6: ein Port, der `model` importiert, läuft rot
  (`wrong-direction`) — das Skript erlaubte das (R1 verbot Ports nur
  I/O/Adapter/yaml). Gleiche Klasse: `coretest`→`model`, `adapters`→`rules`,
  `app`→`coretest`. Fail-closed und hexagon-konsistent, im Ist-Baum ohne
  Wirkung (CLEAN grün, realer Import-Graph vollständig von den `edges`
  gedeckt) — aber der nächste legitime Port-/Helfer-Ausbau läuft rot auf,
  ohne dass ein R1–R6-Verstoß vorliegt, und „Parität belegt" (Sensors-Zeile)
  deckt nur die Unter-Deckungs-Richtung. Der ADR-0029-Re-Eval-Trigger
  („Modul-Layout ändert sich → Config-Anpassung") fängt den Pflege-Fall,
  benennt den Richtungs-Unterschied aber nicht. Die analoge Über-Deckung der
  `tech`-Bannliste (os/-Familie repo-weit statt nur im Kern) ist dagegen im
  Config-Kommentar deklariert — dort ohne Befund.
- **verifizierbar:** ja — Probe 6.

### INFO-2 — Scan-Scope umfasst künftige Nicht-Paket-Go-Dateien (z. B. testdata-Fixtures), die `go list` nie sah

- **kategorie:** INFO
- **quelle:** Maintainability; a-check `extract.go` (WalkDir über alle
  `languages`-Glob-Treffer, nur `.git` geskippt)
- **pfad:** `.a-check.yml:11-17` (`languages`/`exclude`)
- **befund:** `exclude` nimmt nur `**/*_test.go` aus; eine künftige
  `.go`-Fixture unter `testdata/` (heute existiert keine — verifiziert) würde
  gescannt und z. B. bei einem `os`-Import falsch-rot, während das
  paket-basierte Skript sie nie sah. Fail-closed-Richtung, kein Loch; als
  dokumentationswürdige Annahme notiert.
- **verifizierbar:** ja — `testdata/`-Fixture mit `import "os"` anlegen,
  `make arch-check` rot.

---

## Negativbefunde (geprüft, ohne Befund)

- **Fail-closed/Silent-Grün:** fehlende Config → Exit 2 (Probe 8); kaputte
  Config → Exit 2 via `KnownFields(true)` + Pflichtblock-Validierung
  (Probe 9, Quelle a-check `config.go`); Befunde → Exit 1 → make rot (alle
  Rot-Proben); nicht pullbares Image → `docker run` schlägt fehl → make rot;
  fehlendes `a-check.mk` → `include` ist fatal; das Delegations-Konstrukt
  `arch-check: a-check` (Prerequisite ohne Recipe) propagiert den Fehler
  (Proben 3–5: `make: *** [a-check.mk:15: a-check] Fehler 1`). Kein
  Silent-Grün-Pfad gefunden.
- **gate-consistency-Verdrahtung:** Richtung 1+2 grün (Probe 10);
  `arch-check` in der Datei `Makefile` definiert (Parser sieht es), das
  Fragment-Target `a-check` steht in **keiner** Doku-Tabelle (Richtung 1
  kann nicht feuern), Fragment-Name/-Recipe gegenüber `--print-mk`
  unverändert (nur Kommentar-Divergenz — Pin-Hebungs-Hinweis im
  Fragment-Kopf vorhanden). Ohne Befund.
- **Paritäts-Kern (Unter-Deckung, außer MEDIUM-1/-2):** alle
  Skript-Verbotszweige systematisch gegen die Config abgeglichen — R1
  (os, os/*, net exakt, syscall, io/fs, Adapter-Import, yaml; für model als
  `domain`, rules/app/coretest als `app`, ports als `port`), R2a/R2b inkl.
  CLI/cmd via `composition_root: forbid`, R3-Doppel-Erlaubnis
  (Adapter-Liste), R4-Dreifach-Zone (fs-Adapter + Root-Default-allow), R5
  (lateral, Sub-Paket-Semantik deckungsgleich), R6 beide Richtungen. Die
  Implementierer-Matrix (13 Proben, Log eingesehen) prüft Exit **und**
  Grund-Code; die exclude-Mutations-Probe belegt Load-bearing. Meine
  Zusatz-Proben 3–5 schließen die Rest-Zweige rot. Ohne weiteren Befund.
- **Layer-Globs vs. reale Verzeichnisse:** alle sechs Layer decken
  existierende Verzeichnisse (inkl. `coretest` unter `core/`, `port/driven`
  unter `port/**`); keine `.go`-Dateien direkt unter `core/`; keine toten
  Globs; `edges` decken den realen Import-Graphen vollständig (u. a.
  `report`→`app`+`model`, `coretest`→`ports`, `httpcheck`→`net`). Ohne
  Befund.
- **Hermetik-/Pin-Politik (DC-QA-03, ADR-0010/0011):** Lauf
  `--network none` + `:ro`-Mount; Digest-Pin identisch mit der
  v0.8.0-Release-Konstante im Schwester-Repo; die drei
  Umstellungs-Vorbedingungen (tech-Adapter-Liste, `composition_root:
  forbid`, `exclude`) sind im v0.8.0-Quelltext real geliefert und geprobt;
  `make versions` weist den Pin aus; `?=`-Override analog der
  semgrep-Politik. Ohne Befund.
- **Rückbau/Tombstone (ADR-0025):** fünfter `ignore-refs`-Eintrag samt
  Kommentar-Zeile stil-identisch zu den vier Bestands-Tombstones; Skript per
  `git rm`, Dockerfile-Stage + `NO_CACHE_FILTER_ARCH` + `clean`-Zeile +
  Kopfkommentare entfernt; keine `arch-check`-Reste in `.github/` (CI baut
  keine gelöschte Stage). Ohne Befund.
- **Doku-Ehrlichkeit:** AGENTS §4-Zeile und Sensors-Zeile beschreiben die
  neue Mechanik korrekt (Image, Fragment, netzlos/read-only,
  exclude-Erwähnung); `.golangci.yml`-Kommentar nachgezogen; das
  versions-Modul bleibt unberührt (Digest-Pin ohne `:vX.Y.Z` matcht das
  `pin-pattern` nicht — bekannter Blind-Spot, keine Verschlechterung durch
  diesen Commit). Ohne Befund.
- **Referenz-Richtung (SDP)/Marker-Ehrlichkeit:** keine Provenance-Marker im
  Diff; ADR→Slice-Nennungen nur in Geschichte/Kommentaren (Provenance). Ohne
  Befund.

---

## Kategorie-Summary

| Kategorie | Anzahl |
| --- | --- |
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

## Verdikt

**NACHBESSERN (Config-/Beleg-Ebene, vor Closure).** Die Umstellung ist in
Verdrahtung, Fail-closed-Verhalten, Rückbau und Doku-Currency sauber; das Gate
ist auf dem Ist-Baum real wirksam (alle Injektions-Proben rot, sauberer Baum
grün, kein Silent-Grün-Pfad). Kein HIGH: der einzige echte Deckungsverlust
(MEDIUM-1) ist offen deklariert statt still. Beide MEDIUMs sind
`.a-check.yml`-lokal (kein Produkt-Code, kein a-check-CR zwingend): die
`net`-Familien-Abdeckung ist entgegen der dokumentierten Begründung per Config
herstellbar, und die Substring-Falle der `adapter:`-Pfade ist eine belegte,
heute schlafende Lücke in genau der Kapsel-Klasse, die das Gate schützen soll.
LOW-1 gehört in die Closure-Notiz (Proben-Nachtrag), INFO-1/-2 sind
Dokumentations-Kandidaten ohne Handlungszwang.
