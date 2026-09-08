**Vorgang:** slice-216
**Fund:** **Eine Äquivalenz behauptet, wo eine Deckung gemessen war.** Der
Slice maß drei Konfigurationen über einen Mini-Bestand und schrieb daraus,
`refs: ["**"]` erreiche die Wirkung von `exempt-paths` *„exakt"* und *„und
mehr"*. Gemessen war, dass **beide toten Verweise einer Datei** stumm werden
— zwei Verweise, eine Auflösungs-Klasse.

**Die N+1-te Form widerlegt es, und sie steht im Produkt selbst.** Die
**Symlink**-Ablehnung überlebt jedes `refs`-Glob; der Kommentar über dem
Ventil in `internal/hexagon/core/rules/links.go` sagt es wörtlich — die
Symlink-Prüfung habe **Vorrang** und bleibe unberührt. Die beiden Ventile
stehen damit **quer** zueinander, nicht in Dominanz: `exempt-paths` nähme die
Datei aus dem ganzen Modul, `ignore-refs` nimmt Auflösungs-Klassen eines Ziels
heraus.

**Die Behauptung trug den Entscheid.** Sie stand nicht als Beiwerk daneben,
sondern war das Argument gegen den CR — *„das Werkzeug löst die Aufgabe
bereits, und schärfer"*. Der Entscheid kippt dadurch nicht, aber seine
Reichweite: Er gilt dem **beschriebenen Fall**, nicht dem allgemeinen.

**Dasselbe Muster eine Ebene tiefer:** Gemessen wurde mit **einem** Modul
(`links`), entschieden wurde über **zwei** (`links` und `anchors`). Erst die
Nachmessung mit `--enable anchors` zeigte, dass das Rezept dort ebenfalls
wirkt — und nebenbei, dass für einen entfernten Baum gar kein
`anchor-missing` entsteht, sondern `target-missing`.
