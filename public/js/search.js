const prototypeCard = document.getElementById("prototype")

const companyName = prototypeCard.getElementsByClassName("company-name")[0]
const companyStatus = prototypeCard.getElementsByClassName("company-status")[0]
const companyCui = prototypeCard.getElementsByClassName("company-cui")[0]
const companyInregistrare = prototypeCard.getElementsByClassName("company-inregistrare")[0]
const companyTip = prototypeCard.getElementsByClassName("company-tip")[0]
const companyJudet = prototypeCard.getElementsByClassName("company-judet")[0]
const companyDate = prototypeCard.getElementsByClassName("company-date")[0]

const companyAngajati = prototypeCard.getElementsByClassName("company-angajati")?.[0]
const companyProfit = prototypeCard.getElementsByClassName("company-profit")?.[0]
const companyCifraAfaceri = prototypeCard.getElementsByClassName("company-cifra-afaceri")?.[0]

const firstDots = document.getElementById("pageFirstDots")
const secondDots = document.getElementById("pageSecondDots")

const pagesContainer = document.getElementById("paginationContainer")

const leftPages = document.getElementsByClassName("left-page-number")
const middlePages = document.getElementsByClassName("middle-page-number")
const rightPages = document.getElementsByClassName("right-page-number")

const backButton = document.getElementById("pageBack")
const nextButton = document.getElementById("pageNext")

function updatePageNumbers(currentPage, totalResults) {
	if (totalResults <= 20) {
		return
	}

	const totalPages = Math.ceil(totalResults / 20);

	if (currentPage == 1) {
		backButton.style.display = "none"
	} else if (currentPage == totalPages) {
		nextButton.style.display = "none"
	}

	rightPages[rightPages.length - 1].textContent = totalPages
	rightPages[rightPages.length - 2].textContent = totalPages - 1

	if (totalPages <= 9) {
		firstDots.style.display = "none"
		secondDots.style.display = "none"

		for (const [index, item] of [...leftPages, ...middlePages, ...rightPages].entries()) {
			if (index + 1 > totalPages) {
				item.style.display = "none"
				continue
			}

			item.innerText = index + 1

			if (index + 1 == currentPage) {
				item.classList.add("active")
			}
		}

		return
	}

	if (currentPage <= 5) {
		firstDots.style.display = "none"
	}

	if (currentPage > totalPages - 5) {
		secondDots.style.display = "none"
	}

	let startingPage = currentPage - 2
	if (currentPage <= 5) {
		startingPage = 3
	} else if (currentPage > totalPages - 5) {
		startingPage = totalPages - 7
	}

	for (let el of middlePages) {
		el.innerText = startingPage

		if (startingPage == currentPage) {
			el.classList.add("active")
		}

		startingPage = startingPage + 1
	}

	if (currentPage <= 2) {
		leftPages[currentPage-1].classList.add("active")
	} else if (currentPage > totalPages - 2) {
		rightPages[currentPage - totalPages + 1].classList.add("active")
	}

	pagesContainer.style.display = "flex"
}

let pageEndpoint = ''
if (window.location.pathname.includes("top")) {
	pageEndpoint = `top`
} else {
	pageEndpoint = `search`
}

function goToPage(pageButton) {
	const pageNumber = parseInt(pageButton.innerText)
    window.location = `/${pageEndpoint}/${pageNumber}?${getFilterParameters()}`
}

function goToNextPage() {
	const pageNumber = parseInt(window.location.pathname.split('/').pop());
    window.location = `/${pageEndpoint}/${pageNumber+1}?${getFilterParameters()}`
}

function goToPrevPage() {
	const pageNumber = parseInt(window.location.pathname.split('/').pop());
    window.location = `/${pageEndpoint}/${pageNumber-1}?${getFilterParameters()}`
}

async function main() {
	const pageNumber = window.location.pathname.split('/').pop();

    let endpoint = ''

    if (window.location.pathname.includes("top")) {
        endpoint = `/topFirme`
    } else {
        endpoint = `/firme`
    }

    endpoint = `${endpoint}/${pageNumber}?${getFilterParameters()}`

    const response = await fetch(endpoint)
    const data = await response.json()

    console.log(data)

    let resultCount = data.Count

    for (let firma of data.Data) {
        companyName.textContent = firma.Nume
		companyJudet.textContent = firma.Judet

		companyTip.textContent = firma.FormaJuridica
		companyCui.textContent = firma.Cui
		companyInregistrare.textContent = firma.CodInmatriculare
		companyDate.textContent = firma.DataInregistrare

		companyStatus.textContent = firma.Status

		companyStatus.classList.remove("active")
		companyStatus.classList.remove("inactive")
		if (firma.Status === "funcțiune") {
			companyStatus.classList.add("active")
		} else {
			companyStatus.classList.add("inactive")
		}

		if (companyProfit) {
			companyProfit.textContent = formatMoney(firma.ProfitNet)
		}

		if (companyCifraAfaceri) {
			companyCifraAfaceri.textContent = formatMoney(firma.CifraAfaceri)
		}

		if (companyAngajati) {
			companyAngajati.textContent = firma.Angajati ?? 0
		}

        clone = prototypeCard.cloneNode(true)

        clone.style.display = "block"
        clone.removeAttribute("id")

        prototypeCard.before(clone)
    }

    document
        .getElementById("resultCount")
        .textContent = `${resultCount} rezultate`

	if (resultCount === 0) {
		document
			.getElementById("emptyState")
			.style.display = "block"
	}

	updatePageNumbers(pageNumber, resultCount)
}

main()
