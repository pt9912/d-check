# Review slice-243 — Mechanisierungs-Entscheid Trigger-Audit (R1)

- **Review-Art:** Plan/Design-Review — geprüft gegen den Slice-Plan `slice-243` (§1
  Ziel und Abgrenzung, §6 Risiken), die Hard Rules (`AGENTS.md` §3), die Konventionen
  (`MR-073`, `MR-053`) und das Delta `v6.13.0` · `regelwerk/modul-06-roadmap.md`
  §Closure (vierte-Mal-Schwelle); **nicht** gegen die DoD (Verifikation,
  getrennter Kontext).
- **Gegenstand:** Range `920c1c45..HEAD`, soweit `slice-243`: `43bb8daa` (Plan
  angelegt, nur Plan-Datei), `98fc696c` (Beanspruchung — Ruhe-Marker verlässt
  Roadmap, reiner Rename der Plan-Datei), `0a6ba54b` (DoD abgehakt + Evaluierung
  in §7, nur Plan-Datei). Im Range daneben, außerhalb des Gegenstands: die
  slice-241/242-Moves samt Verify-Report (eigenständig reviewt) und zwei
  Deps-Bumps mit eigener Kennung (`4401eda5`, `856fc94f`). Der Plan steht in
  `in-progress/` — die Korrektur aus F-1 ist in dieser Lifecycle-Position noch
  ohne Move möglich.
- **Skill:** `reviewer.md` @ 1.16.0
- **Modell-ID:** glm-5.3-flash (Claude Agent SDK)
- **Datum:** 2026-09-29
- **Eingangs-Kontext:** Slice-Plan `slice-243`; Nachbarplan `slice-241` (Adoption,
  V-Beleg); `AGENTS.md` §3/§5; `harness/conventions.md` (`MR-073`, `MR-053`,
  `MR-056`); `v6.13.0` · `regelwerk/modul-06-roadmap.md` §Closure; vorherige
  Findings am gleichen Vorgang: slice-242 R1 F-1 (`commit-boundary-cross-slice`),
  slice-241 R1 (V-1, Trigger-Zeilen-Bestand).

---

## Findings

### F-1 — MEDIUM

- **Kategorie:** MEDIUM
- **Quelle:** Dokumentations-Regel 14
  ([`harness/rules/dokumentations-regeln/14-messung-form-vor-zaehlung.md`](../../harness/rules/dokumentations-regeln/14-messung-form-vor-zaehlung.md))
  · Anker 17 `BEO-ALL/zaehlmethode-misst-proxy-statt-gegenstand`
- **Pfad:** `docs/plan/planning/in-progress/slice-243-trigger-audit-mechanisierungs-entscheid.md:110`
  („die 97 ADR-Trigger-Konditionale") und `:123` („97 Re-Evaluierungs-Trigger,
  jeder ein einzigartiges Konditional")
- **Befund:** Die ADR-Zeile der Evaluierung (DoD 1) belegt die Klasse mit einer
  Zahl, deren Zähl-Form nicht ausgeschrieben ist und die mit keiner natürlichen
  Form reproduzierbar ist. Gemessen gegen `docs/plan/adr/` (Stand `0a6ba54b`):
  76 Sektionen `## Re-Evaluierungs-Trigger`, 101 Zeilen mit dem Begriff
  „Re-Evaluierungs-Trigger", 55 Konditional-Sätze (Wenn/Falls/Sobald/Sofort) in
  den Sektionen, 188 Nichtleer-Sektions-Zeilen — keine dieser Formen ergibt 97,
  und der Plan nennt keine weitere. Die qualitative Aussage („jeder ein
  einzigartiges Konditional") trägt — Stichprobe je Anker-17-Arbeitsanweisung
  über drei Treffer: ADR-0048 („Wenn eine Floskel-Art dreimal …"), ADR-0066
  („Zeigt sich über mehrere Monate …"), ADR-0072 („Der erste Adopter, dessen
  Workflow-Ablage nicht …") sind je einzigartige Konditionale. Der
  Prosa-Verbleib für die ADR-Klasse steht und fällt nicht an der Zahl — aber
  DoD 1 verspricht je Zeile einen Gegenstand „(Datei/Config-Mechanismus, nicht
  bloße Behauptung)", und für diese Klasse ist der Beleg eine Behauptung mit
  nicht nachvollziehbarer Zählmethode. Folge-Szenario: ein späterer Leser (oder
  der Closure-Note-Review) zählt 76 bzw. 55 und kann nicht unterscheiden, ob
  der Bestand gewachsen ist oder die Messung falsch war — der Vollzugs-Beleg
  „gegen den Baum bewertet" (§7, „Was hat funktioniert") ist für die Klasse
  damit untreu.
- **Verifizierbar:** ja — `grep -rh "Re-Evaluierungs-Trigger" docs/plan/adr/ | wc -l`
  (101); `grep -rh "^## Re-Evaluierungs-Trigger" docs/plan/adr/ | wc -l` (76);
  Konditional-Zählung per awk über die Sektionen (55). Keine Form liefert 97.
- **Klasse:** messung-ohne-form (Zahl ohne ausgeschriebene Zähl-Form; Proxy-Risiko
  der `zaehlmethode`-Familie)

---

## Negativbefunde (geprüft, ohne Befund)

1. **Carveout-Zeile (Prüfpunkt 1):** `docs/plan/carveouts/` trägt nur `done/` und
   `README.md` — kein offener Carveout; CO-001 (Auflösungs-Trigger erfüllt am
   2026-09-17) und CO-002 stehen in `done/`, beide mit eigenem
   Auflösungs-Trigger-Konditional im Körper. Die Zeile „Prosa-Verbleib — ein Gate
   auf ‚keine offenen Carveouts' kriminalisiert legitime Zustände" ist gegen den
   Baum belegt. (Die Pfad-Kurzfassung `carveouts/done/` statt
   `docs/plan/carveouts/done/` ist erkennbar und kein Befund.)
2. **bootstrap-aware Gate — n.a. (Prüfpunkt 1):** `Makefile:39`
   `THRESHOLD ?= 93` mit Kalibrierungs-Bindungs-Kommentar (`Makefile:35`); die
   Schwelle ist eine kalibrierte Konstante ohne Hochschalt-Trigger — die
   n.a.-Begründung ist baum-wahr, der V-Beleg auf slice-241 deckt sich mit dem
   Bestand.
3. **Hard-Rule-Zeile (Prüfpunkt 1):** `AGENTS.md` §3 — Trigger-Zeilen stehen auf
   §3.2, §3.3, §3.4, §3.6, §3.8 (Zeilen 117, 140, 181, 215, 254); §3.1, §3.5,
   §3.7 und §3.9 tragen keine — exakt der Alt-Bestand, den die Zeile nennt. Die
   Grenze der class-Abgrenzung (structure-Form würde genau diese vier
   ankreiden) stimmt.
4. **Vierte-Mal-Schwelle (Prüfpunkt 2):** Das Plan-Zitat ist wörtlich korrekt
   (Ellipsis für „der die Klasse künftig fängt") — `v6.13.0` ·
   `regelwerk/modul-06-roadmap.md` §Closure. Die Begründung „Prosa-Verbleib für
   alle vier" unterschlägt keine Druck-Lage: der Rückkehr-Anlass ist
   ausdrücklich benannt („bei der ersten Audit-Auffälligkeit kehrt die Frage
   zurück — dann mit belegtem Anlass") und damit **strenger** als die
   Baseline-Schwelle (viertes Mal); wiederholte Audit-Fehler und
   Steuersignale sind der benannte Rückkehr-Grund, nicht verschwiegen. Die
   je-Klassen-Begründung (Zustands-Prüfung vs. Urteils-Bedeutung) ist
   substanziell und nicht bloß eine Schwellen-Citation — der planinterne
   Riskik-Punkt R1 (vorzeitige Mechanisierung) ist damit zu Recht abgewehrt:
   die Baseline verlangt den Sensor erst bei ausgeschöpfter Prosa
   („verkörpert heißt nicht zwangsläufig automatisiert", ebenda).
5. **Abgrenzung (Prüfpunkt 3):** gehalten — kein Produkt-Sensor (kein Code,
   keine Config, kein `DC-*`-Delta im Range), keine Automatisierung des Urteils
   (§7 trennt durchgängig Zustands-Prüfung von Urteils-Bedeutung), keine
   Retro-Prüfung des Bestands (§7 nennt den Bestand nur als Zustands-Menge,
   slice-241s Abgrenzung unangetastet).
6. **Kein Sensor versteckt eingebaut (Prüfpunkt 4):** `git diff 920c1c45..HEAD`
   über `.d-check.yml`, `.d-check.closure.yml`, `Makefile`, `tools/`,
   `internal/`, `spec/` — leer. Der Commit `0a6ba54b` berührt ausschließlich die
   Plan-Datei; die Commit-Zerlegung der drei slice-243-Commits ist sauber
   (43bb8daa Plan · 98fc696c Ruhe-Marker+Rename · 0a6ba54b DoD/§7) — kein
   Wiederholungsfall der slice-242-F-1-Klasse.
7. **Sechzehn-Fragen-Sweep:** reiner Planning-Diff — keine Gate-Skripte, Module,
   Imports, Netz-Zugriffe oder Code-Kommentare berührt (Fragen 1–6, 10–18 ohne
   Gegenstand). Zustandsfelder in §7 (Risiken-Ausgänge, Lerneintrag/Folge-Slice/
   Register) tragen Zustand und Beleg, keine Chronik (Frage 7 ohne Befund); die
   drei Paarungen der Closure-Notiz sind semantisch gefüllt (Lerneintrag
   „vier Klassen, keine mechanisierbar ohne Anlass" — Folge-Slice: keiner mit
   Rückkehr-Bedingung — Register: kein neuer Eintrag).

---

## Kategorie-Summary

| Kategorie | Anzahl | Findings |
| --- | --- | --- |
| HIGH | 0 | — |
| MEDIUM | 1 | F-1 (messung-ohne-form, ADR-Klassen-Zeile) |
| LOW | 0 | — |
| INFO | 0 | — |

## Verdikt

**MEDIUM blockiert typischerweise — hier mit engem Lösungsfenster:** F-1 liegt
in der Closure-Notiz eines Slices, der noch in `in-progress/` steht; die
Korrektur (Zähl-Form ausschreiben oder die Zahl durch die baum-treue Form
ersetzen — 76 Sektionen bzw. die qualitative Aussage ohne Zahl) ist eine Zeile
in derselben Lifecycle-Position, bevor der Slice nach `done/` wandert. Die
Entscheidung selbst (Prosa-Verbleib je Klasse, keine Mechanisierung) ist
gegen den Baum belegt, delta-treu zur vierte-Mal-Schwelle und abgrenzungs-treu;
sie wird von F-1 nicht berührt. Negativbefunde 1–4 bestätigen die vier
Klassen-Zeilen je Gegenstand — mit der einen F-1-Ausnahme für die ADR-Zeile.
