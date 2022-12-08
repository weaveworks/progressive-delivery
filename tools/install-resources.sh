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

wait_for_crd() {
  name=${1}
  max_retry=${2:-10}

  counter=1
  while [ ${counter} -le ${max_retry} ]; do
    echo "> [${counter}/${max_retry}] waiting for ${name}"

    kubectl get crd "${name}" >/dev/null 2>&1
    if [ $? -eq 0 ]; then
      echo "> ${name} is ready"
      return
    fi


    counter=$(( counter + 1 ))
    sleep 1
  done

  echo "!!! ${name} is not available after ${max_retry} retries"

  exit 1
}


install_istio() {
  echo "Install Istio"

  kubectl apply -f "$(git_root)/tools/extra-resources/ns/"
  kubectl apply -f "$(git_root)/tools/extra-resources/istio/"

  wait_for_crd "gateways.networking.istio.io"

  kubectl apply -f "$(git_root)/tools/extra-resources/istio/resources/"
}

install_flagger() {
  echo "Install Flagger"

  kubectl apply -f "$(git_root)/tools/extra-resources/ns/"
  kubectl apply -f "$(git_root)/tools/extra-resources/flagger/"

}

install_canaries() {
  echo "Install Canaries"

  wait_for_crd "canaries.flagger.app"

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
