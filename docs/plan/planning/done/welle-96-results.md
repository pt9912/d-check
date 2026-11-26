# Welle welle-96 — `vcs` folgt einer Umbenennung: ein Link-Nachzug auf die umbenannte Datei ist keine Drift — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v6.18.1` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt).

**Welle:** welle-96
**Abschluss:** 2026-11-24
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — Ergebnis, nicht Tätigkeit.

- **slice-291** (Anforderung, Spezifikation, ADR): Lastenheft
  [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
  von 0.113.2 auf 0.114.3 — neuer Absatz „Umbenennung folgen (opt-in)“, eine Liste benannter
  Grenzen, 30 Akzeptanzkriterien (25 neue), geschärfter Out-of-Scope-Satz; Verfahren in der
  Spezifikation (Schritt 2 und 4 von
  [`DC-FA-VCS-001.a`](../../../../spec/specification.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
  und die Schema-Zeile `vcs.follow-renames`); [ADR-0115](../../adr/0115-vcs-folgt-umbenennung.md)
  mit 14 Entscheidungen. Entscheidungen 1–8 bestätigte der Auftraggeber am 2026-11-23, die
  Folge-Entscheidungen 9–14 und die neue Fassung von 9 und 10 am selben Tag, Entscheidung 9
  nochmals am 2026-11-24 mit ihren drei nachsichtigeren Fällen.
- **slice-292** (Kern und git-Adapter): `vcs.follow-renames` (Schalter, Default aus). Der Adapter
  liest die Umbenennungen der Range mit **eigenem Blob-Vergleich** — die Umbenennungs-Erkennung
  von go-git wird nicht benutzt —, die Kette läuft entlang der Abstammung, Mehrdeutigkeit und
  „Basis kein Vorfahr“ sind geregelt, ein unlesbares Objekt ist Exit 2. Der Kern setzt Link-Ziele
  über die Abbildung gleich. Im gestagten Modus zählen jetzt nur noch Datei, ausführbare Datei
  und Symlink (Angleichung an die Spezifikation). Handbuch-Abschnitt und
  `harness/sensors/adr-check.md` sind nachgezogen.
- **Profil und ADR-Übergang:** [ADR-0115](../../adr/0115-vcs-folgt-umbenennung.md) ist `Accepted`;
  im selben Commit schaltet `.d-check.yml` `follow-renames: true` ein (Entscheidung 14), und
  `AGENTS.md` §3.5 und `harness/README.md` nennen sie neben
  [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md).
- **Harness:** neue Dokumentations-Regel 18 in `AGENTS.md` §5 (Steering-Loop-Eintrag unten).
- **Anlass:** Issue #18 des Kurs-Repos — eine Layout-Migration, die `spec/spezifikation.md`
  umbenennt, lässt die `Accepted`-ADRs mit Link darauf nicht mehr rot werden.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Die Rollenfolge lief von Anfang an regelkonform:** Architect-Verdikt vor `in-progress` bei
  beiden Slices, die ADR schrieb nur der Architect, jede Entscheidung in der Sache ging an den
  Auftraggeber.
- **Die Mutation fand, was die Lektüre nicht fand.** In slice-292 setzten Implementer (22),
  Reviewer (22) und Verifier (11) die Mutationen je einzeln; sie deckte beide Test-Lücken auf
  (die Commit-Menge „von HEAD erreichbar, von BASE nicht“ und die doppelte Marke) und zeigte, welcher
  Zweig nur äquivalente Optimierung ist. Der Code brauchte danach zwei Review-Runden und keinen
  Befund im Produktcode.
- **Die Verifikation, die gegen die Spec liest, und das Vorher-Image:** Weil die Blackbox-Probe das
  Modul `vcs` nicht fährt, verglich der Verifier es selbst gegen ein Image aus `v0.92.0`
  (Zahlen unten).
- **Zitieren statt nacherzählen** beendete die Plan-Fehler zu Entscheidung 9; daraus wurde Regel 18.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile möglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Die Spezifikation war ein Vielfaches des Anlasses.** Issue #18 nennt vier Akzeptanzkriterien;
  daraus wurden 30 Kriterien, 14 Entscheidungen und rund 40 Testfälle, vier Review-Runden und rund
  acht Stunden bis zum Spec-Stand. Ursache war die empfohlene Kette über jeden Commit der Range und
  die Kombination mit `ignore-link-targets`; fast jeder Befund ab Runde 2 entstand dort. Der
  Auftraggeber beanstandete den Aufwand, das Zurückschneiden der Spec verwarf er. Der Code folgte der
  fertigen Spec: Implementer-Lauf rund 46 Minuten, die beiden Review-Runden rund 22 Minuten; ob ein
  kleinerer Umfang schneller gewesen wäre, ist nicht gemessen. Konsequenz: Beobachtung
  [`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](../observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md).
- **Zweimal stand die Zusammenfassung des Planners gegen Lastenheft und ADR** (die Richtung von
  Entscheidung 9, dann ihre drei Fälle). Beide Male trug sie eine Vorlage an den Auftraggeber, der
  deshalb dreimal bestätigen musste. Konsequenz: Dokumentations-Regel 18.
- **Drei Commit-Botschaften behaupteten mehr, als ihr Diff trug** (`e3e5d4eb`, `e84632b1` und eine
  Planner-Botschaft); berichtigt in späteren Botschaften, nicht amendet.
- **Die Blackbox-Probe fährt `vcs` nicht** (Default-`PROBE_FORMS` ohne `--range` und `--staged`);
  der Closure-Trigger „ohne den Schalter nichts Neues“ trug deshalb nicht über die Probe, sondern
  über den eigenen Vergleich des Verifiers. Dazu weicht der Eigenlauf der Probe gegen `v0.92.0` ab,
  seit das Profil `follow-renames` einschaltet: ältere Images lehnen den Schlüssel ab (Exit 2).
- **Die Lastenheft-Historie** wurde zuerst in der unveröffentlichten Zeile 0.114.0 fortgeschrieben
  und nach Review auf eigene Zeilen 0.114.1 bis 0.114.3 umgestellt.
- **Ein Altbefund wurde sichtbar und angeglichen:** die Spezifikation versprach für `--staged` nur
  die drei Datei-Arten, der Code nahm jeden Index-Eintrag. Ein Gitlink, der eine immutable Datei
  ersetzt, endete dort mit Exit 2 und meldet jetzt `core-drift-vcs`, wie im `--range`-Modus — die
  einzige Abweichung ohne den Schalter.
- **Hinweis für die spätere Umbenennung der eigenen Spec-Dateien** (nicht Teil dieser Welle): der
  Verifier zeigte im Wegwerf-Klon, dass sie außer den ADR-Links auch den Einbettungs-Pfad im
  Go-Quelltext und `.d-check.yml` braucht.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<KUERZEL>/<slug>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks).

- **Hard Rule (Dokumentations-Regel 18)** ergänzt: eine Festlegung im Slice-Plan nennt die Stelle
  des Lastenhefts oder der ADR, die sie trägt, und zitiert sie, statt sie nachzuerzählen; trägt keine
  Stelle sie, ist sie eine Entscheidung, keine Notiz — liegt in `harness/rules/dokumentations-regeln/18-plan-festlegung-zitiert-ihre-quelle.md`.
  Auslöser:
  [`BEO-ALL/plan-entscheid-steht-gegen-lastenheft`](../observations/BEO-ALL/plan-entscheid-steht-gegen-lastenheft/observation.md)
  (slice-289, slice-290, slice-291 — 3×). Auflösungs-Trigger: permanent.

Keine weitere Beobachtung erreicht in dieser Welle die Schwelle 3× ohne Ausgang.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- Keine neuen Slices. Die Welle für Issue #20 (`linelength` nimmt eingerückte Tabellenzeilen aus)
  steht in der Roadmap unter *Nächste Wellen*, ihr Trigger ist mit dieser Closure erfüllt; ihre
  Slices entstehen bei ihrer Eröffnung.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

Beleg des Verifiers: [`docs/reviews/2026-11-24-welle-96-verifikation.md`](../../../reviews/2026-11-24-welle-96-verifikation.md);
Trigger-Audit des Architects: [`docs/reviews/2026-11-24-welle-96-trigger-audit.md`](../../../reviews/2026-11-24-welle-96-trigger-audit.md).

- Beide Slices liegen in `done/`; `make fullbuild` Exit 0 auf `274c6ba4` — Image-Hash
  `sha256:6d1195bcbb53a57ec3e72056b0d5f680332c9af080e2df7dfcdebaa494e6eadb` (Coverage 94,90 % bei
  Schwelle 93 %, Semgrep 0 findings, `make bench` mit Median 1142 ms, `make completeness-check` mit 58
  Anforderungen und 0 Waisen, `make verify-closure-notes` darin grün). Danach kamen nur noch
  Dokumente und Beobachtungen hinzu (`66eb94a3`); `make gates`, `make image-test`, `make coverage-gate`,
  `make completeness-check` und `make verify-closure-notes` liefen dort erneut mit Exit 0.
  `make fullbuild` selbst wurde auf diesem letzten Stand nicht wiederholt.
- **Ohne den Schalter nichts Neues.** Eigener Vergleich des Moduls `vcs` gegen ein Image aus
  `v0.92.0` (28 Fixtures ohne `follow-renames`, stdout, stderr und Exit getrennt): 60 Vergleiche,
  56 gleich, 4 weichen ab — alle vier derselbe Fall, `--staged` mit einem Gitlink, der eine immutable
  Datei ersetzt: vorher Exit 2 (`object not found`), jetzt Exit 1 mit `core-drift-vcs`. Das ist die
  einzige Abweichung und sie ist im Welle-Plan benannt. `make blackbox-probe REF=v0.92.0`: Fixtures
  `ids`, `links`, `sauber`, `targets` gleich; der Eigenlauf weicht in allen vier Formen ab (Exit 2
  vorher, 0 nachher), weil das Vorher-Image `follow-renames` im Profil ablehnt. In einem
  Wegwerf-Klon ohne diese eine Profilzeile: byte-identisch über 20 Vergleiche.
- **Der Anlass trägt am gebauten Image.** Ein Commit benennt `spec/spezifikation.md` per `git mv` um
  und zieht in 63 ADRs (60 `Accepted` mit Link, [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md) trägt ihn nur im ausgenommenen Abschnitt
  `Geschichte`) nur das Link-Ziel nach:

  | Probe | Ergebnis |
  |---|---|
  | ohne `follow-renames` (`ignore-link-targets` an) | Exit 1, 59 `core-drift-vcs`, je ADR einer |
  | mit `follow-renames: true` | Exit 0, 0 Befunde |
  | Gegenprobe: Datei gelöscht, mit anderem Inhalt neu angelegt, mit Schalter | Exit 1, dieselben 59 ADRs |
  | Vorher-Image `v0.92.0`, Profil ohne Schalter | Exit 1, 59 Befunde, Ausgabe byte-gleich zum Nachher-Image ohne Schalter |

- **Laufzeit.** `make adr-check` über 100 Commits mit einem Merge, dessen eine Seite vor BASE
  abzweigt: Median 1,139 s ohne, 1,562 s mit `follow-renames` (fünf Läufe je Konfiguration), Faktor
  1,37 gegen die Schwelle 2.
- **Trigger-Audit:** ADR 0 offen ([ADR-0115](../../adr/0115-vcs-folgt-umbenennung.md), [ADR-0103](../../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md), [ADR-0024](../../adr/0024-vcs-immutable-gate.md), [ADR-0109](../../adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md)), Hard Rule 0 offen; die eine
  Lockerung der Prüfregel des eigenen ADR-Gates trägt
  [ADR-0115](../../adr/0115-vcs-folgt-umbenennung.md) Entscheidung 7 und sie wurde im selben Commit
  `Accepted`. Carveouts: keine offen. Bootstrap-aware Gates: keines berührt (`coverage-gate` unverändert).
