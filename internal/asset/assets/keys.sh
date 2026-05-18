#!/usr/bin/env sh

set -eu

REMOTE_URL="{{ .PublicUrl }}"
VERSION="{{ .Version }}"

if ! command -v curl >/dev/null 2>&1; then
    echo "Curl is not available. Cannot run.";
    exit 1
fi

case "${1:-}" in
    list)
        case "${2:-}" in
            "")
                curl -H "Accept: text/plain" "$REMOTE_URL"
                ;;
            --label | --command | --key)
                if [ -z "${3:-}" ]; then
                    echo "Missing value for $2 filter." >&2
                    exit 1
                fi
                curl -H "Accept: text/plain" "$REMOTE_URL?${2#--}=$3"
            ;;
            *)
                echo "Invalid list argument. Must be 'label' or 'command' or 'key'."
                exit 1
                ;;
        esac
        ;;

    press)
        case "${2:-default}" in
            default)
                echo "Key not specified." >%2
                exit 1
                ;;
            --version)
                echo "$VERSION"
                exit
                ;;
            *)
                curl -X POST -H "Accept: text/plain" "$REMOTE_URL/trigger/$2"
                exit
                ;;
        esac
        ;;
    *)
        echo "A command-line client for $REMOTE_URL"
        echo ""
        echo "Usage:" >&2
        echo "  press KEY:  Tell the server to press the specified key." >&2
        echo "  list: Show the full list of available keys." >&2
        echo "  list --label VALUE: Show keys whose label starts with VALUE." >&2
        echo "  list --command VALUE: Show keys whose command contains VALUE" >&2
        echo "  list --key VALUE: Show keys whose name contains VALUE" >&2
        exit 1
        ;;
esac
