# Repo-Einstellungen auf GitHub — d-check

Was in den GitHub-Einstellungen des Repos steht und im Klon **nicht** prüfbar
ist. Diese Datei beschreibt den Soll-Zustand und seine Bedienung; welche
Durchsetzungsgrenze er hat, steht in
[`harness/README.md`](../../harness/README.md) §Gate-Taxonomie und
Durchsetzungsgrenzen, die Entscheidung zur CI-Range in
[ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md).

## Das Ruleset `main-nur-admin`

Gilt für `refs/heads/main`, Durchsetzung `active`:

| Regel | Wirkung |
|---|---|
| `deletion` | `main` lässt sich nicht löschen |
| `non_fast_forward` | kein Force-Push auf `main` |
| `pull_request` | Änderungen nur über einen Pull-Request, **ohne** Pflicht-Review (`required_approving_review_count: 0`) |
| `required_status_checks` | der Check `gates + image-test + traceability` muss bestehen |

**Bypass:** die Rolle Repository-Admin, immer (`bypass_mode: always`). Direkte
Pushes eines Admins auf `main` unterliegen dem Check deshalb **nicht** — so
landen die Commits des Harness-Laufs auf `main`. Wer Admins ebenfalls durch
Pull-Requests zwingen will, setzt `bypass_mode: pull_request`; das bremst den
Arbeitsablauf mit direkten Commits und ist hier nicht gewählt.

Zuvor stand an der Stelle der PR-Regel und des Checks die Regel `update`. Sie
sperrte **jede** Änderung an `main` für Nicht-Admins, auch einen Merge mit grüner
CI; jeder Dependabot-PR brauchte einen Admin-Merge. Der Check hängt am Namen des
Jobs: er muss in **beiden** Läufen (Pull-Request und Push) grün sein, deshalb
fällt der Push-Lauf nach einem Force-Push auf die Default-Branch-Range zurück
([ADR-0109](../plan/adr/0109-ci-push-range-faellt-bei-force-push-zurueck.md)).

## Auto-Merge und Dependabot

Im Repo ist **Auto-Merge** eingeschaltet (`allow_auto_merge`). Ein
Dependabot-PR ([ADR-0067](../plan/adr/0067-dependabot-als-hebender-kanal.md)) geht
damit ohne Admin-Recht durch, sobald der Check grün ist:

```bash
gh pr merge <nummer> --auto --rebase
```

Ist der PR veraltet, lässt `@dependabot rebase` als PR-Kommentar ihn neu auf
`main` aufsetzen; die CI läuft danach neu.

## Ansehen und Zurückgehen

```bash
gh api repos/pt9912/d-check/rulesets/24829659        # Ist-Zustand lesen
gh api repos/pt9912/d-check -q .allow_auto_merge     # Auto-Merge
```

Zurück auf den früheren Zustand geht über `PUT` desselben Rulesets mit der
gesicherten Fassung, nachdem `id`, `source`, `source_type`, `node_id`,
`created_at`, `updated_at`, `current_user_can_bypass` und `_links` entfernt sind.
Eine Änderung am Ruleset ist eine Änderung an der Durchsetzung der Gates:
`AGENTS.md` §3.6 gilt auch dort — eine Lockerung gehört in eine ADR.

## Stand und Beobachtung

Umgestellt am 2026-10-20. Beobachtet ist, dass **vor** der Umstellung ein
PR-Merge ohne Admin-Recht trotz grüner Checks abgewiesen wurde. Dass er **seit**
der Umstellung bei grünem Check durchgeht, ist am nächsten Dependabot-PR zu
beobachten und hier nachzutragen.
