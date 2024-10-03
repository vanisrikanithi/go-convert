package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type testRunner struct {
	name  string
	input Node
	want  *harness.Step
}

func prepareFileOpsTest(t *testing.T, filename string, folderName string, step *harness.Step) testRunner {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/fileOps/"+folderName, filename+".json"))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node
	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	return testRunner{
		name:  filename,
		input: inputNode,
		want:  step,
	}
}

func TestConvertFileOpsCopyFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsCopy_snippet", "fileOpsCopy", &harness.Step{
		Id:   "fileOperationsa39de5",
		Name: "fileCopyOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("cp -r src/*.txt dest/"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileCopy(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileCopy() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func extractAnanymousOperation(currentNode Node) map[string]interface{} {
	// Step 1: Extract the 'delegate' map from the 'parameterMap'
	delegate, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'delegate' in parameterMap")
	}

	// Step 2: Extract the 'arguments' map from the 'delegate'
	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'arguments' in delegate map")
	}

	// Step 3: Extract the list of anonymous operations
	anonymousOps, ok := arguments["<anonymous>"].([]interface{})
	if !ok {
		fmt.Println("No anonymous operations found in arguments")
	}
	var extractedOperation map[string]interface{}
	// Step 4: Iterate over each operation and handle based on the 'symbol' type
	for _, op := range anonymousOps {
		// Convert the operation to a map for easy access
		operation, ok := op.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid operation format")
			continue
		}
		extractedOperation = operation
	}
	return extractedOperation
}
