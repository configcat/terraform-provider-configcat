package resource

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	testing "github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testStepNewImportState(ctx context.Context, t testing.T, c TestCase, helper *plugintest.Helper, wd *plugintest.WorkingDir, step TestStep, cfg string) error {
	t.Helper()

	spewConf := spew.NewDefaultConfig()
	spewConf.SortKeys = true

	if step.ResourceName == "" {
		t.Fatal("ResourceName is required for an import state test")
	}

	// get state from check sequence
	var state *terraform.State
	var err error
	err = runProviderCommand(ctx, t, func() error {
		state, err = getState(ctx, t, wd)
		if err != nil {
			return err
		}
		return nil
	}, wd, providerFactories{
		legacy:  c.ProviderFactories,
		protov5: c.ProtoV5ProviderFactories,
		protov6: c.ProtoV6ProviderFactories})
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Determine the ID to import
	var importId string
	switch {
	case step.ImportStateIdFunc != nil:
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateIdFunc for import identifier")

		var err error

		logging.HelperResourceDebug(ctx, "Calling TestStep ImportStateIdFunc")

		importId, err = step.ImportStateIdFunc(state)

		if err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep ImportStateIdFunc")
	case step.ImportStateId != "":
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateId for import identifier")

		importId = step.ImportStateId
	default:
		logging.HelperResourceTrace(ctx, "Using resource identifier for import identifier")

		resource, err := testResource(step, state)
		if err != nil {
			t.Fatal(err)
		}
		importId = resource.Primary.ID
	}

	if step.ImportStateIdPrefix != "" {
		logging.HelperResourceTrace(ctx, "Prepending TestStep ImportStateIdPrefix for import identifier")

		importId = step.ImportStateIdPrefix + importId
	}

	logging.HelperResourceTrace(ctx, fmt.Sprintf("Using import identifier: %s", importId))

	// Create working directory for import tests
	if step.Config == "" {
		logging.HelperResourceTrace(ctx, "Using prior TestStep Config for import")

		step.Config = cfg
		if step.Config == "" {
			t.Fatal("Cannot import state with no specified config")
		}
	}
	importWd := helper.RequireNewWorkingDir(ctx, t)
	defer importWd.Close()
	err = importWd.SetConfig(ctx, step.Config)
	if err != nil {
		t.Fatalf("Error setting test config: %s", err)
	}

	logging.HelperResourceDebug(ctx, "Running Terraform CLI init and import")

	err = runProviderCommand(ctx, t, func() error {
		return importWd.Init(ctx)
	}, importWd, providerFactories{
		legacy:  c.ProviderFactories,
		protov5: c.ProtoV5ProviderFactories,
		protov6: c.ProtoV6ProviderFactories})
	if err != nil {
		t.Fatalf("Error running init: %s", err)
	}

	err = runProviderCommand(ctx, t, func() error {
		return importWd.Import(ctx, step.ResourceName, importId)
	}, importWd, providerFactories{
		legacy:  c.ProviderFactories,
		protov5: c.ProtoV5ProviderFactories,
		protov6: c.ProtoV6ProviderFactories})
	if err != nil {
		return err
	}

	var importState *terraform.State
	err = runProviderCommand(ctx, t, func() error {
		importState, err = getState(ctx, t, importWd)
		if err != nil {
			return err
		}
		return nil
	}, importWd, providerFactories{
		legacy:  c.ProviderFactories,
		protov5: c.ProtoV5ProviderFactories,
		protov6: c.ProtoV6ProviderFactories})
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Go through the imported state and verify
	if step.ImportStateCheck != nil {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateCheck")

		var states []*terraform.InstanceState
		for _, r := range importState.RootModule().Resources {
			if r.Primary != nil {
				is := r.Primary.DeepCopy()
				is.Ephemeral.Type = r.Type // otherwise the check function cannot see the type
				states = append(states, is)
			}
		}

		logging.HelperResourceDebug(ctx, "Calling TestStep ImportStateCheck")

		if err := step.ImportStateCheck(states); err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep ImportStateCheck")
	}

	// Verify that all the states match
	if step.ImportStateVerify {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateVerify")

		newResources := importState.RootModule().Resources
		oldResources := state.RootModule().Resources

		for _, r := range newResources {
			// Find the existing resource
			var oldR *terraform.ResourceState
			for r2Key, r2 := range oldResources {
				// Ensure that we do not match against data sources as they
				// cannot be imported and are not what we want to verify.
				// Mode is not present in ResourceState so we use the
				// stringified ResourceStateKey for comparison.
				if strings.HasPrefix(r2Key, "data.") {
					continue
				}

				if r2.Primary != nil && r2.Primary.ID == r.Primary.ID && r2.Type == r.Type && r2.Provider == r.Provider {
					oldR = r2
					break
				}
			}
			if oldR == nil || oldR.Primary == nil {
				t.Fatalf(
					"Failed state verification, resource with ID %s not found",
					r.Primary.ID)
			}

			// don't add empty flatmapped containers, so we can more easily
			// compare the attributes
			skipEmpty := func(k, v string) bool {
				if strings.HasSuffix(k, ".#") || strings.HasSuffix(k, ".%") {
					if v == "0" {
						return true
					}
				}
				return false
			}

			// Compare their attributes
			actual := make(map[string]string)
			for k, v := range r.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				actual[k] = v
			}

			expected := make(map[string]string)
			for k, v := range oldR.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				expected[k] = v
			}

			// Remove fields we're ignoring
			for _, v := range step.ImportStateVerifyIgnore {
				for k := range actual {
					if strings.HasPrefix(k, v) {
						delete(actual, k)
					}
				}
				for k := range expected {
					if strings.HasPrefix(k, v) {
						delete(expected, k)
					}
				}
			}

			// timeouts are only _sometimes_ added to state. To
			// account for this, just don't compare timeouts at
			// all.
			for k := range actual {
				if strings.HasPrefix(k, "timeouts.") {
					delete(actual, k)
				}
				if k == "timeouts" {
					delete(actual, k)
				}
			}
			for k := range expected {
				if strings.HasPrefix(k, "timeouts.") {
					delete(expected, k)
				}
				if k == "timeouts" {
					delete(expected, k)
				}
			}

			if !reflect.DeepEqual(actual, expected) {
				// Determine only the different attributes
				for k, v := range expected {
					if av, ok := actual[k]; ok && v == av {
						delete(expected, k)
						delete(actual, k)
					}
				}

				t.Fatalf(
					"ImportStateVerify attributes not equivalent. Difference is shown below. Top is actual, bottom is expected."+
						"\n\n%s\n\n%s",
					spewConf.Sdump(actual), spewConf.Sdump(expected))
			}
		}
	}

	return nil
}
