TINYGO_TARGET=waveshare-rp2040-lora

build:
	tinygo build -target ${TINYGO_TARGET} -scheduler tasks -gc conservative -size full --stack-size=4kb -tags "plugin_insecure_rng plugin_beacon" -o ${TINYGO_TARGET}.uf2

flash:
	tinygo flash -target ${TINYGO_TARGET} -scheduler tasks -gc conservative -size full --stack-size=4kb -tags "plugin_insecure_rng plugin_beacon"

monitor:
	tinygo monitor -target ${TINYGO_TARGET}

all: flash monitor
