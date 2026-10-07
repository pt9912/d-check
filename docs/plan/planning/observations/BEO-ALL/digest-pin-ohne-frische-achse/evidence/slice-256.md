**Vorgang:** slice-256
**Fund:** `moby/buildkit` (Builder von `setup-buildx-action`) und
`tonistiigi/binfmt` (`setup-qemu-action`) sind in `release.yml` digest-gepinnt;
Dependabot hebt nur die `uses:`-Referenzen, der Nachtlauf kennt die beiden
Pins nicht. Benannt in [ADR-0102](../../../../../adr/0102-multi-arch-index-und-spiegel-per-index-digest.md) §Konsequenzen; `make versions` gibt sie aus.
