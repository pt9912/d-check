# Slice slice-164: Der Nachtlauf hat keinen Adressaten — wer liest das Rot, und wann?

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml)
(der Nachtlauf und seine bereits benannte Grenze);
[slice-142](../done/slice-142-freshness-weitere-achsen.md) (die zwölf Achsen);
[slice-161](../done/slice-161-sechs-pins-heben.md) (die Hebung, die den
dritten §5-Punkt offen ließ); [`MR-051`](../../../../harness/conventions.md#mr-051)
(die Bump-Nachwirkung, die planmäßig meldet).

**Berührte Spec-Stellen:** — (Harness-Betrieb; keine Anforderung).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

Zwölf Achsen wachen über jeden gepinnten Fremd-Bestand. Sie melden **korrekt**
— und **an niemanden**. Der Job fällt rot aus und ist nur in der
Actions-Übersicht des Repos sichtbar; es gibt keinen Push-Kanal, keine
Benachrichtigung, keinen festen Lese-Schritt.

**Die Lücke ist nicht neu und nicht verdeckt:** der Workflow-Kopf trägt sie als
benannte Grenze — *ein dauerroter Nachtlauf ist wieder derselbe verwaiste
Sensor*. [slice-161](../done/slice-161-sechs-pins-heben.md) hat das Rot geräumt
und den Punkt als **weiter offen** an dieses Register zurückgegeben. Genau
deshalb ist er jetzt fällig: solange alles grün ist, kostet die Entscheidung
nichts; beim nächsten Rot ist sie wieder dringend und wieder unbequem.

**Der nächste rote Lauf ist absehbar**, nicht hypothetisch: nach
[`MR-051`](../../../../harness/conventions.md#mr-051) meldet die Zitat-Prüfung
bei jedem Baseline-Bump planmäßig, und Fremd-Releases erscheinen ohne
Ankündigung.

## 2. Vorgehen

1. **Zuerst die Frage, nicht das Werkzeug:** wer ist der Adressat, und in
   welchem Takt soll er lesen? Eine Kadenz ohne benannten Leser ist dieselbe
   Verwaisung eine Ebene höher.
2. **Den vorhandenen Lese-Schritt prüfen, bevor ein neuer gebaut wird.** Die
   Slice-Planung hat bereits einen Sichtungs-Schritt (Beobachtungs-Register).
   Ob der Nachtlauf-Stand dort hineingehört, ist billiger zu haben als jeder
   Benachrichtigungs-Kanal — und braucht kein Geheimnis, keinen Dienst, keine
   zweite Fehlerquelle.
3. **Erst dann die Kanal-Frage**, falls Schritt 2 nicht trägt: was GitHub ohne
   Zusatz-Dienst kann (Actions-Benachrichtigung je Watch-Einstellung, ein Issue
   je rotem Lauf) gegen das, was einen Dienst bräuchte.
4. **Die Rausch-Frage mitentscheiden.** Ein Kanal, der bei jedem Bump meldet,
   erzieht zum Wegklicken. Ob eine planmäßige Meldung
   ([`MR-051`](../../../../harness/conventions.md#mr-051)) anders behandelt
   wird als eine unerwartete, gehört zur Kadenz.
5. Nur bauen, was die Antwort trägt; eine Entscheidung **gegen** einen Kanal
   ist ebenso auszuweisen wie eine dafür.
6. `make gates`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Auto-Bump.** Die Achsen melden; wer hebt, entscheidet.
- **Keine Aufnahme netzhaltiger Achsen in `gates`.** Der Nachtlauf bleibt der
  Bindepunkt.
- **Kein Fremd-Dienst mit Geheimnis**, solange Schritt 2 oder 3 trägt — ein
  Wächter, der ein Token braucht, hat eine neue Ausfall-Achse.

## 4. Definition of Done

- [x] Der **Adressat** ist benannt, nicht nur der Kanal: der Rolleninhaber der
      Implementer-Rolle beim Planen des nächsten Slice.
- [x] Der vorhandene Lese-Schritt ist geprüft, bevor ein neuer gebaut wird — der
      **Moment** ist wiederverwendet (die Slice-Planung), das **Werkzeug** ist
      neu. Das ist nicht dasselbe und war auch nicht verlangt.
- [x] Die Rausch-Frage ist entschieden — **und die erste Fassung entschied sie
      nur zur Hälfte**: sie übertrug die Unterscheidung dem Leser und schob ihm
      für **jede** Rot-Ursache dieselbe Deutung unter, auch für `cancelled` und
      `timed_out`. Jetzt trennt das Werkzeug Achsen-Urteil von Lauf-Störung.
- [x] Was gebaut wird, ist gefahren — mit gelesener Ausgabe: der Live-Lauf, die
      sechs Selbsttest-Proben und die Bruchprobe des Zeitstempel-Formats.
- [x] `make gates` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Ein Kanal ist billig zu bauen und teuer zu ignorieren** — dieselbe Klasse
  wie bei den Achsen selbst, nur eine Ebene höher. — **Ausgang: entfallen.** Es
  wurde keiner gebaut. **Die Begründung dagegen war allerdings überzogen:** sie
  behauptete, *jeder* ohne Fremd-Dienst verfügbare Kanal falle unter einen von
  zwei Zweigen, und nannte das *„gemessen"*. Beides ist zurückgenommen — drei
  betrachtete Kandidaten mit Grund, die Liste ausdrücklich unvollständig, und
  keine Probe behauptet.
- **Ein Lese-Schritt in der Slice-Planung greift nur, wenn ein Slice geplant
  wird.** In einer Pause liest niemand. — **Ausgang: eingetreten, als Grenze
  akzeptiert.** Sie steht in
  [`MR-053`](../../../../harness/conventions.md#mr-053), in
  [`AGENTS.md`](../../../../AGENTS.md) §4 und §5 und in
  [`harness/README.md`](../../../../harness/README.md) §Sensors — und sie ist
  der `Auflösungs-Trigger` des Eintrags: bleibt ein Rot über mehrere
  Slice-Planungen ungelesen, ist die Kanal-Frage neu zu stellen.
- **Die Entscheidung könnte den Zustand mit dem Werkzeug verwechseln.** *„Wir
  bauen eine Benachrichtigung"* beantwortet nicht, wer sie liest. —
  **Ausgang: eingetreten, und zwar im Werkzeug selbst.** Die
  Adressaten-Entscheidung hielt; das **Werkzeug** aber meldete das Feld
  `conclusion` des jüngsten Lauf-Objekts und nannte das *„den Stand des
  Nachtlaufs"* — ohne zu prüfen, ob dieser noch **läuft** (JSON-`null` ⇒
  Falsch-Rot) und ohne sein **Alter** anzusehen (ein abgeschalteter Nachtlauf
  meldete weiter `gruen`). Das erste ist behoben, das zweite als Grenze
  benannt. Genau die Verwechslung, vor der der Punkt warnt — nur eine Ebene
  tiefer als erwartet.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei.

**Rückführungen:** `in-progress` → `next`, falls die Adressaten-Frage ein
Auftraggeber-Entscheid ist.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** CI (GF), Harness-Prozess (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — über die gemessene Menge reden, nicht über
  die naheliegende; [`BEO-011`](../observations.md) — die Regel aus dem
  Bestand, nicht aus dem Anlass.
- **Nachtlauf-Stand lesen** (`make nightly-state`, dritte Vorprüfung nach
  [`MR-053`](../../../../harness/conventions.md#mr-053) — die dieser Slice
  einführt und deshalb auf sich selbst anwendet): **ROT**, jüngster Lauf
  `2026-08-27T10:49:23Z`. **Gelesen, nicht weggeklickt** — und der Beleg ist
  der `head_sha`, nicht die Uhrzeit: der Lauf trägt `48cf132`, und der
  Bump-Commit der sechs Pin-Hebungen aus
  [slice-161](../done/slice-161-sechs-pins-heben.md) ist dessen **direkter
  Nachfolger** (`git rev-list --count 48cf132..b8e5cc1` = 1), lag zum Dispatch
  also noch nicht auf `origin/main`. Der Lauf hat die Hebungen nicht gesehen.
  Alle zwölf Achsen melden seither einzeln `ok`; **dass der nächste Lauf grün
  ausfällt, ist damit erwartet und nicht gemessen.**

Slice-ID: slice-164. Betroffene IDs: — (Harness-Betrieb; keine Anforderung).
Module: CI, Harness-Prozess. Gates: `make gates`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Betriebs-Entscheid an vorhandener Mechanik.

## 9. Closure-Notiz (nach `done/`)

**Der Nachtlauf hat einen Adressaten — und das Werkzeug, das ihn bedienen soll,
hat dieselbe Verwechslung begangen, vor der §5 warnt.**

**Die Reihenfolge, die §2 verlangt, ist eingehalten: Adressat, dann Takt, dann
Kanal.** Adressat ist der **Rolleninhaber der Implementer-Rolle beim Planen des
nächsten Slice** — nicht *„das Team"*, nicht *„wer hinsieht"*. Takt ist jede
Slice-Planung, als dritter `Vorgelagert`-Block neben den zwei des Kanons. Der
Moment ist wiederverwendet, weil ein neuer dieselbe Verwaisung eine Ebene höher
erzeugt hätte.

**Kein Kanal — und die Begründung dafür war im ersten Anlauf überzogen.** Sie
sagte, *jeder* ohne Fremd-Dienst verfügbare Kanal falle unter einen von zwei
Zweigen, und nannte das *„gemessen"*. Beides hält nicht: ein Status-Badge fällt
unter keinen der beiden, gegen die Issue-Option argumentierte ich in ihrer
**teuersten** Form statt in der milden, und **gemessen** war gar nichts — kein
Kanal wurde aufgesetzt. Die Entscheidung bleibt; ihre Verpackung als
vollständig und gemessen ist zurückgenommen.

**Das Werkzeug liest ohne Fremd-Werkzeug.** `make nightly-state` fragt die
GitHub-API mit `curl` — derselben Erwartung, die die Netz-Skripte ohnehin
tragen. `gh` wäre eine neue gewesen; dass
[`AGENTS.md`](../../../../AGENTS.md) §3.1 sie bis eben nur für **zwei** Skripte
aussprach und *„Beide"* sagte, hatte ich übersehen, während der Kopf meines
Skripts sich auf genau diesen Absatz berief. Derselbe Befund war am selben Tag
schon einmal geschlossen worden.

**Der schwerste Befund traf das Werkzeug an seinem Zweck.** Die API liefert bei
laufendem Job `"conclusion": null` — nicht ein fehlendes Feld. Mein SKIP-Zweig
war damit **toter Code**, und ein laufender Nachtlauf wurde als `ROT (null)`
gemeldet: ein Falsch-Rot in genau dem Werkzeug, dessen erklärter Zweck das
Verhindern von Wegklicken ist. Dazu bekam das **ganze** Nicht-`success`-
Vokabular den *„planmäßige Meldung"*-Rat, auch `cancelled` und `timed_out`, wo
er schlicht falsch ist. Beides behoben; das Werkzeug trennt jetzt Achsen-Urteil
von Lauf-Störung und prüft die **Form** des Zeitstempels, bevor ein Wert als
Stand gilt.

**Der netzlose Prüfeinstieg hat sich sofort bezahlt.** `--parse` und
`--selftest` waren im ersten Anlauf nicht da — die Schwester-Skripte tragen
beides mit ausgeschriebener Begründung, meines nicht. Nachgeholt, und der
**erste** Selbsttest fand einen Fehler, den der Live-Lauf verdeckte: mein `sed`
streifte ein schließendes `"`, aber nicht `"}`. Live fiel das nicht auf, weil
`conclusion` dort nicht das letzte Feld ist.

**Der Zeit-Beleg in §7 widerlegte, was er stützen sollte.** Ich verglich `10:49`
(UTC aus der API) gegen `12:46` (`+0200` aus `git log`). In **einer** Zeitbasis
liegt die Pin-Hebung **2 min 49 s vor** dem Lauf, nicht zwei Stunden danach.
Die Schlussfolgerung stimmt trotzdem — der Lauf trägt `head_sha 48cf132`, und
der Bump-Commit ist dessen **direkter Nachfolger**, lag zum Dispatch also noch
nicht auf `origin/main`. Der Beleg ist ersetzt, und *„durch slice-161 behoben"*
steht jetzt als das da, was es ist: eine **Vorhersage** über den nächsten Lauf.

**Und die Regel des Nachbar-Slice hat sich im selben Commit bewährt.** Die
Änderung an §3.1 verschob die Zeilenspanne der `AGENTS.md`-Zitat-Direktive aus
[slice-163](../done/slice-163-zitate-ausserhalb-des-speichers.md);
`citations` meldete `citation-mismatch`, und
[`MR-051`](../../../../harness/conventions.md#mr-051) sagt, was zu tun ist:
**neu ankern, nicht löschen.** Der Fall, den jener Slice als Grenze benannt
hatte, ist beim ersten Anlass eingetreten — und die Regel trug.

**Vier Grenzen sind benannt statt entdeckt zu werden:** die Pause; das
ungeprüfte **Alter** des Laufs (ein abgeschalteter Nachtlauf meldete weiter
`gruen`); der bei privatem Repo von einer Netzstörung ununterscheidbare `SKIP`;
und der Repo-Slug als **Default** statt Fund — in einem Fork meldete das
Werkzeug sonst den Nachtlauf des Originals.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 563 Dateien, 0 Befunde),
`make fullbuild` (Exit 0, 48 Anforderungen / 0 Waisen, bench Median 695 ms),
`--selftest` (sechs Proben, 0 Fehlschläge), `make nightly-state` gegen den
echten Lauf mit gelesener Ausgabe. Ein unabhängiger Review ist gelaufen; sein
Urteil war *„schließbar nach Nacharbeit"*, seine vierzehn Befunde sind
eingearbeitet, und seine vier HIGH sind eigens nachgemessen statt übernommen.
Bemerkenswert ist, was er **nicht** brechen konnte: das Komma-Parsing hält, weil
das Muster an beide Anführungszeichen verankert ist — der Fehler saß in der
fehlenden Wert-Prüfung, nicht in der Trennung.
