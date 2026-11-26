# `codepaths` prüft einen Pfad nur, wenn der Inline-Code-Span auf einer Zeile steht

**Sub-Area:** `*`

Das Modul `codepaths` urteilt zeilenlokal: Steht ein nicht existierender Pfad in
einem Inline-Code-Span auf **einer** Zeile, meldet es `codepath-missing`; derselbe
Pfad in einem Span, der über einen Zeilenumbruch läuft (auf der Öffnungs- wie auf
der Folgezeile), meldet **nichts**. Das grüne `make doc-check` sagt dann über diese
Pfade nichts, und niemand sieht es: Das Modul `spans` führt den mehrzeiligen Span
als wohlgeformt ([`DC-FA-SPAN-001`](../../../../../../spec/lastenheft.md#dc-fa-span-001--markdown-span-artefakte-modul-spans-opt-in) nennt ihn ausdrücklich), und weder
[`DC-FA-CODE-001`](../../../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) noch die Spezifikation tragen einen Satz zur Zeilengrenze eines
Pfad-Spans.

**Gemessen:** im Wegwerf-Klon mit dem gebauten Image, ein erfundener Pfad einzeilig,
mehrzeilig auf der Öffnungszeile, mehrzeilig auf der Folgezeile — ein Befund, keiner,
keiner. Der Bump-Eintrag eines Slice trug acht solche Spans; ihre Pfade waren von
Hand gelesen und existieren.

**Unterschied zu
[`zeilenblinder-matcher-verschmilzt-unabhaengige-fundstellen`](../zeilenblinder-matcher-verschmilzt-unabhaengige-fundstellen/observation.md):**
Dort sieht ein Matcher **zu viel**, weil er Zeilen verschmilzt; hier sieht der Matcher
**zu wenig**, weil er sie nicht verbindet. Beide teilen die Wurzel, dass die Zeile
die Einheit des Moduls ist und der Span keine ist.

**Verwandt:**
[`module-promise-only-on-scan-axis`](../module-promise-only-on-scan-axis/observation.md)
— dort liest ein Modul Eingaben außerhalb seiner Scan-Menge; hier liegt der Pfad
**innerhalb** der Scan-Menge, nur in einer Form, die das Modul nicht verbindet — und
[`prosa-verschwindet-in-inline-code-spanne`](../prosa-verschwindet-in-inline-code-spanne/observation.md),
dessen Lexik Spans absatzweit paart.
