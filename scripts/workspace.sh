#!/usr/bin/env sh

# Create a tmux workspace.

SESSION_NAME="keys"

set -eu

cd "$(dirname "$0")/../"

tmux new-session -d -s "$SESSION_NAME" "$SHELL"

tmux send-keys -t "$SESSION_NAME" "$EDITOR ." C-m

tmux new-window -a -t "$SESSION_NAME" "$SHELL"

tmux new-window -a -t "$SESSION_NAME" -n "server" "scripts/watch.sh"

tmux select-window -t "$SESSION_NAME":0

tmux attach-session -t "$SESSION_NAME"
