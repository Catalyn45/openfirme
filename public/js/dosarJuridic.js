function openPortalJust(button) {
    let numarDosar = window.location.pathname.split('/').at(-1)
    numarDosar = numarDosar.replaceAll('-', '/')

    let portalJustEndpoint = `https://portal.just.ro/SitePages/cautare.aspx?k=${encodeURIComponent(numarDosar)}&v1=default`
    button.href = portalJustEndpoint
}

async function main() {
	const numarDosar = window.location.pathname.split('/').at(-1)
	const inregistrare = window.location.pathname.split('/').at(-2)

	let profileNumarDosar = document.getElementById("profileNumarDosar")
	let profileObiect = document.getElementById("profileObiect")
	let profileCategorie = document.getElementById("profileCategorie")
	let profileTribunal = document.getElementById("profileTribunal")
	let profileDepartament = document.getElementById("profileDepartament")
	let profileData = document.getElementById("profileData")
	let profileStadiuProcesual = document.getElementById("profileStadiuProcesual")

    let numePartePrototype = document.getElementById("numePartePrototype")
    let calitatePartePrototype = document.getElementById("calitatePartePrototype")
    let numePartePrototypeValue = document.getElementById("numePartePrototypeValue")
    let calitatePartePrototypeValue = document.getElementById("calitatePartePrototypeValue")

	const data = await fetch(`/dosarJuridicFirma/${inregistrare}/${numarDosar}`)
	const json = await data.json()

    console.log(json)

	profileNumarDosar.textContent = json.Numar
    profileObiect.textContent = json.Obiect
    profileCategorie.textContent = json.CategorieCazNume
    profileTribunal.textContent = json.Institutie
    profileDepartament.textContent = json.Departament
    profileData.textContent = json.Data
    profileStadiuProcesual.textContent = json.StadiuProcesualNume

    for (let parte of json.Parti.DosareParte) {
        console.log(parte)

        numePartePrototypeValue.textContent = parte.Nume
        calitatePartePrototypeValue.textContent = parte.CalitateParte

        numeClone = numePartePrototype.cloneNode(true)
        calitateClone = calitatePartePrototype.cloneNode(true)

        numeClone.style.display = "block"
        numeClone.removeAttribute("id")

        calitateClone.style.display = "block"
        calitateClone.removeAttribute("id")

        numePartePrototype.before(numeClone)
        numePartePrototype.before(calitateClone)
    }
}

main().catch(console.error);
