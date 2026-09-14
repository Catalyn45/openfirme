class DosareSearchPage extends SearchPageBase {
    init() {
		const path = window.location.pathname.split("/")
        this.inregistrare = path.at(-2)

        this.numeFirma = decodeURIComponent(new URLSearchParams(window.location.search).get("nume_firma"))

		this.initSearchBar()
		this.initFilters()
		this.initPrototype()
		this.initPages()
    }

    initFilters() {
        this.tribunal = document.getElementById("tribunalFilter")
        this.categorie = document.getElementById("categorieFilter")
        this.stadiuProcesual = document.getElementById("stadiuProcesualFilter")

        this.sortBy = document.getElementById("sortField")
        this.sortOrder = document.getElementById("sortDirection")
    }

    setFilters() {
		const params = new URLSearchParams(window.location.search);

        this.tribunal.value = params.get("tribunal") ?? this.tribunal.value
        this.categorie.value = params.get("categorie") ?? this.categorie.value
        this.stadiuProcesual.value = params.get("stadiu_procesual") ?? this.stadiuProcesual.value

        this.sortBy.value = params.get("sort_by") ?? this.sortBy.value
        this.sortOrder.value = params.get("sort_order") ?? this.sortOrder.value
    }

    getFilters() {
		const params = new URLSearchParams();

        if (this.tribunal.value) {
            params.set("tribunal", this.tribunal.value)
        }

        if (this.categorie.value) {
            params.set("categorie", this.categorie.value)
        }

        if (this.stadiuProcesual.value) {
            params.set("stadiu_procesual", this.stadiuProcesual.value)
        }

		params.set('sort_by', this.sortBy.value)
        params.set('sort_order', this.sortOrder.value)

        return params
    }

    resetFilters() {
        window.location = `${window.location.pathname}?nume_firma=${this.numeFirma}`
    }

    getEmptyMessage() {
        return [
            "Nici un dosar găsit.",
            "Portaljust nu a returnat nici un dosar juridic."
        ]
    }

    getTimeoutMessage() {
        return [
            "Timeout căutare.",
            "Portaljust nu a răspuns în timp util."
        ]
    }

	getLinkForPage(pageNumber) {
		return `/dosareJuridice/${this.inregistrare}/${pageNumber}?nume_firma=${this.numeFirma}&${this.getFilters()}`
	}

	getLinkForDataRequest() {
		return `/dosareJuridiceFirma/${this.inregistrare}/${this.pageNumber}?${this.getFilters()}`
	}

	getSearchTitle() {
		return `Dosare juridice pentru: ${this.numeFirma}`
	}

    populatePrototype(data) {
        let numarDosar = data.Numar.replaceAll("/", "-")
        this.companyVeziProfil.href = `/dosarJuridic/${this.inregistrare}/${numarDosar}`

        this.companyName.textContent = data.Obiect
		this.companyJudet.textContent = data.CategorieCazNume
		this.companyTip.textContent = data.Numar
		this.companyCui.textContent = data.Institutie
		this.companyInregistrare.textContent = data.Departament
		this.companyDate.textContent = formatDateDosare(data.Data)
		this.companyStatus.textContent = data.StadiuProcesualNume
    }

	createResults(data) {
		for (let firma of data.Dosare) {
			this.populatePrototype(firma)

			let clone = this.prototypeCard.cloneNode(true)

			clone.style.display = "block"
			clone.removeAttribute("id")

			this.prototypeCard.before(clone)
		}
	}
}

function CreateComponent() {
    return new DosareSearchPage()
}
