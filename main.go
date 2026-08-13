package main

import (
	"fmt"
	"slices"
)

func resolve_nomenclatura(dataset []map[string]string, nomenclatura []map[string]string) {
	for _, data := range dataset {
		index := slices.IndexFunc(nomenclatura, func (element map[string]string) bool { return data["COD"] == element["COD"]})
		data["STATUS"] = nomenclatura[index]["DENUMIRE"]
	}
}

func main() {
	repository := newRepository()
	repository.Init()

	parsed := read_data("./data/od_firme.csv")
	repository.UpdateFirme(parsed)

	parsed = read_data("./data/od_reprezentanti_legali.csv")
	repository.UpdateReprezentanti(parsed)

	stare_firma := read_data("./data/od_stare_firma.csv")
	nomenclatura := read_data("./data/n_stare_firma.csv")

	resolve_nomenclatura(stare_firma, nomenclatura)

	repository.UpdateStari(stare_firma)

	fmt.Println("Finished")
}