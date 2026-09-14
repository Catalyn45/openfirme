function createAddress(json) {
	let adresa = ""

	if (json.Sector) {
		adresa += `Sector ${json.Sector}, `
	}

	if (json.Strada) {
		adresa += `Str. ${json.Strada}`
	}

	if (json.NrStrada) {
		adresa += `, nr. ${json.NrStrada}`
	}

	if (json.Bloc) {
		adresa += `, bl. ${json.Bloc}`
	}

	if (json.Scara) {
		adresa += `, sc. ${json.Scara}`
	}

	if (json.Etaj) {
		adresa += `, et. ${json.Etaj}`
	}

	if (json.Apartament) {
		adresa += `, ap. ${json.Apartament}`
	}

	return adresa
}

function copyContentToClipboard(event) {
	let targetId = event.dataset.target

	if (targetId === undefined)
		return

	target = document.getElementById(targetId)

	if (target === undefined)
		return

	setClipboard(target.innerHTML)
}

async function setClipboard(text) {
	const type = "text/plain";
	const clipboardItemData = {
		[type]: text,
	};

	const clipboardItem = new ClipboardItem(clipboardItemData);
	await navigator.clipboard.write([clipboardItem]);
}

class InfoPage extends Base {
    init() {
        super.init()

        this.inregistrare = window.location.pathname.split('/').pop();

		this.veziDosareButton = document.getElementsByClassName("view-button")[0]

        this.profileName = document.getElementById("profileName")
        this.profileCompanyName = document.getElementById("profileCompanyName")
        this.profileCompanyCui = document.getElementById("profileCompanyCui")
        this.profileCompanyId = document.getElementById("profileCompanyId")
        this.profileCompanyEuid = document.getElementById("profileCompanyEuid")
        this.profileCompanyForma = document.getElementById("profileCompanyForma")
        this.profileCompanyStatus = document.getElementById("profileCompanyStatus")
        this.profileCompanyData = document.getElementById("profileCompanyData")
        this.profileCompanyJudet = document.getElementById("profileCompanyJudet")
        this.profileCompanyLocalitate = document.getElementById("profileCompanyLocalitate")
        this.profileCompanyAdresa = document.getElementById("profileCompanyAdresa")
        this.profileCompanyCodPostal = document.getElementById("profileCompanyCodPostal")
        this.profileCompanyCaen = document.getElementById("profileCompanyCaen")
        this.profileCompanyPrimaryCaen = document.getElementById("profileCompanyPrimaryCaen")
        this.profileCompanyTva = document.getElementById("profileCompanyTva")

        this.profileCompanyAdministratori = document.getElementById("profileCompanyAdministratori")
        this.profileCompanyAsociati = document.getElementById("profileCompanyAsociati")
    }

    populatePage(data) {
        this.veziDosareButton.href = `/dosareJuridice/${this.inregistrare}/1?nume_firma=${data.Nume}`

        this.profileName.textContent = data.Nume
        this.profileCompanyName.textContent = data.Nume
        this.profileCompanyCui.textContent = data.Cui
        this.profileCompanyEuid.textContent = data.Euid
        this.profileCompanyId.textContent = data.CodInmatriculare
        this.profileCompanyForma.textContent = data.FormaJuridica
        this.profileCompanyStatus.textContent = data.Statusuri.sort((a, b) => {
            return (a === "funcțiune") - (b === "funcțiune")
        }).join(", ")

        this.profileCompanyData.textContent = formatDate(data.DataInregistrare)
        this.profileCompanyJudet.textContent = data.Judet
        this.profileCompanyLocalitate.textContent = data.Localitate

        setValueIfExist(this.profileCompanyAdresa, createAddress(data))
        setValueIfExist(this.profileCompanyCodPostal, data.CodPostal)
        setValueIfExist(this.profileCompanyCaen, data.CoduriCaen?.join(", "))

        if (data.Tva === true) {
            this.profileCompanyTva.textContent = "Da"
        }

        if (data.Reprezentanti) {
            for (let reprezentant of data.Reprezentanti) {
                let [name, role] = reprezentant.split('^')

                if (role === 'administrator') {
                    if (this.profileCompanyAdministratori.children.length === 0) {
                        this.profileCompanyAdministratori.textContent = ""
                    }

                    const link = document.createElement("a");

                    link.href = `/admins/${this.inregistrare}/${name}/1`
                    link.textContent = name

                    this.profileCompanyAdministratori.appendChild(link)
                    this.profileCompanyAdministratori.appendChild(document.createElement("br"))
                } else {
                    if (this.profileCompanyAsociati.children.length === 0) {
                        this.profileCompanyAsociati.textContent = ""
                    }

                    const link = document.createElement("a");
                    link.textContent = name
                    link.href = `/admins/${this.inregistrare}/${name}/1`

                    const span = document.createElement("span");
                    span.textContent = ` - ${role}`

                    this.profileCompanyAsociati.appendChild(link)
                    this.profileCompanyAsociati.appendChild(span)
                    this.profileCompanyAsociati.appendChild(document.createElement("br"))
                }
            }
        }

        if (!data.BilanturiFirma) {
            return
        }

        const bilanturiFirmaContainer = document.getElementById("dateFinanciareContainer")

        const table = document.getElementById("dateFinanciareTable")

        const prototype = document.getElementById("dateFinanciareRowPrototype")

        const financiarAn = prototype.getElementsByClassName("dateFinanciareAn")[0]
        const financiarCifraAfaceri = prototype.getElementsByClassName("dateFinanciareCifraAfaceri")[0]
        const financiarProfit = prototype.getElementsByClassName("dateFinanciareProfit")[0]
        const financiarDatorii = prototype.getElementsByClassName("dateFinanciareDatorii")[0]
        const financiarActiveImobilizate = prototype.getElementsByClassName("dateFinanciareActiveImobilizate")[0]
        const financiarActiveCirculante = prototype.getElementsByClassName("dateFinanciareActiveCirculante")[0]
        const financiarCapitaluriProprii = prototype.getElementsByClassName("DateFinanciareCapitaluriProprii")[0]
        const financiarCapitaluriAngajati = prototype.getElementsByClassName("DateFinanciareAngajati")[0]

        for (let [index, bilant] of data.BilanturiFirma.entries()) {
            if (index == 0) {
                profileCompanyPrimaryCaen.textContent = bilant.Caen
            }

            financiarAn.textContent = bilant.An
            financiarCifraAfaceri.textContent = formatMoney(bilant.CifraAfaceri)
            financiarProfit.textContent = formatMoney(bilant.ProfitNet)
            financiarDatorii.textContent = formatMoney(bilant.Datorii)
            financiarActiveImobilizate.textContent = formatMoney(bilant.ActiveImobilizate)
            financiarActiveCirculante.textContent = formatMoney(bilant.ActiveCirculante)
            financiarCapitaluriProprii.textContent = formatMoney(bilant.Capitaluri)
            financiarCapitaluriAngajati.textContent = bilant.Angajati

            let clone = prototype.cloneNode(true)
            clone.style.display = "table-row"
            clone.removeAttribute("id")

            prototype.before(clone)
        }

        bilanturiFirmaContainer.style.display = "block"
    }

    async Start() {
        const data = await fetch(`/firma/${this.inregistrare}`)
        if (data.status !== 200) {
            await setErrorPage(data.status)
            return
        }

        const json = await data.json()

        console.log(json)

        this.populatePage(json)
    }
}

function CreateComponent() {
    return new InfoPage()
}
