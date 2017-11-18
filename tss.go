package tss

/*
#include <plist/plist.h>
#include "./idevicerestore/src/tss.c"
#include "./idevicerestore/src/common.c"

#cgo LDFLAGS: -lplist -lcurl -L${SRCDIR}/idevicerestore/src
*/
import "C"
import (
	"errors"
	"unicode/utf8"

	"howett.net/plist"
)

type TSSRequest struct {
	plist C.plist_t
}

func NewTSSRequest(overrides interface{}) *TSSRequest {
	return &TSSRequest{
		C.tss_request_new(MarshalCPlist(overrides)),
	}
}

func (t *TSSRequest) AddParametersFromManifest(params interface{}, buildIdentity interface{}) {
	C.tss_parameters_add_from_manifest(MarshalCPlist(params), MarshalCPlist(buildIdentity))
}

func (t *TSSRequest) AddCommonTags(params interface{}, overrides interface{}) {
	C.tss_request_add_common_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *TSSRequest) AddAPTags(params interface{}, overrides interface{}) {
	C.tss_request_add_ap_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *TSSRequest) AddBaseBandTags(params interface{}, overrides interface{}) {
	C.tss_request_add_baseband_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *TSSRequest) AddSETags(params interface{}, overrides interface{}) {
	C.tss_request_add_se_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *TSSRequest) AddAPImg4Tags(params interface{}) {
	C.tss_request_add_ap_img4_tags(t.plist, MarshalCPlist(params))
}

func (t *TSSRequest) AddAPImg3Tags(params interface{}) {
	C.tss_request_add_ap_img3_tags(t.plist, MarshalCPlist(params))
}

func (t *TSSRequest) Bytes() ([]byte, error) {
	var p interface{}

	err := UnmarshalCPlist(t.plist, &p)

	if err != nil {
		return nil, err
	}

	return plist.Marshal(p, plist.XMLFormat)
}

func MarshalCPlist(data interface{}) C.plist_t {
	if data == nil {
		return nil
	}

	val, err := plist.Marshal(data, plist.XMLFormat)

	if err != nil {
		panic(err)
	}

	var x C.plist_t

	C.plist_from_xml(C.CString(string(val)), C.uint32_t(utf8.RuneCount(val)), &x)

	return x
}

func UnmarshalCPlist(cPlist C.plist_t, d interface{}) error {
	str := C.CString("")
	l := C.uint32_t(0)

	C.plist_to_xml(cPlist, &str, &l)

	if l <= 0 {
		return errors.New("invalid plist")
	}

	val := C.GoString(str)

	_, err := plist.Unmarshal([]byte(val), d)

	return err
}
