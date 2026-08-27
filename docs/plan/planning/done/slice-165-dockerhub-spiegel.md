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

- [ ] [`DC-FA-DIST-002`](../../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
      steht im Lastenheft mit Akzeptanzkriterien-Trio; der Out-of-Scope-Satz von
      [`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
      ist zurückgeschnitten, nicht stillschweigend überholt.
- [ ] Die Pipeline spiegelt **fail-closed** und **misst** die Digest-Gleichheit,
      statt sie zu behaupten.
- [ ] Die Vorbedingungen (Hub-Repo, zwei Secrets) sind in `releasing.md` benannt
      — samt dem, was ohne sie passiert.
- [ ] Beide READMEs, Handbuch und `operations.md` zeigen den zweiten Bezugsweg.
- [ ] `make gates` und `make ci` grün (Exit explizit); unabhängiger Review.

## 5. Abnahme-Punkte / Risiken

- **Fail-closed bindet an eine fremde Verfügbarkeit.** Ist Docker Hub gestört
  oder das Token abgelaufen, ist das Release rot — obwohl das Bild auf GHCR
  liegt und gültig ist. Das ist der gewählte Preis, aber er ist erst dann
  ehrlich, wenn die Fehlermeldung sagt, **was schon veröffentlicht ist**.
  — **Ausgang:** *(bei Closure)*
- **Digest-Gleichheit ist eine Annahme, bis sie gemessen ist.** `docker tag` +
  `push` sollte denselben Manifest-Digest liefern; „sollte" ist keine Zusage.
  Misst die Pipeline nicht nach, steht in der Anforderung eine Eigenschaft, die
  niemand prüft. — **Ausgang:** *(bei Closure)*
- **Der Spiegel verdoppelt die Oberfläche, auf der ein Pin veralten kann.**
  `versions` hält heute `ghcr`-präfixierte Pins gegen `version.md#aktuell`; ein
  `docker.io`-Pin fällt nicht darunter und driftet still — dieselbe Klasse wie
  das bare Tag im Handbuch. — **Ausgang:** *(bei Closure)*

## 6. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei; Auftraggeber-Entscheid über
den Rang liegt vor (zugesagt statt Komfort).

**Rückführungen:** `in-progress` → `open`, falls die Vorbedingungen im
Auftraggeber-Konto nicht herstellbar sind — dann ist eine zugesagte Distribution
nicht lieferbar, und die ehrliche Antwort ist der Komfort-Spiegel als eigener
Entscheid, nicht eine Zusage ohne Deckung.

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
