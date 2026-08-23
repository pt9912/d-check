# Slice slice-125: Release-Prep über drei Erweiterungen und Release `v0.63.0`

**Lifecycle:** Der Zustand dieses Slice ist das **Verzeichnis** (`open/`/`next/`/
`in-progress/`/`done/`), bewegt per `git mv` — kein Status-Feld.

**Welle:** [welle-82-config-flaechen](welle-82-config-flaechen.md) (zugeordnet
bei der Eröffnung).

**Bezug:**
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
(Image-Akzeptanzkriterien), [`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus);
die drei Erweiterungen aus slice-122 bis slice-124; die Release-Prozedur in
`docs/user/releasing.md`.

**Berührte Spec-Stellen:** — (Nutzer-Doku und Release-Register; die
Spec-Änderungen liegen in den drei Vorgänger-Slices).

**Verantwortlich:** pt9912. **Autor:** pt9912. **Datum:** 2026-08-22.

---

## 1. Ziel

Die drei Erweiterungen erreichen den Konsumenten: die Nutzer-Doku ist auf
Stand, das Release `v0.63.0` liegt auf GHCR, und der Digest-Pin ist
zurückgeschrieben. Release-Prep ist hier **eigener** Schritt und nicht Anhang
eines Feature-Slice — die drei Feature-Commits fassen die Nutzer-Doku bewusst
nicht an.

## 2. Vorgehen

1. **Doku-Currency über alle drei Erweiterungen**, nach der §4-Checkliste der
   Release-Prozedur und den bekannten Blind-Spots: Benutzerhandbuch (§5
   Konfigurations-Referenz, §6 Modul-Tabelle, §11-Zeile **chronologisch
   unten**, Kopf-Version), `docs/user/operations.md` (Modul-Liste **und**
   Optionen-Tabelle), README **und** README.de (Modul-Listen), `version.md`
   (Verlaufs-Zeile **und** der Anker wandert auf die neue Version — genau
   dafür trägt ihn nur die aktuelle), `CHANGELOG.md`.
2. **Den Handbuch-Satz aus slice-121 korrigieren:** dort steht, `diagrams` habe
   weder Datei- noch Zeilen-Ventil. Mit slice-124 ist das überholt.
3. **Fenced Config-Beispiele prüfen** — sie sind gate-unsichtbar und können dem
   Validator widersprechen; jedes neue Beispiel gegen den Validator laufen
   lassen.
4. **`make fullbuild`**, dann Tag `v0.63.0`, Push, Release-Pipeline beobachten,
   `docker pull`-Beweis (OCI-Label, Digest, Smoke).
5. **Digest-Backfill** committen (Handbuch, Roadmap, Slice-Verweise auf den
   Digest-Pin).
6. Unabhängiger Review **vor** dem Tag; Closure.

## 3. Ausdrücklich NICHT in diesem Slice

- **Keine Funktions-Änderung.** Was hier passiert, ist Doku und Auslieferung.
- **Keine Profil-Umstellung** (das `diagrams`-Scoping, die Beobachtungs-3×-Form)
  — eigene Entscheide nach dem Release.

## 4. Definition of Done

- [x] Alle Doku-Flächen der §4-Checkliste nachgezogen; die überholte
      Ventil-Aussage korrigiert — **weiter reichend als geplant**: nicht nur im
      Handbuch, sondern auch in `operations.md`, beiden READMEs und, nach dem
      Review, in Lastenheft (0.65.1) und Spezifikation. Config-Beispiele gegen
      den Validator geprüft: der §5-Block wie gedruckt (Exit 1) und einmal mit
      einkommentierter `patterns`-Liste statt Kurzform (Exit 1) — kein Exit 2;
      der Reviewer hat den `structure`-Block als N+1-Form nachgefahren.
- [x] `make fullbuild` Exit 0 (48 Anforderungen, 0 Waisen; Closure-Profil 414
      Dateien, 0 Befunde), zweimal gefahren — vor und nach der
      Review-Einarbeitung. Unabhängiger Review **vor** dem Tag: Report
      [2026-08-23](../../../reviews/2026-08-23-slice-125-release-v0630-review.md),
      Verdikt blockierend, 0 HIGH · 2 MEDIUM · 4 LOW · 2 INFO, alle acht
      eingearbeitet.
- [x] Tag `v0.63.0` gepusht, Release-Pipeline `32619522094` completed success,
      `docker pull`-Beweis: Digest
      `sha256:7049cefd2d91b367b72f15f789123ab5f51bf09150ad9cd262a9a945cfceb16e`,
      OCI-Label `org.opencontainers.image.version` = `0.63.0`, netzloser Smoke
      über den Digest-Pin gegen dieses Repo Exit 0 (448 Dateien, 0 Befunde);
      Digest-Backfill committet.
- [x] [`version.md`](../../../../version.md): Verlaufs-Zeile ergänzt und der
      Anker auf die neue Version gewandert — genau einer, vom Reviewer
      nachgezählt.

## 5. Abnahme-Punkte / Risiken

- **Die Doku-Blind-Spots sind gate-unsichtbar:** Modul-Listen, §11-Position,
  Optionen-Tabelle und fenced Config-Beispiele fängt kein Gate. Sie stehen
  deshalb auf der Checkliste und werden einzeln abgehakt. — **Ausgang:**
  *eingetreten, und die Checkliste hat nicht gereicht.* Die Enumerationen
  hielten (kein neues Modul, keine neue Option; `operations.md` führt weiterhin
  20 Module), die §11-Zeile steht chronologisch unten, und beide
  Config-Beispiele bestehen den Validator. **Aber die Checkliste kennt keinen
  Punkt „Prosa-Aussage über eine Modul-Menge"** — und genau dort lagen beide
  MEDIUM-Befunde. Ein Blind-Spot, der auf keiner Liste steht, wird von der
  Liste nicht gefunden.
- **Der Anker-Wechsel in `version.md`** ist neu scharf (seit slice-121): wird
  er vergessen, bleibt ein alter Pin auflösbar. — **Ausgang:** *nicht
  eingetreten.* Der Anker ist gewandert, die `v0.62.0`-Zeile hat ihn verloren,
  und der Reviewer hat unabhängig nachgezählt: genau einer.
- **Drei Erweiterungen in einem Release** heißt drei Handbuch-Abschnitte und
  eine §11-Zeile, die alle drei nennt — nicht drei Zeilen. — **Ausgang:**
  *nicht eingetreten.* Eine §11-Zeile (1.56) trägt alle drei Erweiterungen und
  zusätzlich die Korrektur der Vorgängerzeile 1.55, die sie widerlegt.

## 6. Trigger

**Start** (`open` → `in-progress`): slice-122, slice-123 und slice-124 in
`done/`.

**Rückführungen:** `in-progress` → `next`, falls der Review vor dem Tag eine
Funktions-Auflage findet (dann erst der betroffene Feature-Slice).

## 7. Vorgelagert (vor der Modus-Begründung)

- **Sub-Area prüfen:** Nutzer-Doku (GF), Release-Register (GF),
  Release-Pipeline (GF).
- **Offene Beobachtungen sichten** (Register-Stand 2026-08-22): **BEO-009**
  (Botschaft behauptet ungeprüfte Probe) ist im Release-Flow einschlägig —
  Botschaften über Pipeline-Läufe erst **nach** deren Exit; BEO-006 beim
  pfad-selektiven Commit; BEO-002 als Spiegel-Pflicht.

Slice-ID: slice-125. Betroffene IDs:
[`DC-FA-DIST-001`](../../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[`DC-QA-02`](../../../../spec/lastenheft.md#dc-qa-02--determinismus). Module:
Nutzer-Doku, Release-Register, Release-Pipeline. Gates: `make fullbuild`,
`make image-test`; Release-Pipeline.

## 8. Sub-Area-Modus-Begründung

**GF (Greenfield, Repo-Default)** — Auslieferung nach etablierter Prozedur.

## 9. Closure-Notiz (nach `done/`)

Ausgeliefert ist `v0.63.0` mit den drei Konfigurations-Flächen der Welle;
Digest, OCI-Label und ein netzloser Smoke über den Digest-Pin liegen als Beleg
im Backfill-Commit. Der geplante Teil des Slice — Doku-Currency, Tag,
Pipeline, Backfill — hat gehalten.

**Die Lehre ist die vierte Reichweiten-Lehre dieser Welle, und sie sitzt eine
Ebene höher als die drei davor.** slice-122, slice-123 und slice-124 haben
jeweils eine Exklusivitäts-**Aussage** aus dem Anlass statt aus dem Bestand
gezogen. Dieser Slice hat genau solche Aussagen korrigiert — und dabei eine
neue **Herleitung** erfunden: „der Marker wirkt in den vier Modulen, die eigene
Muster konfigurieren und ihre Befunde an Zeilen hängen". Die Aufzählung stimmt,
das Kriterium nicht: `matrix` konfiguriert `classes[].token` und meldet auf
Zeilen, `structure` konfiguriert vier Muster und meldet auf der
Überschriften-Zeile. In `operations.md` stand der Satz sogar unter der
Überschrift „Kein Zeilen-Marker für `matrix`". **Wer eine falsche
Exklusivität repariert, ist versucht, die reparierte Menge zu begründen — und
eine Begründung ist eine neue, ungezählte Behauptung.** Die Menge steht
deshalb jetzt überall ausdrücklich als *benannte Liste, kein ableitbares
Kriterium*.

**Die zweite Lehre ist eine Reichweite anderer Art: die Spiegel-Liste endete
an der Dokument-Rolle.** Korrigiert waren Handbuch, `operations.md` und beide
READMEs — also alles, was als „Nutzer-Doku" im Slice-Plan stand. Lastenheft
und Spezifikation trugen denselben falschen Satz weiter, das Lastenheft seit
0.8.0. Wer der Source Precedence folgt und im ranghöheren Dokument nachschlägt,
hätte genau die Aussage gelesen, die dieser Slice als Fehler ausweist.
[`MR-025`](../../../../harness/conventions.md#mr-025) verlangt, die Spiegel
**vor** dem Editieren aufzulisten; die Liste war nach
Dokument-Rolle gebildet statt nach Aussage. Der Ableiter ist das `grep` nach
dem alten Wortlaut — repo-weit, nicht nach Verzeichnis.

**Ein Beispiel, das seinen eigenen Gegenstand verfehlte.** Der neue
Handbuch-Absatz zum `diagrams`-Marker zeigte ihn in einem Tilden-Fence
**innerhalb** eines `markdown`-Fence. Der geteilte Fence-Automat schaltet bei
jeder Fence-Zeile um: die Tilden-Zeile schloss den äußeren Block, die
Diagramm-Zeilen waren für alle Prosa-Module Fließtext, und grün war das
Beispiel allein wegen seiner zweistelligen Kennungen. Gefangen hat es der
`doc-check` erst im zweiten Anlauf, benannt der Review. Der Absatz zeigt jetzt
einen echten Fence, nennt beim Zeilen-Ort das Kommentar-Zeichen der
Diagramm-Sprache und weist die schließende Fence-Zeile ausdrücklich als
Nicht-Ort aus — alle drei Zusagen mit einer Fünf-Fälle-Fixture gegengeprobt.

**Eine Zahl in einer Botschaft benannte die falsche Menge.** „Alle 24
bare-Tag-Pins" — 24 ist exakt die Zahl der `ghcr`-**präfixierten**, also
gate-gedeckten Pins; „bare Tag" heißt in der Release-Prozedur *ohne*
`ghcr`-Präfix, und davon trägt das Handbuch genau **einen**. Die einzige
Stelle, für die Handarbeit überhaupt nötig war, verschwand in einer Zahl, die
sie nicht enthält. Der Commit lag bereits auf `origin` und wurde nicht
amendiert; die Korrektur steht im Review-Commit und hier.

**Offen und bewusst nicht angefasst:** der `diagrams.scope` des eigenen Profils
und die 3×-Form von BEO-008 bleiben eigene Entscheide mit eigener Messung —
diese Welle macht sie baubar, sie baut sie nicht. Die Welle selbst bleibt
offen: ihr Closure-Trigger verlangt zusätzlich die Ergebnisnotiz und den
Status-Übergang der begleitenden ADR.
