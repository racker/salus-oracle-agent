Provides an agent that gathers information from exported data from Oracle for reporting to a monitoring system such as Salus.

This is expected to be used in conjunction with Rackspace Managed Oracle hosting, where the oracle DBA's are setting up the DB jobs that export data in the correct format to monitor.

## Usage

```
  -configs string (default: ./config.d)
    	directory containing config files that define continuous monitoring (env AGENT_CONFIGS)
```

## Continuous-Monitoring Config File Format

When running the agent it will periodically collect Oracle telemetry at the interval configured in each config file. The option specifies a directory where any files in that directory that is of the correct format will be processed. The files are expected in JSON format and the structure of those files is:

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
- `interval` : an integer specifying the interval of package collection, in seconds.
- `type` : indicates the type or oracle monitoring to be started. Values must be one of the following "oracle_rman", "oracle_tablespace", or "oracle_dataguard"
- `filePath` : This filePath is the location where the agent expects to find the files associated with the database names. 
- `databaseNames` : This is a list of the names of the databases that should be monitored. These are expected to be .txt files of the same name as in this list in the filePath provided. 
- `errorCodeWhitelist` : this parameter is only valid for oracle_rman configurations. This is a list of expected error codes in the output file. If those error codes are discovered then it will be removed from the response
