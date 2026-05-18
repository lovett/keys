#!/usr/bin/env sh

set -eu

REMOTE_URL="{{ .PublicUrl }}"
VERSION="{{ .Version }}"

if ! command -v curl >/dev/null 2>&1; then
    echo "Curl is not available. Cannot run.";
    exit 1
fi

case "${1:-}" in
    --version)
        echo "$VERSION"
        exit
        ;;

    --help)
        echo "A command-line client for $REMOTE_URL" >&2
        echo "" >&2
        echo "Version: $VERSION" >&2
        echo "" >&2
        echo "Usage:" >&2
        echo "  KEY:  Tell the server to trigger the specified key (by name or physical_key)." >&2
        echo "  list: Show the full list of available keys." >&2
        echo "  list --name VALUE: Show keys whose name starts with VALUE." >&2
        echo "  list --command VALUE: Show keys whose command contains VALUE" >&2
        echo "  list --key VALUE: Show keys whose physical_key is VALUE" >&2
        echo "  --version: application version (both server and client)" >&2
        echo "  --help: this message" >&2
        exit 1
        ;;

    list)
        case "${2:-}" in
            "")
                curl -H "Accept: text/plain" "$REMOTE_URL"
                ;;
            --name | --command | --key)
                if [ -z "${3:-}" ]; then
                    echo "Missing value for $2 filter." >&2
                    exit 1
                fi
                curl -H "Accept: text/plain" "$REMOTE_URL?${2#--}=$3"
            ;;
            *)
                echo "Invalid list filter. Must be 'name' or 'command' or 'key'."
                exit 1
                ;;
        esac
        ;;
    "")
        echo "Key to press not specified." >&2
        exit 1
        ;;
    *)
        curl -X POST -H "Accept: text/plain" "$REMOTE_URL/trigger/$1"
        exit
        ;;
esac
