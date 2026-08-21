# CO-NNN: <Kurztitel des Carveouts>

> **Template-Hinweis.** Vorlage für einen Carveout. Kopiere nach
> `docs/plan/carveouts/CO-<NNN>-<kurztitel-kebab>.md` und ersetze
> Platzhalter. Lösche diesen Block. Vergiss nicht, den Carveout-Index
> in `docs/plan/carveouts/README.md` zu aktualisieren.

**Status:** Aktiv | Aufgelöst (`done/CO-NNN-...md`).

**Datum angelegt:** YYYY-MM-DD. **Letzte Prüfung:** YYYY-MM-DD.

**Betroffenes Gate:** `<make-target>` (z.B. `coverage-gate`, `noqa-gate`,
`arch-check`).

**Geltungsbereich:** <Pfad / Modul / Datei>

**Folge-Slice:** [`slice-<NN>-...md`](../planning/<state>/slice-NN-...md)

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout — ein
Carveout braucht immer einen Auflösungs-Trigger **und** einen Folge-Slice.

---

## Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — technische Begründung, keine
„noch nicht geschafft"-Aussagen.

<!--
Warum kann das Gate hier (jetzt) nicht hart sein? Technische
Begründung, keine "noch nicht geschafft"-Aussagen.

Erlaubte Begründungen:
- "Bibliothek X bietet noch keine API für die Prüfung."
- "Bootstrap: aktueller Reifegrad < Meilenstein-Schwelle."
- "Externe Abhängigkeit ist außerhalb unserer Kontrolle."
-->

<…>

## Auflösungs-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist" ist kein Trigger.

<!--
Wann verschwindet der Carveout? Konkret, prüfbar.

Falsch: "wenn Zeit ist", "perspektivisch".
Richtig: "Bei Meilenstein M3.", "Wenn pkg/x > 500 LOC.",
"Wenn Bibliothek Y v2 ausgeliefert ist."
-->

<…>

## Geltungs-Konfiguration

<!--
Wenn der Carveout im Gate-Tool konfiguriert ist: wo steht es?
Beispiel: pyproject.toml [tool.coverage.report] omit = [...]
mit Kommentar "CO-NNN".
-->

| Datei | Zeile/Section | Wert |
|---|---|---|
| <…> | <…> | <…> |

## Verifikation (nach Auflösung)

<!--
Welche Schritte schließen den Carveout?
-->

- [ ] Gate ist für den Geltungsbereich aktiviert (Gate-Konfiguration aktualisiert).
- [ ] `make gates` grün ohne Ausnahme.
- [ ] Datei wird nach `docs/plan/carveouts/done/` bewegt (reiner `git mv`). <!-- d-check:ignore (done/ entsteht erst bei erster Carveout-Auflösung) -->
- [ ] Folge-Slice geschlossen oder explizit dokumentiert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| YYYY-MM-DD | Angelegt | <Slice-Datei> |
| YYYY-MM-DD | Geprüft, weiterhin gültig | — |
