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

Docker-only — kein Host-Go nötig, dieselbe Toolchain-Disziplin wie das
aufrufende Repo (`AGENTS.md` §3.1, falls vorhanden). Aus diesem Verzeichnis:

```
docker build --target test -f Dockerfile -t archive-wave:test .   # Testsuite
docker build -f Dockerfile -t archive-wave:latest .                # Binary bauen

# Dry-Run (Default) -- schreibt nichts, listet die geplante Operation:
docker run --rm -v "/pfad/zum/repo":/repo:ro \
    archive-wave:latest -root=/repo -welle=welle-42

# Anwenden -- schreibt, Datei-Eigner via -u:
docker run --rm -u "$(id -u):$(id -g)" -v "/pfad/zum/repo":/repo \
    archive-wave:latest -root=/repo -welle=welle-42 -apply
```

`-root` ist optional (Default `.`). In d-check selbst führt
`make archive-wave WELLE=welle-42 [APPLY=1]` denselben Build- und
Run-Schritt aus (siehe `AGENTS.md` §4).

## Grenzen

- Nur explizite Wellen-Zugehörigkeit über das `**Welle:**`-Feld — wellenlose
  Alt-Slices werden nicht eingesammelt; ihre Zuordnung ist eine eigene
  Entscheidung des aufrufenden Repos.
- Review-Zuordnung ist eine Dateinamen-Heuristik — ein Review ohne
  Slice-Kennung im Namen bleibt unzugeordnet.
- Kein Rollback bei einem Fehler mitten in `-apply`: derselbe Umgang wie bei
  jedem `git mv` — der Bediener sieht den Zwischenstand per `git status`.
