package sx126x

const (
	// SX126X physical layer properties
	SX126X_FREQUENCY_STEP_SIZE = 0.9536743164
	SX126X_MAX_PACKET_LENGTH   = 255
	SX126X_CRYSTAL_FREQ        = 32.0
	SX126X_DIV_EXPONENT        = 25

	// SX126X SPI commands
	// operational modes commands
	SX126X_CMD_NOP                      = 0x00
	SX126X_CMD_SET_SLEEP                = 0x84
	SX126X_CMD_SET_STANDBY              = 0x80
	SX126X_CMD_SET_FS                   = 0xC1
	SX126X_CMD_SET_TX                   = 0x83
	SX126X_CMD_SET_RX                   = 0x82
	SX126X_CMD_STOP_TIMER_ON_PREAMBLE   = 0x9F
	SX126X_CMD_SET_RX_DUTY_CYCLE        = 0x94
	SX126X_CMD_SET_CAD                  = 0xC5
	SX126X_CMD_SET_TX_CONTINUOUS_WAVE   = 0xD1
	SX126X_CMD_SET_TX_INFINITE_PREAMBLE = 0xD2
	SX126X_CMD_SET_REGULATOR_MODE       = 0x96
	SX126X_CMD_CALIBRATE                = 0x89
	SX126X_CMD_CALIBRATE_IMAGE          = 0x98
	SX126X_CMD_SET_PA_CONFIG            = 0x95
	SX126X_CMD_SET_RX_TX_FALLBACK_MODE  = 0x93

	// register and buffer access commands
	SX126X_CMD_WRITE_REGISTER = 0x0D
	SX126X_CMD_READ_REGISTER  = 0x1D
	SX126X_CMD_WRITE_BUFFER   = 0x0E
	SX126X_CMD_READ_BUFFER    = 0x1E

	// DIO and IRQ control
	SX126X_CMD_SET_DIO_IRQ_PARAMS         = 0x08
	SX126X_CMD_GET_IRQ_STATUS             = 0x12
	SX126X_CMD_CLEAR_IRQ_STATUS           = 0x02
	SX126X_CMD_SET_DIO2_AS_RF_SWITCH_CTRL = 0x9D
	SX126X_CMD_SET_DIO3_AS_TCXO_CTRL      = 0x97

	// RF, modulation and packet commands
	SX126X_CMD_SET_RF_FREQUENCY          = 0x86
	SX126X_CMD_SET_PACKET_TYPE           = 0x8A
	SX126X_CMD_GET_PACKET_TYPE           = 0x11
	SX126X_CMD_SET_TX_PARAMS             = 0x8E
	SX126X_CMD_SET_MODULATION_PARAMS     = 0x8B
	SX126X_CMD_SET_PACKET_PARAMS         = 0x8C
	SX126X_CMD_SET_CAD_PARAMS            = 0x88
	SX126X_CMD_SET_BUFFER_BASE_ADDRESS   = 0x8F
	SX126X_CMD_SET_LORA_SYMB_NUM_TIMEOUT = 0x0A

	// status commands
	SX126X_CMD_GET_STATUS           = 0xC0
	SX126X_CMD_GET_RSSI_INST        = 0x15
	SX126X_CMD_GET_RX_BUFFER_STATUS = 0x13
	SX126X_CMD_GET_PACKET_STATUS    = 0x14
	SX126X_CMD_GET_DEVICE_ERRORS    = 0x17
	SX126X_CMD_CLEAR_DEVICE_ERRORS  = 0x07
	SX126X_CMD_GET_STATS            = 0x10
	SX126X_CMD_RESET_STATS          = 0x00

	// SX126X register map
	SX126X_REG_WHITENING_INITIAL_MSB = 0x06B8
	SX126X_REG_WHITENING_INITIAL_LSB = 0x06B9
	SX126X_REG_CRC_INITIAL_MSB       = 0x06BC
	SX126X_REG_CRC_INITIAL_LSB       = 0x06BD
	SX126X_REG_CRC_POLYNOMIAL_MSB    = 0x06BE
	SX126X_REG_CRC_POLYNOMIAL_LSB    = 0x06BF
	SX126X_REG_SYNC_WORD_0           = 0x06C0
	SX126X_REG_SYNC_WORD_1           = 0x06C1
	SX126X_REG_SYNC_WORD_2           = 0x06C2
	SX126X_REG_SYNC_WORD_3           = 0x06C3
	SX126X_REG_SYNC_WORD_4           = 0x06C4
	SX126X_REG_SYNC_WORD_5           = 0x06C5
	SX126X_REG_SYNC_WORD_6           = 0x06C6
	SX126X_REG_SYNC_WORD_7           = 0x06C7
	SX126X_REG_NODE_ADDRESS          = 0x06CD
	SX126X_REG_BROADCAST_ADDRESS     = 0x06CE
	SX126X_REG_LORA_SYNC_WORD_MSB    = 0x0740
	SX126X_REG_LORA_SYNC_WORD_LSB    = 0x0741
	SX126X_REG_RANDOM_NUMBER_0       = 0x0819
	SX126X_REG_RANDOM_NUMBER_1       = 0x081A
	SX126X_REG_RANDOM_NUMBER_2       = 0x081B
	SX126X_REG_RANDOM_NUMBER_3       = 0x081C
	SX126X_REG_RX_GAIN               = 0x08AC
	SX126X_REG_OCP_CONFIGURATION     = 0x08E7
	SX126X_REG_XTA_TRIM              = 0x0911
	SX126X_REG_XTB_TRIM              = 0x0912

	// undocumented registers
	SX126X_REG_SENSITIVITY_CONFIG  = 0x0889 // SX1268 datasheet v1.1, section 15.1
	SX126X_REG_TX_CLAMP_CONFIG     = 0x08D8 // SX1268 datasheet v1.1, section 15.2
	SX126X_REG_RTC_STOP            = 0x0920 // SX1268 datasheet v1.1, section 15.3
	SX126X_REG_RTC_EVENT           = 0x0944 // SX1268 datasheet v1.1, section 15.3
	SX126X_REG_IQ_CONFIG           = 0x0736 // SX1268 datasheet v1.1, section 15.4
	SX126X_REG_RX_GAIN_RETENTION_0 = 0x029F // SX1268 datasheet v1.1, section 9.6
	SX126X_REG_RX_GAIN_RETENTION_1 = 0x02A0 // SX1268 datasheet v1.1, section 9.6
	SX126X_REG_RX_GAIN_RETENTION_2 = 0x02A1 // SX1268 datasheet v1.1, section 9.6

	// SX126X SPI command variables
	//SX126X_CMD_SET_SLEEP                                                MSB   LSB   DESCRIPTION
	SX126X_SLEEP_START_COLD = 0b00000000 //  2     2     sleep mode: cold start, configuration is lost (default)
	SX126X_SLEEP_START_WARM = 0b00000100 //  2     2                 warm start, configuration is retained
	SX126X_SLEEP_RTC_OFF    = 0b00000000 //  0     0     wake on RTC timeout: disabled
	SX126X_SLEEP_RTC_ON     = 0b00000001 //  0     0                          enabled

	//SX126X_CMD_SET_STANDBY
	SX126X_STANDBY_RC   = 0x00 //  7     0     standby mode: 13 MHz RC oscillator
	SX126X_STANDBY_XOSC = 0x01 //  7     0                   32 MHz crystal oscillator

	//SX126X_CMD_SET_RX
	SX126X_RX_TIMEOUT_NONE = 0x000000 //  23    0     Rx timeout duration: no timeout (Rx single mode)
	SX126X_RX_TIMEOUT_INF  = 0xFFFFFF //  23    0                          infinite (Rx continuous mode)

	//SX126X_CMD_SET_TX
	SX126X_TX_TIMEOUT_NONE = 0x000000 //  23    0     Tx timeout duration: no timeout (Tx single mode)

	//SX126X_CMD_STOP_TIMER_ON_PREAMBLE
	SX126X_STOP_ON_PREAMBLE_OFF = 0x00 //  7     0     stop timer on: sync word or header (default)
	SX126X_STOP_ON_PREAMBLE_ON  = 0x01 //  7     0                    preamble detection

	//SX126X_CMD_SET_REGULATOR_MODE
	SX126X_REGULATOR_LDO   = 0x00 //  7     0     set regulator mode: LDO (default)
	SX126X_REGULATOR_DC_DC = 0x01 //  7     0                         DC-DC

	//SX126X_CMD_CALIBRATE
	SX126X_CALIBRATE_IMAGE_OFF      = 0b00000000 //  6     6     image calibration: disabled
	SX126X_CALIBRATE_IMAGE_ON       = 0b01000000 //  6     6                        enabled
	SX126X_CALIBRATE_ADC_BULK_P_OFF = 0b00000000 //  5     5     ADC bulk P calibration: disabled
	SX126X_CALIBRATE_ADC_BULK_P_ON  = 0b00100000 //  5     5                             enabled
	SX126X_CALIBRATE_ADC_BULK_N_OFF = 0b00000000 //  4     4     ADC bulk N calibration: disabled
	SX126X_CALIBRATE_ADC_BULK_N_ON  = 0b00010000 //  4     4                             enabled
	SX126X_CALIBRATE_ADC_PULSE_OFF  = 0b00000000 //  3     3     ADC pulse calibration: disabled
	SX126X_CALIBRATE_ADC_PULSE_ON   = 0b00001000 //  3     3                            enabled
	SX126X_CALIBRATE_PLL_OFF        = 0b00000000 //  2     2     PLL calibration: disabled
	SX126X_CALIBRATE_PLL_ON         = 0b00000100 //  2     2                      enabled
	SX126X_CALIBRATE_RC13M_OFF      = 0b00000000 //  1     1     13 MHz RC osc. calibration: disabled
	SX126X_CALIBRATE_RC13M_ON       = 0b00000010 //  1     1                                 enabled
	SX126X_CALIBRATE_RC64K_OFF      = 0b00000000 //  0     0     64 kHz RC osc. calibration: disabled
	SX126X_CALIBRATE_RC64K_ON       = 0b00000001 //  0     0                                 enabled
	SX126X_CALIBRATE_ALL            = 0b01111111 //  6     0     calibrate all blocks

	//SX126X_CMD_CALIBRATE_IMAGE
	SX126X_CAL_IMG_430_MHZ_1 = 0x6B
	SX126X_CAL_IMG_430_MHZ_2 = 0x6F
	SX126X_CAL_IMG_470_MHZ_1 = 0x75
	SX126X_CAL_IMG_470_MHZ_2 = 0x81
	SX126X_CAL_IMG_779_MHZ_1 = 0xC1
	SX126X_CAL_IMG_779_MHZ_2 = 0xC5
	SX126X_CAL_IMG_863_MHZ_1 = 0xD7
	SX126X_CAL_IMG_863_MHZ_2 = 0xDB
	SX126X_CAL_IMG_902_MHZ_1 = 0xE1
	SX126X_CAL_IMG_902_MHZ_2 = 0xE9

	//SX126X_CMD_SET_PA_CONFIG
	SX126X_PA_CONFIG_HP_MAX   = 0x07
	SX126X_PA_CONFIG_PA_LUT   = 0x01
	SX126X_PA_CONFIG_SX1262_8 = 0x00

	//SX126X_CMD_SET_RX_TX_FALLBACK_MODE
	SX126X_RX_TX_FALLBACK_MODE_FS         = 0x40 //  7     0     after Rx/Tx go to: FS mode
	SX126X_RX_TX_FALLBACK_MODE_STDBY_XOSC = 0x30 //  7     0                        standby with crystal oscillator
	SX126X_RX_TX_FALLBACK_MODE_STDBY_RC   = 0x20 //  7     0                        standby with RC oscillator (default)

	//SX126X_CMD_SET_DIO_IRQ_PARAMS
	SX126X_IRQ_TIMEOUT           = 0b1000000000 //  9     9     Rx or Tx timeout
	SX126X_IRQ_CAD_DETECTED      = 0b0100000000 //  8     8     channel activity detected
	SX126X_IRQ_CAD_DONE          = 0b0010000000 //  7     7     channel activity detection finished
	SX126X_IRQ_CRC_ERR           = 0b0001000000 //  6     6     wrong CRC received
	SX126X_IRQ_HEADER_ERR        = 0b0000100000 //  5     5     LoRa header CRC error
	SX126X_IRQ_HEADER_VALID      = 0b0000010000 //  4     4     valid LoRa header received
	SX126X_IRQ_SYNC_WORD_VALID   = 0b0000001000 //  3     3     valid sync word detected
	SX126X_IRQ_PREAMBLE_DETECTED = 0b0000000100 //  2     2     preamble detected
	SX126X_IRQ_RX_DONE           = 0b0000000010 //  1     1     packet received
	SX126X_IRQ_TX_DONE           = 0b0000000001 //  0     0     packet transmission completed
	SX126X_IRQ_ALL               = 0b1111111111 //  9     0     all interrupts
	SX126X_IRQ_NONE              = 0b0000000000 //  9     0     no interrupts

	//SX126X_CMD_SET_DIO2_AS_RF_SWITCH_CTRL
	SX126X_DIO2_AS_IRQ       = 0x00 //  7     0     DIO2 configuration: IRQ
	SX126X_DIO2_AS_RF_SWITCH = 0x01 //  7     0                         RF switch control

	//SX126X_CMD_SET_DIO3_AS_TCXO_CTRL
	SX126X_DIO3_OUTPUT_1_6 = 0x00 //  7     0     DIO3 voltage output for TCXO: 1.6 V
	SX126X_DIO3_OUTPUT_1_7 = 0x01 //  7     0                                   1.7 V
	SX126X_DIO3_OUTPUT_1_8 = 0x02 //  7     0                                   1.8 V
	SX126X_DIO3_OUTPUT_2_2 = 0x03 //  7     0                                   2.2 V
	SX126X_DIO3_OUTPUT_2_4 = 0x04 //  7     0                                   2.4 V
	SX126X_DIO3_OUTPUT_2_7 = 0x05 //  7     0                                   2.7 V
	SX126X_DIO3_OUTPUT_3_0 = 0x06 //  7     0                                   3.0 V
	SX126X_DIO3_OUTPUT_3_3 = 0x07 //  7     0                                   3.3 V

	//SX126X_CMD_SET_PACKET_TYPE
	SX126X_PACKET_TYPE_GFSK = 0x00 //  7     0     packet type: GFSK
	SX126X_PACKET_TYPE_LORA = 0x01 //  7     0                  LoRa

	//SX126X_CMD_SET_TX_PARAMS
	SX126X_PA_RAMP_10U   = 0x00 //  7     0     ramp time: 10 us
	SX126X_PA_RAMP_20U   = 0x01 //  7     0                20 us
	SX126X_PA_RAMP_40U   = 0x02 //  7     0                40 us
	SX126X_PA_RAMP_80U   = 0x03 //  7     0                80 us
	SX126X_PA_RAMP_200U  = 0x04 //  7     0                200 us
	SX126X_PA_RAMP_800U  = 0x05 //  7     0                800 us
	SX126X_PA_RAMP_1700U = 0x06 //  7     0                1700 us
	SX126X_PA_RAMP_3400U = 0x07 //  7     0                3400 us

	//SX126X_CMD_SET_MODULATION_PARAMS
	SX126X_GFSK_FILTER_NONE      = 0x00 //  7     0     GFSK filter: none
	SX126X_GFSK_FILTER_GAUSS_0_3 = 0x08 //  7     0                  Gaussian, BT = 0.3
	SX126X_GFSK_FILTER_GAUSS_0_5 = 0x09 //  7     0                  Gaussian, BT = 0.5
	SX126X_GFSK_FILTER_GAUSS_0_7 = 0x0A //  7     0                  Gaussian, BT = 0.7
	SX126X_GFSK_FILTER_GAUSS_1   = 0x0B //  7     0                  Gaussian, BT = 1
	SX126X_GFSK_RX_BW_4_8        = 0x1F //  7     0     GFSK Rx bandwidth: 4.8 kHz
	SX126X_GFSK_RX_BW_5_8        = 0x17 //  7     0                        5.8 kHz
	SX126X_GFSK_RX_BW_7_3        = 0x0F //  7     0                        7.3 kHz
	SX126X_GFSK_RX_BW_9_7        = 0x1E //  7     0                        9.7 kHz
	SX126X_GFSK_RX_BW_11_7       = 0x16 //  7     0                        11.7 kHz
	SX126X_GFSK_RX_BW_14_6       = 0x0E //  7     0                        14.6 kHz
	SX126X_GFSK_RX_BW_19_5       = 0x1D //  7     0                        19.5 kHz
	SX126X_GFSK_RX_BW_23_4       = 0x15 //  7     0                        23.4 kHz
	SX126X_GFSK_RX_BW_29_3       = 0x0D //  7     0                        29.3 kHz
	SX126X_GFSK_RX_BW_39_0       = 0x1C //  7     0                        39.0 kHz
	SX126X_GFSK_RX_BW_46_9       = 0x14 //  7     0                        46.9 kHz
	SX126X_GFSK_RX_BW_58_6       = 0x0C //  7     0                        58.6 kHz
	SX126X_GFSK_RX_BW_78_2       = 0x1B //  7     0                        78.2 kHz
	SX126X_GFSK_RX_BW_93_8       = 0x13 //  7     0                        93.8 kHz
	SX126X_GFSK_RX_BW_117_3      = 0x0B //  7     0                        117.3 kHz
	SX126X_GFSK_RX_BW_156_2      = 0x1A //  7     0                        156.2 kHz
	SX126X_GFSK_RX_BW_187_2      = 0x12 //  7     0                        187.2 kHz
	SX126X_GFSK_RX_BW_234_3      = 0x0A //  7     0                        234.3 kHz
	SX126X_GFSK_RX_BW_312_0      = 0x19 //  7     0                        312.0 kHz
	SX126X_GFSK_RX_BW_373_6      = 0x11 //  7     0                        373.6 kHz
	SX126X_GFSK_RX_BW_467_0      = 0x09 //  7     0                        467.0 kHz
	SX126X_LORA_BW_7_8           = 0x00 //  7     0     LoRa bandwidth: 7.8 kHz
	SX126X_LORA_BW_10_4          = 0x08 //  7     0                     10.4 kHz
	SX126X_LORA_BW_15_6          = 0x01 //  7     0                     15.6 kHz
	SX126X_LORA_BW_20_8          = 0x09 //  7     0                     20.8 kHz
	SX126X_LORA_BW_31_25         = 0x02 //  7     0                     31.25 kHz
	SX126X_LORA_BW_41_7          = 0x0A //  7     0                     41.7 kHz
	SX126X_LORA_BW_62_5          = 0x03 //  7     0                     62.5 kHz
	SX126X_LORA_BW_125_0         = 0x04 //  7     0                     125.0 kHz
	SX126X_LORA_BW_250_0         = 0x05 //  7     0                     250.0 kHz
	SX126X_LORA_BW_500_0         = 0x06 //  7     0                     500.0 kHz

	//SX126X_CMD_SET_PACKET_PARAMS
	SX126X_GFSK_PREAMBLE_DETECT_OFF         = 0x00 //  7     0     GFSK minimum preamble length before reception starts: detector disabled
	SX126X_GFSK_PREAMBLE_DETECT_8           = 0x04 //  7     0                                                           8 bits
	SX126X_GFSK_PREAMBLE_DETECT_16          = 0x05 //  7     0                                                           16 bits
	SX126X_GFSK_PREAMBLE_DETECT_24          = 0x06 //  7     0                                                           24 bits
	SX126X_GFSK_PREAMBLE_DETECT_32          = 0x07 //  7     0                                                           32 bits
	SX126X_GFSK_ADDRESS_FILT_OFF            = 0x00 //  7     0     GFSK address filtering: disabled
	SX126X_GFSK_ADDRESS_FILT_NODE           = 0x01 //  7     0                             node only
	SX126X_GFSK_ADDRESS_FILT_NODE_BROADCAST = 0x02 //  7     0                             node and broadcast
	SX126X_GFSK_PACKET_FIXED                = 0x00 //  7     0     GFSK packet type: fixed (payload length known in advance to both sides)
	SX126X_GFSK_PACKET_VARIABLE             = 0x01 //  7     0                       variable (payload length added to packet)
	SX126X_GFSK_CRC_OFF                     = 0x01 //  7     0     GFSK packet CRC: disabled
	SX126X_GFSK_CRC_1_BYTE                  = 0x00 //  7     0                      1 byte
	SX126X_GFSK_CRC_2_BYTE                  = 0x02 //  7     0                      2 byte
	SX126X_GFSK_CRC_1_BYTE_INV              = 0x04 //  7     0                      1 byte, inverted
	SX126X_GFSK_CRC_2_BYTE_INV              = 0x06 //  7     0                      2 byte, inverted
	SX126X_GFSK_WHITENING_OFF               = 0x00 //  7     0     GFSK data whitening: disabled
	SX126X_GFSK_WHITENING_ON                = 0x01 //  7     0                          enabled

	//SX126X_CMD_SET_CAD_PARAMS
	SX126X_CAD_ON_1_SYMB  = 0x00 //  7     0     number of symbols used for CAD: 1
	SX126X_CAD_ON_2_SYMB  = 0x01 //  7     0                                     2
	SX126X_CAD_ON_4_SYMB  = 0x02 //  7     0                                     4
	SX126X_CAD_ON_8_SYMB  = 0x03 //  7     0                                     8
	SX126X_CAD_ON_16_SYMB = 0x04 //  7     0                                     16
	SX126X_CAD_GOTO_STDBY = 0x00 //  7     0     after CAD is done, always go to STDBY_RC mode
	SX126X_CAD_GOTO_RX    = 0x01 //  7     0     after CAD is done, go to Rx mode if activity is detected

	//SX126X_CMD_GET_STATUS
	SX126X_STATUS_MODE_STDBY_RC   = 0b00100000 //  6     4     current chip mode: STDBY_RC
	SX126X_STATUS_MODE_STDBY_XOSC = 0b00110000 //  6     4                        STDBY_XOSC
	SX126X_STATUS_MODE_FS         = 0b01000000 //  6     4                        FS
	SX126X_STATUS_MODE_RX         = 0b01010000 //  6     4                        RX
	SX126X_STATUS_MODE_TX         = 0b01100000 //  6     4                        TX
	SX126X_STATUS_DATA_AVAILABLE  = 0b00000100 //  3     1     command status: packet received and data can be retrieved
	SX126X_STATUS_CMD_TIMEOUT     = 0b00000110 //  3     1                     SPI command timed out
	SX126X_STATUS_CMD_INVALID     = 0b00001000 //  3     1                     invalid SPI command
	SX126X_STATUS_CMD_FAILED      = 0b00001010 //  3     1                     SPI command failed to execute
	SX126X_STATUS_TX_DONE         = 0b00001100 //  3     1                     packet transmission done
	SX126X_STATUS_SPI_FAILED      = 0b11111111 //  7     0     SPI transaction failed

	//SX126X_CMD_GET_PACKET_STATUS
	SX126X_GFSK_RX_STATUS_PREAMBLE_ERR    = 0b10000000 //  7     7     GFSK Rx status: preamble error
	SX126X_GFSK_RX_STATUS_SYNC_ERR        = 0b01000000 //  6     6                     sync word error
	SX126X_GFSK_RX_STATUS_ADRS_ERR        = 0b00100000 //  5     5                     address error
	SX126X_GFSK_RX_STATUS_CRC_ERR         = 0b00010000 //  4     4                     CRC error
	SX126X_GFSK_RX_STATUS_LENGTH_ERR      = 0b00001000 //  3     3                     length error
	SX126X_GFSK_RX_STATUS_ABORT_ERR       = 0b00000100 //  2     2                     abort error
	SX126X_GFSK_RX_STATUS_PACKET_RECEIVED = 0b00000010 //  2     2                     packet received
	SX126X_GFSK_RX_STATUS_PACKET_SENT     = 0b00000001 //  2     2                     packet sent

	//SX126X_CMD_GET_DEVICE_ERRORS
	SX126X_PA_RAMP_ERR     = 0b100000000 //  8     8     device errors: PA ramping failed
	SX126X_PLL_LOCK_ERR    = 0b001000000 //  6     6                    PLL failed to lock
	SX126X_XOSC_START_ERR  = 0b000100000 //  5     5                    crystal oscillator failed to start
	SX126X_IMG_CALIB_ERR   = 0b000010000 //  4     4                    image calibration failed
	SX126X_ADC_CALIB_ERR   = 0b000001000 //  3     3                    ADC calibration failed
	SX126X_PLL_CALIB_ERR   = 0b000000100 //  2     2                    PLL calibration failed
	SX126X_RC13M_CALIB_ERR = 0b000000010 //  1     1                    RC13M calibration failed
	SX126X_RC64K_CALIB_ERR = 0b000000001 //  0     0                    RC64K calibration failed

	// SX126X SPI register variables
	//SX126X_REG_LORA_SYNC_WORD_MSB + LSB
	SX126X_SYNC_WORD_PUBLIC  = 0x34 // actually 0x3444  NOTE: The low nibbles in each byte (0x_4_4) are masked out since apparently, they're reserved.
	SX126X_SYNC_WORD_PRIVATE = 0x12 // actually 0x1424        You couldn't make this up if you tried.

	SX126X_LORA_MAC_PUBLIC_SYNCWORD  = 0x3444
	SX126X_LORA_MAC_PRIVATE_SYNCWORD = 0x1424
)
