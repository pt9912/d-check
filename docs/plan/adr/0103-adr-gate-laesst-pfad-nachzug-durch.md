# ADR-0103: Das ADR-Gate lässt einen reinen Pfad-Nachzug durch

**Status:** Accepted

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
(Lastenheft 0.102.0, Schlüssel `vcs.ignore-link-targets`);
[ADR-0016](0016-adr-immutable-gate.md), [ADR-0024](0024-vcs-immutable-gate.md);
`AGENTS.md` §3.5 und §3.6; slice-267 <!-- d-check:status-provenance -->.

**Schärft:** [`DC-FA-VCS-001`](../../../spec/lastenheft.md#dc-fa-vcs-001--git-diff-immutabilität-des-core-über-eine-commit-range-modul-vcs-opt-in)
— welche Änderung an einer `Accepted`-ADR dieses Repo als inhaltlich zählt.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`make adr-check` hält die `Accepted`-ADRs dieses Repos über das Modul `vcs`
unveränderlich; erlaubt waren zwei Dinge: Anhänge an `## Geschichte` und der
Übergang der Status-Zeile. Das Baseline-Regelwerk (`modul-04-adrs.md`
§Nachzug ist keine Überschreibung) nimmt einen dritten Fall ausdrücklich aus:
den **Referenz-/Pfad-Nachzug**. Verschiebt sich das Ziel eines Verweises, darf
die ADR ihm folgen, ohne dass eine Folge-ADR nötig ist — die Entscheidung
bleibt unberührt.

Das Gate kannte den Fall nicht und meldete jeden geänderten Link-Pfad als
Körper-Änderung. Ein Umzug einer verlinkten Datei — Anlass ist
`docs/user/releasing.md`, auf die zwei `Accepted`-ADRs verweisen — ließ damit
nur zwei Wege: tote Links in den ADRs stehen lassen und per `ignore-refs`
ausnehmen, oder das Gate umgehen.

Das Gate dafür zu öffnen, ist eine Lockerung einer Prüfregel; `AGENTS.md` §3.6
verlangt dafür eine ADR.

## Entscheidung

1. Das Modul `vcs` bekommt den opt-in-Schlüssel `vcs.ignore-link-targets`:
   beim Bilden des Core wird das Ziel jedes Inline-Links, jedes Bilds und jeder
   Referenz-Definition geleert. Eine Änderung, die nur Ziele ändert, ergibt
   denselben Core.
2. Dieses Repo schaltet den Schlüssel für seine ADRs ein (`.d-check.yml`,
   Block `vcs`). `make adr-check` lässt damit einen reinen Pfad-Nachzug durch.
3. Weiterhin ein Befund bleibt jede Änderung am Linktext, an der übrigen
   Zeile, am Titel eines Links und das Hinzufügen oder Entfernen eines Links.
4. `AGENTS.md` §3.5 nennt den Pfad-Nachzug neben `## Geschichte` und dem
   Status-Übergang.

## Verglichene Alternativen

| Alternative | Warum nicht |
|---|---|
| Tote Links in den ADRs stehen lassen, per `ignore-refs` ausnehmen | Eine deklarierte Gate-Senkung am Link-Gate statt am ADR-Gate; der Leser einer ADR findet das Ziel nicht mehr, obwohl es existiert |
| Den Umzug lassen | Die Ablage der Doku richtete sich dann nach der Unveränderlichkeit eines Verweises, nicht nach ihrem Gegenstand |
| Jede Zeilenänderung mit gleichem Text außerhalb der Klammern erlauben | Unschärfer als die Ziele selbst — das Muster soll genau das Ziel freigeben, nicht beliebige Zeichen einer Zeile |
| Den Nachzug per Hand und `--no-verify` einspielen | Umgeht den Hook, nicht die CI; und jede Umgehung macht das Gate zur Gewohnheitsfrage |

## Konsequenzen

- Ein Umzug einer verlinkten Datei zieht die Verweise in `Accepted`-ADRs mit
  nach, im selben Commit wie die übrigen Verweise.
- **Grenze:** Das Gate sieht nur die Form. Ob das neue Ziel dieselbe Sache
  meint wie das alte, prüft es nicht — ein Nachzug auf ein anderes Dokument
  sähe genauso aus. Das trägt der Review des Commits.
- **Grenze:** Die Leerung ist zeilenweise und ohne Code-Kontext; ein Link-Ziel
  in Inline-Code oder einem Codeblock wird ebenso geleert.

## Fitness Function (falls maschinell prüfbar)

Der Test `TestVCSIgnoreLinkTargets` hält beide Seiten: ein reiner Nachzug ist
mit dem Schlüssel keine Drift und ohne ihn eine; Linktext, Zeile, Titel und
Zahl der Links bleiben mit dem Schlüssel Drift. `make adr-check` über die
Range des Umzugs (slice-268 <!-- d-check:status-provenance -->) ist der Beleg
am Bestand.

## Re-Evaluierungs-Trigger

- Das Baseline-Regelwerk ändert die Ausnahme für den Referenz-/Pfad-Nachzug.
- Ein Review findet einen Nachzug, der die Aussage einer ADR verschoben hat.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-09 | Proposed → Accepted (Auftraggeber-Entscheid zum Schnitt von slice-267) |
| 2026-10-09 | Grenzen präzisiert in [`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 4: geleert wird nur ein vollständiger Link; was die Erkennung nicht trifft, bleibt Drift; eine Zeile `[Wort]: wort.md` gilt als Referenz-Definition. Alternative nachgetragen: Verweis über eine Kennung statt einer Adresse oder ein Verfall-Vermerk (Baseline-Regelwerk `modul-04-adrs.md` §Nachzug ist keine Überschreibung) — trägt für künftige ADRs, nicht für den Bestand, dessen Verweise schon Adressen sind |
| 2026-10-09 | Berichtigt zur vorigen Zeile: Die Grenzen der Erkennung legt [`DC-FA-VCS-001.a`](../../../spec/spezifikation.md#dc-fa-vcs-001a--git-diff-immutabilität-über-eine-commit-range-vcs) Schritt 4 fest. Das Baseline-Regelwerk (`modul-04-adrs.md` §Nachzug ist keine Überschreibung) nimmt den Referenz-/Pfad-Nachzug ausdrücklich von der Überschreibung aus und empfiehlt, ihn durch einen Verweis über eine Kennung oder einen Verfall-Vermerk zu vermeiden. Beides trägt den Anlass nicht: Eine Kennung fehlt dem Bestand, dessen Verweise Adressen sind, und ein Verfall-Vermerk ließe einen Verweis verfallen, dessen Ziel nicht verschwindet, sondern umzieht |
