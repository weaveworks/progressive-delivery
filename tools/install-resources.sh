#/usr/bin/bash

usage(){
  cat << EOF
usage: $0 [-i] [-f] [-c] [-p]

Install extra resources.

OPTIONS:
   -p|--prometheus   Install Prometheus
   -c|--canaries     Install Canary objects
   -f|--flagger      Install Flagger
   -i|--istio        Install Istio
      --all          Install all
      --uninstall    Uninstall instead of install
   -h|--help         Show this message
EOF
}

git_root() {
  git rev-parse --show-toplevel
}

log() {
  message="${1}"

  echo "  >>> ${message}"
}

install_ns() {
  name="${1}"

  kubectl get ns "${name}" >/dev/null 2>&1 \
    || kubectl apply -f "$(git_root)/tools/extra-resources/ns/${name}.yaml"
}

wait_for_resource() {
  kind=${1}
  name=${2}
  ns=${3}
  max_retry=${4:-10}

  counter=1
  while [ ${counter} -le ${max_retry} ]; do
    log "> [${counter}/${max_retry}] waiting for ${name}"

    if [[ "${ns}" == "" ]]; then
      kubectl get "${kind}" "${name}" >/dev/null 2>&1
    else
      kubectl get "${kind}" --namespace="${ns}" "${name}" >/dev/null 2>&1
    fi
    if [ $? -eq 0 ]; then
      log "> ${name} is ready"
      return
    fi


    counter=$(( counter + 1 ))
    sleep 1
  done

  log "!!! ${name} is not available after ${max_retry} retries" >&2

  exit 1
}

install_prometheus() {
  log "Install Prometheus"

  kubectl apply -f "$(git_root)/tools/extra-resources/prometheus/"
  kubectl -n istio-system rollout status deployment/prometheus
}

install_istio() {
  log "Install Istio"

  install_ns "istio-system"
  install_ns "istio-ingress"

  kubectl apply -f "$(git_root)/tools/extra-resources/istio/"

  wait_for_resource "crd" "gateways.networking.istio.io"
  wait_for_resource "service" "istiod" "istio-system"

  kubectl apply -f "$(git_root)/tools/extra-resources/istio/resources/"
}

install_flagger() {
  log "Install Flagger"

  install_ns "flagger"
  kubectl apply -f "$(git_root)/tools/extra-resources/flagger/"

}

install_canaries() {
  log "Install Canaries"

  wait_for_resource "crd" "canaries.flagger.app"

  install_ns "canary"
  kubectl apply -f "$(git_root)/tools/extra-resources/canary/"
}

uninstall_ns() {
  name="${1}"

  kubectl get ns "${name}" >/dev/null 2>&1 \
    && kubectl delete -f "$(git_root)/tools/extra-resources/ns/${name}.yaml"
}

uninstall_prometheus() {
  log "Uninstall Prometheus"

  kubectl delete -f "$(git_root)/tools/extra-resources/prometheus/" 2>/dev/null
}

uninstall_istio() {
  log "Uninstall Istio"

  kubectl delete -f "$(git_root)/tools/extra-resources/istio/resources/" 2>/dev/null
  kubectl delete -f "$(git_root)/tools/extra-resources/istio/" 2>/dev/null

  uninstall_ns "istio-ingress"
  uninstall_ns "istio-system"
}

uninstall_flagger() {
  log "Uninstall Flagger"

  kubectl apply -f "$(git_root)/tools/extra-resources/flagger/" 2>/dev/null
  uninstall_ns "flagger"
}

uninstall_canaries() {
  log "Uninstall Canaries"

  kubectl delete -f "$(git_root)/tools/extra-resources/canary/"  2>/dev/null
  uninstall_ns "canary"
}



ISTIO=0
FLAGGER=0
CANARIES=0
PROMETHEUS=0
UNINSTALL=0
while [ ! $# -eq 0 ]; do
  case "${1}" in
    -h | --help)
      usage
      exit
      ;;
    -f | --flagger)
      FLAGGER=1
      ;;
    -i | --istio)
      ISTIO=1
      ;;
    -c | --canaries)
      CANARIES=1
      ;;
    -p | --prometheus)
      PROMETHEUS=1
      ;;
    --all)
      ISTIO=1
      FLAGGER=1
      CANARIES=1
      PROMETHEUS=1
      ;;
    --uninstall)
      UNINSTALL=1
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
if [ $UNINSTALL -ne 0 ]; then
  if [ $CANARIES -eq 1 ]; then uninstall_canaries; fi
  if [ $FLAGGER -eq 1 ]; then uninstall_flagger; fi
  if [ $PROMETHEUS -eq 1 ]; then uninstall_prometheus; fi
  if [ $ISTIO -eq 1 ]; then uninstall_istio; fi
else
  if [ $ISTIO -eq 1 ]; then install_istio; fi
  if [ $PROMETHEUS -eq 1 ]; then install_prometheus; fi
  if [ $FLAGGER -eq 1 ]; then install_flagger; fi
  if [ $CANARIES -eq 1 ]; then install_canaries; fi
fi
