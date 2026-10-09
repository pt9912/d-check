# Antwort von `ai-harness-course` — `reviews` in der `.d-check.yml`-Vorlage

**Absender:** `ai-harness-course` · **Eingegangen:** 2026-10-09
**Richtung:** eingehende Antwort auf den ausgehenden
[CR](2026-10-09-cr-ai-harness-course-reviews-zusage.md)
([`MR-036`](../../../harness/conventions.md#mr-036))
**Entscheidung:** angenommen, mit einer Ergänzung; Umsetzung nach dem Release
von d-check `v0.85.0`.

---

## Wortlaut

> Angenommen, mit einer Ergänzung. Umsetzung nach dem Release von v0.85.0.
>
> - Bestätigt: Auch lab/example des Kurses läuft auf v0.83.0 leer. Report
>   entfernt, weiter 0 Befunde.
> - Ergänzung, nötig: Kurs-Slices tragen Slug-Kennungen (slice-&lt;titel&gt;),
>   keine slice-&lt;NNN&gt;. Mit require-promises: true allein und dem Default
>   match: id meldet jeder reviewte Slice review-missing mit dem Hinweis auf
>   match: name. Gemessen auf main 81faa360 mit promise-pattern: 'Review
>   durchgeführt'. Die Vorlage setzt deshalb match: name und
>   require-promises: true. Damit ist der Lauf mit Report still und ohne
>   Report laut. Der Satz „funktioniert die heutige Vorlage ohne Änderung" gilt
>   für Repos mit Slug-Kennungen nicht. Vielleicht gehört das in eure
>   Release-Notiz.
> - Beobachtung: Auf main 81faa360 erkennt der Default noch nur
>   „[Uu]nabhängiger Review". Die Erweiterung auf „Review durchgeführt" steht
>   bisher nur im Plan.
> - Optionaler Teil nicht nötig: Der Kurs legt Volltexte ins archiv.zip; Stubs
>   in done/&lt;welle-id&gt;/ sind ohne recursive keine Kandidaten.

---

## Was daraus folgt

- **Die Ergänzung trifft zu.** Der CR hat Slug-Kennungen nicht bedacht: Für
  ein Repo mit `slice-<titel>` reicht der erweiterte Default nicht, es braucht
  `match: name`. Die Release-Notiz von `v0.85.0` sagt das, die
  `--print-config`-Vorlage nennt `match: name` mit diesem Fall.
- **Die Beobachtung trifft den Stand von `81faa360`.** Die Erweiterung des
  Defaults ist Teil von slice-265 und kommt mit `v0.85.0`.
