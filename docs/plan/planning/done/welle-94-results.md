# Welle welle-94 — Modul `pinsync`: gleiche Versions-Pins in mehreren Dateien gegeneinander prüfen — Closure-Notiz

**Welle:** welle-94
**Abschluss:** 2026-10-31
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-286** (Anforderung, Spezifikation und ADR):
  [`DC-FA-PSYNC-001`](../../../../spec/lastenheft.md#dc-fa-psync-001--gleichheit-eines-wertes-an-mehreren-stellen-modul-pinsync-opt-in)
  im Lastenheft, Verfahren
  [`DC-FA-PSYNC-001.a`](../../../../spec/spezifikation.md#dc-fa-psync-001a--gleichheit-eines-wertes-an-mehreren-stellen-pinsync)
  und [ADR-0113](../../adr/0113-pinsync-modul.md) Accepted. Die Entscheide des Auftraggebers vor dem
  Code festgelegt: kein Schalter `strip-prefix`, nur Gleichheit, kein Netz; dazu die Leseweise
  (zeilenweise, RE2, erste Passung mit nicht leerer Gruppe), Pin ab zwei Quellen, ein Symlink als
  Quelle ist Exit 2.
- **slice-287** (Modulkern): Konfigurationsblock mit Prüfung am Rand (Exit 2), zeilenweises Fangen,
  Vergleich gegen die erste Quelle mit genau einem Treffer, die Grund-Codes `pinsync-mismatch`,
  `pinsync-no-match`, `pinsync-ambiguous`; Registrierung in allen Modul-Aufzählungen und der
  Handbuch-Abschnitt 4.27 mit dem `v?`-Muster.
- **slice-288** (eigenes Profil): zwei Pins zu je zwei Quellen (`GO_VERSION`,
  `GOLANGCI_LINT_VERSION` in `Dockerfile` und `Makefile`) im Gate `make doc-check`, beidseitig belegt.
- **Außerhalb der Slices, in dieser Welle:** Phase (8) in `image-test` — der Wertvergleich unter
  `:ro`, mit fremdem Benutzer und ohne Netz, Exit 0 bei gleichen, Exit 1 bei ungleichen Werten, dieselbe
  Ausgabe wie gegen ein beschreibbares Repo, Symlink-Quelle Exit 2. Der Closure-Trigger verlangt sie; das
  Gate hatte bis dahin keinen `pinsync`-Teil.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Die Entscheide vor dem Code, mit Empfehlung:** der Auftraggeber bestätigte sie in einem Zug; die
  Abweichung (Symlink als Quelle ist Exit 2, `targets` folgt ihm still) stand vor der Anforderung fest,
  weil der Code von `targets` gelesen wurde, ehe sie geschrieben wurde.
- **Das Vorbild `metrics` als Liste der Spiegel:** der Commit des Vorgänger-Moduls zeigte, welche
  Dateien ein neues Modul berührt; Reviewer und Verifier zählten die Aufzählungen unabhängig nach, und
  die Spiegeltests wurden rot gemessen, ehe sie grün waren.
- **Die Brechproben des Verifiers:** über 30 Mutationen in `slice-287` und sieben je Pin und Datei in
  `slice-288`; sie deckten einen CRLF-Teilfall auf, der ohne das Abschneiden des `\r` grün geblieben
  wäre, und zeigten die Trennschärfe jedes Teilfalls.
- **Review in frischem Kontext, danach Verifier, je Slice:** der Review von `slice-287` holte eine
  Grenze (zwei gleichwertige Ausdrücke vergleichen dieselbe Zeile) zurück, die schon im Review von
  `slice-286` als Hinweis stand und nicht aufgenommen worden war.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz, die
daraus schon gezogen wurde.

- **Eine `Accepted`-ADR lässt sich nach ihrem Commit nicht mehr korrigieren** (slice-286): das Review
  änderte den Körper, der Gate-Lauf wies die Korrektur ab; der lokale, ungepushte Commit wurde neu
  gebildet. Konsequenz: Beleg im Register (`adr-accepted-bevor-der-mechanismus-steht`, jetzt 2×).
- **Die Grund-Code-Zeilen der Spezifikation durften nicht vor dem Code stehen** (slice-286): der
  Lockstep-Test hält die Tabelle gegen die Codes; die Zeilen kamen mit slice-287.
- **Die Blackbox-Probe ist mit der heutigen `.d-check.yml` nicht identisch.** Die Fixture `repo`
  liest die Konfiguration des Repos, und das Vorher-Image (`v0.90.0`) kennt den Block `pinsync` nicht
  (Exit 2). Mit der Konfiguration vor `slice-288` ist die Probe byte-identisch über 20 Vergleiche. Das
  ist erwartet für ein neues opt-in-Modul und kein Verschieben eines bisher grünen Laufs.
- **Zwei Vermerke eines Verifiers trafen nicht zu:** `slice-287` V-2 (kein Test für den unbekannten
  Schlüssel; `TestDecode_PinsyncEingabefehler` führt ihn) und `slice-288` V-3 (`FOCUS_DISABLE` ohne
  Wächter; `make test` wird rot, `TestFocusDisable_DecktDCheckYmlModules`). Beide Male fuhr der Verifier
  den engeren Sensor statt des zuständigen. Einzelfälle, kein Register-Beleg.
- **Zwei Prosa-Spiegel der Modulliste** (Makefile-Beschreibung, Modulzahl in der Sensor-Doku) fehlten im
  Plan-Nachtrag von slice-288 und kamen aus dem Review (Beleg im Register).
- **Der Digest-Pin im Handbuch folgt dem Release-Tag erst im Nachgang** (eigener Commit nach der
  Pipeline) — wie bei jedem Release dieses Repos.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen **nur** Beobachtungen, die im Register
3× erreicht haben.

Kein Eintrag über der Schwelle, der noch auf einen Ausgang wartet. Gelesen wurde das Register
vollständig: Einträge mit Stand `offen` und mindestens drei Belegen gibt es keine. Gezählt, nicht
neu verkörpert:

- `BEO-ALL/adr-accepted-bevor-der-mechanismus-steht` steht bei 2×, offen (slice-267, slice-286).
  Tritt sie ein drittes Mal auf, gehört die Frage *steht der Mechanismus, den die ADR beschließt,
  schon, wenn sie `Accepted` wird?* in den Architect-Ablauf.
- `BEO-ALL/gemeinsame-leser-sind-zwei-familien` bleibt bei 2×, offen: der Zeilen-Leser von
  `pinsync` folgt dem Verfahren und gilt nicht als dritter Beleg.
- `BEO-ALL/modulliste-spiegel-ungegated` steht bei 4×, verkörpert (Spiegel-Deckungstests), und bekommt
  einen Beleg (slice-288): zwei Prosa-Spiegel, die der Test bewusst nicht hält.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*: nichts Neues.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur, das
Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine. Offen bleiben benannte Grenzen der Anforderung, keine Arbeit der Welle: eine FIFO,
Gerätedatei oder ein Socket als Quelle lässt den Lauf beim Lesen hängen, statt mit Exit 2 zu enden
(Spezifikation und Handbuch; nicht im Lastenheft); die übrigen Spiegel einer Pin-Hebung bleiben
ungeprüft (`BEO-ALL/pin-bump-mirrors-ungated`).

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices liegen in `done/`; `make fullbuild` Exit 0 am Stand `5b669af0` —
  Image-Hash `sha256:fc0c68c7938e7e866601394bbeffa0ee1fe05e8625c5439306a4230b68745252`
  (`make gates`, `make image-test`, `make bench`, `make completeness-check` und
  `make verify-closure-notes` darin grün; Vollständigkeit: 57 Anforderungen, 0 Waisen).
- **Ohne den Block nichts Neues:** `make blackbox-probe REF=v0.90.0` mit der `.d-check.yml` vor
  slice-288: byte-identisch über 20 Vergleiche (stdout, stderr und Exit getrennt, Kanarienlauf vorher
  und nachher). Mit der heutigen Konfiguration weichen vier Vergleiche der Fixture `repo` ab, weil das
  Vorher-Image den Block `pinsync` nicht kennt (Exit 2); die 16 Vergleiche der übrigen Fixtures sind
  in beiden Fällen gleich.
- **Der Anlass trägt am eigenen Bestand, beidseitig:** `make doc-check` grün mit
  `pinsync: 2 Pin(s) geprüft, 4 Quelle(n) gelesen`. Ein Wert in nur einer Datei gehoben:
  `pinsync-mismatch` an `Makefile:13` (GO_VERSION 1.27.4 gegen 1.27.3) bzw. `Makefile:14`
  (GOLANGCI_LINT_VERSION), Rücknahme grün; `v2.15.0` und `2.15.0` gelten als gleich; Zeile
  auskommentiert `pinsync-no-match`, verdoppelt `pinsync-ambiguous`.
- **Betriebsbedingungen am gebauten Image:** `make image-test` (in `fullbuild`), Phase (8): `:ro`,
  fremder Benutzer (12345:12345), `--network none` — gleiche Werte Exit 0, ungleiche Exit 1 mit
  `pinsync-mismatch`, Ausgabe gleich dem beschreibbaren Lauf, Symlink-Quelle Exit 2.
- Replay-Lauf gegen Golden Set: entfällt, d-check führt keinen (deterministischer Kern,
  [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Coverage gesamt: 94,90 % bei Schwelle 93 % (`make coverage-gate`), offene Carveouts: keine.

**Trigger-Audit (Schritt 2), je Klasse geprüft:**

- **Carveout:** keine aktiven (`docs/plan/carveouts/` trägt nur `done/` und die README).
- **Bootstrap-aware Gate:** `coverage-gate` ist kalibriert (Schwelle 93 %, gemessen 94,90 %,
  Senkung nur per ADR); kein Hochschalt-Trigger gefeuert, keine neue Stufe.
- **ADR:** keine Bedingung feuerte. Gelesen wurde die berührte:
  [ADR-0113](../../adr/0113-pinsync-modul.md) trägt ihren Re-Evaluierungs-Trigger — mehrzeilige Werte,
  Normalisierung der Werte, Quellen über einen Symlink; keiner ist eingetreten. Die Aufnahme von
  `pinsync` in `doc-check` ist keine Senkung eines Gates (zusätzliche Prüfung), die ADR trägt sie
  (`AGENTS.md` §3.6).
- **Hard Rule:** in `AGENTS.md` §3 tragen alle einen Auflösungs-Trigger *permanent* oder benennen die
  ungedeckte Kategorie; §3.8 wurde in slice-286 angewandt (Quellen außerhalb der Scan-Menge benannt).
  Nichts feuerte, keine Zeile wird entfernt.

**Die drei Paarungen:** (a) Anker — kein Steering-Loop-Eintrag mit `liegt in`-Feld, keine Paarung
auszulösen; (b) Folge-Slice — keiner genannt; (c) Register — die zitierten Beobachtungen existieren als
Verzeichnisse und tragen alle Belege.
