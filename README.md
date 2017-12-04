Dockerized-Crawler-Service
==========================

Introduction
------------
In this project, a simple crawler that is providing two searching strategies `(eg. BFS, DFS)` boosting `go routines` was implemented from scratch,
and integrated into `Redis` and `Neo4j`. It is starting fetching from the first target URL and visit the rest of URL either `BFS` or `DFS`.
It is saving and keep track the visited URLs on `Redis` to provide consistency of data amongst a bunch of go routines.
It is either creating a new node and new link for each URL pair, or updating existing nodes and links  on `Neo4j` by using [`Cypher`](https://www.google.com) query, as well.

Configuration
-------------
The crawler is starting fetching the first `target_url` and according to `crawler_method` you defined in configuration file, under `crawler-go/src/data` directory.

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

The all services we are referring in this implementation are working as a docker container.
We wrote different dockerfile and *.yml files in order to manage, test and deploy these containers systematically.
Although we are using 3 different containers in building phase, we are using one more additional container in testing phase to launch a local `http-server`


    To test the services, we are using
        docker-compose.test.yml,
        src/Dockerfile.http-server

    To build the services, we are using
        docker-compose.yml,
        src/Dockerfile

In testing step, docker-compose file

Usage
-----
Please, note that you need to install compatible `docker` and `docker-compose` version before staring usage of the service.
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
You can find makefiles how to test and build the code under the root and src directories,
please take a look at the makefiles for further information. You can type make command on your shell, below, briefly

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

On Neo4j UI, you could write some specific [`Cypher`](https://www.google.com) query to visualize and analysis the graph corresponding to our data that the crawler service collected over the target URL.

When you type query Q1 below, it would return the all patterns related to our data collections.
    URL is class of an instance of URL visited by our crawler.
    link is corresponding to the relationship between two instances of URL.

    Q1:
        MATCH p=(url0:URL)-[r:link]->(url1:URL)
        RETURN p

When you type query Q2 below, you would see filtered result according to size of URL and link.
size is referring how many times that URL and link occur during the crawling.

    Q2:
        MATCH p=(url0:URL)-[r:link]->(url1:URL)
        WHERE url0.size > 5 and url1.size > 4 and r.size > 0
        RETURN p


![alt text](img/screenshot_neo4j%232.png "Logo Title Text 1")
![alt text](img/screenshot_neo4j#2.png?raw=true "Logo Title Text 1")

![adad]()

![adad](img/screenshot_neo4j%231.png)
