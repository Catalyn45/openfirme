package common

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

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
	this.InitBilanturi()
}

func (this *Repository) convert_values_to_string(oldmaps []map[string]string) []map[string]int {
	newMaps := []map[string]int{}

	for _, m := range oldmaps {
		newMap := make(map[string]int)
		for key, value := range m {
			newValue := 0

			if value != "" {
				var err error
				newValue, err = strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
			}

			newMap[key] = newValue
		}

		newMaps = append(newMaps, newMap)
	}

	return newMaps
}

func (this *Repository) Update() {
	fmt.Println("Updating tables")

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

	for i := 2011; ; i++ {
		filePath := "./data/web_bl_bs_sl_an" + strconv.Itoa(i) + ".txt"
		_, err := os.Stat(filePath)
		if err != nil {
			break
		}

		skipIndex := -1
		if i <= 2015 {
			skipIndex = 14
		}

		parsed = read_data_delimiter(filePath, ",", skipIndex)
		intParsed := this.convert_values_to_string(parsed)

		this.UpdateBilanturi(intParsed, i)
	}
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

func convertDate(datetime string) string {
	if datetime == "" {
		return ""
	}

	layouts := []string{
		"02/01/2006",
		"02/01/2006 15:04",
		"02/01/2006 15:04:05",
	}

	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, datetime)
		if err == nil {
			return t.Format("2006-01-02 15:04")
		}
	}

	panic(err)
}

func (this *Repository) UpdateFirme(dataset []map[string]string) {
	fmt.Println("Updating firme")

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
		_, err = preparedStmt.Exec(data["DENUMIRE"], data["CUI"], data["COD_INMATRICULARE"], convertDate(data["DATA_INMATRICULARE"]), data["EUID"], data["FORMA_JURIDICA"], data["ADR_TARA"], data["ADR_JUDET"], data["ADR_LOCALITATE"], data["ADR_DEN_STRADA"], data["ADR_NR_STRADA"], data["ADR_BLOC"], data["ADR_SCARA"], data["ADR_ETAJ"], data["ADR_APARTAMENT"], data["ADR_COD_POSTAL"], data["ADR_SECTOR"], data["ADR_COMPLETARE"], data["WEB"], data["TARA_FIRMA_MAMA"])
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
		ON reprezentanti(cod_inmatriculare, calitate);

		CREATE INDEX IF NOT EXISTS idx_reprezentanti_persoana_imputernicita
		ON reprezentanti(persoana_imputernicita, calitate, cod_inmatriculare);
	`

	_, err := this.db.Exec(createTableStmt)
	if err != nil {
		panic(err)
	}
}

func (this *Repository) UpdateReprezentanti(dataset []map[string]string) {
	fmt.Println("Updating reprezentanti")

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

func (this *Repository) UpdateStari(dataset []map[string]string) {
	fmt.Println("Updating stari")

	stmt := `
		INSERT INTO stari (cod_inmatriculare, cod, status)
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
	fmt.Println("Updating caen")

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
			patrimoniul INTEGER,
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

func (this *Repository) UpdateBilanturi(dataset []map[string]int, an int) {
	fmt.Println("Updating bilanturi")

	stmt := `
		INSERT INTO bilanturi (cui, cod_caen, active_imobilizate, active_circulante, stocuri, creante, casa_si_conturi, cheltuieli_avans, datorii, venituri_avans, provizioane, capitaluri, capital_subscris, patrimoniul, cifra_afaceri, venituri, cheltuieli, profit_brut, profit_net, numar_mediu_salariati, an)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

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
		_, err = preparedStmt.Exec(data["CUI"], data["CAEN"], data["I1"], data["I2"], data["I3"], data["I4"], data["I5"], data["I6"], data["I7"], data["I8"], data["I9"], data["I10"], data["I11"], data["I12"], data["I13"], data["I14"], data["I15"], data["I16"] - data["I17"], data["I18"] - data["I19"], data["I20"], an)
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
	numePartial string
	judet string
	status string
	formaJuridica string
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
	if sortBy == "infiintare" {
		return "firme.data_inmatriculare"
	}

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

func (this *Repository) addFiltersToQuery(stmt string, filters *FirmeFilters, params *[]any) string {
	if filters != nil {
		if filters.numePartial != "" {
			stmt += " AND firme_search MATCH ? "
			*params = append(*params, filters.numePartial)
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
			stmt += " AND statuses LIKE '%' || ? || '%'"
			*params = append(*params, filters.status)
		}

		stmt += "\n"
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

func (this *Repository) executePagedQuery(stmt string, filters *FirmeFilters, ordering *FirmeOrdering, pageNumber int, params ...any) *sql.Rows {
	if filters != nil {
		stmt = this.addFiltersToQuery(stmt, filters, &params)
	}
	
	if ordering != nil {
		sortBy := this.mapSortFiled(ordering.sortBy)
		if sortBy != "" {
			stmt += ` ORDER BY ` + sortBy + " " + ordering.sortOrder + "\n"
		}
	}

	stmt = this.addLimitToQuery(stmt, pageNumber, &params)

	if pageNumber == 0 {
		stmt = "SELECT COUNT(*) FROM ( " + stmt + " )"
	}

	stmt += ";"

	fmt.Println("stmt ", stmt)
	fmt.Println("params ", params)

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Started searching in db...")
	rows, err := preparedStmt.Query(params...)
	if err != nil {
		panic(err)
	}

	return rows
}

func (this *Repository) getCount(stmt string, filters *FirmeFilters, params ...any) int {
	rows := this.executePagedQuery(stmt, filters, nil, 0, params...)
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}

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
					SELECT GROUP_CONCAT(s.status, ',')
					FROM stari s
					WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				) AS statuses
			FROM firme
			JOIN firme_search
				ON firme.rowid = firme_search.rowid
			WHERE 1=1 `

	result := &InfoFirmeResult {
		Count: this.getCount(stmt, filters),
		Data: []*InfoFirmaLight{},
	}

	ordering := FirmeOrdering {
		sortBy: "rank",
		sortOrder: "asc",
	}

	rows := this.executePagedQuery(stmt, filters, &ordering, pageNumber)
	defer rows.Close()

	for rows.Next() {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, ",")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	return result;
}

func (this *Repository) GetTopFirme(filters *FirmeFilters, ordering *FirmeOrdering, pageNumber int) *InfoFirmeResult {
	if filters.formaJuridica != "" && slices.Index(allowedFormeJuridiceForBilanturi, filters.formaJuridica) == -1 && ordering.sortBy != "infiintare" {
		return &InfoFirmeResult {
			Count: 0,
			Data: []*InfoFirmaLight{},
		}
	}

	fields := `
			SELECT
				firme.denumire,
				firme.cod_inmatriculare,
				firme.forma_juridica,
				firme.cui,
				firme.data_inmatriculare,
				firme.judet,
				(
					SELECT GROUP_CONCAT(s.status, ',')
					FROM stari s
					WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				) AS statuses,
		 		bilanturi.cifra_afaceri,
		 		bilanturi.profit_net,
		 		bilanturi.numar_mediu_salariati
			`

	dataInmatriculareOrderEfficientStmt := fields + `
			FROM firme
			LEFT JOIN bilanturi
				ON bilanturi.an = (
					SELECT MAX(b.an)
					FROM bilanturi b
				)
				AND firme.cui = bilanturi.cui
			WHERE 1=1
	`

	bilanturiEfficientStmt := fields + `
			FROM bilanturi
			RIGHT JOIN firme
				ON firme.cui = bilanturi.cui
			WHERE bilanturi.an = (
					SELECT MAX(b.an)
					FROM bilanturi b
				)
	`

	result := &InfoFirmeResult {
		Count: this.getCount(dataInmatriculareOrderEfficientStmt, filters),
		Data: []*InfoFirmaLight{},
	}

	var stmt string
	if ordering.sortBy == "infiintare" {
		stmt = dataInmatriculareOrderEfficientStmt
	} else {
		stmt = bilanturiEfficientStmt
	}

	rows := this.executePagedQuery(stmt, filters, ordering, pageNumber)
	defer rows.Close()

	for rows.Next() {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses, &infoFirma.CifraAfaceri, &infoFirma.ProfitNet, &infoFirma.Angajati)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, ",")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	return result;
}

type InfoFirma struct {
	Nume string
	CodInmatriculare string
	FormaJuridica string
	Cui int
	Administratori []string
	DataInregistrare string
	Judet string
	Localitate string
	Strada string
	NrStrada string
	Bloc string
	Scara string
	Etaj string
	Apartament string
	CodPostal string
	Sector string
	Statusuri []string
	CoduriCaen []string
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
}

func (this *Repository) getInfoFirma(numar_inmatriculare string) *InfoFirma {
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
					SELECT GROUP_CONCAT(s.status, ',')
					FROM stari s
					WHERE s.cod_inmatriculare = firme.cod_inmatriculare
				) AS statuses,
				(
					SELECT GROUP_CONCAT(DISTINCT c.cod_caen)
					FROM caen c
					WHERE c.cod_inmatriculare = firme.cod_inmatriculare
				) AS coduri_caen
			FROM firme
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
		var strada sql.NullString
		var nrStrada sql.NullString
		var bloc sql.NullString
		var scara sql.NullString
		var etaj sql.NullString
		var apartament sql.NullString
		var codPostal sql.NullString
		var sector sql.NullString
		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &administratori, &infoFirma.DataInregistrare, &infoFirma.Judet, &infoFirma.Localitate, &strada, &nrStrada, &bloc, &scara, &etaj, &apartament, &codPostal, &sector, &statuses, &coduriCaen)
		if err != nil {
			panic(err)
		}

		if administratori.Valid {
			infoFirma.Administratori = strings.Split(administratori.String, ",") 
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, ",") 
		}

		if coduriCaen.Valid {
			infoFirma.CoduriCaen = strings.Split(coduriCaen.String, ",")
		}

		if strada.Valid {
			infoFirma.Strada = strada.String
		}

		if nrStrada.Valid {
			infoFirma.NrStrada = nrStrada.String
		}

		if bloc.Valid {
			infoFirma.Bloc = bloc.String
		}

		if scara.Valid {
			infoFirma.Scara = scara.String
		}

		if etaj.Valid {
			infoFirma.Etaj = etaj.String
		}

		if apartament.Valid {
			infoFirma.Apartament = apartament.String
		}

		if codPostal.Valid {
			infoFirma.CodPostal = codPostal.String
		}

		if sector.Valid {
			infoFirma.Sector = sector.String
		}
	}

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
				numar_mediu_salariati
			FROM bilanturi
			WHERE cui = ?
			ORDER BY an desc;`

	preparedStmt, err := this.db.Prepare(stmt)
	if err != nil {
		panic(err)
	}

	rows, err := preparedStmt.Query(cui)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var bilanturiFirma []*BilantFirma
	for rows.Next() {
		var bilantFirma BilantFirma
		err := rows.Scan(&bilantFirma.An, &bilantFirma.CifraAfaceri, &bilantFirma.ProfitNet, &bilantFirma.Datorii, &bilantFirma.ActiveImobilizate, &bilantFirma.ActiveCirculante, &bilantFirma.Capitaluri, &bilantFirma.Angajati)
		if err != nil {
			panic(err)
		}

		bilanturiFirma = append(bilanturiFirma, &bilantFirma)
	}

	return bilanturiFirma
}

func (this *Repository) GetFirma(numar_inmatriculare string) *InfoFirma {
	infoFirma := this.getInfoFirma(numar_inmatriculare)
	
	if slices.Index(allowedFormeJuridiceForBilanturi, infoFirma.FormaJuridica) != -1 {
		infoFirma.BilanturiFirma = this.getBilanturiFirma(infoFirma.Cui)
	}

	return infoFirma
}

func (this *Repository) GetAdminFirme(cod_inmatriculare string, admin string, pageNumber int) *InfoFirmeResult {
	stmt := `
		SELECT
			firme.denumire,
			firme.cod_inmatriculare,
			firme.forma_juridica,
			firme.cui,
			firme.data_inmatriculare,
			firme.judet,
			(
				SELECT GROUP_CONCAT(s.status, ',')
				FROM stari s
				WHERE s.cod_inmatriculare = firme.cod_inmatriculare
			) AS statuses
		FROM firme
		JOIN reprezentanti
			ON reprezentanti.cod_inmatriculare = firme.cod_inmatriculare
		JOIN (
			SELECT *
			FROM reprezentanti
			WHERE cod_inmatriculare = ? AND persoana_imputernicita = ? and calitate = 'administrator'
			LIMIT 1
		) r ON reprezentanti.persoana_imputernicita = r.persoana_imputernicita
			AND reprezentanti.calitate = r.calitate
			AND reprezentanti.data_nastere = r.data_nastere
			AND reprezentanti.judet_nastere = r.judet_nastere
	`

	result := &InfoFirmeResult {
		Count: this.getCount(stmt, nil, cod_inmatriculare, admin),
		Data: []*InfoFirmaLight{},
	}

	rows := this.executePagedQuery(stmt, nil, nil, pageNumber, cod_inmatriculare, admin, )
	defer rows.Close()

	for rows.Next() {
		var infoFirma InfoFirmaLight

		var statuses sql.NullString
		err := rows.Scan(&infoFirma.Nume, &infoFirma.CodInmatriculare, &infoFirma.FormaJuridica, &infoFirma.Cui, &infoFirma.DataInregistrare, &infoFirma.Judet, &statuses)
		if err != nil {
			panic(err)
		}

		if statuses.Valid {
			infoFirma.Statusuri = strings.Split(statuses.String, ",")
		}

		result.Data = append(result.Data, &infoFirma)
	}

	return result;
}
