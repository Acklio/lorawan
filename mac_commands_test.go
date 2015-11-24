package lorawan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedundacy(t *testing.T) {
	Convey("Given an empty Redundacy", t, func() {
		var r Redundacy
		Convey("ChMaskCntl = 0 and NbRep = 0", func() {
			So(r.ChMaskCntl(), ShouldEqual, 0)
			So(r.NbRep(), ShouldEqual, 0)
		})
	})

	Convey("Given I use NewRedundacy to create a new Redundacy", t, func() {
		Convey("An error should be returned when chMaskCntl > 7", func() {
			_, err := NewRedundacy(8, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("An error should be returned when nbRep > 15", func() {
			_, err := NewRedundacy(0, 16)
			So(err, ShouldNotBeNil)
		})
		Convey("Given I call NewRedundacy(5, 11)", func() {
			r, err := NewRedundacy(5, 11)
			So(err, ShouldBeNil)
			Convey("ChMaskCntl() should return 5", func() {
				So(r.ChMaskCntl(), ShouldEqual, 5)
			})
			Convey("NbRep() should return 11", func() {
				So(r.NbRep(), ShouldEqual, 11)
			})
		})
	})
}

func TestDataRateTXPower(t *testing.T) {
	Convey("Given an empty DataRateTXPower", t, func() {
		var dr DataRateTXPower
		Convey("DataRate = 0 and TXPower = 0", func() {
			So(dr.DataRate(), ShouldEqual, 0)
			So(dr.TXPower(), ShouldEqual, 0)
		})
	})

	Convey("Given I use NewDataRateTXPower to create a new DataRateTXPower", t, func() {
		Convey("An error should be returned when dataRate > 15", func() {
			_, err := NewDataRateTXPower(16, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("An error should be returned when txPower > 15", func() {
			_, err := NewDataRateTXPower(0, 16)
			So(err, ShouldNotBeNil)
		})

		Convey("Given I call NewDataRateTXPower(11, 14)", func() {
			dr, err := NewDataRateTXPower(11, 14)
			So(err, ShouldBeNil)
			Convey("DataRate should be 11", func() {
				So(dr.DataRate(), ShouldEqual, 11)
			})
			Convey("TXPower should be 14", func() {
				So(dr.TXPower(), ShouldEqual, 14)
			})
		})
	})
}

func TestChMask(t *testing.T) {
	Convey("Given an empty ChMask", t, func() {
		var mask ChMask
		Convey("It should not contain any channels", func() {
			So(mask.Channels(), ShouldBeEmpty)
		})
	})

	Convey("Given I use NewChMask to create a new ChMask", t, func() {
		Convey("An error should be returned for a channel > 16", func() {
			_, err := NewChMask(1, 2, 17)
			So(err, ShouldNotBeNil)
		})

		Convey("Given I call NewChMask(1, 5, 7, 11)", func() {
			mask, err := NewChMask(1, 5, 7, 11)
			So(err, ShouldBeNil)
			Convey("Channels should return 1, 5, 7 and 11", func() {
				So(mask.Channels(), ShouldResemble, []uint8{1, 5, 7, 11})
			})
		})
	})
}

func TestLinkADRAnsPayload(t *testing.T) {
	Convey("Given an empty LinkADRAnsPayload", t, func() {
		var p LinkADRAnsPayload
		Convey("ChMaskACK, DataRateACK and PowerACK should be false", func() {
			So(p.ChMaskACK(), ShouldBeFalse)
			So(p.DataRateACK(), ShouldBeFalse)
			So(p.PowerACK(), ShouldBeFalse)
		})
	})

	Convey("Given I use NewLinkADRAnsPayload to create a new LinkADRAnsPayload", t, func() {
		Convey("Given I call NewLinkADRAnsPayload(true, false, false)", func() {
			p := NewLinkADRAnsPayload(true, false, false)
			Convey("ChMaskACK should be true", func() {
				So(p.ChMaskACK(), ShouldBeTrue)
				So(p.DataRateACK(), ShouldBeFalse)
				So(p.PowerACK(), ShouldBeFalse)
			})
		})

		Convey("Given I call NewLinkADRAnsPayload(true, true, false)", func() {
			p := NewLinkADRAnsPayload(true, true, false)
			Convey("ChMaskACK and DataRateACK should be true", func() {
				So(p.ChMaskACK(), ShouldBeTrue)
				So(p.DataRateACK(), ShouldBeTrue)
				So(p.PowerACK(), ShouldBeFalse)
			})
		})

		Convey("Given I call NewLinkADRAnsPayload(true, true, true)", func() {
			p := NewLinkADRAnsPayload(true, true, true)
			Convey("ChMaskACK DataRateACK and PowerACK should be true", func() {
				So(p.ChMaskACK(), ShouldBeTrue)
				So(p.DataRateACK(), ShouldBeTrue)
				So(p.PowerACK(), ShouldBeTrue)
			})
		})
	})
}

func TestDutyCycleReqPayload(t *testing.T) {
	Convey("Given I use NewDutyCycleReqPayload to create a new DutyCycleReqPayload", t, func() {
		Convey("A value > 15 should return an error", func() {
			_, err := NewDutyCycleReqPayload(16)
			So(err, ShouldNotBeNil)
		})
		Convey("A value < 255 should return an error", func() {
			_, err := NewDutyCycleReqPayload(254)
			So(err, ShouldNotBeNil)
		})
		Convey("A value < 15 should not return an error", func() {
			p, err := NewDutyCycleReqPayload(14)
			So(err, ShouldBeNil)
			So(p, ShouldEqual, DutyCycleReqPayload(14))
		})
	})
}
