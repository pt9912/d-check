# archive-wave

Setzt [Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur,
Schritt 4](https://github.com/pt9912/ai-harness-course) um: Archivierung
einer geschlossenen Welle in AI-Harness-Kurs-konformen Repos.

Eigenständiges Werkzeug mit eigenem `go.mod` — kein Import aus dem Repo, in
dem es liegt. Portabel: dieses Verzeichnis kann unverändert in jedes Repo mit
demselben `docs/plan/planning/`-Layout kopiert werden.

## Was es tut

Für eine übergebene Wellen-ID (`-welle=welle-NN`):

1. Sammelt alle `docs/plan/planning/done/slice-*.md`, deren `**Welle:**`-Feld
   exakt diese ID nennt.
2. Ordnet `docs/reviews/*.md` per Dateiname (`slice-<NNN>` im Namen) zu, 1:N.
3. Baut `docs/plan/planning/done/<welle-id>/archiv.zip` mit dem Welle-Plan,
   allen gesammelten Slice-Volltexten und allen zugeordneten Review-Reports.
4. Ersetzt Welle-Plan und Slices durch Stubs am neuen Ort
   `docs/plan/planning/done/<welle-id>/` (Review-Reports bekommen keinen
   Stub — sie haben keine Identität jenseits ihres Slice).
5. Zieht repo-weite Markdown-Verweise auf die verschobenen Dateien nach.

## Nutzung

Docker/make-only — kein Host-Go nötig, kein roher `docker`-Aufruf: jede
Operation steht hinter einem Target im eigenen [`Makefile`](Makefile), damit
dieses Verzeichnis dieselbe Disziplin trägt wie das aufrufende Repo
(`AGENTS.md` §3.1, falls vorhanden). Aus diesem Verzeichnis:

```
make test                                    # Testsuite
make run WELLE=welle-42                      # Dry-Run gegen den Repo-Zweig eine Ebene über tools/
make run WELLE=welle-42 APPLY=1              # schreibt
make run WELLE=welle-42 ROOT=/pfad/zum/repo  # gegen ein anderes Repo
```

`ROOT` ist optional (Default: die zwei Verzeichnisebenen über diesem
`Makefile` — die Annahme `<repo>/tools/archive-wave/`; bei abweichender
Ablage explizit setzen). Der Mount ist **nur bei `APPLY=1` beschreibbar** —
im Dry-Run-Default hängt der Container am Repo mit `:ro`, dieselbe
Read-only-Form wie d-checks eigenes Image, nicht ein durchgehend
beschreibbarer Voll-Mount. `-u` bindet den Dateieigner an den aufrufenden
Bediener statt an eine feste Image-UID.

In d-check selbst delegiert `make archive-wave WELLE=welle-42 [APPLY=1]`
(bzw. `make archive-wave-test`) an genau dieses lokale Makefile — eine
Quelle für den Docker-Aufruf, kein Duplikat (siehe `AGENTS.md` §4).

## Grenzen

- Nur explizite Wellen-Zugehörigkeit über das `**Welle:**`-Feld — wellenlose
  Alt-Slices werden nicht eingesammelt; ihre Zuordnung ist eine eigene
  Entscheidung des aufrufenden Repos.
- Review-Zuordnung ist eine Dateinamen-Heuristik — ein Review ohne
  Slice-Kennung im Namen bleibt unzugeordnet.
- Kein Rollback bei einem Fehler mitten in `-apply`: derselbe Umgang wie bei
  jedem `git mv` — der Bediener sieht den Zwischenstand per `git status`.
- Der Verweis-Nachzug löst einen nicht-`/`-präfigierten Verweis relativ zum
  **Verzeichnis der verweisenden Datei** auf (`path.Join(Dir(quelle), ziel)`)
  — die geschwister-relative Form. Ein Verzeichnis-Präfix-Verweis (voller
  Repo-Pfad ab `docs/...`) löst damit nur korrekt auf, wenn die verweisende
  Datei im Repo-Wurzelverzeichnis liegt (heute in diesem Repo für jeden
  bestehenden Verweis dieser Form der Fall). Ein künftiger
  Verzeichnis-Präfix-Verweis aus einer **nicht-Wurzel**-Datei würde still
  übersehen, statt gemeldet zu werden.
