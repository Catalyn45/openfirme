class InfoDosarJuridicPage extends Base {
    init() {
        super.init()

        this.numarDosar = window.location.pathname.split('/').at(-1)
        this.inregistrare = window.location.pathname.split('/').at(-2)

        this.expandContentButton = document.getElementsByClassName("expand-content-button")[0]

        this.profileNumarDosar = document.getElementById("profileNumarDosar")
        this.profileObiect = document.getElementById("profileObiect")
        this.profileCategorie = document.getElementById("profileCategorie")
        this.profileTribunal = document.getElementById("profileTribunal")
        this.profileDepartament = document.getElementById("profileDepartament")
        this.profileData = document.getElementById("profileData")
        this.profileStadiuProcesual = document.getElementById("profileStadiuProcesual")

        this.partePrototype = document.getElementById("partePrototype")
        this.partePrototypeValue = this.partePrototype.getElementsByClassName("partePrototypeValue")[0]

        this.sedintaPrototype = document.getElementById("sedintaPrototype")
        this.sedintaTitle = this.sedintaPrototype.getElementsByClassName("sedinta-title")[0]
        this.sedintaData = this.sedintaPrototype.getElementsByClassName("sedinta-data")[0]
        this.sedintaComplet = this.sedintaPrototype.getElementsByClassName("sedinta-complet")[0]
        this.sedintaDataPronuntare = this.sedintaPrototype.getElementsByClassName("sedinta-data-pronuntare")[0]
        this.sedintaOra = this.sedintaPrototype.getElementsByClassName("sedinta-ora")[0]
        this.sedintaSolutie = this.sedintaPrototype.getElementsByClassName("sedinta-solutie")[0]
        this.sedintaSolutieSumar = this.sedintaPrototype.getElementsByClassName("sedinta-solutie-sumar")[0]

        this.portalJustButton = document.getElementsByClassName("view-button")[0]
    }

    expandSedinta(element) {
        let arrow = element.parentNode.parentNode.getElementsByClassName("sedintaExpand")[0]
        let content = element.parentNode.parentNode.getElementsByClassName("profile-grid")[0]

        let display = content.style.display

        if (!display) {
            content.style.display = "none"
            arrow.textContent = "▸"
        } else {
            content.style.removeProperty("display")
            arrow.textContent = "▾"
        }
    }

    populateWithData(data) {
        let numarDosar = this.numarDosar.replaceAll('-', '/')
        let portalJustEndpoint = `https://portal.just.ro/SitePages/cautare.aspx?k=${encodeURIComponent(numarDosar)}&v1=default`
        this.portalJustButton.href = portalJustEndpoint

        this.profileNumarDosar.textContent = data.Numar
        this.profileObiect.textContent = data.Obiect
        this.profileCategorie.textContent = data.CategorieCazNume
        this.profileTribunal.textContent = normalizeJudecatorie(data.Institutie)
        this.profileDepartament.textContent = data.Departament
        this.profileData.textContent = formatDateDosare(data.Data)
        this.profileStadiuProcesual.textContent = data.StadiuProcesualNume

        if (data.Parti.DosareParte.length > 4) {
            this.expandContentButton.style.display = "block"
        }

        for (let parte of data.Parti.DosareParte) {
            console.log(parte)

            this.partePrototypeValue.textContent = `${parte.Nume} - ${parte.CalitateParte}`

            let clone = this.partePrototype.cloneNode(true)

            clone.style.removeProperty("display")
            clone.removeAttribute("id")

            this.partePrototype.before(clone)
        }

        for (let sedinta of data.Sedinte?.DosareSedinta ?? []) {
            this.sedintaTitle.textContent = `Ședință - ${formatDateSedinta(sedinta.Data)}`
            this.sedintaData.textContent = formatDateSedinta(sedinta.Data)
            this.sedintaComplet.textContent = sedinta.Complet
            this.sedintaDataPronuntare.textContent = formatDateDosare(sedinta.DataPronuntare)
            this.sedintaOra.textContent = sedinta.Ora
            this.sedintaSolutie.textContent = sedinta.Solutie
            this.sedintaSolutieSumar.textContent = sedinta.SolutieSumar

            let clone = this.sedintaPrototype.cloneNode(true)

            clone.style.removeProperty("display")
            clone.removeAttribute("id")

            this.sedintaPrototype.before(clone)
        }
    }

    async Start() {
        const data = await fetch(`/api/dosarJuridic/${this.inregistrare}/${this.numarDosar}`)
        if (data.status !== 200) {
            await setErrorPage(data.status)
            return
        }

        const json = await data.json()

        console.log(json)

        this.populateWithData(json)
    }
}

function CreateComponent() {
    return new InfoDosarJuridicPage()
}
