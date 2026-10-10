**Vorgang:** slice-273
**Fund:** Ein Fremd-Scanner meldete am publizierten Image `v0.86.0` vier
behebbare CVEs in `golang.org/x/net` v0.59.0 — CVE-2026-78663 mit 9.1
(CRITICAL), CVE-2026-78659, -78660, -78669 mit je 7.5 (HIGH). Trivy 0.74.0
führte mit der Schwachstellen-Datenbank vom 2026-10-09 alle vier als
`UNKNOWN`, `make image-scan` endete mit Exit 0. Mit der Datenbank vom
2026-10-10 standen zwei auf `UNKNOWN`, CVE-2026-78663 auf MEDIUM und
CVE-2026-78669 auf HIGH — der Rückstand betrug rund einen Tag (Verifikation
I1, Review F-2).
