# Sofia Traffic Crawler

## Goal

The main goal of the crawler is to crawl only Sofia Traffic websites (http://schedules.sofiatraffic.bg and http://m.sofiatraffic.bg) and first use it as a tool to get the current infrastructure of the public transportation of Sofia by crawling all possible stations, lines, schedules, etc.
It also tries to match the data between listed above sites, which differs *A LOT*.
In future the tool should be able to detect the frequent changes to the structure.
Another goal of the crawler is to poll for the times of *each* active stop which is included in the infrastructure on given operation mode.

That poll can be every couple of seconds at best. The idea of the poll is to proxy it via the API with the value it gives and later on to use it for data learning purposes

## Sources and types of extracted data
The sources for the data are different:
For predefined lines, directions, operation modes, stops and schedules we can use [schedules.sofiatraffic.bg](http://schedules.sofiatraffic.bg/) and some of it 'hidden' services.

After we have some reference data about the structure of the traffic we can match it
from what we can scrape from [m.sofiatraffic.bg](http://m.sofiatraffic.bg/).
I presume that both sites use different databases because of the significant
differences on name where they should be the name - e.g direction names can be horribly different.
Also they use different IDs for the same things and structure of information as a whole.
Nevertheless there is some really useful (not really consistent though) information that can be found here.
Namely - the precious almost real-time-ish data about the time of next arrival of the vehicle on the
stop we are querying for. To get that data we need to make a `POST` request to `http://m.sofiatraffic.bg/schedules/vehicle-vt` with example data `stop=1099&lid=24&vt=0&rid=873`
And after parsing some HTML can get this: `16:50:42,16:43:19,16:59:07`

## Some noted problems
Those times can disappear as if the vehicle has arrived and then reappear.
Another problem is that they don't always disappear after that time or the vehicle has arrived
and can stay for observation more than 15 minutes (possibly more) after the time has passed.

Problem is silent addition to lines with schedules which serve only 1 day purpose.

Another problem is that some lines have 8 different routes that they serve.

Some lines do not operate during the weekend, others during holidays, etc.
Some lines are marked for removal or update of some of the stops, but the notes are
in Bulgarian and are practically un-parsable without NPL.

## Usage
In order to use the crawler - you need redis and you have to pass a redis connection pool to the crawler.

After that you can start 4 types of crawls, which have some dependencies between them, but are cache-able
in redis, so you don't have to run them each time.

* `CrawlLines` will extract all the lines, directions, operation modes, stops information with some internal IDs
* `CrawlSchedules` will use the information from the the previous call and get all the schedules for
each stop in the traffic network for every line, operation mode, and direction.
* `CrawlVirtualTablesLines` will use the information from `CrawlLines` and will match everything it can
from [m.sofiatraffic.bg](http://m.sofiatraffic.bg/) in order to extract the IDs for stops there
* `CrawlVirtualTablesStopsForTimes` will use the information from `CrawlVirtualTablesLines` and will
start parallel query for all the active stops in order to get the desired times.
