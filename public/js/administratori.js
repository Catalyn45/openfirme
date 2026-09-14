class AdministratoriSearchPage extends SearchPageBase {
	init() {
		const path = window.location.pathname.split("/")

        this.numeAdmin = decodeURIComponent(path.at(-2))
        this.inregistrare = path.at(-3)

		this.initFilters()
		this.initPrototype()
		this.initPages()
	}

    setFilters() { }

    getEmptyMessage() {
        return [
            "Nici un administrator găsit.",
            "Administratorul nu este asociat cu nici-o firmă."
        ]
    }

    getTimeoutMessage() {
        return [
            "Nici un administrator găsit.",
            "Administratorul nu este asociat cu nici-o firmă."
        ]
    }

	getLinkForPage(pageNumber) {
		return `/admins/${this.inregistrare}/${this.numeAdmin}/${pageNumber}`
	}

	getLinkForDataRequest() {
		return `/api/admins/${this.inregistrare}/${this.numeAdmin}/${this.pageNumber}`
	}

	getSearchTitle() {
		return `Companii asociate cu: ${this.numeAdmin}`
	}
}

function CreateComponent() {
    return new AdministratoriSearchPage()
}
