package osexitcheckanalyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/xantinium/metrix/cmd/staticlint/osexitcheckanalyzer"
)

func TestOSExitCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexitcheckanalyzer.OSExitCheckAnalyzer, "./...")
}
