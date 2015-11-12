#!/bin/bash

set -e
set -u
set -o pipefail

TDIR=test-rm-me-XXXXXXXX
if [ "$1" = "clean" ]; then
	rm -rf ${TDIR//X/?}
	exit
fi

TESTDIR=$(mktemp -d -p . $TDIR)
echo "Tests:  $1"
echo "Driver: $2"
echo "Dir:    $TESTDIR"

rm_tdir() {
	cd .. && rm -rf $TESTDIR
}
cd $TESTDIR
trap rm_tdir EXIT

. ../"env.sh"
. ../"${2}.sh"
MOS_FILE="./${2}.mos"

for T in ${1//,/ }; do
	. ../"${T}.sh"

	if ! prepare; then
		if $? -eq 22; then
			echo "SKIP"
			exit 0
		else
			echo "FAIL"
			exit 1
		fi
	fi
	run_tests && echo "PASS" || echo "FAIL"
	cleanup
done
