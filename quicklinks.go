package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/quicklinks/config"
)

func main() {
	args := os.Args[1:]

	allConfigs := loadConfigs(args)
	runConfigs := chooseConfigs(allConfigs, args)
	run(runConfigs)
}

// usage prints an helper message
func usage() {
	println("quicklinks file [choices...]")
	println("The documentation is available at:")
	println("https: //github.com/gandrille/quicklinks")
}

func loadConfigs(args []string) []config.Configuration {

	if len(args) == 0 {
		result.PrintRed("Missing parameters")
		usage()
		os.Exit(1)
	}

	if args[0] == "--help" {
		usage()
		os.Exit(0)
	}

	configs, err := config.New(args[0])
	if err != nil {
		result.PrintError(err.Error())
		os.Exit(1)
	}

	return configs
}

func chooseConfigs(configs []config.Configuration, args []string) []config.Configuration {
	keys := getKeys(configs, args)

	if len(keys) == 1 && keys[0] == "all" {
		return configs
	}

	var selected []config.Configuration
	for _, key := range keys {
		selected = append(selected, configByKey(configs, key))
	}

	return selected
}

func getKeys(configs []config.Configuration, args []string) []string {
	if len(args) >= 2 {
		return args[1:]
	} else {
		printConfigs(configs)
		print("Please enter at least one number or key: ")
		input := readLine()
		tokens := splitLine(input)
		return tokens
	}
}

func printConfigs(configs []config.Configuration) {
	for i, conf := range configs {
		fmt.Printf("%2d %s\n", i+1, conf.Key())
	}
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		result.PrintError("Oups! Reading error...")
		os.Exit(0)
	}
	return scanner.Text()
}

func splitLine(line string) []string {
	var tokens []string
	for _, token := range strings.Split(line, " ") {
		t := strings.TrimSpace(token)
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func configByKey(configs []config.Configuration, key string) config.Configuration {
	for i, conf := range configs {
		if conf.Key() == key {
			return conf
		}
		if strconv.Itoa(i+1) == key {
			return conf
		}
	}
	result.PrintError("Can't find line with key " + key)
	os.Exit(0)
	return configs[0] // to make the linter happy after Exit
}

func run(runConfigs []config.Configuration) {

	if len(runConfigs) == 0 {
		result.PrintError("Nothing choosen, bye!")
		os.Exit(0)
		return // to make the linter happy after Exit
	}

	if len(runConfigs) == 1 {
		runSingle(runConfigs[0])
		return
	}

	var results []result.Result
	for _, conf := range runConfigs {
		results = append(results, runSingle(conf))
		fmt.Println("")
	}

	result.PrintInfo("Summary")
	result.NewSet(results, "" /* default message */).Print()
}

func runSingle(conf config.Configuration) result.Result {
	result.PrintInfo(conf.Key())
	runner := func() result.Result { return conf.Run() }
	res := result.Run(runner)
	res.Print()
	return res
}
