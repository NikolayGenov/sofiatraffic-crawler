package crawler

import "strings"

type ScheduleID string

type ScheduleTimes []string

type Schedules map[ScheduleID]ScheduleTimes

func convertToScheduleID(path string) ScheduleID {
	//The path is in this format /server/html/schedule_load/{OperationID}/{DirectionID}/{StopSign}
	//ScheduleID is in this format {OperationID}/{DirectionID}/{StopSign}
	parts := strings.Split(path, "/")
	return ScheduleID(strings.Join(parts[len(parts)-3:], "/"))
}

//func getNormalTimesOfTimes(doc *goquery.Document) [][]string {
//	timesOfTimes := make([][]string, 0)
//	doc.Find(SCHEDULE_TIMES_SELECTOR).Each(func(i1 int, st *goquery.Selection) {
//		times := make([]string, 0)
//		st.Find(SCHEDULE_LINKS_TIMES_SELECTOR).
//			Each(func(i int, s *goquery.Selection) {
//				times = append(times, strings.TrimSpace(s.Text()))
//			})
//		timesOfTimes = append(timesOfTimes, times)
//	})
//	return timesOfTimes
//}

//func advancedTimes(doc *goquery.Document) [][][]string {
//	timesOfTimes := make([][][]string, 0)
//	doc.Find(SCHEDULE_TIMES_SELECTOR).Each(func(i1 int, st *goquery.Selection) {
//		times := make([][]string, 0)
//		st.Find(SCHEDULE_LINKS_TIMES_SELECTOR).Each(func(i int, s *goquery.Selection) {
//			click, _ := s.Attr("onclick")
//			//fmt.Println(click)
//			i2 := strings.LastIndex(click, "'")
//			i1 := strings.LastIndex(click[:i2], "'")
//			reduced := click[i1+1 : i2]
//			splits := strings.Split(reduced, ",")
//			times = append(times, splits)
//		})
//		timesOfTimes = append(timesOfTimes, times)
//	})
//	return timesOfTimes
//}
//func intToTime(c string) string {
//	i, _ := strconv.Atoi(c)
//	return fmt.Sprintf("%v:%02d", i/60, i%60)
//}
//
//func printTimes(times [][]string) {
//	l := len(times[0])
//	for i := 1; i < l; i++ {
//		for _, row := range times {
//
//			time := row[i]
//
//			if time != "" {
//				fmt.Printf("%v\t", intToTime(time))
//			} else {
//				fmt.Print("*****\t")
//			}
//		}
//		fmt.Print("\n")
//	}
//}
