
async function main() {
	const profileId = window.location.href.split('/').pop();

	let profileName = document.getElementById("profileName")
	let profilecCompanyName = document.getElementById("profileCompanyName")
	
	const data = await fetch(`/firme/${profileId}`)
	const json = await data.json()

    console.log(json)

	profileName.textContent = json.Nume
    profileCompanyName.textContent = json.Nume
}

main()
