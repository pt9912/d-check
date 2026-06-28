# ADR-0022 — matrix: Token-basierte Referenz-Richtung, Provenance-Marker, Grandfathering

**Status:** Accepted
**Datum:** 2026-06-28
**Autor:** pt9912
**Bezug:** [`DC-FA-MTX-003`](../../../spec/lastenheft.md#dc-fa-mtx-003--token-basierte-referenz-richtung-mit-provenance-marker-modul-matrix)
(Modul `matrix`); mechanisiert
[§Referenz-Richtung (SDP)](../../../.harness/baseline/v1.4.0/regelwerk/grundlagen-konventionen.md#referenz-richtung-sdp-wer-darf-wen-referenzieren)
des adoptierten Regelwerks, die das Lab bewusst dem Reviewer überließ; setzt die
`matrix`-Mechanik aus [ADR-0021](0021-matrix-klasseninterne-verweisrichtung.md)
fort (dieselbe Linie wie [`MR-006`](../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)/`matrix`).
**Schärft:** [`spec/spezifikation.md` §DC-FA-MTX-001.a](../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)
(neuer Token-/Marker-/Grandfathering-Schritt) sowie Config-Schema
(`matrix.classes[].token`, `matrix.exempt-paths`).

## Kontext

`matrix` ([`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
prüft Referenz-Richtungen nur über **Markdown-Links**: Link-Ziel → Datei →
Klasse. Eine Referenz kann aber auch als **bare ID-Token** im Fließtext stehen
(eine Slice-Kennung in einem ADR-Körper) — die sieht der Link-Scan nicht, und es
gibt für Slice-Kennungen kein `ids`-Linkpflicht-Muster, sie bleiben also freier
Text. Das adoptierte Regelwerk benennt diese Lücke und entscheidet sie für sein
Lab **zugunsten des Reviewers**: ein nackter `slice-NNN`-Grep im ADR-Körper
würde legitime Verifikations-Zeiger („verifiziert in <Slice>") falsch-positiv
flaggen, also nicht grep-bar → Reviewer-Sache. d-checks Identität ist jedoch das
Gegenteil: **mechanisieren, wo die Baseline beim Menschen bleibt** (es hat schon
[`MR-006`](../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs) als minimalen Lab-Grep zu `matrix` mechanisiert). Die Frage ist daher,
ob sich die semantische Unterscheidung doch grep-bar machen lässt.

## Entscheidung

Drei zusammenhängende Mechaniken, alle opt-in/default-aus (byte-identisch ohne
Konfiguration, [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)):

1. **Token-Erkennung.** Eine Klasse kann ein `token`-Regex tragen. `matrix`
   scannt den Körper jedes klassifizierten Dokuments (Prosa außerhalb Fences,
   außerhalb `exclude-sections`, außerhalb Markdown-Links) nach `token` *anderer*
   Klassen; ein Treffer ist eine Token-Referenz, auf die dieselbe
   `{from,to,allow}`-Regel greift. Verbotene Kante ⇒ `matrix-forbidden`
   (Token-Form, kein neuer Grund-Code).
2. **Provenance-Marker.** Der Marker `<!-- d-check:status-provenance -->` auf
   derselben Zeile nimmt eine verbotene Token-Referenz aus. Er verschiebt die
   Unterscheidung von „**Bedeutung**" (nicht grep-bar) auf „**deklariert?**"
   (grep-bar) — das löst genau den Regelwerk-Einwand auf. Es ist `matrix`' erster
   Zeilen-Marker; bewusst **benannt/strukturiert** (näher an
   `allow-supersede-lineage` als an einem generischen Stummschalter), und kehrt
   die „nur strukturelle Ausnahmen"-Haltung damit eng begrenzt um. Was der Marker
   **nicht** grep-bar macht — ob die Deklaration *ehrlich* ist (echte Provenance
   vs. getarnte Entscheidungsgrundlage) — bleibt Reviewer-Backstop.
3. **Grandfathering per `exempt-paths`.** Eine Datei, die `matrix.exempt-paths`
   matcht, wird ganz übersprungen. Bereits `Accepted`-ADRs sind immutabel
   (`make adr-check`) und können **nicht** nachträglich markiert werden; sie
   werden grandfathered (Regelwerk: „Gate prüft nur ab Einführung neu"). Kein
   Eingriff in `adr-check`, kein Editieren immutabler ADRs.

`matrix-forbidden` deckt damit beide Formen; **Links** bleiben markerlos
(legitime Provenance ist ein stabiler Token, kein Datei-Link;
`exclude-sections` deckt Provenance unter `## Geschichte`/Historie ab).

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **Token-Erkennung + Provenance-Marker + Grandfathering (gewählt)** | macht die SDP-Richtung vollständig mechanisch (Link **und** Token); Marker macht die Semantik grep-bar; immutable ADRs unberührt | erster `matrix`-Zeilen-Marker; neue ADRs müssen Verifikations-Zeiger deklarieren |
| Nur Reviewer (Baseline-Wahl) | kein Code, minimal | lässt die Token-Form unmechanisiert — gegen d-checks Mechanisier-Identität |
| Pauschaler `slice-NNN`-Grep ohne Marker | einfachste Mechanik | flaggt legitime Verifikations-Zeiger falsch-positiv (genau die Regelwerk-Warnung) |
| Immutable ADRs nachträglich markieren | volle Abdeckung | braucht `adr-check`-Eingriff (Kommentar-Strip) + churnt 21 eingefrorene ADRs |
| Generischer `d-check:ignore` für `matrix` | nutzt bestehenden Marker | unbenanntes Stummschalten — gegen `matrix`' „strukturell statt zeilenweise"-Linie |

## Fitness Function

- Ohne `token`/`exempt-paths` ist der Befundsatz byte-identisch (Default-aus-Selbsttest, [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Unmarkierter verbotener Token im Körper ⇒ genau ein `matrix-forbidden`; mit Marker auf der Zeile ⇒ kein Befund; in `exclude-sections`/grandfatherter Datei ⇒ kein Befund.
- Token in Markdown-Links und in Fences zählen nicht (Link-Scan bzw. Fence-Opazität).
- Read-only/netzlos unverändert ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Dogfood: `slice`-Klasse trägt `token: 'slice-\d{3}'`, Regel `{from: adr, to: slice}` und `exempt-paths` für die Alt-ADRs in [`.d-check.yml`](../../../.d-check.yml). Dieser Beleg-Verweis auf den umsetzenden Slice slice-051 <!-- d-check:status-provenance --> trägt selbst den Marker und bleibt befundfrei.

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-28 | Entwurf + Annahme mit der slice-051-Closure: Design in der Auftraggeber-Session (Token-Erkennung statt Link-only; Marker macht die Semantik grep-bar; Grandfathering statt Immutable-Edit). Modul-Erweiterung implementiert + getestet, zwei unabhängige Reviews, `make gates` grün. Status Accepted. |
