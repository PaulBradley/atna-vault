![Architecture](/img/vault.png?raw=true)

> noun: vault; a secure store where valuables are kept safe.

# ATNA Audit Vault

ATNA Audit Vault allows you to maintain a secure long-term backup/archive for all your IHE audit messages.

IHE vendors who can provide "filter forward" functionally within their Audit Trail and Node Authentication sub-systems can forward copies of all audit messages to the vault.

The ATNA Vault software contains two commands. Freeze and Thaw.

## Freeze

The `freeze` daemon accepts incoming messages over TLS/TCP listening on port 6514. Each message contains a [Syslog header](https://datatracker.ietf.org/doc/html/rfc5424#section-6) followed by an XML [audit message](https://datatracker.ietf.org/doc/html/rfc3881#section-5). `freeze` performs gzip compression on the raw audit message saving in the region of 50-70%, depending on the message type.

Each message is given a [Universally Unique Lexicographically Sortable Identifier](https://github.com/ulid/spec) as the object filename. The first five characters of the identifier are used as an [S3 prefix](https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-prefixes.html) to group a day's worth of audit messages within a single folder/prefix.

Each message gets persisted to a secure AWS S3 bucket which has encryption at rest enabled. Each object is placed in the [Glacier Instant Retrieval](https://aws.amazon.com/s3/storage-classes/?nc=sn&loc=3#Instant_Retrieval) storage tier, which costs $0.005 per GB. The `freeze` daemon also sets the S3 object lock property with a user-defined retention period so that you can enforce retention policies as an added layer of data protection for your regulatory compliance.

The `freeze` daemon also supports the additional tagging of each object to help with billing analysis. 

### Glacier Instant Retrieval Storage Class

Amazon S3 Glacier Instant Retrieval is an archive storage class that delivers the lowest-cost storage for long-lived data that is rarely accessed and requires retrieval in milliseconds. It also offers the following benefits:

* Data retrieval in milliseconds with the same performance as S3 Standard
* Designed for durability of 99.999999999% of objects across multiple Availability Zones
* Data is resilient in the event of the destruction of one entire Availability Zone
* Designed for 99.9% data availability in a given year

### Architecture Diagram

![Architecture](/img/architecture.png?raw=true)

### Starting the Freeze Daemon

The `freeze` daemon is launched using [socat](http://www.dest-unreach.org/socat/), the multipurpose relay available on most Linux distributions. The command below shows how to launch the daemon on port 6514 with specified SSL certificate files.

``` bash
socat -lf socat.log
    openssl-listen:6514,
    cert=sslcertificate.crt,
    key =sslcertificate.key,
    verify=0,
    fork system:/path/to/anta-vault/freeze &
```

The `freeze` daemon can also been as a UDP listener like:

``` bash
socat -lf socat.dat
    -u udp4-recvfrom:6514,
    reuseport,
    fork system:/path/to/anta-vault/freeze &
```

### Command Line Options

The `freeze` daemon supports the following command line parameters.

| parameter       | description                                                                                             |
| --------------- | ------------------------------------------------------------------------------------------------------- |
| --store-locally | Instructs AV to store a copy of the messages on the local server as well as in S3. Suitable for testing |
| -bucketname     | AWS S3 bucket name to upload files to                                                                   |
| -region         | AWS region that the target bucket lives in e.g. eu-west-2                                               |


### Environment Variables

The `freeze` daemon uses the following environment variables to configure its APIs.

| variable  | description                                                                                |
| ----------| ------------------------------------------------------------------------------------------ |
| S3_KEY    | The access key for the IAM account which has has permissions to write to the target bucket |
| S3_SECRET | The secret key for the IAM account which has has permissions to write to the target bucket |


## Thaw

The `thaw` daemon indexes and parses the syslog headers and audit messages and creates an optimised index which can be use to fullfil IG queries.

## ATNA Message Types

The following ATNA message types will be supported:

* Application Start
* Application Stop
* Audit Log Used
* Cross Gateway Patient Discovery
* Delete Document Set
* Document Metadata Notify
* Document Metadata Subscribe
* Login
* Logout
* Notify XAD-PID Link Change
* Patient Demographics Query
* Patient Identity Feed
* PIX Query
* Provide and Register Document Set-b
* Register Document Set-b
* Registry Stored Query
* Retrieve Document Set
* Security Administrative
* Update Document Set
* User
