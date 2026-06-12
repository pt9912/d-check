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
