package main

type CmdStruct struct {
	Command string `json:"command"`
}

type RespuestaJSON struct {
	Respuesta string `json:"respuesta"`
}

type MBRU struct {
	Mbr_tamano         int
	Mbr_fecha_creacion string
	Mbr_dsk_signature  int
	Dsk_fit            string
	Mbr_partition_1    PartitionU
	Mbr_partition_2    PartitionU
	Mbr_partition_3    PartitionU
	Mbr_partition_4    PartitionU
}

type PartitionU struct {
	Part_status int
	Part_type   string
	Part_fit    string
	Part_start  int
	Part_size   int
	Part_name   string
}

type EBR struct {
	Part_status int
	Part_fit    string
	Part_start  int
	Part_size   int
	Part_next   int
	Part_name   string
}
