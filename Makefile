smoketest: FORCE
	go build .
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/fric10key/
	tinygo build -o /tmp/out.uf2 --target gopher-badge    --size short ./targets/gopher-badge/
	tinygo build -o /tmp/out.uf2 --target macropad-rp2040 --size short ./targets/macropad-rp2040/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkb/left/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkb/right/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkey/
	tinygo build -o /tmp/out.uf2 --target wioterminal     --size short ./targets/wiokey/

FORCE:
