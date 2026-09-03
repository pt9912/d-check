**Vorgang:** slice-195
**Fund:** Die Beobachtungs-Register-Datenmigration überschritt die im
eigenen §5 benannte Ein-Sitzungs-Grenze und den dafür in §6 vorab
festgelegten Rückführungs-Trigger, wurde aber bewusst nicht nach `next`
zurückgeführt, weil eine Aufteilung den Zähler-Diff-Beleg der Migration
selbst zerrissen hätte — die Review-Last wurde stattdessen über einen
gezielten Stichproben-Fokus statt vollständiger Zeile-für-Zeile-Prüfung
aufgefangen.
