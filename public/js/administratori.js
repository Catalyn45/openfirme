class AdministratoriSearchPage extends SearchPage {
    constructor() {
        super()

		const path = window.location.pathname.split("/")

        this.numeAdmin = decodeURIComponent(path.at(-2))
        this.inregistrare = path.at(-3)
    }

    initFilters() {
    }

    setFilters() {
        return ""
    }

    getFilters() {
        return ""
    }

	getLinkForPage(pageNumber) {
		return `/admins/${this.inregistrare}/${this.numeAdmin}/${pageNumber}`
	}

	getLinkForDataRequest() {
		return `/adminsFirme/${this.inregistrare}/${this.numeAdmin}`
	}

	getSearchTitle() {
		return `Companii asociate cu: ${this.numeAdmin}`
	}
}

function CreateComponent() {
    return new AdministratoriSearchPage()
}
