# Realdatenbeleg — slice-080 Modul `sources` (Upstream-Content-Drift)

| Feld | Wert |
|---|---|
| Slice | slice-080 (`DC-FA-SRC-001`, ADR-0046) |
| Datum | 2026-07-19 |
| Quelle | `https://github.com/pt9912/ai-harness-course/releases/download/v1.4.0/lab-regelwerk.zip` — das **reale**, von d-check committet-vendored Kurs-Release-Asset (Archiv-Pfad `unpack: zip`) |
| Lauf | `docker run --rm -v <scratch>:/repo:ro d-check:latest` **mit Netz**; Config-Pin in `.d-check.yml`, `modules: [sources]` |

Belegt den **Archiv-Pfad** (entpacktes Content-Manifest) gegen echte Netz-Daten
sowie den „einmal laufen, gemeldeten Hash pinnen"-Workflow (die dpin-Ergonomie-
Sackgasse ist hier nativ gelöst — voller Ist-Hash statt `shortHash`).

## Läufe

1. **Dummy-Pin (64 Nullen) → `source-drift`, Exit 1.** Befund an `.d-check.yml:5`
   (die `url`-Zeile — Config-Pin-Zeilen-Capture wirkt), Ziel = die URL. Die
   `--json`-`message`: `errechnet sha256:91622e8b8eba9377f7d3d253e51997f7d2f9a5d3a52ecf8d8a4bce4467d0430c`
   `(Pin sha256:0000…0000)` — der **volle** 64-Hex-Ist-Hash als Re-Pin-Vorlage.
2. **Korrekter Pin (`…430c`) → 0 Befunde, Exit 0.** Realer Fetch + entpacktes,
   pfad-sortiertes Content-Manifest == Pin.
3. **Pin in GROSSSCHRIFT → 0 Befunde, Exit 0.** Case-Insensitivität (Vertrag
   F-3, Code-Nit R2-C-2) am realen Hash bestätigt.
4. **Ein-Zeichen-Drift (`…430c` → `…430d`) → `source-drift`, Exit 1.** Kleinste
   Pin-/Inhalts-Abweichung wird gefangen.

## Fazit

Modul `sources` erkennt Upstream-Content-Drift eines realen Archiv-Assets
korrekt — unverändert grün, ein Byte anders rot; der volle Ist-Hash macht das
Pinnen ohne externes Werkzeug möglich. Der Lauf war bewusst **mit** Netz (das
Modul ist die zweite, opt-in Netz-Tür neben `external`); d-checks eigene
netzlose Gate-Läufe bleiben unberührt, da `sources` nie im Default- oder
`.d-check.yml`-Live-Modulset steht.
