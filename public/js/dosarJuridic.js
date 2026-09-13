class InfoDosarJuridicPage extends BaseComponent {
    constructor() {
        super()

        this.numarDosar = window.location.pathname.split('/').at(-1)
        this.inregistrare = window.location.pathname.split('/').at(-2)

        this.profileNumarDosar = document.getElementById("profileNumarDosar")
        this.profileObiect = document.getElementById("profileObiect")
        this.profileCategorie = document.getElementById("profileCategorie")
        this.profileTribunal = document.getElementById("profileTribunal")
        this.profileDepartament = document.getElementById("profileDepartament")
        this.profileData = document.getElementById("profileData")
        this.profileStadiuProcesual = document.getElementById("profileStadiuProcesual")

        this.numePartePrototype = document.getElementById("numePartePrototype")
        this.calitatePartePrototype = document.getElementById("calitatePartePrototype")
        this.numePartePrototypeValue = document.getElementById("numePartePrototypeValue")
        this.calitatePartePrototypeValue = document.getElementById("calitatePartePrototypeValue")
        this.numePartePrototypeClipboardButton = document.getElementById("numePartePrototypeClipboardButton")
        this.calitatePartePrototypeClipboardButton = document.getElementById("calitatePartePrototypeClipboardButton")
        

        this.portalJustButton = document.getElementsByClassName("view-button")[0]
    }

    populateWithData(data) {
        let numarDosar = this.numarDosar.replaceAll('-', '/')
        let portalJustEndpoint = `https://portal.just.ro/SitePages/cautare.aspx?k=${encodeURIComponent(numarDosar)}&v1=default`
        this.portalJustButton.href = portalJustEndpoint

        this.profileNumarDosar.textContent = data.Numar
        this.profileObiect.textContent = data.Obiect
        this.profileCategorie.textContent = data.CategorieCazNume
        this.profileTribunal.textContent = data.Institutie
        this.profileDepartament.textContent = data.Departament
        this.profileData.textContent = formatDateDosare(data.Data)
        this.profileStadiuProcesual.textContent = data.StadiuProcesualNume

        let indexParte = 0
        let numeParteId = this.numePartePrototypeValue.getAttribute("id")
        let calitateParteId = this.calitatePartePrototypeValue.getAttribute("id")
        
        this.numePartePrototypeClipboardButton.removeAttribute("id")
        this.calitatePartePrototypeClipboardButton.removeAttribute("id")
        
        for (let parte of data.Parti.DosareParte) {
            console.log(parte)

            this.numePartePrototypeValue.textContent = parte.Nume
            this.numePartePrototypeValue.setAttribute("id", `${numeParteId}-${indexParte}`)

            this.calitatePartePrototypeValue.textContent = parte.CalitateParte
            this.calitatePartePrototypeValue.setAttribute("id", `${calitateParteId}-${indexParte}`)

            this.numePartePrototypeClipboardButton.dataset.target = this.numePartePrototypeValue.getAttribute("id")
            this.calitatePartePrototypeClipboardButton.dataset.target = this.calitatePartePrototypeValue.getAttribute("id")

            let numeClone = this.numePartePrototype.cloneNode(true)
            let calitateClone = this.calitatePartePrototype.cloneNode(true)

            numeClone.style.display = "block"
            numeClone.removeAttribute("id")

            calitateClone.style.display = "block"
            calitateClone.removeAttribute("id")

            this.numePartePrototype.before(numeClone)
            this.numePartePrototype.before(calitateClone)

            indexParte++
        }

        this.numePartePrototypeValue.setAttribute("id", numeParteId)
        this.calitatePartePrototypeValue.setAttribute("id", calitateParteId)
    }

    initSearchBar() { }
    setSearchBar() { }

    async Start() {
        const data = await fetch(`/dosarJuridicFirma/${this.inregistrare}/${this.numarDosar}`)
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
