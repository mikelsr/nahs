#!/bin/bash

script_dir="$( cd "$( dirname "${0}" )" >/dev/null 2>&1 && pwd )"
nahs_dir=$(dirname ${script_dir})
conf_dir="${nahs_dir}/config"
psk_file="${conf_dir}/private_network.psk"

mkdir "${nahs_dir}/config" 2> /dev/null

header="/key/swarm/psk/1.0.0/\n/base16/\n"

printf ${header} > ${psk_file}
printf "$(hexdump -n 32 -e '8/4 "%08X"' /dev/random | awk '{print tolower($0)}')" >> ${psk_file}
echo "Wrote PSK to file ${psk_file}"
