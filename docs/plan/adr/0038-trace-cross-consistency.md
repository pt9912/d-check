# ADR-0038 — Kreuzverweis-Konsistenz zweier Traceability-Sichten als `trace`-Unterfähigkeit

**Status:** Proposed
**Datum:** 2026-07-16
**Autor:** pt9912
**Schärft:** [`DC-FA-XREF-001.a`](../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency)
**Bezug:** [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in), [`DC-FA-CLI-009`](../../../spec/lastenheft.md#dc-fa-cli-009--requirements-traceability-matrix), [`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code), [`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen), [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in), [ADR-0035](0035-trace-coverage-quellen.md), [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md), [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus), [`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit).

## Kontext

Ein Konsument (grid-gym, Trigger 088) pflegt Anforderung→Design in zwei Sichten.
Die **Rück-Kanten** stehen als `Bezug`-Spalte in mehreren Tabellen von
`spec/architecture.md` (Prinzipien, Ports, Tabus, Komponenten, …), authort dort,
wo das Design lebt. Die **Vorwärts-Sicht** (`traceability.md` §27.1) ist laut
eigenem Intro ein „gegen `spec/architecture.md` **kuratierter**" Spiegel — also
das redundante Artefakt, das driftet. Realer Fall: die Vorwärts-Zeile einer
Architektur-Anforderung nannte `{GG-AR-COMP-CORE, GG-AR-COMP-DOMAIN}`, die
Rück-Kanten derselben Anforderung `{GG-AR-P-005, GG-AR-P-009, GG-AR-COMP-SCHED}` —
Schnittmenge null, von keinem Gate bemerkt.

Realdaten-Befund (bestimmt die Config): der Backward-Namensraum ist einheitlich
`GG-AR-*` und die Artefakt-ID steht als **erste Spalte** jeder `Bezug`-Tabelle
(deren Header aber variiert: `Kennung`/`Port-ID`/`Tabu-ID`/`Komponente`); die
`Bezug`-Zelle trägt Anforderungs-IDs in **Range-/Enum-Notation**
(`GG-SIM-001..009`, `GG-ARCH-006/007/008`) — dieselbe, die
[`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
schon range-aware liest.

Kein bestehendes Modul deckt das ab: `matrix` = Richtung, `trace.coverage` =
Abdeckung (≥1), `ids` = Existenz. Es fehlt der **Mengenabgleich zweier unabhängig
gepflegter Kanten-Sichten**. Die Fähigkeit ist deterministisch, hermetisch und
read-only — sie passt in das `--trace`/`--require-complete`-Idiom
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus),
[`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).

## Entscheidung

1. **Platzierung als `trace`-Unterfähigkeit**, kein neues Top-Level-Modul: opt-in
   `trace.cross-consistency` neben `trace.requirements`/`trace.coverage`,
   verkörpert als
   [`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in).
   Beide Sichten werden über den vorhandenen header-gebundenen Reader
   ([`DC-FA-REQ-001`](../../../spec/lastenheft.md#dc-fa-req-001--anforderungsquellen-als-headings-oder-tabellen))
   + die range-aware Span-Semantik von
   [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)
   gelesen — kein neuer Parser.

2. **Rück-Kanten sind die Quelle der Wahrheit.** Extraktion konsumenten-neutral:
   je Tabelle mit header-gebundener `Bezug`-Spalte ist die Artefakt-ID die **erste
   Spalte** (positionell, via `design-pattern` extrahiert — deren Header variiert,
   `Bezug` nicht); die `Bezug`-Zelle liefert range-aware die Anforderungs-IDs
   (`req-pattern`). Prosa-`Bezug:`-Zeilen ohne Tabelle sind v1-Out-of-Scope.

3. **Vorwärts-Sicht (§27.1) wird konsumenten-seitig auf konkrete `GG-AR-*`-IDs
   restrukturiert** (header-gebundene Anforderungs-/Design-Spalte, keine Familien-
   Wildcards, keine freie Prosa) — d-check spezifiziert nur den Lese-Vertrag, die
   Restrukturierung liegt beim Konsumenten.

4. **Invertieren + Mengen-Diff**, Modi `equal` (beide Differenzen gaten) /
   `superset` (nur `B \ F`). `1:N` Normalfall; Range-/Enum-Notation beidseitig
   expandiert (Reuse von
   [`DC-FA-COV-001`](../../../spec/lastenheft.md#dc-fa-cov-001--kuratierte-coverage-quellen-der-rtm-tracecoverage-opt-in)).

5. **Ableitungssprünge per `exclude-req` (RE2)** als benanntes Ventil (kuratierte
   Kante mit eigener Drift-Gefahr, wie `matrix.exclude-sections`) — **kein**
   gelöstes Problem. *(Nachtrag 2026-07-17: sein **Totalfall** ist inzwischen
   geguardet — ein Ventil, das alle Anforderungen verschluckt, ist Vakuum, siehe
   Entscheidung 8. Die Drift-Gefahr der **Teil**-Ausschlüsse bleibt unverändert
   ungelöst.)*

6. **Eine Pattern-Syntax (RE2), fail-closed, ohne Block byte-identisch.** Der
   Abgleich gatet über das **globale** `--require-complete`
   ([`DC-FA-CLI-011`](../../../spec/lastenheft.md#dc-fa-cli-011--vollständigkeits-prüfung-als-opt-in-exit-code)),
   nicht über einen block-lokalen Schalter.

7. **Generator ist der Ziel-Zustand, aber eigene spätere CR.** Sobald das Gate die
   Treue der beiden Sichten erzwingt, wird §27.1 aus den Rück-Kanten **generiert**
   (Handpflege entfällt) — ein Schreib-Pfad (näher an
   [`DC-FA-CLI-008`](../../../spec/lastenheft.md#dc-fa-cli-008--reparatur-patch)),
   der eigene ADR-Behandlung verdient. Dieses Gate ist sein Korrektheits-Harness
   („erst das Gate, dann der Generator").

8. **Vakuum ≠ Konsistenz — aber eine einseitig leere Sicht ist ein Ergebnis.**
   (Nachtrag, 2026-07-17, Status noch `Proposed`.) Die Namensraum-Kongruenz
   (Entscheidungen 2/4) war zunächst nur als **Vorbedingung beschrieben** —
   beschriebene Vorbedingungen halten nicht. Ein `design-pattern`, das kompiliert,
   aber am Artefakt-Namensraum vorbeigreift, räumt (weil es **geteilt** ist) *beide*
   Sichten leer und liefert `0 Differenz(en)`/Exit 0: eine behauptete, nie geprüfte
   Konsistenz. Der Abgleich prüft daher seine eigene **Vakuität** — keine Kante aus
   beiden Sichten ⇒ Exit 2; ebenso eine kantenleere Rück-Sicht unter
   `mode: superset`, wo allein `B \ F` gatet und damit konstruktionsbedingt nie ein
   Befund entstehen kann.
   **Bewusst nicht** geguardet ist die *einseitig* leere Sicht: der Diff läuft über
   `keys(F) ∪ keys(B)` und ist für `F = ∅` wohldefiniert — eine noch
   unrestrukturierte Vorwärts-Sicht bei gepflegten Rück-Kanten ist genau der
   Bootstrap-Zustand, den Entscheidung 3 dem Konsumenten aufträgt und den
   Entscheidung 7 als Generator-Eingang braucht. Ein symmetrisch je Sicht feuernder
   Guard würgte ihn mit einer Config-**Fehl**diagnose ab, statt die Rück-Kanten laut
   zu melden.
   **Abgrenzung zu [ADR-0037](0037-trace-tabellenquellen-nullmengen-guard.md):**
   dort hängt der Nullmengen-Guard an der *explizit gesetzten Quelle*
   (`strictSource`) und schützt die RTM selbst; hier ist der Bezugspunkt nicht eine
   Sicht, sondern der **Vergleich** — geguardet wird, was nie einen Befund liefern
   kann. Dieselbe Lehre („fail-closed statt irreführender Nullmenge"), enger
   gefasster Auslöser.
   **Der Guard fasst die Wirkung, nicht die Ursache** (Nachtrag 2026-07-17): er
   wird **nach** dem `exclude-req`-Ausschluss (Entscheidung 5) gemessen, denn
   maßgeblich ist, was am Ende tatsächlich verglichen wird. Ein Ventil, das alle
   Anforderungen verschluckt, schaltet das Gate ebenso still ab wie ein
   fehlgreifendes Muster — es ist selbst eine kuratierte, drift-fähige Kante, und
   der normative Satz oben („geguardet wird, was konstruktionsbedingt nie einen
   Befund liefern kann") trägt beide Ursachen. Eine Ursachen**liste** risse bei der
   nächsten unbekannten Ursache erneut; die Wirkungs-Fassung nicht. Weil
   `exclude-req` dasselbe Prädikat auf **beide** Sichten anwendet, heißt
   post-Ausschluss-Leere „jede vergleichbare Anforderung ist ausgenommen" — und
   „hier nicht prüfen" wird korrekt ausgedrückt, indem man den opt-in Block
   **weglässt**, nicht indem man ihn konfiguriert und leerräumt.
   *Provenienz:* das unabhängige Closure-Review (R3) las die schweigende
   Fehlerpräzedenz als „prä-Ausschluss gemeint" und neigte dazu; die Auflösung
   gegen diese Neigung fiel zugunsten der Wirkungs-Fassung — R4 bestätigte sie und
   revidierte die eigene Lesart. Vertrag nachgezogen als Lastenheft-CR 0.44.2.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **`trace`-Unterblock, Reader-Reuse** (gewählt) | Keine Parser-Duplikation; erbt Header-Binding + range-aware Span-Semantik; bleibt Traceability-Familie | `trace`-Konfigurationsfläche wächst |
| Eigenes Modul `symmetry` | Maximal generisch; klar abgegrenzt | Dupliziert Tabellen-Reader-Plumbing; Über-Generalisierung |
| Erweiterung von `matrix` | kennt Klassen/Kanten | `matrix` = Richtung, nicht Mengen-Diff; semantischer Bruch |
| Generator zuerst (iv vor iii) | erreicht den Ziel-Zustand direkt | Kein bewährter Vergleichs-Harness vor der Generierung; verwirft die Treue-Prüfung |
| Prosa-Annotations-Parser (v1-Entwurf) | trifft Annotationen ohne Tabelle | Format-gekoppelt, portabilitätsschwach |

**Fitness-Funktion:**

- Der reale grid-gym-Drift wird geflaggt; konsistentes `1:N` läuft grün;
  Range-/Enum-Notation beidseitig korrekt expandiert.
- In die Mittelschicht verschobene Familien werden nicht fälschlich als Waisen
  gemeldet (Ventil greift).
- Ein Abgleich, der nichts vergleicht, ist **rot**, nicht grün: Muster am
  Namensraum vorbei ⇒ Exit 2; eine einseitig leere Vorwärts-Sicht meldet dagegen
  ihre Rück-Kanten laut (Entscheidung 8).
- Ohne Block jede RTM byte-identisch
  ([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)), kein
  Schreibzugriff/Netz
  ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Der Rückkanten-Vertrag (erste-Spalte-ID + header-gebundene `Bezug`-Spalte,
  scoped auf Abschnitte) ist konsumenten-neutral formuliert.

### Offene Designpunkte

- **Gelöst durch Realdaten:** Rückquelle ist tabellarisch (erste-Spalte-keyed,
  heterogene Header) — Prosa-`Bezug` bleibt v1-Out-of-Scope; Namensraum einheitlich
  `GG-AR-*` (die Kurzschreibweise `COMP-*`/`P-*` war eine Abkürzung); beide Seiten
  range-aware.
- **Konsumenten-Aufgabe (nicht d-check):** §27.1 auf konkrete IDs restrukturieren
  (Wildcards/Prosa entfernen), die exakten Backward-`sections` benennen.
- **Schlüssel-Schreibweise offen (d-check besitzt die Schlüssel):**
  `artifact-id-column: first` als positioneller Sentinel neben header-Namen;
  Bereichskürzel `XREF` / Blockname `cross-consistency`.

## Konsequenzen

- **Positiv:** ein ungefangener Doku-Drift-Typ wird deterministisch geflaggt;
  volle Reader-Wiederverwendung; opt-in, byte-identisch für Nicht-Konsumenten; das
  Gate ist zugleich der Korrektheits-Harness des späteren Generators.
- **Negativ / Kosten:** die `trace`-Fläche wächst; `exclude-req` ist ein
  kuratiertes Ventil mit eigener Drift-Gefahr; der positionelle `first`-Modus ist
  ein neuer Extraktions-Sonderweg (durch heterogene Backward-Header begründet);
  §27.1 muss konsumenten-seitig restrukturiert werden, bevor das Gate grün wird.
- **Verworfen:** eigenes Modul, `matrix`-Erweiterung, Prosa-Parser, block-lokaler
  Gate-Schalter, Generator-zuerst (jeweils oben begründet).

## Geschichte

| Datum | Ereignis |
|---|---|
| 2026-07-16 | Proposed. Change Request grid-gym (Trigger 088, ADR 0080 §4.4 iii), v2 nach Design-Review; Ziel-Architektur „Gate jetzt, Generator später" gegen grid-gyms reale Quellen bestätigt. Umsetzender Slice slice-071. |
| 2026-07-17 | Entscheidung 8 auf die **Wirkungs-Fassung** gezogen (Messung nach dem `exclude-req`-Ausschluss) + Entscheidung 5 um den geguardeten Totalfall annotiert; Status weiterhin `Proposed`. Anlass: Review R3 reproduzierte, dass `exclude-req: '.'` das Gate bei realem Drift still abschaltete — dieselbe Silent-Green-Klasse, andere Ursache. Vertrag nachgezogen als Lastenheft-CR 0.44.2. |
| 2026-07-17 | Entscheidung 8 (Vakuität) nachgetragen, Status weiterhin `Proposed`. Anlass: unabhängiges Closure-Review zu slice-071 — R1 reproduzierte, dass die nur *beschriebene* Namensraum-Vorbedingung ein stilles Grün zuließ (HIGH); R2 wies den ersten, symmetrisch je Sicht feuernden Fix als vertragswidrig nach (er brach den Bootstrap-Zustand aus Entscheidung 3). Vertrag nachgezogen als Lastenheft-CR 0.44.1 ([`DC-FA-XREF-001`](../../../spec/lastenheft.md#dc-fa-xref-001--kreuzverweis-konsistenz-zweier-traceability-sichten-tracecross-consistency-opt-in)) + [`DC-FA-XREF-001.a`](../../../spec/spezifikation.md#dc-fa-xref-001a--kreuzverweis-konsistenz-cross-consistency) Schritt 5. |
