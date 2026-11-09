package arangom

//go:generate go tool mockgen -destination=arangodb_mocks_test.go -package=arangom -mock_names=Database=MockArangoDB,Collection=MockArangoCollection,Graph=MockArangoGraph,View=MockArangoView,ArangoSearchView=MockArangoSearchView,ArangoSearchViewAlias=MockArangoSearchViewAlias,Analyzer=MockArangoSearchAnalyzer github.com/arangodb/go-driver/v2/arangodb Database,Collection,Graph,View,ArangoSearchView,ArangoSearchViewAlias,Analyzer
//go:generate go tool mockgen -destination=logger_mocks_test.go -package=arangom -source=logger.go
