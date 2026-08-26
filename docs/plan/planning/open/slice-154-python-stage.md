# Slice slice-154: Wo ein Skript nötig ist, läuft es als eigene Dockerfile-Stage

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (die Regel, die diesen Weg
benennt); [`MR-040`](../../../../harness/conventions.md#mr-040) (die Sperre auf
dem Host); [ADR-0002](../../adr/0002-distribution-ghcr-image.md) (Multi-Stage);
[ADR-0011](../../adr/0011-digest-pins-build-gate-images.md) (jede `FROM`-Zeile
digest-gepinnt); [`BEO-008`](../observations.md) (Pin-Spiegel).

**Berührte Spec-Stellen:** — (Build- und Werkzeug-Schicht; das Produkt bleibt
unberührt).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-26.

---

## 1. Ziel

Seit [`MR-040`](../../../../harness/conventions.md#mr-040) blockiert der
Tool-Call-Wächter Host-Skript-Interpreter.
[`AGENTS.md`](../../../../AGENTS.md) §3.1 nennt dafür den vorgesehenen Weg:
*„Wo ein Skript wirklich nötig ist, läuft es als **eigene Dockerfile-Stage** wie
jede andere Toolchain dieses Repos, digest-gepinnt und über ein `make`-Target."*
Diese Stage gibt es noch nicht — die Regel verweist auf etwas, das fehlt.

**Die erste Frage ist, ob sie gebraucht wird.** Datei-Änderungen macht das
Harness-Werkzeug ohne Shell; Messungen gehören ins Produkt — das ist die
Identität dieses Repos, es hat `tools/*.sh` durch Go-Module ersetzt. Bleibt ein
Rest, der weder das eine noch das andere ist? Der gehört **benannt**, bevor eine
vierte Toolchain entsteht.

## 2. Vorgehen

1. **Den Bedarf belegen, nicht behaupten.** Aus der Sitzungshistorie die Fälle
   sammeln, in denen ein Skript lief, und je Fall entscheiden: Datei-Werkzeug ·
   Produkt · POSIX-Werkzeug · echter Rest. Nur der Rest begründet die Stage.
2. Trägt der Rest: Stage im **vorhandenen** Dockerfile, Basis-Image per
   `@sha256:`-Digest gepinnt ([ADR-0011](../../adr/0011-digest-pins-build-gate-images.md)),
   plus ein `make`-Target, das den Repo-Baum mountet.
3. **Die Pin-Spiegel mitziehen** — [`BEO-008`](../observations.md) führt drei
   Klassen, und eine neue `FROM`-Zeile berührt `make versions`, die
   Doku-Enumerationen und die Prosa-Pins. Je Klasse eine Zählung.
4. `AGENTS.md` §4 und `harness/README.md` tragen das Target; `gate-consistency`
   grün.
5. Trägt der Rest **nicht**: §3.1 sagt danach nicht mehr, es gebe diesen Weg —
   die Regel wird auf das zurückgeschnitten, was existiert.
6. `make gates`, `make fullbuild`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Rücknahme der Sperre.** Der Host bleibt frei von Interpretern,
  unabhängig davon, wie dieser Slice ausgeht.
- **Kein Netz im Gate.** Die Stage wird beim Build gezogen wie die übrigen; die
  Gates bleiben netzlos.
- **Keine Skript-Sammlung.** Was entsteht, ist eine Ausführungs-Form, kein neues
  `tools/`-Verzeichnis — sonst kommt zurück, was
  [ADR-0028](../../adr/0028-planning-lifecycle-modul.md) abgeschafft hat.

## 4. Definition of Done

- [ ] Der Bedarf ist **belegt** — je Fall eine Zuordnung, und der Rest ist
      benannt oder es gibt keinen.
- [ ] Bei Bau: Stage digest-gepinnt, `make`-Target vorhanden, in `AGENTS.md` §4
      und `harness/README.md` eingetragen, `gate-consistency` grün.
- [ ] Je Pin-Spiegel-Klasse aus [`BEO-008`](../observations.md) eine Zahl.
- [ ] Bei Nicht-Bau: [`AGENTS.md`](../../../../AGENTS.md) §3.1 verweist nicht
      mehr auf einen Weg, den es nicht gibt.
- [ ] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Eine vierte Toolchain ist teurer als sie aussieht.** Sie bringt eine
  `FROM`-Zeile, einen Digest, eine Frische-Achse und eine Doku-Zeile mit — und
  [slice-142](../open/slice-142-freshness-weitere-achsen.md) führt bereits zwei
  ungewachte Pin-Klassen. Der Nutzen gehört gegen diesen Preis gestellt, nicht
  gegen die Bequemlichkeit. — **Ausgang:** *(bei Closure)*
- **Der belegte Bedarf könnte null sein.** Dann ist die ehrliche Lieferung eine
  **gestrichene Zeile in §3.1**, nicht eine Stage, die niemand ruft. Das ist ein
  Ergebnis, kein Fehlschlag. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls der Bedarfs-Beleg ein
Produkt-Delta nahelegt (eine Messung, die ins Produkt gehört) statt einer Stage.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Build-Schicht (GF), Harness-Werkzeug (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-26):
  [`BEO-008`](../observations.md) für die Pin-Spiegel;
  [`BEO-011`](../observations.md) für jede Aussage darüber, wofür ein Skript
  „nötig" sei — der Bedarf gehört gezählt, nicht erinnert.

Slice-ID: slice-154. Betroffene IDs: — (kein `DC-`-Bezug; Build-Schicht).
Module: — . Gates: `make gates`, `make fullbuild`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — neue Stage in einem etablierten
Multi-Stage-Build.

## 9. Closure-Notiz (nach `done/`)

*(wird mit dem Closure-Body gefüllt)*
