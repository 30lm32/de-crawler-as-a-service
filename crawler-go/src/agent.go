package main

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"time"
	"fmt"
	"log"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
)

type Agent struct {
	mu  sync.Mutex
	config *configuration
	conn   bolt.Conn
	client *redis.Client
	containerLinkName string
	containerUrlName string
}

func initializeNeo4jConnection(config *configuration, ch chan message) {
	driver := bolt.NewDriver()

	attempts := int64(0)
	attemptsTimeInterval, err := time.ParseDuration(config.Attempts_Time_Interval)
	if err != nil {
		ch <- message{nil, err}
		return
	}

	for {
		if attempts >= config.Neo4j_Connection_Attempts {
			ch <- message{nil, fmt.Errorf("Neo4j connection could not be established...\n")}
			return
		} else {

			//conn, err :=  driver.OpenNeo("bolt://neo4j:1234@localhost:7687")
			conn, err :=  driver.OpenNeo(config.Neo4j_Bolt_Address)

			if err != nil {
				log.Printf("Connection to Neo4j: Attempt#%d\n", attempts)
				attempts++
			} else {
				log.Printf("Neo4j connection is ready...\n")
				ch <- message{conn, nil}
				return
			}
		}
		time.Sleep(attemptsTimeInterval)
	}

}

func initializeRedisConnection(config *configuration, ch chan message) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis_Address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	attempts := int64(0)
	attemptsTimeInterval, err := time.ParseDuration(config.Attempts_Time_Interval)
	if err != nil {
		ch <- message{nil, err}
		return
	}

	for {
		if attempts >= config.Redis_Connection_Attempts {
			ch <- message{ nil, fmt.Errorf("Redis connection could not be established...\n")}
			return
		} else {
			pong, err := client.Ping().Result()
			if err != nil {
				log.Printf("Connection to Redis: Attempt#%d\n", attempts)
				attempts++
			} else {
				log.Printf("Redis connection is ready... <%s>\n", pong)
				ch <- message{client, nil}
				return
			}
			time.Sleep(attemptsTimeInterval)
		}
	}
}

func NewAgent(config *configuration) (*Agent) {

	agent := &Agent{
		config : config,
		client : nil,
		conn : nil,
		containerLinkName : fmt.Sprintf(keyLink, config.App_Name),
		containerUrlName : fmt.Sprintf(keyUrl, config.App_Name),
	}

	chNeo4j := make(chan message)
	chRedis := make(chan message)

	go initializeNeo4jConnection(config, chNeo4j)
	go initializeRedisConnection(config, chRedis)

	var err error

	// sync all services
	var wg sync.WaitGroup
	wg.Add(2)

	go func () {
		for i := 0; i < 2; i++{
			select {
			case t :=<-chNeo4j:
				err = t.err
				if err != nil{
					log.Printf(err.Error())
				} else {
					agent.conn = t.holder.(bolt.Conn)
				}
				wg.Done()

			case t :=<-chRedis:
				err = t.err
				if err != nil{
					log.Printf(err.Error())
				} else {
					agent.client = t.holder.(*redis.Client)
				}
				wg.Done()

			case <-time.After(1 * time.Minute):
				wg.Done()
			}
		}
	} ()

	//waiting all services until finalizing
	wg.Wait()

	close(chNeo4j)
	close(chRedis)

	//when occur any error related to any service, agent will be nil
	if err != nil {
		return nil
	}

	return agent
}


func (agent Agent) dumbData(key string)  {
	r, err := agent.client.HGetAll(key).Result()
	if err == nil {
		var tmp int64 = 0
		for k, v := range r{
			nv, _ := strconv.ParseInt(v, 10, 64)
			tmp += nv
			log.Printf("%s : %d\n", k, nv)
		}
		log.Printf("Total links we visited: %d", tmp)
	}
}

func (agent Agent) getNumberOfUrls(url string) int64{
	v, err := agent.client.Get(url).Int64()
	if err != nil {
		return int64(0)
	}
	return v
}

func (agent Agent) getNumberOfLinks(lnk link) int64 {
	key := fmt.Sprintf(formatLink, lnk.PUrl, lnk.SUrl)

	v , err := agent.client.HGet(agent.containerLinkName, key).Int64()
	if err != nil {
		return int64(0)
	}
	return v
}


func (agent Agent) getNumberOfUniqueUrls() int {
	r, err := agent.client.HGetAll(agent.containerUrlName).Result()

	if err == nil {
		return len(r)
	}
	return 0
}

func (agent Agent) getNumberOfUniqueLinks() int {
	r, err := agent.client.HGetAll(agent.containerLinkName).Result()

	if err == nil {
		return len(r)
	}
	return 0
}


func (agent Agent) executeRelQuery(query string, lnk link) int64{

	agent.mu.Lock()

	v := map[string]interface{}{
		"name0": lnk.PUrl,
		"name1": lnk.SUrl,
		"size": agent.getNumberOfLinks(lnk),
	}

	//log.Printf("%s %d\n", query, v)
	result, _ := agent.conn.ExecNeo(query,v)
	numResult, _ := result.RowsAffected()

	agent.mu.Unlock()

	return numResult
}

func (agent Agent) insertOrUpdateRel(lnk link) {

	q := "MATCH p=(url0:URL {name: {name0}})-[l:link]->(url1:URL {name: {name1}})" +
		"SET l.size = {size}"

	numResult := agent.executeRelQuery(q, lnk)
	if numResult > -1 {
		log.Printf("Updated the relation between %s and %s %d\n", lnk.PUrl, lnk.SUrl, numResult)
	} else {
		q = "MATCH (url0:URL {name: {name0}})" +
			"MATCH (url1:URL {name: {name1}})" +
			"MERGE (url0)-[:link{size:{size}}]->(url1)"

		numResult := agent.executeRelQuery(q, lnk)
		log.Printf("Created a relation between %s and %s %d\n", lnk.PUrl,lnk.SUrl, numResult)
	}
}

func (agent Agent) insertOrUpdateNode(url string) {

	q := "MATCH (url:URL {name: {name}}) SET url.size = {size}"
	numResult := agent.executeNodeQuery(q, url)

	if numResult > -1 {
		log.Printf("Updated Node: %s\n", url)
	} else {
		q = "CREATE (url:URL {name: {name}, size: {size}})"
		numResult := agent.executeNodeQuery(q, url)

		log.Printf("Created Node: %d %s\n", numResult, url)
	}
}


func (agent Agent) executeNodeQuery(query string, url string) int64 {

	agent.mu.Lock()

	v := map[string]interface{}{
		"name": url,
		"size" : agent.getNumberOfUrls(url),
	}

	//log.Printf("%s %d\n", query, v)
	result, _ := agent.conn.ExecNeo(query, v)
	numResult, _ := result.RowsAffected()

	agent.mu.Unlock()

	return numResult
}


func (agent Agent) CountUrl(url string) {
	agent.client.HSet(agent.containerUrlName, url, agent.client.Incr(url).Val())
	agent.insertOrUpdateNode(url)
}

func (agent Agent) CountLink(lnk link) {
	key := fmt.Sprintf(formatLink, lnk.PUrl, lnk.SUrl)
	agent.client.HSet(agent.containerLinkName, key, agent.client.Incr(key).Val())
	agent.insertOrUpdateRel(lnk)
}

func (agent Agent) linkExists(lnk link) bool {
	key := fmt.Sprintf(formatLink, lnk.PUrl, lnk.SUrl)

	v, _ := agent.client.HGet(agent.containerLinkName, key).Int64()

	if v == 0 {
		return false
	}
	return true
}

func (agent Agent) UrlExists(url string) bool {

	v, _ := agent.client.HGet(agent.containerUrlName, url).Int64()

	if v == 0 {
		return false
	}
	return true
}

func (agent Agent) deleteContainer(containerName string) {

	r, err := agent.client.HGetAll(containerName).Result()
	if err == nil {
		for k, _ := range r{
			agent.client.Del(k)
		}
	}
}

func (agent Agent) resetNeo4j(wg *sync.WaitGroup){
	result, _ := agent.conn.ExecNeo("MATCH (n) DETACH DELETE n", nil)
	numResult, _ := result.RowsAffected()
	log.Printf("Neo4j was reset: <%d>", numResult)
	wg.Done()
}

func (agent Agent) resetRedis(wg *sync.WaitGroup) {

	agent.deleteContainer(agent.containerUrlName)
	agent.deleteContainer(agent.containerLinkName)

	agent.client.Del(agent.containerUrlName)
	agent.client.Del(agent.containerLinkName)

	log.Printf("Redis was reset...")
	wg.Done()
}

func (agent Agent) Reset() {
	var wg sync.WaitGroup
	wg.Add(2)

	go agent.resetNeo4j(&wg)
	go agent.resetRedis(&wg)

	wg.Wait()
}

func (agent Agent) ResetByConfig() {

	var count = 0

	if agent.config.Reset_Neo4j {
		count++
	}

	if agent.config.Reset_Redis {
		count++
	}

	if count == 0 {
		return
	} else {
		var wg sync.WaitGroup
		wg.Add(count)

		if  agent.config.Reset_Neo4j {
			go agent.resetNeo4j(&wg)
		}

		if agent.config.Reset_Redis {
			go agent.resetRedis(&wg)
		}

		wg.Wait()
	}
}


func (agent Agent) Close() {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := agent.client.Close()
		if err != nil {
			log.Printf(err.Error())
		}
		wg.Done()
	}()

	go func() {
		err := agent.conn.Close()
		if err != nil {
			log.Printf(err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
