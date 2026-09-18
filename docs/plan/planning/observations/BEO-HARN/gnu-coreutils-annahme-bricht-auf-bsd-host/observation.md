# GNU-Coreutils-Annahmen in `tools/harness/*.sh` brechen unbemerkt auf BSD/macOS-Hosts

**Sub-Area:** `tools/harness/`

Zwei unabhängige Stellen in `fetch-baseline-cache.sh` nahmen stillschweigend
GNU-Coreutils-Verhalten an: ein Datei-Anzahl-Vergleich über `wc -l` (GNU
gibt bei einzelner Eingabe eine ungepaddete Zahl aus, BSD rechtsbündig mit
führenden Leerzeichen — der String-Vergleich `=` brach nur auf BSD) und
`readlink -e` (ein reines GNU-Flag, das BSD-`readlink` mit „illegal option"
ablehnt). Beide Stellen liefen über Jahre auf Linux-CI grün und fielen erst
beim ersten realen `make gates`-Lauf auf einem macOS-Host auf — kein Gate
prüft die Host-Werkzeug-Klasse gegen BSD, weil lokale Läufe bislang nur auf
Linux/CI gefahren wurden.
