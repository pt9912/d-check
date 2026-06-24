# d-check — Release-Register

> Kanonischer, **auflösender** Link-Ziel-Ort für Erwähnungen **eigener**
> d-check-Releases — etwa [die jeweils aktuelle Version](#aktuell). **Nur die aktuelle Version**
> trägt einen expliziten HTML-Anker `#vX.Y.Z` (wörtlich, mit Punkten — der
> Heading-/Tabellen-Slug verschluckt sie). Beim Release **wandert** der Anker zur
> neuen aktuellen Version; die bisherige Zeile verliert ihn — dadurch *bricht*
> jeder feste Markdown-Link-Pin auf eine veraltete Version (`anchor-missing`),
> und ein vergessener Bump fällt auf (das ist der Zweck dieses Registers).
>
> **Kein Duplikat** der Detail-Changes — die stehen im
> [CHANGELOG](CHANGELOG.md). Hier nur Versions-Koordinaten (Datum, Tag).
> **Fremde** Versionen (Kurs-Baseline, semgrep, Go-Toolchain) gehören
> nicht hierher, sondern verlinken auf ihre eigene Quelle.
>
> Quelle der Liste: Git-Tags + CHANGELOG-Einträge (1:1). Dieses Register
> ist **noch nicht** als `ids`-Linkpflicht scharfgeschaltet — reines,
> bereitgestelltes Link-Ziel (siehe `idee-version-link-pflicht`).

## Aktuell

Aktuelle Version: [`v0.27.0`](#v0.27.0) — 2026-06-23.

Aus anderen Dokumenten stabil referenzierbar als `version.md#aktuell`
(zeigt immer hierher, nie auf eine feste Nummer). Pro Release sind genau
diese Zeile **und** eine neue Tabellen-Zeile im Verlauf nachzuziehen **und der
`<a id>`-Anker auf die neue Version zu verschieben** (die bisherige Zeile verliert
ihn, sonst bleiben veraltete feste Pins auflösbar und die Anker-Kaskade meldet den
vergessenen Bump nicht) — der einzige Bump-Punkt, analog zum README-Versions-Pin.

## Verlauf

| Version                        | Datum      | Release                                                               |
| ------------------------------ | ---------- | --------------------------------------------------------------------- |
| `v0.27.0` <a id="v0.27.0"></a> | 2026-06-23 | [Tag v0.27.0](https://github.com/pt9912/d-check/releases/tag/v0.27.0) |
| `v0.26.0`                      | 2026-06-23 | [Tag v0.26.0](https://github.com/pt9912/d-check/releases/tag/v0.26.0) |
| `v0.25.0`                      | 2026-06-23 | [Tag v0.25.0](https://github.com/pt9912/d-check/releases/tag/v0.25.0) |
| `v0.24.0`                      | 2026-06-23 | [Tag v0.24.0](https://github.com/pt9912/d-check/releases/tag/v0.24.0) |
| `v0.23.0`                      | 2026-06-22 | [Tag v0.23.0](https://github.com/pt9912/d-check/releases/tag/v0.23.0) |
| `v0.22.0`                      | 2026-06-22 | [Tag v0.22.0](https://github.com/pt9912/d-check/releases/tag/v0.22.0) |
| `v0.19.0`                      | 2026-06-20 | [Tag v0.19.0](https://github.com/pt9912/d-check/releases/tag/v0.19.0) |
| `v0.18.0`                      | 2026-06-19 | [Tag v0.18.0](https://github.com/pt9912/d-check/releases/tag/v0.18.0) |
| `v0.17.0`                      | 2026-06-19 | [Tag v0.17.0](https://github.com/pt9912/d-check/releases/tag/v0.17.0) |
| `v0.12.0`                      | 2026-06-18 | [Tag v0.12.0](https://github.com/pt9912/d-check/releases/tag/v0.12.0) |
| `v0.11.0`                      | 2026-06-17 | [Tag v0.11.0](https://github.com/pt9912/d-check/releases/tag/v0.11.0) |
| `v0.10.0`                      | 2026-06-16 | [Tag v0.10.0](https://github.com/pt9912/d-check/releases/tag/v0.10.0) |
| `v0.9.0`                       | 2026-06-15 | [Tag v0.9.0](https://github.com/pt9912/d-check/releases/tag/v0.9.0)   |
| `v0.8.0`                       | 2026-06-13 | [Tag v0.8.0](https://github.com/pt9912/d-check/releases/tag/v0.8.0)   |
| `v0.7.0`                       | 2026-06-13 | [Tag v0.7.0](https://github.com/pt9912/d-check/releases/tag/v0.7.0)   |
| `v0.6.0`                       | 2026-06-13 | [Tag v0.6.0](https://github.com/pt9912/d-check/releases/tag/v0.6.0)   |
| `v0.5.0`                       | 2026-06-13 | [Tag v0.5.0](https://github.com/pt9912/d-check/releases/tag/v0.5.0)   |
| `v0.4.0`                       | 2026-06-13 | [Tag v0.4.0](https://github.com/pt9912/d-check/releases/tag/v0.4.0)   |
| `v0.3.0`                       | 2026-06-12 | [Tag v0.3.0](https://github.com/pt9912/d-check/releases/tag/v0.3.0)   |
| `v0.2.1`                       | 2026-06-12 | [Tag v0.2.1](https://github.com/pt9912/d-check/releases/tag/v0.2.1)   |
| `v0.2.0`                       | 2026-06-12 | [Tag v0.2.0](https://github.com/pt9912/d-check/releases/tag/v0.2.0)   |
| `v0.1.0`                       | 2026-06-11 | [Tag v0.1.0](https://github.com/pt9912/d-check/releases/tag/v0.1.0)   |
