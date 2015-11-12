prepare() {
	mkdir fsimg.loc
	cat > $MOS_FILE << EOF
type: fsimg
location: fsimg.loc
volumeMap: \\([a-z]\\)\\([^_]*\\)_\\(.*\\) \\1/\\2/\\3
EOF
}

cleanup() {
	rmdir fsimg.loc || fail "mosaic dir not empty"
	rm -f $MOS_FILE
}
