# ADR-0086: Der git-Adapter bekommt Zugriff auf `os`/`io/fs` — für Typ-Konformität, nicht für I/O

**Status:** Accepted

**Datum:** 2026-09-17

**Autor:** pt9912

**Bezug:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in),
[ADR-0005](0005-modul-layout-hexagon-ordner.md) (R2, die Tech-Kapseln),
[ADR-0029](0029-arch-check-via-a-check.md) (a-check als Durchsetzung),
[ADR-0072](0072-workflows-modul.md) (Präzedenzfall: dieselbe Kapsel-Erweiterung
für `gopkg.in/yaml.v3`, dort R3)

**Schärft:** [`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs)
Schritt 2 (die Kandidaten-Auflösung — der eingehende CR und sein Fix liegen
in `slice-226` <!-- d-check:status-provenance -->)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Ein eingehender CR (`ai-harness-init`, 2026-09-17,
`docs/plan/cr/2026-09-17-cr-eingehend-ai-harness-init-pack-praefix.md`) meldet:
Ein Pack, dessen Objekte vollständig gültig sind, aber unter einem anderen
Dateinamens-Präfix als `pack-` liegen (z. B. `loose-<hash>.pack`, wie `git
maintenance run --task=loose-objects` es schreibt), macht `vcs` seit `v0.76.1`
mit `Exit 2` abbrechen — obwohl `git cat-file`/`git fsck` dieselben Objekte
anstandslos lesen. Ursache: go-gits `storage/filesystem/dotgit`-Schicht (`v5.19.2`)
entdeckt und öffnet Packs **ausschließlich** über den rekonstruierten Namen
`pack-<hash>.{pack,idx}` (`DotGit.objectPacks`/`objectPackPath`) — ein anderer
Präfix ist für sie unsichtbar, unabhängig vom Objekt-Inhalt.

Der Fix (`slice-226` <!-- d-check:status-provenance -->) ist ein read-only `billy.Filesystem`-Dekorator
(`packAliasFS`, `internal/adapter/driven/git/packalias.go`), der Packs mit
gültigem Hash-Suffix unter ihrem kanonischen Namen zusätzlich sichtbar macht,
ohne das gemountete Repository zu verändern. Das erfordert, `billy.Filesystem`s
`ReadDir`-Methode zu implementieren — ihre Signatur ist von der Bibliothek
vorgegeben: `ReadDir(path string) ([]os.FileInfo, error)`. `os.FileInfo` ist
seit Go 1.16 ein Typ-Alias auf `io/fs.FileInfo`; **welcher** der beiden Namen
importiert wird, ändert am zugrunde liegenden Typ nichts — beide sind laut
[`.a-check.yml`](../../../.a-check.yml) (R4) auf `internal/adapter/driven/fs/`
gekapselt, und der git-Adapter verletzt diese Kapsel mit jeder Variante.

## Entscheidung

Wir erweitern die `os`/`io/fs`-Kapsel in `.a-check.yml` um
`internal/adapter/driven/git/` (Listenform, wie bereits bei `gopkg.in/yaml.v3`
in [ADR-0072](0072-workflows-modul.md) geschehen) — **nicht** als generelle
Freigabe für host-Dateisystem-I/O, sondern **ausschließlich** für die
Typ-Konformität mit `billy.Filesystem`, das der Adapter ohnehin exklusiv
kapselt ([ADR-0024](0024-vcs-immutable-gate.md), R2b). `packAliasFS` selbst
liest nie direkt über `os`/`io/fs` — es liest und schreibt ausschließlich über
das eingebettete `billy.Filesystem`; `os.FileInfo` erscheint nur als
Schnittstellen-Typ in Methoden-Signaturen und einem schlanken Wrapper
(`aliasFileInfo`), der `Name()` überschreibt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, den CR ablehnen | keine Architektur-Änderung | der gemeldete Defekt bleibt: ein Repo nach `git maintenance run --task=loose-objects` — eine von git selbst empfohlene Routine-Wartung — macht `vcs` fail-closed grundlos rot |
| B — `packAliasFS` in `internal/adapter/driven/fs/` statt in `git/` unterbringen | keine Kapsel-Erweiterung nötig | `fs/` kapselt host-Dateisystem-Zugriff für andere Module (Discovery etc.); der Dekorator ist strukturell an `billy.Filesystem`/go-git gebunden (R2b, [ADR-0024](0024-vcs-immutable-gate.md): go-git nur im git-Adapter) — er dort unterzubringen risse ihn aus seiner einzigen Verwendung heraus und bräche stattdessen die go-git-Kapsel |
| C — go-git selbst patchen/forken, um beliebige Pack-Präfixe zu unterstützen | löst das Problem an der Quelle, käme auch anderen Konsumenten zugute | ein Fork ist eine Supply-Chain- und Wartungslast, die dieses Repo nirgends sonst trägt ([ADR-0002](0002-distribution-ghcr-image.md)/[ADR-0024](0024-vcs-immutable-gate.md): reine, ungeforkte Abhängigkeiten); ein Upstream-PR ist nicht in der Hand dieses Repos und löst den Adopter-Fall nicht kurzfristig |
| **D — Kapsel um den git-Adapter erweitern, Dekorator dort (gewählt)** | kleinster Eingriff, der die Zusage hält („go-git nur im git-Adapter“ bleibt exakt wahr); die Erweiterung ist eng begründet (Typ-Konformität, kein I/O) und durch a-check weiterhin geprüft | die `os`/`io/fs`-Kapsel ist ab jetzt zwei statt einen Adapter breit — eine zweite Stelle, an der ein Reviewer nachschlagen muss, warum |

## Konsequenzen

- **Positiv:** Der CR ist strukturell lösbar, ohne go-git zu forken oder eine
  vierte Toolchain einzuführen ([`MR-046`](../../../harness/conventions.md#mr-046)).
- **Positiv:** Die Erweiterung ist eng — `packAliasFS` importiert `os`
  ausschließlich für den `os.FileInfo`-Schnittstellentyp, nie für eine
  `os.Open`/`os.ReadFile`-artige direkte Host-I/O. `make arch-check` prüft
  weiterhin, dass **kein** anderer Import aus `os`/`io/fs` in diesem Adapter
  entsteht.
- **Negativ, benannt:** Die R4-Kapsel ist jetzt zwei Adapter breit statt einem
  — dieselbe Klasse von Zugeständnis wie bei R3 ([ADR-0072](0072-workflows-modul.md)).
  Ein dritter Adopter-Fall, der eine dritte Kapsel-Erweiterung verlangt, wäre
  ein Signal, die Regel selbst zu überdenken (Re-Evaluierungs-Trigger unten).
- **Folgepflicht:** `spec/spezifikation.md` §[`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 2 bekommt
  einen Nachtrag zur Pack-Auflösung (in `slice-226` <!-- d-check:status-provenance -->, nicht hier — diese ADR
  trägt die Architektur-Entscheidung, nicht die Fach-Spezifikation).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `os`/`io/fs` nur in `internal/adapter/driven/fs/` und `internal/adapter/driven/git/`, `composition_root` weiterhin default-allow | `make arch-check` |
| Go-Tests | `TestAllPathsPackUnterFremdemPraefix` (Fix greift), `TestAllPathsPackMitUnbrauchbaremPraefixBleibtFehlerhaft` (Gegenprobe: kein gültiges Hash-Suffix bleibt fail-closed) | `make test` |

## Re-Evaluierungs-Trigger

**Ein dritter Adapter, der aus demselben Grund (externe Interface-Konformität,
nicht eigene I/O) Zugriff auf eine bereits gekapselte `tech`-Klasse braucht.**
Dann ist die Frage nicht mehr "welchen Adaptern erlauben wir das", sondern ob
die Kapsel-Regel selbst zwischen *Typ-Konformität* und *tatsächlicher I/O*
unterscheiden sollte — eine Unterscheidung, die a-check heute nicht kennt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted | `slice-226` |
