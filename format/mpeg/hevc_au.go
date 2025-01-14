package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var hevcAUNALFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.HEVC_AU,
		Description: "H.265/HEVC Access Unit",
		DecodeFn:    hevcAUDecode,
		DefaultInArg: format.HevcAuIn{
			LengthSize: 4,
		},
		RootArray: true,
		RootName:  "access_unit",
		Dependencies: []decode.Dependency{
			{Names: []string{format.HEVC_NALU}, Group: &hevcAUNALFormat},
		},
	})
}

// TODO: share/refactor with avcAUDecode?
func hevcAUDecode(d *decode.D) any {
	var hi format.HevcAuIn
	d.ArgAs(&hi)

	if hi.LengthSize == 0 {
		// TODO: is annexb the correct name?
		annexBDecode(d, hevcAUNALFormat)
		return nil
	}

	for d.NotEnd() {
		d.FieldStruct("nalu", func(d *decode.D) {
			l := int64(d.FieldU("length", int(hi.LengthSize)*8)) * 8
			d.FieldFormatLen("nalu", l, hevcAUNALFormat, nil)
		})
	}

	return nil
}
