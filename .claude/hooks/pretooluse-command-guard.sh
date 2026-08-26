#!/usr/bin/env bash
# pretooluse-command-guard — verbietet Host-Paketmanager, die Host-Go-Toolchain
# und Host-Skript-Interpreter (d-check ist make/Docker-only, AGENTS.md §3.1,
# ADR-0001).
#
# Geprüft wird die Befehlsposition jedes Kommando-Segments (Trennung
# an ; && || | $( ` ( und Zeilenenden) — `git commit -m "docs: pip"`
# oder `docker run img npm test` bleiben erlaubt; `/usr/bin/pip` und
# `sudo pip` werden erkannt. Sub-Shell-Strings (`bash -c "…"`, auch in
# Flag-Bündeln wie `-lc`/`-ec`/`-cx`) werden rekursiv geprüft (MR-005).
#
# SKRIPT-INTERPRETER: `python`/`python3`/`python3.x`, `perl` und `ruby` in
# Befehlsposition sind blockiert. Grund ist nicht die Sprache, sondern der
# Skopus: §3.1 nennt als Host-Klasse `git`, GNU `make`, `bash`, Docker und die
# POSIX-Standardwerkzeuge, die die Gate-Skripte rufen. Ein General-Purpose-
# Interpreter gehört nicht dazu, und er hebelt die Klasse aus — er kann alles,
# was die genannten Werkzeuge können, ohne deren Grenze zu erben.
#
# `node` steht MIT auf der Liste. Der Agent ruft es nicht mehr auf; der Hook
# ruft es weiterhin, denn Hook-Subprozesse laufen nicht durch den Hook.
#
# DAS IST EINE OFFENE INKONSISTENZ, KEINE LÖSUNG: dieser Guard soll erzwingen,
# dass nur die in AGENTS.md §3.1 genannte Host-Klasse benutzt wird — und ist
# selbst in einer Toolchain geschrieben, die nicht dazugehört. Die Ablösung auf
# bash + POSIX-Werkzeuge ist geschnitten; bis dahin bleibt node eine deklarierte
# Host-Abhängigkeit dieses einen Skripts.
#
# BENANNTE GRENZE: `find -exec`, `awk`-Programme und jeder Interpreter, den
# diese Liste nicht kennt, bleiben ungeprüft — der Guard ist ein Stolperdraht
# gegen versehentliche Host-Toolchain-Nutzung, keine Sandbox.
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
    "npm","pnpm","yarn","npx","corepack","cargo","rustup","gem","conda",
    "go","gofmt","golangci-lint","staticcheck", // Host-Go: ADR-0001 + AGENTS §3.1
    "perl","ruby","node","deno","bun"]);        // Skript-Interpreter: AGENTS §3.1
  // python, python3, python3.12 … — die Versions-Suffixe machen eine Liste
  // unvollständig, darum als Muster.
  const BLOCKED_RE = /^python[0-9]*(\.[0-9]+)*$/;
  const PREFIXES = new Set(["sudo","env","command","exec","nice","time",
    "xargs","eval"]);
  const SHELLS = new Set(["bash","sh","zsh","dash","ksh"]);
  const stripQuotes = t => t.replace(/^["'\'']+|["'\'']+$/g, "");

  function scan(cmd, depth) {
    if (depth > 3) return true; // zu tief verschachtelt → fail-closed
    // Auch einfaches & trennt (Hintergrund-Start): `echo x & pip install y`.
    const segments = cmd.split(/(?:;|&&|&|\|\||\||\$\(|`|\(|\r?\n)/);
    for (const seg of segments) {
      const tokens = seg.trim().split(/\s+/).filter(Boolean).map(stripQuotes);
      let i = 0;
      // Brace-Group-Delimiter mit ueberspringen: bei `{ go build; }` waere der
      // Kopf sonst `{` und das Werkzeug an Position 2 entkaeme der Pruefung.
      while (i < tokens.length &&
             (/^[A-Za-z_][A-Za-z0-9_]*=/.test(tokens[i]) || PREFIXES.has(tokens[i]) ||
              tokens[i] === "{" || tokens[i] === "}")) i++;
      if (i >= tokens.length) continue;
      const head = tokens[i].replace(/^.*\//, ""); // /usr/bin/pip → pip
      if (BLOCKED.has(head) || BLOCKED_RE.test(head)) return true;
      if (SHELLS.has(head)) {
        // -c auch in Flag-Bündeln erkennen (-lc, -ec, -cx, …): bei
        // sh/bash ist c das einzige Single-Letter-Flag mit
        // Kommando-String-Semantik, das Bündel ist also eindeutig.
        const cIdx = tokens.findIndex((t, k) => k > i && /^-[a-z]*c[a-z]*$/.test(t));
        if (cIdx !== -1 && cIdx + 1 < tokens.length &&
            scan(tokens.slice(cIdx + 1).join(" "), depth + 1)) return true;
      }
    }
    return false;
  }

  let s = "";
  process.stdin.on("data", d => s += d);
  process.stdin.on("end", () => {
    let cmd = "";
    try {
      const j = JSON.parse(s);
      cmd = String((j.tool_input && j.tool_input.command) || "");
    } catch { process.stdout.write("block"); return; } // unlesbar → fail-closed
    process.stdout.write(scan(cmd, 0) ? "block" : "ok");
  });
')"

if [ "$verdict" = "block" ]; then
  cat <<'JSON'
{
  "decision": "block",
  "reason": "d-check is make/Docker-only (AGENTS.md §3.1, ADR-0001). Use make targets and the POSIX host tools the gate scripts use (grep/sed/awk/find); do not run host package managers, host go, or host script interpreters (apt/brew/pip/npm/cargo/go/python/perl/ruby/node)."
}
JSON
fi
# Pass-Fall: keine Ausgabe — normale Permission-Prüfung übernimmt.
