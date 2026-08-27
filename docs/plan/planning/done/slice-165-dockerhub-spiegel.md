# Slice slice-165: Der Spiegel wird eine Zusage, nicht ein Nebenprodukt

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** — **wellenlos**, solange keine Closure-Bedingung über die eigene DoD
hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(dessen Out-of-Scope den Spiegel heute ausschließt);
[ADR-0002](../../adr/0002-distribution-ghcr-image.md) (GHCR als einziger Weg);
[ADR-0014](../../adr/0014-latest-tag-fuer-stabile-releases.md) (`:latest` nur
stabil); [`release.yml`](../../../../.github/workflows/release.yml).

**Berührte Spec-Stellen:** [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Out-of-Scope-Satz), [`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel) (neu).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-27.

---

## 1. Ziel

**`docker.io/pt9912/d-check` existiert und trägt dieselben Bilder wie GHCR — als
zugesagte Distribution, nicht als Nebenprodukt.**

Der Auftraggeber hat den Rang ausdrücklich gewählt: **zugesagt**, nicht
Komfort. Das hat einen Preis, und er ist Teil der Entscheidung — **ein
fehlgeschlagener Hub-Push bricht das Release**, auch wenn GHCR längst grün ist.
Der Spiegel darf nicht still zurückfallen; genau das unterscheidet ihn vom
fail-open-Muster des Schwester-Repos.

**Der Vertrag steht heute dagegen.**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
nennt `ghcr.io/pt9912/d-check` namentlich und schließt im Out-of-Scope
*„Distributionswege jenseits GHCR"* aus. Der Satz ist nicht nebenbei zu
übergehen — er ist die Stelle, an der die Zusage heute endet.

## 2. Vorgehen

1. **Den Vertrag zuerst, dann die Pipeline.** Neue Anforderung
   [`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
   mit Akzeptanzkriterien-Trio; der Out-of-Scope-Satz von
   [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
   wird auf die verbleibenden Wege (Homebrew, Paketmanager)
   zurückgeschnitten. Anforderungen entstehen nur im Lastenheft, nie per ADR
   ([`AGENTS.md`](../../../../AGENTS.md) §5).
2. **Die Zusage muss prüfbar sein, sonst ist sie Prosa.** Die tragende Größe ist
   der **Digest**: ein Spiegel, der neu baut, liefert ein anderes Bild unter
   demselben Tag. Zugesagt wird deshalb Digest-**Gleichheit**, und die Pipeline
   misst sie nach dem Push, statt sie anzunehmen.
3. **ADR für die Entscheidung**, nicht für die Anforderung: zweite Registry,
   fail-closed, Spiegel-Verfahren (dasselbe lokale Bild neu taggen statt neu
   bauen), verglichene Alternativen.
4. **Erst dann die Doku** — beide READMEs, Handbuch, `operations.md`,
   `releasing.md`. Eine neue Distribution ist eine Zusage an Konsumenten und
   gehört in jede Oberfläche, die den Aufruf zeigt.
5. `make gates`, `make ci`; unabhängiger Review; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Kein Wechsel des Primär-Registrys.** GHCR bleibt die Quelle; Docker Hub ist
  der Spiegel. Die Richtung ist Teil der Zusage, nicht offen.
- **Keine Binary-Distribution.** Der vertagte Goreleaser-Weg aus
  [ADR-0002](../../adr/0002-distribution-ghcr-image.md) Punkt 5 bleibt vertagt;
  sein Trigger ist ein anderer.
- **Kein Anlegen fremder Konten-Artefakte durch den Lauf.** Das Hub-Repo und die
  zwei Secrets entstehen im Konto des Auftraggebers — der Slice liefert die
  Pipeline, die sie benutzt, und benennt sie als Vorbedingung.

## 4. Definition of Done

- [x] [`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
      steht im Lastenheft mit Akzeptanzkriterien-Trio; der Out-of-Scope-Satz von
      [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
      ist zurückgeschnitten. **Die Anforderung musste ein zweites Mal
      geschrieben werden** (0.71.0), weil ihre erste Fassung die falsche
      Prüfgröße zusagte.
- [x] Die Pipeline spiegelt **fail-closed** und **misst** die Gleichheit — aber
      nicht die, die zuerst zugesagt war: der **Config**-Digest, aus **beiden
      Registries** gelesen. Zwei aus demselben lokalen Bild abgeleitete Werte
      wären trivial gleich gewesen.
- [x] Die Vorbedingungen stehen in
      [`releasing.md`](../../../../docs/user/releasing.md) §Vorbedingungen —
      samt der Folge, dass ohne sie **jedes** Release fehlschlägt.
- [x] Beide READMEs, Handbuch und `operations.md` zeigen den zweiten Bezugsweg.
      `operations.md` fehlte im ersten Anlauf und stand im selben DoD-Haken.
- [x] `make gates` und `make ci` grün (Exit explizit); unabhängiger Review,
      Urteil *„nicht schließbar"*, vierzehn Befunde eingearbeitet.

## 5. Abnahme-Punkte / Risiken

- **Fail-closed bindet an eine fremde Verfügbarkeit.** Ist Docker Hub gestört
  oder das Token abgelaufen, ist das Release rot — obwohl das Bild auf GHCR
  liegt und gültig ist. Das ist der gewählte Preis, aber er ist erst dann
  ehrlich, wenn die Fehlermeldung sagt, **was schon veröffentlicht ist**. —
  **Ausgang: eingetreten, und die erste Fassung löste ihn nicht ein.** Die
  Zusage stand im Text, aber die Zugangsdaten wurden in einem `uses:`-Schritt
  verbraucht, der mit seiner eigenen Meldung scheitert; der `trap` mit dem
  GHCR-Digest stand im **nächsten** Schritt und wäre nie angelaufen. Der Login
  ist jetzt ein `run`-Schritt mit Vorab-Prüfung — die Meldung trägt den
  Teil-Zustand in beiden Fällen, fehlend **und** ungültig.
- **Digest-Gleichheit ist eine Annahme, bis sie gemessen ist.** `docker tag` +
  `push` sollte denselben Manifest-Digest liefern; „sollte" ist keine Zusage. —
  **Ausgang: eingetreten, und es ist der teuerste Befund dieses Slice.** Der
  Punkt stand hier, richtig formuliert, und wurde trotzdem **in** der Pipeline
  beantwortet statt **vor** ihr. Ein unabhängiger Review hat gemessen, was hier
  offen stand: drei von drei Tag-Paaren des Schwester-Repos tragen
  **verschiedene** Manifest-Digests bei identischem Config-Digest — neu
  komprimierte Layer-Blobs. Die fail-closed-Prüfung hätte **jedes** Release
  gebrochen, nach erfolgreichem Spiegel und vor der Release-Anlage.
  [ADR-0065](../../adr/0065-spiegel-gleichheit-ist-der-config-digest.md) ersetzt
  ihre Vorgängerin — die Kette bleibt ADR-intern.
- **Der Spiegel verdoppelt die Oberfläche, auf der ein Pin veralten kann.**
  `versions` hält heute `ghcr`-präfixierte Pins gegen `version.md#aktuell`; ein
  `docker.io`-Pin fällt nicht darunter und driftet still. —
  **Ausgang: eingetreten, als Grenze akzeptiert.** Es sind jetzt **zwei** bare
  Tags ohne `ghcr`-Präfix, und beide stehen namentlich auf der
  Release-Prep-Liste in [`releasing.md`](../../../../docs/user/releasing.md) §4
  — der Auffanglinie, die dieses Repo für genau diese Klasse führt. Die
  Muster-Fläche des Moduls zu weiten ist ein eigener Entscheid und steht in
  [ADR-0065](../../adr/0065-spiegel-gleichheit-ist-der-config-digest.md)
  §Konsequenzen als ungedeckt.

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei; Auftraggeber-Entscheid über
den Rang liegt vor (zugesagt statt Komfort).

**Rückführungen:** `in-progress` → `open`, falls die Vorbedingungen im
Auftraggeber-Konto nicht herstellbar sind — dann ist eine zugesagte Distribution
nicht lieferbar, und die ehrliche Antwort ist der Komfort-Spiegel als eigener
Entscheid, nicht eine Zusage ohne Deckung.

**Der Trigger ist nicht eingetreten, und er ist noch offen.** Hub-Repository und
die zwei Secrets liegen im Konto des Auftraggebers und sind zum Closure-Zeitpunkt
nicht gesetzt. Der Slice liefert die Pipeline, die sie benutzt, und benennt sie
als Vorbedingung — **die Probe steht damit aus**: das nächste Release ist der
erste Lauf, der den Spiegel und den Digest-Vergleich wirklich fährt. Bis dahin
ist die Zusage geschrieben und gate-geprüft, aber nicht durchgeführt.

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Distribution/Release-Pipeline (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-27):
  [`BEO-020`](../observations.md) — die gezählte Menge benennen, bevor die Zahl
  fällt; hier: was der Digest-Vergleich vergleicht (Manifest, nicht Inhalt);
  [`BEO-011`](../observations.md) — die Regel aus dem Bestand, nicht aus dem
  Anlass: das Schwester-Repo ist **ein** Muster, kein Kriterium.
- **Nachtlauf-Stand lesen** (`make nightly-state`, dritte Vorprüfung nach
  [`MR-053`](../../../../harness/conventions.md#mr-053)): **ROT**, jüngster Lauf
  `2026-08-27T10:49:23Z`. **Gelesen:** derselbe Lauf, den
  [slice-164](../done/slice-164-nachtlauf-kadenz.md) und
  [slice-160](../done/slice-160-reviewer-skill-kontextlast.md) bereits
  eingeordnet haben — er lief vor den sechs Pin-Hebungen aus
  [slice-161](../done/slice-161-sechs-pins-heben.md). Keine neue Meldung; die
  Probe darauf steht weiter aus.

Slice-ID: slice-165. Betroffene IDs:
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel),
[ADR-0065](../../adr/0065-spiegel-gleichheit-ist-der-config-digest.md).
Module: — (Pipeline, Spec, Doku). Gates: `make gates`, `make ci`.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — die zweite Registry ist ein Zuwachs an einer
vorhandenen Pipeline; kein Fremdsystem, keine Reconciliation.

## 9. Closure-Notiz (nach `done/`)

**Der Spiegel steht als Zusage — und die Zusage, die zuerst dastand, war die
falsche Größe.**

**Der Vertrag stand dagegen, und das war der erste Fund.**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
nennt `ghcr.io/pt9912/d-check` namentlich und schloss im Out-of-Scope
*„Distributionswege jenseits GHCR"* aus — dieser Slice wäre vertragswidrig
gewesen. Der Satz ist zurückgeschnitten, nicht stillschweigend überholt.

**Die tragende Zusage war am Bestand widerlegt, und die Prüfung hätte jedes
Release gebrochen.** Zugesagt war derselbe **Manifest**-Digest. Gemessen am
Schwester-Repo, das mit identischem Verfahren spiegelt: gleiche Image-ID,
gleicher Config-Digest — und **drei von drei** Tag-Paaren mit verschiedenen
Manifest-Digests, weil der zweite Push die Layer-Blobs neu komprimiert. Beide
Pushes gelingen, der Vergleich schlägt fehl, und weil die Release-Anlage danach
steht, existierte das Release nicht, während der Spiegel öffentlich liegt.
[ADR-0065](../../adr/0065-spiegel-gleichheit-ist-der-config-digest.md) ersetzt die
erste Fassung — die Supersedes-Kette bleibt ADR-intern —: die Gleichheit
ist der **Config**-Digest, aus **beiden Registries** gelesen.

**§5 hatte die Frage richtig gestellt.** *„`docker tag` + `push` sollte denselben
Manifest-Digest liefern; ‚sollte' ist keine Zusage"* — der Satz stand hier, und
ich habe ihn **in** der Pipeline beantwortet statt **vor** ihr. Der Fehler saß
nicht im Nichtwissen, sondern in der Reihenfolge. Das ist die Lehre dieses Slice
und der Grund, warum sein Review teuer war.

**Zweimal am Bestand abgelesen statt in der Regel nachgeschlagen — beide Male
falsch.** Die Supersede-Form habe ich von einer Nachbar-ADR abgelesen (dort
`Status: Accepted` plus Index-Vermerk) statt aus der Vorlage; der Kanon schreibt
`Superseded by ADR-NNNN` in der **Datei**. Und beim Korrigieren habe ich den Wert
verlinkt — die Konfiguration kodiert die bare Form längst (`vcs.head-allow`), und
der `pre-commit`-Hook hat es als `core-drift-vcs` abgewiesen. Mein vorheriger
`STAGED=1`-Lauf war grün, weil nichts gestaged war: eine Probe, die nichts prüfte.

**Die Korrektur legte eine Abweichung offen, die kein Schlendrian war.** Zehn
ADRs tragen einen Supersede-Vermerk im Index, **eine** trägt ihn in der Datei.
Die Ursache ist die Konfiguration: `matrix.status.forbidden` verbot **jede**
Referenz auf eine superseded ADR — auch die ADR→ADR-Lineage, die der Kanon
ausdrücklich zulässt. Wer den Status korrekt setzte, brach damit die
Supersede-Kette. Das Produkt kann den Fall längst
(`allow-supersede-lineage`); die Konfiguration schaltete ihn nur nicht ein.
Jetzt eingeschaltet — und `matrix-inactive` feuert seither korrekt. **Die neun
übrigen bleiben offen** und sind ein eigener Schnitt: ob *„Skript-Mechanik
abgelöst"* ein vollständiges Supersede ist, ist Urteil, kein `sed`.

**Der ADR-Index folgt jetzt seiner Vorlage.** Ihm fehlte die Aussage, die den
ganzen Befund erklärt: *„Quelle der Wahrheit sind die ADR-Dateien; dieser Index
ist eine Bequemlichkeits-Sicht."* Das Supersede nur im Register zu führen dreht
Wahrheit und Sicht um. Dazu die zwei Sätze, die wir gar nicht führten — der
Referenz-Richtungs-Zeiger und *„Prozess-ADRs ohne Spec-Stratum tragen `—`"*.

**Die Darstellung ist ausdrücklich keine Zusage.** Kurztext und Overview kommen
aus [`packaging/dockerhub/`](../../../../packaging/dockerhub/README.md), beide
Release-Schritte tragen `continue-on-error`, und das Zeichen-Limit prüft
stattdessen `make gates` als Go-Test — dort fällt es beim **Schreiben** auf statt
beim Veröffentlichen, und es zählt Zeichen statt Bytes. Das Bild ist die Zusage,
der Beschreibungstext ist Präsentation.

**Was aussteht, steht in §6:** Hub-Repository und die zwei Secrets liegen im
Auftraggeber-Konto und sind nicht gesetzt. Die Zusage ist geschrieben und
gate-geprüft, aber **nicht durchgeführt** — das nächste Release ist ihr erster
Lauf.

**Sensors:** `make gates` (Exit 0, zehn Glieder, 570 Dateien, 0 Befunde),
`make ci` (Exit 0, image-test alle vier Fälle), `make completeness-check`
(Exit 0, 49 Anforderungen / 0 Waisen), `make adr-check STAGED=1`,
`make hubdesc-pin-freshness` (`ok` gegen 5.0.0), und die Digest-Gegenprobe über
beide Registries mit gelesener Ausgabe. Ein unabhängiger Review ist gelaufen;
sein Urteil war *„nicht schließbar"*, seine vierzehn Befunde sind eingearbeitet,
und seine zwei HIGH sind eigens nachgemessen statt übernommen. Bemerkenswert ist,
was er **nicht** brechen konnte: der `trap` feuert bei explizitem `exit` nicht
(keine Doppelmeldung), das `${VAR#*@}`-Stripping trägt beidseitig, und das
Boundary-Kriterium — Prerelease lässt `:latest` unberührt — war von Anfang an
gedeckt.
