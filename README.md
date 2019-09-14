Crawler as a Service
==========================

Introduction
------------
A simple crawler service was implemented from scratch, and integrated into `Redis` and `Neo4j` NoSQL systems by using `Docker` and `Docker-compose`.
The crawler service is crawling the first target URL, and then, visiting the rest of URLs in the fetched HTML documents, respectively and recursively.
While crawling a HTML documents corresponding to URLs, it could refer to 1 out of 2 different searching algorithms (`BFS, DFS`).
Those searching algorithms were boosted by `go routines` in `GO` in order to speed up crawling service.

During crawling, there is a possibility that a bunch of go routines that would be created may fetch and process the same HTML documents at the same time.
In this case, the crawler may create inconsistent data. Thus, `Redis` Key-Value NoSQL system was preferred using in this project to solve that problem and build a robust and consistent system.

Each URL may referring to either the other different URL or itself in a HTML document. That relationship between two URLs can call as a Link.
There is a simple easy way to represent those crawled Links and URLs by using a specific data structure, which is graph.
Thus, `Neo4j` Graph NoSQL were used to represent and visualize the graph which consists of URLs and Links.
During crawling, the crawling service is either creating a new node for each URL and new link for each URL pair, or updating existing nodes and links on `Neo4j` by using [`Cypher`](https://neo4j.com/developer/cypher-query-language/) query, as well.

Configuration
-------------
The crawler service is starting crawling the first `target_url` and according to `crawler_method` you defined in configuration file, under `crawler-go/src/data` directory.

    config.json

You could define different parameters related to the service, in this configuration file, below.

```javascript
{
  "app_name" : "crawler",
  "redis_address":  "redis:6379",
  "neo4j_bolt_address": "bolt://neo4j:7687",
  "target_url": "http://tomblomfield.com",
  "filter_url": "http://tomblomfield.com",
  "reset_redis" : true,
  "reset_neo4j" : true,
  "redis_connection_attempts": 15,
  "neo4j_connection_attempts" : 15,
  "attempts_time_interval": "1000ms",
  "crawler_method" : "dfs"
}
```

Integration of Services
-----------------------
In integration of services, `docker` and `docker-compose` are being used.

All services we are using in this project are working as a docker container.
We wrote different dockerfile and *.yml files in order to manage, test and deploy these containers systematically.
Although we are using 3 different containers in building phase, we are using one more additional container corresponding to `http-server` in testing phase to run test cases locally on this http-server.

    To test the services, we are using
        docker-compose.test.yml,
        src/Dockerfile.http-server

    To build the services, we are using
        docker-compose.yml,
        src/Dockerfile


Usage
-----
Please, note that you need to install compatible `docker` and `docker-compose` version before using the service.
You could see the version of `docker` and `docker-compose`, below.

##### docker version

    Client:
     Version:      17.09.0-ce
     API version:  1.32
     Go version:   go1.8.3
     Git commit:   afdb6d4
     Built:        Tue Sep 26 22:42:18 2017
     OS/Arch:      linux/amd64

    Server:
     Version:      17.09.0-ce
     API version:  1.32 (minimum version 1.12)
     Go version:   go1.8.3
     Git commit:   afdb6d4
     Built:        Tue Sep 26 22:40:56 2017
     OS/Arch:      linux/amd64
     Experimental: false

##### docker-compose version

    docker-compose version 1.17.0, build ac53b73
    docker-py version: 2.5.1
    CPython version: 2.7.13
    OpenSSL version: OpenSSL 1.0.1t  3 May 2016

##### Running Test Cases and Building Services
You can find makefiles how to test and build the code under the directories, root and src.
please take a look at the makefiles for further information. You can type make commands on your command prompt, below

To run test cases:

        make run-test

To stop test cases:

        make stop-service

To run service:

        make run-service

To stop service:

        make stop-service

Visualization
-------------

After run the services using the make command above, you could type the following address on your favorite browser to open `Neo4j` Console .

    http://localhost:7474/browser/

On Neo4j UI, you could write some specific [`Cypher`](https://neo4j.com/developer/cypher-query-language/) query to visualize and analyze the graph corresponding to our data that the crawler service collected over the target URL.

When you type query Q1 below, it would return the all patterns related to our data collections.
    URL is class of an instance of URL visited by our crawler.
    link is corresponding to the relationship between two instances of URL.

Q1:

        MATCH p=(url0:URL)-[r:link]->(url1:URL)
        RETURN p

When you type query Q2 below, you would see filtered result according to size of URL and Link.
size is referring how many times that URL and link occur during the crawling.

Q2:

        MATCH p=(url0:URL)-[r:link]->(url1:URL)
        WHERE url0.size > 5 and url1.size > 4 and r.size > 0
        RETURN p

You can see graph of crawled data over the target URL, below.

![Screenshot of Graph-1, http://tomblomfield.com](img/screenshot_neo4j%231.png)

![Screenshot of Graph-2, http://tomblomfield.com/rss](img/screenshot_neo4j%232.png)
