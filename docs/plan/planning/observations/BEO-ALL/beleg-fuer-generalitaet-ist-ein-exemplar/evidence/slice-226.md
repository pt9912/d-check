`slice-226`s Pack-Alias-Mechanismus ist für „jeden Fremdpräfix" konstruiert
(Hash-Suffix + passender Index, präfix-unabhängig), aber nur an einem
Exemplar beobachtet: dem `loose-`-Präfix aus `git
maintenance run --task=loose-objects`, dem konkreten Anlass des eingehenden
CR. Der Slice-Plan (§6) benennt das selbst und zieht ausdrücklich den
Präzedenzfall [ADR-0072](../../../../../adr/0072-workflows-modul.md) heran, der
dieselbe Evidenzlücke für die yaml-Kapsel-Erweiterung akzeptiert.
