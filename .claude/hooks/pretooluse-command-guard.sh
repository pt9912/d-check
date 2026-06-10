#!/usr/bin/env bash
# pretooluse-command-guard — verbietet Host-Paketmanager und
# Toolchain-Installationen (d-check ist make/Docker-only, AGENTS.md
# §3.1; Implementierungssprache via ADR-0001 — Liste danach um die
# sprachkonkreten Build-/Test-Tools erweitern).
#
# Geprüft wird nur die Befehlsposition jedes Kommando-Segments
# (Trennung an ; && || | $( ` ( und Zeilenenden), nicht beliebige
# Argumente — `git commit -m "docs: pip"` oder
# `docker run img npm test` bleiben erlaubt; `/usr/bin/pip` und
# `sudo pip` werden erkannt.
#
# Im Pass-Fall: KEINE Ausgabe — "approve" würde das Permission-System
# überspringen; ohne Ausgabe läuft die normale Permission-Entscheidung.
set -euo pipefail

# Fail-closed: ohne node keine Prüfung möglich → Tool-Call blockieren.
if ! command -v node >/dev/null 2>&1; then
  echo "pretooluse-command-guard: node not found on host — blocking (fail-closed)." >&2
  exit 2
fi

input="$(cat)"

verdict="$(printf '%s' "$input" | node -e '
  const BLOCKED = new Set(["apt","apt-get","brew","pip","pip3","pipx",
    "npm","pnpm","yarn","npx","corepack","cargo","rustup","gem","conda"]);
  const PREFIXES = new Set(["sudo","env","command","exec","nice","time","xargs"]);
  let s = "";
  process.stdin.on("data", d => s += d);
  process.stdin.on("end", () => {
    let cmd = "";
    try {
      const j = JSON.parse(s);
      cmd = String((j.tool_input && j.tool_input.command) || "");
    } catch { process.stdout.write("block"); return; } // unlesbar → fail-closed
    const segments = cmd.split(/(?:;|&&|\|\||\||\$\(|`|\(|\r?\n)/);
    for (const seg of segments) {
      const tokens = seg.trim().split(/\s+/).filter(Boolean);
      let i = 0;
      while (i < tokens.length &&
             (/^[A-Za-z_][A-Za-z0-9_]*=/.test(tokens[i]) || PREFIXES.has(tokens[i]))) i++;
      if (i >= tokens.length) continue;
      const head = tokens[i].replace(/^.*\//, ""); // /usr/bin/pip → pip
      if (BLOCKED.has(head)) { process.stdout.write("block"); return; }
    }
    process.stdout.write("ok");
  });
')"

if [ "$verdict" = "block" ]; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "d-check is make/Docker-only (AGENTS.md §3.1). Use make targets; do not install or run host package managers (apt/brew/pip/npm/cargo/...)."
}
JSON
fi
# Pass-Fall: keine Ausgabe — normale Permission-Prüfung übernimmt.
