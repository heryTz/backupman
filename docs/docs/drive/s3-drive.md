---
title: S3 Drive
description: Store backups on Amazon S3 or S3-compatible services
---

# S3 Drive

The S3 drive allows you to store database backups on Amazon S3 or any S3-compatible storage service (MinIO, DigitalOcean Spaces, Wasabi, etc.).

## Configuration

```yaml
drives:
  - provider: s3
    label: S3 Drive
    bucket: my-backup-bucket
    region: us-east-1
    access_key: YOUR_ACCESS_KEY
    secret_key: YOUR_SECRET_KEY
    endpoint: "" # Optional: for S3-compatible services
    prefix: backups # Optional: folder prefix in bucket
    force_path_style: false # Required for some providers
```

## Configuration Options

| Option             | Required | Description                                |
| ------------------ | -------- | ------------------------------------------ |
| `provider`         | Yes      | Must be `s3`                               |
| `label`            | Yes      | Descriptive name for the drive             |
| `bucket`           | Yes      | S3 bucket name                             |
| `region`           | Yes      | AWS region (e.g., `us-east-1`)             |
| `access_key`       | Yes      | AWS access key ID                          |
| `secret_key`       | Yes      | AWS secret access key                      |
| `endpoint`         | No       | Custom endpoint for S3-compatible services |
| `prefix`           | No       | Folder prefix within bucket                |
| `force_path_style` | No       | URL style (see below)                      |

## URL Style Configuration

The `force_path_style` flag is crucial for provider compatibility:

### Amazon S3

```yaml
drives:
  - provider: s3
    endpoint: "" # Use AWS default
    force_path_style: false
```

### MinIO

```yaml
drives:
  - provider: s3
    endpoint: http://localhost:9000
    force_path_style: true
```

### DigitalOcean Spaces

```yaml
drives:
  - provider: s3
    endpoint: https://nyc3.digitaloceanspaces.com
    force_path_style: true
```

## Supported Providers

### Amazon S3

- Full compatibility
- Uses virtual-hosted style URLs by default
- Supports all AWS regions

### S3-Compatible Services

- **MinIO**: Self-hosted S3-compatible storage
- **DigitalOcean Spaces**: Object storage with S3 API
- **Wasabi**: Cloud storage with S3 compatibility
- **Backblaze B2**: Cloud storage with S3 compatibility
- **Linode Object Storage**: S3-compatible object storage

## Security Best Practices

### IAM Permissions

Your AWS credentials should have these minimum permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject", "s3:ListBucket"],
      "Resource": ["arn:aws:s3:::your-bucket-name", "arn:aws:s3:::your-bucket-name/*"]
    }
  ]
}
```

### Environment Variables

For better security, consider using environment variables:

```yaml
drives:
  - provider: s3
    label: S3 Drive
    bucket: ${S3_BUCKET}
    region: ${AWS_REGION}
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
```

## Examples

### Basic Amazon S3 Setup

```yaml
drives:
  - provider: s3
    label: Production Backups
    bucket: company-backups-prod
    region: us-west-2
    access_key: AKIAIOSFODNN7EXAMPLE
    secret_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    prefix: database-backups
    force_path_style: false
```

### MinIO Development Setup

```yaml
drives:
  - provider: s3
    label: Development Backups
    bucket: dev-backups
    region: us-east-1
    access_key: minioadmin
    secret_key: minioadmin123
    endpoint: http://localhost:9000
    prefix: dev-db-backups
    force_path_style: true
```

### Multiple S3 Providers

```yaml
drives:
  - provider: s3
    label: AWS Primary
    bucket: primary-backups
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
    force_path_style: false

  - provider: s3
    label: MinIO Backup
    bucket: secondary-backups
    region: us-east-1
    access_key: ${MINIO_ACCESS_KEY}
    secret_key: ${MINIO_SECRET_KEY}
    endpoint: http://minio.example.com:9000
    prefix: backup-storage
    force_path_style: true
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**

   - Check access key and secret key
   - Verify IAM permissions
   - Ensure correct region

2. **Bucket Not Found**

   - Verify bucket name spelling
   - Check bucket exists in correct region
   - Ensure proper permissions

3. **Connection Timeout**

   - Check endpoint URL for custom providers
   - Verify `force_path_style` setting
   - Check network connectivity

4. **SSL Certificate Errors**
   - For custom endpoints, consider using HTTP for testing
   - Ensure endpoint URL matches certificate

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
./backupman run --config config.yml --debug
```

## Performance Considerations

- **Multipart Upload**: Large files are automatically handled by AWS SDK
- **Concurrent Uploads**: Multiple drives upload in parallel
- **Compression**: Consider compressing large databases before upload
- **Lifecycle Policies**: Configure S3 lifecycle rules for automatic cleanup

## Cost Optimization

- **Storage Classes**: Use S3 Standard-IA for infrequent access
- **Lifecycle Policies**: Move old backups to Glacier
- **Compression**: Enable database dump compression
- **Cleanup**: Configure retention policies to remove old backups

---

## Examples

### Amazon S3 Production Setup

```yaml
drives:
  - provider: s3
    label: Production S3 Backups
    bucket: company-prod-backups
    region: us-west-2
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
    prefix: database-backups/prod
    force_path_style: false
```

### MinIO Development Environment

```yaml
drives:
  - provider: s3
    label: Development MinIO
    bucket: dev-backups
    region: us-east-1
    access_key: minioadmin
    secret_key: minioadmin123
    endpoint: http://localhost:9000
    prefix: dev-db-backups
    force_path_style: true
```

### DigitalOcean Spaces

```yaml
drives:
  - provider: s3
    label: DigitalOcean Spaces
    bucket: my-backup-space
    region: nyc3
    access_key: ${SPACES_ACCESS_KEY}
    secret_key: ${SPACES_SECRET_KEY}
    endpoint: https://nyc3.digitaloceanspaces.com
    prefix: database-backups
    force_path_style: true
```

### Multi-Region Redundancy

```yaml
drives:
  - provider: s3
    label: US-East Backups
    bucket: company-backups-us-east
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
    prefix: backups/us-east
    force_path_style: false

  - provider: s3
    label: EU-West Backups
    bucket: company-backups-eu-west
    region: eu-west-1
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
    prefix: backups/eu-west
    force_path_style: false
```

### Environment-Based Configuration

```yaml
drives:
  - provider: s3
    label: ${ENV} Backups
    bucket: ${ENV}-company-backups
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY_ID}
    secret_key: ${AWS_SECRET_ACCESS_KEY}
    prefix: backups/${ENV}/database
    force_path_style: false
```

### Docker Compose with MinIO

```yaml
# docker-compose.yml
services:
  backupman:
    image: backupman:latest
    environment:
      - S3_BUCKET=test-backups
      - S3_ACCESS_KEY=minioadmin
      - S3_SECRET_KEY=minioadmin123
      - S3_ENDPOINT=http://minio:9000
    volumes:
      - ./config.yml:/app/config.yml
    depends_on:
      - minio

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin123
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data

volumes:
  minio_data:
```

---

## Troubleshooting

### Authentication Issues

#### Access Denied Error

```
Error: Access Denied: Unable to create S3 client
```

**Solutions:**

1. Verify your access key and secret key are correct
2. Check IAM permissions include required S3 actions
3. Ensure credentials have access to the specified bucket
4. Verify region matches bucket location

#### Invalid Token Error

```
Error: InvalidToken: The provided token is malformed
```

**Solutions:**

1. Regenerate access keys in AWS console
2. Check for extra spaces or special characters in credentials
3. Use environment variables instead of hardcoded values

### Connection Issues

#### Timeout Errors

```
Error: operation error S3: PutObject, exceeded maximum number of attempts
```

**Solutions:**

1. Check network connectivity to S3 endpoint
2. Verify firewall allows HTTPS traffic to AWS
3. For custom endpoints, check endpoint URL is correct
4. Increase timeout in configuration (if supported)

#### SSL Certificate Errors

```
Error: certificate is not trusted
```

**Solutions:**

1. For custom endpoints, verify SSL certificate
2. Use HTTP for testing (not recommended for production)
3. Check endpoint URL matches certificate domain

### Bucket Issues

#### Bucket Not Found

```
Error: NoSuchBucket: The specified bucket does not exist
```

**Solutions:**

1. Verify bucket name spelling
2. Check bucket exists in specified region
3. Ensure bucket name follows S3 naming conventions
4. Verify you have permissions to access the bucket

### Provider-Specific Issues

#### MinIO Connection Issues

```
Error: connection refused
```

**Solutions:**

1. Verify MinIO is running and accessible
2. Check endpoint URL includes port number
3. Ensure `force_path_style: true`
4. Verify MinIO credentials are correct

#### DigitalOcean Spaces Issues

```
Error: InvalidAccessKeyId
```

**Solutions:**

1. Use Spaces access keys, not AWS keys
2. Set correct region-specific endpoint
3. Ensure `force_path_style: true`

### Common Configuration Mistakes

#### 1. Missing Region

```yaml
# Wrong
drives:
  - provider: s3
    bucket: my-bucket
    region: ""  # Missing region

# Correct
drives:
  - provider: s3
    bucket: my-bucket
    region: us-east-1
```

#### 2. Incorrect Force Path Style

```yaml
# Wrong for MinIO
force_path_style: false

# Correct for MinIO
force_path_style: true
```

#### 3. Hardcoded Credentials

```yaml
# Wrong - security risk
access_key: AKIAIOSFODNN7EXAMPLE
secret_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

# Correct - use environment variables
access_key: ${AWS_ACCESS_KEY_ID}
secret_key: ${AWS_SECRET_ACCESS_KEY}
```

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
# Run with debug output
./backupman run --config config.yml --debug

# Test S3 connection
./backupman health --config config.yml
```

### Testing Configuration

```bash
# Test S3 connection
aws s3 ls s3://your-bucket/

# Create test backup
./backupman run --config config.yml --test-mode
```
