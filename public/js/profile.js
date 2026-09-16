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

        data.Statusuri = data.Statusuri.sort((a, b) => {
            return (a === "funcțiune") - (b === "funcțiune")
        })

        this.profileCompanyStatus.textContent = ""
        for (let status of data.Statusuri) {
            const p = document.createElement("p")
            p.textContent = `- ${status}`

            this.profileCompanyStatus.appendChild(p)
        }

        this.profileCompanyData.textContent = formatDate(data.DataInregistrare)
        this.profileCompanyJudet.textContent = data.Judet
        this.profileCompanyLocalitate.textContent = data.Localitate

        setValueIfExist(this.profileCompanyAdresa, createAddress(data))
        setValueIfExist(this.profileCompanyCodPostal, data.CodPostal)

        let caenDescriere = {}
        for (let codCaen of data.CoduriCaen ?? []) {
            let [cod_caen, descriere_caen] = codCaen.split('^')
            caenDescriere[cod_caen] = descriere_caen
            this.profileCompanyCaen.textContent = ""
        }

        for (let [codCaen, descriereCaen] of Object.entries(caenDescriere)) {
            const p = document.createElement("p");
            p.textContent = `${codCaen} - ${descriereCaen}`
            this.profileCompanyCaen.appendChild(p)

        }

        if (data.Tva === true) {
            this.profileCompanyTva.textContent = "Da"
        }

        if (data.Reprezentanti) {
            for (let reprezentant of data.Reprezentanti) {
                let [name, role] = reprezentant.split('^')

                const p = document.createElement("p")
                const link = document.createElement("a");

                if (role === 'administrator') {
                    if (this.profileCompanyAdministratori.children.length === 0) {
                        this.profileCompanyAdministratori.textContent = ""
                    }

                    link.href = `/admins/${this.inregistrare}/${name}/1`
                    link.textContent = name

                    p.appendChild(link)

                    this.profileCompanyAdministratori.appendChild(p)
                } else {
                    if (this.profileCompanyAsociati.children.length === 0) {
                        this.profileCompanyAsociati.textContent = ""
                    }

                    link.textContent = name
                    link.href = `/admins/${this.inregistrare}/${name}/1`

                    const span = document.createElement("span");
                    span.textContent = ` - ${role}`

                    p.appendChild(link)
                    p.appendChild(span)

                    this.profileCompanyAsociati.appendChild(p)
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

                let descriere = caenDescriere[bilant.Caen]
                if (descriere) {
                    profileCompanyPrimaryCaen.textContent += ` - ${caenDescriere[bilant.Caen]}`
                }
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
        const data = await fetch(`/api/profile/${this.inregistrare}`)
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
