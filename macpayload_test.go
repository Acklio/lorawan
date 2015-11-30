package lorawan

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMACPayload(t *testing.T) {
	Convey("Given an empty MACPayload", t, func() {
		var p MACPayload
		Convey("Then MarshalBinary returns []byte{0, 0, 0, 0, 0, 0, 0}", func() {
			b, err := p.MarshalBinary()
			So(err, ShouldBeNil)
			So(b, ShouldResemble, []byte{0, 0, 0, 0, 0, 0, 0})
		})

		Convey("Given FPort=1", func() {
			fPort := uint8(1)
			p.FPort = &fPort
			Convey("Then MarshalBinary returns an error", func() {
				_, err := p.MarshalBinary()
				So(err, ShouldResemble, errors.New("lorawan: FPort should not be set when FRMPayload is empty"))
			})

			Convey("Given FRMPayload contains MACCommand(CID=LinkCheckReq)", func() {
				p.FRMPayload = []Payload{&MACCommand{CID: LinkCheckReq}}
				Convey("Then MarshalBinary returns an error", func() {
					_, err := p.MarshalBinary()
					So(err, ShouldResemble, errors.New("lorawan: a MAC command is only allowed when FPort=0"))
				})
			})
		})

		Convey("Given FHDR(DevAddr=67305985), FPort=1, FRMPayload=[]Payload{DataPayload(Bytes=[]byte{5, 6, 7})}", func() {
			fPort := uint8(1)
			p.FHDR.DevAddr = DevAddr(67305985)
			p.FPort = &fPort
			p.FRMPayload = []Payload{&DataPayload{[]byte{5, 6, 7}}}

			Convey("Then MarshalBinary returns []byte{1, 2, 3, 4, 0, 0, 0, 1, 5, 6, 7}", func() {
				b, err := p.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{1, 2, 3, 4, 0, 0, 0, 1, 5, 6, 7})
			})

			Convey("Given the key []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}", func() {
				key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
				Convey("Then EncryptPayload does not return an error", func() {
					err := p.EncryptPayload(key)
					So(err, ShouldBeNil)

					Convey("Then FRMPayload contains one DataPayload not equal to DataPayload(Bytes=[]byte{5, 6, 7})", func() {
						So(p.FRMPayload, ShouldHaveLength, 1)
						data, ok := p.FRMPayload[0].(*DataPayload)
						So(ok, ShouldBeTrue)
						So(data.Bytes, ShouldNotResemble, []byte{5, 6, 7})
					})

					Convey("Then DecryptPayload does not return an error", func() {
						err := p.DecryptPayload(key)
						So(err, ShouldBeNil)

						Convey("Then FRMPayload contains one DataPayload(Bytes=[]byte{5, 6, 7})", func() {
							So(p.FRMPayload, ShouldHaveLength, 1)
							data, ok := p.FRMPayload[0].(*DataPayload)
							So(ok, ShouldBeTrue)
							So(data.Bytes, ShouldResemble, []byte{5, 6, 7})
						})
					})
				})
			})
		})

		Convey("Given uplink=true, FHDR(DevAddr=67305985), FPort=0, FRMPayload=[]Payload{MACCommand{CID: DevStatusAns, Payload: DevStatusAnsPayload(Battery=10, Margin=20)}}", func() {
			fPort := uint8(0)
			p.uplink = true
			p.FHDR.DevAddr = DevAddr(67305985)
			p.FPort = &fPort
			p.FRMPayload = []Payload{&MACCommand{CID: DevStatusAns, Payload: &DevStatusAnsPayload{Battery: 10, Margin: 20}}}

			Convey("Then MarshalBinary returns []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10, 20}", func() {
				b, err := p.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10, 20})
			})
		})

		Convey("Given the slice []byte{1, 2, 3, 4, 0, 0}", func() {
			b := []byte{1, 2, 3, 4, 0, 0}
			Convey("Then UnmarshalBinary returns an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldResemble, errors.New("lorawan: at least 7 bytes needed to decode FHDR"))
			})
		})

		Convey("Given the slice []byte{1, 2, 3, 4, 3, 0, 0, 0, 0}", func() {
			b := []byte{1, 2, 3, 4, 3, 0, 0, 0, 0}
			Convey("Then UnmarshalBinary returns an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldResemble, errors.New("lorawan: not enough bytes to decode FHDR"))
			})
		})

		Convey("Given the slice []byte{1, 2, 3, 4, 0, 0, 0, 1}", func() {
			b := []byte{1, 2, 3, 4, 0, 0, 0, 1}
			Convey("Then UnmarshalBinary returns an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldResemble, errors.New("lorawan: data contains FPort but no FRMPayload"))
			})
		})

		Convey("Given uplink=true and slice []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10}", func() {
			b := []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10}
			p.uplink = true
			Convey("Then UnmarshalBinary returns an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldResemble, errors.New("lorawan: not enough remaining bytes"))
			})
		})

		Convey("Given uplink=true and slice []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10, 20}", func() {
			b := []byte{1, 2, 3, 4, 0, 0, 0, 0, 6, 10, 20}
			p.uplink = true

			Convey("Then UnmarshalBinary does not return an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldBeNil)

				Convey("Then FHDR(DevAddr=67305985)", func() {
					So(p.FHDR.DevAddr, ShouldEqual, DevAddr(67305985))
				})
				Convey("Then FPort=0", func() {
					So(*p.FPort, ShouldEqual, 0)
				})
				Convey("Then FRMPayload=[]Payload{MACCommand{CID: DevStatusAns, Payload: DevStatusAnsPayload(Battery=10, Margin=20)}}", func() {
					So(p.FRMPayload, ShouldHaveLength, 1)

					mac, ok := p.FRMPayload[0].(*MACCommand)
					So(ok, ShouldBeTrue)
					So(mac.CID, ShouldEqual, DevStatusAns)

					pl, ok := mac.Payload.(*DevStatusAnsPayload)
					So(ok, ShouldBeTrue)
					So(pl.Battery, ShouldEqual, 10)
					So(pl.Margin, ShouldEqual, 20)
				})
			})
		})

		Convey("Given the slice []byte{1, 2, 3, 4, 0, 0, 0, 1, 6, 10, 20}", func() {
			b := []byte{1, 2, 3, 4, 0, 0, 0, 1, 6, 10, 20}

			Convey("Then UnmarshalBinary does not return an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldBeNil)

				Convey("Then FHDR(DevAddr=67305985)", func() {
					So(p.FHDR.DevAddr, ShouldEqual, DevAddr(67305985))
				})
				Convey("Then FPort=1", func() {
					So(*p.FPort, ShouldEqual, 1)
				})
				Convey("Then FRMPayload=[]Payload{DataPayload([]byte{6, 10, 20})}", func() {
					So(p.FRMPayload, ShouldHaveLength, 1)

					pl, ok := p.FRMPayload[0].(*DataPayload)
					So(ok, ShouldBeTrue)
					So(pl.Bytes, ShouldResemble, []byte{6, 10, 20})
				})
			})
		})
	})
}
