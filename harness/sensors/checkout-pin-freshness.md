# `make checkout-pin-freshness` u. a. — meldet, ob die drei Action-Pins der Workflows veraltet sind

## Vertrag

Dieselbe Versions-Frage wie die Toolchain-Achsen, anderer Gegenstand: die
**drei Action-Pins** der Workflows. Der Pin ist ein SHA — aber der
**Tag-Kommentar daneben** trägt den Release-Tag, und das ist genau die Größe,
die der `releases/latest`-Redirect liefert; die Achse brauchte deshalb **keine**
neue Quellen-Form, nur einen anderen Extraktor.

**Sie ergänzt [`workflow-pins`](workflow-pins.md), das die Form hält:** Das
Pinnen schließt ein umgehängtes Tag aus und macht zugleich blind dafür, dass
upstream eine Lücke behoben hat. Ein alter Action-Pin ist damit sehr wohl ein
**Sicherheitsthema**, nicht bloß ein Aktualitäts-Wunsch.

## Grenze — was das Grün nicht abdeckt

1. **Verglichen wird der Kommentar gegen upstream, nicht der SHA gegen den
   Kommentar** — jene Gültigkeitsfrage bleibt die benannte Grenze von
   [`workflow-pins`](workflow-pins.md) ([`AGENTS.md`](../../AGENTS.md) §3.9).
2. **Fail-open** — jede fremde Störung endet als `SKIP`.

## Bindung

kein Gate, **nicht** in `gates`; Netz, gerufen vom Nachtlauf.
[`AGENTS.md`](../../AGENTS.md) §3.9
