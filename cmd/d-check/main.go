// d-check — Doc-Referenz-Checker. Dünner Einstiegspunkt; die
// CLI-Logik lebt im driving Adapter (ADR-0005, u-boot-Konvention).
package main

import (
	"os"

	"github.com/pt9912/d-check/internal/adapter/driving/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
