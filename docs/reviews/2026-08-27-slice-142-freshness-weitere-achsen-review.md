# Review slice-142 — Zwei ungewachte Pin-Klassen

**Gegenstand:** [slice-142](../plan/planning/done/slice-142-freshness-weitere-achsen.md), Stand vor der Nacharbeit.
**Datum:** 2026-08-27. **Reviewer:** unabhängiger Subagent.

**Nachtrag zur Form:** dieser Report ist **nach** dem Lauf aus dem
Sitzungsprotokoll übertragen worden, nicht währenddessen geschrieben. Die
Befund-Nummern und das Urteil sind die des Laufs; die Formulierungen sind
gekürzt.

---

## Urteil

**Nicht schließbar in der vorliegenden Form — schließbar nach Nacharbeit.**

> Blockierend sind nicht die Kanten des Codes, sondern **zwei
> Entscheidungs-Aussagen, die die Messung nicht trägt**. Zum Schließen genügt
> es, diese beiden Stellen ehrlich zu machen (oder die zwei Achsen zu bauen —
> beides wäre je fünf Zeilen).

## Befunde

| # | Rang | Befund |
|---|---|---|
| F-1 | **HIGH** | Die Action-Achse wurde mit zwei Sätzen abgelehnt, die **beide widerlegbar** sind. Sie braucht *keine* neue Quellen-Form — der Tag-Kommentar neben dem SHA ist genau die Größe, die der `releases/latest`-Zweig vergleicht; gegengeprobt mit `PINNED=v6.0.2` ⇒ `VERALTET — upstream v7.0.1`, und `docker/login-action` v4.2.0 ⇒ v4.6.0. Und der alte SHA **ist** ein Sicherheitsproblem: `actions/checkout` v6.1.0 existiert eigens, damit v6-Pinner den Fork-PR-Fix bekommen |
| F-2 | **HIGH** | Die Zusage, die zwei übrigen Basis-Images seien „über die Achsen darüber gewacht", ist **falsch**. `make freshness-go` meldet `ok`, während `golang:1.27.0` `sha256:0ecdc2a9…` trägt gegen unseren Pin `sha256:65b6f280…` — seit fünf Tagen |
| F-3 | MEDIUM | `docker buildx imagetools inspect` kennt keine Zeitgrenze; die Fail-open-Zusage des Kopfes gilt für jeden Zweig |
| F-4 | MEDIUM | Nur die Upstream-Seite wird auf ihre Form geprüft. Der Pin kommt aus einer Textextraktion am `Dockerfile` und ist die fragilere Hälfte |
| F-5 | MEDIUM | Die Zählmethode misst **unsere eigenen Bumps**, nicht die Bewegung upstream — und trägt trotzdem die Entscheidung |
| F-6 | MEDIUM | `semgrep`- und `a-check`-Image fehlen in der Klasse ganz; letzteres trägt gar keinen Tag |
| F-7 | MEDIUM | Der Workflow-Kopf sagt „Kein Docker" und „Drei Achsen" — beides nach dem Einbau falsch |
| F-8 | MEDIUM | Die DoD-Haken waren **vor** dem Review gesetzt |
| F-9 | MEDIUM | Die Quellen-Begründung behauptet Unmöglichkeit des Handbetriebs; ein `curl` auf gcr.io liefert `docker-content-digest` ohne Token |
| F-10 | MEDIUM | Der Risiko-Ausgang „entfallen" steht, während der Nachtlauf jede Nacht rot wird — ohne benannten Adressaten |
| F-11 | LOW | `RUNTIME_BASE` dupliziert die `FROM`-Zeile |
| F-12 | LOW | „VERALTET" behauptet eine Ordnung, die Digests nicht haben |

## Erledigung

Alle zwölf Befunde sind eingearbeitet (`e4471bf`, `6c1d97a`, `2485f42`):

- **F-1/F-2** durch **Bauen** statt Umformulieren — zwei Action-Achsen und drei
  Digest-Achsen kamen dazu, der Nachtlauf trägt jetzt zwölf.
- **F-5** als eigene Register-Klasse [`BEO-020`](../plan/planning/observations.md).
- **F-6** samt der strukturellen Ursache: `A_CHECK_VERSION` steht jetzt als
  Variable statt als Prosa im Kommentar.
- **F-10** als eingetretener Risiko-Ausgang mit Folge-Slice
  [slice-161](../plan/planning/open/slice-161-sechs-pins-heben.md).
- **F-11** entfiel mit dem Muster: Referenz und Digest kommen jetzt aus
  **einer** `FROM`-Zeile.
