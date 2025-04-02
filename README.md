# Mini Port Scanner

`mpscan` is a simple command-line utility that concurrently scans a given IP address or hostname for open ports and reveals any banner information available on those ports. You can specify multiple targets, ports, port ranges as well as output the results to a file in JSON format.

## Video Demonstration

Watch a video demo of the program on [YouTube](https://youtu.be/nHhIzGrHFOQ?si=21GqDZ6JCSVZ9zGK).

## Building & Running From Source

### Prerequisites

- a relatively modern version of Go installed on your system

### Steps

1. Clone the repository

```
git clone https://github.com/andreshungbz/mpscan.git
```

2. Change to the project directory

```
cd mpscan
```

3. Build the project using `make`

```
make
```

> [!NOTE]
> If you don't have `make` installed, you can manually build with `go build -o bin/`

4. Run the program with the following command, adding flags as needed:

```
./bin/mpscan
```

## Example Usage & Output

Usage

```
./bin/mpscan -targets=scanme.nmap.org,go.dev -ports=22,80,443 -json
```

Output

```
[SCAN START]

[go.dev]          [3/3 ports] [=======================================================] 100 %
[scanme.nmap.org] [3/3 ports] [=======================================================] 100 %

[BANNERS]

[go.dev:80] Google Frontend
[scanme.nmap.org:22] SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
[scanme.nmap.org:80] Apache/2.4.7 (Ubuntu)

[SCAN SUMMARY]

[scanme.nmap.org]
Total Ports Scanned: 3
Open Ports Count: 2
Open Ports: [22 80]
Time Taken: 0.262s

[go.dev]
Total Ports Scanned: 3
Open Ports Count: 2
Open Ports: [80 443]
Time Taken: 5.001s

[JSON OUTPUT SAVED: 20250331-215629-mpscan.json]
```

Contents of JSON file `20250331-215629-mpscan.json` (additionally formatted)

```json
[
  {
    "Hostname": "scanme.nmap.org",
    "TotalPortsScanned": 3,
    "OpenPortCount": 2,
    "OpenPorts": [22, 80],
    "TimeTaken": 262353999
  },
  {
    "Hostname": "go.dev",
    "TotalPortsScanned": 3,
    "OpenPortCount": 2,
    "OpenPorts": [80, 443],
    "TimeTaken": 5001092924
  }
]
```

## Options

| Flag               | Description                                                                                                      |
| ------------------ | ---------------------------------------------------------------------------------------------------------------- |
| `-target`          | The hostname or IP address to be scanned.                                                                        |
| `-start-port`      | The lower bound port to begin scanning. (default 1)                                                              |
| `-end-port`        | The upper bound port to finish scanning. (default 1024)                                                          |
| `-workers`         | The number of concurrent goroutines to launch per target. (default 100)                                          |
| `-timeout`         | The maximum time in seconds to wait for connections to be established. (default 5)                               |
| `-ports`           | Comma-separated list of ports (e.g., -ports=22,80,443). Setting this overrides -start-port and -end-port.        |
| `-targets`         | Comma-separated list of targets (e.g., -targets=localhost,scanme.nmap.org). Targets are aggregated with -target. |
| `-json` (boolean)  | Indicates whether to also output a JSON file of the scan results.                                                |
| `-debug` (boolean) | Displays flag values for debugging.                                                                              |

This information can also be reviewed by running

```
./bin/mpscan -help
```

## Notes

- The maximum number of retries to attempt to connect to a port is 1. This can be changed in `connection.attemptScan`, but increasing it will generally make the scan slower to complete. A random 0.0 to 1.0 multiplier is also applied to the backoff timer.
- If you set `-ports`, then `-start-port` and `-end-port` are ignored.
- Setting both `-target` and `-targets` will result in the targets being aggregated. For example, if you set `-target=localhost` and `-targets=scanme.nmap.org`, the program will scan both `localhost` and `scanme.nmap.org`. If both are empty, the program will default to `localhost`.
- `-timeout` sets the timeout for both the establishment of the attempted TCP connection, and the attempt to grab the banner for an open port.
- Outputs in banner grabbing may vary due to timeouts. In that case, set an increased `-timeout` value.
- For multiple targets, the rest of the flags are shared across all targets.

## Cleaning Up

To remove the build folder and generated JSON files, run:

```
make clean
```
