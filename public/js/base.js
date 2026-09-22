class Base {
	constructor() {
		this.init()
	}

    init() {
        this.initSearchBar()
    }

    initSearchBar() {
        this.query = document.getElementById("searchInput")
    }

    searchFirma() {
        let numePartial = this.query.value?.toLowerCase().trim()
        if (numePartial.length < 3) {
            return
        }

        let encoded = encodeURIComponent(numePartial)
        window.location = `/search/${encoded}/1`
    }

    copyToClipboard(el) {
        copyContentToClipboard(
            el,
            el.parentNode.getElementsByClassName("clipboardContent")[0].textContent
        )
    }

    Start() { }
}
