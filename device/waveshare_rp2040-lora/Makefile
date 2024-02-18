TINYGO_TARGET=waveshare-rp2040-lora

build:
	tinygo build -target ${TINYGO_TARGET} -scheduler tasks -gc conservative -size full -o waveshare_rp2040-lora.uf2

flash:
	tinygo flash -target ${TINYGO_TARGET} -scheduler tasks -gc conservative -size full

monitor:
	tinygo monitor -target ${TINYGO_TARGET}