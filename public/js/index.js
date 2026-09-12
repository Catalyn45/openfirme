class BaseComponent {
    constructor() {
        this.judetSearchParam = document.getElementById("judetSearchParam")
        this.statusSearchParam = document.getElementById("statusSearchParam")
        this.formaJuridicaSearchParam = document.getElementById("formaJuridicaSearchParam")

        this.initFilters()
    }

    initFilters() {
        this.query = document.getElementById("searchInput")
        this.county = document.getElementById("countyFilter")
        this.status = document.getElementById("statusFilter")
        this.formaJuridica = document.getElementById("formaFilter")
    }

	getFilters() {
		const params = new URLSearchParams();

		let nume_partial = this.query?.value?.toLowerCase().trim()
		if (nume_partial) {
			params.set('nume_partial', nume_partial)
		}

        params.set('judet', this.county.value)
        params.set('status', this.status.value)
        params.set('forma_juridica', this.formaJuridica.value)

		return params
	}

	setFilters() {
		const params = new URLSearchParams(window.location.search);
		
		if (this.query) {
			this.query.value = params.get("nume_partial") ?? this.query.value
		}

        this.county.value = params.get("judet") ?? this.county.value
        this.status.value = params.get("status") ?? this.status.value
        this.formaJuridica.value = params.get("forma_juridica") ?? this.formaJuridica.value

        return params
	}

    setSearchAction() {
        // TODO: find a better way
        if (this.judetSearchParam) {
            let params = this.getFilters()
            this.judetSearchParam.value = params.get("judet")
            this.statusSearchParam.value = params.get("status")
            this.formaJuridicaSearchParam.value = params.get("forma_juridica")
        }
    }

	onChangeFilter() {
        this.setSearchAction()
    }

    Start() {
        this.setFilters()
        this.setSearchAction()
    }
}

function CreateComponent() {
    return new BaseComponent()
}
