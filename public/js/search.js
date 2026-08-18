const prototypeCard = document.getElementById("prototype")

const companyName = prototypeCard.getElementsByClassName("company-name")[0]
const companyStatus = prototypeCard.getElementsByClassName("company-status")[0]
const companyCui = prototypeCard.getElementsByClassName("company-cui")[0]
const companyInregistrare = prototypeCard.getElementsByClassName("company-inregistrare")[0]
const companyTip = prototypeCard.getElementsByClassName("company-tip")[0]
const companyJudet = prototypeCard.getElementsByClassName("company-judet")[0]
const companyDate = prototypeCard.getElementsByClassName("company-date")[0]
const companyCaen = prototypeCard.getElementsByClassName("company-caen")[0]

async function searchCompanies() {
    const query =
        document
            .getElementById("searchInput")
            .value
            .toLowerCase()
            .trim();

    const county =
        document
            .getElementById("countyFilter")
            .value;

    const status =
        document
            .getElementById("statusFilter")
            .value;

    const industry =
        document
            .getElementById("industryFilter")
            .value;

    let endpoint = `/firme?nume_partial=${query}`

    if (county) {
        endpoint = `${endpoint}&judet=${county}`
    }

    if (status) {
        endpoint = `${endpoint}&status=${status}`
    }

    const response = await fetch(endpoint)

    const data = await response.json()

    const cards = document.getElementsByClassName("company-card")
    for (let el of Array.from(cards)) {
        if (el.id !== "prototype") {
            el.remove()
        }
	}

    console.log(data)

    let resultCount = data.length;

    for (let firma of data) {
        companyName.textContent = firma.Nume
		companyJudet.textContent = firma.Judet

		if (firma.CoduriCaen) {
			companyCaen.textContent = firma.CoduriCaen[0]
		} else {
			companyCaen.textContent = ""
		}

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

        clone = prototypeCard.cloneNode(true)

        clone.style.display = "block"
        clone.removeAttribute("id")

        prototypeCard.before(clone)
    }

    document
        .getElementById("resultCount")
        .textContent = `${resultCount} rezultat${resultCount === 1 ? "" : "e"}`;

    document
        .getElementById("emptyState")
        .style.display = resultCount === 0 ? "block" : "none";
}

