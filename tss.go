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

type Request struct {
	plist C.plist_t
}

func NewRequest(overrides interface{}) *Request {
	return &Request{
		C.tss_request_new(MarshalCPlist(overrides)),
	}
}

func AddParametersFromManifest(p interface{}, buildIdentity interface{}) interface{} {
	pls := MarshalCPlist(p)

	C.tss_parameters_add_from_manifest(pls, MarshalCPlist(buildIdentity))

	UnmarshalCPlist(pls, &p)

	return p
}

func (t *Request) AddCommonTags(params interface{}, overrides interface{}) {
	C.tss_request_add_common_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *Request) AddAPTags(params interface{}, overrides interface{}) {
	C.tss_request_add_ap_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *Request) AddBaseBandTags(params interface{}, overrides interface{}) {
	C.tss_request_add_baseband_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *Request) AddSETags(params interface{}, overrides interface{}) {
	C.tss_request_add_se_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *Request) AddAPImg4Tags(params interface{}) {
	C.tss_request_add_ap_img4_tags(t.plist, MarshalCPlist(params))
}

func (t *Request) AddAPImg3Tags(params interface{}) {
	C.tss_request_add_ap_img3_tags(t.plist, MarshalCPlist(params))
}

func (t *Request) AddSavageTags(params interface{}, overrides interface{}) {
	C.tss_request_add_savage_tags(t.plist, MarshalCPlist(params), MarshalCPlist(overrides))
}

func (t *Request) Send() (map[string]interface{}, error) {
	response := C.tss_request_send(t.plist, nil)

	var p map[string]interface{}

	err := UnmarshalCPlist(response, &p)

	if err != nil {
		return nil, err
	}

	return p, err
}

func (t *Request) Bytes() ([]byte, error) {
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
