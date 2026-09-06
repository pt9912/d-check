**Vorgang:** slice-205
**Fund:** Eine Inventur-Berichtigung zählte Anforderungen per
`grep -cE '^#{3,4} [A-Z][A-Z0-9]*-'` — das trifft jede Überschrift aus
Großbuchstaben und Bindestrich, nicht nur Kennungen. In einem Fremd-Repo waren
alle vier „Treffer" deutsche Komposita (`MVP-Abgrenzung`, `CLI-Ziel`,
`CORS- und CSP-Grundregeln`). Aus deren Präfixen entstand die Behauptung, das
Repo führe eine eigene ID-Familie `CLI`/`CORS`/`MVP` — ein **fünftes
ID-Schema, das es nicht gibt**. Die Korrektur war damit schlimmer als der
Fehler, den sie behob: die erste Fassung hatte falsche Zahlen, die zweite eine
erfundene Kategorie. Mit einer Form-Prüfung (Präfix **plus Ziffernblock**) sind
es vier Schemata, und das Repo führt auf dieser Ebene keine.
