# <img src="https://github.com/user-attachments/assets/5b0a4a37-ecab-4ca6-b955-1a2bbccad0b4" />

### <img src="https://github.com/user-attachments/assets/74a1fa93-5e7e-476d-93cb-be565eca4a59" height=25 /> GigAPI Config

This repository provides the configuration options for Gigapi modules. Below is a detailed breakdown of all the available configuration options and their usage.

## Configuration Structure

The configuration is structured into the following sections:

- **Gigapi**
- **BasicAuth**
- **Global Settings**

### 1. Gigapi Configuration

This section contains configurations specific to Gigapi functionality.

| Key               | Type    | Default | Description                                   |
|--------------------|---------|---------|-----------------------------------------------|
| `gigapi.root`      | String  | `""`    | Root folder for all the data files.           |
| `gigapi.merge_timeout_s` | Integer | `10`    | Base timeout between merges (in seconds).    |
| `gigapi.save_timeout_s`  | Float   | `1`     | Timeout before saving the new data to disk (in seconds). |
| `gigapi.no_merges` | Boolean | `false` | Disable merging functionality.                |

### 2. Basic Authentication Configuration

This section provides credentials for basic authentication.

| Key              | Type    | Default | Description                                   |
|-------------------|---------|---------|-----------------------------------------------|
| `basic_auth.username` | String  | `""`    | Username for basic authentication.           |
| `basic_auth.password` | String  | `""`    | Password for basic authentication.           |

### 3. Global Settings

These are general configurations that apply to the entire application.

| Key                | Type    | Default   | Description                                   |
|---------------------|---------|-----------|-----------------------------------------------|
| `port`             | Integer | `7971`    | HTTP port to listen on.                       |
| `host`             | String  | `0.0.0.0` | Host to bind to (e.g., `0.0.0.0` for all interfaces). |
| `flightsql_port`   | Integer | `8082`    | FlightSQL port to listen on.                  |
| `disable_ui`       | Boolean | `false`   | Disable the UI for the querier.               |
| `loglevel`         | String  | `info`    | Log level (options: `debug`, `info`, `warn`, `error`, `fatal`). |
| `mode`             | String  | `aio`     | Execution mode (options: `readonly`, `writeonly`, `compaction`, `aio`). |

## Configuration Methods

### Environment Variables

You can configure the application using environment variables. Environment variable names are derived from the configuration keys by replacing dots (`.`) with underscores (`_`) and using uppercase letters. For example:

- `GIGAPI_ROOT`
- `GIGAPI_MERGE_TIMEOUT_S`
- `BASIC_AUTH_USERNAME`
- `PORT`

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

basic_auth:
  username: "admin"
  password: "securepassword"

port: 8080
host: "127.0.0.1"
flightsql_port: 9090
disable_ui: true
loglevel: "debug"
mode: "readonly"
```

#### Initialization and Default Handling
The `InitConfig` function initializes the configuration by either reading a file or using environment variables. If no explicit values are provided, defaults are automatically applied via the `setDefaults()` function.
