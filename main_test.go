package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type TestTimeInformation struct {
}

func (t *TestTimeInformation) Now() int64 {
	var time, _ = time.Parse(time.RFC3339, "2019-11-19T22:40:59.201685Z")
	// Our system is expecting timestamps in milliseconds
	return time.Unix() * 1000
}

func (t *TestTimeInformation) getFileInformation(fileName string) int64 {
	var time, _ = time.Parse(time.RFC3339, "2019-11-19T22:40:59.201685Z")
	// Our system is expecting timestamps in milliseconds
	return time.Unix() * 1000
}


type Testconnection struct {
	mock.Mock
}

func (t *Testconnection) WriteToEnvoy(input string) {
	t.Called(input)
}

func (t *Testconnection) Retry() error {
	return nil
}



func TestDataguardOutput(t *testing.T) {
	testObj := new (Testconnection)
	var byteOutput = string("{\"Timestamp\":1574203259,\"Name\":\"oracle_dataguard\",\"Tags\":null,\"Fields\":{\"file_age\":1574203259,\"replication\":1,\"status\":\"success\"}}")
	conn = testObj

	testObj.On("WriteToEnvoy", mock.Anything)

	timestamp = new (TestTimeInformation)
	var value = []string{"1"}
	createDataguardOutput(value, "notUsed", nil)

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestRMANOutput(t *testing.T) {
	testObj := new (Testconnection)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{"RMAN-12345","ORA-123","RMAN-456123"}
	var byteOutput = string("{\"Timestamp\":1574203259,\"Name\":\"oracle_rman\",\"Tags\":null,\"Fields\":{\"error_codes\":[\"RMAN-12345\",\"ORA-123\",\"RMAN-456123\"],\"file_age\":1574203259,\"status\":\"success\"}}")
	testObj.On("WriteToEnvoy", mock.Anything)
	createRMANOutput(value, "notUsed", nil)

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestTablespaceOutput(t *testing.T) {
	testObj := new (Testconnection)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{"SYSTEM", "2.59", "SYSAUX", "3.48"}
	var systemTableOutput = string("{\"Timestamp\":1574203259,\"Name\":\"oracle_tablespace\",\"Tags\":{\"tablespace_name\":\"SYSTEM\"},\"Fields\":{\"file_age\":1574203259,\"status\":\"success\",\"usage\":\"2.59\"}}")
	var sysauxTableOutput = string("{\"Timestamp\":1574203259,\"Name\":\"oracle_tablespace\",\"Tags\":{\"tablespace_name\":\"SYSAUX\"},\"Fields\":{\"file_age\":1574203259,\"status\":\"success\",\"usage\":\"3.48\"}}")

	testObj.On("WriteToEnvoy", mock.Anything)
	createTablespaceOutput(value, "notUsed", nil)

	testObj.AssertCalled(t, "WriteToEnvoy", systemTableOutput)
	testObj.AssertCalled(t, "WriteToEnvoy", sysauxTableOutput)
}

func TestRMANOutputSucceedsWithNoErrorCodes(t *testing.T) {
	testObj := new (Testconnection)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{}
	var byteOutput = string("{\"Timestamp\":1574203259,\"Name\":\"oracle_rman\",\"Tags\":null,\"Fields\":{\"error_codes\":[],\"file_age\":1574203259,\"status\":\"success\"}}")
	testObj.On("WriteToEnvoy", mock.Anything)
	createRMANOutput(value, "notUsed", nil)

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestProcessRMAN(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_RMAN",
		databaseName: "RMAN",
		filePath: "./testdata/",
		errorCodeWhitelist: []string{"RMAN-12345", "ORA-123"},
	}


	result, err := readFile("./testdata/RMAN.txt", config, processRMAN)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result), "RMAN Expected 2 error codes")
	assert.Equal(t, "RMAN-1234", result[0])
	assert.Equal(t, "RMAN-456123", result[1])
}

func TestProcessRMANWhitelistAllErrorCodes(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_RMAN",
		databaseName: "RMAN",
		filePath: "./testdata/",
		errorCodeWhitelist: []string{"RMAN-12345", "ORA-123", "RMAN-456123", "RMAN-1234"},
	}


	result, err := readFile("./testdata/RMAN.txt", config, processRMAN)
	assert.Nil(t, err)
	assert.Empty(t, result, "RMAN expecte to whitelist all error codes")
}

func TestProcessRMANSucceeds(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_RMAN",
		databaseName: "RMAN",
		filePath: "./testdata/",
		errorCodeWhitelist: []string{"RMAN-12345", "ORA-123", "RMAN-1234", "RMAN-456123"},
	}


	result, err := readFile("./testdata/RMAN.txt", config, processRMAN)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result), "All error codes should be in whitelist")
}


func TestProcessTablespace(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_tablespace",
		databaseName: "tablespace",
		filePath: "./testdata/",
		errorCodeWhitelist: nil,
	}

	result, err := readFile("./testdata/tablespace.txt", config, processTablespace)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(result), "tablespace file only has two lines which should result in 4 values in the array")
	assert.Equal(t, "SYSTEM", result[0])
	assert.Equal(t, "2.59", result[1])
	assert.Equal(t, "SYSAUX", result[2])
	assert.Equal(t, "3.48", result[3])
}

func TestProcessDataguard(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_dataguard",
		databaseName: "dataguard",
		filePath: "./testdata/",
		errorCodeWhitelist: nil,
	}

	result, err := readFile("./testdata/dataguard.txt", config, processDataguard)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result), "Dataguard should always only return one result")
	assert.Equal(t, "1", result[0], "Testing Value should be the string \"1\"")
}

func TestProcessNoFile(t *testing.T) {
	config := Configuration {
		interval: 30,
		configType: "oracle_RMAN",
		databaseName: "RMAN",
		filePath: "./testdata/",
		errorCodeWhitelist: []string{"RMAN-12345", "ORA-123"},
	}
	result, err := readFile("./testdata/RMANDoesntExist.txt", config, processRMAN)
	assert.Error(t, err, "File Not Found")
	assert.EqualError(t, err, "open ./testdata/RMANDoesntExist.txt: no such file or directory")
	assert.Nil(t, result, "Expected results to be nil since file does not exist")
}

func TestReadConfigsFromPath(t *testing.T) {
	readConfigsFromPath("./testdata/config")

	assert.NotEmpty(t, tickers)
	assert.NotNil(t, tickers["RMAN.json:RMAN"])
	assert.NotNil(t, tickers["tablespace"])
	assert.NotNil(t, tickers["dataguard"])
}



