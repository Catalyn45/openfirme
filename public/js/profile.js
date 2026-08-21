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

function formatMoney(amount) {
	return new Intl.NumberFormat('de-DE', {
	  style: 'currency',
	  currency: 'RON',
	  maximumFractionDigits: 0
	}).format(amount);
}

async function main() {
	const profileId = window.location.pathname.split('/').pop();

	console.log(profileId)

	let profileName = document.getElementById("profileName")
	let profileCompanyName = document.getElementById("profileCompanyName")
	let profileCompanyCui = document.getElementById("profileCompanyCui")
	let profileCompanyId = document.getElementById("profileCompanyId")
	let profileCompanyForma = document.getElementById("profileCompanyForma")
	let profileCompanyStatus = document.getElementById("profileCompanyStatus")
	let profileCompanyData = document.getElementById("profileCompanyData")
	let profileCompanyJudet = document.getElementById("profileCompanyJudet")
	let profileCompanyLocalitate = document.getElementById("profileCompanyLocalitate")
	let profileCompanyAdresa = document.getElementById("profileCompanyAdresa")
	let profileCompanyCodPostal = document.getElementById("profileCompanyCodPostal")
	let profileCompanyCaen = document.getElementById("profileCompanyCaen")
	let profileCompanyAdministratori = document.getElementById("profileCompanyAdministratori")

	const data = await fetch(`/firma/${profileId}`)
	const json = await data.json()

    console.log(json)

	profileName.textContent = json.Nume
    profileCompanyName.textContent = json.Nume
	profileCompanyCui.textContent = json.Cui
	profileCompanyId.textContent = json.CodInmatriculare
	profileCompanyForma.textContent = json.FormaJuridica
	profileCompanyStatus.textContent = json.Status
	profileCompanyData.textContent = json.DataInregistrare
	profileCompanyJudet.textContent = json.Judet
	profileCompanyLocalitate.textContent = json.Localitate
	profileCompanyAdresa.textContent = createAddress(json)
	profileCompanyCodPostal.textContent = json.CodPostal
	profileCompanyCaen.textContent = json.CoduriCaen?.join(", ") ?? ""
	profileCompanyAdministratori.textContent = json.Administratori?.join(", ") ?? "Fără administratori"

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

	for (let bilant of json.BilanturiFirma) {
		financiarAn.textContent = bilant.An
		financiarCifraAfaceri.textContent = formatMoney(bilant.CifraAfaceri)

		let profitNet = bilant.ProfitNet
		if (profitNet == 0) {
			profitNet = -1 * bilant.PierdereNeta
		}

		financiarProfit.textContent = formatMoney(profitNet)
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

main()
