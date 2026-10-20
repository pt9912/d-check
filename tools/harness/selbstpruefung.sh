#!/usr/bin/env bash
# selbstpruefung.sh — Negativ-Selbsttest des commit-msg-Hooks und des Stop-Hooks,
# adoptiert aus ai-harness-init (Stand v0.10.0) · seit slice-245. EIN KOMMANDO,
# KEIN GATE.
#
# WAS ER BELEGT. Dieses Repo klont sich selbst in ein Wegwerf-Verzeichnis,
# stellt fest, dass der frische Klon in seiner LOKALEN Konfiguration keinen
# core.hooksPath traegt, aktiviert dort den Traeger der Commit-Kennung
# (`make hooks`) und faehrt danach ZWEI Commit-Versuche: einer OHNE Kennung
# faellt und entsteht nicht, einer MIT Kennung geht durch. Zuletzt laeuft das
# Gate-Kommando im Klon, und der Stop-Hook faehrt dort seine Zweige: streng ohne
# Schalter, locker nur mit .harness/stop-gate-locker, streng mit STOP_GATE_STRENG=1
# trotz der Datei, und ein Commit ohne Nachweis blockt im lockeren Zweig. BEIDE
# AUSGAENGE STEHEN IN EINEM LAUF — ein Lauf, der nur den durchgelassenen Commit
# beobachtet, belegt nicht, dass der Traeger ueberhaupt etwas aufhaelt.
#
# WARUM EIN KLON UND NICHT DIESES ARBEITSVERZEICHNIS. Die Aktivierung ist
# lokale Konfiguration (core.hooksPath) und reist nicht mit dem Klon; genau
# dieser Zustand ist der Pruefgegenstand. Und die zwei Commit-Versuche
# schreiben Historie — im Klon ist sie nach dem Lauf weg.
#
# DIE GRENZE. Geprueft sind der Traeger und die zwei Commit-Ausgaenge. NICHT
# geprueft ist, ob der Traeger jeden Commit-Pfad erreicht: `git commit
# --no-verify` umgeht ihn, Commits aus Repo-Werkzeugen tragen ihre eigenen
# Messages, und Aufrufformen ausserhalb der aktivierten Traeger-Form bleiben
# ausserhalb der Zusage. Geprueft ist ferner die ANWESENHEIT einer Kennung in
# der Message, nicht ihre Wahrheit.
#
# VORAUSSETZUNG. Ein git-Repo mit mindestens einem Commit. Docker muss laufen
# (das Gate-Kommando im Klon baut das Bild). Kein Netz.
set -euo pipefail

SELBSTPRUEFUNG_TRAEGER="${SELBSTPRUEFUNG_TRAEGER:-.githooks/commit-msg}"
SELBSTPRUEFUNG_AKTIVIERUNG="${SELBSTPRUEFUNG_AKTIVIERUNG:-make hooks}"
SELBSTPRUEFUNG_GATE="${SELBSTPRUEFUNG_GATE:-make gates}"
SELBSTPRUEFUNG_MSG_ROT="${SELBSTPRUEFUNG_MSG_ROT:-Selbstpruefung ohne Kennung}"
SELBSTPRUEFUNG_MSG_GRUEN="${SELBSTPRUEFUNG_MSG_GRUEN:-Selbstpruefung mit Kennung slice-244}"

fehler() {
	echo "selbstpruefung: FEHLER — $1" >&2
	exit 1
}

command -v git >/dev/null 2>&1 || fehler "git liegt nicht im Pfad — ohne git gibt es keinen Klon, auf dem ein Traeger etwas aufhalten koennte."

quelle="$(git rev-parse --show-toplevel 2>/dev/null)" || fehler "dieses Verzeichnis liegt in keinem git-Repo — die Pruefung klont das Repo, in dem sie lauft."
git -C "$quelle" rev-parse --verify --quiet 'HEAD^{commit}' >/dev/null 2>&1 ||
	fehler "$quelle traegt keinen Commit — ein Klon haette keinen Baum, auf dem ein Commit-Versuch etwas bedeutete. Ein erster Commit macht die Pruefung fahrbar."

echo "selbstpruefung: Marker — Traeger=[$SELBSTPRUEFUNG_TRAEGER] Aktivierung=[$SELBSTPRUEFUNG_AKTIVIERUNG] Gate=[$SELBSTPRUEFUNG_GATE] MsgRot=[$SELBSTPRUEFUNG_MSG_ROT] MsgGruen=[$SELBSTPRUEFUNG_MSG_GRUEN]"
echo "selbstpruefung: Quelle=[$quelle] (der Klon traegt den Stand von HEAD)"

echo "selbstpruefung: Traeger, Aktivierung, die zwei Commit-Ausgaenge, das Gate-Kommando und der Stop-Hook im frischen Klon ..."

arbeit="$(mktemp -d -p "${TMPDIR:-/tmp}")"
trap 'rm -rf "$arbeit"' EXIT
# mktemp -d liefert 0700. Laeuft das Gate-Kommando in einem Container als
# Nicht-Root, kann er einen 0700-Pfad nicht traversieren.
chmod 755 "$arbeit"
klon="$arbeit/klon"

# `file://` ist Pflicht: ueber einen blossen Pfad legt git einen Klon mit
# geteilten Objekten an, und der haengt am Quell-Repo statt fuer sich zu stehen.
git clone -q "file://$quelle" "$klon" || fehler "der Klon von $quelle ist nicht entstanden."
chmod 755 "$klon"

# (1) DER FRISCHE KLON IST UNGEPRUEFT. Gelesen wird die LOKALE Konfiguration
# des Klons — sie ist es, die ein Klon nicht erbt. `--get` ohne Scope naehme
# auch eine globale oder systemweite Setzung mit und meldete einen Zustand des
# Rechners als einen des Klons.
vorher_hooks="$(git -C "$klon" config --get --local core.hooksPath || true)"
if [ -n "$vorher_hooks" ]; then
	fehler "der frische Klon traegt in seiner lokalen Konfiguration bereits core.hooksPath=$vorher_hooks — dann misst der Lauf nicht, was die Aktivierung bewirkt."
fi
echo "selbstpruefung: der frische Klon traegt lokal keinen core.hooksPath — der Traeger reist mit, seine Aktivierung nicht."

if [ ! -f "$klon/$SELBSTPRUEFUNG_TRAEGER" ]; then
	fehler "im Klon liegt kein Traeger unter [$SELBSTPRUEFUNG_TRAEGER] — er ist entweder nicht committet oder liegt woanders; der Marker SELBSTPRUEFUNG_TRAEGER nennt den Pfad."
fi

# (2) DIE AKTIVIERUNG. Gelesen wird der Wert, den git danach lokal fuehrt.
akt_out=""
if ! akt_out="$(cd "$klon" && bash -c "$SELBSTPRUEFUNG_AKTIVIERUNG" 2>&1)"; then
	printf '%s\n' "$akt_out" >&2
	fehler "der Aktivierungsschritt [$SELBSTPRUEFUNG_AKTIVIERUNG] endet nicht mit Exit 0 — der Traeger ist im Klon nicht in Betrieb zu nehmen."
fi
printf '%s\n' "$akt_out" | sed 's/^/selbstpruefung:   /'
nachher_hooks="$(git -C "$klon" config --get --local core.hooksPath || true)"
if [ -z "$nachher_hooks" ]; then
	fehler "nach [$SELBSTPRUEFUNG_AKTIVIERUNG] ist core.hooksPath im Klon weiter leer — der Schritt meldet Erfolg, ohne einen zu haben."
fi
echo "selbstpruefung: aktiviert — core.hooksPath=$nachher_hooks"

# (2b) DER TRAEGER-MARKER NENNT DIE DATEI, DIE AUFHAELT — nicht irgendeine.
# git ruft einen Hook ueber seinen nackten Namen aus core.hooksPath; welche
# Datei das ist, steht damit fest. Ohne diesen Abgleich lenkte der Marker nur
# eine Existenzpruefung, und die Schluss-Zeile schriebe ihm eine Wirkung zu,
# die eine andere Datei hatte.
hook_name="${SELBSTPRUEFUNG_TRAEGER##*/}"
case "$nachher_hooks" in
/*) gerufen="$nachher_hooks/$hook_name" ;;
*) gerufen="$klon/$nachher_hooks/$hook_name" ;;
esac
if [ ! -e "$gerufen" ] || [ ! "$gerufen" -ef "$klon/$SELBSTPRUEFUNG_TRAEGER" ]; then
	fehler "[$SELBSTPRUEFUNG_TRAEGER] ist nicht der Traeger, den der Aktivierungsschritt in Betrieb nimmt: [$SELBSTPRUEFUNG_AKTIVIERUNG] setzt core.hooksPath=$nachher_hooks, und git ruft daraus $gerufen. Der Marker nennt damit eine andere Datei, als aufhaelt."
fi
echo "selbstpruefung: in Betrieb ist [$SELBSTPRUEFUNG_TRAEGER] — git ruft ihn als $nachher_hooks/$hook_name."

# (3) DIE ZWEI COMMIT-AUSGAENGE, in einem Lauf. `--allow-empty` haelt jeden
# Versuch ohne Baum-Aenderung; gelesen werden Exit-Code UND die Lage von HEAD,
# denn eine Meldung allein belegt keinen Abbruch.
vorher_head="$(git -C "$klon" rev-parse HEAD)"
for fall in rot gruen; do
	case "$fall" in
	rot) msg="$SELBSTPRUEFUNG_MSG_ROT" ;;
	gruen) msg="$SELBSTPRUEFUNG_MSG_GRUEN" ;;
	*) fehler "unbekannter Fall [$fall]." ;;
	esac
	out=""
	rc=0
	out="$(git -C "$klon" -c user.email=selbstpruefung@example.invalid -c user.name=selbstpruefung \
		commit -q --allow-empty -m "$msg" 2>&1)" || rc=$?
	if [ "$fall" = rot ]; then
		if [ "$rc" -eq 0 ]; then
			fehler "der Commit '$msg' geht durch — er sollte am Traeger [$SELBSTPRUEFUNG_TRAEGER] fallen. Fuehrt dieses Repo dort seinen EIGENEN Traeger (skip-if-present), nennt SELBSTPRUEFUNG_MSG_ROT eine Message, die jener aufhaelt."
		fi
		nun_head="$(git -C "$klon" rev-parse HEAD)"
		if [ "$nun_head" != "$vorher_head" ]; then
			fehler "der Traeger meldet den Abbruch, der Commit ist aber entstanden (HEAD bewegt sich von $vorher_head nach $nun_head)."
		fi
		echo "selbstpruefung: ROT — die Message [$msg] faellt am Traeger (Exit $rc) und HEAD steht unveraendert. Meldung des Traegers:"
		printf '%s\n' "$out" | sed -n '1,3p' | sed 's/^/selbstpruefung:   /'
	else
		if [ "$rc" -ne 0 ]; then
			printf '%s\n' "$out" >&2
			fehler "der Commit '$msg' faellt am Traeger (Exit $rc) — erwartet war ein Durchgang. Fuehrt dieses Repo an [$SELBSTPRUEFUNG_TRAEGER] seinen EIGENEN Traeger, nennt SELBSTPRUEFUNG_MSG_GRUEN eine Message, die jener durchlaesst."
		fi
		betreff="$(git -C "$klon" log -1 --format=%s)"
		if [ "$betreff" != "$msg" ]; then
			fehler "der Commit '$msg' meldet Erfolg, HEAD traegt aber '$betreff' — der Durchgang ist nicht belegt."
		fi
		echo "selbstpruefung: GRUEN — die Message [$msg] geht durch, HEAD traegt: $betreff"
	fi
done

# (4) DAS GATE-KOMMANDO IM KLON. Es laeuft zuletzt, weil die Schritte davor
# Historie schreiben und der Klon danach wegfaellt.
gate_rc=0
gate_out="$(cd "$klon" && bash -c "$SELBSTPRUEFUNG_GATE" 2>&1)" || gate_rc=$?
if [ "$gate_rc" -ne 0 ]; then
	printf '%s\n' "$gate_out" >&2
	fehler "[$SELBSTPRUEFUNG_GATE] im Klon endet mit Exit $gate_rc — der Klon traegt den Stand von HEAD, und der ist damit nicht gruen."
fi
# DIE GANZE AUSGABE, nicht ihre letzte Zeile: was ein Kommando hinterlaesst, ist
# das einzige, woran von aussen abzulesen ist, WELCHES lief. Eine einzelne Zeile
# stammt aus einem fremden Werkzeug und traegt diese Unterscheidung nicht.
echo "selbstpruefung: GATE — [$SELBSTPRUEFUNG_GATE] im Klon ist Exit 0. Ausgabe:"
printf '%s\n' "$gate_out" | sed 's/^/selbstpruefung:   /'

# (5) DER STOP-HOOK IM KLON, ALLE ZWEIGE. Den Nachweis schreibt record-gates.sh
# direkt, nicht das Gate-Kommando: ein gesetzter Gate-Marker muss record-gates nicht
# fahren, und ohne Stempel misst der Lauf nur den strengen Zweig. Gelesen wird die
# Entscheidung des echten Hooks, Exit 0 vorausgesetzt (MR-079).
#
# DIESES REPO IST STRENG, OHNE SCHALTER: der Hook blockt ein Turn-Ende mit
# ungedecktem Inhalt auch ohne neuen HEAD, solange .harness/stop-gate-locker fehlt.
# Der Lauf misst das zuerst, legt die Datei dann an, damit der lockere Zweig messbar
# ist, und prueft, dass STOP_GATE_STRENG=1 sie aufhebt. Dass der Klon die Datei NICHT
# traegt, ist selbst Teil der Zusage — eine versionierte Lockerung bricht den Lauf.
stop_hook="$klon/.claude/hooks/stop-require-gates.sh"
[ -f "$stop_hook" ] || fehler "im Klon liegt kein Stop-Hook unter .claude/hooks/stop-require-gates.sh."
[ ! -e "$klon/.harness/stop-gate-locker" ] || fehler "im Klon liegt .harness/stop-gate-locker — dieses Repo haelt den Stop-Hook streng (MR-079); die Datei lockert das Handoff-Gate und braucht eine ADR (AGENTS.md §3.6)."
stop_entscheid() {
	local aus rc=0
	aus="$(cd "$klon" && printf '{}' | env -u STOP_GATE_STRENG ${1:+"STOP_GATE_STRENG=$1"} bash .claude/hooks/stop-require-gates.sh 2>&1)" || rc=$?
	[ "$rc" -eq 0 ] || fehler "der Stop-Hook endet mit Exit $rc (erwartet 0 mit einer Entscheidung):
$aus"
	case "$aus" in
	*'"decision":"approve"'*) echo frei ;;
	*'"decision": "block"'*) echo blockiert ;;
	*) fehler "der Stop-Hook liefert keine lesbare Entscheidung:
$aus" ;;
	esac
}
rg_out=""
rg_out="$(cd "$klon" && bash tools/harness/record-gates.sh 2>&1)" || fehler "record-gates.sh im Klon endet nicht mit Exit 0:
$rg_out"
[ -s "$klon/.harness/state/gates-passed.head" ] || fehler "record-gates.sh hat im Klon keinen HEAD-Stempel geschrieben (.harness/state/gates-passed.head)."
entscheid="$(stop_entscheid "")"; [ "$entscheid" = frei ] || fehler "der Stop-Hook haelt einen gedeckten Stand auf — erwartet ist die Freigabe nach record-gates."
printf 'selbstpruefung\n' >"$klon/selbstpruefung-stop.txt"
entscheid="$(stop_entscheid "")"; [ "$entscheid" = blockiert ] || fehler "ohne .harness/stop-gate-locker gibt der Stop-Hook ein Turn-Ende mit ungedeckter Aenderung frei — erwartet ist die Blockade (streng ist der Default)."
echo "selbstpruefung: STOP STRENG — ohne Schalter haelt der Stop-Hook ein Turn-Ende mit ungedeckter Aenderung auf, auch ohne neuen HEAD."
: >"$klon/.harness/stop-gate-locker"
entscheid="$(stop_entscheid "")"; [ "$entscheid" = frei ] || fehler "mit .harness/stop-gate-locker haelt der Stop-Hook ein Turn-Ende ohne neuen Commit auf — der lockere Zweig ist nicht der, den die Selbstpruefung misst."
echo "selbstpruefung: STOP LOCKER — nur mit .harness/stop-gate-locker geht ein Turn-Ende ohne neuen HEAD durch; so waere der Hook, wenn dieses Repo die Datei anlegte (das braucht eine ADR)."
entscheid="$(stop_entscheid 1)"; [ "$entscheid" = blockiert ] || fehler "mit STOP_GATE_STRENG=1 gibt der Stop-Hook ein Turn-Ende mit ungedeckter Aenderung frei, obwohl die Lockerungs-Datei liegt — erwartet ist die Blockade."
echo "selbstpruefung: STOP VARIABLE — STOP_GATE_STRENG=1 haelt dasselbe Turn-Ende trotz der Datei auf."
commit_out=""
commit_out="$(git -C "$klon" add -A 2>&1 &&
	git -C "$klon" -c user.email=selbstpruefung@example.invalid -c user.name=selbstpruefung \
		commit -q -m "$SELBSTPRUEFUNG_MSG_GRUEN" 2>&1)" || fehler "der Commit fuer die Stop-Hook-Probe ist nicht entstanden:
$commit_out"
entscheid="$(stop_entscheid "")"; [ "$entscheid" = blockiert ] || fehler "im lockeren Zweig gibt der Stop-Hook einen Commit ohne Nachweis frei — erwartet ist die Blockade (neuer HEAD und abweichender Hash)."
echo "selbstpruefung: STOP BLOCK — im lockeren Zweig haelt der Stop-Hook einen Commit ohne Nachweis auf (neuer HEAD, abweichender Hash); im strengen haelt ihn schon der Hash auf."
rm -f "$klon/.harness/stop-gate-locker"
entscheid="$(stop_entscheid "")"; [ "$entscheid" = blockiert ] || fehler "ohne Schalter gibt der Stop-Hook einen Commit ohne Nachweis frei — erwartet ist die Blockade."

echo "selbstpruefung: OK — der Traeger [$SELBSTPRUEFUNG_TRAEGER] reist mit dem Klon, seine Aktivierung nicht; [$SELBSTPRUEFUNG_AKTIVIERUNG] nimmt genau ihn in Betrieb, danach faellt '$SELBSTPRUEFUNG_MSG_ROT' und geht '$SELBSTPRUEFUNG_MSG_GRUEN' durch, und [$SELBSTPRUEFUNG_GATE] laeuft im Klon gruen; der Stop-Hook haelt ohne Schalter ein Turn-Ende mit ungedeckter Aenderung auf, gibt es nur mit .harness/stop-gate-locker frei, haelt es mit STOP_GATE_STRENG=1 trotz der Datei auf und haelt einen Commit ohne Nachweis auf. Nicht geprueft: ob der Traeger jeden Commit-Pfad erreicht (--no-verify, Werkzeug-Commits, andere Aufrufformen)."
