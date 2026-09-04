# Ein frischer CVE bleibt severity-los, während ein anderer Vendor ihn schon einstuft

**Sub-Area:** `*`

`make image-scan` macht nur bei CRITICAL/HIGH **mit verfügbarem Fix** rot
(ADR-0066). Ein frisch veröffentlichter CVE trägt in Trivys Vuln-DB zunächst
Severity `UNKNOWN` (die eigene Warnung: „Using severities from other vendors
for some vulnerabilities"), während ein anderer Vendor zum selben Zeitpunkt
bereits eine Severity führt. Der Fund selbst steht im Vollbericht — nur das
Entscheidungsfilter sieht ihn nicht, solange die Klassifikation aussteht.
