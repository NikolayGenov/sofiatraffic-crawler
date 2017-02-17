package crawler

import (
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func loadGoodVTLine() *goquery.Document {
	return loadDocument("../testdata/vt_line_tram_1.html")
}

func TestVtLineCrawler_findRouteVTStops(t *testing.T) {
	vtRoutesStops := make([]VirtualTableStop, 0)
	routeSelection := loadGoodVTLine().Find("form").First()
	routes := Routes{
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
	}
	findRouteVTStops(routeSelection, &vtRoutesStops, routes)
	expected := []VirtualTableStop{
		{StopID: "2675", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "2658", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "2640", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "2606", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "2598", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1067", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1095", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1135", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1155", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1171", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "5353", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "4833", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1227", LineID: "27", RouteID: "1297", TransportationType: "0"},
		{StopID: "1241", LineID: "27", RouteID: "1297", TransportationType: "0"},
	}
	if !reflect.DeepEqual(vtRoutesStops, expected) {
		t.Errorf("Expecting %#v, got %#v", expected, vtRoutesStops)
	}
}
