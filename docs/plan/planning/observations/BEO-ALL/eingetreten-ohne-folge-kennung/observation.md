# Ein eingetretenes Risiko wird im Slice behoben und trägt weder Carveout noch Folge-Kennung

**Sub-Area:** `*`

Der Kanon führt *eingetreten* mit einem Carveout oder einer Folge-Slice-Kennung
(Baseline-Regelwerk `modul-05-planning-harness.md` §Offene Risiken werden bei
Closure aufgelöst). Der Bestand schreibt für ein Risiko, das eintrat und im
selben Slice behoben wurde, *eingetreten — im selben Slice behoben*, ohne
Kennung. Die Form ist im Repo üblich, aber nirgends als Abweichung deklariert;
der Ausgangs-Wächter prüft nur das erste Wort. Ableiter: entweder die Form in
[`MR-049`](../../../../../../harness/conventions.md#mr-049) deklarieren oder den
Fall unter *entfallen* mit Begründung führen.
