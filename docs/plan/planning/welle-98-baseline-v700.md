# Welle welle-98: Baseline-Pin auf `v7.0.0` heben und die Neuerungen dem Auftraggeber vorlegen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-11-25.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der vendorte Baseline-Baum steht auf
[`v7.0.0`](https://github.com/pt9912/ai-harness-course/releases/tag/v7.0.0) (Kurs-Welle 168,
2026-11-23), jeder lebende pin-gebundene Verweis zeigt auf ihn, das Produkt ist unverändert —
und der Auftraggeber hat zu jeder Neuerung des Bundles eine Entscheidungsvorlage, ob und wie sie
übernommen wird. **Übernommen wird in dieser Welle nichts.** Das Akzeptanzkriterium, an dem sich das
Ziel spiegelt: *Given der Pin steht auf `v7.0.0`, when `make baseline-freshness` läuft, then meldet
es den Pin als aktuell (heute Exit 3, „NEUER RELEASE verfügbar (Pin v6.18.1)“, gemessen
2026-11-25), und `make gates` ist grün.*

**Der Zuschnitt ist die kleinste tragende Fassung je Frage** — die Lehre aus
[`BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass`](observations/BEO-ALL/empfehlung-weitet-anforderung-ueber-den-anlass/observation.md):
Der Bump ist mechanisch und gut bekannt (die Pin-Kette von [`MR-011`](../../../harness/conventions.md#mr-011)
bis [`MR-080`](../../../harness/conventions.md#mr-080); zuletzt
[`slice-277`](done/wellenlos/slice-277-baseline-v6181-bump.md)); die Sichtung ist eine
**Entscheidungsvorlage, kein Umbau**. Ein Release ist nicht Teil der Welle (§6).

**Entscheide des Auftraggebers (2026-11-25, bindend):**

1. **`spec/spezifikation.md` bleibt deutsch und wird nicht umbenannt.** Das Regelwerk `v7.0.0` nennt
   `specification.md` im Verzeichnis `spec/`; die Abweichung ist als Adaption festzuhalten. Die Form (eigener MR-Eintrag
   oder Teil des Bump-Eintrags) entscheidet der Architect (§5, Frage 1).
2. **Umfang „Bump plus Neuerungen“:** (a) der Bump selbst, (b) eine Sichtung der Neuerungen, an deren
   Ende der Auftraggeber je Neuerung entscheidet, was übernommen wird. Die Sichtung übernimmt nichts;
   Übernahmen sind eigene Folge-Slices nach seiner Entscheidung.

**Gemessener Ausgangsstand** (Planner, 2026-11-25, am entpackten Release-Asset `lab-regelwerk.zip`
von `v7.0.0` — SHA256 `34a7866d…e0e`, gleich dem `SHA256SUMS` des Release — gegen den vendorten
Baum `.harness/baseline/v6.18.1/`; der Implementer **wiederholt** die Messung am vendorten Baum
von `v7.0.0` und am alten Baum aus `git archive`, die Zahlen hier sind Erwartung, nicht Beleg).
Gegenstand der Zählung ist die **Datei im Bundle**, nicht der Verweis auf sie:

| Gemessen | Zahl | Form der Zählung |
|---|---|---|
| Bundle-Dateien alt / neu | 54 / 55 | `find` über `regelwerk/` und `templates/` ohne `SHA256SUMS`; vendort je +1 für `SHA256SUMS` (55 / 56) |
| gleichnamig in beiden | 47 | Rest: 7 umbenannt (alt) ↔ 7 (neu), 1 neu |
| davon byte-verschieden | 38 | `cmp` je Datei; **862** Diff-Zeilen (`diff`, Zeilen mit `<`/`>`) |
| davon **inhaltlich** verschieden | 31 (19 Regelwerk, 12 Vorlagen) | ohne die `<!-- Quelle:`-Zeilen, Versionsstempel auf einen Platzhalter normalisiert; **804** Diff-Zeilen. Die übrigen 7 ändern nur Stempel und Quell-Verweis |
| umbenannt | 7 | `welle.template.md`→`wave.template.md`, `welle-results.template.md`→`wave-results.template.md`, `archiv-stub-slice`/`-welle`→`archive-stub-slice`/`-wave`, `spezifikation.template.md`→`specification.template.md`, `NNNN-titel`→`NNNN-title`, `MR-NNN-titel`→`MR-NNN-title` (je 2–28 inhaltliche Diff-Zeilen) |
| neu | 1 | `templates/harness/sensors/environment-check.template.md` (74 Zeilen) |

**Größte Bewegungen** (inhaltliche Diff-Zeilen): `modul-14-docker-harness` 258,
`templates/.d-check.yml` 74, `modul-04-adrs` 57, `modul-05-planning-harness` 49,
`templates/README.md` 35. **Drei Befunde prägen den Zuschnitt:**

1. **Die sieben Vorlagen-Umbenennungen treffen das Repo mechanisch.** `plan-welle.md` und
   `close-welle.md` nennen die alten Namen (`git grep -n "welle\.template\|welle-results\.template\|archiv-stub"
   -- .claude`), ebenso `harness/conventions.md` (`MR-<NNN>-titel.template.md`) und zwei Kommentarzeilen in
   `tools/archive-wave/stub.go`. Das gehört zum Bump (slice-295).
2. **Das Regelwerk führt jetzt englische Dateinamen als Regel** (`grundlagen-harness-dateien.md`, neuer
   Absatz *Dateinamen sind englisch*): `specification.md`, aber auch `wave-<Kennung>.md`,
   `wave-<Kennung>-results.md`, `archive.zip` und der Herkunfts-Anker `seit wave-<Kennung>`.
   Entscheid 1 deckt davon **nur** die Spezifikation; `welle-`, `archiv.zip` und `seit welle-` stehen offen und
   sind **Sichtungs-Gegenstand** (slice-296), nicht Teil des Bumps.
3. **Eingebettete Dokumente:** Das Produkt bettet `spec/spezifikation.md` und das Handbuch ein
   ([`ADR-0107`](../adr/0107-mitgelieferte-dokumente-im-modul-root.md), `manual.go`).
   `spec/spezifikation.md` Zeile 8 nennt den Pin (`Baseline-Regelwerk v6.18.1`); das Handbuch nennt ihn nicht
   (`git grep -n "v6\.18\.1" -- docs/user README.md` ohne Treffer). Der Bump ändert damit die Ausgabe von
   `--manual` ([`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen))
   in genau dieser Pin-Nennung — das Produkt-**Verhalten** bleibt, die eingebetteten **Bytes** nicht (§3 (d)).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- **Ein neuer Baseline-Release liegt vor, und der gepinnte Stand ist unverändert.** `make baseline-freshness`
  endet mit Exit 3 und der Zeile „NEUER RELEASE verfügbar (Pin v6.18.1): v7.0.0“, zugleich „gepinnter Tag
  v6.18.1 upstream unverändert (Bytes == vendored SHA256SUMS)“ — gelesen am 2026-11-25. Der Nachtlauf
  `upstream-drift` ist seit demselben Tag rot (`make nightly-state`; die Meldung ist **planmäßig** gelesen, nicht
  weggeklickt: ein Fremd-Release, [`MR-051`](../../../harness/conventions.md#mr-051)). Es ist eine reine
  Fortschreibung, kein Nachholen stiller Drift.
- **Der Auftraggeber hat die Welle angefordert** (2026-11-25, in der Sitzung), mit den zwei Entscheiden aus §1.
- **welle-97 ist geschlossen:** ihre Welle-Datei liegt unter `done/welle-97/`
  (`ls docs/plan/planning/done/welle-97/`), ihre Closure-Notiz als `done/welle-97-results.md`.
- **`git fetch`:** `origin/main` und der lokale Hauptzweig sind gleich (`git rev-list --left-right --count
  origin/main...HEAD` ergibt `0 0`). `in-progress/` trägt nur die Roadmap, `next/` ist leer, `open/` trägt nur
  die zwei Slices dieser Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Beide Slices der Welle sind `done/`.
- `make fullbuild` grün auf dem Endstand.
- **Das Mehr gegenüber den Slice-DoDs.** (d) bis (f) sind Aussagen über das Repo als Ganzes, die kein
  einzelner Slice abhaken kann — sie stehen deshalb **nicht** in einer Slice-DoD
  ([`AGENTS.md`](../../../AGENTS.md) §5 Regel 10). (a) bis (c) liefert
  [`slice-295`](open/slice-295-baseline-v700-bump.md); die Welle **misst sie am Endstand gegen**, nachdem
  [`slice-296`](open/slice-296-baseline-v700-neuerungen-sichtung.md) gelesen und geschrieben hat — ein Nachzug
  dort darf kein Rot erzeugen. Zu jeder Aussage das Rot **ohne** die Arbeit der Welle:
  - **(a) Pin und Bestand gehören zusammen.** `make baseline-verify` Exit 0 **und** `make baseline-freshness`
    Exit 0 mit „Pin v7.0.0 ist der neueste Release-Tag“ und „upstream unverändert“. *Rot ohne die Welle:*
    `make baseline-freshness` endet heute mit Exit 3. Erscheint bis zur Closure ein neuerer Release, ist das
    ein neuer Anlass und nicht Teil dieser Welle; die Closure-Notiz nennt ihn.
  - **(b) Kein lebender Verweis nennt mehr den alten Pin.** Die Treffermenge von
    `git grep -l 'v6\.18\.1' -- . ':!docs/plan/planning/done' ':!docs/reviews' ':!CHANGELOG.md' ':!harness/conventions/done' ':!.harness/baseline'`
    ist **genau** die Menge, die der Bump-Eintrag als Vergangenheits-Aussage oder eingefrorene Nennung
    mit benanntem Grund führt (erwartet: der Bump-Eintrag selbst, die eingefrorenen `Accepted`-ADRs, falls der
    Architect das Ventil wählt, §5 Frage 3). **Und** die Symlinks lösen auf den neuen Baum auf
    (`for l in .claude/rules/*; do readlink "$l"; done | grep -c 'v6\.18\.1'` ergibt 0 — `git grep` liest keinen
    Symlink-Inhalt). *Rot ohne die Welle:* 55 Dateien mit 103 Zeilen (Stand `76cf7d51`, vor Anlage dieser Pläne)
    und 7 Symlinks. Die Zählung nennt ihre Form: gezählt wird die Zeichenkette in Dateien, nicht der Verweis.
  - **(c) Die gehaltenen Zitate stehen auf dem neuen Stand.** `make doc-check` ohne Befund, darunter `citations`
    für die **15** lebenden `d-check:cite`-Spannen auf den Baum (davon verschieben sich **7**, am Vergleich
    derselben Zeilenbereiche beider Bäume gemessen) und die Spannen der Pläne dieser Welle; `ignore-refs` in der
    [`.d-check.yml`](../../../.d-check.yml) wächst nur um die im Bump-Eintrag benannten Einträge
    ([`MR-069`](../../../harness/conventions.md#mr-069)). *Rot ohne die Welle:* nach dem Entfernen des alten Baums
    melden die Pfad-Links und die Spannen Befunde (Pfad-Links `target-missing`, Spannen das Modul `citations`); nach
    dem bloßen Tag-Tausch bleiben die 7 verschobenen Spannen rot.
  - **(d) Das Produkt ist unverändert — und die eine erwartete Abweichung ist benannt.** (d1) Kein
    Quelltext des Produkts ändert sich: `git diff --stat <Vor-Bump>..HEAD -- '*.go' ':!tools' Dockerfile go.mod`
    ist leer. (d2) `make blackbox-probe REF=<Vor-Bump>` — der Commit **vor** dem ersten Bump-Commit — endet mit
    „byte-identisch“ über seine Vergleiche (stdout, stderr, Exit getrennt; Formen: Standard, `--json`, `--yaml`,
    `--doctor`). **Nicht `REF=v0.93.0`:** auch dieses Vorher-Image kennt das Profil (`.d-check.yml` ist seit
    `v0.93.0` unverändert, `git diff --stat v0.93.0..HEAD -- .d-check.yml` ist leer), aber `v0.93.0` trägt gegenüber
    dem Hauptzweig zwei Änderungen an Bauteilen des Images (`git diff --stat v0.93.0..HEAD -- '*.go' Dockerfile go.mod manual.go spec docs/user`:
    `Dockerfile` und das Handbuch), die nicht aus dieser Welle stammen. (d3) Die eingebetteten Dokumente:
    `git diff <Vor-Bump>..HEAD -- spec/spezifikation.md docs/user/benutzerhandbuch.md` zeigt im Handbuch nichts
    und in der Spezifikation **nur die Pin-Nennung** (`v6.18.1` → `v7.0.0`); `make test` hält die
    Einbettung gleich der Quelle (`TestManual_EingebetteteDokumenteGleichDerQuelle`). **Die Grenze:** der
    Probe-Lauf fährt `--manual` nicht (Default-Formen), die Abweichung der `--manual`-Ausgabe belegt (d3), nicht
    (d2). *Rot ohne die Welle:* (d) ist **nicht** rot, wenn nichts geschieht — es hält eine Zusage fest, die der
    Bump verletzen könnte (ein mitgeänderter Quelltext, ein mitgeänderter Handbuch-Satz).
  - **(e) Das Trigger-Audit** (seit [`slice-241`](done/wellenlos/slice-241-trigger-audit-adoption.md)) ist über
    die Closure vollzogen und in `welle-98-results.md` belegt. Fällig ist mindestens der Auflösungs-Trigger von
    [`MR-080`](../../../harness/conventions.md#mr-080) („die nächste Pin-Hebung“); die übrigen Trigger-Klassen
    liest die Closure, sie sind hier nicht vorweggenommen.
  - **(f) Die Entscheidungsvorlage liegt dem Auftraggeber vor.** Jede der **39** Dateien mit Bewegung (31 inhaltlich
    verschieden, 7 umbenannt, 1 neu) hat eine Zeile mit Klassifikation und Empfehlung; die 7 Dateien, die nur ihren
    Stempel ändern, stehen in **einer** Sammelzeile mit ihrem Beleg. Gezählt wird mit dem Kommando aus
    [`slice-296`](open/slice-296-baseline-v700-neuerungen-sichtung.md) §3. **Die Entscheidungen des Auftraggebers
    sind nicht Closure-Bedingung** — die Welle schließt, wenn die Vorlage ihm vorliegt; jede Übernahme ist ein
    eigener Slice nach seiner Antwort (§6). *Rot ohne die Welle:* die Datei existiert nicht.
- Closure-Notiz in `welle-98-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-295`](open/slice-295-baseline-v700-bump.md) | **Bump:** Baseline-Pin auf `v7.0.0` — Re-Vendor, Spiegel, `d-check:cite`-Spannen, Vorlagen-Namen in den Command-Dateien, Bump-Eintrag (Delta vollständig nach Datei, keine Neuerung übernommen) und Adaptions-Eintrag zum Dateinamen der Spezifikation | [`MR-011`](../../../harness/conventions.md#mr-011) (Pin-Kette), [`MR-021`](../../../harness/conventions.md#mr-021) |
| [`slice-296`](open/slice-296-baseline-v700-neuerungen-sichtung.md) | **Sichtung:** je Bewegung der 39 Dateien eine Klassifikation (bereits befolgt · neu und anwendbar · nicht anwendbar) mit Empfehlung, dazu die Frage zu `environment-check` und zur „Umgebungs-Zusage“ gegen die vorhandenen Gates — **eine** Entscheidungsvorlage, keine Übernahme | [`MR-011`](../../../harness/conventions.md#mr-011) |

**Zwei Slices statt dreien.** Der dritte Kandidat wäre der Adaptions-Eintrag zum Dateinamen: ein Slice mit
einer Datei als Liefer-Punkt und ohne eigenen Lieferwert — er hat nur zusammen mit dem Bump einen Sinn. Der
Eintrag ist eine **Entscheidung des Architect** (Rollenfolge: Verdikt vor `in-progress`, den Eintrag schreibt der
Architect); seine **Form** steht im Verdikt vor dem Start, **geschrieben** wird er im Lauf von
[`slice-295`](open/slice-295-baseline-v700-bump.md) nach dem Re-Vendor-Commit — vorher löste sein Verweis auf den
Baum von `v7.0.0` nicht auf. Die Namens-Nachzüge (`.claude/commands/`, `harness/conventions.md`, zwei
Kommentarzeilen) haben denselben Gegenstand wie der Bump — den neuen Baum — und bleiben darin:
[`slice-295`](open/slice-295-baseline-v700-bump.md) hat drei Liefer-Punkte (Vendoring · Spiegel samt Zitaten und
Namen · die zwei Einträge), die Grenze von `modul-05-planning-harness.md` §Ziel-Form: Slice.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine Welle. Jede Übernahme einer Neuerung ist ein eigener Slice **nach** der Antwort des
  Auftraggebers auf die Vorlage (§6) und hängt an dieser Welle nicht.
- Wird blockiert von: keiner.
- **Die Reihenfolge ist fachlich bindend: erst pinnen, dann dem Kanon folgen.** Erst
  [`slice-295`](open/slice-295-baseline-v700-bump.md) (der neue Baum steht), dann
  [`slice-296`](open/slice-296-baseline-v700-neuerungen-sichtung.md): Die Sichtung liest den **vendorten** Baum von
  `v7.0.0`, nicht das Scratchpad, und den alten Baum aus `git archive <Vor-Bump>`; ein Audit gegen zwei Stände,
  von denen einer nirgends im Repo liegt, ist kein Audit (dieselbe Begründung wie in der Vorgänger-Migration
  [`welle-85`](done/welle-85/welle-85-baseline-v5120-migration.md) §5). Das WIP-Limit 1 hält sie ohnehin
  nacheinander. **Rollenfolge:** vor `in-progress` steht das Architect-Verdikt zu den Fragen unten, je Slice
  ein frischer Kontext für Implementer, Reviewer und Verifier.
- **Offene Fragen an den Architect — Antwort vor dem Start von slice-295.** Je Frage steht die kleinste Antwort,
  die den Anlass trägt; die größere steht daneben, damit ihr Mehraufwand vor der Bestätigung sichtbar ist.
  1. *Form des Adaptions-Eintrags zum Dateinamen: eigener MR oder Teil des Bump-Eintrags?* — **kleinste Antwort:
     ein eigener Eintrag.** Der Bump-Eintrag löst sich bei der nächsten Pin-Hebung auf
     ([`MR-080`](../../../harness/conventions.md#mr-080): „die nächste Pin-Hebung“ und `conventions/done/`), die
     Abweichung dagegen gilt über Hebungen hinweg, bis der Auftraggeber sie zurücknimmt oder die Baseline den
     Namen nicht mehr vorgibt — zwei Lebensdauern, zwei Dateien. Preis der größeren Antwort (ein Eintrag): die
     Abweichung wanderte beim nächsten Bump mit nach `done/` und müsste dort neu angelegt werden.
  2. *Wie weit trägt die Adaption?* — **kleinste Antwort: nur `spec/spezifikation.md`** (Entscheid 1). `welle-`,
     `archiv.zip`, `…-results.md` und `seit welle-` bleiben bis zur Antwort des Auftraggebers wie bisher und stehen
     im Bump-Eintrag als *nicht übernommene Bewegung*, nicht als Adaption; die Vorlage führt sie als Punkt 1
     ([`slice-296`](open/slice-296-baseline-v700-neuerungen-sichtung.md)). Die größere Antwort („Dateinamen
     bleiben deutsch“ als Prinzip) trüge eine Entscheidung, die der Auftraggeber nicht getroffen hat. Frage an
     den Architect, ob das Repo in der Zwischenzeit damit eine **unbenannte** Abweichung vom Kanon trägt — und wenn
     ja, ob sie als befristeter Eintrag zu führen ist. (Der nächstliegende Kanon-Satz, der Ausgang *widerspricht* in
     `modul-02-harness-bootstrap.md` §Freshness-Audit, gilt für **bestehende** Adaptionen; hier gibt es keine, deshalb
     eine Frage und kein Befund.)
  3. *Eingefrorene Dateien, die den alten Baum **verlinken** — Ventil oder Nachzug?* Gemessen **7**
     (`git grep -l '\](.*\.harness/baseline/v6\.18\.1' -- docs/plan/adr docs/plan/planning/done docs/reviews
     harness/conventions/done CHANGELOG.md docs/plan/cr`): die `Accepted`-ADRs 0112 bis 0116 (je eine Zeile), die
     Slice-Datei `slice-274` und der Architect-Verdikt `2026-10-30-welle-93-architect-verdikt`. **Kleinste
     Antwort: das Ventil** — je Datei ein `ignore-refs`-Eintrag, die Praxis seit dem Tombstone für `v6.9.0`
     ([`MR-069`](../../../harness/conventions.md#mr-069): „wächst mit jedem Bump“; heute 30 Einträge über 12
     Tags, gezählt mit `grep -c 'refs: \[".harness/baseline' .d-check.yml`). Die Alternative: der reine Pfad-Nachzug, den [`ADR-0103`](../adr/0103-adr-gate-laesst-pfad-nachzug-durch.md)
     für ADRs durchlässt, und für `slice-274` die Ausnahme, die [`MR-080`](../../../harness/conventions.md#mr-080)
     schon einmal genommen hat; er hält das Ventil klein und ändert zugleich, **welchen Stand** der Link behauptet.
  4. *Der alte Baum fällt im Bump, nicht nach dem Review.* Der Kanon sagt: „Das alte Verzeichnis fällt
     erst, wenn der Review durch ist.“ (`modul-02-harness-bootstrap.md` §Freshness-Audit, Absatz „Der Review
     vergleicht auch die Form, nicht nur die Regeln“); die Praxis
     dieses Repos seit [`MR-021`](../../../harness/conventions.md#mr-021) (1) und [`welle-85`](done/welle-85/welle-85-baseline-v5120-migration.md)
     entfernt ihn in der Bump-Etappe, und der alte Baum bleibt über die git-Historie lesbar
     ([`MR-052`](../../../harness/conventions.md#mr-052)). **Kleinste Antwort: die Praxis bleibt**, der Bump-Eintrag
     nennt den Vor-Bump-Commit, aus dem die Sichtung vergleicht. Frage: bestätigen, und ob der Satz im Eintrag
     als Abweichung vom Kanon steht.
  5. *Gehört der Durchgang durch die Adaptions-Liste (Freshness-Audit) in diese Welle?* — **kleinste Antwort: nein.**
     Gelesen: [`MR-060`](../../../harness/conventions/done/MR-060-baseline-v600.md) (das letzte Major),
     [`MR-074`](../../../harness/conventions/done/MR-074-baseline-v6170.md) und
     [`MR-080`](../../../harness/conventions.md#mr-080) beschreiben keinen Durchgang durch die Adaptions-Liste
     (`grep -n Freshness` trifft dort höchstens eine Aussage über das Altern des Baums); als eigene Etappe stand er in
     [`welle-85`](done/welle-85/welle-85-baseline-v5120-migration.md) — die dazwischenliegenden Bump-Slices sind
     nicht gelesen. Gemessen als Dateinamen-Join von
     `Geltungsbereich`/`Ersetzt-Baseline-Regel` der 41 aktiven Einträge mit den 19 inhaltlich geänderten
     Regelwerk-Dateien: **höchstens 17 Treffer** (Obergrenze, ungeprüft — `README` trifft auch `harness/README.md`).
     Die größere Antwort ist ein dritter Slice mit je einer Frage pro Treffer; sie ist die Entscheidung, ob ein
     Major ihn verlangt.
  6. *Wo liegt die Entscheidungsvorlage?* **Eine Konvention für diese Klasse gibt es nicht** (`docs/plan/` trägt
     `adr/`, `carveouts/`, `cr/`, `planning/`). Die Vorgänger-Delta-Audits lebten in §9 ihres Slice
     ([`welle-85`](done/welle-85/welle-85-baseline-v5120-migration.md): `slice-149`) und liegen nach der Closure im
     Archiv — der Auftraggeber antwortet aber **nach** der Closure. **Kleinste Antwort:**
     `baseline-v700-neuerungen.md` unter `docs/plan/planning/`, flach, mit einem Namen, den der Glob `welle-*.md` des
     Moduls `planning` nicht trifft; sie überlebt die Closure und wird bei der Antwort fortgeschrieben.
     Die Frage an den Architect: Ort bestätigen, und ob die Klasse einen Eintrag braucht.
  7. *Die zwei Kommentarzeilen in `tools/archive-wave/stub.go`* nennen die alten Vorlagen-Namen. **Kleinste
     Antwort: im Bump nachziehen** — ein Kommentar, der eine Datei nennt, die es nicht mehr gibt, ist falsch
     ([`AGENTS.md`](../../../AGENTS.md) §3.7), und `make archive-wave-test` hält den Rest. Kein Verhalten ändert sich.
- **Review-Vorgabe der Welle** (Lehre aus welle-96 und welle-97): je Slice **höchstens zwei Runden**,
  eingearbeitet werden nur HIGH und MEDIUM; LOW und INFO bleiben im Report. Jede Nachbesserung zitiert die
  Stelle, die sie trägt, statt sie nachzuerzählen
  ([`AGENTS.md`](../../../AGENTS.md) §5 Regel 18).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Die Umbenennung der Spezifikation.** `spec/spezifikation.md` bleibt (Entscheid 1); daran hängen
  `manual.go`, die Pfade in der [`.d-check.yml`](../../../.d-check.yml) und die Spec-Verweise im ganzen Repo. Der Bump
  schreibt die Abweichung auf, er ändert keinen dieser Pfade.
- **Die Übernahme jeder Neuerung** (Entscheid 2) — auch der kleinen und der mechanischen: die Bezeichner `wave-`,
  `archive.zip`, `seit wave-` in unseren Artefakten, das neue Vorlagen-Paar, die Worktree-Regel samt
  `.claude/worktrees/**` im Scan-Ignore, die Bedrohungs-Ausgänge, `environment-check` und die Umgebungs-Zusage.
  Die Sichtung klassifiziert und empfiehlt; ein Haken „übernommen“ gehört in einen Folge-Slice nach der Antwort.
- **Eine Änderung am Produkt** — Quelltext, Standardwerte, Spezifikation, Handbuch. Die eine Berührung ist die
  Pin-Nennung in `spec/spezifikation.md` Zeile 8 (§1, Befund 3); sie ist in (d) benannt.
- **Ein Durchgang durch die Adaptions-Liste** (Freshness-Audit) — Frage 5 in §5; ohne Antwort bleibt es bei der
  Obergrenze von 17 Treffern als Messung.
- **Ein CR an den Kurs.** Findet die Sichtung einen Widerspruch im Regelwerk selbst, ist ein ausgehender
  Change Request ([`MR-035`](../../../harness/conventions.md#mr-035)) ein eigener Vorgang.
- **Die Tool-Pins** (Go, golangci-lint, semgrep, a-check) und alle Dependabot-Hebungen — eigene Vorgänge mit
  eigener Messung, wie in [`slice-277`](done/wellenlos/slice-277-baseline-v6181-bump.md).
- **Nacharbeit am eingefrorenen Bestand** — `done/`, `docs/reviews/`, `CHANGELOG.md`, `conventions/done/` und der
  Kern der `Accepted`-ADRs nennen den alten Pin als Vergangenheits-Aussage und bleiben byte-stabil; Ausnahmen
  entscheidet allein die Antwort auf Frage 3 in §5.
- **Release, Tag, Release-Prep** (CHANGELOG, READMEs, Handbuch-Kopf, Digest-Pins) — eigener Vorgang nach der
  Welle ([`AGENTS.md`](../../../AGENTS.md) §5 Zeile 17); der Bump berührt keine Distributions-Fläche außer der
  eingebetteten Pin-Nennung in (d).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-98-results.md` (Geschwister im Ruheort `done/`, entsteht bei der Closure)
Zähler: [Beobachtungs-Register](observations/README.md)
