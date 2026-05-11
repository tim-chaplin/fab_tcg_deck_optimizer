param([switch]$DryRun)

$ErrorActionPreference = 'Stop'

# Internal/sim has many references to the now-privatized TurnState fields. This
# script renames the field references in sim-internal Go files only — cards /
# heroes / weapons / testutils / turntests use the new accessor methods and
# stay untouched.

$renames = @(
    @{ Old = 'CardBanished';          New = 'cardBanished' },
    @{ Old = 'ActionPoints';          New = 'actionPoints' },
    @{ Old = 'ArcaneDamageDealt';     New = 'arcaneDamageDealt' },
    @{ Old = 'OpponentMarked';        New = 'opponentMarked' },
    @{ Old = 'Auras';                 New = 'auras' },
    @{ Old = 'Triggers';              New = 'triggers' },
    @{ Old = 'Value';                 New = 'value' },
    @{ Old = 'CardsPlayed';           New = 'cardsPlayed' },
    @{ Old = 'AuraCreated';           New = 'auraCreated' },
    @{ Old = 'CardsRemaining';        New = 'cardsRemaining' },
    @{ Old = 'Pitched';               New = 'pitched' },
    @{ Old = 'Overpower';             New = 'overpower' },
    @{ Old = 'NonAttackActionPlayed'; New = 'nonAttackActionPlayed' },
    @{ Old = 'IncomingDamage';        New = 'incomingDamage' },
    @{ Old = 'ArcaneIncomingDamage';  New = 'arcaneIncomingDamage' },
    @{ Old = 'BlockTotal';            New = 'blockTotal' },
    @{ Old = 'Defenders';             New = 'defenders' },
    @{ Old = 'TriggeringCard';        New = 'triggeringCard' }
)

# Variable spellings that bind a *TurnState: bufs.state, ctx.state, the parameter
# named state, ts, s, prior, p, t. Match each so the patterns catch every site.
$receivers = @('s', 'ts', 'state', 'prior', 'p', 't', 'bufs\.state', 'ctx\.state')

$files = Get-ChildItem -Path 'internal/sim' -Filter *.go -File
Write-Host ("Scanning {0} sim files" -f $files.Count)

$updated = 0
foreach ($f in $files) {
    $orig = Get-Content -Path $f.FullName -Raw -Encoding UTF8
    $text = $orig

    foreach ($r in $renames) {
        # Field access via known receiver names: <recv>.OldName -> <recv>.NewName.
        foreach ($recv in $receivers) {
            $text = [regex]::Replace($text,
                '(' + $recv + ')\.' + $r.Old + '\b',
                '$1.' + $r.New)
        }
        # Struct literal field name: `OldName:` -> `newName:`. Only matches when
        # OldName starts a struct-literal entry (preceded by `{` or `,` plus
        # whitespace/newline) so other identifier uses aren't touched.
        $text = [regex]::Replace($text,
            '(?m)([{,\s])' + $r.Old + ':',
            '$1' + $r.New + ':')
    }

    if ($text -cne $orig) {
        if ($DryRun) {
            Write-Host ("DRY: {0}" -f $f.Name)
        } else {
            [System.IO.File]::WriteAllText($f.FullName, $text, (New-Object System.Text.UTF8Encoding($false)))
            $updated++
        }
    }
}
Write-Host ("Updated {0} files" -f $updated)
