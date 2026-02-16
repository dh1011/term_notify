# PowerShell helper for term_notify (tn)
# Source this file in your $PROFILE to get the `tn` alias and helpers.
#
# Usage:
#   Add to your PowerShell profile (~\Documents\PowerShell\Microsoft.PowerShell_profile.ps1):
#     . "path\to\tn.ps1"
#
#   Then use:
#     tn run npm run build
#     tn notify "Done!"
#     my-command; tn notify "my-command finished"

# Ensure tn.exe is in PATH or set the full path here
$TN_EXE = "tn"

function Invoke-TermNotifyRun {
    <#
    .SYNOPSIS
    Wraps a command and sends a notification when it finishes.
    .EXAMPLE
    tnr npm run build
    #>
    & $TN_EXE run @args
}

function Invoke-TermNotifyDone {
    <#
    .SYNOPSIS
    Sends a notification that the previous command finished.
    Reads $LASTEXITCODE to determine success/failure.
    .EXAMPLE
    npm run build; tnd
    #>
    param(
        [string]$Message = ""
    )

    if ($Message -eq "") {
        if ($LASTEXITCODE -eq 0) {
            $Message = "Previous command succeeded"
        } else {
            $Message = "Previous command failed (exit code $LASTEXITCODE)"
        }
    }

    & $TN_EXE notify $Message
}

# Short aliases
Set-Alias -Name tnr -Value Invoke-TermNotifyRun
Set-Alias -Name tnd -Value Invoke-TermNotifyDone

Write-Host "term_notify loaded. Commands: tn, tnr (run+notify), tnd (notify last result)" -ForegroundColor DarkCyan
