# a-check.mk — Architektur-Gate via a-check (Schwester-Tool), zum `include`
# ins Makefile. Erzeugt aus `a-check --print-mk` und an die Repo-Politik
# angepasst (ADR-0029): Digest-Pin-Politik wie alle Gate-Images (ADR-0011),
# Lauf netzlos + read-only (DC-QA-03). Pin-Hebung ist ein bewusster Commit;
# dabei das Fragment per --print-mk neu erzeugen (das Makefile-Target
# arch-check delegiert hierher und bleibt unberührt).
#
# A_CHECK_VERSION steht als eigene Variable, nicht als Prosa im Kommentar:
# die Version IST der Vergleichsgegenstand der Frische-Achse, und was nur im
# Kommentar steht, kann kein Sensor lesen. Die Referenz fuehrt beides — Tag
# UND Digest, wie die Dockerfile-Stages: der Tag macht die Version les- und
# vergleichbar, gezogen wird trotzdem nach Digest.
#
# Das v0.8.0-Release liefert die drei Vorbedingungen des Architektur-Gates
# (tech.adapter-Liste, composition_root: forbid, exclude).
A_CHECK_VERSION ?= v0.8.0
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check:$(A_CHECK_VERSION)@sha256:a1c9c4d6ae3b9690250c6f7271f87b6bb7d5e8d207386fed35ff064508db8e96

.PHONY: a-check
a-check: ## Architektur: Hexagon-Regeln via a-check (netzlos, read-only).
	docker run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) /src
