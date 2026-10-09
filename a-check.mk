# a-check.mk — Architektur-Gate via a-check (Schwester-Tool), zum `include`
# ins Makefile. Erzeugt aus `a-check --print-mk` und an die Repo-Politik
# angepasst (ADR-0029): Digest-Pin-Politik wie alle Gate-Images (ADR-0011),
# Lauf netzlos + read-only (DC-QA-03). Pin-Hebung ist ein bewusster Commit;
# dabei das Fragment per --print-mk neu erzeugen (das Makefile-Target
# arch-check delegiert hierher und bleibt unberührt).
#
# ZWEI TEILE DES v0.23.0-FRAGMENTS SIND BEWUSST NICHT ADOPTIERT (beide lagen
# schon im v0.17.0-Fragment; das v0.23.0-Fragment ist byte-gleich zum v0.20.0-Fragment), damit die
# Anweisung oben nicht als unbelegte Zusage dasteht:
#   DOCKER ?= docker   Eine Runtime-Indirektion zahlt sich nur repo-weit aus;
#                      die uebrigen Rezepte dieses Repos rufen `docker` hart.
#                      Sie hier allein einzufuehren erzeugte genau die
#                      Halb-und-halb-Lage, gegen die sie gebaut ist — und
#                      Docker ist in AGENTS.md §3.1 als Voraussetzung gesetzt,
#                      nicht als Wahl.
#   a-check-graph      ABGELEHNT, nicht aufgeschoben. Drei Gruende, zwei
#                      gemessen: (a) es fehlt kein Bild — spec/architecture.md
#                      traegt bereits zwei Mermaid-Diagramme; (b) der Graph
#                      entsteht aus .a-check.yml und zeigt damit Modul-Pfade,
#                      die AGENTS.md §3.4 in genau der Sicht verbietet, in die
#                      ein Architektur-Bild gehoert (Verschaerfung MR-033) —
#                      anderswo abgelegt waere er ein Bild ohne Ort; (c) er
#                      koennte nur zeigen, was arch-check ohnehin erzwingt, ist
#                      also redundant oder falsch. Dazu der Preis: ein Target
#                      ist dauerhaft gate-consistency-pflichtig (AGENTS.md §4,
#                      harness/README.md §Sensors) — fuer etwas, das kein Gate
#                      faehrt. Eine erneute Vorlage braucht ein Argument, das
#                      diese drei entkraeftet, nicht nur eine neue Version.
# DOCKER bleibt ein benannter Kandidat; a-check-graph ist entschieden. Ein
# Kandidat, der zwei Fragment-Versionen lang unangetastet dasteht, IST eine
# Entscheidung — sie war nur nicht aufgeschrieben.
#
# A_CHECK_VERSION steht als eigene Variable, nicht als Prosa im Kommentar:
# die Version IST der Vergleichsgegenstand der Frische-Achse, und was nur im
# Kommentar steht, kann kein Sensor lesen. Die Referenz fuehrt beides — Tag
# UND Digest, wie die Dockerfile-Stages: der Tag macht die Version les- und
# vergleichbar, gezogen wird trotzdem nach Digest.
#
# Die drei Vorbedingungen des Architektur-Gates (tech.adapter-Liste,
# composition_root: forbid, exclude) kamen mit v0.8.0 und tragen weiter.
# Vor der Hebung auf v0.23.1 gemessen: derselbe Lauf ueber dieses Repo, 0
# Befunde in beiden Fassungen, und beide melden denselben konstruierten
# Verstoss (app-impurity) an derselben Zeile; --print-mk ist unveraendert.
# Der opt-in-Block shapes (seit v0.21.0) ist hier nicht konfiguriert.
# Der BREAKING Change aus v0.20.0 — ein
# Richtungssegment am Ende eines port-Globs (z. B. `.../ports/outbound/**`)
# schaltet `port-locality` nicht mehr still ab — trifft dieses Repo nicht:
# der einzige port-Glob (`internal/hexagon/port/**`) endet nicht auf einem
# Richtungssegment, und die ports-Schicht in `.a-check.yml` fuehrt ohnehin
# kein `direction`-Feld, die einzige Vorbedingung, unter der die Aenderung
# greift. Gemessen, nicht aus dem Changelog geschlossen.
A_CHECK_VERSION ?= v0.23.1
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check:$(A_CHECK_VERSION)@sha256:4948e1e45a595750e6fe49ef19acde42fbfae1cd8d84ea0cca3ae0611443665d

.PHONY: a-check
a-check: ## Architektur: Hexagon-Regeln via a-check (netzlos, read-only).
	docker run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) /src
