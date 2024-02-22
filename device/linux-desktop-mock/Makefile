TINYGO_TARGET=
OUTPUT_FILE_NAME=

# Function to determine the correct output file name based on TINYGO_TARGET
define determine_output_file_name
$(if $(TINYGO_TARGET),$(TINYGO_TARGET).uf2,default_output.uf2)
endef

build:
	@$(eval OUTPUT_FILE_NAME=$(call determine_output_file_name))
	$(if $(TINYGO_TARGET),tinygo build -target $(TINYGO_TARGET) -scheduler tasks -gc conservative -size full -o $(OUTPUT_FILE_NAME),tinygo build -scheduler tasks -gc conservative -size full -o $(OUTPUT_FILE_NAME))

flash:
	$(if $(TINYGO_TARGET),tinygo flash -target $(TINYGO_TARGET) -scheduler tasks -gc conservative -size full,tinygo flash -scheduler tasks -gc conservative -size full)

monitor:
	$(if $(TINYGO_TARGET),tinygo monitor -target $(TINYGO_TARGET),tinygo monitor)
