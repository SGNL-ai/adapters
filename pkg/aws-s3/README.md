# AWS S3 Adapter

## Entity External IDs

The Entity External ID maps to a CSV filename in an S3 bucket. The adapter constructs the S3 object key as:

```
{prefix}/{entityExternalID}.{fileType}
```

The adapter is open-ended - any filename (without extension) can be used as an entity external ID.

| Entity External ID | S3 Object Key (with prefix `data/internal`) | Description |
|---|---|---|
| `users` | `data/internal/users.csv` | Users CSV file |
| `customers` | `data/internal/customers.csv` | Customers CSV file |
| `sales` | `data/internal/sales.csv` | Sales CSV file |

## How It Works

The Entity External ID is combined with the configured path prefix and file type to form the S3 object key. The adapter then fetches that object from the configured S3 bucket using the AWS SDK's `GetObject` API.

Example: with config `bucket: "my-bucket"`, `prefix: "exports/2024"`, entity external ID `"orders"`, file type `"csv"`:
- S3 path: `s3://my-bucket/exports/2024/orders.csv`

## Configuration

```json
{
  "region": "us-east-1",
  "bucket": "my-data-bucket",
  "prefix": "data/internal",
  "fileType": "csv"
}
```

## Notes

- The only supported file type is `csv`.
- Max page size is 1000.
- Attribute external IDs correspond to CSV column headers.
- Authentication uses AWS access key/secret key credentials.

## Adding New Entity External IDs

**No code changes required.** The adapter uses the entity external ID as the filename (without extension) in the S3 bucket. Simply configure a new entity in SGNL with a name matching the CSV file (e.g., if the file is `departments.csv`, set the entity external ID to `departments`).

Ensure the CSV file exists at `{prefix}/{entityExternalID}.csv` in the configured S3 bucket. Attribute external IDs must match the CSV column headers exactly.
