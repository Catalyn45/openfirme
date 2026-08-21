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
    let endpoint = `/search/1?nume_partial=${query.value.toLowerCase().trim()}`


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
    document
        .getElementById("searchInput")
        .value = "";

    document
        .getElementById("countyFilter")
        .value = "";

    document
        .getElementById("statusFilter")
        .value = "";

    document
        .getElementById("yearFilter")
        .value = "";

	document
		.getElementById("formaFilter")
		.value = "";
}


/* ============================================================
   COMPANY PROFILE
============================================================ */

async function openProfile(button) {
	card = button.parentElement.parentElement;

	let inregistrare = card.getElementsByClassName("company-inregistrare")[0].textContent
	inregistrare = inregistrare.replaceAll("/", "-")

    window.location = `/profile/${inregistrare}`
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
    sortOrder.value = params.get("sort_order") ?? ""
}

updateFilters()
