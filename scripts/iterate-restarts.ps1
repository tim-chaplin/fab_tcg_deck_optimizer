# iterate-restarts.ps1 — run `fabsim iterate` N times with different deck names and rank the
# results at the end. Replaces the former -N / -deck-template flags that lived inside fabsim's
# Go code: a thin shell-level wrapper keeps iterate's surface narrow while still giving the
# user multi-restart sweeps.
#
# Usage:
#   ./scripts/iterate-restarts.ps1 -N 10 -DeckTemplate 'viserai_*' -Incoming 7
#
# Each iteration resolves '*' to 1..N and invokes `go run ./cmd/fabsim iterate -deck <name> ...`.
# Decks on disk at mydecks/<name>.json are preserved as resume state (already-converged runs
# pass through one no-op round); pass the same -DeckTemplate to resume a partial sweep. Delete
# the JSON files manually when you want fresh random starts.

[CmdletBinding()]
param(
    [Parameter(Mandatory)][int]$N,
    [Parameter(Mandatory)][string]$DeckTemplate,
    [int]$Incoming = 0,
    [int]$ShallowShuffles = 100,
    [int]$DeepShuffles = 10000,
    [int]$DeckSize = 40,
    [int]$MaxCopies = 2,
    [string]$Format = "silver_age",
    [switch]$Finalize,
    [switch]$Reevaluate,
    [switch]$IterateDebug
)

$ErrorActionPreference = 'Stop'

# Validate template shape mirrors what fabsim's own validator used to reject — exactly one '*'.
$starCount = ($DeckTemplate.ToCharArray() | Where-Object { $_ -eq '*' } | Measure-Object).Count
if ($starCount -ne 1) {
    throw "-DeckTemplate must contain exactly one '*' placeholder (got '$DeckTemplate')"
}
if ($N -lt 1) {
    throw "-N must be >= 1 (got $N)"
}

# Collect each restart's measured avg so the final ranking can be printed.
$results = New-Object System.Collections.Generic.List[object]

for ($i = 1; $i -le $N; $i++) {
    $deckName = $DeckTemplate.Replace('*', $i.ToString())
    Write-Host "`n=== Restart $i/$N : $deckName ==="

    # Build the flag list dynamically so optional switches only appear when set. PowerShell 5.1
    # flattens @-splatting across array boundaries for native exe calls.
    $goArgs = @(
        'run', './cmd/fabsim', 'iterate',
        '-deck', $deckName,
        '-incoming', $Incoming,
        '-shallow-shuffles', $ShallowShuffles,
        '-deep-shuffles', $DeepShuffles,
        '-deck-size', $DeckSize,
        '-max-copies', $MaxCopies,
        '-format', $Format
    )
    if ($Finalize)     { $goArgs += '-finalize' }
    if ($Reevaluate)   { $goArgs += '-reevaluate' }
    if ($IterateDebug) { $goArgs += '-debug' }

    & go @goArgs

    # iterate persists its final state to mydecks/<deckName>.json; the deck's Stats block carries
    # the avg directly, so we don't have to parse stdout to rank. Missing file means the run
    # failed or was aborted before a single improvement saved — leave it out of the ranking.
    $deckPath = Join-Path 'mydecks' "$deckName.json"
    if (Test-Path $deckPath) {
        $deck = Get-Content $deckPath -Raw | ConvertFrom-Json
        $hands = [double]$deck.Stats.Hands
        $avg = if ($hands -gt 0) { [double]$deck.Stats.TotalValue / $hands } else { 0.0 }
        $results.Add([PSCustomObject]@{ Deck = $deckName; Avg = $avg }) | Out-Null
    } else {
        Write-Warning "No deck file at $deckPath after restart ${i}; excluded from ranking."
    }
}

Write-Host "`n=== Restart ranking ==="
$sorted = $results | Sort-Object -Property Avg -Descending
foreach ($r in $sorted) {
    Write-Host ("  {0,-40} avg {1:F3}" -f $r.Deck, $r.Avg)
}
if ($sorted.Count -gt 0) {
    Write-Host ("`nBest: {0} (avg {1:F3})" -f $sorted[0].Deck, $sorted[0].Avg)
}
