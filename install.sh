#!/bin/bash
set -euo pipefail

APP_NAME="stackroost"
BIN_PATH="/usr/local/bin/$APP_NAME"
BINARY_URL="https://stackroost.krushnam.cloud/stackroost"

# -----------------------------
# Utility functions
# -----------------------------
log()    { echo -e "[INFO] $*"; }
error()  { echo -e "[ERROR] $*" >&2; }
fatal()  { error "$*"; exit 1; }

prompt_yes_no() {
    while true; do
        read -rp "$1 [y/n]: " yn
        case $yn in
            [Yy]*) return 0 ;;
            [Nn]*) return 1 ;;
        esac
    done
}

# -----------------------------
# Root check
# -----------------------------
if [[ "$EUID" -ne 0 ]]; then
    fatal "This script must be run as root or with sudo."
fi

# -----------------------------
# Dependencies check
# -----------------------------
for cmd in curl stat date; do
    command -v "$cmd" >/dev/null 2>&1 || fatal "Required command '$cmd' not found. Install it first."
done

# -----------------------------
# Network check
# -----------------------------
if ! curl -fsSL --head "$BINARY_URL" >/dev/null; then
    fatal "Cannot reach $BINARY_URL. Check network/DNS/firewall."
fi

# -----------------------------
# Backup existing binary
# -----------------------------
if [[ -f "$BIN_PATH" ]]; then
    log "Binary already exists at $BIN_PATH."
    if prompt_yes_no "Do you want to replace it?"; then
        mv "$BIN_PATH" "$BIN_PATH.bak.$(date +%s)"
        log "Existing binary backed up."
    else
        log "Skipping binary installation."
        exit 0
    fi
fi

# -----------------------------
# Download with progress, speed, ETA
# -----------------------------
log "Downloading $APP_NAME binary to $BIN_PATH..."

# Get total size in bytes
TOTAL_SIZE=$(curl -sI "$BINARY_URL" | grep -i Content-Length | awk '{print $2}' | tr -d '\r')
if [[ -z "$TOTAL_SIZE" ]]; then
    log "Cannot detect file size, downloading without progress..."
    curl -L -o "$BIN_PATH" "$BINARY_URL"
else
    START_TIME=$(date +%s)
    curl -L -o "$BIN_PATH" "$BINARY_URL" --progress-bar 2>&1 | while true; do
        if [[ -f "$BIN_PATH" ]]; then
            BYTES_DOWNLOADED=$(stat -c%s "$BIN_PATH" 2>/dev/null || echo 0)
            PERCENT=$(( BYTES_DOWNLOADED * 100 / TOTAL_SIZE ))
            NOW=$(date +%s)
            ELAPSED=$(( NOW - START_TIME ))
            ELAPSED=$((ELAPSED>0 ? ELAPSED : 1))
            SPEED=$((BYTES_DOWNLOADED / ELAPSED))
            SPEED_KB=$((SPEED / 1024))
            REMAINING=$(( (TOTAL_SIZE - BYTES_DOWNLOADED) / SPEED ))
            printf "\rDownloaded: %3d%% | Speed: %5d KB/s | ETA: %3ds" "$PERCENT" "$SPEED_KB" "$REMAINING"
            if (( BYTES_DOWNLOADED >= TOTAL_SIZE )); then
                break
            fi
        fi
        sleep 1
    done
    echo -e "\nDownload completed!"
fi

chmod +x "$BIN_PATH"

# -----------------------------
# Verify binary works
# -----------------------------
if ! "$BIN_PATH" --help >/dev/null 2>&1; then
    fatal "Binary exists but failed to execute. Check architecture and binary file."
fi

log "$APP_NAME installed successfully!"
log "You can now run the CLI using: $APP_NAME"
log "Binary path: $BIN_PATH"
