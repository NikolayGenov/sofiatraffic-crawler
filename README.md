# Sofia Traffic Crawler

The main goal of the crawler is to crawl only Sofia Traffic website and first use it as a tool to get the current infrastructure of the public transportation of Sofia by crawling all possible stations and lines. After that the tool should be use to detect and reflect changes in the infrastructure.
Another goal of the crawler is to poll on predefined period of time for the state of *each* vehicle which is included in the infrastructure


That poll can be every couple of seconds at best. The idea of the poll is to proxy it via the API with the value it gives and later on to use it for data learning purposes


The sources of information can be multiple.


An example can be this [Link](http://m.sofiatraffic.bg/schedules/vehicle?stop=1099&lid=24&vt=0&rid=873)
What can be extracted here is possibly the soft time schedules (but there can be seen at other places as well),
the list of consecutive stations with their names and Urban Mobility Center code (UMCC)


The data for the currently estimated time of arrival can be found by making a POST request to `http://m.sofiatraffic.bg/schedules/vehicle-vt` with example data `stop=1099&lid=24&vt=0&rid=873`

What is returned is some useful data like:

```
<div class="info">
		Информация към 17.01.2017 16:37	</div>
    <br>
	   " Точно време на пристигане: "
    <br>
    <b>16:50:42,16:43:19,16:59:07</b>
```
That data is constantly updating in some unknown period (for now) of time.
The main problem is that, the data is changing a lot during some hours of the day (mostly because of the traffic jams).
Yet another problem is that it's unknown when it's the exact arrival/leaving time of the vehicle without polling constantly and see that something that was there before has disappeared.


 
