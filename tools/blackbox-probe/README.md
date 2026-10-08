# Black-Box-Probe — Fixtures

Jedes Unterverzeichnis von `fixtures/` ist ein kleines Repository, das `make blackbox-probe`
als `/repo` in Vorher- und Nachher-Image mountet. Ein Fixture gehört hierher,
wenn ein Slice eine Zusage „ohne Schalter unverändert" an einem Modul trifft,
das noch keines hat; es zeigt den **Default**-Zustand des Moduls, nicht den
neuen Schalter.

`fixtures/sauber/` und `fixtures/links/` sind zugleich der **Kanarienlauf**: sie
müssen mit Exit 0 bzw. 1 enden. Wer ihren Exit ändert, ändert die Umgebungsprüfung der
Probe mit.
