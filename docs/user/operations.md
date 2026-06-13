# Operations — d-check Aufruf-Referenz

Kurzreferenz für den Betrieb; verbindliche Verträge stehen im
[Lastenheft](../../spec/lastenheft.md) (Anforderungen) und in der
[Spezifikation](../../spec/spezifikation.md) (Schemata, Defaults,
Grund-Codes) — diese Seite dupliziert sie nicht.

## Aufruf

```sh
docker run --rm -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check@sha256:<digest>
```

Das Image prüft das nach `/repo` gemountete Repository (read-only
genügt — das Tool schreibt nie,
[`DC-QA-03`](../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
CLI-Optionen werden als Container-Argumente angehängt.

## Optionen

| Option | Wirkung |
|---|---|
| `--enable <modul>` / `--disable <modul>` | Regelmodule zu-/abschalten (`links`, `anchors`, `ids`, `matrix`, `external`, `codepaths`, `spans`, `hostpaths`); CLI schlägt Konfiguration ([`DC-FA-CLI-002`](../../spec/lastenheft.md#dc-fa-cli-002--regelmodul-auswahl)) |
| `--json` | maschinenlesbare Gesamt-Ausgabe ([Schema](../../spec/spezifikation.md)) |

Default-Module ohne Konfiguration: `links` + `anchors`. Das Modul
`external` ist strikt opt-in (einzige Netzwerk-Tür).

## Exit-Codes

`0` = keine Befunde, `1` = mindestens ein Befund, `2` = Nutzungs-
oder Umgebungsfehler
([`DC-FA-CLI-003`](../../spec/lastenheft.md#dc-fa-cli-003--exit-codes)).

## Konfiguration

Optionale `.d-check.yml` in der Repo-Wurzel; Schema und Defaults in
der [Spezifikation §`.d-check.yml`](../../spec/spezifikation.md#d-checkyml),
Vollvalidierung gemäß
[`DC-FA-CONF-001`](../../spec/lastenheft.md#dc-fa-conf-001--konfigurationsdatei)
(jeder Konfigurationsfehler → Exit 2). Ein Beispiel im Vollausbau ist
die [Selbstkonfiguration dieses Repos](../../.d-check.yml).

## Linkdichte erzwingen — `link-policy: always`

Das Modul `ids` prüft per Default nur, dass *nackte* Kennungen im
Fließtext verlinkt sind; eine Kennung in Inline-Code (`` `…` ``) bleibt
frei. Wer **gut verlinkte Dokumente** als gemessenes Property will
(jede Kennung eine navigierbare Referenz), setzt pro Muster
`link-policy: always` — dann ist auch eine Kennung in Inline-Code
linkpflichtig
([`DC-FA-ID-001`](../../spec/lastenheft.md#dc-fa-id-001--linkpflicht-für-kennungen-modul-ids)):

```yaml
ids:
  patterns:
    - regex: 'ADR-\d{4}'
      target: docs/plan/adr/
      link-policy: always            # prose (Default) | always
      exempt-paths: [CHANGELOG.md, "docs/reviews/**"]
```

Zwei Ventile halten `always` treffsicher:

- **`exempt-paths`** (Glob-Liste, Syntax wie `scan.ignore`): Dateien,
  in denen die strenge Regel nicht gilt — typischerweise literal-schwere
  Artefakte wie Changelogs oder Review-Reports.
- **`d-check:ignore`** (HTML-Kommentar auf der Zeile, Begründung
  empfohlen): nimmt eine Zeile aus — für bewusst illustrative
  Beispiel-Kennungen. Der Marker wirkt auf `ids` und `codepaths`.

**Bewusst opt-in:** Der Default bleibt `prose`, damit bestehende
Konfigurationen byte-identisch laufen und kein Repo ungefragt rote
Läufe bekommt. Das heißt **nicht**, dass ungenügende Verlinkung
unsichtbar bleibt: Sie zu *entdecken* ist eine Bringschuld des
Betreibers — ein `always`-Lauf über die geprüften Repos zeigt die
Lücken, unabhängig davon, ob ein Repo `always` schon aktiviert hat
([`DC-QA-04`](../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Muster).
