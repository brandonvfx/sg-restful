package main

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

// QueryParserManager is
type QueryParserManager struct {
	queries QueryParsersMap
}

// AddQuery accepts a pointer to an object which implements the
// QueryParserI interface and adds it to the internal map of queries
func (qpm *QueryParserManager) AddQuery(name string, qpPtr QueryParserI) {
	qpm.queries[name] = qpPtr
}

// GetQueries returns a list of quaryparsers corresponding with the name of
// queries supplied as a string slice
func (qpm *QueryParserManager) GetQueries(names []string) QueryParsersList {
	returnList := QueryParsersList{}
	for _, key := range names {
		if val, ok := qpm.queries[key]; ok != true {
			panic("invalid key supplied to GetQueries:" + key)
		} else {
			returnList = append(returnList, val)
		}
	}
	return returnList
}

var qPManagerPtr *QueryParserManager

// GetQPManager returns a pointer to the QuerParserManager, which is responsible
// for keeping trakc of query strategies
func GetQPManager() *QueryParserManager {
	if qPManagerPtr == nil {
		qPManagerPtr = &QueryParserManager{}
	}
	return qPManagerPtr
}
