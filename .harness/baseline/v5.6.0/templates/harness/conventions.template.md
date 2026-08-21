# Harness-Konventionen

> **Template-Hinweis.** Diese Datei ist eine Vorlage für
> `harness/conventions.md` deines Repos. Kopiere sie nach
> `harness/conventions.md`, ersetze `<Platzhalter>` und lösche
> diesen Block. Pflichtgliederung folgt
> [Baseline-Regelwerk §harness/conventions.md als Konventionsspeicher](../../regelwerk/grundlagen-harness-dateien.md#harnessconventionsmd-als-konventionsspeicher).
>
> **Was diese Datei trägt:** repo-lokale Strukturregeln und Adaptionen
> ggü. der adoptierten Harnesskonvention (Baseline). Sie ist
> **Pflicht** (Existenz), die Form (Einzeldatei vs. Verzeichnis,
> ADR-artig vs. Prosa) ist **Wahl**.
>
> **Was diese Datei NICHT trägt:** Kurs- oder Baseline-Konventionstext
> wird nicht dupliziert — Pointer reichen. Sonst drift gegen die Quelle.

---

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln dieses Repos
gegenüber der adoptierten Harnesskonvention (Baseline). Sie ist der
Default-Ort für:

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt.
  Der Baseline-Default wird als Teil der `MR-000`-Aussage festgehalten;
  ein abweichendes Präfix oder Schema ist ein eigener `MR`-Eintrag.
- **Zusatzklassen-Deklarationen** für repo-spezifische
  Bindung-Klassen in der Sensors-Tabelle, die über die vier kanonischen
  hinausgehen (ADR, Carveout, Schwelle, Reproduzierbarkeit).
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die
kanonische Quelle (Source Precedence). Diese Datei ist konformitäts-
bringend für *Form*-Fragen, nicht autoritativ über Inhalt.

## Baseline

<!--
Welche Harnesskonvention wird adoptiert? Stand und Datum festhalten,
damit spätere Adaptionen einen Bezugspunkt haben.
-->

- **Konvention:** <Name, z. B. "AI-Harness-Kurs", interner Standard, Industrie-Norm>
- **Stand:** <Datum oder Version, z. B. "Template-Set 2026-06">
- **Datum der Adoption:** <Datum>

## Adoptierte Konventions-Quellen

<!--
Pointer auf die Quellen der Baseline. KEINE Wiederholung des Inhalts —
nur Verweise.

Das Agenten-Regelwerk ist die Quelle, die ein Code-Agent statt des
vollen Lehrmaterials liest (operatives Regelwerk ohne Didaktik). Es
ist derivativ — bei Konflikt gilt das Lehrmaterial.
-->

- **Extern (Lehrmaterial):** <Pfad oder URL>
- **Vendored Baseline (Regelwerk + Templates):** aus dem self-contained
  Release-Asset
  https://github.com/pt9912/ai-harness-course/releases/download/v5.6.0/lab-regelwerk.zip
  nach `.harness/baseline/<tag>/{regelwerk,templates}/` entpackt (netzlos,
  `SHA256SUMS`) — adoptierten Stand notieren (Stand-Zeile in
  `regelwerk/README.md`, z. B. „Kurs-Welle 24 · 2026-07-16"; Wellen-Register:
  CHANGELOG.md im Kurs-Repo); für harte Reproduzierbarkeit das Asset eines Tags
  ziehen statt `latest`.
- **In-Repo (verkörperte Form):** <Pfade zu deinen kopiert-und-ausgefüllten
  Artefakten> — die vendored `.harness/baseline/<tag>/templates/` sind die
  Referenz-Form („Ziel-Form" des Regelwerks); deine eigenen Dateien sind daraus
  kopiert und ausgefüllt.

## Adaptions-Block

Regeln dieser Sektion: Diese Datei trägt den **Index**, nicht die Einträge.
Jede Adaption ist eine eigene Datei unter `harness/conventions/`, kopiert aus
`harness/conventions/MR-NNN-titel.template.md` der vendored Baseline;
ist ihr Auflösungs-Trigger eingetreten, wandert sie per `git mv` nach
`conventions/done/`. Der Zustand ist die Verzeichnis-Position, kein
Status-Feld. Der Grund für den Schnitt: Was hier steht, liest **jeder**
Agentenlauf — aufgelöste Adaptionen gehören nicht in diesen Pfad
(Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/conventions.md als Konventionsspeicher).

### MR-000 — Baseline-Aussage

Bleibt hier: Sie ist keine Adaption, sondern die Adoptions-Erklärung, und
sie gilt für jeden Lauf.

- **Datum:** <Datum>
- **Geltungsbereich:** gesamtes Repo
- **Ersetzt-Baseline-Regel:** — *(keine; dieser Eintrag ist die
  Adoptions-Erklärung, keine Adaption)*
- **Adaption:** *keine inhaltlichen Adaptionen ggü. Baseline-Default
  für Verzeichniskonvention, Lifecycle-Regeln, Carveout-Disziplin,
  ID-Schema (`<PREFIX>-FA-*`, `<PREFIX>-QA-*`, `SPEC-<NNN>`, `ARC-<NNN>`,
  `ADR-<NNNN>`, `CO-<NNN>`, `slice-<NNN>`, `MR-<NNN>`, `BEO-<NNN>`, `RC-<NNN>` — nur das
  Vertrags-Präfix wird repo-weit festgelegt, z. B. `LH`; `SPEC-*` und
  `ARC-*` kodieren das Stratum und sind fest, siehe Baseline-Regelwerk
  `grundlagen-source-precedence.md` §ID-Schema als Klammer;
  bei mehreren gleichzeitig schreibenden Entwicklern für die Artefakte mit
  je eigener Datei (`ADR-*`, `CO-*`, `slice-*`, `welle-*`) zusätzlich das
  Bereichssegment und damit den Zählraum je Sub-Area festlegen —
  `SPEC-*`/`ARC-*` bleiben davon ausgenommen und zählen fortlaufend je
  Datei, siehe Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe).*
- **Begründung:** Initial-Setzung. Spätere Adaptionen werden als
  `MR-<NNN>` nachgetragen.
- **Auflösungs-Trigger:** permanent.

### Aktive Adaptionen

<!-- Eine Zeile je Datei in harness/conventions/. Geltungsbereich und
     Ersetzt-Baseline-Regel stehen hier, damit ein Agent ohne Öffnen
     entscheiden kann, ob der Eintrag ihn betrifft.

     Das <a id="mr-<NNN>"> in der MR-Zelle ist die Adresse, unter der andere
     Dateien diese Adaption referenzieren: conventions.md#mr-<NNN>. Es steht
     hier und nicht in der Eintrags-Datei, weil die Datei bei Auflösung nach
     conventions/done/ wandert und ein Pfad-Link dabei bricht — die Zeile
     wechselt nur die Tabelle, der Anker reist mit. Der Anker trägt die
     Kennung, nie den Titel: Titel werden umformuliert, Kennungen nicht.

     Kam dieses Repo von der Inline-Form (### MR-NNN — Titel), trägt die
     Zeile den alten Überschriften-Slug als ZWEITEN Anker daneben, sonst
     rotten die bereits veröffentlichten Verweise. -->

| MR | Titel | Geltungsbereich | Ersetzt-Baseline-Regel |
|---|---|---|---|
| [\<NNN\>](conventions/MR-<NNN>-<titel>.md) <a id="mr-<NNN>"></a> | <Titel> | <Dateien / Sub-Areas> | <§Abschnitt der Baseline> |

### Aufgelöste Adaptionen

<!-- Eine Zeile je Datei in harness/conventions/done/ — nur ID und
     Nachfolger, damit die Kette auffindbar bleibt, ohne gelesen zu werden.

     Der Anker der Zeile zieht aus der Tabelle oben mit um; er ist der Grund,
     warum ein Verweis auf eine aufgelöste Adaption nicht bricht. -->

| MR | aufgelöst durch |
|---|---|
| [\<NNN\>](conventions/done/MR-<NNN>-<titel>.md) <a id="mr-<NNN>"></a> | [MR-\<NNN\>](conventions/MR-<NNN>-<titel>.md) |

## Zusatzklassen-Deklaration für Sensors-Bindung

<!--
Die vier kanonischen Bindung-Klassen der Sensors-Tabelle in
`harness/README.md` (ADR, Carveout, Schwelle, Reproduzierbarkeit) sind
ohne Deklaration legitim.

Repos können weitere Klassen einführen — z. B. Anforderungs-Bindung
(`LH-...`), Compliance-Bindung (Regulatorik-Artikel), Modell-Version-
Bindung (für KI-Evals). Diese müssen hier deklariert werden, sonst sind
sie für Reviewer nicht von Tippfehlern unterscheidbar.

Eine nicht-deklarierte Zusatzklasse in der Sensors-Tabelle ist eine
stille Setzung und damit Harness-Lüge in derselben Klasse wie ein
halluziniertes Gate (Modul 13).
-->

| Klasse | Form | Bedeutung | Beispiel |
|---|---|---|---|
| <z. B. LH-Bindung> | `LH-<...>` | <z. B. Gate prüft eine bestimmte LH-Anforderung> | <z. B. `LH-QA-01` für Determinismus-Gate> |

<!-- Wenn keine Zusatzklassen verwendet werden: Tabelle entfernen oder
"— keine —" eintragen. -->

## Modus-Deklaration pro Sub-Area

<!--
Pro Modul / Verzeichnis / Sub-Area: Modus festlegen.
- Greenfield (GF): Doc führt, Code folgt. Steady-State.
- Brownfield (BF): Code führt, Doc folgt. Übergangsmodus mit
  Konvergenz-Auftrag zu GF. Graduation-Bedingung benennen.
- Hybrid: gemischt pro Sub-Sub-Area.
- Permanent-BF (selten): nur für Code, der absehbar entfernt wird;
  mit Begründung und Folge-Slice analog zu permanentem Carveout.

Eine Sub-Area in BF *ohne* Graduation-Plan ist eine Harness-Lüge:
"permanente Ausnahme als temporär getarnt" (Modul 7 Analogie).
-->

| Sub-Area (Pfad / Modul) | Modus | Begründung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|
| `*` (Default für gesamtes Repo) | <Greenfield / Brownfield / Hybrid> | <warum> | <Bedingung oder "n/a (GF)" oder "permanent + slice-Ref"> |

## Glossar (optional)

<!--
Repo-spezifische Begriffe, die in den Kernbegriffen des
Baseline-Regelwerks nicht stehen. Nur ergänzen, nicht wiederholen.
-->

| Begriff | Bedeutung |
|---|---|
| <repo-spezifischer Begriff> | <Bedeutung in diesem Repo> |
