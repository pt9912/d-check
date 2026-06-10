// Package configyaml ist der driven Config-Adapter: striktes Decoding
// und statische Validierung von .d-check.yml (DC-FA-CONF-001;
// Strategie gemäß ADR-0003). Den Datei-Inhalt beschafft das CLI über
// den Filesystem-Adapter — dieser Adapter macht kein I/O
// (spec/architecture.md §2). Die Existenz-Prüfung deklarierter
// Scan-Wurzeln erfolgt im Kern beim Lauf-Start (ebenfalls Exit 2).
package configyaml

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/d-check/internal/hexagon/core"
)

// FileName ist der feste Name der Konfigurationsdatei.
const FileName = ".d-check.yml"

// raw bildet das Voll-Schema von .d-check.yml ab
// (spec/spezifikation.md §2). Unbekannte Schlüssel sind durch
// KnownFields(true) Fehler.
type raw struct {
	Scan *struct {
		Roots  []string `yaml:"roots"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"scan"`
	Modules []string `yaml:"modules"`
	IDs     *struct {
		Patterns []struct {
			Regex  string `yaml:"regex"`
			Target string `yaml:"target"`
		} `yaml:"patterns"`
	} `yaml:"ids"`
	Matrix *struct {
		Classes []struct {
			Name  string   `yaml:"name"`
			Paths []string `yaml:"paths"`
		} `yaml:"classes"`
		Rules []struct {
			From  string `yaml:"from"`
			To    string `yaml:"to"`
			Allow bool   `yaml:"allow"`
		} `yaml:"rules"`
		Status *struct {
			Forbidden []string `yaml:"forbidden"`
		} `yaml:"status"`
		ExcludeSections []string `yaml:"exclude-sections"`
	} `yaml:"matrix"`
	External *struct {
		TimeoutSeconds int `yaml:"timeout-seconds"`
		Parallel       int `yaml:"parallel"`
	} `yaml:"external"`
}

// Decode parst und validiert den Datei-Inhalt vollständig — Syntax
// und statische Semantik; jeder Fehler bedeutet Exit 2 ohne Prüfung.
// content == nil bedeutet: keine Konfigurationsdatei → Defaults.
func Decode(content []byte) (core.Config, error) {
	var cfg core.Config
	if content == nil {
		return cfg, nil
	}
	var r raw
	dec := yaml.NewDecoder(bytes.NewReader(content))
	dec.KnownFields(true)
	if err := dec.Decode(&r); err != nil {
		return cfg, fmt.Errorf("%s: %w", FileName, err)
	}

	if r.Scan != nil {
		cfg.Roots = r.Scan.Roots
		cfg.Ignore = r.Scan.Ignore
	}
	cfg.Modules = r.Modules

	if r.IDs != nil {
		for i, p := range r.IDs.Patterns {
			re, err := regexp.Compile(p.Regex)
			if err != nil {
				return cfg, fmt.Errorf("%s: ids.patterns[%d].regex: %v", FileName, i, err)
			}
			if p.Target == "" {
				return cfg, fmt.Errorf("%s: ids.patterns[%d].target fehlt", FileName, i)
			}
			cfg.IDPatterns = append(cfg.IDPatterns, core.IDPattern{Regex: re, Target: p.Target})
		}
	}
	if r.Matrix != nil {
		classes := map[string]bool{}
		for i, c := range r.Matrix.Classes {
			if c.Name == "" || classes[c.Name] {
				return cfg, fmt.Errorf("%s: matrix.classes[%d].name fehlt oder doppelt", FileName, i)
			}
			classes[c.Name] = true
		}
		for i, rule := range r.Matrix.Rules {
			if !classes[rule.From] || !classes[rule.To] {
				return cfg, fmt.Errorf("%s: matrix.rules[%d] referenziert undeklarierte Klasse", FileName, i)
			}
		}
	}
	if r.External != nil {
		if t := r.External.TimeoutSeconds; t != 0 && (t < 1 || t > 300) {
			return cfg, fmt.Errorf("%s: external.timeout-seconds außerhalb 1–300", FileName)
		}
		if p := r.External.Parallel; p != 0 && (p < 1 || p > 16) {
			return cfg, fmt.Errorf("%s: external.parallel außerhalb 1–16", FileName)
		}
	}
	return cfg, nil
}
