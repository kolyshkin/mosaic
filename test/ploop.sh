prepare() {
	# Check if ploop is available
	if test -f /proc/vz/ploop_minor && ploop list >/dev/null; then
		true # ploop seems to be working
	else
		echo "* ploop not available, skipping tests"
		return 22
	fi

	mkdir ploop.dir
	cat > $MOS_FILE << EOF
type: ploop
location: ploop.dir
volumeMap: \\([a-z]\\)\\([a-z][a-z]\\)\\(.*\\) \\1/\\2/\\3
EOF
}

cleanup() {
	rmdir ploop.dir || fail "mosaic dir not empty"
	rm -f $MOS_FILE
}
