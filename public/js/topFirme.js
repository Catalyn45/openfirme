class TopFirmeSearchPage extends SearchPage {
	initPrototype() {
		super.initPrototype()

		this.companyAngajati = this.prototypeCard.getElementsByClassName("company-angajati")[0]
		this.companyProfit = this.prototypeCard.getElementsByClassName("company-profit")[0]
		this.companyCifraAfaceri = this.prototypeCard.getElementsByClassName("company-cifra-afaceri")[0]
	}

    initFilters() {
        super.initFilters()

        this.sortBy = document.getElementById("sortField")
        this.sortOrder = document.getElementById("sortDirection")
    }

	populatePrototype(data) {
        super.populatePrototype(data)

        this.companyProfit.textContent = formatMoney(data.ProfitNet)
        this.companyCifraAfaceri.textContent = formatMoney(data.CifraAfaceri)
        this.companyAngajati.textContent = data.Angajati ?? 0
	}

    setFilters() {
        let params = super.setFilters()

        this.sortBy.value = params.get("sort_by") ?? this.sortBy.value
        this.sortOrder.value = params.get("sort_order") ?? this.sortOrder.value
    }

    getFilters() {
        let params = super.getFilters()

        params.set('sort_by', this.sortBy.value)
        params.set('sort_order', this.sortOrder.value)

        return params
    }

	getLinkForPage(pageNumber) {
		return `/top/${pageNumber}`
	}

	getLinkForDataRequest() {
		return "/topFirme"
	}

    initSearchBar() { }
    setSearchBar() { }
}

function CreateComponent() {
    return new TopFirmeSearchPage()
}
