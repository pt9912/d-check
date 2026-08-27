# extract-command.awk — zieht tool_input.command aus der PreToolUse-Hook-JSON.
#
# ZUSAGE: Stdout ist der dekodierte Befehl (ggf. leer). Exit 0 = gelesen,
# Exit 3 = Parse-Zweifel. Der Aufrufer blockt bei 3 (fail-closed).
#
# WARUM EIN SCANNER UND KEIN GRIFF: ein Kommando ist Freitext und kann die
# Zeichenkette "command" selbst enthalten. Ein Regex-Griff naehme den Treffer
# IM Wert und der Waechter entschiede ueber die falsche Groesse. Darum
# zeichenweise, mit Tiefen- und Key-Stack: nur der Pfad tool_input -> command
# auf Objekt-Tiefe 2 zaehlt.
#
# GRENZE: \u-Escapes im Befehl gelten als Zweifel und blocken. Sie zu dekodieren
# hiesse, den Waechter zu einer Sandbox auszubauen; er ist ein Stolperdraht
# (AGENTS.md §3.1).
#
# POSIX-awk (busybox/gawk/BSD), kein gawk-Spezifikum.

{ doc = (NR == 1) ? $0 : doc "\n" $0 }

END {
  n = length(doc)
  depth = 0       # Verschachtelungstiefe ({}/[])
  instr = 0       # in einem JSON-String?
  esc = 0         # letztes Zeichen war Backslash?
  buf = ""        # aktueller String-Inhalt (dekodiert)
  hadu = 0        # aktueller String enthielt \uXXXX
  sawobj = 0      # je ein Top-Level-Objekt gesehen?
  found = 0
  cmdval = ""

  for (i = 1; i <= n; i++) {
    c = substr(doc, i, 1)

    if (instr) {
      if (esc) {
        esc = 0
        if (c == "\"") buf = buf "\""
        else if (c == "\\") buf = buf "\\"
        else if (c == "/") buf = buf "/"
        else if (c == "n") buf = buf "\n"
        else if (c == "t") buf = buf "\t"
        else if (c == "r") buf = buf "\r"
        else if (c == "b") buf = buf sprintf("%c", 8)
        else if (c == "f") buf = buf sprintf("%c", 12)
        else if (c == "u") {
          # \u verlangt GENAU vier Hex-Stellen. Sonst ist das JSON malformed und
          # ein i+=4 wuerde ueber ein schliessendes " hinweg desynchronisieren --
          # der Waechter ginge dann fail-OPEN.
          if (substr(doc, i + 1, 4) !~ /^[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]$/) exit 3
          hadu = 1; i = i + 4
        }
        else buf = buf c                    # unbekannter Escape: Zeichen behalten
        continue
      }
      if (c == "\\") { esc = 1; continue }
      if (c == "\"") {
        # Stringende: Key oder Value?
        if (depth > 0 && ctype[depth] == "o" && wantkey[depth] == 1) {
          curkey[depth] = buf
        } else if (depth >= 2 && ctype[depth] == "o" && curkey[depth] == "command" &&
                   ctype[depth - 1] == "o" && curkey[depth - 1] == "tool_input") {
          if (hadu) exit 3
          found = 1
          cmdval = buf
        }
        instr = 0
        continue
      }
      buf = buf c
      continue
    }

    # ausserhalb eines Strings
    if (c == "\"") { instr = 1; buf = ""; hadu = 0; continue }
    if (c == "{") { depth++; sawobj = 1; ctype[depth] = "o"; wantkey[depth] = 1; curkey[depth] = ""; continue }
    if (c == "}") { if (depth > 0) depth--; continue }
    if (c == "[") { depth++; ctype[depth] = "a"; continue }
    if (c == "]") { if (depth > 0) depth--; continue }
    if (c == ":") { if (depth > 0 && ctype[depth] == "o") wantkey[depth] = 0; continue }
    if (c == ",") { if (depth > 0 && ctype[depth] == "o") wantkey[depth] = 1; continue }

    # Ausserhalb eines Strings sind nur Struktur-Zeichen, Whitespace und die
    # Zeichen von Zahlen/Literalen zulaessig. Alles andere ist malformes JSON
    # -- und ohne diese Pruefung liefe der Scanner darueber hinweg und lieferte
    # einen ABGESCHNITTENEN Befehl: `"cmd":"bash -lc "go build""` gaebe nur
    # `bash -lc ` zurueck, und der Waechter urteilte ueber die halbe Eingabe.
    if (c !~ /^[ \t\r\neEtrufalsn0-9.+-]$/) exit 3
  }

  if (!sawobj) exit 3                      # kein Objekt -> kein/kaputtes JSON
  if (instr == 1 || depth != 0) exit 3     # abgeschnitten/unbalanciert
  if (found) printf "%s", cmdval
  exit 0
}
