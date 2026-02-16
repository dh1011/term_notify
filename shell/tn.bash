#!/usr/bin/env bash
# Bash/Zsh helper for term_notify (tn)
# Source this file in your .bashrc or .zshrc:
#   source /path/to/tn.bash
#
# Then use:
#   tn run npm run build
#   tn notify "Done!"
#   my-command; tn notify "my-command finished"

# Ensure tn is in PATH or set the full path here
TN_EXE="${TN_EXE:-tn}"

# tnr: wrap a command with term_notify run
tnr() {
    "$TN_EXE" run "$@"
}

# tnd: notify about the last command's result
# Usage: my-command; tnd
# Or:    my-command; tnd "Custom message"
tnd() {
    local last_exit=$?
    local msg="${1:-}"

    if [ -z "$msg" ]; then
        if [ "$last_exit" -eq 0 ]; then
            msg="Previous command succeeded"
        else
            msg="Previous command failed (exit code $last_exit)"
        fi
    fi

    "$TN_EXE" notify "$msg"
}

# Optional: auto-notify for long-running commands
# Uncomment and adjust the threshold (in seconds) to enable.
# Requires bash-preexec or zsh precmd/preexec hooks.
#
# TN_AUTO_THRESHOLD=30  # seconds
#
# _tn_preexec() {
#     _TN_CMD_START=$SECONDS
#     _TN_CMD_NAME="$1"
# }
#
# _tn_precmd() {
#     local elapsed=$(( SECONDS - ${_TN_CMD_START:-$SECONDS} ))
#     if [ "$elapsed" -ge "${TN_AUTO_THRESHOLD:-30}" ] && [ -n "$_TN_CMD_NAME" ]; then
#         "$TN_EXE" notify "Finished: $_TN_CMD_NAME (${elapsed}s)"
#     fi
#     unset _TN_CMD_START _TN_CMD_NAME
# }
#
# # For Bash (requires https://github.com/rcaloras/bash-preexec):
# # preexec_functions+=(_tn_preexec)
# # precmd_functions+=(_tn_precmd)
#
# # For Zsh:
# # preexec() { _tn_preexec "$1"; }
# # precmd()  { _tn_precmd; }

echo "term_notify loaded. Commands: tn, tnr (run+notify), tnd (notify last result)"
