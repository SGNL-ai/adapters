# SGNL Adapters

This repository contains the public adapters that are shipped with SGNL.

## Code Structure

- `pkg/`: Contains the implementation of supported adapters.
- `cmd/adapter/main.go`: Responsible for running all adapters defined within `pkg`. New adapters MUST be registered via `RegisterAdapter`.
- `smoketests/`: Contains smoke tests for all supported adapters. These tests use `go-vcr.v3` to record data from a live instance on the first run, then use the cached responses for subsequent runs.
- `smoketests/fixtures/`: Contains example responses from live test instances for each supported adapter type. Before submitting any commits / PRs please ensure your fixture contains no PII and all secrets are redacted.

## Build

### Building a Docker Image

To build the Docker image for `adapters`, run the following command from the root of the repository:

```bash
docker build -t adapters:latest .
```

To run the container after building, use:

```bash
docker run -d --name adapters adapters:latest
```

### Building a Binary

To build and run the adapter as a binary, use the following commands:

```bash
go build -o adapters ./cmd/adapter
./adapters
```

If you encounter a permission error, make the binary executable:

```bash
chmod +x adapters
```

## Run

**Note:**
The adapter server requires an auth token on startup to initialize successfully. This is provided by the environment variable `AUTH_TOKENS_PATH`. The file must contain a JSON array of strings representing auth tokens. For example:

```json
["this-is-an-auth-token"]
```

To run the adapter locally for development and testing purposes (ran from the root of this repo), execute:

```bash
export AUTH_TOKENS_PATH="./shared-envs/ADAPTER_TOKENS"


# Run main.go
go run cmd/adapter/main.go

# OR if you have a previously built binary, you can run
./adapters

# OR if you choose to run the docker image, you'll need to provide the file and env variable.
# In a typical deployment, it is done by mounting a volume containing the file. For example:
docker run -d --name adapters \
    -v ./shared-envs/ADAPTER_TOKENS:/container/secrets/ADAPTER_TOKENS \
    -p 8080:8080 \
    -e AUTH_TOKENS_PATH=/container/secrets/ADAPTER_TOKENS \
    adapters:latest
```

### Fetch Data from a System of Record

By default, the adapter listens on port 8080. You can use Postman to send a gRPC request to the adapter by following these steps:

1. Create a new Request of type **gRPC**. Optionally save this in a new or existing Collection.

2. Under the **Service definition** tab, paste the following link or import the `GetPage` Protobuf definition: https://github.com/SGNL-ai/adapter-framework/blob/f2cafb0d963b54c350350967906ce59776d720a1/api/adapter/v1/adapter.proto

3. Set the gRPC request URL to the adapter server (e.g., `http://localhost:8080`) and select the `GetPage` method from the dropdown

4. In the **Metadata** tab, add a `token` key and set its value to one of the tokens in the `AUTH_TOKENS_PATH` file.

5. In the **Message** tab, enter the `GetPage` request following the schema defined in step 1.

An example gRPC request:

```json
{
  "cursor": "",
  "datasource": {
    "id": "Okta",
    "type": "Okta-1.0.1",
    "address": "{{address}}",
    "auth": {
      "http_authorization": "Bearer {{token}}"
    },
    "config": "{{b64_encoded_string}}"
  },
  "entity": {
    "attributes": [
      {
        "external_id": "id",
        "type": "ATTRIBUTE_TYPE_STRING",
        "id": "id"
      }
    ],
    "external_id": "User",
    "id": "User",
    "ordered": false
  },
  "page_size": "100"
}
```

The `config` should be a base64 encoded string of the `Config` struct defined in `config.go`. For example, if the `Config` struct is:

```go
type Config struct {
    APIVersion string `json:"apiVersion,omitempty"`
}
```

then the `config` field should be:

```json
{
  "apiVersion": "v1"
}
```

which is base64 encoded to `eyJhcGlWZXJzaW9uIjoidjEifQ==`.

Hit Send!
