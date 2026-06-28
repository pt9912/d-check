# ADR-0021 — matrix: Verweisrichtung innerhalb einer geordneten Klasse über `order`/`direction` (Glob-Rang, fail-closed)

**Status:** Accepted
**Datum:** 2026-06-28
**Autor:** pt9912
**Bezug:** [`DC-FA-MTX-002`](../../../spec/lastenheft.md#dc-fa-mtx-002--verweisrichtung-innerhalb-einer-geordneten-dokumentklasse-modul-matrix)
(Modul `matrix`); kodiert intern dieselbe Richtung wie
[`MR-006`](../../../harness/conventions.md#mr-006--referenzrichtung-spec-straten-verweisen-nie-abwärts-auf-adrs)
(Spec-Straten verweisen nie abwärts), die bislang nur **zwischen** Klassen
(`spec-straten → adr/slice`) maschinell war.
**Schärft:** [`spec/spezifikation.md` §DC-FA-MTX-001.a](../../../spec/spezifikation.md#dc-fa-mtx-001a--klassen--und-status-auflösung)
(neuer Schritt: klasseninterne Rang-/Richtungsprüfung) sowie Config-Schema
(`matrix.classes[].order`/`.direction`) und Grund-Code (`matrix-downward`).

## Kontext

Das Modul `matrix` ([`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix))
prüft Referenzrichtungen **zwischen** deklarierten Klassen über `{from, to,
allow}`-Paare. Die Source-Precedence-Schichtung **innerhalb** eines Stratums —
`architecture` → `spezifikation` → `lastenheft`, Verweise nur aufwärts zur
autoritativeren Schicht — war damit nicht ausdrückbar: Alle drei Spec-Dateien
liegen in **einer** Klasse `spec-straten`, und eine klasseninterne Referenz hat
keine Regel.

Der naheliegende Workaround — `spec-straten` in drei Einzelklassen auflösen und
die Richtung als paarweise `allow: false`-Regeln ausschreiben — scheitert an
zwei Transparenz-Brüchen: (1) `classOf` ordnet eine Datei der **ersten**
passenden Klasse zu (First-Match); kombinierte Klasse **und** Einzelklassen für
dieselben Dateien schließen sich aus, die später deklarierten Regeln feuern nie
(stille tote Regel, grün trotz nicht-greifender Politik). (2) Selbst wenn es
feuerte, ist eine Rangordnung, verstreut über *n·(n−1)/2* flache Paare, für
einen Leser **nicht als Richtung erkennbar**. Hinzu kommt der Konsumenten-Bedarf:
ein reales Spec-Verzeichnis (d-migrate: 23 Dateien) lässt sich nicht als Liste
einzelner Dateien ordnen.

## Entscheidung

Eine Klasse trägt optional `order` (Liste von **Pfad-Globs**, autoritativste
Schicht zuerst) und `direction` (Politik). Der **Rang** einer Datei ist der
Index des ersten `order`-Globs, den sie matcht (First-Match wie `classOf`,
glob-fähig → eine Schicht fasst viele Dateien). Bei `direction: no-downward`
erzeugt eine **klasseninterne** Referenz von Rang *i* auf Rang *j > i* einen
Befund `matrix-downward`; der transitive Sprung ist automatisch erfasst (keine
Paar-Aufzählung). Rangfreie Mitglieder (kein `order`-Treffer) nehmen nicht teil.

Fehlkonfiguration ist **fail-closed**: unbekannter `direction`-Wert, `direction`
ohne `order` und `order` ohne `direction` sind Konfigurationsfehler (Exit 2) —
eine Richtungs-Deklaration darf nicht still wirkungslos sein. Das ist die
direkte Lehre aus dem toten-Regel-Bruch oben: lieber laut scheitern als grün
lügen. Default (beide Felder leer) ist byte-identisch zum bisherigen Verhalten
([`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)). Verteilung
im gepinnten Image, Konsum über `doc-check` — kein kopiertes Skript.

## Verglichene Alternativen

| Alternative | Pro | Contra |
| --- | --- | --- |
| **`order`/`direction` an der Klasse, Glob-Rang (gewählt)** | Gruppe bleibt intakt + lesbar; feuert (keine Verschattung); ein Eintrag = die Absicht; transitiv automatisch; glob-fähig → viele Dateien je Schicht | neues Klassen-Feld + Algorithmus-Schritt |
| Einzelklassen + paarweise `allow:false`-Regeln | nutzt vorhandene `rules`-Maschinerie | First-Match verschattet die Klassen → tote Regeln (stilles Grün); Richtung über *n·(n−1)/2* Paare nicht erkennbar; 23-Datei-Listing unwartbar |
| `order` als Liste **einzelner Dateien** (statt Globs) | maximal explizit | generalisiert nicht auf vielfiles-`spec/` (d-migrate: 23 Zeilen) |
| Mehrfach-Klassenzugehörigkeit in `classOf` | spart das neue Feld | ändert die [`DC-FA-MTX-001`](../../../spec/lastenheft.md#dc-fa-mtx-001--referenzmatrix-zwischen-dokumentklassen-modul-matrix)-Kernsemantik (eine Datei, eine Klasse) für alle Module-Nutzer |
| Default-on / `direction` ohne `order` als No-op tolerieren | weniger Config-Fehler | reintroduziert genau die stille Wirkungslosigkeit, die der Auslöser war |

## Fitness Function

- Ohne `order`/`direction` ist der Befundsatz byte-identisch (Default-aus-Selbsttest, [`DC-QA-02`](../../../spec/lastenheft.md#dc-qa-02--determinismus)).
- Klasseninterner Abwärtsverweis (höher- → niederrangig, auch transitiv) ⇒ genau ein `matrix-downward`; aufwärts ⇒ kein Befund.
- Rangfreies Mitglied (kein `order`-Treffer) ⇒ kein `matrix-downward`.
- Fail-closed-Config: `order` ohne `direction`, `direction` ohne `order`, unbekannter `direction`-Wert ⇒ Konfigurationsfehler (Exit 2), kein stiller No-op.
- Read-only/netzlos unverändert ([`DC-QA-03`](../../../spec/lastenheft.md#dc-qa-03--seiteneffektfreiheit-und-netzwerk-sparsamkeit)).
- Dogfooding: `spec-straten` in [`.d-check.yml`](../../../.d-check.yml) trägt `order`/`direction`; ein eingebauter Abwärtsverweis würde `make doc-check` röten (Negativ-Probe).

## Geschichte

| Datum | Ereignis |
| --- | --- |
| 2026-06-28 | Entwurf + Annahme mit der slice-050-Closure: ausgelöst durch den toten-Regel-Bruch der additiven Variante (Negativ-Probe grün statt rot) und den d-migrate-Konsumenten-Bedarf (23 Spec-Dateien). Modul-Erweiterung implementiert + getestet, `make gates` grün. Status Accepted. |
