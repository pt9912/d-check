# Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-19.

**Form:** folgt [Kurs-Modul 6](../../../../.harness/baseline/v1.4.0/regelwerk/modul-06-roadmap.md).

---

## Aktuelle Welle

**Keine aktive Welle.** welle-66-release-prep-aufgabenregel **abgeschlossen** —
[`slice-082`](../done/slice-082-release-prep-aufgabenregel.md) (Release-Prep-Regel
in [`releasing.md`](../../../user/releasing.md) §Release-Prep: ein neuer
Handbuch-§4-Abschnitt für ein Feature ist eine eigene Aufgabe, keine Anhängung —
schließt den strukturellen §4-Erosions-Bindepunkt aus slice-072; reine
Prozess-Doku, kein Release/ADR).

**Vorgänger:** welle-65-handbuch-aufgaben
([`slice-072`](../done/slice-072-handbuch-aufgabenorientierung.md) — §4 des
Benutzerhandbuchs aufgabenorientiert (§4.12-Monolith in §4.12–§4.16 aufgetrennt),
**Handbuch 1.42**, reine Doku). Davor welle-64-dpin-ergonomie
([`slice-081`](../done/slice-081-pins-hash-ergonomie.md) — `pins`/dpin voller Ist-Hash im
`link-stale`-Befund, **v0.51.1**). Davor welle-63-sources
([`slice-080`](../done/slice-080-sources-modul.md) — 19. Modul `sources`,
[ADR-0046](../../adr/0046-sources-upstream-content-drift.md) `Accepted`, **v0.51.0**).
Davor welle-62-zitat-verifikation
([`slice-079`](../done/slice-079-zitat-verifikation.md) — 18. Modul `citations`,
[ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md)
`Accepted`, **v0.50.0**). Davor welle-61-referenz-ventil-quell-skopus
([`slice-078`](../done/slice-078-ignore-refs-quell-skopus.md), **v0.49.0**) und die
welle-60-Kette ([`slice-071`](../done/slice-071-trace-cross-consistency-gate.md)/073/075/076,
v0.44–v0.47).

## Nächste Wellen

**Im Backlog (`next/`):** leer.

**Im Eingang (`open/`), auf Wellen-Einplanung wartend:**
[`slice-083`](../open/slice-083-regelwerk-v500-migration-analyse.md) — Delta-Analyse
der Baseline-Migration `v1.4.0` → `v5.0.0` mit Etappen-Vorschlag A–D (Vendoring ·
Modul-Delta lesen · Adaptions-Bereinigung · Form-Konformität). **Analyse zur
Abnahme**, keine Umsetzung: offen sind der Etappen-Schnitt, der Umgang mit drei
historischen Verweisen auf den vendorten Alt-Stand und die Form-Frage, ob neue
Artefakte ab sofort der Baseline-Vorlage folgen. Auslöser: Auftraggeber-Vorgabe
2026-07-25 (neubasiert 2026-08-01 auf `v5.0.0`, zwei weitere Majors über v3.5.2) —
vollständige Migration nach `v5.0.0`, der Baseline-Default sticht die repo-lokale
Adaption.

**Kandidat (noch kein Slice, auf Freigabe wartend):** der **RTM-Generator** (RTM
aus den Rückwärts-`Bezug`-Kanten erzeugen; von
[ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 als spätere
CR sequenziert, slice-071 ist sein Korrektheits-Harness). Auftraggeber-Nachreichung
2026-07-17: zusätzlich Artefakt-Titel + Kanten-Anmerkung. Freigabe und Scope offen.

Ferner ein `--print-version-md`-Scaffold, das ein `version.md`-Skelett mit
Platzhaltern auf stdout ausgibt (Familie `--print-config`/`--print-mk`/
`--suggest-config`; read-only, deterministisch). Produkt-Feature ⇒ Change Request
(`DC-FA-CLI-*` im Lastenheft) + Slice + Spezifikation-`.a`, **kein** ADR (additive
CLI-Ausgabe). Anlass: Nutzer-Frage 2026-07-04 zum Nachbau von `version.md` in
Fremd-Repos (der Aufbau selbst ist seit Handbuch 1.21 dokumentiert).

## Historische Trigger-Verschiebungen

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-11 | slice-012-Trigger: „slice-011 done" → „slice-011 **und** slice-013 done" | Der [`DC-QA-04`](../../../../spec/lastenheft.md#dc-qa-04--migrationsabdeckung-der-alt-tools)-Vergleichslauf gegen das erweiterte `docs-check.js` zeigte die Inline-Code-Pfad-Prüfung als Konsolidierungs-Lücke; Change Request [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in) (Lastenheft 0.3.0) als slice-013 eingeschoben |
| 2026-07-17 | **WIP-Limit wiederhergestellt:** slice-071 `in-progress`→`open` (Blocker), slice-076 `in-progress`→`next`; welle-60 führt nur noch slice-073 in Arbeit. Reihenfolge danach: slice-073 zu Ende (vier offene R1-Befunde + bestätigender Review) → Closure → slice-075 | `in-progress/` trug **drei** Slices gleichzeitig; Modul 5: „WIP-Limit pro Implementer = 1 ist eine harte Größe, kein Vorschlag" und `next→in-progress` verlangt „WIP-Limit frei". Bei slice-076 wurde die Bedingung beim Einplanen schlicht nicht geprüft (`6d60094`); slice-071 war bereits blockiert und hätte nach Modul 5 längst zurückgeführt gehört — beides still, bis der Auftraggeber die Regel einforderte. slice-075 erhält Vorrang vor slice-076, weil er produktiv verdrahtetes `trace.coverage` **verfälscht** (Auftraggeber-Meldung grid-gym), während slice-076 Blindheit ohne Falschaussage ist |
| 2026-07-17 | slice-074 aus welle-60 zurückgestellt (`in-progress/` → `open/`), Implementierung zurückgenommen; slice-076 in welle-60 nachgenommen | Drei unabhängige Reviews belegten an fünf aufeinanderfolgenden Fassungen dieselbe Klasse, zuletzt einen Stilles-Grün-Pfad (R3-F-1). Der Realdatenbeleg für slice-071 ist damit weiter blockiert — offen ausgewiesen statt still weitergeschoben. slice-076 kam aus dem Spike, den die Rücknahme ausgelöst hat |
| 2026-07-18 | slice-071 wieder aufgenommen (`open/`→`in-progress/`), welle-60 wieder aktiv | Blocker aufgelöst: die Direktiven-Toleranz aus slice-074 (v0.48.1) lässt den `17. Testarchitektur`-Abschnitt durchlaufen, statt an `architecture.md:913` abzubrechen. WIP-Slot frei (Modul 5), daher `open→next→in-progress` in einem Zug. Der von [ADR-0038](../../adr/0038-trace-cross-consistency.md) Entscheidung 7 geforderte Realdatenbeleg gegen grid-gym ist damit fahrbar |
| 2026-07-18 | **welle-61-referenz-ventil-quell-skopus eröffnet**; slice-078 `open`→`in-progress` (WIP-Slot frei nach welle-60-Abschluss) | §4-Vorfrage vom Auftraggeber entschieden: das erweiterte `ignore-refs`-Ventil (Quell-Skopus `in:`, `refs`/`keep`) wohnt als **neues geteiltes Bereichskürzel** (Ziel-Achsen-Pendant zu [`DC-FA-SCAN-001`](../../../../spec/lastenheft.md#dc-fa-scan-001--datei-auswahl-und-ignorier-regeln)), nicht als Änderung dreier Anforderungen — vermeidet die Verdreifachung der Ventil-Spezifikation, `codepaths.ignore-refs` bleibt Alias. Konsumenten-CR `ai-harness-course` |
| 2026-07-18 | **welle-61 abgeschlossen**; slice-078 `in-progress`→`done`, **v0.49.0 veröffentlicht** | Vollständige Kette umgesetzt: Lastenheft [`DC-FA-REF-001`](../../../../spec/lastenheft.md#dc-fa-ref-001--geteiltes-referenz-ventil-ignore-refs-mit-quell-skopus) + Spec + [ADR-0044](../../adr/0044-geteiltes-referenz-ventil-quell-skopus.md) `Accepted`, Code über `links`/`anchors`/`codepaths` + Alias (Mutations-gepinnt), Realdatenbeleg gegen `ai-harness-course` (Baseline 42 → Ventil 0, nicht durch Wegschauen), Review R1 ACCEPT-WITH-NITS (Nits eingearbeitet). WIP-Slot wieder frei |
| 2026-07-18 | **welle-62-zitat-verifikation eröffnet**; slice-079 `open`→`in-progress` (WIP-Slot frei nach welle-61) | Beide §4-Vorfragen entschieden: Adopter-Rückfrage **empirisch** (33/33 `datei:zeile`-Zitate in `ai-harness-init` in Inline-Code, null Prosa) ⇒ `codepaths`-Erweiterung; Zuschnitt Form (c) (Auftraggeber): Stufe 1/2 als Erweiterung von [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in), Stufe 3 (`verbatim`) als eigenes Modul mit Direktive `d-check:cite` (durch slice-074 entblockt). Adopter-CR `ai-harness-init` |
| 2026-07-18 | **welle-62 abgeschlossen**; slice-079 `in-progress`→`done`, **v0.50.0 veröffentlicht** | Vollständige Kette: [`DC-FA-CODE-001`](../../../../spec/lastenheft.md#dc-fa-code-001--explizite-pfade-in-inline-code-modul-codepaths-opt-in)-Erweiterung (`check-lines`) + neues [`DC-FA-CITE-001`](../../../../spec/lastenheft.md#dc-fa-cite-001--verbatim-zitat-verifikation-modul-citations-opt-in) (18. Modul `citations`) + Spec + [ADR-0045](../../adr/0045-zitat-verifikation-codepaths-erweiterung-und-citations-modul.md) `Accepted`; opt-in `codepaths.check-lines` (`citation-out-of-range`/`citation-inverted-range`) + Modul `citations` (`d-check:cite`, whitespace-normalisierter Teilstring, `citation-mismatch`), mutations-gepinnt; Realdatenbeleg gegen `ai-harness-init` (korrekt grün, Drift rot, Baseline 0 Direktiven belegt den Substrat-Caveat); Pre-Release-Review R1 BLOCK auf F-1 (Absturz bei `von=0`) → gefixt (1-basierte Untergrenze in beiden Schwestern) → ACCEPT. WIP-Slot wieder frei |
| 2026-07-19 | **welle-63-sources eröffnet**; slice-080 in `in-progress/` angelegt (WIP-Slot frei nach welle-62) | §4-Vorfragen (Nutzer) entschieden: Pin-Deklaration **beides** (Marker + Config), Quelltypen **Einzeldatei + Archiv** (`unpack: zip`); der `pins`/dpin-Hash-Ergonomie-Fix bleibt separat ([`slice-072`](../done/slice-072-handbuch-aufgabenorientierung.md)). Anlass: Nutzer-Frage „Drift gegen Upstream in d-check einbauen" — produktisiert `check_regelwerk_drift.py` als reusables Modul [`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz), erweitert [`DC-QA-03`](../../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit) um eine zweite Netz-Tür |
| 2026-07-19 | **welle-63 abgeschlossen**; slice-080 `in-progress`→`done`, **v0.51.0 veröffentlicht** | Vollständige Kette: [`DC-FA-SRC-001`](../../../../spec/lastenheft.md#dc-fa-src-001--upstream-content-drift-externer-quellen-modul-sources-opt-in-netz) + [ADR-0046](../../adr/0046-sources-upstream-content-drift.md) `Accepted` + Spec-Algorithmus-Sektion + 19. Modul `sources` (Marker+Config, Einzeldatei+Archiv, byte-genaues Content-Manifest). Doppel-Review: R1-doc BLOCK auf Manifest-Kern → ACCEPT, R2-code ACCEPT-WITH-NITS (Marker-Hash-64/Limits/Golden). Realdatenbeleg gegen echtes `lab-regelwerk.zip`, Config-Surface-Doku. Digest `sha256:9197fcf0…1d98`. WIP-Slot wieder frei |
| 2026-07-19 | **welle-64-dpin-ergonomie eröffnet**; slice-081 in `in-progress/` angelegt (WIP-Slot frei nach welle-63) | Nutzer-Entscheid: der in slice-079/080 als „separat" ausgewiesene `pins`/dpin-Ergonomie-Retrofit (voller Ist-Hash im `link-stale`-Befund) wird **eigener** Kleinst-Slice — NICHT in slice-072 (reine §4-Doku, kein Release) quergeschnitten. Nur die nicht stabilitätsgarantierte Befund-`message`, kein Vertragsdelta/ADR; Patch v0.51.1 |
| 2026-07-19 | **welle-64 abgeschlossen**; slice-081 `in-progress`→`done`, **v0.51.1 veröffentlicht** | dpin-Ergonomie: `pins.go` emittiert den vollen Ist-Hash im `link-stale`-Befund (mutations-echter Test); neue Handbuch-Aufgaben-Sektion §5 „Link-Inhalt pinnen". Patch (nur nicht stabilitätsgarantierte Befund-`message`, kein Vertragsdelta/ADR). Digest `sha256:fede3d02…d03c`. WIP-Slot wieder frei |
| 2026-07-19 | **welle-65-handbuch-aufgaben eröffnet**; slice-072 `open`→`in-progress` (WIP-Slot frei nach welle-64) | Der lange als Backlog geführte §4-Aufgabenorientierungs-Slice wird aufgenommen — nachdem der dpin-Ergonomie-Fix (fälschlich zunächst hierher geroutet, vom Auftraggeber korrigiert) als eigener slice-081 erledigt ist, ist slice-072 wieder **reine §4-Doku** (B-1…B-8, kein Release/ADR). No-Regress via die Handbuch-Test-Harnesse (`docexamples_test`/`handbook_examples_test`) in `make gates` |
| 2026-07-19 | **welle-65 abgeschlossen**; slice-072 `in-progress`→`done` | §4 des Benutzerhandbuchs aufgabenorientiert nachgezogen (Benutzerhandbuch-Standard §2/§5): der §4.12-`--trace`-Monolith in §4.12–§4.16 aufgetrennt (RTM/Coverage/Modalität/Kreuzverweis, `--print-mk`→§4.16), Grammatik-/Modalitäts-Referenz→§5, „0 Anforderungen"-Fehlerbild→§7, Versions-Prosa→§11, doppelte WAISE-Definition raus; §4.4-Straffung, `citations`-§5-Titel auf Task-Form; Handbuch 1.42. Zwei unabhängige Reviews: Cluster-Review ACCEPT, Abschluss-Gegenprobe BESTANDEN (12 erledigt/2 bewusst-Ref/0 offen). Reine Doku, kein Release/ADR; `make gates` grün (265/0). WIP-Slot wieder frei |
| 2026-07-19 | **welle-66-release-prep-aufgabenregel eröffnet**; slice-082 in `in-progress/` angelegt (WIP-Slot frei nach welle-65) | slice-072 räumte §4 redaktionell auf, beseitigte aber nicht die **Ursache** (jeder Feature-Slice hängte seine Fähigkeit an §4.12 an, statt eine Aufgabe zu schreiben). Auftraggeber-Entscheid: die billigste dauerhafte Sicherung — eine Release-Prep-Regel „neuer §4-Abschnitt = eigene Aufgabe" — als eigener Folge-Slice, nicht in slice-072 quergeschnitten. Reine Prozess-Doku, kein Release/ADR |
| 2026-07-19 | **welle-66 abgeschlossen**; slice-082 `in-progress`→`done` | Release-Prep-Regel in releasing.md §Release-Prep Punkt 4: ein neuer Handbuch-§4-Abschnitt für ein Feature ist eine **eigene** Aufgabe (nach Benutzerhandbuch-Standard §5), keine Anhängung — schließt den strukturellen §4-Erosions-Bindepunkt, den slice-072 nur redaktionell umging. Ehrlich unenforced (kein Gate), begründet mit Inline-Fakt (§4.12 wuchs auf ~330 Zeilen/8 Themen). Plan-Review ACCEPT-WITH-NITS eingearbeitet (Faktenfehler + Inline-Fakt-Altitude). Reine Prozess-Doku, kein Release/ADR; `make gates` grün (266/0). WIP-Slot wieder frei |
