# ADR-0107: Die mitgelieferten Dokumente liegen in einem Paket im Modul-Root

**Status:** Proposed

**Datum:** 2026-10-09

**Autor:** pt9912

**Bezug:** [`DC-FA-CLI-013`](../../../spec/lastenheft.md#dc-fa-cli-013--handbuch-und-spezifikation-aus-dem-werkzeug-lesen)
(die Anforderung, die das Paket trägt); [ADR-0005](0005-modul-layout-hexagon-ordner.md)
(das Modul-Layout, das hier um einen Ort erweitert wird);
[ADR-0029](0029-arch-check-via-a-check.md) (die Import-Prüfung über a-check);
slice-266 <!-- d-check:status-provenance -->.

**Schärft:** — *(keine Spec-Stelle; eine Entscheidung über das Code-Layout)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`--manual` gibt Abschnitte aus Benutzerhandbuch und Spezifikation aus, ohne
Netz und byte-gleich zu den Quelldateien des Build-Stands. Der Weg dahin ist
`go:embed`: Es legt Dateien zur Übersetzungszeit ins Binary. `go:embed`
erreicht aber nur Dateien unter dem Verzeichnis des einbettenden Pakets — die
Dokumente liegen unter `docs/user/` und `spec/`, ein Paket unter `internal/`
kann sie nicht einbetten.

[ADR-0005](0005-modul-layout-hexagon-ordner.md) kennt kein Paket im
Modul-Root. Der Review zeigte, dass `make arch-check` es als „in keiner
Schicht“ ungeprüft ließ, und dass a-check den Importpfad des Modul-Roots
nicht auflöst: Eine Kante vom Kern zum Root-Paket blieb auch nach der
Schicht-Zuordnung unsichtbar.

## Entscheidung

1. Ein Paket `dcheck` im Modul-Root (`manual.go`) bettet Handbuch und
   Spezifikation ein und trägt sonst nichts: keine Logik, keine Importe außer
   `embed`.
2. Nur die Composition Root (`internal/adapter/driving/cli`) importiert es.
   Die Abschnitts-Suche liegt im Kern (`app`) und bekommt die Dokumente als
   Werte übergeben.
3. `.a-check.yml` führt das Paket als eigene Schicht `docs` ohne Kanten, damit
   es nicht als „in keiner Schicht“ durchfällt. Die Rolle ist `domain`, die
   strengste Einstufung: keine ausgehende Kante zu einer anderen Schicht. Dass
   a-check einen Verstoß als „Kern importiert“ meldet, betrifft nur das Wort —
   zum Kern gehört das Paket nicht.
4. Weil a-check den Importpfad des Modul-Roots nicht auflöst, hält ein Test in
   der Composition Root die Kante: Kein Paket unter `internal/hexagon` oder
   `internal/adapter/driven` importiert das Root-Paket. Ein zweiter Test hält
   Entscheidung 1: das Root-Paket importiert nur `embed` und trägt keine
   Funktion — a-check lässt Importe der Standardbibliothek zu, und das
   Coverage-Gate misst nur `internal/`.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon (Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR
(MADR)).

| Option | Pro | Contra |
|---|---|---|
| Nichts tun: keine Dokumente im Binary | Kein neuer Ort im Layout | Ein netzloses Repo kann die Regeln eines Moduls nur durch Probieren ermitteln (CR eines Adopters) |
| Dokumente beim Build unter `internal/` kopieren | Das Layout bleibt unverändert | Eine zweite Fassung im Repo, die gegen die Quelle driften kann, oder ein Kopierschritt in jeder Dockerfile-Stage, die das Paket übersetzt |
| Dateien ins Image legen, ohne Einbettung | Kein Code nötig | Das distroless-Image hat kein `cat`; gelesen würde per `docker cp`, und das Binary allein trüge sie nicht |
| **Gewählt:** Paket im Modul-Root, eigene Schicht, Kante per Test | Eine Quelle, byte-gleich per Konstruktion; die Kante ist geprüft | Ein Ort außerhalb von `internal/`; die Prüfung der Kante liegt in einem Test statt in a-check |

## Konsequenzen

- Jede Änderung an Handbuch oder Spezifikation ändert das Binary.
- **Grenze:** a-check prüft die Kante Kern → Root-Paket nicht; sie hält der
  Test `TestManual_RootPaketNurAusDerCompositionRoot`; dass das Paket nur
  `embed` importiert und keine Funktion trägt, hält
  `TestManual_RootPaketTraegtNurDieEinbettung` — weder a-check noch das
  Coverage-Gate sehen es. Löst a-check den
  Importpfad des Modul-Roots künftig auf, ist der Test entbehrlich.

## Fitness Function (falls maschinell prüfbar)

`TestManual_RootPaketNurAusDerCompositionRoot` wird rot, sobald ein Paket des
Kerns oder ein getriebener Adapter das Root-Paket importiert (gegengeprüft
mit einem Import in `app`). `TestManual_EingebetteteDokumenteGleichDerQuelle`
hält Pfad-Konstante und Einbettungs-Ziel zusammen, damit die Kopfzeile der
Ausgabe die richtige Datei nennt; dass die Einbettung zum Build-Stand passt,
sichert `go:embed` selbst.

## Re-Evaluierungs-Trigger

- a-check löst den Importpfad des Modul-Roots auf.
- Ein weiteres Dokument soll mitgeliefert werden, das nicht unter dem
  Modul-Root liegt.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-10-09 | Proposed (Review R1 zu slice-266: das Root-Paket lag außerhalb von Layout und Schichten) |
