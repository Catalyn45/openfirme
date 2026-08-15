package common

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"slices"
)

type Repository struct {
	dbName string
	db     *sql.DB
}

func NewRepository() *Repository {
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

func (this *Repository) Update() {
	parsed := read_data("./data/od_firme.csv")
	this.UpdateFirme(parsed)

	parsed = read_data("./data/od_reprezentanti_legali.csv")
	this.UpdateReprezentanti(parsed)

	stare_firma := read_data("./data/od_stare_firma.csv")
	nomenclatura := read_data("./data/n_stare_firma.csv")

	resolve_nomenclatura(stare_firma, nomenclatura)

	this.UpdateStari(stare_firma)
}

func resolve_nomenclatura(dataset []map[string]string, nomenclatura []map[string]string) {
	for _, data := range dataset {
		index := slices.IndexFunc(nomenclatura, func (element map[string]string) bool { return data["COD"] == element["COD"]})
		data["STATUS"] = nomenclatura[index]["DENUMIRE"]
	}
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

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}
	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["DENUMIRE"], data["CUI"], data["COD_INMATRICULARE"], data["DATA_INMATRICULARE"], data["EUID"], data["FORMA_JURIDICA"], data["ADR_TARA"], data["ADR_JUDET"], data["ADR_LOCALITATE"], data["ADR_DEN_STRADA"], data["ADR_NR_STRADA"], data["ADR_BLOC"], data["ADR_SCARA"], data["ADR_ETAJ"], data["ADR_APARTAMENT"], data["ADR_COD_POSTAL"], data["ADR_SECTOR"], data["ADR_COMPLETARE"], data["WEB"], data["TARA_FIRMA_MAMA"])

		if err != nil {
			panic(err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
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

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["PERSOANA_IMPUTERNICITA"], data["CALITATE"], data["DATA_NASTERE"], data["LOCALITATE_NASTERE"], data["JUDET_NASTERE"], data["TARA_NASTERE"], data["LOCALITATE"], data["JUDET"], data["TARA"])

		if err != nil {
			panic(err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
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

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["COD"], data["STATUS"])

		if err != nil {
			panic(err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

type InfoFirmaLight struct {
	Nume string
}

func (this *Repository) GetFirme(partialNumeFirma string) []*InfoFirmaLight {
	stmt := `SELECT denumire from firme where denumire LIKE '%' || ? || '%';`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	rows, err := preparedStmt.Query(partialNumeFirma)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	listaNume := []*InfoFirmaLight{}
	for rows.Next() {
		var nume string

		err := rows.Scan(&nume)
		if err != nil {
			panic(err)
		}

		listaNume = append(listaNume, &InfoFirmaLight{Nume: nume})
	}

	return listaNume;
}

type InfoFirma struct {
	Nume string
	FormaJuridica string
	Cui int
	Administrator string
	DataInregistrare string
	Judet string
	Status string
}

func (this *Repository) GetFirma(nume_firma string) *InfoFirma {
	stmt := `SELECT
				firme.denumire,
				firme.forma_juridica,
				firme.cui,
				reprezentanti.persoana_imputernicita,
				firme.data_inmatriculare,
				firme.judet,
				stari.status
			from firme
			left join reprezentanti on firme.cod_inmatriculare = reprezentanti.cod_inmatriculare
									and reprezentanti.calitate = 'administrator'
			join stari on firme.cod_inmatriculare = stari.cod_inmatriculare
			where firme.denumire = ?`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	rows, err := preparedStmt.Query(nume_firma)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var infoFirma InfoFirma
	for rows.Next() {
		var administrator sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.FormaJuridica, &infoFirma.Cui, &administrator, &infoFirma.DataInregistrare, &infoFirma.Judet, &infoFirma.Status)
		if err != nil {
			panic(err)
		}

		if administrator.Valid {
			infoFirma.Administrator = administrator.String
		}
	}

	return &infoFirma
}
