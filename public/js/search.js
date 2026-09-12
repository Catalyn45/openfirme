class SearchPage extends BaseComponent {
	constructor() {
        super()

		this.initPrototype()
		this.initPages()
	}

	initPrototype() {
		this.prototypeCard = document.getElementById("prototype")

		this.companyName = this.prototypeCard.getElementsByClassName("company-name")[0]
		this.companyStatus = this.prototypeCard.getElementsByClassName("company-status")[0]
		this.companyCui = this.prototypeCard.getElementsByClassName("company-cui")[0]
		this.companyInregistrare = this.prototypeCard.getElementsByClassName("company-inregistrare")[0]
		this.companyTip = this.prototypeCard.getElementsByClassName("company-tip")[0]
		this.companyJudet = this.prototypeCard.getElementsByClassName("company-judet")[0]
		this.companyDate = this.prototypeCard.getElementsByClassName("company-date")[0]

		this.companyVeziProfil = this.prototypeCard.getElementsByClassName("view-button")[0]

		this.resultCount = document.getElementById("resultCount")
		this.emptyState = document.getElementById("emptyState")
	}

	initPages() {
		this.firstDots = document.getElementById("pageFirstDots")
		this.secondDots = document.getElementById("pageSecondDots")

		this.pagesContainer = document.getElementById("paginationContainer")

		this.leftPage = document.getElementsByClassName("left-page-number")[0]
		this.middlePages = document.getElementsByClassName("middle-page-number")
		this.rightPage = document.getElementsByClassName("right-page-number")[0]

		this.backButton = document.getElementById("pageBack")
		this.nextButton = document.getElementById("pageNext")

		this.pageNumber = parseInt(window.location.pathname.split("/").pop())
	}

	onChangeFilter() {
		const params = this.getFilters()
		const endpoint = this.getLinkForPage(1)

		window.location = `${endpoint}?${params}`
	}

	getLinkForPage(pageNumber) {
		return `/search/${pageNumber}`
	}

	getLinkForDataRequest() {
		return "/firme"
	}

	populatePrototype(data) {
		this.companyName.textContent = data.Nume
		this.companyJudet.textContent = data.Judet

		this.companyTip.textContent = data.FormaJuridica
		this.companyCui.textContent = data.Cui
		this.companyInregistrare.textContent = data.CodInmatriculare
		this.companyDate.textContent = formatDate(data.DataInregistrare)

		for (let child of this.companyStatus.parentNode.children) {
			if (child !== this.companyStatus) {
				child.remove()
			}
		}

		data.Statusuri.sort((a, b) => {
			return (a === "funcțiune") - (b === "funcțiune")
		})

		for (let [index, statusFirma] of data.Statusuri.entries()) {
			if (index > 0) {
				let statusClone = this.companyStatus.cloneNode(true)
				this.companyStatus.before(statusClone)
			}

			this.companyStatus.textContent = statusFirma
			this.companyStatus.title = statusFirma

			this.companyStatus.classList.remove("active")
			this.companyStatus.classList.remove("inactive")
			if (statusFirma === "funcțiune") {
				this.companyStatus.classList.add("active")
			} else {
				this.companyStatus.classList.add("inactive")
			}
		}

		let inregistrare = data.CodInmatriculare.replaceAll("/", "-")
		this.companyVeziProfil.href = `/profile/${inregistrare}?${this.getFilters()}`
	}

	createResults(data) {
		for (let firma of data.Data) {
			this.populatePrototype(firma)

			let clone = this.prototypeCard.cloneNode(true)

			clone.style.display = "block"
			clone.removeAttribute("id")

			this.prototypeCard.before(clone)
		}
	}

	getSearchTitle() {
		return null
	}

	setPageNumber(pageElement, pageNumber) {
		pageElement.textContent = pageNumber
		pageElement.href = `${this.getLinkForPage(pageNumber)}?${this.getFilters()}`
	}

	updatePageNumbers(totalResults) {
		if (totalResults <= 20) {
			return
		}

		const totalPages = Math.ceil(totalResults / 20);

		if (this.pageNumber == 1) {
			this.backButton.style.display = "none"
		} else if (this.pageNumber == totalPages) {
			this.nextButton.style.display = "none"
		}

		this.setPageNumber(this.rightPage, totalPages)
        this.setPageNumber(this.leftPage, 1)

		if (totalPages <= 5) {
			this.firstDots.style.display = "none"
			this.secondDots.style.display = "none"

			for (const [index, item] of [this.leftPage, ...this.middlePages, this.rightPage].entries()) {
				if (index + 1 > totalPages) {
					item.style.display = "none"
					continue
				}

				this.setPageNumber(item, index+1)

				if (index + 1 == this.pageNumber) {
					item.classList.add("active")
				}
			}

			this.pagesContainer.style.display = "flex"

			return
		}

		if (this.pageNumber <= 3) {
			this.firstDots.style.display = "none"
		}

		if (this.pageNumber > totalPages - 3) {
			this.secondDots.style.display = "none"
		}

		let startingPage = this.pageNumber - 1
		if (this.pageNumber <= 3) {
			startingPage = 2
		} else if (this.pageNumber > totalPages - 3) {
			startingPage = totalPages - 3
		}

		for (let el of this.middlePages) {
			this.setPageNumber(el, startingPage)

			if (startingPage == this.pageNumber) {
				el.classList.add("active")
			}

			startingPage = startingPage + 1
		}

		if (this.pageNumber <= 1) {
			this.leftPage.classList.add("active")
		} else if (this.pageNumber > totalPages - 1) {
			this.rightPage.classList.add("active")
		}

		this.pagesContainer.style.display = "flex"
	}

	showEmpty(title, description) {
		if (title) {
			const emptyStateTitle = document.getElementById("emptyStateTitle")
			emptyStateTitle.textContent = title
		}

		if (description) {
			const emptyStateDescription = document.getElementById("emptyStateDescription")
			emptyStateDescription.textContent = description
		}
	}

	async Start() {
        super.Start()

		let dataRequestLink = `${this.getLinkForDataRequest()}/${this.pageNumber}?${this.getFilters()}`

		console.log("searching")
		const response = await fetch(dataRequestLink)
		if (response.status === 422) {
			this.showEmpty("Căutarea este prea generică", "Încearcă să folosești un nume mai specific sau modifică filtrele de căutare.")
			return
		}

		const searchTitle = this.getSearchTitle()
		if (searchTitle) {
			document.getElementById("company-administrator").textContent = searchTitle
		}

		const data = await response.json()
		console.log(data)

		let resultCount = data.Count
		console.log(resultCount)

		this.createResults(data)

		this.resultCount.textContent = `${resultCount} rezultate`

		if (resultCount === 0) {
			this.showEmpty("Nu s-a găsit nici un rezultat", "Nici un rezultat găsit, incearcă să schimbi filtrele de căutare.")
		} else {
			this.emptyState.style.display = "none"
		}

		this.updatePageNumbers(resultCount)
	}
}

function CreateComponent() {
    return new SearchPage()
}
