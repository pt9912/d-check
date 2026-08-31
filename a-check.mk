# a-check.mk — Architektur-Gate via a-check (Schwester-Tool), zum `include`
# ins Makefile. Erzeugt aus `a-check --print-mk` und an die Repo-Politik
# angepasst (ADR-0029): Digest-Pin-Politik wie alle Gate-Images (ADR-0011),
# Lauf netzlos + read-only (DC-QA-03). Pin-Hebung ist ein bewusster Commit;
# dabei das Fragment per --print-mk neu erzeugen (das Makefile-Target
# arch-check delegiert hierher und bleibt unberührt).
#
# ZWEI NEUERUNGEN DES v0.19.0-FRAGMENTS SIND BEWUSST NICHT ADOPTIERT, damit die
# Anweisung oben nicht als unbelegte Zusage dasteht:
#   DOCKER ?= docker   Eine Runtime-Indirektion zahlt sich nur repo-weit aus;
#                      die uebrigen Rezepte dieses Repos rufen `docker` hart.
#                      Sie hier allein einzufuehren erzeugte genau die
#                      Halb-und-halb-Lage, gegen die sie gebaut ist — und
#                      Docker ist in AGENTS.md §3.1 als Voraussetzung gesetzt,
#                      nicht als Wahl.
#   a-check-graph      Ein neues Target ist gate-consistency-pflichtig (AGENTS.md
#                      §4 und harness/README.md §Sensors) und damit ein eigener
#                      Entscheid, kein Nebeneffekt einer Pin-Hebung.
# Beides bleibt ein benannter Kandidat, kein Versehen.
#
# A_CHECK_VERSION steht als eigene Variable, nicht als Prosa im Kommentar:
# die Version IST der Vergleichsgegenstand der Frische-Achse, und was nur im
# Kommentar steht, kann kein Sensor lesen. Die Referenz fuehrt beides — Tag
# UND Digest, wie die Dockerfile-Stages: der Tag macht die Version les- und
# vergleichbar, gezogen wird trotzdem nach Digest.
#
# Die drei Vorbedingungen des Architektur-Gates (tech.adapter-Liste,
# composition_root: forbid, exclude) kamen mit v0.8.0 und tragen weiter.
# Vor der Hebung auf v0.19.0 gemessen: derselbe Lauf ueber dieses Repo, 0
# Befunde — und beide Fassungen melden denselben konstruierten Verstoss
# (app-impurity) an derselben Zeile. Das Gruen kommt also aus dem Bestand,
# nicht aus einer Fassung, die weniger prueft.
# Der BREAKING Change aus v0.18.0 — die Richtungs-Vokabel gilt je Rolle
# (port: inbound/outbound, adapter: driving/driven), die falsche ist Exit 2 —
# trifft dieses Repo nicht: keine der sechs Schichten fuehrt ein
# direction-Feld. Gemessen, nicht aus dem Changelog geschlossen.
A_CHECK_VERSION ?= v0.19.0
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check:$(A_CHECK_VERSION)@sha256:34d3dfb50e44d99ea735186a35e1040589c4681dcfa2a51ed0f2aaea718cdd2d

.PHONY: a-check
a-check: ## Architektur: Hexagon-Regeln via a-check (netzlos, read-only).
	docker run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) /src
