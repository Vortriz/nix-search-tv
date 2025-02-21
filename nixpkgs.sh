#!/usr/bin/env bash

# === Change keybinds or add more here ===

declare -a INDEXES=(
    "nixpkgs ctrl-n"
    "home-manager ctrl-h"
    "all ctrl-a"
)

OPEN_SOURCE_KEY="ctrl-s"
OPEN_HOMEPAGE_KEY="enter"

OPENER="xdg-open"
if [[ "$(uname)" == 'Darwin' ]]; then
    OPENER="open"
fi

# ========================================

# for debug / development
CMD="${NIX_SEARCH_TV:-nix-search-tv}"

# bind_index binds the given $key to the given $index
bind_index() {
    local key="$1"
    local index="$2"

    local prompt=""
    local indexes_flag=""
    if [[ -n "$index" && "$index" != "all" ]]; then
        indexes_flag="--indexes $index"
        prompt=$index
    fi

    local preview="$CMD preview $indexes_flag"
    local print="$CMD print $indexes_flag"

    echo "$key:change-prompt($prompt> )+change-preview($preview {})+reload($print)"
}

STATE_FILE="/tmp/nix-search-tv-fzf"

# save_state saves the currently displayed index
# to the $STATE_FILE. This file serves as an external script state
# for communication between "print" and "preview" commands
save_state() {
    local index="$1"

    local indexes_flag=""
    if [[ -n "$index" && "$index" != "all" ]]; then
        indexes_flag="--indexes $index"
    fi

    echo "execute(echo $indexes_flag > $STATE_FILE)"
}

HEADER="$OPEN_HOMEPAGE_KEY  - open homepage
$OPEN_SOURCE_KEY - open source code
"

FZF_BINDS=""
for e in "${INDEXES[@]}"; do
    index=$(echo "$e" | awk '{ print $1 }')
    keybind=$(echo "$e" | awk '{ print $2 }')

    fzf_bind=$(bind_index "$keybind" "$index")
    fzf_save_state=$(save_state "$index")
    FZF_BINDS="$FZF_BINDS --bind '$fzf_bind+$fzf_save_state'"

    newline=$'\n'
    HEADER="$HEADER$keybind - $index$newline"
done

# reset the state
echo "" > /tmp/nix-search-tv-fzf

eval "$CMD print | fzf \
    --preview '$CMD preview \$(cat $STATE_FILE) {}' \
    --bind '$OPEN_SOURCE_KEY:execute($CMD source \$(cat $STATE_FILE) {} | xargs $OPENER)' \
    --bind '$OPEN_HOMEPAGE_KEY:execute($CMD homepage \$(cat $STATE_FILE) {} | xargs $OPENER)' \
    --layout reverse \
    --scheme history \
    --preview-window=wrap \
    --header '$HEADER' \
    --header-first \
    --header-border \
    --header-label \"Help\" \
    $FZF_BINDS
"
