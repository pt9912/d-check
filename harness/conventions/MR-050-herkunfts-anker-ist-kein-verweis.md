# MR-050 — Der Herkunfts-Anker einer Hard Rule ist kein Slice-Verweis (schärft MR-045)

- **Status:** Accepted
- **Ersetzt-Baseline-Regel:** keine — im Gegenteil, dieser Eintrag **stellt die
  Baseline wieder her**. [`MR-045`](../conventions.md#mr-045) formulierte
  „`AGENTS.md` trägt **keine** `slice-<NNN>`-Verweise" so breit, dass er auch
  den Herkunfts-Anker erfasste, den der Kanon ausdrücklich verlangt:
  [`modul-09-implementierung.md`](../../.harness/baseline/v5.15.0/regelwerk/modul-09-implementierung.md)
  — *„Eine Hard Rule, die aus dem Steering Loop entstand …, trägt den
  Herkunfts-Anker `(seit welle-<NN>)` — ohne Welle `(seit slice-<NNN>)`"*.
- **Datum:** 2026-08-27
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) und
  [`harness/README.md`](../README.md).
- **Adaption:** Der **Herkunfts-Anker** einer Hard Rule ist von
  [`MR-045`](../conventions.md#mr-045) **ausgenommen**. Ein wellenloser Slice
  hinterlässt `(seit slice-<NNN>)` als Feld — kein Link, kein Pfad, kein Zeiger
  auf das Planungs-Artefakt.

  **Die Trennlinie ist die, die `MR-045` selbst zieht:** verboten ist der
  *„Verweis auf das Planungs-Artefakt"*, der mit ihm durch vier Verzeichnisse
  wandert und am Ende ins Leere zeigt oder gepflegt werden muss. Ein
  Herkunfts-Anker wandert nicht — er nennt eine Nummer, die sich nie ändert, und
  löst nichts auf.
- **Begründung:** [`MR-045`](../conventions.md#mr-045) ist gegen einen echten
  Fall geschrieben (ein Slice-Link in §3.1, der mit dem Lifecycle-Move
  gebrochen wäre) und hat daraus eine Regel gemacht, die den Anlass
  überschreitet — der Bestand kannte damals beide Formen nicht getrennt. Wer
  dem breiten Wortlaut folgt, lässt den vom Kanon verlangten Anker weg und
  merkt es nicht, weil kein Gate ihn prüft.

  **Der Widerspruch ist hiermit gemeldet, nicht stillschweigend aufgelöst**
  ([`AGENTS.md`](../../AGENTS.md) §1): die höherrangige Quelle ist der Kanon,
  die niedriger rangierte wird angepasst — hier durch diesen Nachtrag, nicht
  durch einen Edit an [`MR-045`](../conventions.md#mr-045).
- **Auflösungs-Trigger:** der Kanon lässt den Herkunfts-Anker fallen, oder
  dieses Repo führt wieder Wellen für jede Hard Rule. Dann ist die Ausnahme
  gegenstandslos.
