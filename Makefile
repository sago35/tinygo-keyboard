smoketest: FORCE
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040 --size short ./targets/sgkb/left/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040 --size short ./targets/sgkb/right/

FORCE:
