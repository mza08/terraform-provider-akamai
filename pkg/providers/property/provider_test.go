package property

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

var testProvider *schema.Provider

func TestMain(m *testing.M) {
	testProvider = akamai.Provider(Subprovider())()
	testAccProviders = map[string]*schema.Provider{
		"akamai": testProvider,
	}
	if err := akamai.TFTestSetup(); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := akamai.TFTestTeardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client papi.PAPI, hapiClient hapi.HAPI, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	origHapi := inst.hapiClient
	inst.hapiClient = hapiClient

	defer func() {
		inst.client = orig
		inst.hapiClient = origHapi
		clientLock.Unlock()
	}()

	f()
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// loadFixtureString returns the entire contents of the given file as a string
func loadFixtureString(format string, args ...interface{}) string {
	return string(loadFixtureBytes(fmt.Sprintf(format, args...)))
}

// suppressLogging prevents logging output during the given func unless TEST_LOGGING env var is not empty. Use this
// to keep log messages from polluting test output. Not thread-safe.
func suppressLogging(t *testing.T, f func()) {
	t.Helper()

	if os.Getenv("TEST_LOGGING") == "" {
		orig := hclog.SetDefault(hclog.NewNullLogger())
		defer func() { hclog.SetDefault(orig) }()
		t.Log("Logging is suppressed. Set TEST_LOGGING=1 in env to see logged messages during test")
	}

	f()
}

// Wrapper to intercept the papi.Mock's call of t.FailNow(). The Terraform test driver runs the provider code on
// goroutines other than the one created for the test. When t.FailNow() is called from any other goroutine, it causes
// the test to hang because the TF test driver is still waiting to serve requests. Mockery's failure message neglects to
// inform the user which test had failed. Use this struct to wrap a *testing.T when you call mock.Test(T{t}) and the
// mock's failure will print the failling test's name. Such failures are usually caused by the provider invoking an
// unexpected call on the mock.
//
// NB: You should only need to use this where your test uses the Terraform test driver
type T struct{ *testing.T }

// Overrides testing.T.FailNow() so when a test mock fails an assertion, we see which test had failed before it hangs
func (t T) FailNow() {
	t.T.Fatalf("FAIL: %s", t.T.Name())
}
