package main

import (
	"flag"
	"fmt"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/partition/gpt"
	"github.com/openebs/node-disk-manager/pkg/blkid"
)

const (
	// BytesRequiredForGPTPartitionEntries is the total bytes required to store the GPT partition
	// entries. 128 bytes are required per partition, total no of partition supported by GPT = 128
	// Therefore, total bytes = 128*128
	BytesRequiredForGPTPartitionEntries = 16384

	// GPTPartitionStartByte is the byte on the disk at which the first partition starts.
	// Normally partition starts at 1MiB, (as done by fdisk utility). This is done to
	// align the partition start to physical block sizes on the disk.
	GPTPartitionStartByte = 1048576

	// NoOfLogicalBlocksForGPTHeader is the no. of logical blocks for the GPT header.
	NoOfLogicalBlocksForGPTHeader = 1

	// OpenEBSNDMPartitionName is the name meta info for openEBS created partitions.
	OpenEBSNDMPartitionName = "OpenEBS_NDM"
)

// Disk struct represents a disk which needs to be partitioned
type Disk struct {
	// DevPath is the /dev/sdX entry of the disk
	DevPath string
	// DiskSize is size of disk in bytes
	DiskSize uint64
	// LogicalBlockSize is the block size of the disk normally 512 or 4k
	LogicalBlockSize uint64

	table *gpt.Table

	disk *disk.Disk
}

func GetDiskInfo(diskName string) *Disk {

	disk := Disk{}

	fd, err := diskfs.Open(diskName)
	if err != nil {
		panic(err)
	}

	if _, err := fd.GetPartitionTable(); err == nil {
		panic("Disk contains partition table")
	}

	deviceIdentifier := blkid.DeviceIdentifier{DevPath: diskName}

	if fs := deviceIdentifier.GetOnDiskFileSystem(); len(fs) != 0 {
		panic("Disk is already partitioned")
	}

	disk.disk = fd
	disk.LogicalBlockSize = uint64(fd.LogicalBlocksize)
	disk.DiskSize = uint64(fd.Size)
	disk.DevPath = diskName
	disk.table = &gpt.Table{
		LogicalSectorSize: int(disk.LogicalBlockSize),
		ProtectiveMBR:     true,
	}

	// create partition
	var startSector, endSector uint64

	startSector = GPTPartitionStartByte / disk.LogicalBlockSize

	PrimaryPartitionTableSize := BytesRequiredForGPTPartitionEntries/disk.LogicalBlockSize + NoOfLogicalBlocksForGPTHeader
	endSector = (disk.DiskSize / disk.LogicalBlockSize) - PrimaryPartitionTableSize - 1

	disk.table.Partitions = append(disk.table.Partitions, &gpt.Partition{
		Name:  OpenEBSNDMPartitionName,
		Type:  gpt.LinuxFilesystem,
		Start: startSector,
		End:   endSector,
	})

	fmt.Println("{")
	fmt.Println("  \"device\": {")
	fmt.Println("    \"path\": \"" + disk.DevPath + "\",")
	fmt.Println("    \"size\": \"" + fmt.Sprintf("%d", disk.DiskSize) + "\",")
	fmt.Println("    \"logical_block_size\": \"" + fmt.Sprintf("%d", disk.LogicalBlockSize) + "\",")
	fmt.Println("    \"physical_block_size\": \"" + fmt.Sprintf("%d", disk.LogicalBlockSize) + "\",")
	fmt.Println("    \"partition_table\": \"gpt\"")
	fmt.Println("  },")
	fmt.Println("  \"table\": {")
	fmt.Println("    \"logical_sector_size\": \"" + fmt.Sprintf("%d", disk.table.LogicalSectorSize) + "\",")
	fmt.Println("    \"partitions\": [")
	for i, p := range disk.table.Partitions {
		fmt.Println("      {")
		fmt.Println("        \"name\": \"" + p.Name + "\",")
		fmt.Println("        \"type\": \"" + string(p.Type) + "\",")
		fmt.Println("        \"start\": \"" + fmt.Sprintf("%d", p.Start) + "\",")
		fmt.Println("        \"end\": \"" + fmt.Sprintf("%d", p.End) + "\",")
		fmt.Println("        \"attributes\": \"" + fmt.Sprintf("%d", p.Attributes) + "\",")
		fmt.Println("        \"number\": \"" + fmt.Sprintf("%d", i+1) + "\"")
		fmt.Println("      }")
		if i != len(disk.table.Partitions)-1 {
			fmt.Println(",")
		}
	}
	fmt.Println("    ]")
	fmt.Println("  }")
	fmt.Println("}")

	return &disk

}

func CreatePartition(disk *Disk) {

	err := disk.disk.Partition(disk.table)
	if err != nil {
		panic(err)
	}

}

func main() {

	diskName := flag.String("disk", "", "disk name")
	create := flag.Bool("create", false, "create partition")
	flag.Parse()

	disk := GetDiskInfo(*diskName)

	if *create {
		CreatePartition(disk)
	}

}
