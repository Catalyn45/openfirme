package common

import (
	"context"
	"database/sql"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	config *DBConfig
	db     *sql.DB
}

func NewRepository(dbPath string) *Repository {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	return &Repository{
		config: &config.DBConfig,
		db:     db,
	}
}

func (this *Repository) Init() {
	this.InitMetadata()
	this.InitFirme()
	this.InitReprezentanti()
	this.InitStari()
	this.InitCaen()
	this.InitDescriereCaen()
	this.InitDateIdentificare()
	this.InitBilanturi()
}

func (this *Repository) InitMetadata() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS metadata (
			tablename TEXT PRIMARY KEY,
			dataset TEXT NOT NULL
		);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) DeleteFromTable(transaction *sql.Tx, tableName string) {
	stmt := "DELETE FROM " + tableName + ";"
	_, err := transaction.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) InitFirme() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS firme (
			denumire TEXT NOT NULL,
			denumire_norm TEXT NOT NULL,
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

		CREATE INDEX IF NOT EXISTS idx_firme_cui
		ON firme(cui);

		CREATE INDEX IF NOT EXISTS idx_firme_judet
		ON firme(judet);

		CREATE INDEX IF NOT EXISTS idx_firme_forma_juridica
		ON firme(forma_juridica);

		CREATE INDEX IF NOT EXISTS idx_firme_forma_data_inmatriculare
		ON firme(data_inmatriculare);

		CREATE INDEX IF NOT EXISTS idx_firme_judet_forma_juridica
		ON firme(judet, forma_juridica);

		CREATE INDEX IF NOT EXISTS idx_firme_judet_data_inmatriculara
		ON firme(judet, data_inmatriculare);

		CREATE INDEX IF NOT EXISTS idx_firme_forma_juridica_data_inmatriculara
		ON firme(forma_juridica, data_inmatriculare);

		CREATE INDEX IF NOT EXISTS idx_firme_judet_forma_juridica_data_inmatriculara
		ON firme(judet, forma_juridica, data_inmatriculare);

		CREATE VIRTUAL TABLE IF NOT EXISTS firme_search USING fts5(
			denumire_norm,
			content='firme',
			content_rowid='rowid',
			tokenize='trigram'
		);

		-- Insert
		CREATE TRIGGER IF NOT EXISTS firme_ai AFTER INSERT ON firme BEGIN
			INSERT INTO firme_search(rowid, denumire_norm)
			VALUES (new.rowid, new.denumire_norm);
		END;

		-- Delete
		CREATE TRIGGER IF NOT EXISTS firme_ad AFTER DELETE ON firme BEGIN
			INSERT INTO firme_search(firme_search, rowid, denumire_norm)
			VALUES ('delete', old.rowid, old.denumire_norm);
		END;

		-- Update
		CREATE TRIGGER IF NOT EXISTS firme_au AFTER UPDATE OF denumire_norm ON firme BEGIN
			INSERT INTO firme_search(firme_search, rowid, denumire_norm)
			VALUES ('delete', old.rowid, old.denumire_norm);

			INSERT INTO firme_search(rowid, denumire_norm)
			VALUES (new.rowid, new.denumire_norm);
		END;
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateMetadata(transaction *sql.Tx, datasetName string, tableName string) {
	metadataStmt := `
		INSERT OR REPLACE INTO metadata (tablename, dataset)
		VALUES ('` + tableName + "' , '" + datasetName + "')"

	_, err := transaction.Exec(metadataStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) isOnDataset(datasetName string, table string) bool {
	metadataStmt := `
		SELECT dataset
		FROM metadata
		WHERE tablename = '` + table + "';"

	rows, err := this.db.Query(metadataStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var existingDataset string
	err = rows.Scan(&existingDataset)
	if err != nil {
		panic(err)
	}

	if existingDataset != datasetName {
		return false
	}

	return true
}

func (this *Repository) IsFirmeOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "firme")
}

func (this *Repository) UpdateFirme(dataset []map[string]string, datasetName string) {
	log.Println("Updating firme")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "firme")

	stmt := `
		INSERT OR REPLACE INTO firme (denumire, denumire_norm, cui, cod_inmatriculare, data_inmatriculare, euid, forma_juridica, tara, judet, localitate, strada, nr_strada, bloc, scara, etaj, apartament, cod_postal, sector, completare, web, tara_firma_mama)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?, ?);`

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}
	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["DENUMIRE"], normalize(data["DENUMIRE"]), data["CUI"], data["COD_INMATRICULARE"], convertDate(data["DATA_INMATRICULARE"]), data["EUID"], data["FORMA_JURIDICA"], data["ADR_TARA"], data["ADR_JUDET"], data["ADR_LOCALITATE"], data["ADR_DEN_STRADA"], data["ADR_NR_STRADA"], data["ADR_BLOC"], data["ADR_SCARA"], data["ADR_ETAJ"], data["ADR_APARTAMENT"], data["ADR_COD_POSTAL"], data["ADR_SECTOR"], data["ADR_COMPLETARE"], data["WEB"], data["TARA_FIRMA_MAMA"])
		if err != nil {
			log.Println(data["DENUMIRE"])
			panic(err)
		}
	}

	this.UpdateMetadata(transaction, datasetName, "firme")

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
			persoana_imputernicita_norm TEXT NOT NULL,
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

		CREATE INDEX IF NOT EXISTS idx_reprezentanti_persoana_imputernicita
		ON reprezentanti(persoana_imputernicita_norm, cod_inmatriculare);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) IsReprezentantiOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "reprezentanti")
}

func (this *Repository) normalizeNumeReprezentant(nume string) string {
	splitted := strings.Fields(nume)

	result := []string{}
	for _, item := range splitted {
		if strings.HasSuffix(item, ".")  {
			continue
		}

		result = append(result, strings.ToLower(item))
	}

	return strings.Join(result, " ")
}

func (this *Repository) UpdateReprezentanti(dataset []map[string]string, datasetName string) {
	log.Println("Updating reprezentanti")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "reprezentanti")

	stmt := `
		INSERT INTO reprezentanti (cod_inmatriculare, persoana_imputernicita, persoana_imputernicita_norm, calitate, data_nastere, localitate_nastere, judet_nastere, tara_nastere, localitate, judet, tara)
		VALUES (?,?,?,?,?,?,?,?,?,?,?);`

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["COD_INMATRICULARE"], data["PERSOANA_IMPUTERNICITA"], this.normalizeNumeReprezentant(data["PERSOANA_IMPUTERNICITA"]), data["CALITATE"], strings.Split(data["DATA_NASTERE"], " ")[0], data["LOCALITATE_NASTERE"], strings.ToUpper(data["JUDET_NASTERE"]), data["TARA_NASTERE"], data["LOCALITATE"], data["JUDET"], data["TARA"])
		if err != nil {
			panic(err)
		}
	}

	this.UpdateMetadata(transaction, datasetName, "reprezentanti")

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

func (this *Repository) InitStari() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS stari (
			cod_inmatriculare TEXT NOT NULL,
			cod INTEGER NOT NULL,
			status TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_stari_cod_status
		ON stari(cod_inmatriculare, status);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) IsStariOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "stari")
}

func (this *Repository) UpdateStari(dataset []map[string]string, datasetName string) {
	log.Println("Updating stari")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "stari")

	stmt := `
		INSERT INTO stari (cod_inmatriculare, cod, status)
		VALUES (?,?,?);`

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

	this.UpdateMetadata(transaction, datasetName, "stari")

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

func (this *Repository) IsCaenOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "caen")
}

func (this *Repository) UpdateCaen(dataset []map[string]string, datasetName string) {
	log.Println("Updating caen")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "caen")

	stmt := `
		INSERT INTO caen (cod_inmatriculare, cod_caen, versiune_caen)
		VALUES (?,?,?);`

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

	this.UpdateMetadata(transaction, datasetName, "caen")

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

func (this *Repository) InitDescriereCaen() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS descriere_caen (
			sectiunea TEXT,
			subsectiunea TEXT,
			diviziunea INTEGER,
			grupa INTEGER,
			clasa INTEGER,
			denumire TEXT,
			versiune_caen TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_descriere_caen_clasa_versiune
		ON descriere_caen(versiune_caen, clasa);

		CREATE INDEX IF NOT EXISTS idx_descriere_caen_sectiunea_clasa_versiune
		ON descriere_caen(sectiunea, versiune_caen, clasa);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) IsDescriereCaenOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "descriere_caen")
}

func (this *Repository) UpdateDescriereCaen(dataset []map[string]string, datasetName string) {
	log.Println("Updating descriere_caen")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "descriere_caen")

	stmt := `
		INSERT INTO descriere_caen (sectiunea, subsectiunea, diviziunea, grupa, clasa, denumire, versiune_caen)
		VALUES (?,?,?,?,?,?,?);`

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["SECTIUNEA"], data["SUBSECTIUNEA"], data["DIVIZIUNEA"], data["GRUPA"], data["CLASA"], data["DENUMIRE"], data["VERSIUNE_CAEN"])
		if err != nil {
			panic(err)
		}
	}

	this.UpdateMetadata(transaction, datasetName, "descriere_caen")

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

func (this *Repository) InitDateIdentificare() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS dateidentificare (
			cui INTEGER NOT NULL,
			tva INTEGER NOT NULL,
			impozitare_profit INTEGER NOT NULL,
			impozitare_venit INTEGER NOT NULL,
			data_stare TEXT NOT NULL,
			stare TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_dateidentificare_cui
		ON dateidentificare(cui);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) IsDateIdentificareOnDataset(dataset string) bool {
	return this.isOnDataset(dataset, "dateidentificare")
}

func (this *Repository) UpdateDateIdentificare(dataset []map[string]string, datasetName string) {
	log.Println("Updating dateidentificare")

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteFromTable(transaction, "dateidentificare")

	stmt := `
		INSERT INTO dateidentificare (cui, tva, impozitare_profit, impozitare_venit, data_stare, stare)
		VALUES (?,?,?,?,?,?);`

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		if data["COD_FISCAL"] == "" {
			continue
		}

		_, err = preparedStmt.Exec(data["COD_FISCAL"], data["TVA"] == "DA", data["IMP100"] == "DA", data["IMP120"] == "DA", data["DATA_STARE"], data["STARE"])
		if err != nil {
			panic(err)
		}
	}

	this.UpdateMetadata(transaction, datasetName, "dateidentificare")

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

func (this *Repository) InitBilanturi() {
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS bilanturi (
			cui INTEGER NOT NULL,
			cod_caen INTEGER,
			active_imobilizate INTEGER,
			active_circulante INTEGER,
			stocuri INTEGER,
			creante INTEGER,
			casa_si_conturi INTEGER,
			cheltuieli_avans INTEGER,
			datorii INTEGER,
			venituri_avans INTEGER,
			provizioane INTEGER,
			capitaluri INTEGER,
			capital_subscris INTEGER,
			cifra_afaceri INTEGER,
			venituri INTEGER,
			cheltuieli INTEGER,
			profit_brut INTEGER,
			profit_net INTEGER,
			numar_mediu_salariati INTEGER,
			an INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_bilanturi_cui_an
		ON bilanturi(cui, an);

		CREATE INDEX IF NOT EXISTS idx_bilanuri_an_cifra_afaceri_cui
		ON bilanturi(an, cifra_afaceri DESC, cui);

		CREATE INDEX IF NOT EXISTS idx_bilanuri_an_profit_cui
		ON bilanturi(an, profit_net DESC, cui);

		CREATE INDEX IF NOT EXISTS idx_bilanuri_an_active_imobilizate_cui
		ON bilanturi(an, active_imobilizate DESC, cui);

		CREATE INDEX IF NOT EXISTS idx_bilanuri_an_angajati_cui
		ON bilanturi(an, numar_mediu_salariati DESC, cui);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) DoesAnExist(an int) bool {
	stmt := `
		SELECT 1
		FROM bilanturi
		WHERE bilanturi.an = ` + strconv.Itoa(an) + `
		LIMIT 1;`

	rows, err :=  this.db.Query(stmt)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	return rows.Next()
}

func (this *Repository) DeleteAnFromBilanturi(transaction *sql.Tx, an int) {
	stmt := "DELETE FROM bilanturi WHERE an = " + strconv.Itoa(an) + ";"
	_, err := transaction.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateBilanturi(dataset []map[string]int, an int) {
	log.Println("Updating bilanturi an: ", an)

	transaction, err := this.db.Begin()
	if err != nil {
		panic(err)
	}

	defer transaction.Rollback()

	this.DeleteAnFromBilanturi(transaction, an)

	stmt := `
		INSERT INTO bilanturi (cui, cod_caen, active_imobilizate, active_circulante, stocuri, creante, casa_si_conturi, cheltuieli_avans, datorii, venituri_avans, provizioane, capitaluri, capital_subscris, cifra_afaceri, venituri, cheltuieli, profit_brut, profit_net, numar_mediu_salariati, an)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	preparedStmt, err := transaction.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	defer preparedStmt.Close()

	for _, data := range dataset {
		_, err = preparedStmt.Exec(data["CUI"], data["CAEN"], data["I1"], data["I2"], data["I3"], data["I4"], data["I5"], data["I6"], data["I7"], data["I8"], data["I9"], data["I10"], data["I11"], data["I12"], data["I13"], data["I14"], data["I15"] - data["I16"], data["I17"] - data["I18"], data["I19"], an)
		if err != nil {
			panic(err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}
}

type FirmeFilters struct {
	cui int
	numePartial string
	judet string
	status string
	formaJuridica string
	domeniu string
}

type FirmeOrdering struct {
	sortBy string
	sortOrder string
}

type InfoFirmaLight struct {
	Nume string
	CodInmatriculare string
	FormaJuridica string
	Cui int
	DataInregistrare string
	Judet string
	Statusuri []string
	ProfitNet *int
	CifraAfaceri *int
	Angajati *int
}

type InfoFirmeResult struct {
	Count int
	Data []*InfoFirmaLight
}

func (this *Repository) mapSortFiled(sortBy string) string {
	if sortBy == "profit" {
		return "bilanturi.profit_net"
	}

	if sortBy == "angajati" {
		return "bilanturi.numar_mediu_salariati"
	}

	if sortBy == "cifra_afaceri" {
		return "bilanturi.cifra_afaceri"
	}

	if sortBy == "rank" {
		return "bm25(firme_search)"
	}

	return ""
}

func (this *Repository) encodeFts5(word string) string {
	return `"` + strings.ReplaceAll(word, `"`, `""`) + `"`
}

func (this *Repository) processNumePartial(numePartial string) string {
	words := strings.Fields(numePartial)

	for i, word := range words {
		words[i] = this.encodeFts5(word)

		if slices.Contains(allFormeJuridice, strings.ToUpper(word)) {
			wordWithPoints := strings.Join(strings.Split(word, ""), ".")
			wordWithPoints = this.encodeFts5(wordWithPoints)

			words[i] = "(" + words[i] + " OR " + wordWithPoints  + ")"
		}
	}

	ftsQuery := strings.Join(words, " AND ")

	return ftsQuery
}

func (this *Repository) addFiltersToQuery(stmt string, filters *FirmeFilters, params *[]any) string {
	if filters != nil {
		if filters.numePartial != "" {

			stmt += " AND firme_search MATCH ? "
			*params = append(*params, this.processNumePartial(filters.numePartial))
		}

		if filters.cui != 0 {
			stmt += " AND firme.cui = ? "
			*params = append(*params, filters.cui)
		}

		if filters.judet != "" {
			stmt += " AND firme.judet = ? "
			*params = append(*params, filters.judet)
		}

		if filters.formaJuridica != "" {
			stmt += " AND firme.forma_juridica = ? "
			*params = append(*params, filters.formaJuridica)
		}

		if filters.status != "" {
			stmt += ` AND EXISTS(
				SELECT 1
				FROM stari s
				WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				AND s.status = ?
			) `
			*params = append(*params, filters.status)
		}

		if filters.domeniu != "" {
			stmt += ` AND EXISTS(
				SELECT 1
				FROM descriere_caen
				WHERE descriere_caen.clasa = bilanturi.cod_caen
				AND descriere_caen.versiune_caen = 3
				AND descriere_caen.sectiunea = ?
			) `
			*params = append(*params, filters.domeniu)
		}

		stmt += "\n"
	}

	return stmt
}

func (this *Repository) addOrderingToQuery(stmt string, ordering *FirmeOrdering) string {
	sortBy := this.mapSortFiled(ordering.sortBy)
	if sortBy != "" {
		stmt += ` ORDER BY ` + sortBy + " " + ordering.sortOrder + "\n"
	}

	return stmt
}

func (this *Repository) addLimitToQuery(stmt string, pageNumber int, params *[]any) string {
	if pageNumber != 0 {
		stmt += "LIMIT 20 OFFSET ?"
		*params = append(*params, (pageNumber - 1) * 20)

		stmt += "\n"
	} else {
		stmt += "LIMIT 200"
		stmt += "\n"
	}

	return stmt
}

func (this *Repository) constructPagedQuery(stmt string, params *[]any, filters *FirmeFilters, ordering *FirmeOrdering, pageNumber int) string {
	if filters != nil {
		stmt = this.addFiltersToQuery(stmt, filters, params)
	}
	
	if ordering != nil {
		stmt = this.addOrderingToQuery(stmt, ordering)
	}

	stmt = this.addLimitToQuery(stmt, pageNumber, params)

	return stmt
}

type QueryOptions struct {
	timeout *time.Duration
}

type QueryResult struct {
	rows *sql.Rows
	cancel context.CancelFunc
}

func (this *QueryResult) Close() {
	this.rows.Close()
	this.cancel()
}

func (this *Repository) executeQuery(stmt string, params []any, readRowCallback func(*sql.Rows), opt *QueryOptions) {
	stmt += ";"

	if this.config.LogQueries {
		log.Println("stmt ", stmt)
		log.Println("params ", params)
	}

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	log.Println("Started searching in db...")

	var timeout time.Duration

	if opt != nil && opt.timeout != nil {
		timeout = *opt.timeout
	} else {
		timeout = time.Duration(this.config.DefaultTimeoutInSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := preparedStmt.QueryContext(ctx, params...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	log.Println("finished searching in db...")

	for rows.Next() {
		readRowCallback(rows)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Println("finished scanning data...")
}

func (this *Repository) executePagedQuery(stmt string, params []any, filters *FirmeFilters, ordering *FirmeOrdering, pageNumber int, readRowCallback func(*sql.Rows), opt *QueryOptions) {
	stmt = this.constructPagedQuery(stmt, &params, filters, ordering, pageNumber)
	this.executeQuery(stmt, params, readRowCallback, opt)
}

func (this *Repository) getCount(stmt string, params []any, filters *FirmeFilters, opt *QueryOptions) int {
	stmt = this.constructPagedQuery(stmt, &params, filters, nil, 0)
	stmt = "SELECT COUNT(*) FROM ( " + stmt + " )"

	count := 0
	this.executeQuery(stmt, params, func(rows *sql.Rows) {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}, opt)

	return count
}

func (this *Repository) GetFirme(filters *FirmeFilters, pageNumber int) *InfoFirmeResult {
	stmt := `SELECT
				firme.denumire,
				firme.cod_inmatriculare,
				firme.forma_juridica,
				firme.cui,
				firme.data_inmatriculare,
				firme.judet,
				(
					SELECT GROUP_CONCAT(s.status, '^')
					FROM stari s
					WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				) AS statuses
			FROM firme
			JOIN firme_search
				ON firme.rowid = firme_search.rowid
			`

	if filters.domeniu != "" {
		stmt += `
			LEFT JOIN bilanturi
				ON bilanturi.an = (
					SELECT MAX(b.an)
					FROM bilanturi b
				)
				AND firme.cui = bilanturi.cui
		`
	}

	stmt += `
	WHERE 1=1 `

	var ordering *FirmeOrdering = nil

	if filters.numePartial != "" {
		ordering = &FirmeOrdering {
			sortBy: "rank",
			sortOrder: "asc",
		}
	}
	timeout := time.Duration(this.config.SearchFirmeTimeoutInSeconds) * time.Second
	opt := &QueryOptions {
		timeout: &timeout,
	}

	result := &InfoFirmeResult {
		Count: this.getCount(stmt, nil, filters, opt),
		Data: []*InfoFirmaLight{},
	}

	rowCallback := func (rows *sql.Rows) {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, "^")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	this.executePagedQuery(
		stmt,
		nil,
		filters,
		ordering,
		pageNumber,
		rowCallback,
		opt)

	return result;
}

var allowedFormeJuridiceForBilanturi = []string{
	"ALT",
	"CA",
	"GEIE",
	"GIE",
	"INCD",
	"N/A",
	"OC2",
	"OCC",
	"OCM",
	"OCR",
	"RA",
	"SA",
	"SC",
	"SCA",
	"SCE",
	"SCS",
	"SE",
	"SNC",
	"SRL",
}

func (this *Repository) GetTopFirme(filters *FirmeFilters, ordering *FirmeOrdering, pageNumber int) *InfoFirmeResult {
	if filters.formaJuridica != "" &&
		!slices.Contains(allowedFormeJuridiceForBilanturi, filters.formaJuridica) {
		return &InfoFirmeResult {
			Count: 0,
			Data: []*InfoFirmaLight{},
		}
	}

	stmt := `
			SELECT
			firme.denumire,
			firme.cod_inmatriculare,
			firme.forma_juridica,
			firme.cui,
			firme.data_inmatriculare,
			firme.judet,
			(
				SELECT GROUP_CONCAT(s.status, '^')
				FROM stari s
				WHERE s.cod_inmatriculare = firme.cod_inmatriculare
			) AS statuses,
			bilanturi.cifra_afaceri,
			bilanturi.profit_net,
			bilanturi.numar_mediu_salariati
			FROM bilanturi
			RIGHT JOIN firme
				ON firme.cui = bilanturi.cui
			WHERE bilanturi.an = (
				SELECT MAX(b.an)
				FROM bilanturi b
			)
		`

	result := &InfoFirmeResult {
		Count: this.getCount(stmt, nil, filters, nil),
		Data: []*InfoFirmaLight{},
	}

	rowCallback := func (rows *sql.Rows) {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses, &infoFirma.CifraAfaceri, &infoFirma.ProfitNet, &infoFirma.Angajati)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, "^")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	this.executePagedQuery(
		stmt,
		nil,
		filters,
		ordering,
		pageNumber,
		rowCallback,
		nil)

	return result;
}

type InfoFirma struct {
	UpdatedDate string
	Nume string
	CodInmatriculare string
	Euid string
	FormaJuridica string
	Cui int
	Reprezentanti []string
	DataInregistrare string
	Judet string
	Localitate string
	Strada *string
	NrStrada *string
	Bloc *string
	Scara *string
	Etaj *string
	Apartament *string
	CodPostal *string
	Sector *string
	Statusuri []string
	CoduriCaen []string
	Tva *bool
	ImpozitareVenit *bool
	ImpozitareProfit *bool
	BilanturiFirma []*BilantFirma
}

type BilantFirma struct {
	An int
	CifraAfaceri int
	ProfitNet int
	Datorii int
	ActiveImobilizate int
	ActiveCirculante int
	Capitaluri int
	Angajati int
	Caen int
}

func (this *Repository) getInfoFirma(numar_inmatriculare string) *InfoFirma {
	stmt := `SELECT
				(
					SELECT dataset
					FROM METADATA
					WHERE tablename = 'firme'
				) as dataset,
				firme.denumire,
				firme.cod_inmatriculare,
				firme.euid,
				firme.forma_juridica,
				firme.cui,
				(
					SELECT GROUP_CONCAT(r.persoana_imputernicita || '^' || r.calitate, '@')
					FROM reprezentanti r
					WHERE r.cod_inmatriculare = firme.cod_inmatriculare
				) AS persoane_imputernicite,
				firme.data_inmatriculare,
				firme.judet,
				firme.localitate,
				firme.strada,
				firme.nr_strada,
				firme.bloc,
				firme.scara,
				firme.etaj,
				firme.apartament,
				firme.cod_postal,
				firme.sector,
				(
					SELECT GROUP_CONCAT(s.status, '^')
					FROM stari s
					WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				) AS statuses,
				(
					SELECT GROUP_CONCAT(c.cod_caen || '^' || dc.denumire, '@')
					FROM caen c
					LEFT JOIN descriere_caen dc
					ON c.cod_caen = dc.clasa AND c.versiune_caen = dc.versiune_caen
					WHERE c.cod_inmatriculare = firme.cod_inmatriculare
				) AS coduri_caen,
				dateidentificare.tva,
				dateidentificare.impozitare_profit,
				dateidentificare.impozitare_venit
			FROM firme
			LEFT JOIN dateidentificare
				ON firme.cui = dateidentificare.cui
			WHERE firme.cod_inmatriculare = ?`

	var infoFirma InfoFirma
	rowCallback := func (rows *sql.Rows) {
		var reprezentanti sql.NullString
		var coduriCaen sql.NullString
		var statuses sql.NullString
		var dataset string

		err := rows.Scan(&dataset, &infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.Euid, &infoFirma.FormaJuridica, &infoFirma.Cui, &reprezentanti, &infoFirma.DataInregistrare, &infoFirma.Judet, &infoFirma.Localitate, &infoFirma.Strada, &infoFirma.NrStrada, &infoFirma.Bloc, &infoFirma.Scara, &infoFirma.Etaj, &infoFirma.Apartament, &infoFirma.CodPostal, &infoFirma.Sector, &statuses, &coduriCaen, &infoFirma.Tva, &infoFirma.ImpozitareProfit, &infoFirma.ImpozitareVenit)
		if err != nil {
			panic(err)
		}

		infoFirma.UpdatedDate = dataset[6:16]

		if reprezentanti.Valid {
			infoFirma.Reprezentanti = strings.Split(reprezentanti.String, "@")
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, "^")
		}

		if coduriCaen.Valid {
			infoFirma.CoduriCaen = strings.Split(coduriCaen.String, "@")
		}
	}

	this.executeQuery(
		stmt,
		[]any{ numar_inmatriculare },
		rowCallback,
		nil)

	return &infoFirma
}

func (this *Repository) getBilanturiFirma(cui int) []*BilantFirma {
	stmt := `SELECT
				an,
				cifra_afaceri,
				profit_net,
				datorii,
				active_imobilizate,
				active_circulante,
				capitaluri,
				numar_mediu_salariati,
				cod_caen
			FROM bilanturi
			WHERE cui = ?
			ORDER BY an desc`

	var bilanturiFirma []*BilantFirma
	rowCallback := func (rows *sql.Rows) {
		var bilantFirma BilantFirma
		err := rows.Scan(&bilantFirma.An, &bilantFirma.CifraAfaceri, &bilantFirma.ProfitNet, &bilantFirma.Datorii, &bilantFirma.ActiveImobilizate, &bilantFirma.ActiveCirculante, &bilantFirma.Capitaluri, &bilantFirma.Angajati, &bilantFirma.Caen)
		if err != nil {
			panic(err)
		}

		bilanturiFirma = append(bilanturiFirma, &bilantFirma)
	}

	this.executeQuery(
		stmt,
		[]any { cui },
		rowCallback,
		nil)

	return bilanturiFirma
}

func (this *Repository) GetFirma(numar_inmatriculare string) *InfoFirma {
	infoFirma := this.getInfoFirma(numar_inmatriculare)
	
	if slices.Contains(allowedFormeJuridiceForBilanturi, infoFirma.FormaJuridica) {
		infoFirma.BilanturiFirma = this.getBilanturiFirma(infoFirma.Cui)
	}

	return infoFirma
}

func (this *Repository) GetAdminFirme(cod_inmatriculare string, admin string, filters *FirmeFilters, pageNumber int) *InfoFirmeResult {
	stmt := `
		SELECT
			firme.denumire,
			firme.cod_inmatriculare,
			firme.forma_juridica,
			firme.cui,
			firme.data_inmatriculare,
			firme.judet,
			(
				SELECT GROUP_CONCAT(s.status, '^')
				FROM stari s
				WHERE s.cod_inmatriculare = firme.cod_inmatriculare
			) AS statuses
		FROM firme
		JOIN reprezentanti
			ON reprezentanti.cod_inmatriculare = firme.cod_inmatriculare
		JOIN (
			SELECT *
			FROM reprezentanti
			WHERE cod_inmatriculare = ? AND persoana_imputernicita = ?
			LIMIT 1
		) r ON reprezentanti.persoana_imputernicita_norm = r.persoana_imputernicita_norm
			AND reprezentanti.data_nastere = r.data_nastere
			AND reprezentanti.judet_nastere = r.judet_nastere
		WHERE 1=1
	`

	params := []any{ cod_inmatriculare, admin }

	result := &InfoFirmeResult {
		Count: this.getCount(stmt, params, filters, nil),
		Data: []*InfoFirmaLight{},
	}

	rowCallback := func (rows *sql.Rows) {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, "^")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	this.executePagedQuery(
		stmt,
		params,
		filters,
		nil,
		pageNumber,
		rowCallback,
		nil)

	return result;
}

func (this *Repository) GetNumeFirma(numarInmatriculare string) string {
	stmt := `SELECT
				firme.denumire
			FROM firme
			WHERE firme.cod_inmatriculare = ?`

	nume := ""
	rowCallback := func (rows *sql.Rows) {
		err := rows.Scan(&nume)
		if err != nil {
			panic(err)
		}
	}

	this.executeQuery(
		stmt,
		[]any { numarInmatriculare },
		rowCallback,
		nil)

	return nume
}
