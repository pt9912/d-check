# Black-Box-Probe — Fixtures

Jedes Unterverzeichnis von `fixtures/` ist ein kleines Repository, das `make blackbox-probe`
als `/repo` in Vorher- und Nachher-Image mountet. Ein Fixture gehört hierher,
wenn ein Slice eine Zusage „ohne Schalter unverändert" an einem Modul trifft,
das noch keines hat; es zeigt den **Default**-Zustand des Moduls, nicht den
neuen Schalter.

`fixtures/sauber/` ist zugleich der **Kanarienlauf**: er muss Exit 0 und genau
eine geprüfte Datei liefern. Wer ihn ändert, ändert die Umgebungsprüfung der
Probe mit.
