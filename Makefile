smoketest: FORCE
	go build .
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040 --size short ./targets/sgkb/left/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040 --size short ./targets/sgkb/right/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040 --size short ./targets/sgkey/
	tinygo build -o /tmp/out.uf2 --target wioterminal --size short ./targets/wiokey/

FORCE:
