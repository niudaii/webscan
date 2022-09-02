package webscan

import (
	"github.com/imroc/req/v3"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/wappalyzergo"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	reqClient        *req.Client
	wappalyzerClient *wappalyzer.Wappalyze
	threads          int
	noColor          bool
	fingerRules      []*FingerRule
}

func NewEngine(proxy string, timeout, threads int, headers []string, noColor bool, fingerRules []*FingerRule) (*Engine, error) {
	wappalyzerClient, err := wappalyzer.New()
	if err != nil {
		return nil, err
	}
	return &Engine{
		reqClient:        NewReqClient(proxy, timeout, headers),
		wappalyzerClient: wappalyzerClient,
		threads:          threads,
		noColor:          noColor,
		fingerRules:      fingerRules,
	}, nil
}

func NewReqClient(proxy string, timeout int, headers []string) *req.Client {
	reqClient := req.C()
	reqClient.GetTLSClientConfig().InsecureSkipVerify = true
	reqClient.SetCommonHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Cookie":     "rememberMe=1", // check shiro
	})
	reqClient.SetRedirectPolicy(req.AlwaysCopyHeaderRedirectPolicy("Cookie"))
	if proxy != "" {
		reqClient.SetProxyURL(proxy)
	}
	var key, value string
	for _, header := range headers {
		tokens := strings.SplitN(header, ":", 2)
		if len(tokens) < 2 {
			continue
		}
		key = strings.TrimSpace(tokens[0])
		value = strings.TrimSpace(tokens[1])
		reqClient.SetCommonHeader(key, value)
	}
	reqClient.SetTimeout(time.Duration(timeout) * time.Second)
	return reqClient
}

func (e *Engine) Run(urls []string) (results Results) {
	// RunTask
	wg := &sync.WaitGroup{}
	taskChan := make(chan string, e.threads)
	for i := 0; i < e.threads; i++ {
		go func() {
			for task := range taskChan {
				resp, err := e.Webinfo(task)
				if err != nil {
					gologger.Debug().Msgf("%v", err)
				} else {
					// 判断蜜罐匹配大量指纹的情况
					if len(resp.Fingers) > 5 {
						gologger.Warning().Msgf("%v 可能为蜜罐", resp.Url)
					} else {
						gologger.Silent().Msgf(FmtResult(resp, e.noColor))
						results = append(results, resp)
					}
				}
				wg.Done()
			}
		}()
	}

	for _, task := range urls {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()

	gologger.Info().Msgf("扫描结束")

	return
}

func (e *Engine) Webinfo(url string) (result *Result, err error) {
	r := e.reqClient.R()
	resp, err := FirstGet(r, url)
	if err != nil {
		return
	}
	// 处理js跳转, 上限5次
	for i := 0; i < 5; i++ {
		jumpurl := Jsjump(resp)
		if jumpurl == "" {
			break
		}
		resp, err = r.Get(jumpurl)
	}
	if err != nil {
		return
	}
	result = &Result{
		Url:           resp.Request.URL.String(),
		StatusCode:    resp.StatusCode,
		ContentLength: len(resp.String()),
		Title:         getTitle(resp),
		Fingers:       e.getFinger(resp),
	}
	result.Favicon, result.IconHash = e.getFavicon(resp)
	result.Wappalyzer = e.wappalyzerClient.Fingerprint(resp.Header, resp.Bytes())
	return
}

func FirstGet(r *req.Request, url string) (resp *req.Response, err error) {
	var scheme string
	if !strings.HasPrefix(url, "http") {
		scheme = "http://"
		resp, err = r.Get(scheme + url)
		if err != nil {
			scheme = "https://"
		} else {
			if strings.Contains(resp.String(), "sent to HTTPS port") || strings.Contains(resp.String(), "This combination of host and port requires TLS") || strings.Contains(resp.String(), "Instead use the HTTPS scheme to") {
				scheme = "https://"
			}
		}
	}
	resp, err = r.Get(scheme + url)
	return
}
