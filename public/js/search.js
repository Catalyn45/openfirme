class SearchPage extends SearchPageBase {
	init() {
        super.init()

		this.numePartial = window.location.pathname.split("/").at(-2)
		this.query.value = decodeURIComponent(this.numePartial)
	}

    searchFirma() {
        let numePartial = this.query.value?.toLowerCase().trim()
        if (numePartial.length < 3) {
            return
        }

        let encoded = encodeURIComponent(numePartial)
        window.location = `/search/${encoded}/1?${this.getFilters()}`
    }

	getLinkForPage(pageNumber) {
		return `/search/${this.numePartial}/${pageNumber}?${this.getFilters()}`
	}

	getLinkForDataRequest() {
		return `/api/search/${this.numePartial}/${this.pageNumber}?${this.getFilters()}`
	}

	getSearchTitle() {
		return `Rezultate căutare: ${decodeURIComponent(this.numePartial)}`
	}
}

function CreateComponent() {
    return new SearchPage()
}
