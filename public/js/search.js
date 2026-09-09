const prototypeCard = document.getElementById("prototype")

const companyName = prototypeCard.getElementsByClassName("company-name")[0]
const companyStatus = prototypeCard.getElementsByClassName("company-status")[0]
const companyCui = prototypeCard.getElementsByClassName("company-cui")[0]
const companyInregistrare = prototypeCard.getElementsByClassName("company-inregistrare")[0]
const companyTip = prototypeCard.getElementsByClassName("company-tip")[0]
const companyJudet = prototypeCard.getElementsByClassName("company-judet")[0]
const companyDate = prototypeCard.getElementsByClassName("company-date")[0]
const companyVeziProfil = prototypeCard.getElementsByClassName("view-button")[0]

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

function openDosar(button) {
    let inregistrare = window.location.pathname.split('/').at(-2)

    let numar_dosar = button.parentNode.parentNode.getElementsByClassName("company-tip")[0].textContent
	numar_dosar = numar_dosar.replaceAll("/", "-")

    button.href = `/dosarJuridic/${inregistrare}/${numar_dosar}`
}

function changeDosarFilter(el) {
    let inregistrare = window.location.pathname.split("/").at(-2)

    const params = getFilterParameters()
	window.location = `/dosareJuridice/${inregistrare}/1?${params}`
}

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

		pagesContainer.style.display = "flex"

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

let pageEndpoint = window.location.pathname.split('/').slice(1, -1).join("/")

function goToPage(pageButton) {
	const pageNumber = parseInt(pageButton.innerText)
    pageButton.href = `/${pageEndpoint}/${pageNumber}?${getFilterParameters()}`
}

function goToNextPage(pageButton) {
	const pageNumber = parseInt(window.location.pathname.split('/').pop());
    pageButton.href = `/${pageEndpoint}/${pageNumber+1}?${getFilterParameters()}`
}

function goToPrevPage(pageButton) {
	const pageNumber = parseInt(window.location.pathname.split('/').pop());
    pageButton.href = `/${pageEndpoint}/${pageNumber-1}?${getFilterParameters()}`
}

function showEmpty(title, description) {
	if (title) {
		const emptyStateTitle = document.getElementById("emptyStateTitle")
		emptyStateTitle.textContent = title
	}

	if (description) {
		const emptyStateDescription = document.getElementById("emptyStateDescription")
		emptyStateDescription.textContent = description
	}
}

function populateSearch(data) {
    for (let firma of data.Data) {
        companyName.textContent = firma.Nume
		companyJudet.textContent = firma.Judet

		companyTip.textContent = firma.FormaJuridica
		companyCui.textContent = firma.Cui
		companyInregistrare.textContent = firma.CodInmatriculare
		companyDate.textContent = formatDate(firma.DataInregistrare)

		for (let child of companyStatus.parentNode.children) {
			if (child !== companyStatus) {
				child.remove()
			}
		}

		firma.Statusuri.sort((a, b) => {
			return (a === "funcțiune") - (b === "funcțiune")
		})

		for (let [index, statusFirma] of firma.Statusuri.entries()) {
			if (index > 0) {
				let statusClone = companyStatus.cloneNode(true)
				companyStatus.before(statusClone)
			}

			companyStatus.textContent = statusFirma
			companyStatus.title = statusFirma
			companyStatus.classList.remove("active")
			companyStatus.classList.remove("inactive")
			if (statusFirma === "funcțiune") {
				companyStatus.classList.add("active")
			} else {
				companyStatus.classList.add("inactive")
			}
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
}

function populateDosare(data) {
    for (let dosar of data.Dosare) {
        companyName.textContent = dosar.Obiect
		companyJudet.textContent = dosar.CategorieCazNume
		companyTip.textContent = dosar.Numar
		companyCui.textContent = dosar.Institutie
		companyInregistrare.textContent = dosar.Departament
		companyDate.textContent = dosar.Data
		companyStatus.textContent = dosar.StadiuProcesualNume

        clone = prototypeCard.cloneNode(true)

        clone.style.display = "block"
        clone.removeAttribute("id")

        prototypeCard.before(clone)
    }
}

async function main() {
	let path = window.location.pathname.split('/')

	const pageNumber = path.pop();

    let endpoint = ''

    if (pageEndpoint.includes("top")) {
        endpoint = `/topFirme`
    } else if (pageEndpoint.includes("search")) {
        endpoint = `/firme`
    } else if (pageEndpoint.includes("admin")){
		endpoint = `/adminsFirme/${path[2]}/${path[3]}`
		document.getElementById("company-administrator").textContent = `Companii admnistrate de: ${decodeURIComponent(path[3])}`
	} else if (pageEndpoint.includes("dosareJuridice")) {
		endpoint = `/dosareJuridiceFirma/${path[2]}`
		document.getElementById("company-administrator").textContent = `Dosare juridice pentru: ${decodeURIComponent(path[2])}`
	}

    endpoint = `${endpoint}/${pageNumber}?${getFilterParameters()}`

	console.log("searching")
    const response = await fetch(endpoint)
	if (response.status === 422) {
		showEmpty("Căutarea este prea generică", "Încearcă să folosești un nume mai specific sau modifică filtrele de căutare.")
		return
	}

    const data = await response.json()
    console.log(data)

    let resultCount = data.Count
	console.log(resultCount)

	if (pageEndpoint.includes("dosareJuridice")) {
		populateDosare(data)
	} else {
		populateSearch(data)
	}

    document
        .getElementById("resultCount")
        .textContent = `${resultCount} rezultate`

	if (resultCount === 0) {
		showEmpty("Nu am găsit nicio companie", "Încearcă o altă denumire sau modifică filtrele de căutare.")
	} else {
		document
			.getElementById("emptyState")
			.style.display = "none"
    }

	updatePageNumbers(pageNumber, resultCount)
}

main().catch(console.error);
