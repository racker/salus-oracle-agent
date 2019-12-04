package main

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"time"
)

type TestTimeInformation struct {
}

func (t *TestTimeInformation) Now() time.Time {
  var time, _ = time.Parse(time.RFC3339, "2019-11-19T22:40:59.201685Z")
	return time
}

func (t *TestTimeInformation) getFileInformation(fileName string) time.Time {
  var time, _ = time.Parse(time.RFC3339, "2019-11-19T22:40:59.201685Z")
	return time
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
	var byteOutput = "{\"Timestamp\":\"2019-11-19T22:40:59.201685Z\",\"Name\":\"dataguard\",\"Tags\":null,\"Fields\":{\"file_age\":\"2019-11-19T22:40:59.201685Z\",\"replication\":1,\"status\":\"success\"}}"
	conn = testObj

	testObj.On("WriteToEnvoy", mock.Anything)

	timestamp = new (TestTimeInformation)
	var value = []string{"1"}
	createDataguardOutput(value, "notUsed")

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestRMANOutput(t *testing.T) {
	testObj := new (Testconnection)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{"RMAN-12345","ORA-123","RMAN-456123"}
	var byteOutput = "{\"Timestamp\":\"2019-11-19T22:40:59.201685Z\",\"Name\":\"RMAN\",\"Tags\":null,\"Fields\":{\"error_codes\":[\"RMAN-12345\",\"ORA-123\",\"RMAN-456123\"],\"file_age\":\"2019-11-19T22:40:59.201685Z\",\"status\":\"success\"}}"
	testObj.On("WriteToEnvoy", mock.Anything)
	createRMANOutput(value, "notUsed")

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestTablespaceOutput(t *testing.T) {
	testObj := new (Testconnection)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{"SYSTEM", "2.59", "SYSAUX", "3.48"}
	var systemTableOutput = "{\"Timestamp\":\"2019-11-19T22:40:59.201685Z\",\"Name\":\"Tablespace\",\"Tags\":{\"tablespace_name\":\"SYSTEM\"},\"Fields\":{\"file_age\":\"2019-11-19T22:40:59.201685Z\",\"status\":\"success\",\"usage\":\"2.59\"}}"
	var sysauxTableOutput = "{\"Timestamp\":\"2019-11-19T22:40:59.201685Z\",\"Name\":\"Tablespace\",\"Tags\":{\"tablespace_name\":\"SYSAUX\"},\"Fields\":{\"file_age\":\"2019-11-19T22:40:59.201685Z\",\"status\":\"success\",\"usage\":\"3.48\"}}"

	testObj.On("WriteToEnvoy", mock.Anything)
	createTablespaceOutput(value, "notUsed")
	testObj.AssertCalled(t, "WriteToEnvoy", systemTableOutput)
	testObj.AssertCalled(t, "WriteToEnvoy", sysauxTableOutput)
}

func TestProcessRMAN(t *testing.T) {
	result := readFile("./testdata/RMAN.txt",processRMAN)

	assert.Equal(t, 3, len(result), "RMAN file should have 3 error codes in it")
	assert.Equal(t, "RMAN-12345", result[0])
	assert.Equal(t, "ORA-123", result[1])
	assert.Equal(t, "RMAN-456123", result[2])
}

func TestProcessTablespace(t *testing.T) {

	result := readFile("./testdata/tablespace.txt", processTablespace)
	assert.Equal(t, 4, len(result), "tablespace file only has two lines which should result in 4 values in the array")
	assert.Equal(t, "SYSTEM", result[0])
	assert.Equal(t, "2.59", result[1])
	assert.Equal(t, "SYSAUX", result[2])
	assert.Equal(t, "3.48", result[3])
}

func TestProcessDataguard(t *testing.T) {
	result := readFile("./testdata/dataguard.txt", processDataguard)

	assert.Equal(t, 1, len(result), "Dataguard should always only return one result")
	assert.Equal(t, "1", result[0], "Testing Value should be the string \"1\"")
}

func TestProcessNoFile(t *testing.T) {
	result := readFile("./testdata/RMANDoesntExist.txt", processRMAN)
	assert.Nil(t, result, "Expected results to be nil since file does not exist")
}

