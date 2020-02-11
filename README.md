Provides an agent that gathers information from exported data from oracle for reporting to a monitoring systems such as Salus.

## Usage

```
  -configs string (default: ./config.d)
    	directory containing config files that define continuous monitoring (env AGENT_CONFIGS)
```

## Continuous-Monitoring Config File Format

When running the agent it will periodically collect package telemetry at the interval configured in each config file. The option specifies a directory where any files in that directory that is of the correct format will be processed. The files are expected in JSON format and the structure of those files is:

```json
{
  "interval":  30,
  "type":"oracle_rman",
  "databaseNames":["RMAN"],
  "filePath":"./testdata/",
  "errorCodeWhitelist":  ["RMAN-1234"]
}
```

where:
- `interval` : a Go int specifying the interval of package collection, in seconds.
- `type` : indicates the type or oracle monitoring you would like to do. Values must be one of the following "oracle_rman", "oracle_tablespace", or "oracle_dataguard"
- `filePath` : This filePath is the location of where we expect to find the files associated with the database names. 
- `databaseNames` : This is a list of the names of the databases that we want to monitor. These are expected to be .txt files of the same name as in this list in the filePath provided. 
- `errorCodeWhitelist` : this parameter is only valid for oracle_rman configurations. This is a list of expected error codes in the output file. If those error codes are discovered then it will be removed from the response

