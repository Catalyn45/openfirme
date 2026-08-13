package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	dbName string
	db     *sql.DB
}

func newRepository() *Repository {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}

	return &Repository{
		dbName: "foo.db",
		db:     db,
	}
}

func (this *Repository) Init() {
	this.InitFirme()
	this.InitReprezentanti()
	this.InitStari()
}

func (this *Repository) InitFirme() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS firme (
			denumire TEXT NOT NULL UNIQUE,
			cui INTEGER NOT NULL,
			cod_inmatriculare TEXT PRIMARY KEY,
			data_inmatriculare TEXT NOT NULL,
			euid TEXT,
			forma_juridica TEXT NOT NULL,
			tara TEXT,
			judet TEXT,
			localitate TEXT,
			strada TEXT,
			nr_strada TEXT,
			bloc TEXT,
			scara TEXT,
			etaj TEXT,
			apartament TEXT,
			cod_postal TEXT,
			sector TEXT,
			completare TEXT,
			web TEXT,
			tara_firma_mama TEXT
		);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateFirme(dataset []map[string]string) {
	stmt := `
		INSERT OR REPLACE INTO firme (denumire, cui, cod_inmatriculare, data_inmatriculare, euid, forma_juridica, tara, judet, localitate, strada, nr_strada, bloc, scara, etaj, apartament, cod_postal, sector, completare, web, tara_firma_mama) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	for index, data := range dataset {
		_, err = preparedStmt.Exec(data["DENUMIRE"], data["CUI"], data["COD_INMATRICULARE"], data["DATA_INMATRICULARE"], data["EUID"], data["FORMA_JURIDICA"], data["ADR_TARA"], data["ADR_JUDET"], data["ADR_LOCALITATE"], data["ADR_DEN_STRADA"], data["ADR_NR_STRADA"], data["ADR_BLOC"], data["ADR_SCARA"], data["ADR_ETAJ"], data["ADR_APARTAMENT"], data["ADR_COD_POSTAL"], data["ADR_SECTOR"], data["ADR_COMPLETARE"], data["WEB"], data["TARA_FIRMA_MAMA"])

		if err != nil {
			panic(err)
		}


		if index > 200 {
			break
		}
	}
}

func (this *Repository) InitReprezentanti() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS reprezentanti (
			cod_inmatriculare TEXT PRIMARY KEY,
			persoana_imputernicita TEXT NOT NULL,
			calitate TEXT NOT NULL,
			data_nastere TEXT,
			localitate_nastere TEXT,
			judet_nastere TEXT,
			tara_nastere TEXT,
			localitate TEXT,
			judet TEXT,
			tara TEXT
		);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateReprezentanti(dataset []map[string]string) {
	stmt := `
		INSERT OR REPLACE INTO reprezentanti (cod_inmatriculare, persoana_imputernicita, calitate, data_nastere, localitate_nastere, judet_nastere, tara_nastere, localitate, judet, tara) VALUES (?,?,?,?,?,?,?,?,?,?);`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	for index, data := range dataset {
		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["PERSOANA_IMPUTERNICITA"], data["CALITATE"], data["DATA_NASTERE"], data["LOCALITATE_NASTERE"], data["JUDET_NASTERE"], data["TARA_NASTERE"], data["LOCALITATE"], data["JUDET"], data["TARA"])

		if err != nil {
			panic(err)
		}

		if index > 200 {
			break
		}
	}
}

func (this *Repository) InitStari() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS stari (
			cod_inmatriculare TEXT PRIMARY KEY,
			cod INTEGER NOT NULL,
			status TEXT NOT NULL
		);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateStari(dataset []map[string]string) {
	stmt := `
		INSERT OR REPLACE INTO stari (cod_inmatriculare, cod, status) VALUES (?,?,?);`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	for index, data := range dataset {
		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["COD"], data["STATUS"])

		if err != nil {
			panic(err)
		}

		if index > 200 {
			break
		}
	}
}
