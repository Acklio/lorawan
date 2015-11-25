package lorawan

import (
	"encoding/binary"
	"errors"
)

// CID defines the MAC command identifier.
type CID byte

// MAC commands as specified by the LoRaWAN R1.0 specs. Note that each *Req / *Ans
// has the same value. Based on the fact if a message is uplink or downlink
// you should use on or the other.
const (
	LinkCheckReq     CID = 0x02
	LinkCheckAns     CID = 0x02
	LinkADRReq       CID = 0x03
	LinkADRAns       CID = 0x03
	DutyCycleReq     CID = 0x04
	DutyCycleAns     CID = 0x04
	RXParamSetupReq  CID = 0x05
	RXParamSetupAns  CID = 0x05
	DevStatusReq     CID = 0x06
	DevStatusAns     CID = 0x06
	NewChannelReq    CID = 0x07
	NewChannelAns    CID = 0x07
	RXTimingSetupReq CID = 0x08
	RXTimingSetupAns CID = 0x08
	// 0x80 to 0xFF reserved for proprietary network command extensions
)

// LinkCheckAnsPayload represents the LinkCheckAns payload.
type LinkCheckAnsPayload struct {
	Margin uint8
	GwCnt  uint8
}

// ChMask represents the channel mask.
type ChMask [2]byte

// NewChMask returns a new ChMask for the given channel numbers (max. is 16).
func NewChMask(chans ...uint8) (ChMask, error) {
	var mask ChMask
	for _, c := range chans {
		if c > 16 {
			return mask, errors.New("lorawan: the max. channel number is 16")
		}
		c = c - 1 // make it zero based indexed
		i := c / 8
		b := c % 8
		mask[i] = mask[i] ^ 1<<b
	}
	return mask, nil
}

// Channels returns the channels active in the channel mask.
func (m ChMask) Channels() []uint8 {
	var chans []uint8
	for c := uint8(0); c < 16; c++ {
		i := c / 8
		b := c % 8
		if m[i]&(1<<b) > 0 {
			chans = append(chans, c+1)
		}
	}
	return chans
}

// Redundacy represents the redundacy field.
type Redundacy byte

// NewRedundacy returns a new Redundacy. Max allowed value for chMaskCntl is 7,
// max allowed value for nbRep is 15.
func NewRedundacy(chMaskCntl, nbRep uint8) (Redundacy, error) {
	var r Redundacy
	if chMaskCntl > 7 {
		return r, errors.New("lorawan: max value of chMaskCntl is 7")
	}
	if nbRep > 15 {
		return r, errors.New("lorawan: max value of nbRep is 15")
	}

	return Redundacy((chMaskCntl << 4) ^ nbRep), nil
}

// ChMaskCntl (channel mask control) controls the interpretation of the ChMask
// bit field.
func (r Redundacy) ChMaskCntl() uint8 {
	var mask uint8 = (1 << 6) ^ (1 << 5) ^ (1 << 4)
	return (uint8(r) & mask) >> 4
}

// NbRep returns the number of repetition for each uplink message.
func (r Redundacy) NbRep() uint8 {
	var mask uint8 = (1 << 3) ^ (1 << 2) ^ (1 << 1) ^ (1 << 0)
	return uint8(r) & mask
}

// DataRateTXPower represents the requested data rate and TX output power.
type DataRateTXPower byte

// NewDataRateTXPower returns a new DataRateTXPower. Max allowed values for
// dataRate and txPower are 15.
func NewDataRateTXPower(dataRate, txPower uint8) (DataRateTXPower, error) {
	var dr DataRateTXPower
	if dataRate > 15 {
		return dr, errors.New("lorawan: max value for dataRate is 15")
	}
	if txPower > 15 {
		return dr, errors.New("lorawan: max value for txPower is 15")
	}
	return DataRateTXPower((dataRate << 4) ^ txPower), nil
}

// DataRate returns the requested data rate.
func (dr DataRateTXPower) DataRate() uint8 {
	return uint8(dr) >> 4
}

// TXPower returns the requested TX output power.
func (dr DataRateTXPower) TXPower() uint8 {
	var mask uint8 = (1 << 3) ^ (1 << 2) ^ (1 << 1) ^ (1 << 0)
	return uint8(dr) & mask
}

// LinkADRReqPayload represents the LinkADRReq payload.
type LinkADRReqPayload struct {
	DataRateTXPower DataRateTXPower
	ChMask          ChMask
	Redundacy       Redundacy
}

// LinkADRAnsPayload represents the LinkADRAns payload.
type LinkADRAnsPayload byte

// NewLinkADRAnsPayload returns a new LinkADRAnsPayload containing the given options.
func NewLinkADRAnsPayload(chMaskACK, dataRateACK, powerACK bool) LinkADRAnsPayload {
	var p LinkADRAnsPayload
	if chMaskACK {
		p = p ^ (1 << 0)
	}
	if dataRateACK {
		p = p ^ (1 << 1)
	}
	if powerACK {
		p = p ^ (1 << 2)
	}
	return p
}

// ChMaskACK returns if the channel mask sent was successfully interpreted.
func (p LinkADRAnsPayload) ChMaskACK() bool {
	return p&(1<<0) > 0
}

// DataRateACK returns if the data rate was successfylly set.
func (p LinkADRAnsPayload) DataRateACK() bool {
	return p&(1<<1) > 0
}

// PowerACK returns if the power level was successfully set.
func (p LinkADRAnsPayload) PowerACK() bool {
	return p&(1<<2) > 0
}

// DutyCycleReqPayload contains the MaxDCycle value.
type DutyCycleReqPayload uint8

// NewDutyCycleReqPayload returns a new DutyCycleReqPayload for the given MaxDCycle.
func NewDutyCycleReqPayload(maxDCycle uint8) (DutyCycleReqPayload, error) {
	if maxDCycle > 15 && maxDCycle < 255 {
		return 0, errors.New("lorawan: only a MaxDCycle value of 0 - 15 and 255 is allowed")
	}
	return DutyCycleReqPayload(maxDCycle), nil
}

// DLsettings represents the downlink settings.
type DLsettings byte

// RX2DataRate returns the requested data rate.
func (s DLsettings) RX2DataRate() uint8 {
	var mask DLsettings = (1 << 3) ^ (1 << 2) ^ (1 << 1) ^ (1 << 0)
	return uint8(s & mask)
}

// RX1DRoffset returns the offset between uplink data rate and the downlink data rate.
func (s DLsettings) RX1DRoffset() uint8 {
	var mask DLsettings = (1 << 6) ^ (1 << 5) ^ (1 << 4)
	return uint8(s&mask) >> 4
}

// NewDLsettings returns a new DLsettings for the given RX2DataRate and RX1DRoffset.
func NewDLsettings(rx2DataRate, rx1DRoffset uint8) (DLsettings, error) {
	if rx2DataRate > 15 {
		return 0, errors.New("lorawan: max value for rx2DataRate is 15")
	}
	if rx1DRoffset > 7 {
		return 0, errors.New("lorawan: max value for rx1DRoffset is 7")
	}
	return DLsettings(rx2DataRate ^ (rx1DRoffset << 4)), nil
}

// Frequency defines the frequency which is a 24 bits unsigned integer.
type Frequency [3]byte

// NewFrequency returns a new Frequency. Note that the max. allowed value is
// 24 bit (thus 2^24 - 1).
func NewFrequency(frequency uint32) (Frequency, error) {
	var freq Frequency
	if frequency >= 2^24 {
		return freq, errors.New("lorawan: max value for frequency is 2^24-1")
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, frequency)
	for i := 0; i < 3; i++ {
		freq[0] = b[0]
	}
	return freq, nil
}

// Uint32 returns the frequency value as an uint32.
func (f Frequency) Uint32() uint32 {
	b := make([]byte, 4)
	for i, v := range f {
		b[i] = v
	}
	return binary.LittleEndian.Uint32(b)
}

// RX2SetupReqPayload represents the second receive window parameters.
type RX2SetupReqPayload struct {
	DLsettings DLsettings
	Frequency  Frequency
}

// RX2SetupAnsPayload represents payload send by the RXParamSetupAns command.
type RX2SetupAnsPayload byte

// NewRX2SetupAnsPayload returns a new RX2SetupAnsPayload.
func NewRX2SetupAnsPayload(channelACK, rx2DataRateACK, rx1DRoffsetACK bool) RX2SetupAnsPayload {
	var p RX2SetupAnsPayload
	if channelACK {
		p = p ^ (1 << 0)
	}
	if rx2DataRateACK {
		p = p ^ (1 << 1)
	}
	if rx1DRoffsetACK {
		p = p ^ (1 << 2)
	}
	return p
}

// ChannelACK returns if the RX2 slot was successfully set.
func (p RX2SetupAnsPayload) ChannelACK() bool {
	return p&(1<<0) > 0
}

// RX2DataRateACK returns if the RX2 slot data rate was successfully set.
func (p RX2SetupAnsPayload) RX2DataRateACK() bool {
	return p&(1<<1) > 0
}

// RX1DRoffsetACK return if the RX1DRoffset was successfully set.
func (p RX2SetupAnsPayload) RX1DRoffsetACK() bool {
	return p&(1<<2) > 0
}
