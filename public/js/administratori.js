class AdministratoriSearchPage extends SearchPageBase {
	init() {
        super.init()

		const path = window.location.pathname.split("/")

        this.numeAdmin = decodeURIComponent(path.at(-2))
        this.inregistrare = path.at(-3)
	}

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
		return `/admins/${this.inregistrare}/${this.numeAdmin}/${pageNumber}?${this.getFilters()}`
	}

	getLinkForDataRequest() {
		return `/api/admins/${this.inregistrare}/${this.numeAdmin}/${this.pageNumber}?${this.getFilters()}`
	}

	getSearchTitle() {
		return `Companii asociate cu: ${this.numeAdmin}`
	}
}

function CreateComponent() {
    return new AdministratoriSearchPage()
}
