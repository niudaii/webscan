package runner

import (
	"fmt"
	"github.com/niudaii/webscan/pkg/webscan"
	"github.com/projectdiscovery/gologger"
	"sort"
)

type Runner struct {
	options       *Options
	webscanRunner *webscan.Runner
}

func NewRunner(options *Options) (*Runner, error) {
	webscanRunner, err := webscan.NewRunner(options.Proxy, options.Threads, options.Timeout, options.Headers, options.NoColor, options.FingerRules)
	if err != nil {
		return nil, fmt.Errorf("NewEngine err, %v", err)
	}
	runner := &Runner{
		options:       options,
		webscanRunner: webscanRunner,
	}
	return runner, nil
}

func (r *Runner) Run() {
	if len(r.options.Targets) == 0 {
		gologger.Info().Msgf("目标为空")
		return
	}
	gologger.Info().Msgf("指纹数量: %v", len(r.options.FingerRules))
	gologger.Info().Msgf("目标数量: %v", len(r.options.Targets))
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
			// 非重要
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
