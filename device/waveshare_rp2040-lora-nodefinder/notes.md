SX1262 DIO2 should be configured as `SetDIO2AsRfSwitchCtrl` - pin will go high when radio TX
DIO2 is tied to antenna switch along with GPIO17. Both are pulled low. When GPIO17 is high, 
antenna switch routes to RX circuit on radio.