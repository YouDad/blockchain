#!/bin/bash

files=`ls $@`

for file in ${files}; do
	echo -e "\n{{{{{ [TEST{${file}}] }}}}}\n" |\
		ack --flush --passthru --color --color-match "underline bold white" ".*"

	bash "${file}" debug
	if [[ "$?" != "0" ]]; then
		exit "$?"
	fi
done
echo
