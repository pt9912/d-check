# a-check.mk — Architektur-Gate via a-check (Schwester-Tool), zum `include`
# ins Makefile. Erzeugt aus `a-check --print-mk` und an die Repo-Politik
# angepasst (ADR-0029): Digest-Pin-Politik wie alle Gate-Images (ADR-0011),
# Lauf netzlos + read-only (DC-QA-03). Pin-Hebung ist ein bewusster Commit;
# dabei das Fragment per --print-mk neu erzeugen (das Makefile-Target
# arch-check delegiert hierher und bleibt unberührt).
#
# A_CHECK_IMAGE ist auf den v0.8.0-Release digest-gepinnt — das Release
# liefert die drei slice-058-Vorbedingungen (tech.adapter-Liste,
# composition_root: forbid, exclude).
A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:a1c9c4d6ae3b9690250c6f7271f87b6bb7d5e8d207386fed35ff064508db8e96

.PHONY: a-check
a-check: ## Architektur: Hexagon-Regeln via a-check (netzlos, read-only).
	docker run --rm --network none -v "$(CURDIR)":/src:ro $(A_CHECK_IMAGE) /src
