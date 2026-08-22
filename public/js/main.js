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

const an =
    document
        .getElementById("yearFilter")

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

    let endpoint = `/search/1?nume_partial=${}`


    if (county.value) {
        endpoint = `${endpoint}&judet=${county.value}`
    }

    if (status.value) {
        endpoint = `${endpoint}&status=${status.value}`
    }

	if (formaJuridica.value) {
        endpoint = `${endpoint}&forma_juridica=${formaJuridica.value}`
	}

    if (an.value) {
        [after, before] = an.value.split("/")

        if (after) {
            endpoint = `${endpoint}&data_after=${after}`
        }

        if (before) {
            endpoint = `${endpoint}&data_before=${before}`
        }
    }

	if (sortBy.value) {
		endpoint = `${endpoint}&sort_by=${sortBy.value}`

		if (sortOrder) {
			endpoint = `${endpoint}&sort_order=${sortOrder.value}`
		}
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

    sortBy.value = ""
    sortOrder.value = "asc"

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

    query.value = params.get("nume_partial") ?? ""

    county.value = params.get("judet") ?? ""
    status.value = params.get("status") ?? ""
    formaJuridica.value = params.get("forma_juridica") ?? ""

    let anFormatted = `${params.get('data_after') ?? ''}/${params.get('data_before') ?? ''}`
    if (anFormatted == "/") {
        anFormatted = ""
    }

    an.value = anFormatted

    sortBy.value = params.get("sort_by") ?? ""
    sortOrder.value = params.get("sort_order") ?? "asc"
}

updateFilters()
