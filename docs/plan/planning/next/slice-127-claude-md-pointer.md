# Slice slice-127: Zwei heimatlose Hard Rules nach AGENTS.md umziehen — danach ist CLAUDE.md ein Pointer

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** ohne Welle — ein einzelner Harness-Hygiene-Punkt mit eigener DoD und
ohne gemeinsame Closure-Bedingung mit anderer Arbeit (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** Baseline-Regelwerk
[`modul-09-implementierung.md` §AGENTS.md-Regeln](../../../../.harness/baseline/v5.9.0/regelwerk/modul-09-implementierung.md#agentsmd-regeln-modul-9)
(Zwei-Quadranten-Regel; AGENTS.md gehört in jeden Lauf-Kontext) und
[`grundlagen-source-precedence.md`](../../../../.harness/baseline/v5.9.0/regelwerk/grundlagen-source-precedence.md)
(Konflikt-Hard-Rule); [`AGENTS.md`](../../../../AGENTS.md) §1 und §6;
[`MR-015`](../../../../harness/conventions.md#mr-015) (AGENTS.md **routet**,
spiegelt nicht) und [`MR-012`](../../../../harness/conventions.md#mr-012).

**Berührte Spec-Stellen:** — (Harness-Dateien; keine Anforderung, kein
Spec-Stratum, kein Produkt-Code).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-23.

---

## 1. Ziel

**Neuschnitt nach der Rückführung.** Der erste Anlauf wollte `CLAUDE.md` auf
einen Pointer kürzen und stützte sich auf die Zensur *„jede Zeile steht schon
woanders"*. Der unabhängige Review hat sie widerlegt: **zwei** Aussagen stehen
im ganzen Repo nur dort. Damit ist die Aufgabe ein **Umzug**, keine Kürzung —
und die Kürzung ist ihre **Folge**, nicht ihr Inhalt.

Das Regelwerk hat für beide einen Ort, und wir haben ihn nicht benutzt:

> *„Jede Hard Rule liegt in zwei Quadranten: inferential feedforward (**steht in
> AGENTS.md**) + computational feedback (Fitness Function/Linter-Gate).“*
> — `modul-09-implementierung.md` §AGENTS.md-Regeln

Die beiden Waisen:

- **Die Konfliktregel.** `CLAUDE.md` sagt *„Konflikt melden und der
  höherrangigen Quelle folgen"*. [`AGENTS.md`](../../../../AGENTS.md) §1 trägt
  davon nur die halbe Aussage — und zwar über einen engeren Fall (*diese Datei*
  gegen eine kanonische Quelle, nicht zwei kanonische gegeneinander).
  **Der Kanon ist dabei stärker als beide:** *„dass bei Konflikt die niedriger
  rangierte Quelle **angepasst** wird, ist universal (Hard Rule)"*, und ein
  Widerspruch *„gehört benannt"*. Der Umzug ist deshalb keine Verschiebung,
  sondern eine **Angleichung**.
- **Die Benenn-Pflicht.** *„Vor der Implementierung benennen: Slice-ID,
  betroffene `DC-*`-IDs, ADR-IDs, betroffene Module, auszuführende Gates."*
  Der kanonische Schritt 3 verlangt nur, Requirement-/ADR-IDs zu
  **identifizieren**. **Das Delta ist schmaler, als es aussieht:** die
  Baseline-Slice-Vorlage trägt `Bezug:`, `Berührte Spec-Stellen:` und
  `## 3. Plan (vor Code)`, deckt also IDs und Plan bereits ab; übrig bleiben
  **Module und Gates vorab benennen**. Auch dafür gibt es einen kanonischen
  Anker an anderer Stelle — die Baseline-Vorlage der Projekt-README verlangt
  einen `Gates:`-Punkt mit derselben Warnung („halluzinierte Gates sind die
  häufigste Form von Harness-Lüge"). Kanonisch ist also die **Sorge**,
  repo-lokal ihre Verortung **pro Slice**. Heute lebt die volle Form nur in
  `.claude/commands/implement-slice.md` — außerhalb jeder gerankten Quelle.

**Der Kanon stützt die Pointer-Form ausdrücklich.** Die Baseline-Vorlage der
Projekt-README schreibt für ihren Harness-Signal-Block: *„Pointer auf die
kanonischen Quellen — **Inhalt nicht wiederholen**."* Dieselbe Regel, eine
Datei weiter.

**Dass die Baseline `CLAUDE.md` nicht kennt, ist kein Loch.** Sie verlangt, dass
AGENTS.md *in jedem Lauf-Kontext* liegt; welche Datei ein Agenten-Tool dafür
automatisch lädt, ist Werkzeug-Sache. Genau deshalb darf dort nichts **nur**
stehen: eine Datei ohne Rang in der Source Precedence ist kein Ort für eine
Hard Rule.

**Korrektur einer Prämisse des ersten Anlaufs:** die dort als „Spannung"
gemeldete Leseordnung ist **kein Konflikt**. Die Baseline weist §Leseordnung
ausdrücklich *„für den neuen Menschen"* aus, und
[`harness/README.md`](../../../../harness/README.md) nennt sie selbst die
„Menschen-Hälfte des Einstiegs"; der 8-Schritt-Pfad in
[`AGENTS.md`](../../../../AGENTS.md) §6 ist der **Agenten**-Workflow. Zwei
Adressaten, zwei Ordnungen. Es gibt nichts aufzulösen und nichts zu melden.

## 2. Vorgehen

1. **[`AGENTS.md`](../../../../AGENTS.md) §1 — Konfliktregel auf den Kanon
   ziehen:** die niedriger rangierte Quelle wird **angepasst** (nicht nur
   „gilt nicht"), und der Widerspruch **gehört benannt**. Gilt auch zwischen
   zwei kanonischen Quellen, nicht nur gegen diese Datei.
2. **[`AGENTS.md`](../../../../AGENTS.md) §6 — Schritt 3 verschärfen:** vor der
   Implementierung Slice-ID, `DC-*`-IDs, ADR-IDs, betroffene Module und
   auszuführende Gates **benennen**. Als Verschärfung des kanonischen Schritts
   markiert und mit Herkunfts-Anker `(seit slice-127)` versehen
   (`modul-09` §AGENTS.md-Regeln: Hard Rules aus dem Steering Loop tragen ihn;
   ohne Welle die Slice-Form).
3. **`MR-`Eintrag im Konventionsspeicher** für die Verschärfung. Der Kanon sagt
   „identifizieren", wir sagen „benennen" — das ist ein Delta, und der
   Konventionsspeicher ist der Ort, an dem ein Delta seinen
   **Auflösungs-Trigger** bekommt: verschärft die Baseline den Schritt selbst,
   fällt unsere Adaption weg. Ohne MR bliebe das Delta beim nächsten
   Baseline-Bump unsichtbar.
   **Kein Upstream-CR an den Kurs** — noch nicht: ein CR aus **einem** Repo
   wäre eine Aussage aus dem Anlass statt aus dem Bestand
   ([`BEO-011`](../observations.md)). Der MR trägt deshalb den Trigger *tritt
   die Klasse in einem zweiten Repo auf, wird daraus ein Konsumenten-CR* —
   die Form der Upstream-Notiz, die dieses Repo schon fährt.
4. **Erst danach `CLAUDE.md` auf den Pointer reduzieren** — und die Zeile so
   formulieren, dass sie **stimmt**: der erste Anlauf versprach Routing zur
   „Leseordnung", ein Wort, das in AGENTS.md gar nicht vorkommt.
5. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine neue Hard Rule.** Beide Regeln existieren bereits; sie wechseln den
  Ort und werden dabei auf den Kanon angeglichen.
- **Kein Gate für die beiden.** Nach der Zwei-Quadranten-Regel bleiben sie
  damit **halb durchgesetzt** — das ist als Grenze zu benennen, nicht zu
  verschweigen und nicht durch einen Heuristik-Wächter zu übertünchen.
- **Keine `dpin`-Absicherung**, kein Aktivieren des Moduls `pins` (zöge den
  `BEO-010`-Nachzug nach).
- **Keine Änderung an den `.claude/`-Hooks** und keine Aussage darüber, was sie
  abdecken — der erste Anlauf hat sie überdehnt (der Guard greift nur bei
  `Bash`, `stop-require-gates.sh` hat einen zweiten Freigabepfad).

## 4. Definition of Done

- [ ] [`AGENTS.md`](../../../../AGENTS.md) trägt beide Regeln, die Konfliktregel
      in der **kanonischen** (stärkeren) Fassung, die Verschärfung mit
      Herkunfts-Anker.
- [ ] Der `MR-`Eintrag liegt, mit Geltungsbereich, Ersetzt-Baseline-Regel und
      Auflösungs-Trigger; die Index-Zeile in
      [`harness/conventions.md`](../../../../harness/conventions.md) ist
      ergänzt.
- [ ] `CLAUDE.md` trägt Titel und genau eine Anweisung; jede Aussage der
      Vorfassung ist danach in einer **gerankten** Quelle nachweisbar —
      zeilenweise belegt, nicht behauptet (`BEO-011`).
- [ ] Die Pointer-Zeile verspricht nur, was AGENTS.md einlöst.
- [ ] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Der Umzug macht AGENTS.md länger, und AGENTS.md soll routen statt
  spiegeln** ([`MR-015`](../../../../harness/conventions.md#mr-015)). Der
  Unterschied ist, dass diese zwei Aussagen **nirgendwo sonst** stehen — sie
  sind kein Spiegel. Wird das verwechselt, wächst AGENTS.md wieder zur
  Sammelstelle. — **Ausgang:** *(bei Closure)*
- **Zwei Hard Rules ohne Gate.** Der Kanon nennt das halb durchgesetzt. Es
  wird besser als heute (aus null gerankten Quadranten wird einer), aber nicht
  gut. — **Ausgang:** *(bei Closure)*
- **Die Vollständigkeits-Aussage ist erneut die Nagelprobe.** Genau sie ist im
  ersten Anlauf gekippt; diesmal muss der Beleg **vor** dem Löschen stehen. —
  **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`next` → `in-progress`): sofort — `in-progress/` trägt keinen Slice,
und die Rückführung ist mit dem Neuschnitt aufgelöst.

**Rückführungen:** `in-progress` → `next`, falls der Review eine **dritte**
Aussage findet, die nirgends sonst steht — dann ist die Zensur wieder falsch
gewesen und der Schnitt taugt nicht.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Harness-Dateien (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-23):
  [`BEO-011`](../observations.md) ist die **zentrale** Beobachtung dieses
  Slice — sein erster Anlauf ist die siebte Instanz der Klasse. Die
  Vollständigkeits-Aussage ist deshalb zeilenweise zu belegen, **bevor**
  gelöscht wird, und der Beleg gehört in die Closure-Notiz.
  [`BEO-002`](../observations.md) für die Ränder (Verweise auf CLAUDE.md,
  Zitate ihres Inhalts in frozen Reports).

Slice-ID: slice-127. Betroffene IDs:
[`MR-015`](../../../../harness/conventions.md#mr-015),
[`MR-012`](../../../../harness/conventions.md#mr-012). Module: Harness-Dateien.
Gates: `make doc-check`, `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Harness-Dokumentation nach kanonischer
Form; die Zwei-Quadranten-Regel des Regelwerks gibt den Ort vor.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
