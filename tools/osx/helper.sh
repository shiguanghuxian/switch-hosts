#!/bin/bash

###############################################################################
# Debug
###############################################################################

# _debug()
#
# Usage:
#   _debug printf "Debug info. Variable: %s\n" "$0"
#
# A simple function for executing a specified command if the `$_USE_DEBUG`
# variable has been set. The command is expected to print a message and
# should typically be either `echo`, `printf`, or `cat`.
__DEBUG_COUNTER=0
_debug() {
  if [[ "${_USE_DEBUG:-"0"}" -eq 1 ]]
  then
    __DEBUG_COUNTER=$((__DEBUG_COUNTER+1))
    # Prefix debug message with "bug (U+1F41B)"
    printf "🐛  %s " "${__DEBUG_COUNTER}"
    "${@}"
    printf "――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――\n"
  fi
}
# debug()
#
# Usage:
#   debug "Debug info. Variable: $0"
#
# Print the specified message if the `$_USE_DEBUG` variable has been set.
#
# This is a shortcut for the _debug() function that simply echos the message.
debug() {
  _debug echo "${@}"
}

###############################################################################
# Die
###############################################################################

# _die()
#
# Usage:
#   _die printf "Error message. Variable: %s\n" "$0"
#
# A simple function for exiting with an error after executing the specified
# command. The command is expected to print a message and should typically
# be either `echo`, `printf`, or `cat`.
_die() {
  # Prefix die message with "cross mark (U+274C)", often displayed as a red x.
  printf "❌  "
  "${@}" 1>&2
  exit 1
}
# die()
#
# Usage:
#   die "Error message. Variable: $0"
#
# Exit with an error and print the specified message.
#
# This is a shortcut for the _die() function that simply echos the message.
die() {
  _die echo "${@}"
}


###############################################################################
# Help
###############################################################################

# _print_help()
#
# Usage:
#   _print_help
#
# Print the program help information.
_print_help() {
  cat <<HEREDOC
 __  __  ___  _    __  __  ___
|  \/  |/ _ \| |  |  \/  |/ __|
| |\/| | (_) | |__| |\/| | (__
|_|  |_|\___/|____|_|  |_|\___|

Usage:
  ${_ME} [--options] [<arguments>]
  ${_ME} -n AppName -i AppIcon -b BackgroundImage -v DmgVersion
  ${_ME} -h
Options:
  -h --help  Display this help information.
HEREDOC
}

confirm() {
  (($force)) && return 0;

  read -p "$1 [y/N] " -n 1;
  [[ $REPLY =~ ^[Yy]$ ]];
}



###############################################################################
# _spinner()
#
# Usage:
#   _spinner <pid>
#
# Description:
#   Display an ascii spinner while <pid> is running.
#
# Example Usage:
#   ```
#   _spinner_example() {
#     printf "Working..."
#     (sleep 1) &
#     _spinner $!
#     printf "Done!\n"
#   }
#   (_spinner_example)
#   ```
#
# More Information:
#   http://fitnr.com/showing-a-bash-spinner.html
_spinner() {
  local _pid="${1:-}"
  local _delay=0.75
  local _spin_string="|/-\\"

  if [[ -z "${_pid}" ]]
  then
    printf "Usage: _spinner <pid>\n"
    return 1
  fi

  while ps a | awk '{print $1}' | grep -q "${_pid}"
  do
    local _temp="${_spin_string#?}"
    printf " [%c]  " "${_spin_string}"
    _spin_string="${_temp}${_spin_string%${_temp}}"
    sleep ${_delay}
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
}
