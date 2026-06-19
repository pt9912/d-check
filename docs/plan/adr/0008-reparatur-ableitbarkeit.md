# ADR-0008 — Reparatur-Modus: nur deterministisch ableitbare Fixes; Best-Guess review-pflichtig

**Status:** Accepted
**Datum:** 2026-06-19
**Autor:** pt9912
**Bezug:** [`DC-FA-CLI-008`](../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch),
[`DC-FA-CLI-007`](../../../spec/lastenheft.md#dc-fa-cli-007--diagnose-modus)
(die bedienten Anforderungen),
[`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit),
[ADR-0002](0002-distribution-ghcr-image.md), [ADR-0005](0005-modul-layout-hexagon-ordner.md)
**Schärft:** [`DC-FA-CLI-008.a`](../../../spec/spezifikation.md#dc-fa-cli-008a--reparatur-patch)
und [`DC-FA-CLI-007.a`](../../../spec/spezifikation.md#dc-fa-cli-007a--diagnose-modus)
(begründet, warum nur bestimmte Befundarten einen Fix-Kandidaten bzw.
Edit liefern) — nicht das Lastenheft.

## Kontext

`--doctor` (slice-025) und `--repair` (slice-026) leiten ihre
Fix-Kandidaten aus derselben Quelle ab (`core.FixCandidateFor`);
`--repair-broad` ergänzt einen Best-Guess, `--doctor --json` (slice-029)
ist ein drittes Rendering desselben Modells. In der Umsetzung liefern nur
zwei der vierzehn Befundarten einen Edit: `id-unlinked` (konservativ) und
`target-missing` als Datei-Move (breit). Warum gerade diese — und warum
nicht mehr, etwa über die git-Historie — war bislang nur im Algorithmus
(Spezifikation/Code) sichtbar, nicht als Entscheidung festgehalten. Diese
ADR zieht die in slice-026/029 getroffene Schranke retroaktiv nach.

## Entscheidung

Ein Reparatur-Edit (und ein Diagnose-`fixCandidate`) entsteht **nur, wenn
die Korrektur aus Befund, gescanntem Baum und Konfiguration
deterministisch ableitbar ist** — auf einer von zwei Stufen:

1. **Konservativ (eindeutig):** die Korrektur ist eindeutig bestimmt. In
   dieser Version `id-unlinked` → Markdown-Link auf das in der
   `ids`-Regel deklarierte Definitions-`target` (nur nackte
   Prosa-Vorkommen).
2. **Breit (eine überprüfbare Vermutung, review-pflichtig):** genau ein
   starkes, im Baum überprüfbares Signal. In dieser Version
   `target-missing` → eine im Scan-Bestand **eindeutig** gleichnamige
   Markdown-Datei (Verschiebung). Marker auf stderr; der Mensch prüft.

Alle übrigen Befundarten liefern **keinen** Edit und bleiben Befund. Neue
reparierbare Fälle erfordern einen Change Request am Lastenheft, der diese
Schranke erneut anlegt.

## Verglichene Alternativen

| Alternative | Pro | Contra |
|---|---|---|
| **Nur deterministisch ableitbar; Best-Guess review-pflichtig (gewählt)** | kein „grün ≠ richtig"-Fehlpatch; hermetisch; deterministisch ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)) | deckt nur 2 von 14 Befundarten ab |
| Fuzzy-/Nächster-Treffer-Reparatur (z. B. `anchor-missing` → ähnlichster Slug) | mehr Auto-Fixes | rät auf den falschen Abschnitt, rendert **still falsch** — die gefährlichste Fehlerart |
| **git-/VCS-historienbasierte Move-/Rename-Erkennung** | fände auch Umbenennungen (über Inhalts-Ähnlichkeit, nicht nur den Namen) | Eingabe ginge über den gescannten, read-only gemounteten Baum hinaus ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)); das distroless/static-Image ([ADR-0002](0002-distribution-ghcr-image.md)) enthält kein git; `.git` ist beim read-only-Mount nicht garantiert; bräuchte einen neuen nicht-hermetischen Port ([ADR-0005](0005-modul-layout-hexagon-ordner.md)), analog `external` |
| Heuristisches/LLM-gestütztes Reparieren | breit anwendbar | nicht deterministisch, nicht hermetisch, nicht auditierbar |
| Auto-Entfernen von Policy-Befunden (`matrix-forbidden`, `repo-escape`, …) | „grünerer" Lauf | verdeckt echte Architektur-/Sicherheitssignale, statt sie zu melden |

## Konsequenzen

- Die nicht reparierten Befundarten zerfallen in vier Klassen, deren Fix
  nicht deterministisch ableitbar ist:
  - **externer Laufzeit-Zustand** (`external-status`, `external-timeout`,
    `external-redirects`) — kein Quell-Edit behebt das Remote;
  - **Policy-/Architektur-Signale** (`matrix-inactive`,
    `matrix-forbidden`, `repo-escape`, `hostpath-forbidden`, `symlink`) —
    ein Auto-Fix verdeckte das Problem;
  - **mehrdeutiges Ziel** (`anchor-missing`, `codepath-missing`) — eine
    falsche Vermutung wäre still falsch;
  - **mehrdeutige Umstrukturierung** (`span-unclosed`, `span-nested-link`)
    — der Autoren-Intent ist nötig.
- Die breite Move-Erkennung bleibt bewusst basisnamen-basiert (hermetisch,
  deterministisch) statt git-basiert; sie erkennt Verschiebungen, keine
  Umbenennungen.
- Ein künftiges VCS-Modul wäre die Tür für git-basierte Erkennung — als
  opt-in (analog `external`), mit eigener Anforderung; es würde die
  Determinismus-Zusage für dieses Modul ausdrücklich relativieren.

## Absicherung

Keine dedizierte Fitness Function. Die Schranke ist getragen durch die
konservative Implementierung (`core.FixCandidateFor` liefert nur für
`id-unlinked`, der breite Best-Guess nur für `target-missing`), die
Akzeptanztests (`make test`) und die Determinismus-Zusage
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)); neue
reparierbare Fälle erfordern einen Change Request.

## Re-Evaluierungs-Trigger

- Bedarf an git-basierter Move-/Rename-Erkennung → Change Request + opt-in
  VCS-Modul (eigene Anforderung), das diese Entscheidung für das Modul
  relativiert (ggf. Folge-ADR mit `Supersedes`).
- Eine neue Befundart mit eindeutig ableitbarem Fix → Aufnahme in die
  konservative bzw. breite Stufe per Change Request.

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-06-19 | Proposed → Accepted (nachgezogene Dokumentation der in slice-026/029 getroffenen Reparatur-Schranke) |
