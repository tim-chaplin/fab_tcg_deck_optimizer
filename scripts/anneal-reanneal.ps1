# anneal-reanneal.ps1 - run simulated-annealing anneal on the same deck over and over, each
# pass starting from the current on-disk best. Stops only when you Ctrl-C.
#
# Every pass invokes `go run ./cmd/fabsim anneal -deck <Deck> -start-temp <T> ...`. anneal
# converges to a local maximum and writes the best-ever deck to mydecks/<Deck>.json; the next
# pass picks up from that deck (re-evaluates its baseline) and starts annealing again from the
# full -StartTemp. Repeated high-T walks are the whole point: each escape attempts to dislodge
# the deck from the basin of attraction the previous pass left it in.
#
# Usage:
#   ./scripts/anneal-reanneal.ps1 -Deck viserai_annealed -StartTemp 1 -Incoming 7

[CmdletBinding()]
param(
    [Parameter(Mandatory)][string]$Deck,
    [Parameter(Mandatory)][double]$StartTemp,
    [int]$Incoming = 0,
    [double]$TempDecay = 0.95,
    [double]$MinTemp = 0,
    [int]$ShallowShuffles = 100,
    [int]$DeepShuffles = 10000,
    [int]$DeckSize = 40,
    [int]$MaxCopies = 2,
    [string]$Format = "silver_age",
    [switch]$Reevaluate,
    [switch]$IterateDebug
)

$ErrorActionPreference = 'Stop'

$deckPath = Join-Path 'mydecks' "$Deck.json"

# Seed bestSeen from disk if the deck already exists, so "new best" lines in the log reflect
# real gains across the whole reannealing session rather than just the first pass.
$bestSeen = 0.0
if (Test-Path $deckPath) {
    $d = Get-Content $deckPath -Raw | ConvertFrom-Json
    $hands = [double]$d.Stats.Hands
    if ($hands -gt 0) {
        $bestSeen = [double]$d.Stats.TotalValue / $hands
    }
}

Write-Host "=== anneal-reanneal: $Deck, startTemp=$StartTemp, incoming=$Incoming ==="
Write-Host ("Starting from {0:F3} - Ctrl-C to stop.`n" -f $bestSeen)

$pass = 0
while ($true) {
    $pass++
    Write-Host "--- Pass $pass ---"

    $goArgs = @(
        'run', './cmd/fabsim', 'anneal',
        '-deck', $Deck,
        '-incoming', $Incoming,
        '-start-temp', $StartTemp,
        '-temp-decay', $TempDecay,
        '-min-temp', $MinTemp,
        '-shallow-shuffles', $ShallowShuffles,
        '-deep-shuffles', $DeepShuffles,
        '-deck-size', $DeckSize,
        '-max-copies', $MaxCopies,
        '-format', $Format
    )
    if ($Reevaluate)   { $goArgs += '-reevaluate' }
    if ($IterateDebug) { $goArgs += '-debug' }
    # The starting card list is identical every pass once we're iterating on the same deck file,
    # so suppress it from pass 2 onward — the user already saw it in pass 1 and reprinting ~40
    # lines each pass buries the progress signal.
    if ($pass -gt 1)   { $goArgs += '-quiet-load' }

    & go @goArgs
    # Exit 130 is anneal's "user pressed Enter" signal. Break the outer loop so the whole
    # session stops instead of kicking off another pass; any other non-zero is also treated
    # as a reason to stop.
    if ($LASTEXITCODE -ne 0) {
        Write-Host "anneal exited $LASTEXITCODE; ending reanneal session."
        break
    }

    if (Test-Path $deckPath) {
        $d = Get-Content $deckPath -Raw | ConvertFrom-Json
        $hands = [double]$d.Stats.Hands
        $avg = if ($hands -gt 0) { [double]$d.Stats.TotalValue / $hands } else { 0.0 }
        if ($avg -gt $bestSeen) {
            Write-Host ("Pass {0}: new best avg {1:F3} (was {2:F3}, delta +{3:F3})" -f $pass, $avg, $bestSeen, ($avg - $bestSeen))
            $bestSeen = $avg
        } else {
            Write-Host ("Pass {0}: avg {1:F3} (best seen {2:F3})" -f $pass, $avg, $bestSeen)
        }
    } else {
        Write-Warning "No deck file at $deckPath after pass $pass; skipping summary."
    }
    Write-Host ""
}
