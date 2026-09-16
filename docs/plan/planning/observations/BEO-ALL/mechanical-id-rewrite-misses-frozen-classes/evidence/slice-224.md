**Vorgang:** slice-224
**Fund:** **Die Frozen-Liste war diesmal korrekt und stand vor der ersten
Ersetzung — und traf trotzdem eine Stelle nicht, die keine ihrer fünf Klassen
vorsah.** Eingefroren waren acht Dateien (Lauf-Belege), acht wie erwartet.
Ein pauschaler `sed` über die **lebende** Datei `harness/conventions.md`
traf aber zusätzlich die **Tabellenzeile**, mit der diese Datei selbst
[`MR-071`](../../../../../../../harness/conventions.md#mr-071) (die Hebung
auf den alten Pin) im Adaptions-Index führt — dieselbe Über-Hebungs-Klasse
wie bei [`MR-067`](../../../../../../../harness/conventions.md#mr-067) in
slice-222, nur eine Ebene tiefer: nicht die *Datei selbst* (die stand korrekt
auf der Frozen-Liste), sondern eine *Zeile* in einer anderen, unbestreitbar
lebenden Datei, die über diese Datei berichtet.

**Der Fehler war diesmal schnell sichtbar**, nicht erst im Review: Der nächste
Blick auf `harness/conventions.md` zeigte die Zeile mit dem falschen Pfad,
noch bevor der nächste Schritt begann. Behoben im selben Arbeitsschritt.

**Bestätigt, zum fünften Mal an derselben Stelle:** [`MR-070`](../../../../../../../harness/conventions.md#mr-070)
deckt Datei-Eigenschaften, keine Tabellenzeilen innerhalb einer Datei, die
selbst einen Adaptions-Index trägt. Die in slice-222 benannte Grenze („eine
Frozen-Liste über Datei-Eigenschaften kann das prinzipiell nicht fangen")
gilt unverändert — hier nur mit einem glücklicheren Fundzeitpunkt.
