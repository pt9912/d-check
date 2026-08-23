# Slice slice-127: CLAUDE.md auf einen Pointer eindampfen — eine Datei, die nur Kopien trug

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** ohne Welle — ein einzelner Harness-Hygiene-Punkt mit eigener
DoD und ohne gemeinsame Closure-Bedingung mit anderer Arbeit (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §1 (*„trägt Hard Rules und
Pointer … und dupliziert deren Inhalt nicht — sonst entsteht Drift"*) und §3.1;
[`MR-015`](../../../../harness/conventions.md#mr-015) (Auflösung der
[`MR-012`](../../../../harness/conventions.md#mr-012)-Pointer-Drift:
AGENTS.md **routet**, spiegelt nicht mehr) —
CLAUDE.md hat diese Behandlung nie bekommen;
[`MR-004`](../../../../harness/conventions.md#mr-004) und
[`MR-005`](../../../../harness/conventions.md#mr-005) (die `.claude/`-Hooks, die
die beiden Regeln **mechanisch** tragen).

**Berührte Spec-Stellen:** — (Harness-Einstiegsdatei; keine Anforderung, kein
Spec-Stratum, kein Produkt-Code).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

`CLAUDE.md` ist vom 2026-06-23 und trägt **nichts Eigenes**. Jede ihrer Zeilen
steht bereits in [`AGENTS.md`](../../../../AGENTS.md) (§1 Konflikt-Regel, §2
Source Precedence, §3.1 make-only, §6.3 ID-Nennung, §6.6 Gate-Lauf, §6.8 keine
Erfolgsmeldung ohne Gate) oder in
[`harness/README.md`](../../../../harness/README.md) §Leseordnung bzw.
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Sie ist zu hundert Prozent Kopie — und prompt gedriftet:

- **Alterung:** die Toolchain-Verbotsliste lässt die Host-Sprachtoolchain aus,
  die [`AGENTS.md`](../../../../AGENTS.md) §3.1 an erster Stelle nennt; die
  Leseliste kennt weder die Roadmap noch das Beobachtungs-Register
  [`observations.md`](../observations.md) (Pflicht-Vorprüfung **jedes**
  Slice-Plans) noch das vendorte Regelwerk.
- **Verfälschung beim Kopieren, schon am ersten Tag:**
  [`AGENTS.md`](../../../../AGENTS.md) §6.6 sagt Gate-Lauf vor **Handoff**,
  CLAUDE.md macht daraus „vor dem **Abschluss**" — und ebnet damit die scharfe
  Trennung Handoff (`make gates`) / Closure (`make fullbuild`) ein.

**Der Einwand gegen einen reinen Pointer ist geprüft und trägt nicht.**
CLAUDE.md ist die Datei, die der Agent automatisch geladen bekommt, AGENTS.md
nicht — die Sorge ist also berechtigt, dass eine Regel ohne zweite Fundstelle
niemanden mehr bindet. Gemessen am Bestand hält sie nicht: die **zwei** Regeln,
die auch ohne AGENTS.md gelten müssten, hängen an **Hooks**, nicht an Prosa —
`.claude/hooks/pretooluse-command-guard.sh` weist Host-Toolchain-Aufrufe ab
(in dieser Sitzung zweimal ausgelöst), und
`.claude/hooks/stop-require-gates.sh` gibt den Stop nur frei, wenn der
Repo-**Inhalt** durch einen erfolgreichen `make gates`-Lauf gedeckt ist. Die
Kopien erzwingen nichts; sie beschreiben nur, was ohnehin erzwungen wird.

## 2. Vorgehen

1. `CLAUDE.md` auf Titel plus **eine** Anweisung reduzieren: vor jeder Änderung
   an Code oder Dokumentation zuerst [`AGENTS.md`](../../../../AGENTS.md) lesen
   und befolgen, sie trägt die Hard Rules und routet weiter.
2. Prüfen, ob eine Fundstelle auf CLAUDE.md zeigt, die durch die Kürzung ins
   Leere liefe (`grep`, beide Richtungen).
3. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Nachziehen der drei Defekte.** Sie sind das Beweismaterial, nicht die
  Aufgabe — wer sie einzeln repariert, hat in zwei Monaten denselben Befund.
- **Keine `dpin`-Absicherung** der Duplikation (Modul `pins` im eigenen Profil
  aktivieren). Sie wäre gate-gedeckt, zöge aber den `BEO-010`-Nachzug nach
  (`FOCUS_DISABLE` im Makefile, die Netzlos-Modulliste im Go-Test, drei
  Prosa-Beschreibungen) — Mechanik für eine Datei, die gar keinen Eigeninhalt
  hat.
- **Keine Auflösung der Leseordnungs-Spannung.** [`AGENTS.md`](../../../../AGENTS.md)
  §6 sagt „1. `harness/README.md` lesen", [`harness/README.md`](../../../../harness/README.md)
  §Leseordnung stellt AGENTS.md auf Platz 1. Zwei kanonische Quellen
  widersprechen sich; das ist **zu melden**, nicht nebenbei zu entscheiden.
- **Kein neuer Eigeninhalt.** Der eine Kandidat — dass der PreToolUse-Guard auch
  Bash-Heredocs mit dem Toolchain-Token abweist — ist ein *Zuwachs*, kein
  Currency-Fix, und wäre ein eigener Entscheid.

## 4. Definition of Done

- [ ] `CLAUDE.md` trägt Titel und genau eine Anweisung; keine Regel-Kopie mehr.
- [ ] Kein Verweis läuft ins Leere (`make doc-check` grün, plus `grep` auf
      Fundstellen, die auf gekürzte Abschnitte zeigen).
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.
- [ ] Die Leseordnungs-Spannung ist als Konflikt **gemeldet** (Closure-Notiz),
      nicht stillschweigend aufgelöst.

## 5. Abnahme-Punkte / Risiken

- **Die Kürzung nimmt dem automatisch geladenen Kontext Inhalt weg.** Wenn ein
  Agent die eine Anweisung ignoriert, steht ihm nichts mehr im Weg — außer den
  beiden Hooks. Ob das reicht, zeigt sich erst im Betrieb, nicht im Gate. —
  **Ausgang:** *(bei Closure)*
- **`make gates` deckt Prosa nicht.** Dass die verbleibende Zeile inhaltlich
  richtig ist, prüft kein Modul. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): sofort — welle-82 ist geschlossen,
`in-progress/` trägt keinen Slice.

**Rückführungen:** `in-progress` → `next`, falls der Review zeigt, dass
CLAUDE.md doch Eigeninhalt trägt, den AGENTS.md nicht abdeckt — dann ist es
kein Kürzungs-, sondern ein Umzugs-Slice.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Einstiegsdateien (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) ist unmittelbar einschlägig — dieser Slice
  begründet eine **Entfernung** mit einer Zensur über den Bestand („jede Zeile
  steht schon woanders"), und genau solche Vollständigkeits-Aussagen sind die
  Klasse; der Beleg ist zeilenweise zu führen, nicht zu behaupten.
  [`BEO-002`](../observations.md) für die Ränder (Verweise auf CLAUDE.md).

Slice-ID: slice-127. Betroffene IDs:
[`MR-015`](../../../../harness/conventions.md#mr-015),
[`MR-012`](../../../../harness/conventions.md#mr-012). Module:
Harness-Einstiegsdateien. Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Harness-Dokumentation nach etablierter
Konventions-Form; die Regel, der gefolgt wird, ist bereits als
[`MR-015`](../../../../harness/conventions.md#mr-015) gefällt.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
