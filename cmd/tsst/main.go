package main

import (
	"fmt"
	"github.com/cj123/go-tss"
)

func main() {
	/*client := api.NewIPSWClientLatest()

	ip, err := ipsw.NewIPSWWithIdentifierBuild(client, "iPhone3,1", "8B117")

	if err != nil {
		panic(err)
	}

	m, err := ip.BuildManifest()

	if err != nil {
		panic(err)
	}

	raw, err := ip.RawManifest()

	if err != nil {
		panic(err)
	}*/

	params := map[string]interface{}{
		"ApECID":           20,
		"ApNonce":          "",
		"ApSepNonce":       "",
		"ApProductionMode": "",
	}

	tssReq := tss.NewTSSRequest(nil)

	tssReq.AddCommonTags(params, nil)
	b, err := tssReq.Bytes()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
