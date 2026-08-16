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

    const response = await fetch(`/firme?nume_partial=${query}`);


    const cards = document.getElementsByClassName("company-card")

    let prototypeCard = null
    for (let el of Array.from(cards)) {
        if (el.id === "prototype") {
            prototypeCard = el
        } else {
            el.remove()
        }
    }

    const data = await response.json()

    console.log(data)

    let visible = data.length;

    for (let firma of data) {
        clone = prototypeCard.cloneNode(true)

        clone.style.display = "block"
        clone.removeAttribute("id")

        let companyName = clone.getElementsByClassName("company-name")[0]

        companyName.textContent = firma.Nume

        prototypeCard.before(clone)
    }

    /*


    cards.forEach(card => {
        const name =
            card.dataset.name.toLowerCase();

        const cui =
            card.dataset.cui.toLowerCase();

        const cardCounty =
            card.dataset.county;

        const cardStatus =
            card.dataset.status;

        const cardIndustry =
            card.dataset.industry;


        const matchesSearch =
            query === "" ||
            name.includes(query) ||
            cui.includes(query);


        const matchesCounty =
            county === "" ||
            cardCounty === county;


        const matchesStatus =
            status === "" ||
            cardStatus === status;


        const matchesIndustry =
            industry === "" ||
            cardIndustry === industry;


        const visibleCard =
            matchesSearch &&
            matchesCounty &&
            matchesStatus &&
            matchesIndustry;


        if (visibleCard) {
            card.style.display = "block";
            visible++;
        } else {
            card.style.display = "none";
        }

    });

    */

    document
        .getElementById("resultCount")
        .textContent = `${visible} rezultat${visible === 1 ? "" : "e"}`;


    document
        .getElementById("emptyState")
        .style.display = visible === 0 ? "block" : "none";
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

    searchCompanies();

}


/* ============================================================
   COMPANY PROFILE
============================================================ */

function openProfile(companyName) {

    document
        .getElementById("profileCompanyName")
        .textContent = companyName;

    document
        .getElementById("profileName")
        .textContent = companyName;

    document
        .getElementById("profileModal")
        .style.display = "block";

    document.body.style.overflow = "hidden";

}


function closeProfile() {

    document
        .getElementById("profileModal")
        .style.display = "none";

    document.body.style.overflow = "auto";

}


function closeModalOutside(event) {
    if (
        event.target ===
        document.getElementById("profileModal")
    ) {
        closeProfile();
    }
}


/* ============================================================
   ESC KEY
============================================================ */

document.addEventListener(
    "keydown",
    function(event) {

        if (event.key === "Escape") {
            closeProfile();
        }

    }
);

