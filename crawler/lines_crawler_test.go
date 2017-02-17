package crawler

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func loadBadLine() *goquery.Document {
	return loadDocument("../testdata/bad_line.html")
}

func loadGoodLine() *goquery.Document {
	return loadDocument("../testdata/line_tram_1.html")
}

func loadDocument(filename string) *goquery.Document {
	r, _ := os.Open(filename)
	rd := bufio.NewReader(r)
	doc, _ := goquery.NewDocumentFromReader(rd)
	return doc
}

func TestLinesCrawler_getOperationsMap(t *testing.T) {
	doc := loadGoodLine()
	fmt.Print()
	result := getOperationsMap(doc)
	expected := OperationIDMap{Normal: "6671", Holiday: "6672"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestLinesCrawler_getOperationsMapMissingIdentifiers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic because there are missing schedule_active_list_tabs identifers")
		}
	}()
	doc := loadBadLine()
	getOperationsMap(doc)
}

func TestLinesCrawler_getOperationIDRoutesMap(t *testing.T) {
	doc := loadGoodLine()
	fmt.Print()
	result := getOperationIDRoutesMap(doc)
	expected := OperationIDRoutesMap{
		"6671": Routes{
			Route{
				Direction{
					Name: "Надлез Надежда - ж.к. Иван Вазов",
					ID:   "2696",
				},
				Stops{
					Stop{CapitalName: "НАДЛЕЗ НАДЕЖДА", Sign: "1113", ID: "959"},
					Stop{CapitalName: "ПЕТА ГРАДСКА БОЛНИЦА", Sign: "1254", ID: "565"},
					Stop{CapitalName: "ЦЕНТРАЛНА ГАРА", Sign: "1331", ID: "586"},
					Stop{CapitalName: "УЛ. КЛОКОТНИЦА", Sign: "1993", ID: "904"},
					Stop{CapitalName: "БУЛ. СЛИВНИЦА", Sign: "0377", ID: "175"},
					Stop{CapitalName: "УЛ. ПИРОТСКА", Sign: "2112", ID: "958"},
					Stop{CapitalName: "ПЛ. МАКЕДОНИЯ", Sign: "1282", ID: "575"},
					Stop{CapitalName: "БУЛ. ПРАГА", Sign: "0364", ID: "171"},
					Stop{CapitalName: "НДК", Sign: "1134", ID: "515"},
					Stop{CapitalName: "БУЛ. ПЕНЧО СЛАВЕЙКОВ", Sign: "0356", ID: "168"},
					Stop{CapitalName: "14-ТИ ДКЦ", Sign: "0011", ID: "5"},
					Stop{CapitalName: "ЧИТАЛИЩЕ Д-Р П. БЕРОН", Sign: "0929", ID: "417"},
					Stop{CapitalName: "Ж.К. ИВАН ВАЗОВ", Sign: "0625", ID: "286"}},
			},
			Route{
				Direction{
					Name: "Ж.к. Иван Вазов - надлез Надежда",
					ID:   "2697",
				},
				Stops{
					Stop{CapitalName: "Ж.К. ИВАН ВАЗОВ", Sign: "0625", ID: "286"},
					Stop{CapitalName: "ЧИТАЛИЩЕ Д-Р П. БЕРОН", Sign: "0928", ID: "417"},
					Stop{CapitalName: "14-ТИ ДКЦ", Sign: "0010", ID: "5"},
					Stop{CapitalName: "БУЛ. ПЕНЧО СЛАВЕЙКОВ", Sign: "0357", ID: "168"},
					Stop{CapitalName: "НДК", Sign: "1135", ID: "515"},
					Stop{CapitalName: "БУЛ. ПРАГА", Sign: "0365", ID: "171"},
					Stop{CapitalName: "ПЛ. МАКЕДОНИЯ", Sign: "1281", ID: "575"},
					Stop{CapitalName: "УЛ. ПИРОТСКА", Sign: "2113", ID: "958"},
					Stop{CapitalName: "БУЛ. СЛИВНИЦА", Sign: "0376", ID: "175"},
					Stop{CapitalName: "УЛ. КЛОКОТНИЦА", Sign: "1994", ID: "904"},
					Stop{CapitalName: "ЦЕНТРАЛНА АВТОГАРА", Sign: "2665", ID: "1275"},
					Stop{CapitalName: "ЦЕНТРАЛНА ГАРА", Sign: "1332", ID: "586"},
					Stop{CapitalName: "ПЕТА ГРАДСКА БОЛНИЦА", Sign: "1255", ID: "565"},
					Stop{CapitalName: "НАДЛЕЗ НАДЕЖДА", Sign: "1113", ID: "959"}}},
		},
		"6672": Routes{Route{
			Direction{
				Name: "Надлез Надежда - ж.к. Иван Вазов",
				ID:   "2696",
			},
			Stops{Stop{CapitalName: "НАДЛЕЗ НАДЕЖДА", Sign: "1113", ID: "959"},
				Stop{CapitalName: "ПЕТА ГРАДСКА БОЛНИЦА", Sign: "1254", ID: "565"},
				Stop{CapitalName: "ЦЕНТРАЛНА ГАРА", Sign: "1331", ID: "586"},
				Stop{CapitalName: "УЛ. КЛОКОТНИЦА", Sign: "1993", ID: "904"},
				Stop{CapitalName: "БУЛ. СЛИВНИЦА", Sign: "0377", ID: "175"},
				Stop{CapitalName: "УЛ. ПИРОТСКА", Sign: "2112", ID: "958"},
				Stop{CapitalName: "ПЛ. МАКЕДОНИЯ", Sign: "1282", ID: "575"},
				Stop{CapitalName: "БУЛ. ПРАГА", Sign: "0364", ID: "171"},
				Stop{CapitalName: "НДК", Sign: "1134", ID: "515"},
				Stop{CapitalName: "БУЛ. ПЕНЧО СЛАВЕЙКОВ", Sign: "0356", ID: "168"},
				Stop{CapitalName: "14-ТИ ДКЦ", Sign: "0011", ID: "5"},
				Stop{CapitalName: "ЧИТАЛИЩЕ Д-Р П. БЕРОН", Sign: "0929", ID: "417"},
				Stop{CapitalName: "Ж.К. ИВАН ВАЗОВ", Sign: "0625", ID: "286"}},
		},
			Route{
				Direction{
					Name: "Ж.к. Иван Вазов - надлез Надежда",
					ID:   "2697",
				},
				Stops{
					Stop{CapitalName: "Ж.К. ИВАН ВАЗОВ", Sign: "0625", ID: "286"},
					Stop{CapitalName: "ЧИТАЛИЩЕ Д-Р П. БЕРОН", Sign: "0928", ID: "417"},
					Stop{CapitalName: "14-ТИ ДКЦ", Sign: "0010", ID: "5"},
					Stop{CapitalName: "БУЛ. ПЕНЧО СЛАВЕЙКОВ", Sign: "0357", ID: "168"},
					Stop{CapitalName: "НДК", Sign: "1135", ID: "515"},
					Stop{CapitalName: "БУЛ. ПРАГА", Sign: "0365", ID: "171"},
					Stop{CapitalName: "ПЛ. МАКЕДОНИЯ", Sign: "1281", ID: "575"},
					Stop{CapitalName: "УЛ. ПИРОТСКА", Sign: "2113", ID: "958"},
					Stop{CapitalName: "БУЛ. СЛИВНИЦА", Sign: "0376", ID: "175"},
					Stop{CapitalName: "УЛ. КЛОКОТНИЦА", Sign: "1994", ID: "904"},
					Stop{CapitalName: "ЦЕНТРАЛНА АВТОГАРА", Sign: "2665", ID: "1275"},
					Stop{CapitalName: "ЦЕНТРАЛНА ГАРА", Sign: "1332", ID: "586"},
					Stop{CapitalName: "ПЕТА ГРАДСКА БОЛНИЦА", Sign: "1255", ID: "565"},
					Stop{CapitalName: "НАДЛЕЗ НАДЕЖДА", Sign: "1113", ID: "959"}},
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expecting %#v, got %#v", expected, result)
	}
}

func TestLinesCrawler_getOperationIDRoutesMapMissingDirectionInformation(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic because there are missing schedule_view_direction_tabs text and/or href")
		}
	}()
	doc := loadBadLine()
	getOperationIDRoutesMap(doc)
}
