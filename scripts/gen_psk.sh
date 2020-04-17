#!/bin/bash

script_dir="$( cd "$( dirname "${0}" )" >/dev/null 2>&1 && pwd )"
nahs_dir=$(dirname ${script_dir})
conf_dir="${nahs_dir}/config"
psk_file="${conf_dir}/private_network.psk"

mkdir "${nahs_dir}/config" 2> /dev/null
echo $(hexdump -n 16 -e '4/4 "%08X"' /dev/random | awk '{print tolower($0)}') >\
	"${nahs_dir}/config/private_network.psk"
echo "Wrote PSK to file ${psk_file}"
