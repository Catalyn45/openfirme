class TopFirmeSearchPage extends SearchPageBase {
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

    setMetadata() {
        if (this.county.value) {
            this.addMetadata(`OpenFirme România - Top companii din județul ${this.county.value}`, `Vezi topul companiilor din județul ${this.county.value} după cifra de afaceri, profit și numărul de angajați.`)
        } else if (this.formaJuridica.value) {
            this.addMetadata(`OpenFirme România - Top companii ${this.formaJuridica.value}`, `Vezi topul companiilor cu forma juridică ${this.formaJuridica.value} după cifra de afaceri, profit și numărul de angajați.`)
        } else if (this.domeniu.value) {
            let domeniuMapped = domeniiMap[this.domeniu.value]
            this.addMetadata(`OpenFirme România - Top companii din domeniul ${domeniuMapped}`, `Vezi topul companiilor din domeniul ${domeniuMapped} după cifra de afaceri, profit și numărul de angajați.`)
        } else {
            this.addMetadata("OpenFirme România - Top companii din România", "Vezi topul companiilor din România după cifra de afaceri, profit și numărul de angajați.")
        }
    }

    getFilters() {
        let params = super.getFilters()

        params.set('sort_by', this.sortBy.value)
        params.set('sort_order', this.sortOrder.value)

        return params
    }

	getLinkForPage(pageNumber) {
		return `/top/${pageNumber}?${this.getFilters()}`
	}

	getLinkForDataRequest() {
		return `/api/top/${this.pageNumber}?${this.getFilters()}`
	}

	getSearchTitle() {
		return `Top companii ordonate după: ${this.sortBy.options[this.sortBy.selectedIndex].text}`
	}
}

function CreateComponent() {
    return new TopFirmeSearchPage()
}
