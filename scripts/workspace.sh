#!/usr/bin/env sh

set -eu

cd "$(dirname "$0")/../"

workspace() {
    SESSION_NAME="keys"
	tmux new-session -d -s "$SESSION_NAME" "$SHELL"
	tmux send-keys -t "$SESSION_NAME" "$EDITOR ." C-m

	tmux new-window -a -t "$SESSION_NAME" "$SHELL"

	tmux new-window -a -t "$SESSION_NAME" -n "server" "scripts/watch.sh"

	tmux select-window -t "$SESSION_NAME":0
	tmux attach-session -t "$SESSION_NAME"
}

case "${1:-default}" in
    --help)
        echo "Create a tmux workspace."
        ;;
    default)
        workspace
        ;;
    *)
        echo "Unknown argument." >&2
        ;;
esac
