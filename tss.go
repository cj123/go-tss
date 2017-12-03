package tss

/*
#define HAVE_STRSEP 1
#include <plist/plist.h>
#include "idevicerestore/src/tss.c"
#include "idevicerestore/src/common.c"

#cgo LDFLAGS: -lplist -lcurl -L${SRCDIR}/idevicerestore/src

void disableMessages() {
	error_disabled = 1;
	debug_disabled = 1;
	info_disabled = 1;
}

*/
import "C"
import (
	"errors"
	"unicode/utf8"
	"unsafe"

	"howett.net/plist"
)

func DisableMessages() {
	C.disableMessages()
}

type Request struct {
	plist C.plist_t
}

var ErrInvalidStatus = errors.New("tss: invalid idevicerestore/tss status, is less than zero")

func NewRequest(overrides interface{}) (*Request, error) {
	pp, err := MarshalCPlist(overrides)

	if err != nil {
		return nil, err
	}

	return &Request{
		C.tss_request_new(pp),
	}, nil
}

func AddParametersFromManifest(p interface{}, buildIdentity interface{}) (interface{}, error) {
	pls, err := MarshalCPlist(p)

	if err != nil {
		return nil, err
	}

	id, err := MarshalCPlist(buildIdentity)

	if err != nil {
		return nil, err
	}

	status := C.tss_parameters_add_from_manifest(pls, id)

	if status < 0 {
		return nil, ErrInvalidStatus
	}

	err = UnmarshalCPlist(pls, &p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Request) AddCommonTags(params interface{}, overrides interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	over, err := MarshalCPlist(overrides)

	if err != nil {
		return err
	}

	status := C.tss_request_add_common_tags(t.plist, p, over)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddAPTags(params interface{}, overrides interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	over, err := MarshalCPlist(overrides)

	if err != nil {
		return err
	}

	status := C.tss_request_add_ap_tags(t.plist, p, over)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddBaseBandTags(params interface{}, overrides interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	over, err := MarshalCPlist(overrides)

	if err != nil {
		return err
	}

	status := C.tss_request_add_baseband_tags(t.plist, p, over)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddSETags(params interface{}, overrides interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	over, err := MarshalCPlist(overrides)

	if err != nil {
		return err
	}

	status := C.tss_request_add_se_tags(t.plist, p, over)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddAPImg4Tags(params interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	status := C.tss_request_add_ap_img4_tags(t.plist, p)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddAPImg3Tags(params interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	status := C.tss_request_add_ap_img3_tags(t.plist, p)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
}

func (t *Request) AddSavageTags(params interface{}, overrides interface{}) error {
	p, err := MarshalCPlist(params)

	if err != nil {
		return err
	}

	over, err := MarshalCPlist(overrides)

	if err != nil {
		return err
	}

	status := C.tss_request_add_savage_tags(t.plist, p, over)
	if status < 0 {
		return ErrInvalidStatus
	}

	return nil
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

func MarshalCPlist(data interface{}) (C.plist_t, error) {
	if data == nil {
		return nil, nil
	}

	val, err := plist.Marshal(data, plist.XMLFormat)

	if err != nil {
		return nil, err
	}

	var x C.plist_t

	cVal := C.CString(string(val))
	defer C.free(unsafe.Pointer(cVal))

	C.plist_from_xml(cVal, C.uint32_t(utf8.RuneCount(val)), &x)

	return x, nil
}

func UnmarshalCPlist(cPlist C.plist_t, d interface{}) error {
	str := C.CString("")
	defer C.free(unsafe.Pointer(str))
	defer C.free(unsafe.Pointer(cPlist))

	l := C.uint32_t(0)

	C.plist_to_xml(cPlist, &str, &l)

	if l <= 0 {
		return errors.New("invalid plist")
	}

	val := C.GoString(str)

	_, err := plist.Unmarshal([]byte(val), d)

	return err
}
