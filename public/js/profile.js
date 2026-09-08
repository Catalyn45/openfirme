function openDosareJuridice(button) {
	let inregistrare = document.getElementById("profileCompanyId").textContent
	inregistrare = inregistrare.replaceAll("/", "-")

    button.href = `/dosareJuridice/${inregistrare}/1`
}


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

function setValueIfExist(element, value) {
	if (value) {
		element.textContent = value
	}
}

async function main() {
	const profileId = window.location.pathname.split('/').pop();

	let profileName = document.getElementById("profileName")
	let profileCompanyName = document.getElementById("profileCompanyName")
	let profileCompanyCui = document.getElementById("profileCompanyCui")
	let profileCompanyId = document.getElementById("profileCompanyId")
	let profileCompanyEuid = document.getElementById("profileCompanyEuid")
	let profileCompanyForma = document.getElementById("profileCompanyForma")
	let profileCompanyStatus = document.getElementById("profileCompanyStatus")
	let profileCompanyData = document.getElementById("profileCompanyData")
	let profileCompanyJudet = document.getElementById("profileCompanyJudet")
	let profileCompanyLocalitate = document.getElementById("profileCompanyLocalitate")
	let profileCompanyAdresa = document.getElementById("profileCompanyAdresa")
	let profileCompanyCodPostal = document.getElementById("profileCompanyCodPostal")
	let profileCompanyCaen = document.getElementById("profileCompanyCaen")
	let profileCompanyPrimaryCaen = document.getElementById("profileCompanyPrimaryCaen")
	let profileCompanyTva = document.getElementById("profileCompanyTva")

	let profileCompanyAdministratori = document.getElementById("profileCompanyAdministratori")
	let profileCompanyAsociati = document.getElementById("profileCompanyAsociati")

	const data = await fetch(`/firma/${profileId}`)
	const json = await data.json()

    console.log(json)

	profileName.textContent = json.Nume
    profileCompanyName.textContent = json.Nume
	profileCompanyCui.textContent = json.Cui
	profileCompanyEuid.textContent = json.Euid
	profileCompanyId.textContent = json.CodInmatriculare
	profileCompanyForma.textContent = json.FormaJuridica
	profileCompanyStatus.textContent = json.Statusuri.sort((a, b) => {
		return (a === "funcțiune") - (b === "funcțiune")
	}).join(", ")

	profileCompanyData.textContent = formatDate(json.DataInregistrare)
	profileCompanyJudet.textContent = json.Judet
	profileCompanyLocalitate.textContent = json.Localitate

	setValueIfExist(profileCompanyAdresa, createAddress(json))
	setValueIfExist(profileCompanyCodPostal, json.CodPostal)
	setValueIfExist(profileCompanyCaen, json.CoduriCaen?.join(", "))

	if (json.Tva === true) {
		profileCompanyTva.textContent = "Da"
	}

	if (json.Reprezentanti) {
		for (let reprezentant of json.Reprezentanti) {
			let [name, role] = reprezentant.split('^')

			if (role === 'administrator') {
				if (profileCompanyAdministratori.children.length === 0) {
					profileCompanyAdministratori.textContent = ""
				}

				const link = document.createElement("a");

				link.href = `/admins/${profileId}/${name}/1`
				link.textContent = name

				profileCompanyAdministratori.appendChild(link)
				profileCompanyAdministratori.appendChild(document.createElement("br"))
			} else {
				if (profileCompanyAsociati.children.length === 0) {
					profileCompanyAsociati.textContent = ""
				}

				const span = document.createElement("span");
				span.textContent = `${name} - ${role}`

				profileCompanyAsociati.appendChild(span)
				profileCompanyAsociati.appendChild(document.createElement("br"))
			}
		}
	}

	if (!json.BilanturiFirma) {
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

	for (let [index, bilant] of json.BilanturiFirma.entries()) {
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

main().catch(console.error);
