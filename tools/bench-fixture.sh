#!/usr/bin/env bash
# bench-fixture.sh — DC-QA-01-Benchmark (spec/spezifikation.md
# §DC-QA-01.a): generiert ein deterministisches Fixture-Repo
# (1.000 Markdown-Dateien, ≤ 20 MB, definierter Link-/Heading-Mix),
# misst N=3 Läufe der Default-Module im Runtime-Container
# (read-only, netzlos) und prüft den Median gegen < 5 s.
set -euo pipefail

cd "$(dirname "$0")/.."

FIXTURE="${BENCH_DIR:-/tmp/d-check-bench-fixture}"
IMAGE="${IMAGE:-d-check}"
RUNS=3

rm -rf "$FIXTURE"
mkdir -p "$FIXTURE/docs"

filler="$(printf 'x%.0s' {1..400})"
for i in $(seq -w 1 1000); do
  next=$(printf 'f%04d.md' $(((10#$i % 1000) + 1)))
  {
    echo "# Datei $i"
    for h in 1 2 3 4 5 6 7 8 9 10; do
      echo
      echo "## Abschnitt $h"
      echo
      echo "Fließtext zu Abschnitt $h in Datei $i: [Querverweis]($next) und [Anker]($next#abschnitt-$h)."
      echo "Fülltext $filler"
    done
  } > "$FIXTURE/docs/f$i.md"
done

size_mb=$(du -sm "$FIXTURE" | cut -f1)
echo "bench: Fixture generiert — 1000 Dateien, ${size_mb} MB"
if [ "$size_mb" -gt 20 ]; then
  echo "bench: FAIL — Fixture größer als 20 MB (DC-QA-01)" >&2
  exit 2
fi

times=()
for r in $(seq 1 "$RUNS"); do
  start=$(date +%s%N)
  docker run --rm --network none -v "$FIXTURE":/repo:ro "$IMAGE":latest > /dev/null
  end=$(date +%s%N)
  ms=$(((end - start) / 1000000))
  times+=("$ms")
  echo "bench: Lauf $r — ${ms} ms"
done

median=$(printf '%s\n' "${times[@]}" | sort -n | sed -n 2p)
echo "bench: Median ${median} ms (Pass-Kriterium < 5000 ms — DC-QA-01)"
if [ "$median" -ge 5000 ]; then
  echo "bench: FAIL — Median über 5 s" >&2
  exit 1
fi
echo "bench: OK"
