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
	Mbr_Ext            int8 //sirve para ver si ya tiene una extendida
}

//ESTRUCTURA DE UNA PARTICION
type PARTITION struct {
	//PARTICION PRIMARIA
	Part_status  int8
	Part_type    byte
	Part_fit     byte
	Part_start   int64
	Part_end     int64
	Part_size    int64
	Part_name    [16]byte
	Part_isEmpty int8
}

type EXTENDED struct {
	Part_status    int8
	Part_type      byte
	Part_fit       byte
	Part_start     int64
	Part_end       int64
	Part_size      int64
	Part_name      [16]byte
	Part_partition []PARTITION
	Part_ebr       []EBR
}

//ESTRUCTURA DEL Extended Boot Record
type EBR struct {
	Part_status int8
	Part_fit    byte
	Part_start  int64
	Part_end    int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}
type MOUNT struct {
	Mount_id        string
	Mount_path      string
	Mount_particion string
	Mount_usrtxt    string
	Mount_estado    bool
	Mount_part      PARTITION
}

type USER struct {
	User_id       string
	User_type     string
	User_group    string
	User_username string
	User_password string
	User_isLoged  bool
}

//===============  LWH

type SUPERBOOT struct {
	SB_hd_name                   [20]byte
	SB_date                      [20]byte
	SB_date_lstmount             [20]byte
	SB_AVD_count                 int64
	SB_AVD_details_count         int64
	SB_Inodes_count              int64
	SB_blocks_count              int64
	SB_AVD_free                  int64
	SB_Inodes_free               int64
	SB_blocks_free               int64
	SB_mount_count               int64
	SB_ap_bitmap_tree_dir        int64
	SB_ap_tree_dir               int64
	SB_ap_bitmap_detail_dir      int64
	SB_ap_detail_dir             int64
	SB_ap_bitmap_table_inode     int64
	SB_ap_table_inode            int64
	SB_ap_bitmap_blocks          int64
	SB_ap_blocks                 int64
	SB_ap_log                    int64
	SB_size_struct_tree_dir      int64
	SB_size_struct_detail_dir    int64
	SB_size_struct_inodo         int64
	SB_size_struct_block         int64
	SB_first_free_bit_tree_dir   int64
	SB_first_free_bit_detail_dir int64
	SB_first_free_bit_table_dir  int64
	SB_first_free_bit_block      int64
	SB_magic_num                 int64
}

type ARBOLVIRTUALDIR struct {
	Avd_fecha_creacion              [20]byte
	Avd_nombre_directorio           [20]byte
	Avd_ap_array_subdirectorios     []int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      [20]byte
}

type FILE struct {
	DD_file_nombre            [20]byte
	DD_file_ap_inodo          int64
	DD_file_date_creacion     [20]byte
	DD_file_date_modificacion [20]byte
}
type DIRECTORYDETAIL struct {
	DD_array_files           [5]FILE
	DD_ap_detalle_directorio int64
}

type TABLEINODE struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           []int64
	I_ap_indirecto            int64
	I_id_proper               int64
}

type DATABLOCK struct {
	DB_data [50]byte
}

type LOG struct {
	Log_tipo_operacion int64
	log_tipo           int64
	log_nombre         [20]byte
	log_contenido      [50]byte
	log_fecha          [20]byte
}
