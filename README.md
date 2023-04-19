# Golang Disk Partitioning Tool

This is a golang tool for partitioning disks. The logic for this tool has been extracted from the OpenEBS NDM project (https://github.com/openebs/node-disk-manager).

Before running the code, please ensure that `GLIBC_2.35` is installed on your device. This can be done by running commands like `apt install libc6`.

## Usage

To use this tool, run the following command:

```
go run main.go -disk=<disk_name> -create
```

Replace `<disk_name>` with the name of the disk you want to partition.

The `-create` flag indicates that you want to create a new partition table on the specified disk.

## Output

When you run this tool, it will output JSON data containing information about the disk and its partitions. Here's an example of what it might look like:

```
{
  "device": {
    "path": "/dev/sda",
    "size": "10737418240",
    "logical_block_size": "512",
    "physical_block_size": "512",
    "partition_table": "gpt"
  },
  "table": {
    "logical_sector_size": "",
    "partitions": [
      {
        ...
      }
      ...
    ]
  }
}
```


