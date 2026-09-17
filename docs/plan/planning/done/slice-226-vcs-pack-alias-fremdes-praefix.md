# Slice slice-226: `vcs` erkennt Packs unter fremdem Namens-Präfix

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`DC-FA-VCS-001`](../../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[ADR-0086](../../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md).

**Berührte Spec-Stellen:** [`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 (die Kandidaten-Auflösung — Nachtrag zur Pack-Erkennung, kein
Lastenheft-Bump).

**Verantwortlich:** pt9912 (Implementer-Rolle).

**Autor:** pt9912. **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

**Ziel.** `vcs` (und jedes Modul, das denselben Objektspeicher über den
VCS-Port liest) erkennt einen Pack, dessen Objekte gültig sind, auch dann,
wenn seine Datei nicht unter dem kanonischen Präfix `pack-` liegt — etwa
`loose-<hash>.pack`, wie `git maintenance run --task=loose-objects` es
schreibt. Ein wirklich fehlendes oder beschädigtes Objekt bricht weiterhin
mit Exit 2 ab.

**Anlass ist ein eingehender CR** des Adopters `ai-harness-init`
(`docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-pack-praefix.md`):
unter `v0.76.1` bricht `vcs` an einem Repo ab, dessen Objektspeicher Packs
unter dem Präfix `loose-` trägt, obwohl `git cat-file`/`git fsck` dieselben
Objekte anstandslos lesen. Ursache, empirisch am eigenen Probe-Repo
nachvollzogen: go-gits `storage/filesystem/dotgit`-Schicht (`v5.19.2`)
entdeckt und öffnet Packs **ausschließlich** über den rekonstruierten Namen
`pack-<hash>.{pack,idx}` — ein anderer Präfix ist für sie unsichtbar,
unabhängig vom Objekt-Inhalt.

**Warum Recherche vor einem geschriebenen Plan stand — benannt, nicht
verschwiegen.** Der Entwurf (ein read-only `billy.Filesystem`-Dekorator, der
Packs unter ihrem kanonischen Namen zusätzlich sichtbar macht) war ohne
Kenntnis von go-gits internem `DotGit.objectPacks`/`objectPackPath`-Mechanismus
nicht seriös zu planen — „lies einfach jede Pack-Datei" wäre eine Behauptung
ohne Deckung gewesen. Die Recherche (go-gits Quelltext für
`storage/filesystem/dotgit`, `storage/filesystem/object.go`, `repository.go`
sowie go-billys `Filesystem`-Interface) **war** die Plan-Arbeit dieses Slice;
der Code entstand im selben Zug, bevor diese Datei geschrieben wurde. Das ist
eine Abweichung von der kanonischen Reihenfolge „Plan vor Code"
(`AGENTS.md` §6 Schritt 4) — §7 trägt sie als Lerneintrag.

**Abgrenzung — drei Punkte, jeder mit Grund:**

1. **Keine Änderung an go-git selbst.** Ein Fork oder Patch der Bibliothek
   wäre eine neue Supply-Chain-Last, die dieses Repo nirgends sonst trägt
   ([ADR-0002](../../adr/0002-distribution-ghcr-image.md)/[ADR-0024](../../adr/0024-vcs-immutable-gate.md))
   — *es wäre ein anderer Vorgang*, siehe [ADR-0086](../../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md)
   Option C.
2. **Kein Repack/Umbenennen im gemounteten Repository.** Der Adapter bleibt
   rein lesend ([`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)); der Dekorator legt keine Datei an und benennt
   keine um — *Bestand bleibt bewusst stehen*, dieselbe Abgrenzung wie in
   [slice-220](../done/slice-220-vcs-pfadmenge-statt-diff.md).
3. **Keine Lockerung des fail-closed-Abbruchs für einen wirklich
   unauflösbaren Bestand.** Ein Pack ohne gültiges Hash-Suffix oder ohne
   passende `.idx`-Datei bleibt unsichtbar und der Lauf bricht ab, wie vom CR
   ausdrücklich verlangt („Was ausdrücklich nicht gebeten ist") —
   *Schicht-Abgrenzung*: dieser Slice erweitert die Erkennung, er lockert
   nicht die Bedingung, unter der sie abbricht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(1)** Ein Pack mit gültigem Hash-Suffix und passender `.idx`-Datei
      unter einem anderen Präfix als `pack-` wird von `AllPaths`/`FileAt`
      korrekt gelesen — gemessen an einem Repo, dessen Objekte **nur** noch
      im umbenannten Pack liegen (lose Kopien entfernt), nicht nur an einem
      Repo mit redundanter loser Kopie.
- [x] **(2)** Ein Pack **ohne** gültiges Hash-Suffix (Garbage-Name) bleibt
      unsichtbar und der Lauf bricht weiterhin mit Exit 2 ab, wenn er die
      einzige Quelle eines Objekts ist — die Gegenprobe aus dem CR.
- [x] **(3)** Die Architektur-Grenze ist explizit entschieden, nicht
      stillschweigend verletzt: `make arch-check` bleibt grün, die
      `os`/`io/fs`-Kapsel-Erweiterung trägt eine eigene ADR
      ([ADR-0086](../../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md)).
- [x] `make gates` grün.
- [x] Unabhängiger Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
- [x] Der eingehende CR trägt eine `## Antwort`-Sektion, `Stand: entschieden
      und umgesetzt`.
- [x] `spec/spezifikation.md` §[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
      Schritt 2 nennt den Alias-Mechanismus.

## 3. Plan (vor Code)

**Abweichung von der Reihenfolge, siehe §1.** Die folgenden Schritte
beschreiben, was geschah, nicht was noch bevorsteht — mit Ausnahme von
Schritt 4.

1. go-gits `storage/filesystem/dotgit`/`object.go` und go-billys
   `Filesystem`-Interface lesen: Packs werden ausschließlich über
   `ReadDir("objects/pack")` entdeckt und über den aus dem Dateinamen
   rekonstruierten Pfad `objects/pack/pack-<hash>.{pack,idx}` geöffnet.
2. `packAliasFS` als schlanker, lesender `billy.Filesystem`-Dekorator bauen:
   `ReadDir` spiegelt jeden aliasfähigen Pack unter seinem kanonischen Namen,
   `Open` löst einen kanonischen Pfad auf die reale Datei auf. Jede andere
   Methode geht unverändert durch (eingebettetes Interface).
3. `Adapter.Open` von `gogit.PlainOpen` auf den manuellen Aufbau
   (`osfs` → `packAliasFS` → `filesystem.NewStorage` → `gogit.Open`)
   umstellen, damit die Objekt-Storage den Dekorator trägt.
4. Verifizieren: eigenes Probe-Repo, Objekte real gepackt
   (`Repository.RepackObjects`, keine Shell-Abhängigkeit), lose Kopien
   entfernt, Pack umbenannt — Fix greift. Gegenprobe mit Garbage-Namen bleibt
   fehlerhaft. `make gates`. CR-Antwort und Spec-Nachtrag schreiben.

## 4. Trigger

**Beanspruchung:** der Auftraggeber gibt die Umsetzung des CR frei (bereits
erteilt, 2026-09-17).

**Rückführung nach `open/`** (`in-progress→open`): falls sich beim Review
zeigt, dass die Kapsel-Erweiterung ([ADR-0086](../../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md))
eine andere, hier nicht bedachte Fehlerklasse öffnet.

## 5. Closure-Trigger

DoD (1) bis (3) sowie `make gates`, Review, Closure-Notiz, Register,
Risiko-Ausgänge, Paarungen, CR-Antwort und Spec-Nachtrag — alle abgehakt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Alias-Mechanismus ist an go-gits v5.19.2-internem Verhalten
  geeicht.** `packAliasFS` verlässt sich darauf, dass `DotGit` Packs
  ausschließlich über `ReadDir`+rekonstruierten Pfad anspricht (keine
  öffentliche API, kein SemVer-Vertrag). Ein Bump von go-git könnte diesen
  Mechanismus intern ändern, ohne die öffentliche API zu brechen — dann liefe
  der Dekorator ins Leere, ohne dass ein Compile-Fehler warnt. — **Ausgang:**
  weiter offen →
  [`BEO-ALL/undokumentierte-bibliotheks-interna-als-vertrag`](../observations/BEO-ALL/undokumentierte-bibliotheks-interna-als-vertrag/observation.md).
- **Der Hash-Suffix-Regex ist eine Heuristik, kein Beleg über den
  Pack-Inhalt.** Ein Dateiname, der zufällig mit 40/64 Hex-Zeichen endet, aber
  nicht der tatsächliche Content-Hash des Packs ist, würde fälschlich als
  Alias-Ziel akzeptiert. `git maintenance` benennt Packs korrekt (der Hash-Teil
  ist der echte Content-Hash), ein von Hand falsch benanntes Pack wäre ein
  Angriffs- oder Fehlerfall außerhalb des CR-Anlasses. Der unabhängige Review
  (R1-F-2) fand dazu einen konkreten, seither behobenen
  Determinismus-Mangel bei zwei kollidierenden Namen — die Heuristik selbst
  bleibt aber bestehen. — **Ausgang:** weiter offen →
  [`BEO-ALL/dateiname-als-inhalts-beleg-akzeptiert`](../observations/BEO-ALL/dateiname-als-inhalts-beleg-akzeptiert/observation.md).
- **Nur ein Adopter-Fall belegt die Klasse.** Wie bei
  [ADR-0072](../../adr/0072-workflows-modul.md) gilt: eine Aussage über „jedes
  Präfix" ist an einem Exemplar (`loose-`) geeicht. — **Ausgang:** weiter offen →
  [`BEO-ALL/beleg-fuer-generalitaet-ist-ein-exemplar`](../observations/BEO-ALL/beleg-fuer-generalitaet-ist-ein-exemplar/observation.md).

## 7. Closure-Notiz

**Geliefert.** `vcs` (und jedes Modul, das denselben git-Port liest) löst
seither einen Pack mit gültigem SHA1/SHA256-Hash-Suffix und passender
`.idx`-Datei unabhängig vom Namens-Präfix auf — ein read-only
`billy.Filesystem`-Dekorator (`packAliasFS`,
`internal/adapter/driven/git/packalias.go`), der go-gits kanonische
Pack-Namensauflösung nicht ändert, sondern ihr zusätzliche Sichtbarkeit
vorschaltet. Der eingehende CR trägt seine `## Antwort`, `Stand: entschieden
und umgesetzt`; `spec/spezifikation.md`
§[`DC-FA-VCS-001.a`](../../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 nennt den Mechanismus; die Architektur-Ausnahme (`os`/`io/fs` auch
im git-Adapter) trägt
[ADR-0086](../../adr/0086-pack-alias-fs-os-kapsel-erweiterung.md).

**Was funktionierte.** Die empirische Verifikation VOR dem geschriebenen
Plan (go-gits Quelltext lesen, ein Probe-Repo bauen) traf die richtige
Ursache im ersten Anlauf — kein Fehlversuch, keine verworfene Konstruktion.
Der unabhängige Review (R1) verifizierte das unabhängig noch einmal: eigene
End-to-End-Proben mit echtem `git`-Binary und gebautem Image, ein
bewusstes Zurücknehmen des Fixes (Modul 11) gegen den Parent-Commit, und
eine Lektüre des tatsächlichen go-git-Quelltexts — alle drei bestätigten die
Implementierung unabhängig von der Commit-Botschaft.

**Was anders lief.** Die Recherche ersetzte den geschriebenen Plan (§1), eine
bewusste, benannte Abweichung von „Plan vor Code" (`AGENTS.md` §6 Schritt 4)
— **Lerneintrag:** Wo die korrekte Lösung ohne Kenntnis des internen
Verhaltens einer Fremdbibliothek nicht seriös planbar ist, IST die Recherche
die Plan-Arbeit; das Slice-Kopf-Feld sollte das dann so benennen, statt es
als Verzug zu behandeln. Am Ergebnis dieses Laufs ist keine durch die
Reihenfolge verursachte Qualitätslücke sichtbar (bestätigt durch R1) — die
Abweichung bleibt aber eine, die sich wiederholen kann, sobald ein Slice auf
undokumentiertes Fremdverhalten trifft.

**Review-Runde 1 (drei Findings, alle behoben):**

- **R1-F-1 (MEDIUM):** vier Doku-Stellen (`harness/sensors/adr-check.md`,
  `harness/sensors/trace-check.md`, zwei Abschnitte
  `docs/user/benutzerhandbuch.md`) beschrieben weiterhin den
  Vor-slice-226-Stand. Nachgezogen. **Zweiter Treffer** der Registerklasse
  [`BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`](../observations/BEO-ALL/guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen/observation.md)
  — noch nicht die 3×-Schwelle, aber ein zweiter Treffer in zwei
  aufeinanderfolgenden Slices (`slice-219`, `slice-226`); wird bei einem
  dritten Treffer verkörperungspflichtig.
- **R1-F-2 (MEDIUM):** `packAliases()` löste eine Namenskollision (zwei
  alias-fähige Dateien mit identischem Hash-Suffix) nicht deterministisch
  auf (Go-Map-Iteration). Behoben durch sortierte Verarbeitung mit
  Erstbelegungs-Regel; Regressionstest über fünf Läufe.
- **R1-F-3 (LOW):** ein Kommentar behauptete eine falsche technische
  Begründung für den fest verdrahteten `/`-Trenner in `splitPackPath`
  (tatsächlicher Grund: Linux-only-Distribution, nicht
  Trennzeichenunabhängigkeit von go-git/go-billy). Korrigiert.

**Register.** Drei neue Einträge
([`BEO-ALL/undokumentierte-bibliotheks-interna-als-vertrag`](../observations/BEO-ALL/undokumentierte-bibliotheks-interna-als-vertrag/observation.md),
[`BEO-ALL/dateiname-als-inhalts-beleg-akzeptiert`](../observations/BEO-ALL/dateiname-als-inhalts-beleg-akzeptiert/observation.md),
[`BEO-ALL/beleg-fuer-generalitaet-ist-ein-exemplar`](../observations/BEO-ALL/beleg-fuer-generalitaet-ist-ein-exemplar/observation.md))
für die drei §6-Risiken, die weiter offen bleiben; eine weitere Evidence-Datei
für den bestehenden Eintrag `guide-doku-ausserhalb-spec-bleibt-bei-mechanismuswechsel-stehen`
(R1-F-1, jetzt 2×).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung. **Der Abschnitt entfällt nie**; bedingt
ist allein der Modus-Block am Ende.

Dieses Repo führt **drei** Prüfungen — die zwei kanonischen und, als
Adaption, den Nachtlauf-Stand ([`MR-053`](../../../../harness/conventions.md#mr-053)).

**Die drei Vorprüfungen sind bei der Beanspruchung am 2026-09-17 gelesen.**

**Vorgelagert — Sub-Area-Wahl prüfen:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:363-364 -->

> **Sub-Area-Wahl prüfen.** Jede Sub-Area, die der Slice als berührt führt,
> muss das Inklusionskriterium erfüllen — drei Achsen, Schwelle ≥ 2

**Eine** Sub-Area: `internal/adapter/driven/git/` (der neue Dekorator und die
Kapsel-Erweiterung selbst) — trägt keine eigene Modus-Deklaration in
`harness/conventions.md` und fällt damit unter den Default `*`. Greenfield,
wie der Rest des Produkts.

**Vorgelagert — offene Beobachtungen sichten:**

<!-- d-check:cite .harness/baseline/v6.9.0/regelwerk/modul-05-planning-harness.md:369-369 -->

> **Offene Beobachtungen sichten.** Das

Register durchgegangen (gemergter Stand, 42 Verzeichnisse). **Ein Eintrag
ist einschlägig:**

- [`racily-clean-git-fixture`](../observations/BEO-ALL/racily-clean-git-fixture/observation.md)
  (Sub-Area `internal/adapter/driven/git`, unter der Schwelle) — die neuen
  Testfixtures schreiben Dateien über den bestehenden `put()`-Helfer, der
  `coretest.GitFixtureRewriteHazard` bereits trägt; kein neuer Fundort.

**Keine** der übrigen Einträge trifft `internal/adapter/driven/git/` als
Sub-Area.

**Vorgelagert — Nachtlauf-Stand lesen**
([`MR-053`](../../../../harness/conventions.md#mr-053)):

`make nightly-state` am 2026-09-17 gelesen: `image-scan.yml` grün
(2026-09-17T08:41:27Z). `upstream-drift.yml` **rot** (2026-09-17T05:38:51Z)
— laut eigener Meldung eine **planmäßige** Fremd-Release-Benachrichtigung
([`MR-051`](../../../../harness/conventions.md#mr-051)), keine unerwartete;
betrifft gepinnte Fremd-Bestände, nicht diesen Slice.

**Modus-Begründung:** die einzige berührte Sub-Area ist GF — kein
Begründungsblock nötig.
