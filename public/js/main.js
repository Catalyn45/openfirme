/* ============================================================
   SEARCH
============================================================ */


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

    searchCompanies();

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
