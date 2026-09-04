**Vorgang:** slice-201
**Fund:** Trivy stufte `CVE-2026-56855`/`CVE-2026-78662` (`golang.org/x/crypto`) im Nachtlauf vom 2026-09-04 als `UNKNOWN` ein und ließ `make image-scan` grün — Docker Scout führte zum selben Zeitpunkt bereits `7.5 HIGH`, obwohl der Fix (`v0.56.0`) längst existierte.
