package webscan

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

type Result struct {
	Url           string
	StatusCode    int
	ContentLength int
	Favicon       string
	IconHash      string
	Title         string
	Fingers       []*FingerRule
	Wappalyzer    map[string]struct{}
}

// Results 按照contentLength排序
type Results []*Result

func (s Results) Len() int {
	return len(s)
}
func (s Results) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Results) Less(i, j int) bool {
	return s[i].ContentLength < s[j].ContentLength
}

func FmtResult(result *Result, noColor bool) (res string) {
	if noColor {
		res = fmt.Sprintf("%v [%v] [%v] [%v] [%v] [%v] [%v]\n", result.Url, result.StatusCode, result.ContentLength, result.IconHash, result.Title, GetWappalyzerString(result.Wappalyzer), GetFingerString(result.Fingers))
	} else {
		res = fmt.Sprintf("%v [%v] [%v] [%v] [%v] [%v] [%v]\n", result.Url, aurora.Red(result.StatusCode), aurora.Blue(result.ContentLength), aurora.Yellow(result.IconHash), aurora.Green(result.Title), GetWappalyzerString(result.Wappalyzer), aurora.Red(GetFingerString(result.Fingers)))
	}
	return
}

func GetFingerString(fingers []*FingerRule) (fingerString string) {
	for _, finger := range fingers {
		fingerString += fmt.Sprintf("(%v", finger.Name)
		if len(finger.Tags) != 0 {
			fingerString += fmt.Sprintf(" | %v", finger.Tags)
		}
		if len(finger.PocTags) != 0 {
			fingerString += fmt.Sprintf(" | %v", finger.PocTags)
		}
		if finger.Desc != "" {
			fingerString += fmt.Sprintf(" | %v", finger.Desc)
		}
		fingerString += ")"
	}
	return
}

func GetWappalyzerString(result map[string]struct{}) (wappalyzerString string) {
	for k := range result {
		wappalyzerString += fmt.Sprintf("(%v)", k)
	}
	return
}
