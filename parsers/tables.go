package parsers

var createTableMap = map[string]string{
	"BPA": `CREATE TABLE IF NOT EXISTS bpa
	(
		"ID" PRIMARY KEY,
		"CNPJ_CIA" varchar(20),
		"DT_REFER" integer,
		"VERSAO" integer,
		"DENOM_CIA" varchar(100),
		"CD_CVM" integer,
		"GRUPO_DFP" varchar(206),
		"MOEDA" varchar(4),
		"ESCALA_MOEDA" varchar(7),
		"ORDEM_EXERC" varchar(9),
		"DT_FIM_EXERC" integer,
		"CD_CONTA" varchar(18),
		"DS_CONTA" varchar(100),
		"VL_CONTA" real
	);`,

	"BPP": `CREATE TABLE IF NOT EXISTS bpp
	(
		"ID" PRIMARY KEY,
		"CNPJ_CIA" varchar(20),
		"DT_REFER" integer,
		"VERSAO" integer,
		"DENOM_CIA" varchar(100),
		"CD_CVM" integer,
		"GRUPO_DFP" varchar(206),
		"MOEDA" varchar(4),
		"ESCALA_MOEDA" varchar(7),
		"ORDEM_EXERC" varchar(9),
		"DT_FIM_EXERC" integer,
		"CD_CONTA" varchar(18),
		"DS_CONTA" varchar(100),
		"VL_CONTA" real
	);`,
}
