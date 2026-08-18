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

    const industry =
        document
            .getElementById("industryFilter")
            .value;

    // let endpoint = `/firme?nume_partial=${query}`
    let endpoint = `/search?nume_partial=${query}`

    if (county) {
        endpoint = `${endpoint}&judet=${county}`
    }

    if (status) {
        endpoint = `${endpoint}&status=${status}`
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
        .getElementById("industryFilter")
        .value = "";

    document
        .getElementById("yearFilter")
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
