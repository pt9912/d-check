# ADR-0064 — Docker-Hub-Spiegel: derselbe Digest, fail-closed

**Status:** Accepted
**Datum:** 2026-08-27
**Autor:** pt9912
**Bezug:** [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel),
[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image),
[ADR-0002](0002-distribution-ghcr-image.md), [ADR-0014](0014-latest-tag-fuer-stabile-releases.md)
**Schärft:** [`DC-FA-DIST-002`](../../../spec/lastenheft.md#dc-fa-dist-002--docker-hub-spiegel)
— das *Wie* der Spiegelung; das *Was* steht im Lastenheft.

## Kontext

[ADR-0002](0002-distribution-ghcr-image.md) legte GHCR als einzigen
Distributionsweg fest, und
[`DC-FA-DIST-001`](../../../spec/lastenheft.md#dc-fa-dist-001--docker-image)
schloss alles andere aus. Der Auftraggeber verlangt einen zweiten Bezugsweg
über Docker Hub und hat dessen **Rang** ausdrücklich entschieden: eine
**zugesagte** Distribution, kein unverbindlicher Komfort-Spiegel.

Damit sind drei Fragen offen, die das Lastenheft bewusst nicht beantwortet:
**wie** gespiegelt wird (Kopie oder Zweitbau), **woran** die Gleichheit gemessen
wird, und **was** bei einem Fehlschlag geschieht.

## Entscheidung

1. **Spiegeln heißt kopieren, nicht neu bauen.** Dasselbe lokal gebaute Bild
   wird für beide Registries getaggt und gepusht. Ein zweiter Bau erzeugte ein
   zweites Bild — gleiche Quelle, anderer Digest —, und die Zusage
   *„derselbe Digest"* wäre nicht haltbar.
2. **Der Manifest-Digest ist die Prüfgröße, und er wird gemessen.** Nach dem
   Push liest die Pipeline den Digest beider Referenzen und vergleicht sie;
   Ungleichheit bricht ab. Dass `docker push` aus demselben Bild denselben
   Digest erzeugt, ist plausibel und wird trotzdem **nachgesehen** — eine
   Eigenschaft, die eine Anforderung zusagt und niemand prüft, ist ein stiller
   Grün-Pfad.
3. **Fail-closed, und die Meldung nennt den Teil-Zustand.** Der Hub-Push steht
   **nach** dem GHCR-Push; scheitert er, ist der GHCR-Stand bereits
   veröffentlicht. Der Abbruch sagt das ausdrücklich, statt den Betrachter den
   Zustand raten zu lassen.
4. **Reihenfolge GHCR zuerst.** Die Quelle wird zuerst veröffentlicht, der
   Spiegel danach. Umgekehrt trüge Docker Hub kurzzeitig ein Bild, das die
   Quelle noch nicht führt — der Spiegel wäre dann keiner.
5. **Für den Spiegel kein `continue-on-error` und kein Secret-Sondierungs-Schritt.** Fehlt das
   Token, schlägt der Login fehl und das Release ist rot. Das ist die
   Konsequenz aus Punkt 3 und der Unterschied zum fail-open-Muster des
   Schwester-Repos `d-migrate` — dort ist der Spiegel Komfort, hier Zusage.
6. **Bild-Name über `vars` überschreibbar, mit Default.** `DOCKERHUB_IMAGE`
   trägt `pt9912/d-check` als Default; ein Fork spiegelte sonst in fremdes
   Namensgebiet.
7. **Die Darstellung wird aus dem Repo gepflegt — und ist ausdrücklich
   *nicht* fail-closed.** Kurztext und Overview-Seite liegen als Dateien unter
   [`packaging/dockerhub/`](../../../packaging/dockerhub/README.md) und werden
   beim Release gesetzt; die Kategorie hat keinen Action-Input und steht dort
   als **Text**, statt still im Web-UI zu leben. Der Schritt läuft mit
   `continue-on-error`: das **Bild** ist die Zusage, der **Beschreibungstext**
   ist Präsentation. Ein Release rot zu machen, weil ein Kurztext nicht gesetzt
   werden konnte, verwechselte beides. Das Byte-Limit der Description prüft die
   Pipeline dagegen **fail-fast**, statt der Action das stille Abschneiden zu
   überlassen — ein abgeschnittener Text ist ein stiller Defekt, ein
   fehlgeschlagener Upload ein lauter.

## Verglichene Alternativen

| Alternative | Vorteil | Warum nicht |
| ----------- | ------- | ----------- |
| **Fail-open** (Notice statt Abbruch, wie `d-migrate`) | Release bleibt von fremder Verfügbarkeit unabhängig | Der Auftraggeber hat den Rang *zugesagt* gewählt; ein Spiegel, der still zurückfallen darf, ist keine Zusage, sondern eine Gewohnheit |
| **Zweitbau auf Docker Hub** (eigener Build-Schritt) | unabhängig von der GHCR-Pipeline | Zwei Bilder, zwei Digests — die tragende Zusage *„derselbe Digest"* entfiele, und ein `docker.io`-Pin wäre schwächer als ein `ghcr.io`-Pin |
| **`buildx imagetools create`** (Manifest aus GHCR kopieren) | ausdrückliche Kopier-Operation, digest-erhaltend per Definition | zusätzlicher Netz-Umlauf und eine Abhängigkeit von der GHCR-Lesbarkeit unmittelbar nach dem Push; dasselbe lokale Bild liegt ohnehin vor |
| **Registry-seitige Spiegelung** (Docker Hub zieht von GHCR) | keine Pipeline-Änderung | Docker Hub bietet das für diesen Fall nicht als steuerbaren, prüfbaren Schritt; die Gleichheit wäre nicht messbar |
| **Spiegel vor GHCR pushen** | symmetrisch | Der Spiegel führte kurzzeitig ein Bild, das die Quelle nicht hat — Punkt 4 |

## Konsequenzen

**Positiv:** ein zweiter Bezugsweg mit derselben Belastbarkeit wie der erste;
`docker.io`- und `ghcr.io`-Pins sind austauschbar, weil der Digest gleich ist.
Die Gleichheit ist gemessen, nicht behauptet.

**Negativ, und benannt:** das Release bindet an die Verfügbarkeit eines fremden
Dienstes und an ein Token, das ablaufen kann. Ein gestörter Docker Hub macht ein
gültiges, bereits auf GHCR liegendes Release rot. Das ist der Preis der
gewählten Zusage — er wird nicht durch einen Notausgang aufgeweicht, weil ein
Notausgang die Zusage zurücknähme.

**Eine ungedeckte Fläche bleibt:** das Modul `versions` hält
`ghcr`-präfixierte Pins gegen `version.md#aktuell`. Ein `docker.io`-Pin fällt
heute nicht darunter und veraltet still — dieselbe Klasse wie das bare Tag im
Handbuch. Der Slice benennt sie als Abnahme-Punkt; diese ADR schließt sie
nicht, weil die Muster-Fläche des Moduls ein eigener Entscheid ist.

## Re-Evaluierungs-Trigger

Fällt das Release **zweimal** an einer Docker-Hub-Störung aus, während GHCR
grün war, ist die fail-closed-Entscheidung neu zu stellen — dann kostet die
Zusage mehr, als sie trägt, und die ehrliche Antwort wäre eine Rücknahme der
Anforderung, nicht ein stiller Notausgang.
