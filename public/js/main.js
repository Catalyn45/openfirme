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

function getFilterParameters() {
    const params = new URLSearchParams();

	let nume_partial = query?.value?.toLowerCase().trim()
	if (nume_partial) {
        params.set('nume_partial', nume_partial)
	}

    if (county.value) {
        params.set('judet', county.value)
    }

    if (status.value) {
        params.set('status', status.value)
    }

	if (formaJuridica.value) {
        params.set('forma_juridica', formaJuridica.value)
    }

	if (sortBy?.value) {
		params.set('sort_by', sortBy.value)

		if (sortOrder) {
            params.set('sort_order', sortOrder.value)
		}
	}

    return params
}

function searchWithFilters(endpoint) {
    const params = getFilterParameters()

	window.location = `/${endpoint}/1?${params}`
}

function searchCompanies() {
	let nume_partial = query.value.toLowerCase().trim()
	if (!nume_partial) {
		return
	}

    searchWithFilters('search')
}

function changeFilter(el) {
    let endpoint = window.location.pathname.split("/")[1]

    searchWithFilters(endpoint)
}

/* ============================================================
   ENTER KEY
============================================================ */

function handleEnter(event) {

    if (event.key === "Enter") {
        searchCompanies();
    }

}

function formatMoney(amount) {
	return new Intl.NumberFormat('de-DE', {
	  style: 'currency',
	  currency: 'RON',
	  maximumFractionDigits: 0
	}).format(amount);
}

/* ============================================================
   RESET FILTERS
============================================================ */

function resetFilters() {
    window.location = window.location.pathname
}


/* ============================================================
   COMPANY PROFILE
============================================================ */

async function openProfile(button) {
	card = button.parentElement.parentElement;

	let inregistrare = card.getElementsByClassName("company-inregistrare")[0].textContent
	inregistrare = inregistrare.replaceAll("/", "-")

    window.location = `/profile/${inregistrare}?${getFilterParameters()}`
}

function updateFilters() {
    const params = new URLSearchParams(window.location.search);
	
	if (query) {
		query.value = params.get("nume_partial") ?? query.value
	}

    county.value = params.get("judet") ?? county.value
    status.value = params.get("status") ?? status.value
    formaJuridica.value = params.get("forma_juridica") ?? formaJuridica.value

	if (sortBy) {
		sortBy.value = params.get("sort_by") ?? sortBy.value
		sortOrder.value = params.get("sort_order") ?? sortOrder.value
	}
}

updateFilters()
