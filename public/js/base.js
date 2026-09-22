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

    getSearchIfValid() {
        let numePartial = this.query.value?.toLowerCase().trim()

        let words = numePartial.split(" ")
        for (let word of words) {
            if (word.length >= 3) {
                return encodeURIComponent(numePartial)
            }
        }

        this.query.setCustomValidity("Trebuie să aveți măcar un cuvânt de cel puțin 3 caractere");
        this.query.reportValidity();

        setTimeout(() => {
            this.query.setCustomValidity("");
        }, 5000);

        return null
    }

    searchFirma() {
        let numePartial = this.getSearchIfValid()
        if (!numePartial) {
            return
        }

        window.location = `/search/${numePartial}/1`
    }

    copyToClipboard(el) {
        copyContentToClipboard(
            el,
            el.parentNode.getElementsByClassName("clipboardContent")[0].textContent
        )
    }

    Start() { }
}
