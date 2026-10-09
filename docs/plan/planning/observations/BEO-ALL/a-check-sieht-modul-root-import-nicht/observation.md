# a-check löst den Importpfad des Modul-Roots nicht auf

**Sub-Area:** `*`

Ein Paket im Modul-Root (Importpfad = Modulpfad) ist für a-check keine Datei
einer Schicht: Importiert der Kern es, bleibt `make arch-check` grün, auch
wenn das Paket in `.a-check.yml` eine Schicht hat. Die Kante hält in diesem
Repo ein eigener Test
([ADR-0107](../../../../adr/0107-mitgelieferte-dokumente-im-modul-root.md)).
Ableiter: bei einem Paket außerhalb von Unterverzeichnissen eine
Gegenprobe fahren, bevor man der Schicht-Zuordnung traut; ein Fall für einen
Hinweis an das Schwester-Repo a-check.
