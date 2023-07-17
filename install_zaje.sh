#!/bin/sh -e

log_use_fancy_output () { 
    TPUT=/usr/bin/tput
    EXPR=/usr/bin/expr
    if  [ -t 1 ] && 
        [ "x${TERM:-}" != "x" ] && 
        [ "x${TERM:-}" != "xdumb" ] && 
        [ -x $TPUT ] && [ -x $EXPR ] && 
        $TPUT hpa 60 >/dev/null 2>&1 &&
        $TPUT setaf 1 >/dev/null 2>&1 
    then
        [ -z "$FANCYTTY" ] && FANCYTTY=1 || true 
    else 
        FANCYTTY=0
    fi   
    case "$FANCYTTY" in
        1|Y|yes|true)   true;;
        *)              false;;
    esac 
}

# Only do the fancy stuff if we have an appropriate terminal
# and if /usr is already mounted
RED=''
YELLOW=''
BLUE=''
NORMAL=''
BOLD=''
UNSET=''
if log_use_fancy_output; then
    RED=$( $TPUT setaf 1)
    YELLOW=$( $TPUT setaf 3)
    BLUE=$( $TPUT setaf 6)
    NORMAL=$( $TPUT setaf 2)
    BOLD=$($TPUT bold)
    UNSET=$( $TPUT op)
fi

if [ ! "$(which curl 2>/dev/null)" ];then
	printf '%b' "${RED}Need to install curl.${NORMAL}\n"
	exit 2
fi   

NAME="zaje"
HIGHLIGHT_REPO_NAME="highlight"
GH_SPACE="jessp01"
LATEST_VER=$(curl -s "https://api.github.com/repos/$GH_SPACE/$NAME/releases/latest"| grep tag_name|sed 's@\s*"tag_name": "\(.*\)".*@\1@')
OS=$(uname)
ARCH=$(uname -m)
BIN_ARCHIVE="zaje_${OS}_${ARCH}.tar.gz"

# we need this for the lexers
LATEST_HIGHLIGHT_VER=$(curl -s "https://api.github.com/repos/$GH_SPACE/$HIGHLIGHT_REPO_NAME/releases/latest"| grep tag_name|sed 's@\s*"tag_name": "\(.*\)".*@\1@')
HIGHLIGHT_SOURCE_ARCHIVE="${LATEST_HIGHLIGHT_VER}.tar.gz"

CONFIG_DIR="$HOME/.config/$NAME"
LEXERS_DIR="$CONFIG_DIR/syntax_files"
TMP_DIR="/tmp/$NAME"
FUNCTIONS_RC_FILE="$CONFIG_DIR/${NAME}_functions.rc"


printf '%b' "${BOLD}${NORMAL}\nWelcome to ${BLUE}$NAME ($LATEST_VER)${NORMAL}'s installation script:)\n"

mkdir -p "$CONFIG_DIR" "$TMP_DIR"
cd $TMP_DIR

printf '%b' "${NORMAL}Fetching sources...\n\n"
curl -Ls "https://github.com/$GH_SPACE/$NAME/releases/download/${LATEST_VER}/${BIN_ARCHIVE}" --output "${BIN_ARCHIVE}"
curl -Ls "https://github.com/$GH_SPACE/$HIGHLIGHT_REPO_NAME/archive/refs/tags/${HIGHLIGHT_SOURCE_ARCHIVE}" --output "${HIGHLIGHT_SOURCE_ARCHIVE}"

tar zxf "$BIN_ARCHIVE"
mkdir -p ~/bin
mv "$NAME" ~/bin
mv README.md LICENSE "$CONFIG_DIR"


TIMESTAMP=$(date +%s)

if [ -f "$CONFIG_DIR/${NAME}_functions.rc" ];then
    printf '%b' "${BOLD}${YELLOW}$FUNCTIONS_RC_FILE already exists...\n${NORMAL}I'll place the new copy under ${BLUE}${FUNCTIONS_RC_FILE}.${TIMESTAMP}${NORMAL}\n\n"
    FUNCTIONS_RC_FILE="${FUNCTIONS_RC_FILE}.${TIMESTAMP}"
fi

curl -Ls "https://github.com/$GH_SPACE/$NAME/raw/$LATEST_VER/utils/functions.rc" -o "$FUNCTIONS_RC_FILE"

if [ -d "$LEXERS_DIR" ];then
    printf '%b' "${YELLOW}$LEXERS_DIR already exists...\n${NORMAL}I'll place the new lexers under ${BLUE}${LEXERS_DIR}.${TIMESTAMP}${NORMAL}\n\n"
    LEXERS_DIR="${LEXERS_DIR}.${TIMESTAMP}"
fi

tar zxf "$HIGHLIGHT_SOURCE_ARCHIVE"
VERSION_NO_V=$(echo "$LATEST_HIGHLIGHT_VER" | sed 's/^v\(.*\)/\1/')
mv "$HIGHLIGHT_REPO_NAME-$VERSION_NO_V/syntax_files" "$LEXERS_DIR"

printf '%b' "All sorted:)\n\n${BLUE}* $NAME${NORMAL} binary is in ~/bin/${NAME}\n"
printf "* Useful helper functions are under ${BLUE}$FUNCTIONS_RC_FILE\n${NORMAL}  Source them with ${BLUE}'. $FUNCTIONS_RC_FILE'${NORMAL}.\n"
printf "* Lexers are under ${BLUE}$LEXERS_DIR${NORMAL}\n\n"
printf "Downloaded archives are available in ${BLUE}$TMP_DIR${NORMAL}.. Feel free to discard them.${UNSET}\n"

if log_use_fancy_output ;then
    $TPUT sgr0
fi
