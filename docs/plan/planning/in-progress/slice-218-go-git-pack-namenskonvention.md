# Slice slice-218: Die go-git-Pack-Namenskonvention als Grenze der git-lesenden Module

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
und [`DC-FA-COMMITS-001`](../../../../spec/lastenheft.md#dc-fa-commits-001--traceability-kennung-in-commit-messages-über-eine-commit-range-modul-commits-opt-in)
— **als die zwei Zusagen, deren Reichweite hier benannt wird**, nicht als
geänderte Anforderungen. Keine ADR: Der Slice trifft keine Entscheidung, er
schreibt eine gemessene Eigenschaft auf.

**Berührte Spec-Stellen:** — (der Slice berührt keine Spec-Stelle).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912.

## 1. Ziel und Abgrenzung

**Ziel.** Eine gemessene Produkt-Grenze bekommt einen Ort: **Die git-lesenden
Module finden Objekte nur in Pack-Dateien, die git's kanonischen Namen tragen**
(`pack-<Hash-des-Packs>.{idx,pack}`). Weicht der Name ab, ist der Pack
unsichtbar, und Objekte, die nur dort liegen, sind *nicht gefunden*.

**Der Anlass war ein realer Ausfall, kein Gedankenspiel.** Während slice-217
brach der `pre-commit`-Hook mit `HEAD-Tree nicht lesbar: object not found` ab,
nachdem `git maintenance` Packs mit `loose-`-Präfix angelegt hatte —
`git repack -A -d` behob es. **Das trifft jeden Adopter**, der `vcs` oder
`commits` fährt und dessen git eine Wartungs-Aufgabe laufen lässt; es ist
keine Eigenschaft dieses Repos.

**Sechs Proben, isoliert gefahren** (Probe-Repo, alle Objekte im Pack, `.idx`,
`.pack` und `.rev` gemeinsam umbenannt, sonst nichts verändert):

| Pack-Name | Ergebnis |
|---|---|
| `pack-<eigener Hash>` (Kontrolle) | `2 Datei(en) geprüft, 0 Befund(e)`, Exit **0** |
| `loose-<eigener Hash>` | `staged-Basis "HEAD" nicht auflösbar`, Exit **2** |
| `pack-0123…4567` (gültige Form, **fremder** Hash) | Exit **2** |
| `pack-zzzzzzzz` | Exit **2** |
| `packXYZ` | Exit **2** |
| `xpack-abc` | Exit **2** |

**Damit zählen beide Namensteile**, und das ist mehr, als der Anlass zeigte:
Nicht nur das Präfix (`loose-` scheitert mit korrektem Hash), sondern auch der
Hash (die gültige Form scheitert mit fremdem Hash). **`git` selbst liest
weiter** — es enumeriert `*.idx` unabhängig vom Namen; nur go-git tut es
nicht.

**Abgrenzung — vier Punkte, jeder mit Grund:**

1. **Kein Fix im Produkt, und der Grund ist ein Kernvertrag.** Die
   naheliegende Abhilfe wäre, dass d-check selbst repackt — das schriebe in
   das geprüfte Repository und bräche
   [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).
   *Es wäre ein anderer Vorgang*, und zwar einer mit ADR-Last.
2. **Auch keine bessere Fehlermeldung.** Sie wäre Produkt-Code, verdoppelte
   den Slice und verlangte eine eigene Entscheidung darüber, wie viel ein
   Werkzeug über die Objektdatenbank seines Wirts vermuten darf —
   *Schicht-Abgrenzung*: dieser Slice ändert Dokumentation, keinen Code.
3. **Kein Sensor darauf.** Ein Wächter über Pack-Namen wäre ein neues Modul
   mit eigener Scan-Zusage. *Bestand bleibt bewusst stehen.*
4. **Keine Lastenheft-/Spezifikations-Änderung.** Eine aufgeschriebene Grenze
   ist keine geänderte Anforderung; die Zusage war nie eine andere. Ein
   Spec-Eintrag käme erst in Frage, wenn das Produkt sein Verhalten ändert
   (Punkt 1 oder 2) — *ein Folge-Slice übernähme es dann.*

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(1)** [`harness/sensors/adr-check.md`](../../../../harness/sensors/adr-check.md)
      und [`harness/sensors/trace-check.md`](../../../../harness/sensors/trace-check.md)
      tragen die Grenze in ihrem Abschnitt *Grenze — was das Grün nicht
      abdeckt*: **Symptom, Exit-Code, Auslöser und Abhilfe**, jeweils mit der
      gemessenen Ausgabe statt einer Beschreibung.
- [ ] **(2)** Das [Benutzerhandbuch](../../../../docs/user/benutzerhandbuch.md)
      nennt sie bei **beiden** git-basierten Modulen (§`vcs`, §`commits`) —
      denn sie ist eine Eigenschaft des **Produkts**, nicht dieses Repos, und
      ein Adopter liest das Handbuch, nicht unsere Sensor-Dateien.
- [ ] **(3)** **Das Negativ-Ergebnis steht dabei**, und zwar als gemessenes:
      `tracked` ist **nicht** betroffen (es liest den git-**Index**, nicht die
      Objektdatenbank). Belegt durch **Positiv-Kontrolle** — derselbe
      `target-untracked`-Befund unter beiden Pack-Namen —, nicht durch einen
      grünen Lauf, der auch „nichts geprüft" heißen könnte.
- [ ] `make gates` grün.
- [ ] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

**Die Messung ist vor dem Plan gelaufen und steht in §1** — sie ist die
Grundlage, nicht ein Schritt. Was folgt, ist das Aufschreiben.

**Die Form der Aussage, bevor sie an drei Orte geht.** Eine Grenze trägt hier
vier Teile, und jeder muss aus einer Probe stammen:

1. **Symptom** — die wörtliche Meldung, und dass sie **variiert**:
   `staged-Basis "HEAD" nicht auflösbar: reference not found` (Probe-Repo, alle
   Objekte im Pack) gegen `HEAD-Tree nicht lesbar: object not found` (dieses
   Repo, wo der HEAD-Commit noch lose lag). **Ein Leser, der nur die eine
   kennt, erkennt die andere nicht wieder.**
2. **Exit-Code** — **2**, also fail-closed. Das ist die entlastende Hälfte des
   Befunds und gehört genauso hin wie die belastende: Es gibt **kein stilles
   Grün**.
3. **Auslöser** — ein Pack-Name, der von `pack-<eigener Hash>` abweicht;
   erzeugt etwa von `git maintenance run --task=loose-objects`.
4. **Abhilfe** — `git repack -A -d` (die `-A`-Form, damit unerreichbare
   Objekte lose werden statt verworfen).

**Schritte:**

1. Die zwei Sensor-Dateien ergänzen — als **nummerierter Punkt** in ihrem
   vorhandenen Grenzen-Abschnitt, in dessen Form.
2. Das Handbuch bei §`vcs` und §`commits` ergänzen. Beide Stellen brauchen die
   Aussage, weil ein Adopter nur eines der Module fahren kann.
3. `make doc-check`, dann `make gates`, dann Handoff an den unabhängigen
   Review.

**Was der Plan nicht mitnimmt**, steht in §1 — vier Punkte, die den Lauf
binden.

## 4. Trigger

**Beanspruchung:** WIP-Limit frei (`in-progress/` leer nach der Closure von
slice-217), und der Auftraggeber hat den Zuschnitt am 2026-09-08 beauftragt,
nachdem der Befund in slice-217 ausdrücklich als eigener Vorgang benannt
worden war.

**Rückführung nach `next/`** (`in-progress→next`): wenn sich beim Schreiben
zeigt, dass die Grenze eine **Verhaltens**-Antwort verlangt statt einer
Beschreibung — etwa weil die Meldung so irreführend ist, dass ein Adopter sie
für einen Defekt seines Repos hält. Dann ist der Gegenstand Produkt-Code und
der Slice falsch geschnitten (§1 Punkt 2 schließt Code aus).

**Rückführung nach `open/`** (`in-progress→open`): wenn die Messung sich beim
Aufschreiben als nicht reproduzierbar erweist — etwa weil ein weiterer
Faktor mitspielte, den die sechs Proben nicht isoliert haben.

## 5. Closure-Trigger

DoD (1) bis (3) abgehakt, `make gates` grün mit echter Ausgabe, unabhängiger
Review durchgeführt und seine Befunde eingearbeitet, Closure-Notiz geschrieben,
Register fortgeschrieben, jedes Risiko aus §6 mit einem der drei Ausgänge.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Mechanismus ist erschlossen, nicht gelesen.** Gemessen ist das
  **Verhalten** über sechs Proben; **warum** go-git so entscheidet, steht in
  seiner Quelle, und die ist ohne Host-Go nicht lesbar (§3.1). Der Slice darf
  deshalb nur sagen, *was* passiert, und muss die Erklärung als das kennzeichnen,
  was sie ist. **Die erste Erklärung war bereits falsch:** *„go-git listet
  Packs nach dem `pack-`-Präfix"* — widerlegt durch die Probe mit gültiger Form
  und fremdem Hash, die ebenfalls scheitert. Wer den zweiten Erklärungsversuch
  für sicherer hält als den ersten, hat nichts gelernt. — **Ausgang:** \<offen\>
- **Drei Orte, eine Aussage — das ist eine Drift-Quelle.** Sensor-Datei,
  zweite Sensor-Datei und Handbuch tragen dieselbe Grenze; ändert sich das
  Verhalten (go-git-Bump), müssen alle drei mit. **Kein Sensor koppelt sie**,
  und `mentions` wäre das falsche Werkzeug: Es misst *erwähnt*, nicht
  *übereinstimmend*. — **Ausgang:** \<offen\>
- **Die Grenze könnte weiter reichen als die drei gemessenen Module.** Geprüft
  sind `vcs`, `commits` und `tracked`. Ob eine andere Stelle des Produkts die
  Objektdatenbank liest, ist **nicht** geprüft — und §3.8 verlangt genau diese
  Frage. Der Slice darf deshalb nicht *„die git-basierten Module"* schreiben,
  wenn er drei gemessen hat. — **Ausgang:** \<offen\>

## 7. Closure-Notiz

\<wird vor dem `git mv` nach `done/` gefüllt\>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende. Dieses Repo führt **drei** Prüfungen — die
zwei kanonischen und, als Adaption, den Nachtlauf-Stand
([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:268-269 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `*` (Repo-Default). Der Slice ändert drei Dokumente. Der
**Gegenstand** ist zwar Produkt-Verhalten, aber es wird **gemessen, nicht
geändert** — und für Produkt-Code führt dieses Repo keine eigene Sub-Area; die
Modus-Deklaration kennt `*` und `tools/harness/`. Eine dritte hier zu erfinden
wäre eine Sub-Area ohne Deklaration.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md:274-274 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, **39** Verzeichnisse über beide
Kürzel; `BEO-HARN` einzeln geöffnet — der eine offene Eintrag dort betrifft
`--check-latest` und berührt diesen Slice nicht). **Vier** Einträge sind
einschlägig:

- [`module-promise-only-on-scan-axis`](../observations/BEO-ALL/module-promise-only-on-scan-axis/observation.md)
  (2×) — **der tragende, und er ist die Quelle von §6 Risiko 3.** Die Frage
  *„welche Eingaben liest ein Modul, die es nicht scannt?"* ist hier wörtlich
  der Gegenstand: `vcs` und `commits` lesen die **Objektdatenbank**, die in
  keiner Scan-Wurzel steht. Die Grenze **ist** die Antwort auf diese Frage.
- [`wortlaut-behauptet-pruefung-die-fehlt`](../observations/BEO-ALL/wortlaut-behauptet-pruefung-die-fehlt/observation.md)
  (9×) — die Gegenrichtung: Eine Grenze, die mehr Module nennt als gemessen
  wurden, behauptet eine Prüfung, die nicht stattfand. §6 führt es als drittes
  Risiko.
- [`grenzen-liste-wird-als-vollstaendig-gelesen`](../observations/BEO-ALL/grenzen-liste-wird-als-vollstaendig-gelesen/observation.md)
  (8×) — dieser Slice **schreibt** Grenzen-Listen, an drei Orten. Der Ableiter
  gilt doppelt: den Vertrags-Teil umdrehen, und wo der Gegenstand Verhalten
  ist, gegen das Verhalten prüfen statt gegen die Prosa.
- [`form-vom-nachbarn-statt-von-der-vorlage`](../observations/BEO-ALL/form-vom-nachbarn-statt-von-der-vorlage/observation.md)
  (1×, gestern angelegt) — die **Form** des Grenzen-Punkts wird aus dem
  vorhandenen Abschnitt derselben Datei übernommen. Das ist hier **richtig**:
  Der Eintrag warnt vor der Form-Übernahme, wo eine **Vorlage** existiert; für
  einen Punkt in einer bestehenden Liste ist die Liste selbst die Form. Notiert,
  damit die Übernahme eine Entscheidung ist und keine Gewohnheit.

**Keiner der vier erreicht mit diesem Slice die Schwelle erstmalig.**

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-08 gelesen: **beide Nachtläufe grün** —
`upstream-drift.yml` (jüngster Lauf 2026-09-08T05:31:31Z) und `image-scan.yml`
(2026-09-07T08:21:32Z). **Für diesen Slice ausnahmsweise nicht bezugslos:**
Der Gegenstand ist an `go-git v5.19.2` gebunden, und ein Dependabot-Bump dieser
Abhängigkeit könnte das gemessene Verhalten ändern. Die Grenze nennt deshalb
die Version.

**Modus-Begründungsblock.** Alle berührten Sub-Areas GF — ein Block genügt.

### Sub-Area: `*`

- **Modus:** GF (Repo-Default).
- **Konventions-Dichte:** hoch für die **Form** einer Sensor-Grenze — beide
  Zieldateien führen den Abschnitt *Grenze — was das Grün nicht abdeckt* mit
  nummerierten Punkten bereits.
- **Phase-Reife:** Phase 5. Die Sensor-Dateien sind gewachsen und wurden in
  slice-212 zuletzt systematisch auf fehlende Grenzen durchgesehen.
- **Evidenz-/Diskrepanz-Risiko:** **niedrig für den Bestand, mittel für die
  Aussage.** Zu inventarisieren ist nichts. Das Risiko sitzt in der Reichweite:
  drei gemessene Module, aber eine Formulierung, die leicht für alle spricht —
  und in der Erklärung, die schon einmal falsch war.
- **Reconciliation-Aufwand:** keiner (GF).
