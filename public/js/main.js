/* ============================================================
   SEARCH
============================================================ */

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

    const formaJuridica =
        document
            .getElementById("formaFilter")
            .value;

    const an =
        document
            .getElementById("yearFilter")
            .value;

    const sortBy =
        document
            .getElementById("sortField")
            .value;

    const sortOrder =
        document
            .getElementById("sortDirection")
            .value;

    let endpoint = `/search/1?nume_partial=${query}`

    if (county) {
        endpoint = `${endpoint}&judet=${county}`
    }

    if (status) {
        endpoint = `${endpoint}&status=${status}`
    }

	if (formaJuridica) {
        endpoint = `${endpoint}&forma_juridica=${formaJuridica}`
	}

    if (an) {
        [after, before] = an.split("/")

        if (after) {
            endpoint = `${endpoint}&data_after=${after}`
        }

        if (before) {
            endpoint = `${endpoint}&data_before=${before}`
        }
    }

	if (sortBy) {
		endpoint = `${endpoint}&sort_by=${sortBy}`

		if (sortOrder) {
			endpoint = `${endpoint}&sort_order=${sortOrder}`
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
