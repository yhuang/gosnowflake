// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"time"
)

type resultFormat string

const (
	jsonFormat  resultFormat = "json"
	arrowFormat resultFormat = "arrow"
)

type execBindParameter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type execRequest struct {
	SQLText      string                       `json:"sqlText"`
	AsyncExec    bool                         `json:"asyncExec"`
	SequenceID   uint64                       `json:"sequenceId"`
	IsInternal   bool                         `json:"isInternal"`
	DescribeOnly bool                         `json:"describeOnly,omitempty"`
	Parameters   map[string]interface{}       `json:"parameters,omitempty"`
	Bindings     map[string]execBindParameter `json:"bindings,omitempty"`
}

type execResponseRowType struct {
	Name       string `json:"name"`
	ByteLength int64  `json:"byteLength"`
	Length     int64  `json:"length"`
	Type       string `json:"type"`
	Precision  int64  `json:"precision"`
	Scale      int64  `json:"scale"`
	Nullable   bool   `json:"nullable"`
}

type execResponseChunk struct {
	URL              string `json:"url"`
	RowCount         int    `json:"rowCount"`
	UncompressedSize int64  `json:"uncompressedSize"`
	CompressedSize   int64  `json:"compressedSize"`
}

// make all data field optional
type execResponseData struct {
	// succeed query response data
	Parameters         []nameValueParameter  `json:"parameters,omitempty"`
	RowType            []execResponseRowType `json:"rowtype,omitempty"`
	RowSet             [][]*string           `json:"rowset,omitempty"`
	RowSetBase64       string                `json:"rowsetbase64,omitempty"`
	Total              int64                 `json:"total,omitempty"`    // java:long
	Returned           int64                 `json:"returned,omitempty"` // java:long
	QueryID            string                `json:"queryId,omitempty"`
	SQLState           string                `json:"sqlState,omitempty"`
	DatabaseProvider   string                `json:"databaseProvider,omitempty"`
	FinalDatabaseName  string                `json:"finalDatabaseName,omitempty"`
	FinalSchemaName    string                `json:"finalSchemaName,omitempty"`
	FinalWarehouseName string                `json:"finalWarehouseName,omitempty"`
	FinalRoleName      string                `json:"finalRoleName,omitempty"`
	NumberOfBinds      int                   `json:"numberOfBinds,omitempty"`   // java:int
	StatementTypeID    int64                 `json:"statementTypeId,omitempty"` // java:long
	Version            int64                 `json:"version,omitempty"`         // java:long
	Chunks             []execResponseChunk   `json:"chunks,omitempty"`
	Qrmk               string                `json:"qrmk,omitempty"`
	ChunkHeaders       map[string]string     `json:"chunkHeaders,omitempty"`

	// ping pong response data
	GetResultURL      string        `json:"getResultUrl,omitempty"`
	ProgressDesc      string        `json:"progressDesc,omitempty"`
	QueryAbortTimeout time.Duration `json:"queryAbortsAfterSecs,omitempty"`
	ResultIDs         string        `json:"resultIds,omitempty"`
	ResultTypes       string        `json:"resultTypes,omitempty"`
	QueryResultFormat string        `json:"queryResultFormat,omitempty"`

	// async response placeholders
	AsyncResult *snowflakeResult `json:"asyncResult,omitempty"`
	AsyncRows   *snowflakeRows   `json:"asyncRows,omitempty"`
}

type execResponse struct {
	Data    execResponseData `json:"Data"`
	Message string           `json:"message"`
	Code    string           `json:"code"`
	Success bool             `json:"success"`
}

// QueryStatusFromServer status returned from server
type QueryStatusFromServer int

// Query Status defined at server side
const (
	SFQueryRunning = iota
	SFQueryAborting
    SFQuerySuccess
	SFQueryFailedWithError
	SFQueryAborted
	SFQueryQueued
	SFQueryFailedWithIncident
	SFQueryDisconnected
	SFQueryResumingWarehouse
    // SFQueryQueueRepairingWarehouse present in QueryDTO.java.
	SFQueryQueueRepairingWarehouse
	SFQueryRestarted

    // SFQueryBlocked the state when a statement is waiting on a lock on resource held by another statement.
	SFQueryBlocked
	SFQueryNoData
)

var sfQueryStrStatusMap = map[string]QueryStatusFromServer{
	"RUNNING":SFQueryRunning, "ABORTING":SFQueryAborting, "SUCCESS":SFQuerySuccess,
	"FAILED_WITH_ERROR":SFQueryFailedWithError, "ABORTED":SFQueryAborted, "QUEUED":SFQueryQueued,
	  "FAILED_WITH_INCIDENT":SFQueryFailedWithIncident, "DISCONNECTED":SFQueryDisconnected,
	  "RESUMING_WAREHOUSE":SFQueryResumingWarehouse, "QUEUED_REPAIRING_WAREHOUSE":SFQueryQueueRepairingWarehouse,
	  "RESTARTED":SFQueryRestarted,"BLOCKED":SFQueryBlocked, "NO_DATA":SFQueryNoData}

var dummy = struct{}{}
var sfqueryStatusRunning = map[QueryStatusFromServer]struct{} {
	SFQueryRunning:dummy, SFQueryResumingWarehouse:dummy, SFQueryQueued:dummy,
	SFQueryQueueRepairingWarehouse:dummy, SFQueryNoData:dummy}

var sfqueryStatusError = map[QueryStatusFromServer]struct{} {
	 SFQueryAborting:dummy, SFQueryFailedWithError:dummy, SFQueryAborted:dummy,
	SFQueryFailedWithIncident:dummy, SFQueryDisconnected:dummy, SFQueryBlocked: dummy}

type retStatus struct {
	Status  string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode int `json:"errorCode"`
}

type statusResponse struct {
	Data    struct {
		Queries []retStatus `json:"queries"`
	} `json:"data"`
	Message string           `json:"message"`
	Code    string           `json:"code"`
	Success bool             `json:"success"`
}

func strToSFQueryStatus(in string) QueryStatusFromServer {
	return sfQueryStrStatusMap[in]
}

func sfqStatusIsAnError(status QueryStatusFromServer) bool {
	_, exist := sfqueryStatusError[status]
	return exist
}

func sfqStatusIsStillRunning(status QueryStatusFromServer) bool {
	_, exist := sfqueryStatusRunning[status]
	return exist
}


