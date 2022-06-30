#/usr/bin/bash

usage(){
  cat << EOF
usage: $0 [-i] [-f] [-c]

Install extra resources.

OPTIONS:
   -c|--canaries     Install Canary objects
   -f|--flagger      Install Flagger
   -i|--istio        Install Istio
   -h|--help         Show this message
EOF
}

git_root() {
  git rev-parse --show-toplevel
}

tools_bin="$(git_root)/tools/bin"

tool() {
  name=${1}

  if [ ! -f "${tools_bin}/${name}" ]; then
    echo "'${name}' is not available; use 'make dependencies'" >&2

    exit 1
  fi

  echo "${tools_bin}/${name}"
}

install_istio() {
  echo "Install Istio"

  kubectl apply -f "$(git_root)/tools/extra-resources/ns/"
  kubectl apply -f "$(git_root)/tools/extra-resources/istio/"
}

install_flagger() {
  echo "Install Flagger"

  kubectl apply -f "$(git_root)/tools/extra-resources/ns/"
  kubectl apply -f "$(git_root)/tools/extra-resources/flagger/"
}

install_canaries() {
  echo "Install Canaries"

  kubectl apply -f "$(git_root)/tools/extra-resources/ns/"
  kubectl apply -f "$(git_root)/tools/extra-resources/canary/"
}


INSTALL_ISTIO=0
INSTALL_FLAGGER=0
INSTALL_CANARIES=0
while [ ! $# -eq 0 ]; do
  case "${1}" in
    -h | --help)
      usage
      exit
      ;;
    -f | --flagger)
      INSTALL_FLAGGER=1
      ;;
    -i | --istio)
      INSTALL_ISTIO=1
      ;;
    -c | --canaries)
      INSTALL_CANARIES=1
      ;;
    *)
      usage
      exit 1
      ;;
  esac
  shift
done

# They are called here and not inside the 'case' above, becasue if more than one
# defined, the order is important.
if [ $INSTALL_ISTIO -eq 1 ]; then install_istio; fi
if [ $INSTALL_FLAGGER -eq 1 ]; then install_flagger; fi
if [ $INSTALL_CANARIES -eq 1 ]; then install_canaries; fi
