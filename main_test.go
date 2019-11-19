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
	return time.Unix(1574201685, 1574201685000) // return specific timestamp
}

func (t *TestTimeInformation) getFileInformation(fileName string) time.Time {
	return time.Unix(1574201685, 1574201685000)
}


type Testconnection struct {
	mock.Mock
	value []byte
}

func (t *Testconnection) WriteToEnvoy(input []byte) {
	t.value = input
	t.Called(input)
}

func (t *Testconnection) Retry() error {

	return nil
}




// set the data for testing against here.


func TestDataguardOutput(t *testing.T) {
	testObj := new (Testconnection)
	var byteOutput = []byte("{\"Timestamp\":\"2019-11-19T14:40:59.201685-08:00\",\"Name\":\"dataguard\",\"Tags\":null,\"Fields\":{\"file_age\":\"2019-11-19T14:40:59.201685-08:00\",\"replication\":1,\"status\":\"success\"}}")

	testObj.On("WriteToEnvoy", byteOutput)
	conn = testObj
	timestamp = new (TestTimeInformation)
	var value = []string{"1"}
	createDataguardOutput(value, "notUsed")

	if string(testObj.value) != string(byteOutput) {
		t.Fail()
	}

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestRMANOutput(t *testing.T) {
	testObj := new (Testconnection)
	otherTestObj := new (TestTimeInformation)
	conn = testObj
	timestamp = otherTestObj
	var value = []string{"RMAN-12345","ORA-123","RMAN-456123"}
	var byteOutput = []byte("{\"Timestamp\":\"2019-11-19T14:40:59.201685-08:00\",\"Name\":\"RMAN\",\"Tags\":null,\"Fields\":{\"error_codes\":[\"RMAN-12345\",\"ORA-123\",\"RMAN-456123\"],\"file_age\":\"2019-11-19T14:40:59.201685-08:00\",\"status\":\"success\"}}")
	testObj.On("WriteToEnvoy", byteOutput)
	createRMANOutput(value, "notUsed")

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
}

func TestTablespaceOutput(t *testing.T) {
	testObj := new (Testconnection)
	otherTestObj := new (TestTimeInformation)
	conn = testObj
	timestamp = otherTestObj
	var value = []string{"SYSTEM", "2.59", "SYSAUX", "3.48"}
	var byteOutput = []byte("{\"Timestamp\":\"2019-11-19T14:40:59.201685-08:00\",\"Name\":\"Tablespace\",\"Tags\":{\"tablespace_name\":\"SYSTEM\"},\"Fields\":{\"file_age\":\"2019-11-19T14:40:59.201685-08:00\",\"status\":\"success\",\"usage\":\"2.59\"}}")
	var secondOutput = []byte("{\"Timestamp\":\"2019-11-19T14:40:59.201685-08:00\",\"Name\":\"Tablespace\",\"Tags\":{\"tablespace_name\":\"SYSAUX\"},\"Fields\":{\"file_age\":\"2019-11-19T14:40:59.201685-08:00\",\"status\":\"success\",\"usage\":\"3.48\"}}")

	testObj.On("WriteToEnvoy", byteOutput)
	testObj.On("WriteToEnvoy", secondOutput)
	//testObj.On("WriteToEnvoy", byteOutput)
	createTablespaceOutput(value, "notUsed")

	testObj.AssertCalled(t, "WriteToEnvoy", byteOutput)
	testObj.AssertCalled(t, "WriteToEnvoy", secondOutput)
}

func TestProcessRMAN(t *testing.T) {
	result := readFile("./testdata/RMAN.txt",processRMAN)

	if len(result) != 3 {
		t.Fail()
	}
}

func TestProcessTablespace(t *testing.T) {

	result := readFile("./testdata/tablespace.txt", processTablespace)
	assert.Equal(t, 4, len(result), "tablespace file only has two lines which should result in 4 values in the array")
	assert.Equal(t, "SYSTEM", result[0], "")
	assert.Equal(t, "2.59", result[1], "")
	assert.Equal(t, "SYSAUX", result[2], "")
	assert.Equal(t, "3.48", result[3] )
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

