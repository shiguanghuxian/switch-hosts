# Exit immediately on error
set -e

# Detect whether output is piped or not.
[[ -t 1 ]] && piped=0 || piped=1

# Defaults
force=0
quiet=0
verbose=0
interactive=0

# Print help if no arguments were passed.
[[ $# -eq 0 ]] && set -- "-h"

# Set to the program's basename.
_ME=$(basename "${0}")
_CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
_USE_DEBUG=1

source ./helper.sh;
source ./makeapp.sh;
source ./makedmg.sh

_main() {

  _debug printf ">> Performing operation...\n"

  ## Read the options and set stuff
  #while [[ $1 = -?* ]]; do
  #  case $1 in
  #    -h|--help) usage >&2; safe_exit ;;
  #    --version) out "$(basename $0) $version"; safe_exit ;;
  #    -u|--username) shift; username=$1 ;;
  #    -p|--password) shift; password=$1 ;;
  #    -v|--verbose) verbose=1 ;;
  #    -q|--quiet) quiet=1 ;;
  #    -i|--interactive) interactive=1 ;;
  #    -f|--force) force=1 ;;
  #    --endopts) shift; break ;;
  #    *) die "invalid option: $1" ;;
  #  esac
  #  shift
  #done


  # parse cmd args
  while getopts "n:i:b:v:h" arg
  do
    case $arg in
      n)
        name=$OPTARG
        ;;
      i)
        icon=$OPTARG
        ;;
      b)
        img=$OPTARG
        ;;
      v)
        ver=$OPTARG
        ;;
      h)
        _print_help
        die "print help info"
        ;;
      ?)
        _print_help
        die "print help info"
        ;;
    esac
  done

  debug name:$name icon:$icon image:$img version:$ver

  sleep 5

  if [ ! -n "$name" ]; then
    die "invalid args"
  fi

  if [ ! -n "$icon" ]; then
    debug "no icon, just use default icon.png"
    icon="icon.png"
  fi

  targetDir="${name}.app/Contents/MacOS"
  sourceZip="${name}.zip"
  target="${name} ${ver}.dmg"

  if [ -d "$targetDir" ]; then
    debug "cleanning ${name}.app"
    sudo rm -rf ${name}.app
  fi

  if [ -f "$target" ]; then
    debug "clearnning  $target"
    sudo rm -rf $target
  fi

  makeapp $name $icon

  cd $_CUR_DIR

  [[ ! -d "$targetDir" ]] && die "failed to make .app, no $targetDir exist !!!"
  [[ ! -f "$sourceZip" ]] && die "make sure there $sourceZip exist !!!"

  unzip $sourceZip -d $targetDir

  [[ ! -n "$img" ]] && debug "no background image file, use defalut image.png" && img=image.png
  [[ ! -n "$ver" ]] && debug "no version, use defalut 1.0.0" && ver="1.0.0"

  makepack $name $ver $img
}

# Call `_main` after everything has been defined.
_main "${@:-}"
