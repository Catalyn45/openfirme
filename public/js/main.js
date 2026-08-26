/* ============================================================
   SEARCH
============================================================ */
const query =
    document
        .getElementById("searchInput")

const county =
    document
        .getElementById("countyFilter")

const status =
    document
        .getElementById("statusFilter")

const formaJuridica =
    document
        .getElementById("formaFilter")

const sortBy =
    document
        .getElementById("sortField")

const sortOrder =
    document
        .getElementById("sortDirection")

async function searchCompanies() {
	let nume_partial = query.value.toLowerCase().trim()
	if (!nume_partial) {
		return
	}

    let endpoint = `/search/1?nume_partial=${nume_partial}`

    if (county.value) {
        endpoint = `${endpoint}&judet=${county.value}`
    }

    if (status.value) {
        endpoint = `${endpoint}&status=${status.value}`
    }

	if (formaJuridica.value) {
        endpoint = `${endpoint}&forma_juridica=${formaJuridica.value}`
	}

	window.location = endpoint
}

/* ============================================================
   ENTER KEY
============================================================ */

function handleEnter(event) {

    if (event.key === "Enter") {
        searchCompanies();
    }

}


/* ============================================================
   RESET FILTERS
============================================================ */

function resetFilters() {
    query.value = "";
    county.value = "";
    status.value = ""
    an.value = ""
    formaJuridica.value = ""

    window.location = window.location.pathname
}


/* ============================================================
   COMPANY PROFILE
============================================================ */

async function openProfile(button) {
	card = button.parentElement.parentElement;

	let inregistrare = card.getElementsByClassName("company-inregistrare")[0].textContent
	inregistrare = inregistrare.replaceAll("/", "-")

    window.location = `/profile/${inregistrare}${window.location.search}`
}

function updateFilters() {
    const params = new URLSearchParams(window.location.search);
	
	if (query) {
		query.value = params.get("nume_partial") ?? ""
	}

    county.value = params.get("judet") ?? ""
    status.value = params.get("status") ?? ""
    formaJuridica.value = params.get("forma_juridica") ?? ""

	if (sortBy) {
		sortBy.value = params.get("sort_by") ?? "infiintare"
		sortOrder.value = params.get("sort_order") ?? "desc"
	}
}

updateFilters()
