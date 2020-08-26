package structures

/*


 */

type MBR struct {
	Mbr_size           int64
	Mbr_creation_date  [19]byte
	Mbr_disk_signature int64
	Mbr_count          int8
	Mbr_partition_1    PARTITION
	Mbr_partition_2    PARTITION
	Mbr_partition_3    PARTITION
	Mbr_partition_4    PARTITION
}

//ESTRUCTURA DE UNA PARTICION
type PARTITION struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
	Part_active bool
}

//ESTRUCTURA DEL Extended Boot Record
type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}
type MOUNT struct {
	Mount_id        string
	Mount_path      string
	Mount_particion string
	Mount_estado    bool
}
