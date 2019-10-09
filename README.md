# prtg-grpc-health-sensor
PRTG Custom Sensor for integrating with a `grpc_health_v1` endpoint

## How to use
1. Build the sensor by running `go build` (or download the latest release 
[here](https://github.com/vanti-public/prtg-grpc-health-sensor/releases))
2. Copy the `prtg-grpc-health-sensor.exe` file into your `C:\Program Files (x86)\PRTG Network Monitor\Custom Sensors\EXEXML`
 folder 
3. In PRTG, add a new `EXE/Script Advanced` sensor
    1. Set the Name and override any other defaults you want
    2. The 'Parameters' setting is passed to the sensor as command line arguments. The following options are available:
        - `-address=<host:port>` sets the address of the gRPC server (defaults to `localhost:9001`)
        - `-service=<service_name>` (Optional) Sets the name of the service to query (defaults to none)
        
That's it. When everything is working correctly the sensor will provide the response time for the health check, or an
 error if the service is unreachable, or the health is set to `NOT_SERVING`