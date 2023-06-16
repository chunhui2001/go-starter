package starter

import (
	"github.com/chunhui2001/go-starter/actions"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/chunhui2001/go-starter/core/middleware"
)

func init() {

	// simples
	if APP_SETTINGS.DemoEnable {

		// commons simples
		AppendRouter("GET", []string{"/labs-panic"}, actions.PanicRouter)
		AppendRouter("GET", []string{"/httpclient-simple"}, middleware.AccessInterceptor(true), actions.HttpClientSimpleRouter)
		AppendRouter("GET", []string{"/httpclient-timeout"}, actions.HttpClientTimeOutRouter)

		AppendRouter("GET", []string{"/labs-bigint"}, actions.BigRouter)
		AppendRouter("GET", []string{"/labs-ytld"}, actions.YtIdRouter)
		AppendRouter("GET", []string{"/labs-pem"}, actions.PemRouter)
		AppendRouter("GET", []string{"/labs-leftpad"}, actions.PadLeftRouter)
		AppendRouter("POST", []string{"/labs-upload-file"}, actions.UploadFileRouterOne)
		AppendRouter("GET", []string{"/labs-update-struct-pointer"}, actions.UpdateStructPointer)

		// redis simples
		AppendRouter("GET", []string{"/labs-redis-get"}, actions.RedisGetRouter)
		AppendRouter("GET", []string{"/labs-redis-set"}, actions.RedisSetRouter)
		AppendRouter("GET", []string{"/labs-redis-getset"}, actions.RedisGetSetRouter)
		AppendRouter("GET", []string{"/labs-redis-lpush"}, actions.RedisLpushRouter)
		AppendRouter("GET", []string{"/labs-redis-lrange"}, actions.RedisLrangeRouter)
		AppendRouter("GET", []string{"/labs-redis-del"}, actions.RedisDelRouter)
		AppendRouter("GET", []string{"/labs-redis-hset"}, actions.RedisHsetRouter)
		AppendRouter("GET", []string{"/labs-redis-hsetnx"}, actions.RedisDelRouter)
		AppendRouter("GET", []string{"/labs-redis-zincr"}, actions.RedisIncrRouter)
		AppendRouter("GET", []string{"/labs-redis-expire"}, actions.RedisExpireRouter)
		AppendRouter("GET", []string{"/labs-redis-setnx"}, actions.RedisSetNxRouter)
		AppendRouter("GET", []string{"/labs-redis-ttl"}, actions.RedisTtlRouter)
		AppendRouter("GET", []string{"/labs-redis-exists"}, actions.RedisExistsRouter)
		AppendRouter("POST", []string{"/labs-redis-pub"}, actions.RedisPubRouter)
		AppendRouter("POST", []string{"/labs-redis-producer"}, actions.RedisQueueProducerRouter)
		AppendRouter("GET", []string{"/labs-redis-consumer"}, actions.RedisQueueConsumerRouter)

		// xslt
		// AppendRouter("GET", []string{"/labs-xslt-demo"}, actions.XsltDemoRouter)

		// latex
		AppendRouter("GET", []string{"/labs-latex-demo"}, actions.LatexDemoRouter)

		// funcMaps
		AppendRouter("GET", []string{"/labs-funcMaps"}, actions.FuncMapsRouter)

		// validator data binding simples
		AppendRouter("POST", []string{"/demo/album-create"}, actions.AlbumCreateRouter)
		AppendRouter("GET", []string{"/demo/album-get"}, actions.AlbumGetRouter)
		AppendRouter("POST", []string{"/demo/binding-body"}, actions.BodyBindHandler)

		// elastic search simples
		AppendRouter("POST", []string{"/demo/els-create-or-Update"}, actions.ElsCreateOrUpdateRouter)
		AppendRouter("GET", []string{"/demo/els-cat-Indices"}, actions.ElsCatIndicesRouter)
		AppendRouter("POST", []string{"/demo/els-searcher"}, actions.ElsSearcherRouter)
		AppendRouter("POST", []string{"/demo/els-searcher-Advance"}, actions.ElsSearcherAdvanceRouter)
		AppendRouter("POST", []string{"/demo/els-dsl-Templdate"}, actions.ElsDslTemplateRouter)

		// mysql books
		AppendRouter("GET", []string{"/demo/books-query"}, actions.QueryBooksRouter)

		// mysql trx
		AppendRouter("GET", []string{"/demo/mysql-trans"}, actions.MySqlTxnsRouter)
		AppendRouter("GET", []string{"/demo/mysql-trans-Lock1"}, actions.MySqlTrxLocks1)

		// read congiguration demo
		AppendRouter("GET", []string{"/demo/ReadCacheKey"}, actions.ReadCacheKey)

		// other simples
		AppendRouter("POST", []string{"/websocket-client-simple"}, actions.WsClientSimple)
		AppendRouter("GET", []string{"/demo/ribbon-png"}, actions.RibbonDiagramsRouter)
		AppendRouter("GET", []string{"/demo/defer-func"}, actions.DeferRouter)

		// mysql transaction demo
		AppendRouter("GET", []string{"/demo/transactions"}, controller.TransactionRouter)
		AppendRouter("POST", []string{"/demo/sqlstatement"}, actions.RawSql)

		// aws
		AppendRouter("POST", []string{"/demo/awsv2-sign-simple"}, actions.AwsV2SignSimpleRouter)
		AppendRouter("GET", []string{"/demo/awsv2-sign-http"}, actions.AwsV2SignHttpClientRouter)

		// opensearch
		AppendRouter("GET", []string{"/go-board/es/indices"}, actions.OpenSearchIndicesRouter)
		AppendRouter("GET", []string{"/go-board/lastSnapshotDate"}, actions.OpenSearchLastSnapshotDateRouter)
		AppendRouter("POST", []string{"/go-board/dynamicQuery"}, actions.OpenSearchDynamicQueryRouter)
		AppendRouter("POST", []string{"/go-board/distinctQuery"}, actions.OpenSearchDistinctQueryRouter)
		AppendRouter("POST", []string{"/go-board/aggsMultipleQuery"}, actions.OpenSearchMultipleAggsQueryRouter)
		AppendRouter("POST", []string{"/go-board/ngJsonGenerator"}, actions.NdJsonHandler)
		AppendRouter("POST", []string{"/go-board/bulkQuery"}, actions.BulkQuery)

		// decimal
		AppendRouter("GET", []string{"/demo/NewFromString"}, actions.DecimalNewFromStringHandler)

		// ParallelLoop
		AppendRouter("GET", []string{"/demo/ParallelLoop0"}, actions.ParallelLoop0Handler)
		AppendRouter("GET", []string{"/demo/ParallelLoop1"}, actions.ParallelLoop1Handler)
		AppendRouter("GET", []string{"/demo/ParallelLoop2"}, actions.ParallelLoop2Handler)

		// mysql cat
		AppendRouter("POST", []string{"/demo/mysqlcat/listTables"}, actions.ListTables)

	}

}
