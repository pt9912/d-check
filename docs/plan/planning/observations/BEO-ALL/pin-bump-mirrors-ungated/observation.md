# Die Spiegel einer Pin-Hebung sind vier Klassen, gehoben wird nur die grep-bare

**Sub-Area:** `*`

`baseline/<tag>`-Pfad-Verweise sind gate-gedeckt und werden beim Bump
retargetet; Release-/Tree-**URLs** mit dem Tag, **Prosa-/Ellipsen-Pins**
und der **zitierende Verweis** (ein Verweis zitiert den Wortlaut einer
Datei, deren Pfad mitgehoben wird, während das Zitat daneben am neuen Ziel
nicht mehr existiert) deckt kein Gate und kein Pfad-`grep`. Gate-blind in
beide Richtungen: vergessene Hebung ebenso wie Über-Hebung (eine
Vergangenheits-Aussage wird mitgehoben). Mechanische Form seit slice-122
baubar (`versions.patterns`), im eigenen Profil bewusst nicht scharfgeschaltet.
