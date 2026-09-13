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

function formatDateDosare(d) {
    const date = new Date(d);

    const formatted = date.toLocaleDateString("en-GB") + " " + date.toLocaleTimeString("en-GB", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false
    })

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

async function copyContentToClipboard(button) {
	let targetId = button.dataset.target

	if (targetId === undefined)
		return

	let target = document.getElementById(targetId)

	if (target === undefined)
		return

	let success = await setClipboard(target.innerText)

    if (!success)
        return

    let originalText = button.innerText
    button.innerText = "Copiat"
    button.classList.add("content-copy-button-success")

    setTimeout(() => {
        button.innerText = originalText
        button.classList.remove("content-copy-button-success")
    }, 500)

}

async function setClipboard(text) {
	const type = "text/plain";
	const clipboardItemData = {
		[type]: text,
	};

    const clipboardItem = new ClipboardItem(clipboardItemData);

	try {
        await navigator.clipboard.write([clipboardItem]);
        return true;
    } catch (error) {
        console.error("Clipboard write failed:", error);
        return false;
    }
}
