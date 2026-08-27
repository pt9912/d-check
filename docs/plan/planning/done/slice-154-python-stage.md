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

- [x] Der Bedarf ist **belegt** — je Fall eine Zuordnung, und der Rest ist
      benannt oder es gibt keinen. *(Neun Fälle mit Zuordnung in
      [`MR-046`](../../../../harness/conventions.md#mr-046); ein Rest benannt und
      als [slice-157](../in-progress/slice-157-config-gueltigkeit.md) ausgetragen.)*
- [x] ~~Bei Bau: Stage digest-gepinnt, `make`-Target vorhanden, in `AGENTS.md` §4
      und `harness/README.md` eingetragen, `gate-consistency` grün.~~
      **Entfällt — kein Bau.**
- [x] ~~Je Pin-Spiegel-Klasse aus [`BEO-008`](../observations.md) eine Zahl.~~
      **Entfällt — kein Bau, also keine neue `FROM`-Zeile und kein neuer Pin.**
      Der Haken stand im DoD unbedingt, in §2 Schritt 3 aber ausdrücklich am
      Bau; ausgetragen statt still abgehakt.
- [x] Bei Nicht-Bau: [`AGENTS.md`](../../../../AGENTS.md) §3.1 verweist nicht
      mehr auf einen Weg, den es nicht gibt.
- [x] `make gates` grün (Exit explizit), `make fullbuild` grün; unabhängiger
      Review.

## 5. Abnahme-Punkte / Risiken

- **Eine vierte Toolchain ist teurer als sie aussieht.** Sie bringt eine
  `FROM`-Zeile, einen Digest, eine Frische-Achse und eine Doku-Zeile mit — und
  [slice-142](../open/slice-142-freshness-weitere-achsen.md) führt bereits zwei
  ungewachte Pin-Klassen. Der Nutzen gehört gegen diesen Preis gestellt, nicht
  gegen die Bequemlichkeit. — **Ausgang: entfallen.** Die Abwägung wurde nicht
  gebraucht, weil kein Nutzen auf der anderen Seite stand: acht der neun Fälle
  sind in der erlaubten Klasse lösbar, der neunte trägt keine Toolchain.
  [slice-142](../open/slice-142-freshness-weitere-achsen.md) bleibt von diesem
  Slice unberührt.
- **Der belegte Bedarf könnte null sein.** Dann ist die ehrliche Lieferung eine
  **gestrichene Zeile in §3.1**, nicht eine Stage, die niemand ruft. Das ist ein
  Ergebnis, kein Fehlschlag. — **Ausgang: eingetreten, mit einer Korrektur am
  eigenen Satz.** Der Bedarf ist **fast** null, nicht null: ein Fall bleibt —
  die Gültigkeit einer JSON-Konfiguration. Er trägt keine Stage, aber er ist
  auch nicht nichts, und er wandert deshalb als
  [slice-157](../in-progress/slice-157-config-gueltigkeit.md) weiter statt in eine
  Notiz. Die gestrichene Zeile in §3.1 ist geliefert.

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

**Die Stage wird nicht gebaut, und §3.1 behauptet nicht mehr, es gebe sie.**
Die Regel verwies auf eine digest-gepinnte Dockerfile-Stage als vorgesehenen
Weg — es gab keine. Der Schaden war doppelt: wer den Weg gehen wollte, fand
nichts; wer die Regel las, hielt die Frage für beantwortet.
[`MR-046`](../../../../harness/conventions.md#mr-046) hält den Befund fest,
samt der Bedingungen, unter denen die Stage doch entstünde — und mit dem
Hinweis, dass ihre **Form** offen ist: die zwei jüngsten Fremd-Toolchains
dieses Repos sind keine Stage, sondern digest-gepinnte externe Images
([ADR-0010](../../adr/0010-semgrep-hermetisches-gate.md),
[ADR-0029](../../adr/0029-arch-check-via-a-check.md)).

**Der erste Anlauf hat die falsche Population gemessen — der Review hat es
gefangen.** §2 Schritt 1 verlangt die Fälle **aus der Arbeit**, je Fall eine
Zuordnung; gemessen wurde der committete **Bestand**. Die zwei sind nicht
dieselbe Menge: ein Interpreter-Aufruf darf seit
[`MR-040`](../../../../harness/conventions.md#mr-040) gar nicht mehr
eingecheckt werden, der Bestands-Befund „leer" ist für die Bedarfsfrage also
fast tautologisch — und der Anlass jenes Eintrags selbst, ein Arbeitstag mit
Host-Python, steht in keiner Datei. Nachgeholt als Tabelle mit neun Fällen.
Der Bestand stützt das Ergebnis weiterhin, aber als **zweiter** Beleg, nicht
als erster.

**Ein Rest ist geblieben und hat eine Kennung bekommen.** Die Gültigkeit einer
JSON-Konfiguration ist in der erlaubten Klasse nicht prüfbar; an ihrer Stelle
stand ein von Hand geschriebener `awk`-Klammerzähler, der Balance beantwortet
und nicht Gültigkeit. Er trägt keine Toolchain — er trägt einen Sensor, und der
ist [slice-157](../in-progress/slice-157-config-gueltigkeit.md). Das Gewicht ist
während dieses Slice gewachsen: seit
[`MR-047`](../../../../harness/conventions.md#mr-047) hängen **beide**
Durchsetzungs-Schichten an derselben Datei.

**Der Slice hat einen Nachbar-Befund mitgenommen, der schwerer wog als sein
eigener.** Die Permission-Sperrliste war als zweite Durchsetzungs-Schicht
eingezogen worden — ohne Eintrag, gegen einen wörtlichen Kanon-Satz. Sie hat
jetzt [`MR-047`](../../../../harness/conventions.md#mr-047), und dort steht
auch, was die Commit-Botschaft nur halbseitig sagte: die genannten
Durchfall-Klassen betreffen alle die Interpreter-Hälfte, die der Wächter
zweitdeckt; die **git-/docker-Hälfte ist einschichtig**, hat eigene Lücken und
ein gutes Dutzend Formen ganz ohne Eintrag. Dazu zwei weitere Nicht-Zusagen:
beide Schichten sind über `make` konstruktionsbedingt durchlässig, und für die
neue gibt es **keine** wiederholbare Probe.

**Drei eigene Aussagen waren am Bestand widerlegt.** §3.1 sagte „Messungen
macht das Produkt" — sechs Messungen dieses Repos macht nicht das Produkt; und
„was ein Gate-Skript braucht, tragen `bash` und die POSIX-Werkzeuge" — zwei
Netz-Targets brauchen `curl` und `unzip`, die in der Host-Klasse von Absatz 1
nicht stehen. Beides ersetzt durch eine Rangfolge und eine benannte
Zusatz-Erwartung. Und „zwölf ausführbare Skripte" passte zu keiner messbaren
Menge: **zehn** sind ausführbar gesetzt, **zwölf** heißen `*.sh`, **vierzehn**
tragen einen bash-Shebang. Die zwei außerhalb von `*.sh` — die `.githooks/` —
fielen aus Zählung **und** Such-Enumeration heraus, dieselbe Achse zweimal
ausgelassen ([`BEO-009`](../observations.md) Richtung b).

**Sensors:** `make gates` (Exit 0, zehn Glieder, 531 Dateien, 0 Befunde),
`make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen, Closure-Profil 479
Dateien / 0 Befunde), `make guard-probe` (40 Proben, 0 Fehlschläge, Exit 0).
Ein unabhängiger Review ist gelaufen; seine zwei HIGH, vier MEDIUM und drei LOW
sind in `f6e72cb` eingearbeitet.
