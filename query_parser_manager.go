package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

/*
The QueryParserI type is an interface for queryParser query parser strategies.

The QueryParserManager stores a map of named QueryParserI pointers, one per parsing strategy supported by sg-restful.

 To add a new parsing strategy, one must implement the interface and register an instance with the QueryParserManager.

func init() {
   manager := GetQPManager()
   manager.AddQuery("mynewstrategy", &MyNewStrategy{})
}

sg-restful should be configured with a slice of strategies to use during the string parsing stage. This slice of strings is passed to the QueryParserManager.GetQueries method to retrieve a list of parsing strategies. For each one, sg-restful calls its CanParseString() method. If the strategy can parse the string, it is used and no further parsers are invoked, regardless of the success or failure of the parsing operation.
*/

// QueryParserI interface must be satisfied to define a parsing strategy for
// an input string. sg-restful may be configured to attempt to apply one or more
// parsing strategies depending upon the needs  of the deployment.
type QueryParserI interface {
	// determine if the QueryParser implementation can parse the string
	CanParseString(string) bool
	// Parse the string
	ParseString(string) (readFilters, error)
}

// QueryParsersList list of QyeryParserI
type QueryParsersList []QueryParserI

// QueryParsersMap maps strings t queryParserI interface pointers
type QueryParsersMap map[string]QueryParserI

// QueryParserManager keeps track of query parsers.
type QueryParserManager struct {
	queries QueryParsersMap
	active  []string
}

// AddParser accepts a pointer to an object which implements the
// QueryParserI interface and adds it to the internal map of queries
func (qpm *QueryParserManager) AddParser(name string, qpPtr QueryParserI) {
	if qpm.queries == nil {
		qpm.queries = make(QueryParsersMap)
	}
	qpm.queries[name] = qpPtr
}

// SetActiveParsers provides the mechanism by which one or more filters are
// identified by label to be active.
func (qpm *QueryParserManager) SetActiveParsers(names ...string) error {
	invalid := []string{}
	qpm.active = []string{}
	for _, name := range names {
		if _, ok := qpm.queries[name]; ok == true {
			log.Info("SetActiveQueries adding ", name)
			qpm.active = append(qpm.active, name)
		} else {
			invalid = append(invalid, name)
		}
	}
	if len(invalid) > 0 {
		panic(fmt.Sprintf("SetActiveQueries() - Invalid active queries set:%s", invalid))
	}
	return nil
}

// GetActiveParsers returns a list of quaryparsers corresponding with the name of queries supplied as a string slice
func (qpm *QueryParserManager) GetActiveParsers() ([]string, QueryParsersList) {
	parserList := QueryParsersList{}
	keyList := []string{}
	for _, key := range qpm.active {
		if val, ok := qpm.queries[key]; ok != true {
			panic("invalid key supplied to GetQueries:" + key)
		} else {
			keyList = append(keyList, key)
			parserList = append(parserList, val)
		}
	}
	log.Infof("GetActiveParsers() returning: (%s,%s)", keyList, parserList)
	return keyList, parserList
}

// GetActiveParserNames returns a slice of strings
func (qpm *QueryParserManager) GetActiveParserNames() []string {
	return qpm.active
}

// GetParsers returns a list of quaryparsers corresponding with the name of
// queries supplied as a string slice
func (qpm *QueryParserManager) GetParsers(names []string) QueryParsersList {

	returnList := QueryParsersList{}
	for _, key := range names {
		if val, ok := qpm.queries[key]; ok != true {
			panic("invalid key supplied to GetQueries:" + key)
		} else {
			returnList = append(returnList, val)
		}
	}
	log.Infof("GetParsers() returning %s", returnList)
	return returnList
}

// ResetActive clears all active parsers
func (qpm *QueryParserManager) ResetActive() {
	log.Debug("manager.ResetActive() called")
	qpm.active = []string{}
}

var qPManagerPtr *QueryParserManager

// GetQPManager returns a pointer to the QuerParserManager, which is responsible
// for keeping trakc of query strategies
func GetQPManager() *QueryParserManager {
	if qPManagerPtr == nil {
		qPManagerPtr = &QueryParserManager{}
	}
	log.Debug("GetQPManager() returning", qPManagerPtr)
	return qPManagerPtr
}
