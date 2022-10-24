#!/bin/bash

# We don't need return codes for "$(command)", only stdout is needed.
# Allow `[[ -n "$(command)" ]]`, `func "$(command)"`, pipes, etc.
# shellcheck disable=SC2312

set -u

abort() {
  printf "%s\n" "$@" >&2
  exit 1
}

# Fail fast with a concise message when not using bash
# Single brackets are needed here for POSIX compatibility
# shellcheck disable=SC2292
if [ -z "${BASH_VERSION:-}" ]
then
  abort "Bash is required to interpret this script."
fi

# Check if script is run with force-interactive mode in CI
if [[ -n "${CI-}" && -n "${INTERACTIVE-}" ]]
then
  abort "Cannot run force-interactive mode in CI."
fi

# Check if both `INTERACTIVE` and `NONINTERACTIVE` are set
# Always use single-quoted strings with `exp` expressions
# shellcheck disable=SC2016
if [[ -n "${INTERACTIVE-}" && -n "${NONINTERACTIVE-}" ]]
then
  abort 'Both `$INTERACTIVE` and `$NONINTERACTIVE` are set. Please unset at least one variable and try again.'
fi

# Check if script is run in POSIX mode
if [[ -n "${POSIXLY_CORRECT+1}" ]]
then
  abort 'Bash must not run in POSIX mode. Please unset POSIXLY_CORRECT and try again.'
fi

# string formatters
if [[ -t 1 ]]
then
  tty_escape() { printf "\033[%sm" "$1"; }
else
  tty_escape() { :; }
fi
tty_mkbold() { tty_escape "1;$1"; }
tty_underline="$(tty_escape "4;39")"
tty_blue="$(tty_mkbold 34)"
tty_red="$(tty_mkbold 31)"
tty_bold="$(tty_mkbold 39)"
tty_reset="$(tty_escape 0)"

shell_join() {
  local arg
  printf "%s" "$1"
  shift
  for arg in "$@"
  do
    printf " "
    printf "%s" "${arg// /\ }"
  done
}

chomp() {
  printf "%s" "${1/"$'\n'"/}"
}

ohai() {
  printf "${tty_blue}==>${tty_bold} %s${tty_reset}\n" "$(shell_join "$@")"
}

warn() {
  printf "${tty_red}Warning${tty_reset}: %s\n" "$(chomp "$1")"
}

# Check if script is run non-interactively (e.g. CI)
# If it is run non-interactively we should not prompt for passwords.
# Always use single-quoted strings with `exp` expressions
# shellcheck disable=SC2016
if [[ -z "${NONINTERACTIVE-}" ]]
then
  if [[ -n "${CI-}" ]]
  then
    warn 'Running in non-interactive mode because `$CI` is set.'
    NONINTERACTIVE=1
  elif [[ ! -t 0 ]]
  then
    if [[ -z "${INTERACTIVE-}" ]]
    then
      warn 'Running in non-interactive mode because `stdin` is not a TTY.'
      NONINTERACTIVE=1
    else
      warn 'Running in interactive mode despite `stdin` not being a TTY because `$INTERACTIVE` is set.'
    fi
  fi
else
  ohai 'Running in non-interactive mode because `$NONINTERACTIVE` is set.'
fi

# USER isn't always set so provide a fall back for the installer and subprocesses.
if [[ -z "${USER-}" ]]
then
  USER="$(chomp "$(id -un)")"
  export USER
fi

# First check OS.
OS="$(uname)"
if [[ "${OS}" == "Linux" ]]
then
  LIFERAY_CLI_ON_LINUX=1
elif [[ "${OS}" == "Darwin" ]]
then
  LIFERAY_CLI_ON_MACOS=1
else
  abort "This installer is only supported on macOS and Linux."
fi

# Required installation paths. To install elsewhere, you can download
# a binary for your platform and place it anywhere you like.
if [[ -n "${LIFERAY_CLI_ON_MACOS-}" ]]
then
  UNAME_MACHINE="$(/usr/bin/uname -m)"

  LIFERAY_CLI_PREFIX="/usr/local"
  STAT_PRINTF=("stat" "-f")
  PERMISSION_FORMAT="%A"
  CHOWN=("/usr/sbin/chown")
  CHGRP=("/usr/bin/chgrp")
  GROUP="admin"
  INSTALL=("/usr/bin/install" -o "root" -g "wheel" -m "0755")
else
  UNAME_MACHINE="$(uname -m)"

  LIFERAY_CLI_PREFIX="/usr/local"
  STAT_PRINTF=("stat" "--printf")
  PERMISSION_FORMAT="%a"
  CHOWN=("/bin/chown")
  CHGRP=("/bin/chgrp")
  GROUP="$(id -gn)"
  INSTALL=("/usr/bin/install" -o "root" -g "root" -m "0755")
fi

CHMOD=("/bin/chmod")
MKDIR=("/bin/mkdir" "-p")

unset HAVE_SUDO_ACCESS # unset this from the environment

have_sudo_access() {
  if [[ ! -x "/usr/bin/sudo" ]]
  then
    return 1
  fi

  local -a SUDO=("/usr/bin/sudo")
  if [[ -n "${SUDO_ASKPASS-}" ]]
  then
    SUDO+=("-A")
  elif [[ -n "${NONINTERACTIVE-}" ]]
  then
    SUDO+=("-n")
  fi

  if [[ -z "${HAVE_SUDO_ACCESS-}" ]]
  then
    if [[ -n "${NONINTERACTIVE-}" ]]
    then
      "${SUDO[@]}" -l mkdir &>/dev/null
    else
      "${SUDO[@]}" -v && "${SUDO[@]}" -l mkdir &>/dev/null
    fi
    HAVE_SUDO_ACCESS="$?"
  fi

  if [[ -n "${LIFERAY_CLI_ON_MACOS-}" ]] && [[ "${HAVE_SUDO_ACCESS}" -ne 0 ]]
  then
    abort "Need sudo access on macOS (e.g. the user ${USER} needs to be an Administrator)!"
  fi

  return "${HAVE_SUDO_ACCESS}"
}

execute() {
  if ! "$@"
  then
    abort "$(printf "Failed during: %s" "$(shell_join "$@")")"
  fi
}

execute_sudo() {
  local -a args=("$@")
  if have_sudo_access
  then
    if [[ -n "${SUDO_ASKPASS-}" ]]
    then
      args=("-A" "${args[@]}")
    fi
    ohai "/usr/bin/sudo" "${args[@]}"
    execute "/usr/bin/sudo" "${args[@]}"
  else
    ohai "${args[@]}"
    execute "${args[@]}"
  fi
}

getc() {
  local save_state
  save_state="$(/bin/stty -g)"
  /bin/stty raw -echo
  IFS='' read -r -n 1 -d '' "$@"
  /bin/stty "${save_state}"
}

ring_bell() {
  # Use the shell's audible bell.
  if [[ -t 1 ]]
  then
    printf "\a"
  fi
}

wait_for_user() {
  local c
  echo
  echo "Press ${tty_bold}RETURN${tty_reset}/${tty_bold}ENTER${tty_reset} to continue or any other key to abort:"
  getc c
  # we test for \r and \n because some stuff does \r instead
  if ! [[ "${c}" == $'\r' || "${c}" == $'\n' ]]
  then
    exit 1
  fi
}

major_minor() {
  echo "${1%%.*}.$(
    x="${1#*.}"
    echo "${x%%.*}"
  )"
}

version_gt() {
  [[ "${1%.*}" -gt "${2%.*}" ]] || [[ "${1%.*}" -eq "${2%.*}" && "${1#*.}" -gt "${2#*.}" ]]
}
version_ge() {
  [[ "${1%.*}" -gt "${2%.*}" ]] || [[ "${1%.*}" -eq "${2%.*}" && "${1#*.}" -ge "${2#*.}" ]]
}
version_lt() {
  [[ "${1%.*}" -lt "${2%.*}" ]] || [[ "${1%.*}" -eq "${2%.*}" && "${1#*.}" -lt "${2#*.}" ]]
}

check_run_command_as_root() {
  [[ "${EUID:-${UID}}" == "0" ]] || return

  # Allow Azure Pipelines/GitHub Actions/Docker/Concourse/Kubernetes to do everything as root (as it's normal there)
  [[ -f /.dockerenv ]] && return
  [[ -f /proc/1/cgroup ]] && grep -E "azpl_job|actions_job|docker|garden|kubepods" -q /proc/1/cgroup && return

  abort "Don't run this as root!"
}

get_permission() {
  "${STAT_PRINTF[@]}" "${PERMISSION_FORMAT}" "$1"
}

user_only_chmod() {
  [[ -d "$1" ]] && [[ "$(get_permission "$1")" != 75[0145] ]]
}

exists_but_not_writable() {
  [[ -e "$1" ]] && ! [[ -r "$1" && -w "$1" && -x "$1" ]]
}

get_owner() {
  "${STAT_PRINTF[@]}" "%u" "$1"
}

file_not_owned() {
  [[ "$(get_owner "$1")" != "$(id -u)" ]]
}

get_group() {
  "${STAT_PRINTF[@]}" "%g" "$1"
}

file_not_grpowned() {
  [[ " $(id -G "${USER}") " != *" $(get_group "$1") "* ]]
}

test_curl() {
  if [[ ! -x "$1" ]]
  then
    return 1
  fi

  local curl_version_output curl_name_and_version
  curl_version_output="$("$1" --version 2>/dev/null)"
  curl_name_and_version="${curl_version_output%% (*}"
  version_ge "$(major_minor "${curl_name_and_version##* }")" "$(major_minor "${REQUIRED_CURL_VERSION}")"
}

# Search for the given executable in PATH (avoids a dependency on the `which` command)
which() {
  # Alias to Bash built-in command `type -P`
  type -P "$@"
}

# Search PATH for the specified program that satisfies Homebrew requirements
# function which is set above
# shellcheck disable=SC2230
find_tool() {
  if [[ $# -ne 1 ]]
  then
    return 1
  fi

  local executable
  while read -r executable
  do
    if "test_$1" "${executable}"
    then
      echo "${executable}"
      break
    fi
  done < <(which -a "$1")
}

# Invalidate sudo timestamp before exiting (if it wasn't active before).
if [[ -x /usr/bin/sudo ]] && ! /usr/bin/sudo -n -v 2>/dev/null
then
  trap '/usr/bin/sudo -k' EXIT
fi

####################################################################### script
# shellcheck disable=SC2016
ohai 'Checking for `sudo` access (which may request your password)...'

if [[ -n "${LIFERAY_CLI_ON_MACOS-}" ]]
then
  have_sudo_access
elif ! [[ -w "${LIFERAY_CLI_PREFIX}" ]] &&
     ! have_sudo_access
then
  abort "$(
    cat <<EOABORT
Insufficient permissions to install liferay cli to \"${LIFERAY_CLI_PREFIX}\" (the default prefix).
EOABORT
  )"
fi

check_run_command_as_root

if [[ -d "${LIFERAY_CLI_PREFIX}" && ! -x "${LIFERAY_CLI_PREFIX}" ]]
then
  abort "$(
    cat <<EOABORT
The liferay cli prefix ${tty_underline}${LIFERAY_CLI_PREFIX}${tty_reset} exists but is not searchable.
If this is not intentional, please restore the default permissions and
try running the installer again:
    sudo chmod 775 ${LIFERAY_CLI_PREFIX}
EOABORT
  )"
fi

if [[ -n "${LIFERAY_CLI_ON_MACOS-}" ]]
then
  # On macOS, support 64-bit Intel and ARM
  if [[ "${UNAME_MACHINE}" != "arm64" ]] && [[ "${UNAME_MACHINE}" != "x86_64" ]]
  then
    abort "liferay cli is only supported on Intel and ARM processors!"
  fi
else
  # On Linux, support 64-bit Intel and ARM
  if  [[ "${UNAME_MACHINE}" != "aarch64" ]] && [[ "${UNAME_MACHINE}" != "x86_64" ]]
  then
    abort "liferay cli is only supported on Intel and ARM processors!"
  fi
fi

ohai "This script will install:"
echo "${LIFERAY_CLI_PREFIX}/bin/liferay"

group_chmods=()
directories=(bin)
mkdirs=()
for dir in "${directories[@]}"
do
  if ! [[ -d "${LIFERAY_CLI_PREFIX}/${dir}" ]]
  then
    mkdirs+=("${LIFERAY_CLI_PREFIX}/${dir}")
  fi
done

user_chmods=()
mkdirs_user_only=()
chmods=()
if [[ "${#group_chmods[@]}" -gt 0 ]]
then
  chmods+=("${group_chmods[@]}")
fi
if [[ "${#user_chmods[@]}" -gt 0 ]]
then
  chmods+=("${user_chmods[@]}")
fi

chowns=()
chgrps=()
if [[ "${#chmods[@]}" -gt 0 ]]
then
  for dir in "${chmods[@]}"
  do
    if file_not_owned "${dir}"
    then
      chowns+=("${dir}")
    fi
    if file_not_grpowned "${dir}"
    then
      chgrps+=("${dir}")
    fi
  done
fi

if [[ "${#group_chmods[@]}" -gt 0 ]]
then
  ohai "The following existing directories will be made group writable:"
  printf "%s\n" "${group_chmods[@]}"
fi
if [[ "${#user_chmods[@]}" -gt 0 ]]
then
  ohai "The following existing directories will be made writable by user only:"
  printf "%s\n" "${user_chmods[@]}"
fi
if [[ "${#chowns[@]}" -gt 0 ]]
then
  ohai "The following existing directories will have their owner set to ${tty_underline}${USER}${tty_reset}:"
  printf "%s\n" "${chowns[@]}"
fi
if [[ "${#chgrps[@]}" -gt 0 ]]
then
  ohai "The following existing directories will have their group set to ${tty_underline}${GROUP}${tty_reset}:"
  printf "%s\n" "${chgrps[@]}"
fi
if [[ "${#mkdirs[@]}" -gt 0 ]]
then
  ohai "The following new directories will be created:"
  printf "%s\n" "${mkdirs[@]}"
fi

additional_shellenv_commands=()

if [[ -z "${NONINTERACTIVE-}" ]]
then
  ring_bell
  wait_for_user
fi

if [[ -d "${LIFERAY_CLI_PREFIX}" ]]
then
  if [[ "${#chmods[@]}" -gt 0 ]]
  then
    execute_sudo "${CHMOD[@]}" "u+rwx" "${chmods[@]}"
  fi
  if [[ "${#group_chmods[@]}" -gt 0 ]]
  then
    execute_sudo "${CHMOD[@]}" "g+rwx" "${group_chmods[@]}"
  fi
  if [[ "${#user_chmods[@]}" -gt 0 ]]
  then
    execute_sudo "${CHMOD[@]}" "go-w" "${user_chmods[@]}"
  fi
  if [[ "${#chowns[@]}" -gt 0 ]]
  then
    execute_sudo "${CHOWN[@]}" "${USER}" "${chowns[@]}"
  fi
  if [[ "${#chgrps[@]}" -gt 0 ]]
  then
    execute_sudo "${CHGRP[@]}" "${GROUP}" "${chgrps[@]}"
  fi
else
  execute_sudo "${INSTALL[@]}" "${LIFERAY_CLI_PREFIX}"
fi

if [[ "${#mkdirs[@]}" -gt 0 ]]
then
  execute_sudo "${MKDIR[@]}" "${mkdirs[@]}"
  execute_sudo "${CHMOD[@]}" "ug=rwx" "${mkdirs[@]}"
  if [[ "${#mkdirs_user_only[@]}" -gt 0 ]]
  then
    execute_sudo "${CHMOD[@]}" "go-w" "${mkdirs_user_only[@]}"
  fi
  execute_sudo "${CHOWN[@]}" "${USER}" "${mkdirs[@]}"
  execute_sudo "${CHGRP[@]}" "${GROUP}" "${mkdirs[@]}"
fi

ohai "Downloading and installing liferay cli..."
(
  LIFERAY_CLI_RELEASE=$(curl -fsSL https://api.github.com/repos/liferay/liferay-cli/releases | jq -r 'sort_by(.published_at) | .[-1]')
  LIFERAY_CLI_CHECKSUMS_FILE=$(jq -r '.assets[] | select(.name=="checksums.txt") | .browser_download_url' <<< $LIFERAY_CLI_RELEASE)
  LIFERAY_CLI_CHECKSUMS=$(curl -fsSL $LIFERAY_CLI_CHECKSUMS_FILE)

  if [[ -n "${LIFERAY_CLI_ON_MACOS-}" ]]
  then
    if [[ "${UNAME_MACHINE}" == "arm64" ]]
    then
      LIFERAY_CLI_BINARY="liferay-darwin-arm64"
    elif [[ "${UNAME_MACHINE}" == "x86_64" ]]
    then
      LIFERAY_CLI_BINARY="liferay-darwin-amd64"
    fi
  else
    if  [[ "${UNAME_MACHINE}" == "aarch64" ]]
    then
      LIFERAY_CLI_BINARY="liferay-linux-arm64"
    elif [[ "${UNAME_MACHINE}" == "x86_64" ]]
    then
      LIFERAY_CLI_BINARY="liferay-linux-amd64"
    fi
  fi

  LIFERAY_CLI_BINARY_FILE=$(jq -r ".assets[] | select(.name==\"${LIFERAY_CLI_BINARY}\") | .browser_download_url" <<< $LIFERAY_CLI_RELEASE)
  curl -fsSL "$LIFERAY_CLI_BINARY_FILE" -o "/tmp/${LIFERAY_CLI_BINARY}"
  LIFERAY_CLI_BINARY_SUM=$(cd /tmp && shasum -a 256 "${LIFERAY_CLI_BINARY}")
  LIFERAY_CLI_CHECKSUM=$(echo "$LIFERAY_CLI_CHECKSUMS" | grep "$LIFERAY_CLI_BINARY")

  if [[ "${LIFERAY_CLI_BINARY_SUM}" == "${LIFERAY_CLI_CHECKSUM}" ]]
  then
    execute_sudo "${INSTALL[@]}" "/tmp/${LIFERAY_CLI_BINARY}" "${LIFERAY_CLI_PREFIX}/bin/liferay"

    rm "/tmp/${LIFERAY_CLI_BINARY}"
  else

    abort "$(
      cat <<EOABORT
liferay cli failed checksum!
  Calculated Sum: ${LIFERAY_CLI_BINARY_SUM}
  Downloaded Sum: ${LIFERAY_CLI_CHECKSUM}
EOABORT
  )"

  fi
) || exit 1

ohai "Installation successful!"
echo

ring_bell
