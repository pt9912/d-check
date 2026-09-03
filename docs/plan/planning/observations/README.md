# Beobachtungs-Register

Verzeichnis-Form seit `v6.0.0` (migriert in slice-195; vorher eine Tabelle
an dieser Stelle, `observations.md`). Regeln: Baseline-Regelwerk
[`modul-06-roadmap.md`](../../../../.harness/baseline/v6.0.0/regelwerk/modul-06-roadmap.md)
§Das Beobachtungs-Register.

Je Beobachtung ein Verzeichnis `BEO-<KUERZEL>/<slug>/` mit drei Dateien,
drei Lebensdauern: `observation.md` (unveränderlich ab Anlage — Bezeichnung,
Sub-Area), `state.md` (veränderlich — `offen` oder einer von drei Ausgängen:
`verkörpert`/`geplant`/`gestrichen`), `evidence/<vorgangs-id>.md` (eine
Datei je abgeschlossenem Vorgang — Slice, Welle oder Review-Report; der
Zähler ist die Zahl dieser Dateien, kein gepflegtes Feld).

**Wer schreibt:** die Slice-Closure — neues Verzeichnis oder eine weitere
Evidence-Datei. **Wer liest:** die Welle-Closure (Lese-Schritt bei 3×,
ohne Wellen-Betrieb löst die Slice-Closure selbst aus) und die
Slice-Planung (Sichtungs-Schritt darunter, §7/§8 jedes Slice-Plans).

**Gestrichen heißt nicht gelöscht** — das Verzeichnis bleibt liegen,
`state.md` trägt `gestrichen` mit Begründung; wer still löscht, macht eine
Beobachtung ununterscheidbar von einer, die es nie gab.

**Ist nichts offen**, steht hier nur diese Datei — ein leeres Verzeichnis
führt `git` nicht, und die leere Ablage ist selbst die Aussage.

Ein Review-Report hat in diesem Repo keinen Lifecycle-Ruheort (`docs/reviews/`
ist flach); für ihn als Beleg genügt Existenz. Eingefrorene Bestände
(`done/`, `docs/reviews/`, `harness/conventions/done/`) zitieren die
Register-Form ihrer Zeit und werden nicht nachgezogen.
