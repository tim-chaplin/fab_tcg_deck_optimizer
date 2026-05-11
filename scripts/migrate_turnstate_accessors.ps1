param([switch]$DryRun)

$ErrorActionPreference = 'Stop'

# Walk every .go file outside internal/sim/*.go (non-test) and migrate TurnState
# field accesses to accessor methods. Receivers are detected per-file by scanning
# for declarations / parameters that bind a TurnState variable; only those
# receivers' references are rewritten. Struct-literal field names inside
# TurnStateSpec / TurnState literals are intentionally left alone — the spec
# struct keeps uppercase exported field names.

$privatized = @(
    'CardBanished', 'ActionPoints', 'ArcaneDamageDealt', 'OpponentMarked',
    'Auras', 'Triggers', 'Value', 'CardsPlayed', 'AuraCreated',
    'CardsRemaining', 'Pitched', 'Overpower', 'NonAttackActionPlayed',
    'IncomingDamage', 'ArcaneIncomingDamage', 'BlockTotal', 'Defenders',
    'TriggeringCard'
)

# Skip these dirs entirely — sim-internal files were already handled.
$skipDirs = @('internal/sim')

$paths = @(
    'internal/cards',
    'internal/heroes',
    'internal/weapons',
    'internal/testutils',
    'turntests'
) + @(Get-ChildItem -Path 'internal/sim' -Filter '*_test.go' -File | ForEach-Object { $_.FullName })

$files = New-Object System.Collections.ArrayList
foreach ($p in $paths) {
    if (Test-Path $p -PathType Container) {
        Get-ChildItem -Path $p -Filter '*.go' -File -Recurse | ForEach-Object { [void]$files.Add($_.FullName) }
    } elseif (Test-Path $p -PathType Leaf) {
        [void]$files.Add($p)
    }
}

Write-Host ("Scanning {0} files" -f $files.Count)

function Get-TurnStateReceivers($text) {
    $names = New-Object System.Collections.Generic.HashSet[string]
    # Declaration patterns: `name := sim.TurnState{...}`, `name := &sim.TurnState{...}`,
    # `name := sim.NewTurnStateFromSpec(...)`, `var name sim.TurnState`,
    # `var name *sim.TurnState`.
    foreach ($m in [regex]::Matches($text, '(?m)(\b\w+)\s*:=\s*&?sim\.TurnState\{')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    foreach ($m in [regex]::Matches($text, '(?m)(\b\w+)\s*:=\s*sim\.NewTurnStateFromSpec\(')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    foreach ($m in [regex]::Matches($text, '(?m)\bvar\s+(\w+)\s+\*?sim\.TurnState\b')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    # Parameter / receiver pattern: `name *sim.TurnState` inside () in function sigs.
    foreach ($m in [regex]::Matches($text, '(\b\w+)\s+\*sim\.TurnState\b')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    # Same patterns but inside `internal/sim` package — receivers reference TurnState
    # directly without the sim. prefix.
    foreach ($m in [regex]::Matches($text, '(?m)(\b\w+)\s*:=\s*&?TurnState\{')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    foreach ($m in [regex]::Matches($text, '(?m)\bvar\s+(\w+)\s+\*?TurnState\b')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    foreach ($m in [regex]::Matches($text, '(\b\w+)\s+\*TurnState\b')) {
        [void]$names.Add($m.Groups[1].Value)
    }
    return $names
}

$updated = 0
foreach ($path in $files) {
    $orig = Get-Content -Path $path -Raw -Encoding UTF8
    if (-not $orig) { continue }
    $text = $orig

    $receivers = Get-TurnStateReceivers $text
    if ($receivers.Count -eq 0) { continue }

    foreach ($recv in $receivers) {
        $recvPat = [regex]::Escape($recv)

        # 1. Compound mutations on ActionPoints (++ only — no decrements in cards-side).
        $text = [regex]::Replace($text, $recvPat + '\.ActionPoints\+\+', $recv + '.AddActionPoints(1)')

        # 2. Decrement Value (Test of Strength's clash-loss). Other -- decrements not used.
        $text = [regex]::Replace($text, $recvPat + '\.Value--', $recv + '.SetValue(' + $recv + '.Value() - 1)')

        # 3. Plain assignments: `recv.Field = expr` (not == or !=, not part of < <= > >= := += -=)
        foreach ($f in $privatized) {
            $pat = '(?m)^(\s*)' + $recvPat + '\.' + $f + '\s*=\s*([^=].*?)$'
            $text = [regex]::Replace($text, $pat, ('${1}' + $recv + '.Set' + $f + '(${2})'))
        }

        # 4. Reads: `recv.Field` -> `recv.Field()`. Skip if already followed by `(`
        #    (already a method call) or `:` (struct-literal field name).
        foreach ($f in $privatized) {
            $pat = $recvPat + '\.' + $f + '\b(?![\(:])'
            $text = [regex]::Replace($text, $pat, $recv + '.' + $f + '()')
        }
    }

    if ($text -cne $orig) {
        if ($DryRun) {
            Write-Host ("DRY: {0}" -f $path)
        } else {
            [System.IO.File]::WriteAllText($path, $text, (New-Object System.Text.UTF8Encoding($false)))
            $updated++
        }
    }
}
Write-Host ("Updated {0} files" -f $updated)
