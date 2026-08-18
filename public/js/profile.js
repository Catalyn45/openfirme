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

async function main() {
	const profileId = window.location.href.split('/').pop();

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

	const data = await fetch(`/firme/${profileId}`)
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
}

main()
