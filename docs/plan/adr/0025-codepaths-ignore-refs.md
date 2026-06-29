# ADR-0025 — Referenz-Ventil `codepaths.ignore-refs`: Tombstone-Register entfernter Code-Pfade löst die Frozen-Doc-Refactoring-Falle und entfernt adr-immutable-check.sh

**Status:** Proposed
**Datum:** 2026-06-29
**Autor:** pt9912
**Bezug:** [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)
(Modul `codepaths`); **Supersedes** die „Skript bleibt pfad-stabil"-Teilentscheidung
von [ADR-0024](0024-vcs-immutable-gate.md) (deren VCS-Port-/Modul-`vcs`-/Dogfood-Kern
**gültig bleibt** — Teil-Supersede, wie [ADR-0002](0002-distribution-ghcr-image.md)/[ADR-0014](0014-latest-tag-fuer-stabile-releases.md));
Bezug [ADR-0016](0016-adr-immutable-gate.md) (die immutable Policy, deren Inline-Referenz
die Falle erst auslöste) und [ADR-0008](0008-reparatur-ableitbarkeit.md)
(VCS-Port als eigene künftige Anforderung — nicht dieser Mechanismus);
Ventil-Vorbild das datei-weite `codepaths.exempt-paths` und der
zeilenweise `d-check:ignore`-Marker.
**Schärft:** die Algorithmus-Sektion
[`spec/spezifikation.md` §DC-FA-CODE-001.a](../../../spec/spezifikation.md#dc-fa-code-001a--pfade-in-inline-code)
(Referenz-Ventil-Schritt vor `codepath-missing`) und das §2-Schema
(`codepaths.ignore-refs`).

## Kontext

Das Modul `codepaths`
([`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in))
prüft Inline-Code-Pfade als **lebende** Zeiger: der Pfad muss *jetzt* existieren.
Immutable bzw. historische Doku trägt aber **eingefrorene** Verweise — `Accepted`-ADRs
(maschinell unveränderlich, [ADR-0016](0016-adr-immutable-gate.md)) und abgeschlossene
`done/`-Slices zitieren Pfade, die zum Zeitpunkt der Entscheidung galten. Wird der Code
später refaktoriert oder gelöscht, **dangelt** der eingefrorene Verweis →
`codepath-missing`. Das ist **nicht fixbar**: die ADR ist immutabel, ein Edit an ihrem
Core bräche `make adr-check` ([ADR-0024](0024-vcs-immutable-gate.md)). Die Falle tritt
bei **jedem** solchen Refactoring auf.

Konkret scheiterte zuletzt die Entfernung von
`tools/adr-immutable-check.sh` (durch das Modul `vcs` abgelöst) genau daran, dass
die immutablen [ADR-0016](0016-adr-immutable-gate.md)/[ADR-0024](0024-vcs-immutable-gate.md)
das Skript in Inline-Code referenzieren. [ADR-0024](0024-vcs-immutable-gate.md) entschied
darum, das Skript **pfad-stabil zu behalten** — ein Workaround gegen das Symptom, kein
Mechanismus gegen die Ursache. Der Auftraggeber verlangte „eine bessere Lösung, dass
solch ein Problem … nicht immer wieder auftaucht".

Die beiden vorhandenen Ventile passen nicht:

- Der zeilenweise `d-check:ignore`-Marker müsste **in die ADR-Zeile** — ein Core-Edit
  an immutabler Doku (→ `adr-check` FAIL) und für bereits `Accepted`-ADRs nachträglich
  unerreichbar; zudem ist sein Default („dieser Verweis ist lebendig, nur diese Zeile
  nicht prüfen") für ein **historisches Record** falsch herum.
- Das datei-weite `codepaths.exempt-paths` nähme die **ganze** ADR-Datei aus der Prüfung
  — und verlöre damit die `codepaths`-Abdeckung ihrer **übrigen**, noch lebenden
  Verweise. Zu grob.

## Entscheidung

Ein Config-Schlüssel `codepaths.ignore-refs` (Glob-Liste, `matchGlob` wie
`codepaths.exempt-paths`/`scan.ignore`). Ein in der Auflösung Wurzel-relativ bestimmter
**Ziel-Pfad**, der einem Eintrag entspricht, wird **vor** der Existenz-/Escape-/Anker-
Prüfung **übersprungen** — `codepath-missing` (und `anchor-missing`/`repo-escape` für
diesen Pfad) entfallen. Das Ventil wirkt **referenz-weit**: unabhängig von Datei und
Zeile, anders als das datei-weite `exempt-paths` und der zeilenweise `d-check:ignore`.
Es ist ein **Tombstone-Register** bewusst entfernter/historischer Artefakte, deren Pfad
eingefrorene Doku noch zitiert. Default leer → ohne `ignore-refs` byte-identisch
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).

**Bewusster Akt mit Gate, nicht stiller Default.** `ignore-refs` ist ein Opt-in-Register:
wer einen Pfad entfernt und den Eintrag **vergisst**, fällt weiter auf `codepath-missing`.
Der Gate **erzwingt** die Tombstone-Deklaration, statt sie still durchzulassen — dieselbe
Disziplin wie die Versions-/Digest-Pins. Nichts dangelt unbemerkt.

**Per-Pfad, nicht per-Doc-Klasse.** Bewusst **kein** pauschales Ausnehmen von
`docs/plan/adr/**`/`done/**` — das verlöre die Abdeckung der **übrigen** Verweise dieser
Dateien. `ignore-refs` nimmt nur den **einen** entfernten Pfad aus; alles andere bleibt
geprüft.

**Config statt Doc-Marker.** Das Register lebt in [`.d-check.yml`](../../../.d-check.yml),
nicht in der Doku → **kein** Edit an immutabler Doku nötig, und es **retrofittet**
bestehende `Accepted`-ADRs. (Der zeilenweise `d-check:ignore`-Marker bleibt für
historische Verweise in **lebenden** Dateien.)

**Erster Anwendungsfall: `tools/adr-immutable-check.sh` wird entfernt.** `git rm` +
Eintrag in `codepaths.ignore-refs`; die eingefrorenen
[ADR-0016](0016-adr-immutable-gate.md)/[ADR-0024](0024-vcs-immutable-gate.md)-Inline-
Referenzen werden so zu **deklarierten** Tombstones. Damit nimmt diese ADR die
„Skript bleibt pfad-stabil"-Teilentscheidung von [ADR-0024](0024-vcs-immutable-gate.md)
zurück; deren VCS-Port, das Modul `vcs` und der Dogfood-Ersatz des `adr-check`-Gates
bleiben **unberührt**. Der Negativ-Selbsttest des Skripts lebt als
Akzeptanztest im Modul `vcs` weiter — kein Garantieverlust.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **`codepaths.ignore-refs` (gewählt)** | retrofittet immutable Doku **ohne** Edit; per-Pfad → übrige Verweise bleiben geprüft; bewusster Akt **mit Gate**; verteilt als Config (kein Code/Modul/Dependency) | Register wächst über die Zeit; globaler (referenz-weiter) Scope greift auch in lebenden Dateien |
| `d-check:ignore`-Zeilenmarker in der ADR | nutzt vorhandenes Ventil | Core-Edit an immutabler ADR → `adr-check` FAIL; erreicht `Accepted`-ADRs nicht nachträglich; Default falsch herum für Records |
| `exempt-paths` (ganze Datei) | vorhanden | nimmt die **ganze** Datei aus → übrige lebende Verweise ungeprüft; zu grob |
| klassenweites Ausnehmen (`adr/**`, `done/**`) | einfach | verliert die `codepaths`-Abdeckung ganzer Doc-Klassen; verschleiert echte Drift |
| Skript behalten (Status quo [ADR-0024](0024-vcs-immutable-gate.md)) | schon da | Workaround gegen das Symptom; die Falle bleibt bei **jedem** künftigen Refactoring |

## Konsequenzen

- **Default-Lauf byte-identisch.** Kein neues Modul, keine neue Dependency; ohne
  gesetztes `ignore-refs` ändert sich nichts
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- **Die Frozen-Doc-Falle ist für die `codepaths`-Achse (Inline-Code-Pfade) aufgelöst.**
  Künftiges Refactoring/Löschen, das in **Inline-Code** zitierte Pfade betrifft, ist sauber
  deklarierbar — ohne Edit an immutabler Doku, ohne ganze Doc-Klassen aus dem Check zu
  nehmen. **Nicht** abgedeckt ist die `links`-Achse: ein entfernter Pfad, der als
  **Markdown-Link** (statt Inline-Code) in immutabler Doku steht, feuert weiter
  `target-missing` des Moduls `links` — `ignore-refs` ist `codepaths`-lokal, `links` hat
  kein referenz-weites Pendant (s. Re-Evaluierungs-Trigger). Aktuell rein latent: kein
  immutables Dokument verlinkt den entfernten Pfad in Link-Form.
- **`adr-immutable-check.sh` ist entfernt.** Die `codepaths`-Abdeckung seiner übrigen
  Datei-Nachbarn bleibt; nur der eine Tombstone-Pfad ist still.
- **Register wächst und ist ein sichtbarer Diff-Akt.** Bewusst — ein ehrliches
  Tombstone-Verzeichnis, bounded auf tatsächlich entfernte, noch zitierte Artefakte. Ein
  zu breiter Glob könnte echte Drift verstecken → per-Pfad-Disziplin, der Glob ist eng zu
  halten.
- **Globaler Scope.** `ignore-refs` greift überall (auch in lebenden Dateien). Akzeptabel
  — ein bewusst entfernter Pfad ist überall „nicht da"; lebende Verweise räumt man ohnehin
  auf.
- **Auffindbarkeit.** `--print-config`
  ([`DC-FA-CLI-005`](../../../spec/lastenheft.md#dc-fa-cli-005--konfigurations-gerüst-ausgeben))
  und das Benutzerhandbuch dokumentieren `ignore-refs`; ein Adopter findet das Ventil
  ohne Quelltext-Studium.

## Fitness Function

- Mit `tools/adr-immutable-check.sh` in `codepaths.ignore-refs` läuft `make doc-check`
  **grün**, obwohl die Datei entfernt ist — die eingefrorenen
  [ADR-0016](0016-adr-immutable-gate.md)/[ADR-0024](0024-vcs-immutable-gate.md)-Inline-
  Referenzen brechen nicht.
- **Negativprobe:** entfernt man den `ignore-refs`-Eintrag, feuert `codepath-missing` an
  jenen Referenzen — der Gate erzwingt die Tombstone-Deklaration.
- Ein **nicht**-ignorierter fehlender Pfad feuert weiter `codepath-missing` (per-Pfad,
  kein klassenweites Loch).
- Ohne gesetztes `ignore-refs` ist der Befundsatz byte-identisch (opt-in-Selbsttest).
- `--print-config` und das Benutzerhandbuch zeigen `ignore-refs` (Auffindbarkeit).

## Re-Evaluierungs-Trigger

- Das Register wächst stark / wird missbraucht (zu breite Globs) → periodisches Audit,
  ggf. pro Eintrag eine Begründung verlangen.
- Bedarf, einen entfernten Pfad nur in **einer** Datei (nicht referenz-weit) auszunehmen
  → das deckt das vorhandene `exempt-paths`/`d-check:ignore` ab; `ignore-refs` bleibt
  bewusst referenz-weit.
- Eine **immutable** Doku zitiert einen entfernten Pfad als **Markdown-Link** (statt
  Inline-Code) → `target-missing` des Moduls `links` an einer uneditierbaren Datei;
  `ignore-refs` greift dort nicht. Dann braucht die `links`-Achse ein eigenes
  referenz-weites Tombstone-Ventil (eigene Anforderung) — bis dahin bleibt die Mitigation,
  Link-Form-Verweise auf wandernde Artefakte in **lebenden** Dateien zu Inline-Code zu
  entlinken.
- Wunsch nach einer git-basierten „wann/ob der Pfad wirklich entfernt wurde"-Verifikation
  → eigener VCS-Port-Anwendungsfall ([ADR-0008](0008-reparatur-ableitbarkeit.md)/[ADR-0024](0024-vcs-immutable-gate.md)),
  eigene Anforderung — nicht dieser hermetische Config-Mechanismus.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-29 | Entwurf nach Auftraggeber-Auftrag („eine bessere Lösung, dass Refactoring/Löschen nicht immer wieder die `codepaths`-/Immutabilitäts-Falle auslöst"; Idee `ignore-code-ref` im `.d-check.yml`). Erweitert [`DC-FA-CODE-001`](../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) um das Referenz-Ventil `ignore-refs`; **Supersedes** die „Skript bleibt pfad-stabil"-Teilentscheidung von [ADR-0024](0024-vcs-immutable-gate.md) → entfernt `tools/adr-immutable-check.sh`. Begleitet slice-054. Status Proposed. |
