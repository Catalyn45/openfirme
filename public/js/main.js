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

const tribunal =
    document
        .getElementById("tribunalFilter")

const categorie =
    document
        .getElementById("categorieFilter")

const stadiuProcesual =
    document
        .getElementById("stadiuProcesualFilter")

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

    if (county?.value) {
        params.set('judet', county.value)
    }

    if (status?.value != null) {
        params.set('status', status.value)
    }

	if (formaJuridica?.value) {
        params.set('forma_juridica', formaJuridica.value)
    }

	if (tribunal?.value) {
		params.set("tribunal", tribunal.value)
	}

	if (categorie?.value) {
		params.set("categorie", categorie.value)
	}

	if (stadiuProcesual?.value) {
		params.set("stadiu_procesual", stadiuProcesual.value)
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

function searchCompanies(event, el) {
	let nume_partial = query.value.toLowerCase().trim()
	if (!nume_partial) {
        event.preventDefault()
		return
	}

    const params = getFilterParameters()
    el.href = `/search/1?${params}`
}

function changeFilter(el) {
    let endpoint = window.location.pathname.split("/")[1]
    if (!endpoint || endpoint === "profile") {
        return
    }

    searchWithFilters(endpoint)
}

/* ============================================================
   ENTER KEY
============================================================ */

function handleEnter(event) {
    if (event.key === "Enter") {
        let nume_partial = query.value.toLowerCase().trim()
        if (!nume_partial) {
            return
        }

        const params = getFilterParameters()
        window.location = `/search/1?${params}`
    }
}

function formatMoney(amount) {
	return new Intl.NumberFormat('de-DE', {
	  style: 'currency',
	  currency: 'RON',
	  maximumFractionDigits: 0
	}).format(amount);
}

function formatDate(d) {
	const [date, time] = d.split(" ");
	const [year, month, day] = date.split("-");

	let formatted = `${day}/${month}/${year}`
	if (time) {
		formatted = `${formatted} ${time}`
	}

	return formatted
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

function openProfile(button) {
	card = button.parentElement.parentElement;

	let inregistrare = card.getElementsByClassName("company-inregistrare")[0].textContent
	inregistrare = inregistrare.replaceAll("/", "-")

    button.href = `/profile/${inregistrare}?${getFilterParameters()}`
}

function updateFilters() {
    const params = new URLSearchParams(window.location.search);
	
	if (query) {
		query.value = params.get("nume_partial") ?? query.value
	}

	if (county) {
		county.value = params.get("judet") ?? county.value
	}

	if (status) {
		status.value = params.get("status") ?? status.value
	}

	if (formaJuridica) {
		formaJuridica.value = params.get("forma_juridica") ?? formaJuridica.value
	}

	if (tribunal) {
		tribunal.value = params.get("tribunal") ?? tribunal.value
	}

	if (categorie) {
		categorie.value = params.get("categorie") ?? categorie.value
	}

	if (stadiuProcesual) {
		stadiuProcesual.value = params.get("stadiu_procesual") ?? stadiuProcesual.value
	}

	if (sortBy) {
		sortBy.value = params.get("sort_by") ?? sortBy.value
		sortOrder.value = params.get("sort_order") ?? sortOrder.value
	}
}

updateFilters()
