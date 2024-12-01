package stage

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/nlm/adventofcode2023/internal/utils"
	"github.com/stretchr/testify/assert"
)

var flagStage = flag.Uint("stage", 1, "stage to run")
var flagInput = flag.String("input", "input", "input to read")
var flagRuns = flag.Uint("runs", 1, "number of runs to do")
var flagVerbose = flag.Bool("v", false, "verbose output")
var flagDebug = flag.Bool("debug", false, "debug")
var flagCpuProf = flag.String("cpuprof", "", "profile cpu")

var ErrUnimplemented = fmt.Errorf("unimplemented")

// StageFunc is the stage function signature.
type StageFunc func(input io.Reader) (any, error)

// RunCLI is a CLI helper to run stages.
func RunCLI(input any, fns ...StageFunc) {
	if !flag.Parsed() {
		flag.Parse()
	}
	// Find stage
	stage := int(*flagStage)
	if stage == 0 || stage > len(fns) {
		log.Fatalf("stage %d not found", stage)
	}
	fn := fns[stage-1]

	// Run
	// read input.txt if input is nil
	if input == nil {
		input = Open(*flagInput + ".txt")
	}
	if *flagRuns > 1 {
		input = utils.Must(io.ReadAll(utils.Must(Reader(input))))
	}
	var (
		minDuration   time.Duration
		maxDuration   time.Duration
		totalDuration time.Duration
	)

	// pprof
	if *flagCpuProf != "" {
		f := utils.Must(os.Create(*flagCpuProf))
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	fmt.Printf("%-6v %-20v %-20v %-20v\n", "STAGE", "INPUT", "RESULT", "TIME")
	for i := 0; i < int(*flagRuns); i++ {
		// Prepare reader
		reader, err := Reader(input)
		if err != nil {
			log.Fatal(err)
		}
		// Run
		start := time.Now()
		res, err := fn(reader)
		duration := time.Since(start)
		if err != nil {
			log.Fatal(fmt.Errorf("error: %v", err))
		}
		// Account times
		totalDuration += duration
		if minDuration == 0 || duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
		// Report completion
		if i == 0 {
			fmt.Printf("%-6v %-20v %-20v %-20v\n", stage, *flagInput, res, duration)
		}
	}

	if *flagRuns > 1 {
		fmt.Println()
		fmt.Printf("%-6v %-20v %-20v %-20v\n", "RUNS", "MIN", "MAX", "AVG")
		fmt.Printf("%-6v %-20v %-20v %-20v\n", *flagRuns, minDuration, maxDuration, totalDuration/time.Duration(*flagRuns))
	}
}

// TestCase represents the input and expected result of a test.
// Tests will individually run under the provided name.
// Input supports different types that will be converted to io.Reader.
// Result can be of any type, but it has to match the function result to succeed.
// Err is usually nil, but if we expect one, it can be given here.
type TestCase struct {
	// Name of the test
	Name string
	// Input can be []byte, string, or io.Reader
	Input any
	// Result can be of any type
	Result any
	// Error
	Err error
}

// Test runs the provided test cases against the StageFunc.
func Test(t *testing.T, fn StageFunc, cases []TestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Input == nil {
				tc.Input = Open(tc.Name + ".txt")
			}
			reader, err := Reader(tc.Input)
			if !assert.NoError(t, err) {
				return
			}
			res, err := fn(reader)
			if tc.Err != nil && !assert.ErrorIs(t, err, tc.Err) {
				return
			}
			assert.Equal(t, tc.Result, res)
		})
	}
}

// Test runs the provided test cases against the StageFunc.
func Benchmark(b *testing.B, fn StageFunc, cases []TestCase) {
	for _, tc := range cases {
		if tc.Input == nil {
			tc.Input = Open(tc.Name + ".txt")
		}
		reader, err := Reader(tc.Input)
		if !assert.NoError(b, err) {
			return
		}
		data, err := io.ReadAll(reader)
		if !assert.NoError(b, err) {
			return
		}
		b.Run(tc.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				res, err := fn(bytes.NewReader(data))
				if tc.Err != nil && !assert.ErrorIs(b, err, tc.Err) {
					return
				}
				assert.Equal(b, tc.Result, res)
			}
		})
	}
}

// Reader converts some classic input type as io.Reader.
func Reader(input any) (io.Reader, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if reader, ok := input.(io.Reader); ok {
		return reader, nil
	}
	switch in := input.(type) {
	case []byte:
		return bytes.NewReader(in), nil
	case string:
		return strings.NewReader(in), nil
	default:
		return nil, fmt.Errorf("unsupported input type: %t", input)
	}
}

var stageFS fs.FS

func SetFS(f fs.FS) {
	stageFS = utils.Must(fs.Sub(f, "data"))
}

func Open(name string) io.Reader {
	return utils.Must(stageFS.Open(name))
}

func Verbose() bool {
	return *flagVerbose
}

func Debug() bool {
	return *flagDebug
}
