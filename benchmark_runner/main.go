package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

const (
	actionRunAnother = "another"
	actionChangeLang = "language"
	actionExit       = "exit"
)

// discoverLanguages finds all language folders.
// A language folder is a directory containing a run.sh script.
func discoverLanguages() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read repo root: %w", err)
	}

	var languages []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name[0] == '.' || name == "benchmark_runner" {
			continue
		}

		// A language folder must have a run.sh script
		runScript := filepath.Join(name, "run.sh")
		if _, err := os.Stat(runScript); err == nil {
			languages = append(languages, name)
		}
	}

	sort.Strings(languages)
	return languages, nil
}

// discoverBenchmarks finds all benchmark folders within a language folder.
// A benchmark folder is a subdirectory (excluding hidden directories).
func discoverBenchmarks(language string) ([]string, error) {
	entries, err := os.ReadDir(language)
	if err != nil {
		return nil, fmt.Errorf("failed to read language folder: %w", err)
	}

	var benchmarks []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name[0] == '.' {
			continue
		}
		benchmarks = append(benchmarks, name)
	}

	sort.Strings(benchmarks)
	return benchmarks, nil
}

// runBenchmark executes the language's run.sh script with the benchmark name as argument.
func runBenchmark(language, benchmark string, prefs *Preferences) error {
	fmt.Printf("\n🚀 Running benchmark: %s/%s (min-runs: %s)\n\n", language, benchmark, prefs.Get("min_runs"))

	cmd := exec.Command("bash", "run.sh", benchmark)
	cmd.Dir = language
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	env := os.Environ()
	env = append(env, "MIN_RUNS="+prefs.Get("min_runs"))
	if prefs.GetBool("save_results") {
		env = append(env, "EXPORT_RESULTS=1")
	}
	cmd.Env = env

	return cmd.Run()
}

func main() {
	languages, err := discoverLanguages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering languages: %v\n", err)
		os.Exit(1)
	}

	if len(languages) == 0 {
		fmt.Println("No language folders with benchmarks found.")
		os.Exit(0)
	}

	selectedLanguage, err := selectLanguage(languages)
	if err != nil {
		if err.Error() == "user quit" {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	prefs := newPreferences()

	for {
		benchmarks, err := discoverBenchmarks(selectedLanguage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error discovering benchmarks: %v\n", err)
			os.Exit(1)
		}

		if len(benchmarks) == 0 {
			fmt.Printf("No benchmarks found in %s.\n", selectedLanguage)
			os.Exit(0)
		}

		selectedBenchmark, err := selectBenchmark(selectedLanguage, benchmarks, prefs)
		if err != nil {
			if err.Error() == "user quit" {
				return
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runBenchmark(selectedLanguage, selectedBenchmark, prefs); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Benchmark failed: %v\n", err)
		}

		action, err := promptNextAction()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		switch action {
		case actionRunAnother:
			continue
		case actionChangeLang:
			selectedLanguage, err = selectLanguage(languages)
			if err != nil {
				if err.Error() == "user quit" {
					return
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		case actionExit:
			return
		}
	}
}
