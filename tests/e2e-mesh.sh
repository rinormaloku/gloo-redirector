#!/bin/bash

set -e

GIT_ROOT=$(git rev-parse --show-toplevel)

main(){
  apply_mesh_configuration
  # wait for propagation
  sleep 5

  URL=$(kubectl --context=cluster1 get svc -n istio-gateways -l istio=ingressgateway -o jsonpath="{..loadBalancer.ingress[0].ip}")
  if [[ $URL == "" ]] ; then
    echo "Test failed, no URL found"
    exit 1
  fi
  URL="http://${URL}"

  filename="$GIT_ROOT/samples/csv/status_301.csv"

  while read csv_line; do
      IFS=',' read -a csv_values <<< "$csv_line"

      initial_url=${csv_values[0]}
      initial_path=$(path_parse "$initial_url" path)
      domain=$(domain_parse "$initial_url" host)
      final_url=${csv_values[1]}
      redirect_code=${csv_values[2]}

      echo "[TEST] ${initial_url} -> ${final_url} with code ${redirect_code}. domain '$domain', initial_path '$initial_path', redirect_code '$redirect_code'"
      status_code=$(curl --write-out %{http_code} --silent --output /dev/null "${URL}${initial_path}" -H "Host: $domain")
      redirect_location=$(curl --write-out %{redirect_url} --silent --output /dev/null "${URL}${initial_path}" -H "Host: $domain")

      if [[ "${status_code}" == "${redirect_code}" && "${redirect_location}" == "${final_url}" ]] ; then
        echo "[SUCCESS] URL was updated from '$initial_url' to '$redirect_location' with status code '$status_code'"
      else
        echo "[FAILED]"
        echo "status code should be ${redirect_code} it was $status_code"
        echo "location should be '$final_url' it was $redirect_location"
        echo "To try it out execute: "
        echo "  curl -v \"${URL}${initial_path}\" -H \"Host: $domain\""
        exit 1
      fi

      echo ""
  done < "$filename"
}


path_parse() {
  local -r URL=$1
  URL_NOPRO=${URL:7}
  URL_REL=${URL_NOPRO#*/}
  echo "/${URL_REL%%\?*}"
}

domain_parse() {
  local -r tmp_url=$1
  echo $(echo "$tmp_url" | awk -F/ '{print $3}')
}

apply_mesh_configuration() {
 go run $GIT_ROOT/main.go mesh generate --source /tmp/redirections.csv \
  | kubectl --context=cluster1 apply -f -
}

main
