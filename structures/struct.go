package structures

/*


 */

type MBR struct {
	Mbr_size           int64
	Mbr_creation_date  [30]byte
	Mbr_disk_signature int64
	Mbr_count          int8
	Mbr_partition_1    PARTITION
	Mbr_partition_2    PARTITION
	Mbr_partition_3    PARTITION
	Mbr_partition_4    PARTITION
}

/*type MBAR struct {
	Mbr_size           int64
	Mbr_disk_signature int64
	Mbr_creation_date  [20]byte
}*/

//ESTRUCTURA DE UNA PARTICION
type PARTITION struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
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
