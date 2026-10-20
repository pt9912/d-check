# Verifikation — slice-277: Baseline-Pin auf v6.18.1, mechanisch

- **Rolle:** Verifier (Modul 11) — DoD, Plan und MR gegen den Stand
- **Gegenstand:** `docs/plan/planning/done/wellenlos/slice-277-baseline-v6181-bump.md` §1, §2, §3 (samt den Plan-Änderungen), §4, §6; [MR-080](../../harness/conventions/MR-080-baseline-v6181.md)
- **Range:** `b9211511..HEAD` (Spitze `6dc40149`); der Review-Report `2026-10-20-slice-277-baseline-v6181.md` nur als Kontext
- **Datum:** 2026-10-20

## Messungen (Kommando und Ergebnis)

| # | Kommando | Ergebnis |
|---|---|---|
| M1 | `make baseline-verify` (selbst gefahren) | Exit 0; alle 54 Dateien `OK`, `verify ok (54 Dateien, vollständig)` |
| M2 | `bash tools/harness/fetch-baseline-cache.sh --verify` | Exit 0; `verify ok (54 Dateien, vollständig)`, Alias-Auflösung ok |
| M3 | `make baseline-freshness` | Exit 0; `check-latest OK (Currency) — Pin v6.18.1 ist der neueste Release-Tag`, `OK (Content) — … unverändert (Bytes == vendored SHA256SUMS)` |
| M4 | unabhängig: Asset `lab-regelwerk.zip` von `v6.18.1` geladen, entpackt, `diff -r` gegen `.harness/baseline/v6.18.1` | einziger Unterschied: `SHA256SUMS` nur vendort (erzeugt); 54 Dateien byte-gleich |
| M5 | `ls .harness/baseline` | nur `v6.18.1`; `harness/conventions/done/MR-074-…` vorhanden, `harness/conventions/MR-080-…` aktiv |
| M6 | `make gates` ohne Pipe, Exit-Code erfasst | `GATES_RC=0`; `baseline-verify + workflow-pins + doc-check + lint + test + arch-check + coverage-gate + semgrep + gate-consistency + planning-check green`; doc-check `1187 Datei(en), 0 Befund(e)`; `git status` danach leer |
| M7 | Re-Vendor-Weg in Wegwerf-Worktree, Baum entfernt, `fetch-baseline-cache.sh v6.18.1` auf macOS | Exit 0, `verify ok (54 Dateien, vollständig)`; der neu erzeugte Baum ist gegenüber dem Commit unverändert (`git status` leer) |
| M8 | dasselbe mit der Zeile zurückgenommen (`[ "$src_n" = "$dst_n" ]`) | bricht ab: `Quelle (54) != vendored (      54) — Kopierschritt hat Dateien übergangen` — aus dem richtigen Grund (BSD-`wc` füllt auf) |

## Nachzählung MR-080 (Plan §3, Prüfpunkt 3)

Alter Baum aus `git archive b9211511`, neuer Baum aus dem Arbeitsstand; Vergleich ohne `<!-- Quelle:`-Zeilen, Versionsstempel normalisiert.

| Behauptung | Nachzählung |
|---|---|
| 55 Dateien vorher wie nachher, keine neu, keine entfallen | 55 / 55, Dateiliste identisch (54 + `SHA256SUMS`) |
| 8 Delta-Dateien, 63 Diff-Zeilen | 8 Dateien, 63 Zeilen; je Datei 15 / 17 / 11 / 10 / 4 / 2 / 2 / 2 wie genannt |
| 104 Dateien / 262 Vorkommen am Vorzustand | 104 / 262 (`git grep -c` auf `b9211511`, ohne `.harness/baseline/`) |
| Frozen 51 Dateien / 156 Vorkommen | `done/` 24/68, Reviews 25/85, CR 1/1, CHANGELOG 1/2 = 51/156 |
| 49 lebend retargetet, davon 35 aktive `MR-*` | 14 Nicht-MR-Dateien + 35 `MR-*` = 49; im Diff stehen genau diese; `MR-074` (Gegenstand, verschoben), `MR-075`, `evidence/slice-261.md` und der Plan bleiben bewusst stehen |
| 17 lebende `d-check:cite`-Direktiven, die zwei im Plan zählen nicht | **15** ohne die zwei im Plan; 17 nur mit ihnen (siehe LOW 1) |
| eine Zitat-Verschiebung `modul-10` 87→99 | stimmt; Zeile 87 alt und Zeile 99 neu tragen denselben Text (`„geprüft, ohne Befund"; eigene Sektion unten).`) |

## DoD-Punkte

- **Re-Vendor — erfüllt.** M1–M5. Der Baum ist byte-gleich zum Upstream-Asset (M4, unabhängig von der Content-Prüfung des Skripts). `SHA256SUMS` stimmt (M1), der alte Baum ist entfernt, `conventions.md` §Baseline nennt `v6.18.1` über `MR-080`, `MR-074` liegt in `done/` und steht in der Aufgelöste-Zeile mit Nachfolger `MR-080`.
- **Spiegel — erfüllt.** `git grep "v6\.17\.0"` außer `.harness/baseline/`: Reste nur in `docs/reviews/` (26, davon 25 Bestand und der Report dieses Slice), `done/` (24), `harness/conventions/done/` (1, `MR-074`), `docs/plan/cr/`, `CHANGELOG.md` — alle Frozen —, dazu `MR-075` (Tatsache über die Fassung, in der die Umnummerierung kam), `evidence/slice-261.md` (unveränderlich ab Merge), `MR-080` und der Plan (Gegenstand). Keine Nennung ohne Zeichen `v` (`6.17.0`/`6.17`) in lebenden Dateien; `tools/`, `Makefile`, `.d-check.yml`, `.github/`, `.claude/` ohne Rest. Die acht Symlinks unter `.claude/rules/` lösen gegen `v6.18.1` auf; die Schleife meldet kein `DEFEKT`. Die Cite-Spannen sind neu geankert (eine Verschiebung, 15 Direktiven geprüft, `doc-check` grün).
- **`make gates` grün; Freshness aktuell — erfüllt** (M6, M3).
- **Review — Report liegt vor**; Verifikation mit diesem Bericht.
- **Offen, keine Verletzung:** Beobachtungs-Register, Closure-Notiz §7, Risiko-Ausgänge §6 (alle drei noch „offen bis Closure"; Gegenprobe Risiko 1 mit M1/M6 und Nachzählung belegbar: 55 / 55, jeder lebende Anker löst auf, `doc-check` 0 Befunde), drei Paarungen, `git mv` nach `done/wellenlos/`.

## Bewusstes Brechen (Wegwerf-Worktree, danach entfernt; Repo `git status` leer)

| Zurückgenommen | Lauf | Rot durch |
|---|---|---|
| ein Byte an `regelwerk/modul-13-quality-gates.md` angehängt | `make baseline-verify` | Exit 2; `modul-13-quality-gates.md: FAILED`, `1 computed checksum did NOT match` — die verlangte Ursache |
| lebender Link auf `.harness/baseline/v6.17.0/regelwerk/modul-13-quality-gates.md` in `README.md` | `make doc-check` | `README.md:387 … target-missing Linkziel existiert nicht`, 1 Befund; ohne den Link 0 Befunde |
| Zitat-Spanne `reviewer.md:237` auf `modul-10-review-harness.md:87-87` | `make doc-check` | `.harness/skills/reviewer.md:237 … citation-mismatch Zitattext ist kein zusammenhängender Teilstring der Quell-Spanne`, 1 Befund |

Der dritte Beleg ruht nur auf `doc-check`; eine Verschiebung wird gemeldet, solange die Direktive eine Zeilenspanne trägt (alle 15 lebenden tun es, 0 tag-only).

## Plan-vs-Code (Abgrenzung §1)

Dateien im Diff: Baum `v6.18.1` (neu) und `v6.17.0` (entfernt), 35 aktive `MR-*` plus `MR-080` neu, `MR-074` nach `done/`, 14 lebende Dokumente, 8 Symlinks, `fetch-baseline-cache.sh`, `slice-274` (eine Zeile), `slice-277`, `slice-278` (neu), Review-Report. Jede Datei ist durch §3, die Plan-Änderung „Under-Copy-Barriere" oder die Plan-Änderung „Frozen-Ausnahme `slice-274` / Link-Tiefen `MR-074` / Zitat-Spanne" gedeckt, mit einer Ausnahme (LOW 2). Eingehalten: keine Adoption der Deltas (Delta nur nach Datei im MR beschrieben, `.d-check.yml` unverändert, `reviews`-Block nicht angefasst); `tools/harness/selbstpruefung.sh` unverändert; keine Tool-Pins im Diff; keine Distribution-Fläche. Frozen-Dateien byte-stabil bis auf `slice-274` (eine Zeile, deklarierte Ausnahme, Text unberührt).

## Befunde

- **HIGH / MEDIUM:** keine.
- **LOW 1:** `MR-080` §`d-check:cite`-Neu-Ankern sagt „gemessen **17** lebende Direktiven … (die zwei Direktiven im Plan dieses Slice zählen nicht mit) … Die übrigen 16". Nachgezählt sind es **15** ohne die zwei im Plan (17 nur mit ihnen); „die übrigen" wären 14. Das Ergebnis (eine Verschiebung) trägt, die Zahl und die Ausschlussklammer widersprechen sich (`harness/conventions/MR-080-baseline-v6181.md`).
- **LOW 2:** `slice-278` ist neu angelegt und von `MR-080` als Folge-Slice genannt, steht aber weder in der §3-Tabelle noch in einer Plan-Änderungsnotiz des Plans (§7 „Folge-Slices" ist noch `—`). Der Slice entsteht aus der Delta-Messung und ist in §1 nicht ausgeschlossen, also keine Ausweitung, aber der Plan-Wortlaut nennt ihn nicht; beim Befüllen von §7 nachtragen (`docs/plan/planning/done/wellenlos/slice-277-baseline-v6181-bump.md`).
- **INFO:** Die Plan-Zahl „52 Dateien außerhalb der Frozen-Verzeichnisse … beim Anlegen" (§3) und die MR-Zahl 104 messen verschiedene Mengen (ohne bzw. mit Frozen); beide sind für sich richtig, nur der Plan nennt die 104 nicht.
- **INFO:** Die Frozen-Ausnahme `slice-274` ist im Plan nachträglich nachgetragen (Plan nennt es selbst), im MR als einzige Ausnahme geführt; die Alternative `ignore-refs` wäre eine Gate-Senkung gewesen — die Wahl ist nachvollziehbar.

## Urteil

DoD-Liefer-Punkte „Re-Vendor", „Spiegel" und „`make gates` grün" erfüllt, durch eigene Läufe und unabhängige Nachzählung belegt; das Bewusste Brechen zeigt für Integrität, toten Link und Zitat-Verschiebung jeweils die verlangte Ursache, und die Reparatur an `fetch-baseline-cache.sh` läuft auf macOS-BSD-`wc` durch und bricht ohne sie genau an der behaupteten Stelle. Plan und MR stimmen mit dem Baum überein bis auf die zwei LOW-Nachträge. Keine DoD-Verletzung. Offen sind nur die Closure-Punkte.
