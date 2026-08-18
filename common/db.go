package common

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	dbPath string
	db     *sql.DB
}

func NewRepository(dbPath string) *Repository {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	return &Repository{
		dbPath: dbPath,
		db:     db,
	}
}

func (this *Repository) Init() {
	this.InitFirme()
	this.InitReprezentanti()
	this.InitStari()
	this.InitCaen()
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

	parsed = read_data("./data/od_caen_autorizat.csv")
	this.UpdateCaen(parsed)
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
			denumire TEXT NOT NULL,
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

		CREATE VIRTUAL TABLE IF NOT EXISTS firme_search USING fts5(
			denumire,
			content='firme',
			content_rowid='rowid',
			tokenize='trigram'
		);

		-- Insert
		CREATE TRIGGER IF NOT EXISTS firme_ai AFTER INSERT ON firme BEGIN
			INSERT INTO firme_search(rowid, denumire)
			VALUES (new.rowid, new.denumire);
		END;

		-- Delete
		CREATE TRIGGER IF NOT EXISTS firme_ad AFTER DELETE ON firme BEGIN
			INSERT INTO firme_search(firme_search, rowid, denumire)
			VALUES ('delete', old.rowid, old.denumire);
		END;

		-- Update
		CREATE TRIGGER IF NOT EXISTS firme_au AFTER UPDATE OF denumire ON firme BEGIN
			INSERT INTO firme_search(firme_search, rowid, denumire)
			VALUES ('delete', old.rowid, old.denumire);

			INSERT INTO firme_search(rowid, denumire)
			VALUES (new.rowid, new.denumire);
		END;
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateFirme(dataset []map[string]string) {
	stmt := `
		INSERT OR REPLACE INTO firme (denumire, cui, cod_inmatriculare, data_inmatriculare, euid, forma_juridica, tara, judet, localitate, strada, nr_strada, bloc, scara, etaj, apartament, cod_postal, sector, completare, web, tara_firma_mama)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

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
			fmt.Println(data["DENUMIRE"])
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
			cod_inmatriculare TEXT NOT NULL,
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
		CREATE INDEX IF NOT EXISTS idx_reprezentanti_cod_inmatriculare
		ON reprezentanti(cod_inmatriculare);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateReprezentanti(dataset []map[string]string) {
	stmt := `
		INSERT INTO reprezentanti (cod_inmatriculare, persoana_imputernicita, calitate, data_nastere, localitate_nastere, judet_nastere, tara_nastere, localitate, judet, tara)
		VALUES (?,?,?,?,?,?,?,?,?,?);`

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
		INSERT OR REPLACE INTO stari (cod_inmatriculare, cod, status)
		VALUES (?,?,?);`

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

func (this *Repository) InitCaen() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS caen (
			cod_inmatriculare TEXT NOT NULL,
			cod_caen INTEGER NOT NULL,
			versiune_caen INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_caen_cod_inmatriculare
		ON caen(cod_inmatriculare);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateCaen(dataset []map[string]string) {
	stmt := `
		INSERT INTO caen (cod_inmatriculare, cod_caen, versiune_caen)
		VALUES (?,?,?);`

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
		// some caen codes doesn't have any company registered
		if data["COD_INMATRICULARE"] == "" {
			continue
		}

		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["COD_CAEN_AUTORIZAT"], data["VER_CAEN_AUTORIZAT"])
		if err != nil {
			panic(err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

type InfoFirma struct {
	Nume string
	CodInmatriculare string
	FormaJuridica string
	Cui int
	Administratori []string
	DataInregistrare string
	Judet string
	Status string
	CoduriCaen []string
}

func (this *Repository) GetFirme(partialNumeFirma string, judetFilter string, statusFilter string) []*InfoFirma {
	stmt := `SELECT
				firme.denumire,
				firme.cod_inmatriculare,
				firme.forma_juridica,
				firme.cui,
				(
					SELECT GROUP_CONCAT(r.persoana_imputernicita, ',')
					FROM reprezentanti r
					WHERE r.cod_inmatriculare = firme.cod_inmatriculare
					  AND r.calitate = 'administrator'
				) AS persoane_imputernicite,
				firme.data_inmatriculare,
				firme.judet,
				stari.status,
				(
					SELECT GROUP_CONCAT(DISTINCT c.cod_caen)
					FROM caen c
					WHERE c.cod_inmatriculare = firme.cod_inmatriculare
				) AS coduri_caen
			FROM firme
			JOIN firme_search
				ON firme.rowid = firme_search.rowid
			JOIN stari
				ON firme.cod_inmatriculare = stari.cod_inmatriculare
			WHERE 1=1 `

	params := []any{}

	if partialNumeFirma != "" {
		stmt += " AND firme_search MATCH ? "
		params = append(params, partialNumeFirma)
	}

	if judetFilter != "" {
		stmt += " AND firme.judet = ? "
		params = append(params, judetFilter)
	}

	if statusFilter != "" {
		stmt += " AND stari.status = ? "
		params = append(params, statusFilter)
	}

	stmt += "\n"

	stmt += "LIMIT 20";

	fmt.Println(stmt)
	fmt.Println(params)

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Started searching in db...")
	rows, err := preparedStmt.Query(params...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	listaFirme := []*InfoFirma{}
	for rows.Next() {
		var infoFirma InfoFirma

		var administratori sql.NullString
		var coduriCaen sql.NullString

		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &administratori, &infoFirma.DataInregistrare, &infoFirma.Judet, &infoFirma.Status, &coduriCaen)
		if err != nil {
			panic(err)
		}

		if administratori.Valid {
			infoFirma.Administratori = strings.Split(administratori.String, ",") 
		}

		if coduriCaen.Valid{
			infoFirma.CoduriCaen = strings.Split(coduriCaen.String, ",")
		}

		listaFirme = append(listaFirme, &infoFirma)
	}

	fmt.Println("finished searching in db!")

	return listaFirme;
}

func (this *Repository) GetFirma(numar_inmatriculare string) *InfoFirma {
	stmt := `SELECT
				firme.denumire,
				firme.cod_inmatriculare,
				firme.forma_juridica,
				firme.cui,
				(
					SELECT GROUP_CONCAT(r.persoana_imputernicita, ',')
					FROM reprezentanti r
					WHERE r.cod_inmatriculare = firme.cod_inmatriculare
					  AND r.calitate = 'administrator'
				) AS persoane_imputernicite,
				firme.data_inmatriculare,
				firme.judet,
				stari.status,
				(
					SELECT GROUP_CONCAT(DISTINCT c.cod_caen)
					FROM caen c
					WHERE c.cod_inmatriculare = firme.cod_inmatriculare
				) AS coduri_caen
			FROM firme
			JOIN stari
				ON firme.cod_inmatriculare = stari.cod_inmatriculare
			WHERE firme.cod_inmatriculare = ?;`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	rows, err := preparedStmt.Query(numar_inmatriculare)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var infoFirma InfoFirma
	for rows.Next() {
		var administratori sql.NullString
		var coduriCaen sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &administratori, &infoFirma.DataInregistrare, &infoFirma.Judet, &infoFirma.Status, &coduriCaen)
		if err != nil {
			panic(err)
		}

		if administratori.Valid {
			infoFirma.Administratori = strings.Split(administratori.String, ",") 
		}

		if coduriCaen.Valid {
			infoFirma.CoduriCaen = strings.Split(coduriCaen.String, ",")
		}
	}

	return &infoFirma
}
