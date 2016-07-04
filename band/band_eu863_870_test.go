package band

import (
	"testing"

	"github.com/brocaar/lorawan"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEU863Band(t *testing.T) {
	Convey("Given the EU 863-870 band is selected", t, func() {
		band, err := GetConfig(EU_863_870)
		So(err, ShouldBeNil)

		Convey("Then GetRX1Channel returns the uplink channel", func() {
			for i := 0; i < 3; i++ {
				rx1Chan := band.GetRX1Channel(i)
				So(rx1Chan, ShouldEqual, i)
			}
		})

		Convey("Given a CFList", func() {
			cFlist := lorawan.CFList{
				867100000,
				867300000,
				867500000,
				867700000,
				867900000,
			}

			Convey("Then GetFrequency takes the CFList into consideration", func() {
				tests := []int{
					868100000,
					868300000,
					868500000,
					867100000,
					867300000,
					867500000,
					867700000,
					867900000,
				}

				for expChannel, expFreq := range tests {
					freq, err := band.GetDownlinkFrequency(expChannel, &cFlist)
					So(err, ShouldBeNil)
					So(freq, ShouldEqual, expFreq)

					channel, err := band.GetChannel(expFreq, &cFlist)
					So(err, ShouldBeNil)
					So(channel, ShouldEqual, expChannel)
				}
			})
		})
	})
}
