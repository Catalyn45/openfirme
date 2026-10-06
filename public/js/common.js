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
    if (!d) {
        return d
    }

    const date = new Date(d);

    const formatted = date.toLocaleDateString("en-GB") + " " + date.toLocaleTimeString("en-GB", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false
    })

    return formatted
}

function formatDateSedinta(d) {
    if (!d) {
        return d
    }

    const date = new Date(d);
    const formatted = date.toLocaleDateString("en-GB")

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

async function copyContentToClipboard(button, content) {
	let success = await setClipboard(content)

    if (!success)
        return

    let originalText = button.textContent
    button.textContent = "Copiat"
    button.classList.add("content-copy-button-success")

    setTimeout(() => {
        button.textContent = originalText
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

const LABEL_TOGGLE_COLLAPSIBLE_CONTENT_BUTTON_SHOW = "Arată mai mult"
const LABEL_TOGGLE_COLLAPSIBLE_CONTENT_BUTTON_HIDE = "Arată mai puțin"

async function toggleCollapsibleContent(button) {
	let targetId = button.dataset.target

	if (targetId === undefined)
		return

	let target = document.getElementById(targetId)

	if (target === undefined)
		return

	if (target.classList.contains("collapsible-content-active"))
        button.textContent = LABEL_TOGGLE_COLLAPSIBLE_CONTENT_BUTTON_SHOW
    else
        button.textContent = LABEL_TOGGLE_COLLAPSIBLE_CONTENT_BUTTON_HIDE
    

    target.classList.toggle("collapsible-content-active")
}

const domeniiMap = {
    "A":"Agricultură, silvicultură și pescuit",
    "B":"Industria extractivă",
    "C":"Industria prelucrătoare",
    "D":"Producția și furnizarea de energie electrică",
    "E":"Distribuția apei, salubritate si gestionarea deșeurilor",
    "F":"Construcții",
    "G":"Comerț cu ridicata și cu amănuntul",
    "H":"Transport și depozitare",
    "I":"Hoteluri și restaurante",
    "J":"Activități de editare",
    "K":"Telecomunicații",
    "L":"Intermedieri financiare și asigurări",
    "M":"Tranzacții imobiliare",
    "N":"Activități profesionale, științifice și tehnice",
    "O":"Activități de servicii administrative",
    "P":"Administrație publică și apărare",
    "Q":"Învățământ",
    "R":"Sănătate și asistență socială",
    "S":"Activități de spectacole, sportive și recreative",
    "T":"Alte activități de servicii",
    "U":"Activități ale gospodăriilor",
    "V":"Activități ale organizațiilor extrateritoriale"
}
