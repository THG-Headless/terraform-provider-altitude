package provider

import (
	"crypto/rand"
	"math/big"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"altitude": providerserver.NewProtocol6WithError(New("test")()),
}

func randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}

	return "terraform-acc-test-" + string(ret)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALTITUDE_CLIENT_ID"); v == "" {
		t.Fatal("ALTITUDE_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ALTITUDE_CLIENT_SECRET"); v == "" {
		t.Fatal("ALTITUDE_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("ALTITUDE_MODE"); v == "" {
		t.Fatal("ALTITUDE_MODE must be set for acceptance tests")
	}
}
