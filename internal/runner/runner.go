package runner

import (
	"fmt"
	"github.com/niudaii/webscan/pkg/webscan"
	"github.com/projectdiscovery/gologger"
	"sort"
	"time"
)

type Runner struct {
	options       *Options
	webscanRunner *webscan.Runner
}

func NewRunner(options *Options) (*Runner, error) {
	webscanOptions := &webscan.Options{
		Proxy:       options.Proxy,
		Threads:     options.Threads,
		Timeout:     options.Timeout,
		Headers:     options.Headers,
		NoColor:     options.NoColor,
		FingerRules: options.FingerRules,
	}
	webscanRunner, err := webscan.NewRunner(webscanOptions)
	if err != nil {
		return nil, fmt.Errorf("webscan.NewRunner() err, %v", err)
	}
	return &Runner{
		options:       options,
		webscanRunner: webscanRunner,
	}, nil
}

func (r *Runner) Run() {
	start := time.Now()
	gologger.Info().Msgf("当前时间: %v", start.Format("2006-01-02 15:04:05"))
	// 目标解析
	if len(r.options.Targets) == 0 {
		gologger.Info().Msgf("目标为空")
		return
	}
	gologger.Info().Msgf("指纹数量: %v", len(r.options.FingerRules))
	gologger.Info().Msgf("目标数量: %v", len(r.options.Targets))
	// web扫描
	results := r.webscanRunner.Run(r.options.Targets)
	if len(results) == 0 {
		gologger.Info().Msgf("结果为空")
		return
	}
	// 排序并筛选重点指纹
	sort.Sort(results)
	var res string
	var fingerRes string
	var fingerNum int
	for _, result := range results {
		res += webscan.FmtResult(result, r.options.NoColor)
		// 显示重点指纹
		if len(result.Fingers) > 0 {
			// 筛选tags
			if result.Fingers = checkTags(result.Fingers); len(result.Fingers) == 0 {
				continue
			}
			fingerNum += 1
			fingerRes += webscan.FmtResult(result, r.options.NoColor)
		}
	}
	gologger.Info().Msgf("存活数量: %v", len(results))
	gologger.Print().Msgf("%v", res)
	gologger.Info().Msgf("重点指纹: %v", fingerNum)
	gologger.Print().Msgf("%v", fingerRes)
	gologger.Info().Msgf("运行时间: %v", time.Since(start))
}

func checkTags(fingers []*webscan.FingerRule) (newFingers []*webscan.FingerRule) {
	for _, finger := range fingers {
		flag := true
		for _, tag := range finger.Tags {
			if tag == "非重要" {
				flag = false
				break
			}
		}
		if flag {
			newFingers = append(newFingers, finger)
		}
	}
	return
}
