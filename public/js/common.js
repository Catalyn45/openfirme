function resetFilters() {
    window.location = window.location.pathname
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

function setValueIfExist(element, value) {
	if (value) {
		element.textContent = value
	}
}

async function setErrorPage(statusCode) {
    let errorPage = await fetch("/error")
    let errorPageContent = await errorPage.text()

    document.open();
    document.write(errorPageContent);
    document.close();

    let errorText = document.getElementById("profileCompanyName")
    errorText.textContent += statusCode
}
