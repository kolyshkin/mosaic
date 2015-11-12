prepare() {
	mkdir plain.dir
	cat > $MOS_FILE << EOF
type: plain
location: plain.dir
volumeMap: \\([a-z]\\)\\([^_]*\\)_\\(.*\\) \\1/\\2/\\3
EOF
}

cleanup() {
	rmdir plain.dir || fail "mosaic dir not empty"
	rm -f $MOS_FILE
}
