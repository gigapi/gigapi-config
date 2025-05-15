# <img src="https://github.com/user-attachments/assets/5b0a4a37-ecab-4ca6-b955-1a2bbccad0b4" />

### <img src="https://github.com/user-attachments/assets/74a1fa93-5e7e-476d-93cb-be565eca4a59" height=25 /> GigAPI Config

This repository provides the configuration options for Gigapi modules. Below is a detailed breakdown of all the available configuration options and their usage.

## Configuration Parameters

| Key                        | Type    | Default      | Description                                                           |
|----------------------------|---------|--------------|-----------------------------------------------------------------------|
| `gigapi.root`              | String  | `""`         | Root folder for all the data files.                                   |
| `gigapi.merge_timeout_s`   | Integer | `10`         | Base timeout between merges (in seconds).                            |
| `gigapi.save_timeout_s`    | Float   | `1`          | Timeout before saving the new data to disk (in seconds).             |
| `gigapi.no_merges`         | Boolean | `false`      | Disable merging functionality.                                       |
| `gigapi.ui`                | Boolean | `true`       | Enable UI for querier.                                               |
| `gigapi.mode`              | String  | `aio`        | Execution mode (`readonly`, `writeonly`, `compaction`, `aio`).       |
| `http.port`                | Integer | `7971`       | Port to listen on for the HTTP server.                               |
| `http.host`                | String  | `0.0.0.0`    | Host to bind to (0.0.0.0 for all interfaces).                        |
| `http.basic_auth.username` | String  | `""`         | Basic authentication username.                                       |
| `http.basic_auth.password` | String  | `""`         | Basic authentication password.                                       |
| `flightsql.port`           | Integer | `8082`       | Port to run the FlightSQL server.                                    |
| `flightsql.enable`         | Boolean | `true`       | Enable FlightSQL server.                                             |
| `loglevel`                 | String  | `info`       | Log level (`debug`, `info`, `warn`, `error`, `fatal`).               |

## Configuration Methods

### Environment Variables

You can configure the application using environment variables. Environment variable names are derived from the configuration keys by replacing dots (`.`) with underscores (`_`) and using uppercase letters. For example:

- `GIGAPI_ROOT`
- `GIGAPI_MERGE_TIMEOUT_S`
- `BASIC_AUTH_USERNAME`
- `HTTP_PORT`

### Configuration File

You can also provide a configuration file in any format supported by [Viper](https://github.com/spf13/viper) (e.g., JSON, YAML, TOML). Use the `InitConfig` function to specify the configuration file path.

### Example Configuration (YAML)
Below is an example configuration file in YAML format:

```YAML
gigapi:
  root: "/data"
  merge_timeout_s: 15
  save_timeout_s: 2.5
  no_merges: true
  ui: true
  mode: "aio"

http:
  port: 8080
  host: "127.0.0.1"
  basic_auth:
    username: "admin"
    password: "securepassword"

flightsql:
  port: 9090
  enable: true

loglevel: "debug"
```

#### Initialization and Default Handling
The `InitConfig` function initializes the configuration by either reading a file or using environment variables. If no explicit values are provided, defaults are automatically applied via the `setDefaults()` function.
