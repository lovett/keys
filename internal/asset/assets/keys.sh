#!/usr/bin/env sh

set -eu

REMOTE_URL="{{ .PublicUrl }}"

if ! command -v curl >/dev/null 2>&1; then
    echo "Curl is not available. Cannot run.";
    exit 1
fi

list_keys() {
    endpoint="$REMOTE_URL"
    if [ -n "$1" ]; then
        endpoint="$endpoint?$1"
    fi
    curl -H "Accept: text/plain" "$endpoint"
}

press_key() {
    curl -X POST -H "Accept: text/plain" "$REMOTE_URL/trigger/$1"
}

usage() {
    echo "A command-line client for $REMOTE_URL"
    echo ""
    echo "Usage:" >&2
    echo "  press KEY:  Tell the server to press the key KEY." >&2
    echo "  list: Display the full list of available keys." >&2
    echo ""
    echo "Filters for list command:"
    echo "  --label VALUE: Key label starts with VALUE" >&2
    echo "  --command VALUE: Key command contains VALUE" >&2
    echo "  --key VALUE: Key name contains VALUE" >&2
}

if [ $# -eq 0 ] || [ "$1" = "--help" ]; then
    usage
    exit 1
fi

if [ "$1" = "press" ]; then
    if [ -z "$2" ]; then
        echo "Key not specified." >%2
        exit 1
    fi

    press_key "$2"
fi

if [ "$1" = "list" ]; then
    case "$2" in
        "--label"|"--command"|"--key")
            if [ -z "$3" ]; then
                echo "Missing value for $2 filter." >&2
                exit 1
            fi
            q="${2#--}=$3"
            ;;
        "")
            q=""
            ;;
        *)
            echo "Invalid argument '$2'. Must be 'label' or 'command' or 'key'."
            ;;
    esac
    list_keys "$q"
fi
