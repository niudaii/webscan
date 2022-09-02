package runner

import (
	"encoding/json"
	"fmt"
	"github.com/niudaii/webscan/internal/utils"
	"github.com/niudaii/webscan/pkg/webscan"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"strings"
	"time"
)

type Options struct {
	// input
	Input     string
	InputFile string
	// config
	Threads    int
	Timeout    int
	Proxy      string
	Header     string
	FingerFile string
	// output
	OutputFile string
	NoColor    bool
	// debug
	Silent bool
	Debug  bool

	Targets     []string
	Headers     []string
	FingerRules []*webscan.FingerRule `json:"-"`
}

func ParseOptions() *Options {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Webscanner`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Input, "input", "i", "", "url input(example: -i 'http://www.baidu.com', -i '192.168.243.11:81')"),
		flagSet.StringVarP(&options.InputFile, "input-file", "f", "", "urls file(example: -f 'xxx.txt')"),
	)

	flagSet.CreateGroup("config", "Config",
		flagSet.IntVar(&options.Threads, "threads", 1, "number of threads"),
		flagSet.IntVar(&options.Timeout, "timeout", 10, "timeout in seconds"),
		flagSet.StringVarP(&options.Proxy, "proxy", "p", "", "proxy(example: -p 'http://127.0.0.1:8080')"),
		flagSet.StringVar(&options.Header, "header", "", "add custom headers(example: -header 'User-Agent: xxx, ')"),
		flagSet.StringVar(&options.FingerFile, "finger-file", "", "use your finger file(example: -finger-file 'fingers.json')"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.OutputFile, "output", "o", "webscan.txt", "output file to write found results"),
		flagSet.BoolVarP(&options.NoColor, "no-color", "nc", false, "disable colors in output"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&options.Silent, "silent", false, "show only results in output"),
		flagSet.BoolVar(&options.Debug, "debug", false, "show debug output"),
	)

	if err := flagSet.Parse(); err != nil {
		gologger.Fatal().Msgf("Program exiting: %v", err)
	}

	options.configureOutput()

	showBanner()

	if err := options.validateOptions(); err != nil {
		gologger.Fatal().Msgf("Program exiting: %v", err)
	}

	if err := options.configureOptions(); err != nil {
		gologger.Fatal().Msgf("Program exiting: %v", err)
	}

	return options
}

func (o *Options) configureOutput() {
	if o.NoColor {
		gologger.DefaultLogger.SetFormatter(formatter.NewCLI(true))
	}

	if o.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}

	if o.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}

	gologger.DefaultLogger.SetWriter(utils.NewCLI(o.OutputFile))
}

func (o *Options) validateOptions() error {
	if o.Input == "" && o.InputFile == "" {
		return fmt.Errorf("No service input provided")
	}
	if o.Debug && o.Silent {
		return fmt.Errorf("Both debug and silent mode specified")
	}
	if o.FingerFile != "" && !utils.FileExists(o.FingerFile) {
		return fmt.Errorf("File %v does not exist", o.FingerFile)
	}

	return nil
}

func (o *Options) configureOptions() error {
	if o.Input != "" {
		o.Targets = append(o.Targets, o.Input)
	} else {
		lines, err := utils.ReadLines(o.InputFile)
		if err != nil {
			return err
		}
		o.Targets = append(o.Targets, lines...)
	}

	if o.Header != "" {
		o.Headers = strings.Split(o.Header, ",")
	}

	if o.Proxy == "bp" {
		o.Proxy = "http://127.0.0.1:8080"
	}

	// 读取指纹
	if o.FingerFile != "" {
		bytes, err := utils.ReadFile(o.FingerFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &o.FingerRules)
		if err != nil {
			return err
		}
	}

	o.Targets = utils.RemoveDuplicate(o.Targets)

	gologger.Info().Msgf("当前时间: %v", time.Now().Format("2006-01-02 15:04:05"))
	opt, _ := json.Marshal(o)
	gologger.Debug().Msgf("当前配置: %v", string(opt))
	gologger.Info().Msgf("指纹数量: %v", len(o.FingerRules))

	return nil
}
