package crawler

import "fmt"

type LineBasicInfo struct {
	Name string
	URL  string
	Transportation
}

func (l LineBasicInfo) String() string {
	return fmt.Sprintf("{%v '%v'}", l.Transportation, l.Name)
}

type LinesBasicInfo []LineBasicInfo
