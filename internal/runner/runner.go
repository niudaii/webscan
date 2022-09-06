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
			// 过滤tags
			if result.Fingers = filterTags(result.Fingers, r.options.FilterTags); len(result.Fingers) > 0 {
				fingerNum += 1
				fingerRes += webscan.FmtResult(result, r.options.NoColor)
			}
		}
	}
	gologger.Info().Msgf("存活数量: %v", len(results))
	gologger.Print().Msgf("%v", res)
	gologger.Info().Msgf("重点指纹: %v", fingerNum)
	gologger.Print().Msgf("%v", fingerRes)
}

func filterTags(fingers []*webscan.FingerRule, filterTags []string) (newFingers []*webscan.FingerRule) {
	for _, finger := range fingers {
		flag := true
		for _, fingerTag := range finger.Tags {
			for _, tag := range filterTags {
				if fingerTag == tag {
					flag = false
					break
				}
			}
		}
		if flag {
			newFingers = append(newFingers, finger)
		}
	}
	return
}
